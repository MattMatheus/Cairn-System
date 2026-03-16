package registry

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"toolcli/internal/types"

	"gopkg.in/yaml.v3"
)

const (
	SupportTierApproved  = "approved"
	SupportTierLocal     = "local"
	AvailabilityRequired = "required"
	AvailabilityDefault  = "default"
	AvailabilityScoped   = "scoped"
	StatusActive         = "active"
	StatusPlanned        = "planned"
)

type LoadOptions struct {
	ApprovedPath string
	LocalPath    string
	IncludeLocal bool
}

func ResolvePaths(cwd string) (approvedPath, localPath string, err error) {
	if override := strings.TrimSpace(os.Getenv("CAIRN_TOOL_REGISTRY")); override != "" {
		return override, "", nil
	}

	root, err := findRepoRoot(cwd)
	if err != nil {
		return "", "", err
	}
	approvedPath = filepath.Join(root, "products", "tool-cli", "registry", "approved-tools.yaml")
	cairnHome := strings.TrimSpace(os.Getenv("CAIRN_HOME"))
	if cairnHome == "" {
		cairnHome = filepath.Join(root, ".cairn")
	}
	localPath = filepath.Join(cairnHome, "tools", "registry.yaml")
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
	var reg types.Registry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return types.Registry{}, fmt.Errorf("%s: %w", path, err)
	}
	for i := range reg.Systems {
		reg.Systems[i].Status = normalizeStatus(reg.Systems[i].Status)
		reg.Systems[i].SupportTier = tier
		for j := range reg.Systems[i].Capabilities {
			reg.Systems[i].Capabilities[j].Status = normalizeStatus(reg.Systems[i].Capabilities[j].Status)
			reg.Systems[i].Capabilities[j].Availability = normalizeAvailability(reg.Systems[i].Capabilities[j].Availability)
		}
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
	validAvailability := map[string]struct{}{
		AvailabilityRequired: {},
		AvailabilityDefault:  {},
		AvailabilityScoped:   {},
	}
	validStatus := map[string]struct{}{
		StatusActive:  {},
		StatusPlanned: {},
	}
	for _, system := range reg.Systems {
		if strings.TrimSpace(system.ID) == "" {
			return errors.New("tool system id is required")
		}
		if _, ok := seen[system.ID]; ok {
			return fmt.Errorf("duplicate tool system id: %s", system.ID)
		}
		seen[system.ID] = struct{}{}
		if strings.TrimSpace(system.Name) == "" {
			return fmt.Errorf("tool system %s missing name", system.ID)
		}
		if strings.TrimSpace(system.Description) == "" {
			return fmt.Errorf("tool system %s missing description", system.ID)
		}
		if _, ok := validStatus[normalizeStatus(system.Status)]; !ok {
			return fmt.Errorf("tool system %s has invalid status value: %s", system.ID, system.Status)
		}
		if len(system.Capabilities) == 0 {
			return fmt.Errorf("tool system %s missing capabilities", system.ID)
		}
		capSeen := map[string]struct{}{}
		for _, capability := range system.Capabilities {
			if strings.TrimSpace(capability.ID) == "" {
				return fmt.Errorf("tool system %s has capability missing id", system.ID)
			}
			if _, ok := capSeen[capability.ID]; ok {
				return fmt.Errorf("tool system %s has duplicate capability id: %s", system.ID, capability.ID)
			}
			capSeen[capability.ID] = struct{}{}
			if strings.TrimSpace(capability.Name) == "" {
				return fmt.Errorf("tool system %s capability %s missing name", system.ID, capability.ID)
			}
			if strings.TrimSpace(capability.Description) == "" {
				return fmt.Errorf("tool system %s capability %s missing description", system.ID, capability.ID)
			}
			if strings.TrimSpace(capability.Call.Type) == "" {
				return fmt.Errorf("tool system %s capability %s missing call.type", system.ID, capability.ID)
			}
			if _, ok := validAvailability[normalizeAvailability(capability.Availability)]; !ok {
				return fmt.Errorf("tool system %s capability %s has invalid availability value: %s", system.ID, capability.ID, capability.Availability)
			}
			if _, ok := validStatus[normalizeStatus(capability.Status)]; !ok {
				return fmt.Errorf("tool system %s capability %s has invalid status value: %s", system.ID, capability.ID, capability.Status)
			}
			for _, stage := range capability.StageAffinity {
				if _, ok := validStages[stage]; !ok {
					return fmt.Errorf("tool system %s capability %s has invalid stage_affinity value: %s", system.ID, capability.ID, stage)
				}
			}
			for _, field := range capability.Schema {
				if strings.TrimSpace(field.Name) == "" {
					return fmt.Errorf("tool system %s capability %s has schema entry missing name", system.ID, capability.ID)
				}
				if strings.TrimSpace(field.Type) == "" {
					return fmt.Errorf("tool system %s capability %s schema field %s missing type", system.ID, capability.ID, field.Name)
				}
			}
		}
	}
	return nil
}

