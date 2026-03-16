package registry

import (
	"os"
	"path/filepath"
	"testing"

	"toolcli/internal/types"
)

func TestLoadFileParsesToolSystems(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "registry.yaml")
	content := `version: 1
tools:
  - id: cairn
    name: Cairn
    description: Cairn internal tool system
    status: active
    tags: [cairn, internal]
    guidance: Use for supported Cairn workflows
    capabilities:
      - id: intake.inspect_source
        name: Inspect Intake Source
        description: Inspect a source before staging
        status: active
        availability: scoped
        tags: [intake]
        stage_affinity: [pm]
        guidance: Use before staging
        call:
          type: exec
          command: intake-cli inspect
        schema:
          - name: source
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
	if len(reg.Systems) != 1 {
		t.Fatalf("expected 1 system, got %d", len(reg.Systems))
	}
	system := reg.Systems[0]
	if system.ID != "cairn" || system.SupportTier != SupportTierApproved {
		t.Fatalf("unexpected system: %+v", system)
	}
	if len(system.Capabilities) != 1 || system.Capabilities[0].ID != "intake.inspect_source" {
		t.Fatalf("unexpected capabilities: %+v", system.Capabilities)
	}
}

func TestFilterForContextExcludesScopedCapabilitiesByDefault(t *testing.T) {
	systems := []types.ToolSystem{
		{
			ID:     "cairn",
			Status: StatusActive,
			Capabilities: []types.Capability{
				{ID: "smoke", Status: StatusActive, Availability: AvailabilityDefault},
				{ID: "inspect", Status: StatusActive, Availability: AvailabilityScoped},
			},
		},
	}
	got := FilterForContext(systems, false, false, "")
	if len(got) != 1 || len(got[0].Capabilities) != 1 || got[0].Capabilities[0].ID != "smoke" {
		t.Fatalf("unexpected filtered context: %+v", got)
	}
}

func TestFilterByStageKeepsOnlyMatchingCapabilities(t *testing.T) {
	systems := []types.ToolSystem{
		{
			ID: "cairn",
			Capabilities: []types.Capability{
				{ID: "smoke", StageAffinity: []string{"engineering"}},
				{ID: "inspect", StageAffinity: []string{"pm"}},
			},
		},
	}
	got := FilterByStage(systems, "pm")
	if len(got) != 1 || len(got[0].Capabilities) != 1 || got[0].Capabilities[0].ID != "inspect" {
		t.Fatalf("unexpected stage filter result: %+v", got)
	}
}

func TestDiscoverMatchesCapabilityMetadata(t *testing.T) {
	systems := []types.ToolSystem{
		{
			ID:          "gitnexus",
			Name:        "GitNexus",
			Description: "Repository graph analysis system",
			Capabilities: []types.Capability{
				{
					ID:          "code_graph",
					Name:        "Code Graph",
					Description: "Analyze repository structure",
					Guidance:    "Use for code archaeology and impact analysis",
					Call: types.ToolCall{
						Type:    "exec",
						Command: "gitnexus query",
					},
				},
			},
		},
	}
	got := Discover(systems, "archaeology")
	if len(got) != 1 || got[0].ID != "gitnexus" {
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
	if len(reg.Systems) != 4 {
		t.Fatalf("expected 4 tool systems, got %d", len(reg.Systems))
	}
}
