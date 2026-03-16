package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"toolcli/internal/registry"
	"toolcli/internal/telemetry"
	"toolcli/internal/types"
	"go.opentelemetry.io/otel/attribute"
)

func runValidate(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "validate")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	registryPath := fs.String("registry", "", "approved registry path")
	includeLocal := fs.Bool("include-local", false, "validate local overlay in addition to approved registry")
	if err := fs.Parse(args); err != nil {
		return err
	}

	reg, approvedPath, localPath, err := loadRegistry(ctx, *registryPath, *includeLocal)
	if err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("use.registry.approved_path", approvedPath),
		attribute.String("use.registry.local_path", localPath),
		attribute.Bool("use.registry.include_local", *includeLocal),
		attribute.Int("use.registry.system_count", len(reg.Systems)),
	)

	if err := registry.Validate(reg); err != nil {
		return err
	}
	fmt.Printf("registry valid: %d tool systems\n", len(reg.Systems))
	return nil
}

func runList(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "list")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	registryPath := fs.String("registry", "", "approved registry path")
	stage := fs.String("stage", "", "stage filter")
	tag := fs.String("tag", "", "tag filter")
	format := fs.String("format", "text", "output format: text|json|yaml")
	includeLocal := fs.Bool("include-local", false, "include local overlay tool systems")
	if err := fs.Parse(args); err != nil {
		return err
	}

	systems, err := filteredSystems(ctx, *registryPath, *includeLocal, *stage, *tag, "")
	if err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("use.stage", *stage),
		attribute.String("use.tag", *tag),
		attribute.String("use.format", *format),
		attribute.Int("use.results.count", len(systems)),
	)
	return writeSystemList(systems, *format)
}

func runDiscover(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "discover")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("discover", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	registryPath := fs.String("registry", "", "approved registry path")
	format := fs.String("format", "text", "output format: text|json|yaml")
	includeLocal := fs.Bool("include-local", false, "include local overlay tool systems")
	stage := fs.String("stage", "", "optional stage filter")
	if err := fs.Parse(args); err != nil {
		return err
	}
	query := strings.Join(fs.Args(), " ")
	if strings.TrimSpace(query) == "" {
		return fmt.Errorf("discover query is required")
	}
	systems, err := filteredSystems(ctx, *registryPath, *includeLocal, *stage, "", query)
	if err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("use.query", query),
		attribute.String("use.stage", *stage),
		attribute.String("use.format", *format),
		attribute.Int("use.results.count", len(systems)),
	)
	return writeSystemList(systems, *format)
}