func FilterByStage(systems []types.ToolSystem, stage string) []types.ToolSystem {
	stage = strings.TrimSpace(stage)
	if stage == "" {
		return cloneSystems(systems)
	}
	filtered := make([]types.ToolSystem, 0, len(systems))
	for _, system := range systems {
		matched := filterCapabilitiesByStage(system.Capabilities, stage)
		if len(matched) == 0 {
			continue
		}
		system.Capabilities = matched
		filtered = append(filtered, system)
	}
	return filtered
}

func FilterByTag(systems []types.ToolSystem, tag string) []types.ToolSystem {
	tag = strings.TrimSpace(strings.ToLower(tag))
	if tag == "" {
		return cloneSystems(systems)
	}
	filtered := make([]types.ToolSystem, 0, len(systems))
	for _, system := range systems {
		if containsFold(system.Tags, tag) || capabilityTagMatch(system.Capabilities, tag) {
			filtered = append(filtered, system)
		}
	}
	return filtered
}

func FilterForContext(systems []types.ToolSystem, includeScoped, includePlanned bool, query string) []types.ToolSystem {
	filtered := make([]types.ToolSystem, 0, len(systems))
	for _, system := range systems {
		if normalizeStatus(system.Status) == StatusPlanned && !includePlanned {
			continue
		}
		kept := make([]types.Capability, 0, len(system.Capabilities))
		for _, capability := range system.Capabilities {
			if normalizeStatus(capability.Status) == StatusPlanned && !includePlanned {
				continue
			}
			if !includeScoped && strings.TrimSpace(query) == "" && normalizeAvailability(capability.Availability) == AvailabilityScoped {
				continue
			}
			kept = append(kept, capability)
		}
		if len(kept) == 0 {
			continue
		}
		system.Capabilities = kept
		filtered = append(filtered, system)
	}
	return filtered
}

func FindByID(systems []types.ToolSystem, id string) (types.ToolSystem, bool) {
	id = strings.TrimSpace(id)
	for _, system := range systems {
		if system.ID == id {
			return system, true
		}
	}
	return types.ToolSystem{}, false
}

func Discover(systems []types.ToolSystem, query string) []types.ToolSystem {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return cloneSystems(systems)
	}
	terms := strings.Fields(query)
	type scored struct {
		system types.ToolSystem
		score  int
	}
	var matches []scored
	for _, system := range systems {
		score := 0
		for _, term := range terms {
			score += matchSystemTermScore(system, term)
		}
		if score > 0 {
			matches = append(matches, scored{system: system, score: score})
		}
	}
	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].score == matches[j].score {
			return matches[i].system.ID < matches[j].system.ID
		}
		return matches[i].score > matches[j].score
	})
	out := make([]types.ToolSystem, 0, len(matches))
	for _, match := range matches {
		out = append(out, match.system)
	}
	return out
}

