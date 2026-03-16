package retrieval

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"memorycli/internal/governance"
	"memorycli/internal/index"
	"memorycli/internal/types"
)

type candidate struct {
	Entry     types.IndexEntry
	Meta      types.MetadataFile
	Body      string
	Haystack  string
	Score     float64
	Lexical   float64
	Embedding float64
	Fused     float64
	Reason    string
	Freshness float64
	HasVector bool
}

type RetrieveOptions struct {
	// Mode controls retrieval strategy:
	// - "classic": legacy behavior (semantic confidence gate + deterministic fallback)
	// - "hybrid": lexical + embedding reciprocal-rank fusion with top-k traces
	Mode string
	// TopK controls how many candidate traces are returned in RetrieveResult.
	TopK int
	// RRFK controls reciprocal-rank fusion smoothing constant.
	RRFK float64
}

type cachedCandidates struct {
	candidates []candidate
	skipped    int
}

type cachedEmbeddings struct {
	embeddings map[string]types.EmbeddingRecord
}

type cachedQueryEmbedding struct {
	vector  []float64
	expires time.Time
}

var retrievalCacheState = struct {
	mu              sync.RWMutex
	candidateByKey  map[string]cachedCandidates
	embeddingByKey  map[string]cachedEmbeddings
	queryEmbeddings map[string]cachedQueryEmbedding
}{
	candidateByKey:  map[string]cachedCandidates{},
	embeddingByKey:  map[string]cachedEmbeddings{},
	queryEmbeddings: map[string]cachedQueryEmbedding{},
}

const (
	maxCachedQueryEmbeddings = 256
	queryEmbeddingTTL        = 10 * time.Minute
)

func clearRetrievalCaches() {
	retrievalCacheState.mu.Lock()
	defer retrievalCacheState.mu.Unlock()
	retrievalCacheState.candidateByKey = map[string]cachedCandidates{}
	retrievalCacheState.embeddingByKey = map[string]cachedEmbeddings{}
	retrievalCacheState.queryEmbeddings = map[string]cachedQueryEmbedding{}
}

func normalizeRetrieveOptions(options RetrieveOptions) RetrieveOptions {
	mode := strings.ToLower(strings.TrimSpace(options.Mode))
	if mode == "" {
		mode = "classic"
	}
	topK := options.TopK
	if topK <= 0 {
		topK = 5
	}
	if topK > 50 {
		topK = 50
	}
	rrfK := options.RRFK
	if rrfK <= 0 {
		rrfK = 60.0
	}
	return RetrieveOptions{
		Mode: mode,
		TopK: topK,
		RRFK: rrfK,
	}
}

func Retrieve(root, query, domain string) (types.RetrieveResult, error) {
	result, _, err := RetrieveWithOptionsAndEndpointAndSession(
		root,
		query,
		domain,
		DefaultEmbeddingEndpoint,
		"",
		RetrieveOptions{},
	)
	return result, err
}

func RetrieveWithEmbeddingEndpoint(root, query, domain, embeddingEndpoint string) (types.RetrieveResult, string, error) {
	return RetrieveWithOptionsAndEndpointAndSession(
		root,
		query,
		domain,
		embeddingEndpoint,
		"",
		RetrieveOptions{},
	)
}

func RetrieveWithEmbeddingEndpointAndSession(root, query, domain, embeddingEndpoint, sessionID string) (types.RetrieveResult, string, error) {
	return RetrieveWithOptionsAndEndpointAndSession(
		root,
		query,
		domain,
		embeddingEndpoint,
		sessionID,
		RetrieveOptions{},
	)
}

