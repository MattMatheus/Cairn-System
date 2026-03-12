package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"athenause/internal/types"
)

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
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	if err := runValidate([]string{"--registry", registryPath}); err != nil {
		t.Fatalf("run validate: %v", err)
	}
	_ = w.Close()
	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "registry valid") {
		t.Fatalf("expected validation output, got %s", string(out))
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

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	if err := runContext([]string{"--registry", registryPath, "--stage", "engineering", "--format", "json"}); err != nil {
		t.Fatalf("run context json: %v", err)
	}
	_ = w.Close()
	out, _ := io.ReadAll(r)

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
	if err := json.Unmarshal(out, &payload); err != nil {
		t.Fatalf("unmarshal context payload: %v\n%s", err, string(out))
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

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	if err := runContext([]string{"--registry", registryPath, "--stage", "engineering", "--format", "yaml"}); err != nil {
		t.Fatalf("run context yaml: %v", err)
	}
	_ = w.Close()
	out, _ := io.ReadAll(r)
	text := string(out)
	if !strings.Contains(text, "schema_summary: root(req)") {
		t.Fatalf("expected schema_summary in yaml output, got %s", text)
	}
	if !strings.Contains(text, "parameters:") || !strings.Contains(text, "name: root") || !strings.Contains(text, "required: true") {
		t.Fatalf("expected parameters section in yaml output, got %s", text)
	}
}
