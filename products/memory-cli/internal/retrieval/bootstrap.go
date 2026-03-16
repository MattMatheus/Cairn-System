package retrieval

import (
	"sort"
	"strings"
	"time"

	"memorycli/internal/index"
	"memorycli/internal/types"
)

func Bootstrap(root, repo, sessionID, scenario string) (types.BootstrapPayload, error) {
	payload := types.BootstrapPayload{
		Repo:          strings.TrimSpace(repo),
		SessionID:     strings.TrimSpace(sessionID),
		Scenario:      strings.TrimSpace(scenario),
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		MemoryEntries: []types.BootstrapMemoryEntry{},
	}

	idx, err := index.LoadIndex(root)
	if err != nil {
		return payload, err
	}

	query := strings.TrimSpace(repo + " " + scenario)
	tokens := tokenSet(query)
	if len(tokens) > 0 {
		type scoredEntry struct {
			entry types.IndexEntry
			score float64
		}
		scored := []scoredEntry{}
		for _, e := range idx.Entries {
			if e.Status != "approved" || e.Type != "instruction" {
				continue
			}
			hay := strings.ToLower(strings.Join([]string{e.ID, e.Title, e.Domain, e.Path}, " "))
			hits := 0
			for tok := range tokens {
				if strings.Contains(hay, tok) {
					hits++
				}
			}
			score := float64(hits) / float64(len(tokens))
			if score > 0 {
				scored = append(scored, scoredEntry{entry: e, score: score})
			}
		}
		sort.SliceStable(scored, func(i, j int) bool {
			if scored[i].score == scored[j].score {
				return scored[i].entry.ID < scored[j].entry.ID
			}
			return scored[i].score > scored[j].score
		})
		for i, s := range scored {
			if i >= 5 {
				break
			}
			payload.MemoryEntries = append(payload.MemoryEntries, types.BootstrapMemoryEntry{
				ID:            s.entry.ID,
				SelectionMode: "semantic",
				SourcePath:    s.entry.Path,
				Confidence:    s.score,
				Reason:        "bootstrap contextual match",
				Type:          s.entry.Type,
				Domain:        s.entry.Domain,
				Title:         s.entry.Title,
			})
		}
	}

	return payload, nil
}
