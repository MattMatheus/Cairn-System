package registry

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"athenause/internal/types"
)

const (
	SupportTierApproved = "approved"
	SupportTierLocal    = "local"
)

type LoadOptions struct {
	ApprovedPath string
	LocalPath    string
	IncludeLocal bool
}

func ResolvePaths(cwd string) (approvedPath, localPath string, err error) {
	if override := strings.TrimSpace(os.Getenv("ATHENA_USE_REGISTRY")); override != "" {
		return override, "", nil
	}

	root, err := findRepoRoot(cwd)
	if err != nil {
		return "", "", err
	}
	approvedPath = filepath.Join(root, "products", "athena-use", "registry", "approved-tools.yaml")
	athenaHome := strings.TrimSpace(os.Getenv("ATHENA_HOME"))
	if athenaHome == "" {
		athenaHome = filepath.Join(root, ".athena")
	}
	localPath = filepath.Join(athenaHome, "tools", "registry.yaml")
	return approvedPath, localPath, nil
}

func LoadMerged(opts LoadOptions) (types.Registry, error) {
	approved, err := LoadFile(opts.ApprovedPath, SupportTierApproved)
	if err != nil {
		return types.Registry{}, err
	}
	if !opts.IncludeLocal || strings.TrimSpace(opts.LocalPath) == "" {
		return approved, nil
	}
	if _, err := os.Stat(opts.LocalPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return approved, nil
		}
		return types.Registry{}, err
	}
	local, err := LoadFile(opts.LocalPath, SupportTierLocal)
	if err != nil {
		return types.Registry{}, err
	}
	return mergeRegistries(approved, local)
}

func LoadFile(path, tier string) (types.Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return types.Registry{}, err
	}
	reg, err := parseYAML(data)
	if err != nil {
		return types.Registry{}, fmt.Errorf("%s: %w", path, err)
	}
	for i := range reg.Tools {
		reg.Tools[i].SupportTier = tier
	}
	return reg, nil
}

func Validate(reg types.Registry) error {
	if reg.Version != 1 {
		return fmt.Errorf("unsupported registry version: %d", reg.Version)
	}
	seen := map[string]struct{}{}
	validStages := map[string]struct{}{
		"planning": {}, "architect": {}, "engineering": {}, "qa": {}, "pm": {}, "cycle": {}, "release": {},
	}
	for _, tool := range reg.Tools {
		if strings.TrimSpace(tool.ID) == "" {
			return errors.New("tool id is required")
		}
		if _, ok := seen[tool.ID]; ok {
			return fmt.Errorf("duplicate tool id: %s", tool.ID)
		}
		seen[tool.ID] = struct{}{}
		if strings.TrimSpace(tool.Name) == "" {
			return fmt.Errorf("tool %s missing name", tool.ID)
		}
		if strings.TrimSpace(tool.Description) == "" {
			return fmt.Errorf("tool %s missing description", tool.ID)
		}
		if strings.TrimSpace(tool.Call.Type) == "" {
			return fmt.Errorf("tool %s missing call.type", tool.ID)
		}
		for _, stage := range tool.StageAffinity {
			if _, ok := validStages[stage]; !ok {
				return fmt.Errorf("tool %s has invalid stage_affinity value: %s", tool.ID, stage)
			}
		}
		for _, field := range tool.Schema {
			if strings.TrimSpace(field.Name) == "" {
				return fmt.Errorf("tool %s has schema entry missing name", tool.ID)
			}
			if strings.TrimSpace(field.Type) == "" {
				return fmt.Errorf("tool %s schema field %s missing type", tool.ID, field.Name)
			}
		}
	}
	return nil
}

func FilterByStage(tools []types.Tool, stage string) []types.Tool {
	stage = strings.TrimSpace(stage)
	if stage == "" {
		return append([]types.Tool(nil), tools...)
	}
	var filtered []types.Tool
	for _, tool := range tools {
		if len(tool.StageAffinity) == 0 {
			filtered = append(filtered, tool)
			continue
		}
		for _, candidate := range tool.StageAffinity {
			if candidate == stage {
				filtered = append(filtered, tool)
				break
			}
		}
	}
	return filtered
}

func FilterByTag(tools []types.Tool, tag string) []types.Tool {
	tag = strings.TrimSpace(strings.ToLower(tag))
	if tag == "" {
		return append([]types.Tool(nil), tools...)
	}
	var filtered []types.Tool
	for _, tool := range tools {
		for _, candidate := range tool.Tags {
			if strings.ToLower(candidate) == tag {
				filtered = append(filtered, tool)
				break
			}
		}
	}
	return filtered
}

func Discover(tools []types.Tool, query string) []types.Tool {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return append([]types.Tool(nil), tools...)
	}
	terms := strings.Fields(query)
	type scored struct {
		tool  types.Tool
		score int
	}
	var matches []scored
	for _, tool := range tools {
		haystack := strings.ToLower(strings.Join([]string{
			tool.ID,
			tool.Name,
			tool.Description,
			strings.Join(tool.Tags, " "),
			strings.Join(tool.StageAffinity, " "),
		}, " "))
		score := 0
		for _, term := range terms {
			if strings.Contains(haystack, term) {
				score++
			}
		}
		if score > 0 {
			matches = append(matches, scored{tool: tool, score: score})
		}
	}
	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].score == matches[j].score {
			return matches[i].tool.ID < matches[j].tool.ID
		}
		return matches[i].score > matches[j].score
	})
	out := make([]types.Tool, 0, len(matches))
	for _, match := range matches {
		out = append(out, match.tool)
	}
	return out
}

