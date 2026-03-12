package main

import (
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
