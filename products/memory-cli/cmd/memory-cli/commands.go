package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"memorycli/internal/governance"
	"memorycli/internal/index"
	"memorycli/internal/retrieval"
	"memorycli/internal/snapshot"
	"memorycli/internal/telemetry"
	"memorycli/internal/types"
	"go.opentelemetry.io/otel/attribute"
)

func runWrite(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "write")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("write", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	id := fs.String("id", "", "entry id")
	title := fs.String("title", "", "entry title")
	typeValue := fs.String("type", "", "entry type: prompt|instruction|note")
	domain := fs.String("domain", "", "entry domain")
	body := fs.String("body", "", "entry body")
	bodyFile := fs.String("body-file", "", "path to markdown body")
	sourceRef := fs.String("source-ref", "", "optional provenance source path or reference")
	sourceKind := fs.String("source-kind", "", "optional provenance kind, for example obsidian-note")
	sourceType := fs.String("source-type", "", "optional source note type, for example decision or artifact")
	stage := fs.String("stage", "", "workflow stage: planning|architect|pm")
	sessionID := fs.String("session-id", "session-local", "telemetry session identifier")
	scenarioID := fs.String("scenario-id", "scenario-manual", "telemetry scenario identifier")
	memoryType := fs.String("memory-type", "semantic", "telemetry memory type: procedural|state|semantic")
	operatorVerdict := fs.String("operator-verdict", "not_scored", "telemetry operator verdict")
	telemetryFile := fs.String("telemetry-file", "", "optional telemetry output file (default: <root>/telemetry/events.jsonl)")
	reviewer := fs.String("reviewer", "", "reviewer identity")
	approved := fs.Bool("approved", false, "legacy flag: equivalent to --decision=approved")
	decision := fs.String("decision", "", "review decision: approved|rejected")
	notes := fs.String("notes", "", "decision notes")
	reason := fs.String("reason", "", "reason for change")
	risk := fs.String("risk", "", "risk and mitigation note")
	reworkNotes := fs.String("rework-notes", "", "required when --decision=rejected")
	reReviewedBy := fs.String("re-reviewed-by", "", "required when --decision=rejected")
	embeddingEndpoint := fs.String("embedding-endpoint", retrieval.DefaultEmbeddingEndpoint, "embedding service endpoint for write-time indexing")

	if err := fs.Parse(args); err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("memory.id", *id),
		attribute.String("memory.domain", *domain),
		attribute.String("memory.type", *typeValue),
		attribute.String("session.id", *sessionID),
	)
	startedAt := time.Now().UTC()
	traceID := fmt.Sprintf("trace-%d", startedAt.UnixNano())
	defer func() {
		result := "success"
		errorCode := ""
		reason := ""
		if err != nil {
			result = "fail"
			errorCode = telemetry.TelemetryErrorCode(err)
			reason = err.Error()
		}
		emitErr := telemetry.Emit(*root, *telemetryFile, types.TelemetryEvent{
			EventName:       "memory.write",
			EventVersion:    telemetry.EventSchema,
			TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
			SessionID:       telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:         traceID,
			ScenarioID:      telemetry.NormalizeTelemetryValue(*scenarioID, "scenario-manual"),
			Operation:       "write",
			Result:          result,
			PolicyGate:      "medium",
			MemoryType:      telemetry.NormalizeMemoryType(*memoryType),
			LatencyMS:       time.Since(startedAt).Milliseconds(),
			OperatorVerdict: telemetry.NormalizeOperatorVerdict(*operatorVerdict),
			ErrorCode:       errorCode,
			Reason:          reason,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	_, enforceSpan := telemetry.StartSpan(ctx, "memory.constraint_check")
	err = governance.EnforceConstraintChecks("write", *sessionID, *scenarioID, traceID)
	telemetry.EndSpan(enforceSpan, err)
	if err != nil {
		return err
	}

	var policy types.WritePolicyDecision
	_, policySpan := telemetry.StartSpan(ctx, "memory.policy.enforce")
	policy, err = governance.EnforceWritePolicy(types.WritePolicyInput{
		Stage:        *stage,
		Reviewer:     *reviewer,
		ApprovedFlag: *approved,
		Decision:     *decision,
		Notes:        *notes,
		Reason:       *reason,
		Risk:         *risk,
		ReworkNotes:  *reworkNotes,
		ReReviewedBy: *reReviewedBy,
	})
	telemetry.EndSpan(policySpan, err)
	if err != nil {
		return err
	}

	_, upsertSpan := telemetry.StartSpan(ctx, "memory.index.upsert")
	err = index.UpsertEntry(*root, types.UpsertEntryInput{
		ID:         *id,
		Title:      *title,
		Type:       *typeValue,
		Domain:     *domain,
		Body:       *body,
		BodyFile:   *bodyFile,
		Stage:      *stage,
		SourceRef:  *sourceRef,
		SourceKind: *sourceKind,
		SourceType: *sourceType,
	}, policy)
	telemetry.EndSpan(upsertSpan, err)
	if err != nil {
		return err
	}
	_, embedSpan := telemetry.StartSpan(ctx, "memory.embedding.index_entry")
	warning, err := retrieval.IndexEntryEmbedding(*root, *id, *embeddingEndpoint, *sessionID)
	telemetry.EndSpan(embedSpan, err)
	if err != nil {
		return err
	}
	if strings.TrimSpace(warning) != "" {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
	}

	dirName := "prompts"
	if *typeValue == "instruction" {
		dirName = "instructions"
	}
	if *typeValue == "note" {
		dirName = "notes"
	}
	fmt.Printf("wrote entry %s at %s\n", *id, fmt.Sprintf("%s/%s/%s.md", dirName, *domain, *id))
	return nil
}

func runRetrieve(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "retrieve")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("retrieve", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	query := fs.String("query", "", "natural language query")
	domain := fs.String("domain", "", "optional domain filter")
	sessionID := fs.String("session-id", "session-local", "telemetry session identifier")
	scenarioID := fs.String("scenario-id", "scenario-manual", "telemetry scenario identifier")
	memoryType := fs.String("memory-type", "semantic", "telemetry memory type: procedural|state|semantic")
	operatorVerdict := fs.String("operator-verdict", "not_scored", "telemetry operator verdict")
	telemetryFile := fs.String("telemetry-file", "", "optional telemetry output file (default: <root>/telemetry/events.jsonl)")
	embeddingEndpoint := fs.String("embedding-endpoint", retrieval.DefaultEmbeddingEndpoint, "embedding service endpoint")
	mode := fs.String("mode", "classic", "retrieval mode: classic|hybrid")
	topK := fs.Int("top-k", 5, "number of candidate traces to return (1-50)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("memory.query", *query),
		attribute.String("memory.domain", *domain),
		attribute.String("session.id", *sessionID),
	)
	startedAt := time.Now().UTC()
	traceID := fmt.Sprintf("trace-%d", startedAt.UnixNano())
	telemetryResult := types.RetrieveResult{
		SelectedID:    "__none__",
		SelectionMode: "fallback_path_priority",
		SourcePath:    "__none__",
	}
	defer func() {
		result := "success"
		errorCode := ""
		reason := telemetryResult.Reason
		if err != nil {
			result = "fail"
			errorCode = telemetry.TelemetryErrorCode(err)
			reason = err.Error()
		}
		semanticHit := telemetryResult.SelectionMode == "semantic" ||
			telemetryResult.SelectionMode == "embedding_semantic" ||
			telemetryResult.SelectionMode == "hybrid_rrf"
		fallbackUsed := strings.HasPrefix(telemetryResult.SelectionMode, "fallback_")
		semanticRate := 0.0
		fallbackRate := 0.0
		if rate, ferr := telemetry.EmitRetrievalMetric(*root, telemetryResult); ferr == nil {
			semanticRate = rate.SemanticHitRate
			fallbackRate = rate.FallbackRate
		}
		emitErr := telemetry.Emit(*root, *telemetryFile, types.TelemetryEvent{
			EventName:       "memory.retrieve",
			EventVersion:    telemetry.EventSchema,
			TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
			SessionID:       telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:         traceID,
			ScenarioID:      telemetry.NormalizeTelemetryValue(*scenarioID, "scenario-manual"),
			Operation:       "retrieve",
			Result:          result,
			PolicyGate:      "none",
			MemoryType:      telemetry.NormalizeMemoryType(*memoryType),
			LatencyMS:       time.Since(startedAt).Milliseconds(),
			SelectedID:      telemetryResult.SelectedID,
			SelectionMode:   telemetryResult.SelectionMode,
			SourcePath:      telemetryResult.SourcePath,
			OperatorVerdict: telemetry.NormalizeOperatorVerdict(*operatorVerdict),
			ErrorCode:       errorCode,
			Reason:          reason,
			FallbackUsed:    fallbackUsed,
			SemanticHit:     semanticHit,
			PrecisionProxy:  telemetryResult.PrecisionHint,
			SemanticHitRate: semanticRate,
			FallbackRate:    fallbackRate,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	_, enforceSpan := telemetry.StartSpan(ctx, "memory.constraint_check")
	err = governance.EnforceConstraintChecks("retrieve", *sessionID, *scenarioID, traceID)
	telemetry.EndSpan(enforceSpan, err)
	if err != nil {
		return err
	}

	_, retrieveSpan := telemetry.StartSpan(ctx, "memory.retrieve.execute")
	result, warning, err := retrieval.RetrieveWithOptionsAndEndpointAndSession(
		*root,
		*query,
		*domain,
		*embeddingEndpoint,
		*sessionID,
		retrieval.RetrieveOptions{
			Mode: *mode,
			TopK: *topK,
		},
	)
	telemetry.EndSpan(retrieveSpan, err)
	if err != nil {
		return err
	}
	if strings.TrimSpace(warning) != "" {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
	}
	telemetryResult = result
	return printResult(result)
}

func runBootstrap(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "bootstrap")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	repo := fs.String("repo", "", "repository identifier")
	sessionID := fs.String("session-id", "", "telemetry session identifier")
	scenario := fs.String("scenario", "", "scenario identifier")
	memoryType := fs.String("memory-type", "state", "telemetry memory type: procedural|state|semantic")
	operatorVerdict := fs.String("operator-verdict", "not_scored", "telemetry operator verdict")
	telemetryFile := fs.String("telemetry-file", "", "optional telemetry output file (default: <root>/telemetry/events.jsonl)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*repo) == "" {
		return errors.New("--repo is required")
	}
	if strings.TrimSpace(*sessionID) == "" {
		return errors.New("--session-id is required")
	}
	if strings.TrimSpace(*scenario) == "" {
		return errors.New("--scenario is required")
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("repo.id", *repo),
		attribute.String("session.id", *sessionID),
		attribute.String("scenario.id", *scenario),
	)

	startedAt := time.Now().UTC()
	traceID := fmt.Sprintf("trace-%d", startedAt.UnixNano())
	telemetryResult := types.BootstrapPayload{}
	defer func() {
		result := "success"
		errorCode := ""
		reason := ""
		selectedID := "__none__"
		selectionMode := "bootstrap_empty"
		sourcePath := "__none__"
		if len(telemetryResult.MemoryEntries) > 0 {
			selectedID = telemetryResult.MemoryEntries[0].ID
			selectionMode = telemetryResult.MemoryEntries[0].SelectionMode
			sourcePath = telemetryResult.MemoryEntries[0].SourcePath
		}
		if err != nil {
			result = "fail"
			errorCode = telemetry.TelemetryErrorCode(err)
			reason = err.Error()
		}
		emitErr := telemetry.Emit(*root, *telemetryFile, types.TelemetryEvent{
			EventName:       "memory.bootstrap",
			EventVersion:    telemetry.EventSchema,
			TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
			SessionID:       telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:         traceID,
			ScenarioID:      telemetry.NormalizeTelemetryValue(*scenario, "scenario-bootstrap"),
			Operation:       "bootstrap",
			Result:          result,
			PolicyGate:      "none",
			MemoryType:      telemetry.NormalizeMemoryType(*memoryType),
			LatencyMS:       time.Since(startedAt).Milliseconds(),
			SelectedID:      selectedID,
			SelectionMode:   selectionMode,
			SourcePath:      sourcePath,
			OperatorVerdict: telemetry.NormalizeOperatorVerdict(*operatorVerdict),
			ErrorCode:       errorCode,
			Reason:          reason,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	_, enforceSpan := telemetry.StartSpan(ctx, "memory.constraint_check")
	err = governance.EnforceConstraintChecks("retrieve", *sessionID, *scenario, traceID)
	telemetry.EndSpan(enforceSpan, err)
	if err != nil {
		return err
	}

	_, bootstrapSpan := telemetry.StartSpan(ctx, "memory.bootstrap.execute")
	payload, err := retrieval.Bootstrap(*root, *repo, *sessionID, *scenario)
	telemetry.EndSpan(bootstrapSpan, err)
	if err != nil {
		return err
	}
	telemetryResult = payload
	return printResult(payload)
}

func runSnapshot(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: memory-cli snapshot <create|list|restore> [flags]")
	}
	switch args[0] {
	case "create":
		return runSnapshotCreate(args[1:])
	case "list":
		return runSnapshotList(args[1:])
	case "restore":
		return runSnapshotRestore(args[1:])
	default:
		return fmt.Errorf("unknown snapshot subcommand: %s", args[0])
	}
}

func runVerify(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: memory-cli verify <embeddings|health> [flags]")
	}
	switch args[0] {
	case "embeddings":
		return runVerifyEmbeddings(args[1:])
	case "health":
		return runVerifyHealth(args[1:])
	default:
		return fmt.Errorf("unknown verify subcommand: %s", args[0])
	}
}

func runVerifyEmbeddings(args []string) error {
	fs := flag.NewFlagSet("verify embeddings", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	showMissing := fs.Bool("show-missing", true, "include missing entry ids")
	if err := fs.Parse(args); err != nil {
		return err
	}

	idx, err := index.LoadIndex(*root)
	if err != nil {
		return err
	}
	embeddings, err := index.GetEmbeddingRecords(*root, nil)
	if err != nil {
		return err
	}

	missing := make([]string, 0)
	for _, e := range idx.Entries {
		if rec, ok := embeddings[e.ID]; !ok || len(rec.Vector) == 0 {
			missing = append(missing, e.ID)
		}
	}
	sort.Strings(missing)

	provider := "ollama_or_custom_endpoint"
	if strings.TrimSpace(os.Getenv("AZURE_OPENAI_ENDPOINT")) != "" {
		provider = "azure_openai"
	}

	report := struct {
		Root              string   `json:"root"`
		ProviderHint      string   `json:"provider_hint"`
		IndexedEntries    int      `json:"indexed_entries"`
		StoredEmbeddings  int      `json:"stored_embeddings"`
		MissingEmbeddings int      `json:"missing_embeddings"`
		MissingIDs        []string `json:"missing_ids,omitempty"`
	}{
		Root:              *root,
		ProviderHint:      provider,
		IndexedEntries:    len(idx.Entries),
		StoredEmbeddings:  len(embeddings),
		MissingEmbeddings: len(missing),
	}
	if *showMissing {
		report.MissingIDs = missing
	}

	return printResult(report)
}

func runVerifyHealth(args []string) error {
	fs := flag.NewFlagSet("verify health", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	query := fs.String("query", "memory lifecycle", "health-check query")
	domain := fs.String("domain", "", "optional domain filter")
	sessionID := fs.String("session-id", "session-local", "session identifier")
	embeddingEndpoint := fs.String("embedding-endpoint", retrieval.DefaultEmbeddingEndpoint, "embedding service endpoint")
	if err := fs.Parse(args); err != nil {
		return err
	}

	report, err := retrieval.EvaluateSemanticHealth(*root, *query, *domain, *embeddingEndpoint, *sessionID)
	if err != nil {
		return err
	}
	health := struct {
		Root              string  `json:"root"`
		Query             string  `json:"query"`
		IndexedEntries    int     `json:"indexed_entries"`
		StoredEmbeddings  int     `json:"stored_embeddings"`
		MissingEmbeddings int     `json:"missing_embeddings"`
		CoverageOK        bool    `json:"coverage_ok"`
		SelectionMode     string  `json:"selection_mode"`
		SemanticAvailable bool    `json:"semantic_available"`
		Warning           string  `json:"warning,omitempty"`
		SelectedID        string  `json:"selected_id"`
		Confidence        float64 `json:"confidence"`
		Pass              bool    `json:"pass"`
	}{
		Root:              *root,
		Query:             *query,
		IndexedEntries:    report.IndexedEntries,
		StoredEmbeddings:  report.StoredEmbeddings,
		MissingEmbeddings: report.MissingEmbeddings,
		CoverageOK:        report.CoverageOK,
		SelectionMode:     report.SelectionMode,
		SemanticAvailable: report.SemanticAvailable,
		Warning:           report.Warning,
		SelectedID:        report.SelectedID,
		Confidence:        report.Confidence,
		Pass:              report.Pass,
	}
	return printResult(health)
}

func runSnapshotCreate(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "snapshot.create")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("snapshot create", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	createdBy := fs.String("created-by", "", "snapshot actor")
	reason := fs.String("reason", "", "snapshot rationale")
	scope := fs.String("scope", "full", "snapshot scope")
	sessionID := fs.String("session-id", "session-local", "session identifier")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*createdBy) == "" {
		return errors.New("--created-by is required")
	}
	if strings.TrimSpace(*reason) == "" {
		return errors.New("--reason is required")
	}
	if strings.TrimSpace(*scope) != "full" {
		return errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: only scope=full is supported in v0.2 MVP")
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("session.id", *sessionID),
		attribute.String("snapshot.scope", *scope),
	)

	traceID := fmt.Sprintf("trace-%d", time.Now().UTC().UnixNano())
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:  "snapshot.create.requested",
		SnapshotID: "pending",
		SessionID:  telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})

	_, snapshotSpan := telemetry.StartSpan(ctx, "memory.snapshot.create")
	manifest, err := snapshot.CreateSnapshot(*root, *createdBy, *reason)
	telemetry.EndSpan(snapshotSpan, err)
	if err != nil {
		_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
			EventName:   "snapshot.create.completed",
			SnapshotID:  "pending",
			SessionID:   telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetry.TelemetryErrorCode(err),
			Description: err.Error(),
		})
		return err
	}
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:  "snapshot.create.completed",
		SnapshotID: manifest.SnapshotID,
		SessionID:  telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})
	return printResult(manifest)
}