func mergeRegistries(approved, local types.Registry) (types.Registry, error) {
	merged := types.Registry{Version: 1}
	merged.Tools = append(merged.Tools, approved.Tools...)
	seen := map[string]struct{}{}
	for _, tool := range approved.Tools {
		seen[tool.ID] = struct{}{}
	}
	for _, tool := range local.Tools {
		if _, ok := seen[tool.ID]; ok {
			return types.Registry{}, fmt.Errorf("duplicate tool id across approved and local registries: %s", tool.ID)
		}
		merged.Tools = append(merged.Tools, tool)
	}
	return merged, nil
}

func findRepoRoot(start string) (string, error) {
	dir := start
	for {
		candidate := filepath.Join(dir, "products", "athena-use", "registry", "approved-tools.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("unable to resolve AthenaPlatform repo root")
		}
		dir = parent
	}
}

func parseYAML(data []byte) (types.Registry, error) {
	lines := preprocessLines(data)
	reg := types.Registry{}
	var currentTool *types.Tool
	var currentSchema *types.SchemaField
	inTools := false
	inSchema := false
	inCall := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		indent := countIndent(line)
		switch {
		case trimmed == "tools:":
			inTools = true
			inSchema = false
			inCall = false
		case strings.HasPrefix(trimmed, "version:"):
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "version:"))
			version, err := strconv.Atoi(value)
			if err != nil {
				return types.Registry{}, fmt.Errorf("parse version: %w", err)
			}
			reg.Version = version
		case inTools && trimmed == "tools: []":
			reg.Tools = nil
		case inTools && indent == 2 && strings.HasPrefix(trimmed, "- "):
			inSchema = false
			inCall = false
			currentSchema = nil
			reg.Tools = append(reg.Tools, types.Tool{})
			currentTool = &reg.Tools[len(reg.Tools)-1]
			if err := assignToolField(currentTool, strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))); err != nil {
				return types.Registry{}, err
			}
		case currentTool != nil && indent == 4 && trimmed == "call:":
			inCall = true
			inSchema = false
			currentSchema = nil
		case currentTool != nil && indent == 4 && trimmed == "schema:":
			inSchema = true
			inCall = false
		case currentTool != nil && inCall && indent >= 6:
			if err := assignCallField(&currentTool.Call, trimmed); err != nil {
				return types.Registry{}, err
			}
		case currentTool != nil && inSchema && indent >= 6 && strings.HasPrefix(strings.TrimSpace(line), "- "):
			currentTool.Schema = append(currentTool.Schema, types.SchemaField{})
			currentSchema = &currentTool.Schema[len(currentTool.Schema)-1]
			if err := assignSchemaField(currentSchema, strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "- "))); err != nil {
				return types.Registry{}, err
			}
		case currentTool != nil && inSchema && currentSchema != nil && indent >= 8:
			if err := assignSchemaField(currentSchema, trimmed); err != nil {
				return types.Registry{}, err
			}
		case currentTool != nil && indent >= 4:
			inCall = false
			inSchema = false
			currentSchema = nil
			if err := assignToolField(currentTool, trimmed); err != nil {
				return types.Registry{}, err
			}
		}
	}

	if reg.Version == 0 {
		return types.Registry{}, errors.New("version is required")
	}
	return reg, nil
}

func preprocessLines(data []byte) []string {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "#"); idx >= 0 {
			line = line[:idx]
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

func countIndent(line string) int {
	count := 0
	for _, r := range line {
		if r != ' ' {
			break
		}
		count++
	}
	return count
}

func assignToolField(tool *types.Tool, raw string) error {
	key, value, err := splitKeyValue(raw)
	if err != nil {
		return err
	}
	switch key {
	case "id":
		tool.ID = value
	case "name":
		tool.Name = value
	case "description":
		tool.Description = value
	case "tags":
		tool.Tags = parseInlineList(value)
	case "stage_affinity":
		tool.StageAffinity = parseInlineList(value)
	case "credential":
		tool.CredentialRef = value
	case "support_tier":
		tool.SupportTier = value
	}
	return nil
}

func assignCallField(call *types.ToolCall, raw string) error {
	key, value, err := splitKeyValue(raw)
	if err != nil {
		return err
	}
	switch key {
	case "type":
		call.Type = value
	case "method":
		call.Method = value
	case "url":
		call.URL = value
	case "command":
		call.Command = value
	}
	return nil
}

func assignSchemaField(field *types.SchemaField, raw string) error {
	key, value, err := splitKeyValue(raw)
	if err != nil {
		return err
	}
	switch key {
	case "name":
		field.Name = value
	case "type":
		field.Type = value
	case "required":
		field.Required = strings.EqualFold(value, "true")
	case "description":
		field.Description = value
	case "enum":
		field.Enum = parseInlineList(value)
	}
	return nil
}

func splitKeyValue(raw string) (string, string, error) {
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected key:value pair, got %q", raw)
	}
	return strings.TrimSpace(parts[0]), trimQuotes(strings.TrimSpace(parts[1])), nil
}

func trimQuotes(v string) string {
	v = strings.TrimSpace(v)
	v = strings.Trim(v, "\"")
	v = strings.Trim(v, "'")
	return v
}

func parseInlineList(v string) []string {
	v = strings.TrimSpace(v)
	if v == "" || v == "[]" {
		return nil
	}
	v = strings.TrimPrefix(v, "[")
	v = strings.TrimSuffix(v, "]")
	if strings.TrimSpace(v) == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		item := trimQuotes(part)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}