func RetrieveWithOptionsAndEndpointAndSession(
	root,
	query,
	domain,
	embeddingEndpoint,
	sessionID string,
	options RetrieveOptions,
) (types.RetrieveResult, string, error) {
	options = normalizeRetrieveOptions(options)
	startedAt := time.Now()
	if strings.TrimSpace(query) == "" {
		return types.RetrieveResult{}, "", errors.New("--query is required")
	}

	idx, err := index.LoadIndex(root)
	if err != nil {
		return types.RetrieveResult{}, "", err
	}
	if len(idx.Entries) == 0 {
		return types.RetrieveResult{}, "", errors.New("memory index has no entries")
	}

	candidates, skippedCandidates, err := loadCandidatesCached(root, idx, domain)
	if err != nil {
		return types.RetrieveResult{}, "", err
	}
	if len(candidates) == 0 {
		if skippedCandidates > 0 {
			return types.RetrieveResult{}, "", fmt.Errorf("no candidates found for query/domain; skipped %d invalid entries", skippedCandidates)
		}
		return types.RetrieveResult{}, "", errors.New("no candidates found for query/domain")
	}

	q := strings.ToLower(strings.TrimSpace(query))
	warnings := []string{}
	if skippedCandidates > 0 {
		warnings = append(warnings, fmt.Sprintf("skipped %d invalid entries during candidate load", skippedCandidates))
	}
	embeddingScoresApplied := false
	queryEmbedding, embedErr := getQueryEmbeddingCached(embeddingEndpoint, q)
	if embedErr != nil {
		warnings = append(warnings, fmt.Sprintf("embedding unavailable; using token-overlap scoring: %v", embedErr))
	}
	profile := ActiveEmbeddingProfile(embeddingEndpoint)
	embeddings, embLoadErr := loadEmbeddingsCached(root, idx, candidates)
	if embLoadErr != nil {
		warnings = append(warnings, fmt.Sprintf("embedding store unavailable; using token-overlap scoring: %v", embLoadErr))
	}

	for i := range candidates {
		candidates[i].Lexical = semanticScore(q, candidates[i])
		candidates[i].Score = candidates[i].Lexical
		candidates[i].Reason = "token_overlap"
		if len(queryEmbedding) > 0 {
			if rec, ok := embeddings[candidates[i].Entry.ID]; ok && len(rec.Vector) > 0 {
				if isEmbeddingCompatible(profile, len(queryEmbedding), rec) {
					freshness := embeddingFreshnessBonus(rec, sessionID)
					candidates[i].Freshness = freshness
					candidates[i].HasVector = true
					candidates[i].Embedding = cosineSimilarity(queryEmbedding, rec.Vector) + freshness
					if options.Mode == "classic" {
						candidates[i].Score = candidates[i].Embedding
					}
					candidates[i].Reason = "embedding_similarity_with_freshness_bonus"
					embeddingScoresApplied = true
				}
			}
		}
	}
	if len(queryEmbedding) > 0 && !embeddingScoresApplied {
		warnings = append(warnings, "embedding unavailable for candidate entries; using token-overlap scoring")
	}

	if options.Mode == "hybrid" {
		hybrid := append([]candidate(nil), candidates...)
		assignHybridFusedScores(hybrid, options.RRFK)
		sortByHybridScore(hybrid)

		// When all signal scores are zero, retain legacy deterministic fallback behavior.
		top := hybrid[0]
		if top.Lexical <= 0 && top.Embedding <= 0 {
			for i := range hybrid {
				hybrid[i].Score = hybrid[i].Lexical
			}
			chosen := chooseDeterministicFallback(hybrid)
			result := types.RetrieveResult{
				SelectedID:    chosen.Entry.ID,
				SelectionMode: "fallback_path_priority",
				SourcePath:    chosen.Entry.Path,
				Confidence:    chosen.Score,
				Reason:        "hybrid found no semantic signal; deterministic fallback used",
				FallbackUsed:  true,
				SemanticHit:   false,
				PrecisionHint: 0,
				Candidates:    toRetrieveCandidates(hybrid, "hybrid_rrf", options.TopK),
			}
			return result, joinWarnings(warnings), nil
		}

		semanticHit := top.Lexical > 0 || top.Embedding > 0
		result := types.RetrieveResult{
			SelectedID:    top.Entry.ID,
			SelectionMode: "hybrid_rrf",
			SourcePath:    top.Entry.Path,
			Confidence:    top.Fused,
			Reason:        "hybrid reciprocal-rank fusion selected top candidate",
			FallbackUsed:  false,
			SemanticHit:   semanticHit,
			PrecisionHint: map[bool]float64{true: 1, false: 0}[semanticHit],
			Candidates:    toRetrieveCandidates(hybrid, "hybrid_rrf", options.TopK),
		}
		return result, joinWarnings(warnings), nil
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Score == candidates[j].Score {
			if candidates[i].Freshness == candidates[j].Freshness {
				return candidates[i].Entry.ID < candidates[j].Entry.ID
			}
			return candidates[i].Freshness > candidates[j].Freshness
		}
		return candidates[i].Score > candidates[j].Score
	})

	top := candidates[0]
	second := 0.0
	if len(candidates) > 1 {
		second = candidates[1].Score
	}

	confident := IsSemanticConfident(top.Score, second)
	if embeddingScoresApplied {
		confident = IsEmbeddingSemanticConfident(top.Score, second)
	}

	if governance.IsLatencyDegraded(time.Since(startedAt).Milliseconds()) {
		chosen := chooseDeterministicFallback(candidates)
		return types.RetrieveResult{
			SelectedID:    chosen.Entry.ID,
			SelectionMode: "fallback_path_priority",
			SourcePath:    chosen.Entry.Path,
			Confidence:    chosen.Score,
			Reason:        "latency degradation policy forced deterministic fallback",
			FallbackUsed:  true,
			SemanticHit:   false,
			PrecisionHint: 0,
			Candidates:    toRetrieveCandidates(candidates, "classic", options.TopK),
		}, joinWarnings(warnings), nil
	}

	if confident {
		mode := "semantic"
		if embeddingScoresApplied {
			mode = "embedding_semantic"
		}
		return types.RetrieveResult{
			SelectedID:    top.Entry.ID,
			SelectionMode: mode,
			SourcePath:    top.Entry.Path,
			Confidence:    top.Score,
			Reason:        "semantic confidence gate passed",
			FallbackUsed:  false,
			SemanticHit:   true,
			PrecisionHint: 1,
			Candidates:    toRetrieveCandidates(candidates, mode, options.TopK),
		}, joinWarnings(warnings), nil
	}

	for _, c := range candidates {
		if strings.EqualFold(c.Entry.ID, q) {
			return types.RetrieveResult{
				SelectedID:    c.Entry.ID,
				SelectionMode: "fallback_exact_key",
				SourcePath:    c.Entry.Path,
				Confidence:    c.Score,
				Reason:        "semantic confidence gate failed; exact-key fallback matched",
				FallbackUsed:  true,
				SemanticHit:   false,
				PrecisionHint: 0,
				Candidates:    toRetrieveCandidates(candidates, "classic", options.TopK),
			}, joinWarnings(warnings), nil
		}
	}

	chosen := chooseDeterministicFallback(candidates)
	return types.RetrieveResult{
		SelectedID:    chosen.Entry.ID,
		SelectionMode: "fallback_path_priority",
		SourcePath:    chosen.Entry.Path,
		Confidence:    chosen.Score,
		Reason:        "semantic confidence gate failed; deterministic path-priority fallback used",
		FallbackUsed:  true,
		SemanticHit:   false,
		PrecisionHint: 0,
		Candidates:    toRetrieveCandidates(candidates, "classic", options.TopK),
	}, joinWarnings(warnings), nil
}

