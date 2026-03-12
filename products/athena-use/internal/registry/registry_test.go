package registry

import (
	"os"
	"path/filepath"
	"testing"

	"athenause/internal/types"
)

func TestLoadFileParsesRegistry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "registry.yaml")
	content := `version: 1
tools:
  - id: github.create_pr
    name: Create Pull Request
    description: Opens a pull request
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
      - name: repo
        type: string
        required: true
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	reg, err := LoadFile(path, SupportTierApproved)
	if err != nil {
		t.Fatalf("load registry: %v", err)
	}
	if reg.Version != 1 {
		t.Fatalf("expected version 1, got %d", reg.Version)
	}
	if len(reg.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(reg.Tools))
	}
	tool := reg.Tools[0]
	if tool.ID != "github.create_pr" || tool.SupportTier != SupportTierApproved {
		t.Fatalf("unexpected tool: %+v", tool)
	}
	if tool.Call.Type != "http" || tool.Call.Method != "POST" {
		t.Fatalf("unexpected call contract: %+v", tool.Call)
	}
	if len(tool.Schema) != 2 || !tool.Schema[0].Required {
		t.Fatalf("unexpected schema: %+v", tool.Schema)
	}
}

func TestDiscoverScoresMatches(t *testing.T) {
	tools := []types.Tool{
		{
			ID:          "github.create_pr",
			Name:        "Create Pull Request",
			Description: "Open a pull request on GitHub",
		},
		{
			ID:          "internal.deploy_service",
			Name:        "Deploy Service",
			Description: "Deploy a service to staging",
		},
	}
	got := Discover(tools, "pull request")
	if len(got) != 1 || got[0].ID != "github.create_pr" {
		t.Fatalf("unexpected discover result: %+v", got)
	}
}

func TestApprovedRegistryFileIsValid(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	approvedPath, _, err := ResolvePaths(wd)
	if err != nil {
		t.Fatalf("resolve paths: %v", err)
	}

	reg, err := LoadFile(approvedPath, SupportTierApproved)
	if err != nil {
		t.Fatalf("load approved registry: %v", err)
	}
	if err := Validate(reg); err != nil {
		t.Fatalf("validate approved registry: %v", err)
	}
	if len(reg.Tools) < 5 {
		t.Fatalf("expected at least 5 approved tools, got %d", len(reg.Tools))
	}
}