func normalizeStatus(value string) string {
	if strings.TrimSpace(value) == "" {
		return StatusActive
	}
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeAvailability(value string) string {
	if strings.TrimSpace(value) == "" {
		return AvailabilityDefault
	}
	return strings.ToLower(strings.TrimSpace(value))
}

func mergeRegistries(approved, local types.Registry) (types.Registry, error) {
	merged := types.Registry{Version: 1}
	merged.Systems = append(merged.Systems, approved.Systems...)
	seen := map[string]struct{}{}
	for _, system := range approved.Systems {
		seen[system.ID] = struct{}{}
	}
	for _, system := range local.Systems {
		if _, ok := seen[system.ID]; ok {
			return types.Registry{}, fmt.Errorf("duplicate tool system id across approved and local registries: %s", system.ID)
		}
		merged.Systems = append(merged.Systems, system)
	}
	return merged, nil
}

func findRepoRoot(start string) (string, error) {
	dir := start
	for {
		candidate := filepath.Join(dir, "products", "tool-cli", "registry", "approved-tools.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("unable to resolve Cairn repo root")
		}
		dir = parent
	}
}

func cloneSystems(systems []types.ToolSystem) []types.ToolSystem {
	cloned := make([]types.ToolSystem, 0, len(systems))
	for _, system := range systems {
		system.Capabilities = append([]types.Capability(nil), system.Capabilities...)
		cloned = append(cloned, system)
	}
	return cloned
}

func filterCapabilitiesByStage(capabilities []types.Capability, stage string) []types.Capability {
	filtered := make([]types.Capability, 0, len(capabilities))
	for _, capability := range capabilities {
		if len(capability.StageAffinity) == 0 || containsExact(capability.StageAffinity, stage) {
			filtered = append(filtered, capability)
		}
	}
	return filtered
}

func capabilityTagMatch(capabilities []types.Capability, tag string) bool {
	for _, capability := range capabilities {
		if containsFold(capability.Tags, tag) {
			return true
		}
	}
	return false
}

func containsFold(values []string, target string) bool {
	for _, value := range values {
		if strings.ToLower(value) == target {
			return true
		}
	}
	return false
}

func containsExact(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func matchSystemTermScore(system types.ToolSystem, term string) int {
	score := 0
	score += weightedContains(system.ID, term, 10)
	score += weightedContains(system.Name, term, 8)
	score += weightedContains(system.Description, term, 4)
	score += weightedContains(system.Guidance, term, 3)
	score += weightedContains(strings.Join(system.Tags, " "), term, 6)
	score += weightedContains(strings.Join(system.Complements, " "), term, 2)
	for _, capability := range system.Capabilities {
		score += weightedContains(capability.ID, term, 5)
		score += weightedContains(capability.Name, term, 4)
		score += weightedContains(capability.Description, term, 3)
		score += weightedContains(capability.Guidance, term, 2)
		score += weightedContains(strings.Join(capability.Tags, " "), term, 2)
		score += weightedContains(strings.Join(capability.StageAffinity, " "), term, 1)
		score += weightedContains(capability.Call.Type, term, 1)
		score += weightedContains(capability.Call.Method, term, 1)
		score += weightedContains(capability.Call.URL, term, 1)
		score += weightedContains(capability.Call.Command, term, 2)
		for _, field := range capability.Schema {
			score += weightedContains(field.Name, term, 2)
			score += weightedContains(field.Type, term, 1)
			score += weightedContains(field.Description, term, 2)
			score += weightedContains(strings.Join(field.Enum, " "), term, 1)
		}
	}
	return score
}

func weightedContains(value, term string, weight int) int {
	if weight <= 0 {
		return 0
	}
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" || term == "" {
		return 0
	}
	if value == term {
		return weight * 2
	}
	if strings.Contains(value, term) {
		return weight
	}
	return 0
}