func assignHybridFusedScores(candidates []candidate, k float64) {
	lexical := append([]candidate(nil), candidates...)
	sort.SliceStable(lexical, func(i, j int) bool {
		if lexical[i].Lexical == lexical[j].Lexical {
			return lexical[i].Entry.Path < lexical[j].Entry.Path
		}
		return lexical[i].Lexical > lexical[j].Lexical
	})
	lexicalRank := make(map[string]int, len(lexical))
	for i, c := range lexical {
		lexicalRank[c.Entry.ID] = i + 1
	}

	embedding := make([]candidate, 0, len(candidates))
	for _, c := range candidates {
		if c.HasVector {
			embedding = append(embedding, c)
		}
	}
	sort.SliceStable(embedding, func(i, j int) bool {
		if embedding[i].Embedding == embedding[j].Embedding {
			return embedding[i].Entry.Path < embedding[j].Entry.Path
		}
		return embedding[i].Embedding > embedding[j].Embedding
	})
	embeddingRank := make(map[string]int, len(embedding))
	for i, c := range embedding {
		embeddingRank[c.Entry.ID] = i + 1
	}

	for i := range candidates {
		fused := 0.0
		if rank, ok := lexicalRank[candidates[i].Entry.ID]; ok {
			fused += 1.0 / (k + float64(rank))
		}
		if rank, ok := embeddingRank[candidates[i].Entry.ID]; ok {
			fused += 1.0 / (k + float64(rank))
		}
		candidates[i].Fused = fused
	}
}

func sortByHybridScore(candidates []candidate) {
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Fused == candidates[j].Fused {
			if candidates[i].Embedding == candidates[j].Embedding {
				if candidates[i].Lexical == candidates[j].Lexical {
					return candidates[i].Entry.Path < candidates[j].Entry.Path
				}
				return candidates[i].Lexical > candidates[j].Lexical
			}
			return candidates[i].Embedding > candidates[j].Embedding
		}
		return candidates[i].Fused > candidates[j].Fused
	})
}

