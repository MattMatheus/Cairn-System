package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"athenause/internal/registry"
	"athenause/internal/types"
)

func captureStdout(t *testing.T, fn func() error) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	if err := fn(); err != nil {
		t.Fatal(err)
	}
	_ = w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestSchemaSummary(t *testing.T) {
	got := schemaSummary([]types.SchemaField{
		{Name: "owner", Required: true},
		{Name: "body", Required: false},
	})
	if got != "owner(req), body(opt)" {
		t.Fatalf("unexpected schema summary: %s", got)
	}
}

func TestRunValidateAgainstFixture(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: github.create_pr
    name: Create Pull Request
    description: Opens a pull request against a target branch
    tags: [github, vcs]
    stage_affinity: [engineering, qa]
    credential: GITHUB_TOKEN
    call:
      type: http
      method: POST
      url: "https://api.github.com/repos/{owner}/{repo}/pulls"
    schema:
      - name: owner
        type: string
        required: true
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() error {
		return runValidate([]string{"--registry", registryPath})
	})
	if !strings.Contains(out, "registry valid") {
		t.Fatalf("expected validation output, got %s", out)
	}
}

func TestRunContextJSONIncludesStructuredParameters(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: athena.memory.verify_health
    name: Verify AthenaMind Health
    description: Checks local AthenaMind retrieval health for a memory root and query
    tags: [athena, memory, validation]
    stage_affinity: [engineering]
    credential: ""
    call:
      type: exec
      command: ./verify-health
    schema:
      - name: root
        type: string
        required: true
        description: Memory root path
      - name: query
        type: string
        required: true
      - name: domain
        type: string
        required: false
        enum: [platform, work]
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() error {
		return runContext([]string{"--registry", registryPath, "--stage", "engineering", "--format", "json"})
	})

	var payload struct {
		ToolContext struct {
			Tools []struct {
				ID            string `json:"id"`
				Schema        string `json:"schema"`
				SchemaSummary string `json:"schema_summary"`
				Parameters    []struct {
					Name        string   `json:"name"`
					Type        string   `json:"type"`
					Required    bool     `json:"required"`
					Description string   `json:"description"`
					Enum        []string `json:"enum"`
				} `json:"parameters"`
			} `json:"tools"`
		} `json:"tool_context"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("unmarshal context payload: %v\n%s", err, out)
	}
	if len(payload.ToolContext.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(payload.ToolContext.Tools))
	}
	tool := payload.ToolContext.Tools[0]
	if tool.SchemaSummary != "root(req), query(req), domain(opt)" {
		t.Fatalf("unexpected schema summary: %s", tool.SchemaSummary)
	}
	if tool.Schema != tool.SchemaSummary {
		t.Fatalf("expected schema alias to match schema_summary, got %q vs %q", tool.Schema, tool.SchemaSummary)
	}
	if len(tool.Parameters) != 3 {
		t.Fatalf("expected 3 parameters, got %d", len(tool.Parameters))
	}
	if tool.Parameters[0].Name != "root" || !tool.Parameters[0].Required || tool.Parameters[0].Description != "Memory root path" {
		t.Fatalf("unexpected first parameter: %+v", tool.Parameters[0])
	}
	if got := strings.Join(tool.Parameters[2].Enum, ","); got != "platform,work" {
		t.Fatalf("unexpected enum values: %s", got)
	}
}

func TestRunContextYAMLIncludesParametersSection(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: athena.memory.verify_health
    name: Verify AthenaMind Health
    description: Checks local AthenaMind retrieval health for a memory root and query
    tags: [athena, memory, validation]
    stage_affinity: [engineering]
    credential: ""
    call:
      type: exec
      command: ./verify-health
    schema:
      - name: root
        type: string
        required: true
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	text := captureStdout(t, func() error {
		return runContext([]string{"--registry", registryPath, "--stage", "engineering", "--format", "yaml"})
	})
	if !strings.Contains(text, "schema_summary: root(req)") {
		t.Fatalf("expected schema_summary in yaml output, got %s", text)
	}
	if !strings.Contains(text, "parameters:") || !strings.Contains(text, "name: root") || !strings.Contains(text, "required: true") {
		t.Fatalf("expected parameters section in yaml output, got %s", text)
	}
}

func TestApprovedRegistryContextIncludesRichParameterMetadata(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	approvedPath, _, err := registry.ResolvePaths(wd)
	if err != nil {
		t.Fatalf("resolve approved path: %v", err)
	}

	out := captureStdout(t, func() error {
		return runContext([]string{"--registry", approvedPath, "--stage", "engineering", "--format", "json"})
	})

	var payload struct {
		ToolContext struct {
			Tools []struct {
				ID         string `json:"id"`
				Parameters []struct {
					Name        string   `json:"name"`
					Description string   `json:"description"`
					Enum        []string `json:"enum"`
				} `json:"parameters"`
			} `json:"tools"`
		} `json:"tool_context"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("unmarshal approved context payload: %v\n%s", err, out)
	}

	toolsByID := map[string][]struct {
		Name        string
		Description string
		Enum        []string
	}{}
	for _, tool := range payload.ToolContext.Tools {
		params := make([]struct {
			Name        string
			Description string
			Enum        []string
		}, len(tool.Parameters))
		for i, param := range tool.Parameters {
			params[i] = struct {
				Name        string
				Description string
				Enum        []string
			}{
				Name:        param.Name,
				Description: param.Description,
				Enum:        param.Enum,
			}
		}
		toolsByID[tool.ID] = params
	}

	verifyHealth, ok := toolsByID["athena.memory.verify_health"]
	if !ok || len(verifyHealth) == 0 || verifyHealth[0].Description == "" {
		t.Fatalf("expected verify_health descriptions in approved context, got %+v", verifyHealth)
	}

	taskMetadata, ok := toolsByID["athena.workspace.validate_task_metadata"]
	if !ok || len(taskMetadata) == 0 {
		t.Fatalf("expected task metadata parameters in approved context, got %+v", taskMetadata)
	}
	if got := strings.Join(taskMetadata[0].Enum, ","); got != "changed,all" {
		t.Fatalf("expected mode enum changed,all in approved context, got %q", got)
	}
}