func runSnapshotList(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "snapshot.list")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("snapshot list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	commandSpan.SetAttributes(attribute.String("memory.root", *root))
	_, listSpan := telemetry.StartSpan(ctx, "memory.snapshot.list")
	rows, err := snapshot.ListSnapshots(*root)
	telemetry.EndSpan(listSpan, err)
	if err != nil {
		return err
	}
	return printResult(rows)
}

func runSnapshotRestore(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "snapshot.restore")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("snapshot restore", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	snapshotID := fs.String("snapshot-id", "", "snapshot id")
	sessionID := fs.String("session-id", "session-local", "session identifier")
	stage := fs.String("stage", "pm", "policy stage: planning|architect|pm")
	reviewer := fs.String("reviewer", "", "reviewer identity")
	approved := fs.Bool("approved", false, "legacy flag: equivalent to --decision=approved")
	decision := fs.String("decision", "", "review decision: approved|rejected")
	notes := fs.String("notes", "", "decision notes")
	reason := fs.String("reason", "", "reason for restore")
	risk := fs.String("risk", "", "risk and mitigation note")
	reworkNotes := fs.String("rework-notes", "", "required when --decision=rejected")
	reReviewedBy := fs.String("re-reviewed-by", "", "required when --decision=rejected")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*snapshotID) == "" {
		return errors.New("--snapshot-id is required")
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("snapshot.id", *snapshotID),
		attribute.String("session.id", *sessionID),
	)

	traceID := fmt.Sprintf("trace-%d", time.Now().UTC().UnixNano())
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:  "snapshot.restore.requested",
		SnapshotID: *snapshotID,
		SessionID:  telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})

	_, policySpan := telemetry.StartSpan(ctx, "memory.policy.enforce")
	policy, policyErr := governance.EnforceWritePolicy(types.WritePolicyInput{
		Stage:        *stage,
		Reviewer:     *reviewer,
		ApprovedFlag: *approved,
		Decision:     *decision,
		Notes:        *notes,
		Reason:       *reason,
		Risk:         *risk,
		ReworkNotes:  *reworkNotes,
		ReReviewedBy: *reReviewedBy,
	})
	telemetry.EndSpan(policySpan, policyErr)
	if policyErr != nil {
		err = fmt.Errorf("ERR_SNAPSHOT_RESTORE_POLICY_BLOCKED: %w", policyErr)
		_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
			EventName:   "snapshot.restore.policy_decision",
			SnapshotID:  *snapshotID,
			SessionID:   telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetry.TelemetryErrorCode(err),
			Description: err.Error(),
		})
		_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
			EventName:   "snapshot.restore.failed",
			SnapshotID:  *snapshotID,
			SessionID:   telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetry.TelemetryErrorCode(err),
			Description: err.Error(),
		})
		return err
	}
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:   "snapshot.restore.policy_decision",
		SnapshotID:  *snapshotID,
		SessionID:   telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:     traceID,
		Result:      "success",
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Description: policy.Decision,
	})

	_, restoreSpan := telemetry.StartSpan(ctx, "memory.snapshot.restore")
	err = snapshot.RestoreSnapshot(*root, *snapshotID)
	telemetry.EndSpan(restoreSpan, err)
	if err != nil {
		_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
			EventName:   "snapshot.restore.failed",
			SnapshotID:  *snapshotID,
			SessionID:   telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetry.TelemetryErrorCode(err),
			Description: err.Error(),
		})
		return err
	}
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:  "snapshot.restore.completed",
		SnapshotID: *snapshotID,
		SessionID:  telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})

	return nil
}

func printResult(r any) error {
	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
