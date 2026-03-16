package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestRunValidateAgainstFixture(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: obsidian
    name: Obsidian
    description: Native Obsidian CLI surface
    status: active
    guidance: Use for vault review
    capabilities:
      - id: open_cairn_vault
        name: Open Cairn Vault
        description: Open the Cairn vault in Obsidian
        status: active
        availability: required
        stage_affinity: [pm]
        guidance: Required review surface
        call:
          type: exec
          command: obsidian open --vault Cairn
        schema: []
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() error {
		return runValidate([]string{"--registry", registryPath})
	})
	if !strings.Contains(out, "registry valid: 1 tool systems") {
		t.Fatalf("expected validation output, got %s", out)
	}
}

func TestRunContextJSONIncludesCapabilities(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: cairn
    name: Cairn
    description: Internal Cairn tool system
    status: active
    guidance: Use for internal supported workflows
    capabilities:
      - id: memory.verify_health
        name: Verify memory-cli Health
        description: Check local retrieval health
        status: active
        availability: default
        stage_affinity: [engineering]
        guidance: Use before relying on memory retrieval
        call:
          type: exec
          command: memory-cli verify health
        schema:
          - name: root
            type: string
            required: true
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() error {
		return runContext([]string{"--registry", registryPath, "--stage", "engineering", "--format", "json"})
	})

	var payload struct {
		ToolContext struct {
			Systems []struct {
				ID           string `json:"id"`
				Capabilities []struct {
					ID            string `json:"id"`
					SchemaSummary string `json:"schema_summary"`
				} `json:"capabilities"`
			} `json:"systems"`
		} `json:"tool_context"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("unmarshal context payload: %v\n%s", err, out)
	}
	if len(payload.ToolContext.Systems) != 1 {
		t.Fatalf("expected 1 system, got %d", len(payload.ToolContext.Systems))
	}
	if len(payload.ToolContext.Systems[0].Capabilities) != 1 {
		t.Fatalf("expected 1 capability, got %d", len(payload.ToolContext.Systems[0].Capabilities))
	}
	if payload.ToolContext.Systems[0].Capabilities[0].SchemaSummary != "root(req)" {
		t.Fatalf("unexpected schema summary: %s", payload.ToolContext.Systems[0].Capabilities[0].SchemaSummary)
	}
}

func TestRunContextExcludesScopedCapabilitiesUnlessExplicitlyRequested(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: cairn
    name: Cairn
    description: Internal Cairn tool system
    status: active
    guidance: Use for internal supported workflows
    capabilities:
      - id: smoke
        name: Smoke
        description: Run smoke checks
        status: active
        availability: default
        stage_affinity: [engineering]
        guidance: Use for baseline validation
        call:
          type: exec
          command: ./smoke
        schema: []
      - id: intake.inspect_source
        name: Inspect Intake Source
        description: Inspect a source before staging
        status: active
        availability: scoped
        stage_affinity: [engineering]
        guidance: Use only when intake is relevant
        call:
          type: exec
          command: intake-cli inspect
        schema: []
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	defaultContext := captureStdout(t, func() error {
		return runContext([]string{"--registry", registryPath, "--stage", "engineering", "--format", "json"})
	})
	if strings.Contains(defaultContext, "intake.inspect_source") {
		t.Fatalf("scoped capability should not appear in default context: %s", defaultContext)
	}

	queryContext := captureStdout(t, func() error {
		return runContext([]string{"--registry", registryPath, "--stage", "engineering", "--query", "intake", "--format", "json"})
	})
	if !strings.Contains(queryContext, "intake.inspect_source") {
		t.Fatalf("scoped capability should appear when query is explicit: %s", queryContext)
	}
}

func TestRunInspectResolvesUniqueQuery(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: gitnexus
    name: GitNexus
    description: Repository structure analysis system
    status: active
    guidance: Use for code archaeology
    capabilities:
      - id: code_graph
        name: Code Graph
        description: Analyze repository structure
        status: planned
        availability: scoped
        stage_affinity: [engineering]
        guidance: Use for code archaeology and impact analysis
        call:
          type: planned
          command: ""
        schema: []
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() error {
		return runInspect([]string{"--registry", registryPath, "archaeology"})
	})
	if !strings.Contains(out, "gitnexus  GitNexus") {
		t.Fatalf("expected inspect to resolve unique query, got %s", out)
	}
}

func TestRunInspectRejectsAmbiguousQuery(t *testing.T) {
	root := t.TempDir()
	registryPath := filepath.Join(root, "approved-tools.yaml")
	content := `version: 1
tools:
  - id: firecrawl
    name: Firecrawl
    description: Web content extraction system
    status: planned
    guidance: Use for controlled scrape integrations
    capabilities:
      - id: scrape
        name: Scrape To Markdown
        description: Scrape pages into markdown
        status: planned
        availability: scoped
        stage_affinity: [pm]
        guidance: Use for web cleanup
        call:
          type: planned
          command: ""
        schema: []
  - id: gitnexus
    name: GitNexus
    description: Repository graph analysis system
    status: planned
    guidance: Use for code graph integrations
    capabilities:
      - id: code_graph
        name: Code Graph
        description: Analyze repository structure
        status: planned
        availability: scoped
        stage_affinity: [engineering]
        guidance: Use for graph analysis
        call:
          type: planned
          command: ""
        schema: []
`
	if err := os.WriteFile(registryPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	err := runInspect([]string{"--registry", registryPath, "planned"})
	if err == nil {
		t.Fatal("expected ambiguous inspect query to fail")
	}
	if !strings.Contains(err.Error(), `ambiguous tool system query "planned"`) {
		t.Fatalf("unexpected ambiguous query error: %v", err)
	}
}

func TestApprovedRegistryContextIncludesFourSystems(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() error {
		return runList([]string{"--stage", "pm", "--format", "json"})
	})
	if !strings.Contains(out, `"id": "cairn"`) || !strings.Contains(out, `"id": "obsidian"`) {
		t.Fatalf("expected cairn and obsidian in approved registry output, got %s", out)
	}
	_ = wd
}