func toRetrieveCandidates(candidates []candidate, mode string, topK int) []types.RetrieveCandidate {
	if topK <= 0 {
		return nil
	}
	limit := topK
	if len(candidates) < limit {
		limit = len(candidates)
	}
	out := make([]types.RetrieveCandidate, 0, limit)
	for i := 0; i < limit; i++ {
		c := candidates[i]
		conf := c.Score
		if mode == "hybrid_rrf" {
			conf = c.Fused
		}
		out = append(out, types.RetrieveCandidate{
			ID:             c.Entry.ID,
			SourcePath:     c.Entry.Path,
			SelectionMode:  mode,
			Confidence:     conf,
			LexicalScore:   c.Lexical,
			EmbeddingScore: c.Embedding,
			FusedScore:     c.Fused,
			HasVector:      c.HasVector,
			Reason:         c.Reason,
		})
	}
	return out
}

func getQueryEmbeddingCached(embeddingEndpoint, query string) ([]float64, error) {
	cacheKey := strings.TrimSpace(embeddingEndpoint) + "|" + query
	now := time.Now()
	retrievalCacheState.mu.RLock()
	if cached, ok := retrievalCacheState.queryEmbeddings[cacheKey]; ok && now.Before(cached.expires) {
		retrievalCacheState.mu.RUnlock()
		return append([]float64(nil), cached.vector...), nil
	}
	retrievalCacheState.mu.RUnlock()

	vec, err := GenerateEmbedding(embeddingEndpoint, query)
	if err != nil {
		return nil, err
	}

	retrievalCacheState.mu.Lock()
	if len(retrievalCacheState.queryEmbeddings) >= maxCachedQueryEmbeddings {
		// Simple bounded cache eviction by removing expired entries first, then one arbitrary key.
		for key, record := range retrievalCacheState.queryEmbeddings {
			if now.After(record.expires) {
				delete(retrievalCacheState.queryEmbeddings, key)
			}
		}
		if len(retrievalCacheState.queryEmbeddings) >= maxCachedQueryEmbeddings {
			for key := range retrievalCacheState.queryEmbeddings {
				delete(retrievalCacheState.queryEmbeddings, key)
				break
			}
		}
	}
	retrievalCacheState.queryEmbeddings[cacheKey] = cachedQueryEmbedding{
		vector:  append([]float64(nil), vec...),
		expires: now.Add(queryEmbeddingTTL),
	}
	retrievalCacheState.mu.Unlock()
	return vec, nil
}

func loadEmbeddingsCached(
	root string,
	idx types.IndexFile,
	candidates []candidate,
) (map[string]types.EmbeddingRecord, error) {
	cacheKey := embeddingCacheKey(root, idx.UpdatedAt, candidates)
	retrievalCacheState.mu.RLock()
	if cached, ok := retrievalCacheState.embeddingByKey[cacheKey]; ok {
		retrievalCacheState.mu.RUnlock()
		return cached.embeddings, nil
	}
	retrievalCacheState.mu.RUnlock()

	ids := make([]string, 0, len(candidates))
	for _, c := range candidates {
		ids = append(ids, c.Entry.ID)
	}
	records, err := index.GetEmbeddingRecords(root, ids)
	if err != nil {
		return nil, err
	}
	retrievalCacheState.mu.Lock()
	retrievalCacheState.embeddingByKey[cacheKey] = cachedEmbeddings{embeddings: records}
	retrievalCacheState.mu.Unlock()
	return records, nil
}

func embeddingCacheKey(root, updatedAt string, candidates []candidate) string {
	ids := make([]string, 0, len(candidates))
	for _, c := range candidates {
		ids = append(ids, c.Entry.ID)
	}
	sort.Strings(ids)
	return fmt.Sprintf("%s|%s|%s", root, updatedAt, strings.Join(ids, ","))
}

func isEmbeddingCompatible(profile EmbeddingProfile, queryDim int, rec types.EmbeddingRecord) bool {
	if queryDim <= 0 || len(rec.Vector) == 0 {
		return false
	}
	if rec.Dim > 0 && rec.Dim != queryDim {
		return false
	}
	if len(rec.Vector) != queryDim {
		return false
	}
	if strings.TrimSpace(profile.ModelID) != "" && strings.TrimSpace(rec.ModelID) != "" && profile.ModelID != rec.ModelID {
		return false
	}
	if strings.TrimSpace(profile.Provider) != "" && strings.TrimSpace(rec.Provider) != "" && profile.Provider != rec.Provider {
		return false
	}
	return true
}