func runContext(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "context")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("context", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	registryPath := fs.String("registry", "", "approved registry path")
	stage := fs.String("stage", "", "stage filter")
	query := fs.String("query", "", "optional intent filter")
	format := fs.String("format", "yaml", "output format: yaml|json")
	includeLocal := fs.Bool("include-local", false, "include local overlay tool systems")
	includeScoped := fs.Bool("include-scoped", false, "include scoped capabilities even when no query is provided")
	includePlanned := fs.Bool("include-planned", false, "include planned tool systems and capabilities in emitted context")
	if err := fs.Parse(args); err != nil {
		return err
	}
	systems, err := filteredSystems(ctx, *registryPath, *includeLocal, *stage, "", *query)
	if err != nil {
		return err
	}
	systems = registry.FilterForContext(systems, *includeScoped, *includePlanned, *query)
	commandSpan.SetAttributes(
		attribute.String("use.stage", *stage),
		attribute.String("use.query", *query),
		attribute.String("use.format", *format),
		attribute.Bool("use.include_scoped", *includeScoped),
		attribute.Bool("use.include_planned", *includePlanned),
		attribute.Int("use.results.count", len(systems)),
	)
	switch strings.ToLower(strings.TrimSpace(*format)) {
	case "json":
		payload := map[string]any{
			"tool_context": map[string]any{
				"stage":   *stage,
				"systems": summarizeContextSystems(systems),
			},
		}
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	case "yaml":
		fmt.Println("tool_context:")
		fmt.Printf("  stage: %s\n", strings.TrimSpace(*stage))
		fmt.Println("  systems:")
		for _, system := range summarizeContextSystems(systems) {
			fmt.Printf("    - id: %s\n", system["id"])
			fmt.Printf("      name: %s\n", system["name"])
			fmt.Printf("      description: %s\n", system["description"])
			fmt.Printf("      status: %s\n", system["status"])
			fmt.Printf("      guidance: %s\n", system["guidance"])
			fmt.Printf("      support_tier: %s\n", system["support_tier"])
			fmt.Printf("      capabilities_summary: %s\n", system["capabilities_summary"])
			capabilities := system["capabilities"].([]map[string]any)
			if len(capabilities) == 0 {
				fmt.Println("      capabilities: []")
				continue
			}
			fmt.Println("      capabilities:")
			for _, capability := range capabilities {
				fmt.Printf("        - id: %s\n", capability["id"])
				fmt.Printf("          name: %s\n", capability["name"])
				fmt.Printf("          status: %s\n", capability["status"])
				fmt.Printf("          availability: %s\n", capability["availability"])
				fmt.Printf("          guidance: %s\n", capability["guidance"])
				fmt.Printf("          schema_summary: %s\n", capability["schema_summary"])
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", *format)
	}
}

func runInspect(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "inspect")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	registryPath := fs.String("registry", "", "approved registry path")
	stage := fs.String("stage", "", "optional stage to evaluate capability fit against")
	format := fs.String("format", "text", "output format: text|json|yaml")
	includeLocal := fs.Bool("include-local", false, "include local overlay tool systems")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("inspect requires exactly one tool system id or query")
	}
	query := strings.TrimSpace(fs.Args()[0])
	reg, _, _, err := loadRegistry(ctx, *registryPath, *includeLocal)
	if err != nil {
		return err
	}
	if err := registry.Validate(reg); err != nil {
		return err
	}
	system, resolvedID, err := resolveInspectSystem(reg.Systems, query)
	if err != nil {
		return err
	}
	summary := summarizeInspectSystem(system, *stage)
	commandSpan.SetAttributes(
		attribute.String("use.tool.id", resolvedID),
		attribute.String("use.tool.query", query),
		attribute.String("use.stage", *stage),
		attribute.String("use.format", *format),
		attribute.String("use.tool.status", system.Status),
	)
	return writeInspectSummary(summary, *format)
}

func resolveInspectSystem(systems []types.ToolSystem, query string) (types.ToolSystem, string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return types.ToolSystem{}, "", fmt.Errorf("tool system query is required")
	}
	if system, ok := registry.FindByID(systems, query); ok {
		return system, system.ID, nil
	}
	matches := registry.Discover(systems, query)
	switch len(matches) {
	case 0:
		return types.ToolSystem{}, "", fmt.Errorf("unknown tool system id or query: %s", query)
	case 1:
		return matches[0], matches[0].ID, nil
	default:
		ids := make([]string, 0, min(len(matches), 5))
		for _, system := range matches[:min(len(matches), 5)] {
			ids = append(ids, system.ID)
		}
		return types.ToolSystem{}, "", fmt.Errorf("ambiguous tool system query %q; matches: %s", query, strings.Join(ids, ", "))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func loadRegistry(ctx context.Context, override string, includeLocal bool) (types.Registry, string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return types.Registry{}, "", "", err
	}
	approvedPath := strings.TrimSpace(override)
	localPath := ""
	if approvedPath == "" {
		approvedPath, localPath, err = registry.ResolvePaths(cwd)
		if err != nil {
			return types.Registry{}, "", "", err
		}
	}
	if approvedPath != "" {
		approvedPath = filepath.Clean(approvedPath)
	}
	if localPath != "" {
		localPath = filepath.Clean(localPath)
	}

	_, loadSpan := telemetry.StartSpan(ctx, "use.registry.load")
	defer func() {
		telemetry.EndSpan(loadSpan, err)
	}()
	loadSpan.SetAttributes(
		attribute.String("use.registry.approved_path", approvedPath),
		attribute.String("use.registry.local_path", localPath),
		attribute.Bool("use.registry.include_local", includeLocal),
	)

	reg, err := registry.LoadMerged(registry.LoadOptions{
		ApprovedPath: approvedPath,
		LocalPath:    localPath,
		IncludeLocal: includeLocal,
	})
	if err != nil {
		return types.Registry{}, approvedPath, localPath, err
	}
	return reg, approvedPath, localPath, nil
}

func filteredSystems(ctx context.Context, registryPath string, includeLocal bool, stage, tag, query string) ([]types.ToolSystem, error) {
	reg, _, _, err := loadRegistry(ctx, registryPath, includeLocal)
	if err != nil {
		return nil, err
	}
	if err := registry.Validate(reg); err != nil {
		return nil, err
	}
	systems := reg.Systems
	systems = registry.FilterByStage(systems, stage)
	systems = registry.FilterByTag(systems, tag)
	systems = registry.Discover(systems, query)
	return systems, nil
}

