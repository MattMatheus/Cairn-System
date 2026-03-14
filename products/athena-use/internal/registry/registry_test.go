package registry

import (
	"os"
	"path/filepath"
	"strings"
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
	if tool.Status != StatusActive {
		t.Fatalf("expected active status, got %q", tool.Status)
	}
	if tool.Availability != AvailabilityDefault {
		t.Fatalf("expected default availability, got %q", tool.Availability)
	}
	if len(tool.Schema) != 2 || !tool.Schema[0].Required {
		t.Fatalf("unexpected schema: %+v", tool.Schema)
	}
}

func TestFilterForContextExcludesScopedToolsByDefault(t *testing.T) {
	tools := []types.Tool{
		{ID: "obsidian", Status: StatusActive, Availability: AvailabilityRequired},
		{ID: "platform.check", Status: StatusActive, Availability: AvailabilityDefault},
		{ID: "gitnexus", Status: StatusActive, Availability: AvailabilityScoped},
	}
	got := FilterForContext(tools, false, false, "")
	if len(got) != 2 {
		t.Fatalf("expected 2 tools after default context filter, got %d", len(got))
	}
	for _, tool := range got {
		if tool.ID == "gitnexus" {
			t.Fatalf("scoped tool should not appear in default context: %+v", got)
		}
	}
}

func TestFilterForContextKeepsScopedToolsWhenQueryPresent(t *testing.T) {
	tools := []types.Tool{
		{ID: "gitnexus", Status: StatusActive, Availability: AvailabilityScoped},
	}
	got := FilterForContext(tools, false, false, "code archaeology")
	if len(got) != 1 || got[0].ID != "gitnexus" {
		t.Fatalf("expected query to preserve scoped tool, got %+v", got)
	}
}

func TestFilterForContextExcludesPlannedToolsByDefault(t *testing.T) {
	tools := []types.Tool{
		{ID: "obsidian", Status: StatusActive, Availability: AvailabilityRequired},
		{ID: "athena.intake.normalize_source", Status: StatusPlanned, Availability: AvailabilityScoped},
	}
	got := FilterForContext(tools, true, false, "intake markdown")
	if len(got) != 1 || got[0].ID != "obsidian" {
		t.Fatalf("expected planned tool to remain excluded by default, got %+v", got)
	}
}

func TestFilterForContextCanIncludePlannedToolsExplicitly(t *testing.T) {
	tools := []types.Tool{
		{ID: "athena.intake.normalize_source", Status: StatusPlanned, Availability: AvailabilityScoped},
	}
	got := FilterForContext(tools, true, true, "intake markdown")
	if len(got) != 1 || got[0].ID != "athena.intake.normalize_source" {
		t.Fatalf("expected planned tool when explicitly included, got %+v", got)
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

	toolsByID := map[string]types.Tool{}
	for _, tool := range reg.Tools {
		toolsByID[tool.ID] = tool
	}

	verifyHealth, ok := toolsByID["athena.memory.verify_health"]
	if !ok {
		t.Fatalf("approved registry missing athena.memory.verify_health")
	}
	if len(verifyHealth.Schema) != 3 || strings.TrimSpace(verifyHealth.Schema[0].Description) == "" {
		t.Fatalf("expected verify_health schema descriptions, got %+v", verifyHealth.Schema)
	}

	taskMetadata, ok := toolsByID["athena.workspace.validate_task_metadata"]
	if !ok {
		t.Fatalf("approved registry missing athena.workspace.validate_task_metadata")
	}
	if len(taskMetadata.Schema) != 2 {
		t.Fatalf("expected task metadata schema entries, got %+v", taskMetadata.Schema)
	}
	if got := strings.Join(taskMetadata.Schema[0].Enum, ","); got != "changed,all" {
		t.Fatalf("expected mode enum changed,all, got %q", got)
	}

	intakeTool, ok := toolsByID["athena.intake.normalize_source"]
	if !ok || intakeTool.Status != StatusActive || intakeTool.Call.Type != "exec" {
		t.Fatalf("expected planned intake tool in approved registry, got %+v", intakeTool)
	}
	if len(intakeTool.Schema) != 4 {
		t.Fatalf("expected staged intake tool schema entries, got %+v", intakeTool.Schema)
	}
	intakeInspect, ok := toolsByID["athena.intake.inspect_source"]
	if !ok || intakeInspect.Status != StatusActive || intakeInspect.Call.Type != "exec" {
		t.Fatalf("expected active intake inspect tool in approved registry, got %+v", intakeInspect)
	}
}