func embeddingFreshnessBonus(rec types.EmbeddingRecord, sessionID string) float64 {
	bonus := 0.0
	if strings.TrimSpace(sessionID) != "" && strings.TrimSpace(rec.SessionID) == strings.TrimSpace(sessionID) {
		bonus += 0.02
	}
	when, err := time.Parse(time.RFC3339, strings.TrimSpace(rec.GeneratedAt))
	if err != nil {
		when, err = time.Parse(time.RFC3339, strings.TrimSpace(rec.LastUpdated))
	}
	if err != nil {
		return bonus
	}
	age := time.Since(when)
	if age <= 24*time.Hour {
		bonus += 0.01
	} else if age <= 7*24*time.Hour {
		bonus += 0.005
	}
	return bonus
}

func loadCandidatesCached(root string, idx types.IndexFile, domain string) ([]candidate, int, error) {
	cacheKey := fmt.Sprintf("%s|%s|%s|%d", root, domain, idx.UpdatedAt, len(idx.Entries))
	retrievalCacheState.mu.RLock()
	if cached, ok := retrievalCacheState.candidateByKey[cacheKey]; ok {
		retrievalCacheState.mu.RUnlock()
		return cloneCandidateBase(cached.candidates), cached.skipped, nil
	}
	retrievalCacheState.mu.RUnlock()

	base, skipped, err := loadCandidatesBase(root, idx.Entries, domain)
	if err != nil {
		return nil, 0, err
	}
	retrievalCacheState.mu.Lock()
	retrievalCacheState.candidateByKey[cacheKey] = cachedCandidates{candidates: base, skipped: skipped}
	retrievalCacheState.mu.Unlock()
	return cloneCandidateBase(base), skipped, nil
}

func loadCandidatesBase(root string, entries []types.IndexEntry, domain string) ([]candidate, int, error) {
	filtered := make([]types.IndexEntry, 0, len(entries))
	for _, e := range entries {
		if domain != "" && e.Domain != domain {
			continue
		}
		if e.Status != "approved" {
			continue
		}
		filtered = append(filtered, e)
	}
	if len(filtered) == 0 {
		return []candidate{}, 0, nil
	}

	workerCount := runtime.GOMAXPROCS(0)
	if workerCount < 2 {
		workerCount = 2
	}
	type loadResult struct {
		index int
		item  candidate
		skip  bool
	}
	jobs := make(chan int, len(filtered))
	results := make(chan loadResult, len(filtered))

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idxPos := range jobs {
				e := filtered[idxPos]
				data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(e.Path)))
				if err != nil {
					results <- loadResult{index: idxPos, skip: true}
					continue
				}
				meta, err := loadMetadata(root, e)
				if err != nil {
					results <- loadResult{index: idxPos, skip: true}
					continue
				}
				body := string(data)
				results <- loadResult{
					index: idxPos,
					item: candidate{
						Entry:    e,
						Meta:     meta,
						Body:     body,
						Haystack: buildCandidateHaystack(e, meta, body),
					},
				}
			}
		}()
	}

	for i := range filtered {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	close(results)

	ordered := make([]candidate, 0, len(filtered))
	skipped := 0
	orderedMap := make(map[int]candidate, len(filtered))
	for result := range results {
		if result.skip {
			skipped++
			continue
		}
		orderedMap[result.index] = result.item
	}
	for i := 0; i < len(filtered); i++ {
		if c, ok := orderedMap[i]; ok {
			ordered = append(ordered, c)
		}
	}
	return ordered, skipped, nil
}

func buildCandidateHaystack(entry types.IndexEntry, meta types.MetadataFile, body string) string {
	return strings.ToLower(strings.Join([]string{entry.Title, entry.ID, entry.Domain, meta.Title, body}, " "))
}

func cloneCandidateBase(base []candidate) []candidate {
	out := make([]candidate, len(base))
	for i, c := range base {
		out[i] = c
		out[i].Score = 0
		out[i].Lexical = 0
		out[i].Embedding = 0
		out[i].Fused = 0
		out[i].Reason = ""
		out[i].Freshness = 0
		out[i].HasVector = false
	}
	return out
}