func writeSystemList(systems []types.ToolSystem, format string) error {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		data, err := json.MarshalIndent(summarizeSystems(systems), "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	case "yaml":
		for _, system := range summarizeSystems(systems) {
			fmt.Printf("- id: %s\n", system["id"])
			fmt.Printf("  name: %s\n", system["name"])
			fmt.Printf("  description: %s\n", system["description"])
			fmt.Printf("  status: %s\n", system["status"])
			fmt.Printf("  support_tier: %s\n", system["support_tier"])
			fmt.Printf("  capabilities_summary: %s\n", system["capabilities_summary"])
		}
		return nil
	case "text":
		for _, system := range systems {
			fmt.Printf("%s  %s\n", system.ID, system.Name)
			fmt.Printf("  %s\n", system.Description)
			fmt.Printf("  Status: %s\n", system.Status)
			fmt.Printf("  Capabilities: %s\n", capabilitySummary(system.Capabilities))
			fmt.Printf("  Tier: %s\n", system.SupportTier)
		}
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func summarizeSystems(systems []types.ToolSystem) []map[string]any {
	summary := make([]map[string]any, 0, len(systems))
	for _, system := range systems {
		summary = append(summary, map[string]any{
			"id":                   system.ID,
			"name":                 system.Name,
			"description":          system.Description,
			"status":               system.Status,
			"guidance":             system.Guidance,
			"complements":          append([]string{}, system.Complements...),
			"capabilities_summary": capabilitySummary(system.Capabilities),
			"support_tier":         system.SupportTier,
		})
	}
	return summary
}

func summarizeContextSystems(systems []types.ToolSystem) []map[string]any {
	summary := make([]map[string]any, 0, len(systems))
	for _, system := range systems {
		summary = append(summary, map[string]any{
			"id":                   system.ID,
			"name":                 system.Name,
			"description":          system.Description,
			"status":               system.Status,
			"guidance":             system.Guidance,
			"complements":          append([]string{}, system.Complements...),
			"capabilities_summary": capabilitySummary(system.Capabilities),
			"capabilities":         summarizeCapabilities(system.Capabilities),
			"support_tier":         system.SupportTier,
		})
	}
	return summary
}

func summarizeInspectSystem(system types.ToolSystem, stage string) map[string]any {
	return map[string]any{
		"id":                   system.ID,
		"name":                 system.Name,
		"description":          system.Description,
		"status":               system.Status,
		"guidance":             system.Guidance,
		"complements":          append([]string{}, system.Complements...),
		"support_tier":         system.SupportTier,
		"stage":                strings.TrimSpace(stage),
		"capabilities_summary": capabilitySummary(system.Capabilities),
		"capabilities":         summarizeInspectCapabilities(system.Capabilities, stage),
	}
}

func writeInspectSummary(summary map[string]any, format string) error {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		data, err := json.MarshalIndent(map[string]any{"tool_system": summary}, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	case "yaml":
		fmt.Println("tool_system:")
		writeInspectYAML(summary, "  ")
		return nil
	case "text":
		fmt.Printf("%s  %s\n", summary["id"], summary["name"])
		fmt.Printf("  %s\n", summary["description"])
		fmt.Printf("  Status: %s\n", summary["status"])
		if guidance, ok := summary["guidance"].(string); ok && guidance != "" {
			fmt.Printf("  Guidance: %s\n", guidance)
		}
		fmt.Printf("  Capabilities: %s\n", summary["capabilities_summary"])
		if complements, ok := summary["complements"].([]string); ok && len(complements) > 0 {
			fmt.Printf("  Complements: %s\n", strings.Join(complements, ", "))
		}
		capabilities := summary["capabilities"].([]map[string]any)
		for _, capability := range capabilities {
			fmt.Printf("  Capability: %s\n", capability["id"])
			fmt.Printf("    Name: %s\n", capability["name"])
			fmt.Printf("    Status: %s\n", capability["status"])
			fmt.Printf("    Availability: %s\n", capability["availability"])
			fmt.Printf("    Context: %s\n", capability["context_note"])
			fmt.Printf("    Schema: %s\n", capability["schema_summary"])
		}
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func writeInspectYAML(summary map[string]any, prefix string) {
	fmt.Printf("%sid: %s\n", prefix, summary["id"])
	fmt.Printf("%sname: %s\n", prefix, summary["name"])
	fmt.Printf("%sdescription: %s\n", prefix, summary["description"])
	fmt.Printf("%sstatus: %s\n", prefix, summary["status"])
	if guidance, ok := summary["guidance"].(string); ok && guidance != "" {
		fmt.Printf("%sguidance: %s\n", prefix, guidance)
	}
	fmt.Printf("%ssupport_tier: %s\n", prefix, summary["support_tier"])
	fmt.Printf("%scapabilities_summary: %s\n", prefix, summary["capabilities_summary"])
	if complements, ok := summary["complements"].([]string); ok && len(complements) > 0 {
		fmt.Printf("%scomplements: [%s]\n", prefix, strings.Join(complements, ", "))
	}
	capabilities := summary["capabilities"].([]map[string]any)
	if len(capabilities) == 0 {
		fmt.Printf("%scapabilities: []\n", prefix)
		return
	}
	fmt.Printf("%scapabilities:\n", prefix)
	for _, capability := range capabilities {
		fmt.Printf("%s  - id: %s\n", prefix, capability["id"])
		fmt.Printf("%s    name: %s\n", prefix, capability["name"])
		fmt.Printf("%s    status: %s\n", prefix, capability["status"])
		fmt.Printf("%s    availability: %s\n", prefix, capability["availability"])
		fmt.Printf("%s    context_note: %s\n", prefix, capability["context_note"])
		fmt.Printf("%s    schema_summary: %s\n", prefix, capability["schema_summary"])
	}
}

func summarizeCapabilities(capabilities []types.Capability) []map[string]any {
	summary := make([]map[string]any, 0, len(capabilities))
	for _, capability := range capabilities {
		summary = append(summary, map[string]any{
			"id":             capability.ID,
			"name":           capability.Name,
			"description":    capability.Description,
			"status":         capability.Status,
			"availability":   capability.Availability,
			"guidance":       capability.Guidance,
			"schema_summary": schemaSummary(capability.Schema),
		})
	}
	return summary
}

func summarizeInspectCapabilities(capabilities []types.Capability, stage string) []map[string]any {
	summary := make([]map[string]any, 0, len(capabilities))
	for _, capability := range capabilities {
		summary = append(summary, map[string]any{
			"id":             capability.ID,
			"name":           capability.Name,
			"status":         capability.Status,
			"availability":   capability.Availability,
			"context_note":   inspectContextNote(capability, stage),
			"schema_summary": schemaSummary(capability.Schema),
		})
	}
	return summary
}

func inspectContextNote(capability types.Capability, stage string) string {
	stage = strings.TrimSpace(stage)
	stageMatch := stage == "" || matchesStage(capability, stage)
	if capability.Status == registry.StatusPlanned {
		if stage != "" && !stageMatch {
			return "planned capability; not a fit for this stage without explicit future implementation"
		}
		return "planned capability; visible for planning but not active by default"
	}
	switch capability.Availability {
	case registry.AvailabilityRequired:
		if stage != "" && !stageMatch {
			return "required capability, but current stage does not match its affinity"
		}
		return "required capability; appropriate to include by default"
	case registry.AvailabilityScoped:
		if stage != "" && !stageMatch {
			return "scoped capability; not a fit for this stage without an explicit override"
		}
		return "scoped capability; discoverable and inspectable, but not included by default"
	default:
		if stage != "" && !stageMatch {
			return "default capability, but current stage does not match its affinity"
		}
		return "default capability; appropriate for normal context injection"
	}
}

func matchesStage(capability types.Capability, stage string) bool {
	stage = strings.TrimSpace(stage)
	if stage == "" || len(capability.StageAffinity) == 0 {
		return true
	}
	for _, candidate := range capability.StageAffinity {
		if candidate == stage {
			return true
		}
	}
	return false
}

func schemaSummary(fields []types.SchemaField) string {
	if len(fields) == 0 {
		return "(none)"
	}
	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		required := "opt"
		if field.Required {
			required = "req"
		}
		parts = append(parts, fmt.Sprintf("%s(%s)", field.Name, required))
	}
	return strings.Join(parts, ", ")
}

func capabilitySummary(capabilities []types.Capability) string {
	if len(capabilities) == 0 {
		return "(none)"
	}
	parts := make([]string, 0, len(capabilities))
	for _, capability := range capabilities {
		parts = append(parts, capability.ID)
	}
	return strings.Join(parts, ", ")
}