func TestRunContextExcludesScopedToolsUnlessExplicitlyRequested(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: obsidian.open_vault
    name: Open Athena Vault
    description: Opens the Athena vault in Obsidian
    availability: required
    tags: [obsidian, docs]
    stage_affinity: [pm]
    guidance: Required for reviewing working notes and shared markdown
    credential: ""
    call:
      type: exec
      command: obsidian open --vault Athena
    schema: []
  - id: firecrawl.scrape_markdown
    name: Scrape Page To Markdown
    description: Scrapes a page and converts it to markdown for review
    availability: scoped
    tags: [research, docs, ingestion]
    stage_affinity: [pm]
    guidance: Use when external web content needs cleanup before review
    credential: FIRECRAWL_API_KEY
    call:
      type: http
      method: POST
      url: "https://api.firecrawl.dev/v1/scrape"
    schema:
      - name: url
        type: string
        required: true
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	defaultContext := captureStdout(t, func() error {
		return runContext([]string{"--registry", registryPath, "--stage", "pm", "--format", "json"})
	})
	if strings.Contains(defaultContext, "firecrawl.scrape_markdown") {
		t.Fatalf("scoped tool should not appear in default context: %s", defaultContext)
	}
	if !strings.Contains(defaultContext, "obsidian.open_vault") {
		t.Fatalf("required tool missing from default context: %s", defaultContext)
	}

	queryContext := captureStdout(t, func() error {
		return runContext([]string{"--registry", registryPath, "--stage", "pm", "--query", "markdown scrape", "--format", "json"})
	})
	if !strings.Contains(queryContext, "firecrawl.scrape_markdown") {
		t.Fatalf("scoped tool should appear when query is explicit: %s", queryContext)
	}
}