func joinWarnings(warnings []string) string {
	unique := map[string]struct{}{}
	out := make([]string, 0, len(warnings))
	for _, warning := range warnings {
		trimmed := strings.TrimSpace(warning)
		if trimmed == "" {
			continue
		}
		if _, exists := unique[trimmed]; exists {
			continue
		}
		unique[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return strings.Join(out, "; ")
}

func loadMetadata(root string, entry types.IndexEntry) (types.MetadataFile, error) {
	path := filepath.Join(root, filepath.FromSlash(entry.MetadataPath))
	data, err := os.ReadFile(path)
	if err != nil {
		return types.MetadataFile{}, err
	}
	var meta types.MetadataFile
	if err := json.Unmarshal(data, &meta); err != nil {
		return types.MetadataFile{}, fmt.Errorf("ERR_SCHEMA_VALIDATION: cannot parse metadata %s: %w", path, err)
	}
	if strings.TrimSpace(meta.SchemaVersion) == "" {
		return types.MetadataFile{}, errors.New("ERR_SCHEMA_VERSION_INVALID: metadata schema_version is required")
	}
	if err := index.ValidateSchemaVersion(meta.SchemaVersion); err != nil {
		return types.MetadataFile{}, err
	}
	if err := validateMetadata(meta, entry, path); err != nil {
		return types.MetadataFile{}, err
	}
	return meta, nil
}

func validateMetadata(meta types.MetadataFile, entry types.IndexEntry, path string) error {
	if strings.TrimSpace(meta.ID) == "" || strings.TrimSpace(meta.Title) == "" || strings.TrimSpace(meta.Status) == "" || strings.TrimSpace(meta.UpdatedAt) == "" {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s missing required fields", path)
	}
	if meta.ID != entry.ID {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata id %s does not match entry id %s", meta.ID, entry.ID)
	}
	if !index.IsValidStatus(meta.Status) {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s has invalid status", path)
	}
	if _, err := time.Parse(time.RFC3339, meta.UpdatedAt); err != nil {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s has invalid updated_at", path)
	}
	if strings.TrimSpace(meta.Review.ReviewedBy) == "" {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s missing review.reviewed_by", path)
	}
	if meta.Review.Decision != "approved" && meta.Review.Decision != "rejected" {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s review.decision must be approved or rejected", path)
	}
	if meta.Review.Decision == "approved" && strings.TrimSpace(meta.Review.ReviewedAt) == "" {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s approved records must include review.reviewed_at", path)
	}
	if strings.TrimSpace(meta.Review.ReviewedAt) != "" {
		if _, err := time.Parse(time.RFC3339, meta.Review.ReviewedAt); err != nil {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s has invalid review.reviewed_at", path)
		}
	}
	return nil
}

func semanticScore(query string, c candidate) float64 {
	qTokens := tokenSet(query)
	if len(qTokens) == 0 {
		return 0
	}

	hits := 0
	for tok := range qTokens {
		if strings.Contains(c.Haystack, tok) {
			hits++
		}
	}
	return float64(hits) / float64(len(qTokens))
}

func tokenSet(s string) map[string]struct{} {
	clean := strings.NewReplacer(".", " ", ",", " ", ":", " ", ";", " ", "/", " ", "-", " ", "_", " ").Replace(strings.ToLower(s))
	parts := strings.Fields(clean)
	out := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		if len(p) > 1 {
			out[p] = struct{}{}
		}
	}
	return out
}

func IsSemanticConfident(top, second float64) bool {
	const minConfidence = 0.34
	const minMargin = 0.15
	if top < minConfidence {
		return false
	}
	if top-second < minMargin {
		return false
	}
	return true
}

func IsEmbeddingSemanticConfident(top, second float64) bool {
	// Embedding cosine scores are typically lower and closer together than token-overlap scores.
	const minConfidence = 0.20
	const minMargin = 0.02
	if top < minConfidence {
		return false
	}
	if top-second < minMargin {
		return false
	}
	return true
}

func chooseDeterministicFallback(candidates []candidate) candidate {
	ordered := append([]candidate(nil), candidates...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].Score == ordered[j].Score {
			if ordered[i].Freshness == ordered[j].Freshness {
				return ordered[i].Entry.Path < ordered[j].Entry.Path
			}
			return ordered[i].Freshness > ordered[j].Freshness
		}
		return ordered[i].Score > ordered[j].Score
	})
	return ordered[0]
}
