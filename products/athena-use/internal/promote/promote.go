package promote

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type VaultNote struct {
	Path          string
	ID            string
	Title         string
	Domain        string
	NoteType      string
	Status        string
	Sensitivity   string
	SourceOfTruth string
	Body          string
}

type Plan struct {
	MemoryID   string `json:"memory_id" yaml:"memory_id"`
	Title      string `json:"title" yaml:"title"`
	Domain     string `json:"domain" yaml:"domain"`
	MemoryType string `json:"memory_type" yaml:"memory_type"`
	SourceRef  string `json:"source_ref" yaml:"source_ref"`
	SourceKind string `json:"source_kind" yaml:"source_kind"`
	SourceType string `json:"source_type" yaml:"source_type"`
	Body       string `json:"body" yaml:"body"`
}

type frontmatter struct {
	ID            string `yaml:"id"`
	Type          string `yaml:"type"`
	Status        string `yaml:"status"`
	Domain        string `yaml:"domain"`
	Sensitivity   string `yaml:"sensitivity"`
	SourceOfTruth string `yaml:"source_of_truth"`
}

type Options struct {
	ID     string
	Title  string
	Domain string
}

func LoadVaultNote(path string) (VaultNote, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return VaultNote{}, err
	}
	fm, body, err := splitFrontmatter(string(data))
	if err != nil {
		return VaultNote{}, err
	}
	var parsed frontmatter
	if err := yaml.Unmarshal([]byte(fm), &parsed); err != nil {
		return VaultNote{}, fmt.Errorf("parse frontmatter: %w", err)
	}
	title := firstHeading(body)
	if strings.TrimSpace(title) == "" {
		title = fallbackTitle(path)
	}
	return VaultNote{
		Path:          path,
		ID:            strings.TrimSpace(parsed.ID),
		Title:         title,
		Domain:        strings.TrimSpace(parsed.Domain),
		NoteType:      strings.TrimSpace(parsed.Type),
		Status:        strings.TrimSpace(parsed.Status),
		Sensitivity:   strings.TrimSpace(parsed.Sensitivity),
		SourceOfTruth: strings.TrimSpace(parsed.SourceOfTruth),
		Body:          trimTopHeading(strings.TrimSpace(body)),
	}, nil
}

func BuildPlan(note VaultNote, opts Options) (Plan, error) {
	id := firstNonEmpty(strings.TrimSpace(opts.ID), note.ID, slugify(note.Title))
	title := firstNonEmpty(strings.TrimSpace(opts.Title), note.Title)
	domain := firstNonEmpty(strings.TrimSpace(opts.Domain), note.Domain)
	if id == "" || title == "" || domain == "" {
		return Plan{}, errors.New("promotion requires id, title, and domain")
	}
	if strings.EqualFold(note.Sensitivity, "private_personal") {
		return Plan{}, errors.New("refusing to promote private_personal note")
	}
	if strings.TrimSpace(note.Body) == "" {
		return Plan{}, errors.New("note body is empty after frontmatter and title trimming")
	}
	return Plan{
		MemoryID:   id,
		Title:      title,
		Domain:     domain,
		MemoryType: "note",
		SourceRef:  note.Path,
		SourceKind: "obsidian-note",
		SourceType: note.NoteType,
		Body:       note.Body,
	}, nil
}

func splitFrontmatter(text string) (string, string, error) {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	if !strings.HasPrefix(normalized, "---\n") {
		return "", normalized, errors.New("note is missing frontmatter")
	}
	rest := strings.TrimPrefix(normalized, "---\n")
	idx := strings.Index(rest, "\n---\n")
	if idx < 0 {
		return "", normalized, errors.New("note frontmatter is not terminated")
	}
	return rest[:idx], rest[idx+5:], nil
}

func firstHeading(body string) string {
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
		}
	}
	return ""
}

func trimTopHeading(body string) string {
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
		}
		break
	}
	return body
}

func fallbackTitle(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

var slugPattern = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	lowered := strings.ToLower(strings.TrimSpace(value))
	lowered = slugPattern.ReplaceAllString(lowered, "-")
	return strings.Trim(lowered, "-")
}