func TestRunContextExcludesPlannedToolsUnlessExplicitlyIncluded(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: obsidian.open_athena_vault
    name: Open Athena Vault
    status: active
    description: Opens the Athena vault in Obsidian
    availability: required
    tags: [obsidian, docs]
    stage_affinity: [pm]
    guidance: Required review surface for working notes
    credential: ""
    call:
      type: exec
      command: obsidian open --vault Athena
    schema: []
  - id: athena.intake.normalize_source
    name: Normalize Source Into Intake Markdown
    status: planned
    description: Planned intake utility for markdown normalization
    availability: scoped
    tags: [research, docs, ingestion]
    stage_affinity: [pm]
    guidance: Use when external material needs cleanup before review
    credential: ""
    call:
      type: planned
      command: ""
    schema:
      - name: source
        type: string
        required: true
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	defaultContext := captureStdout(t, func() error {
		return runContext([]string{"--registry", registryPath, "--stage", "pm", "--query", "intake markdown", "--format", "json"})
	})
	if strings.Contains(defaultContext, "athena.intake.normalize_source") {
		t.Fatalf("planned tool should not appear without --include-planned: %s", defaultContext)
	}

	plannedContext := captureStdout(t, func() error {
		return runContext([]string{"--registry", registryPath, "--stage", "pm", "--query", "intake markdown", "--include-planned", "--format", "json"})
	})
	if !strings.Contains(plannedContext, "athena.intake.normalize_source") {
		t.Fatalf("planned tool should appear when explicitly included: %s", plannedContext)
	}
}

func TestApprovedRegistryIncludesActiveIntakeToolInScopedQueries(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	approvedPath, _, err := registry.ResolvePaths(wd)
	if err != nil {
		t.Fatalf("resolve approved path: %v", err)
	}

	out := captureStdout(t, func() error {
		return runContext([]string{"--registry", approvedPath, "--stage", "pm", "--query", "intake markdown", "--format", "json"})
	})
	if !strings.Contains(out, "athena.intake.normalize_source") {
		t.Fatalf("expected active intake tool in query-narrowed context, got %s", out)
	}
}

func TestApprovedRegistryIncludesIntakeInspectToolInScopedQueries(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	approvedPath, _, err := registry.ResolvePaths(wd)
	if err != nil {
		t.Fatalf("resolve approved path: %v", err)
	}

	out := captureStdout(t, func() error {
		return runContext([]string{"--registry", approvedPath, "--stage", "pm", "--query", "inspect intake source", "--format", "json"})
	})
	if !strings.Contains(out, "athena.intake.inspect_source") {
		t.Fatalf("expected intake inspect tool in query-narrowed context, got %s", out)
	}
}

func TestRunInspectExplainsScopedToolFit(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: gitnexus.repo_graph
    name: Inspect Repository Graph
    description: Uses GitNexus to explore structural relationships in a messy repository
    availability: scoped
    tags: [engineering, architecture, code]
    stage_affinity: [engineering, pm]
    guidance: Use when code archaeology or impact analysis is needed
    complements: [athena-mind, obsidian]
    credential: ""
    call:
      type: exec
      command: gitnexus query
    schema:
      - name: repo_path
        type: string
        required: true
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() error {
		return runInspect([]string{"--registry", registryPath, "--stage", "pm", "gitnexus.repo_graph"})
	})
	if !strings.Contains(out, "Availability: scoped") {
		t.Fatalf("expected scoped availability in inspect output, got %s", out)
	}
	if !strings.Contains(out, "Context: scoped tool; discoverable and inspectable, but not included by default") {
		t.Fatalf("expected scoped context note, got %s", out)
	}
	if !strings.Contains(out, "Complements: athena-mind, obsidian") {
		t.Fatalf("expected complements in inspect output, got %s", out)
	}
}

func TestRunInspectExplainsPlannedToolFit(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: gitnexus.code_graph
    name: Analyze Repository Structure With GitNexus
    status: planned
    description: Planned thin code graph contract for messy repositories
    availability: scoped
    tags: [code, architecture, impact]
    stage_affinity: [pm]
    guidance: Use when code archaeology or impact analysis is needed
    complements: [obsidian.open_athena_vault, athena.memory.verify_health]
    credential: ""
    call:
      type: planned
      command: ""
    schema:
      - name: repo
        type: string
        required: true
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() error {
		return runInspect([]string{"--registry", registryPath, "--stage", "pm", "gitnexus.code_graph"})
	})
	if !strings.Contains(out, "Status: planned") {
		t.Fatalf("expected planned status in inspect output, got %s", out)
	}
	if !strings.Contains(out, "Context: planned tool; discoverable for planning, but not included in active context by default") {
		t.Fatalf("expected planned context note, got %s", out)
	}
}
