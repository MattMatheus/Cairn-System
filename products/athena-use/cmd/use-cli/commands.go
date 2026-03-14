package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"athenause/internal/registry"
	"athenause/internal/telemetry"
	"athenause/internal/types"
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
		attribute.Int("use.registry.tool_count", len(reg.Tools)),
	)

	if err := registry.Validate(reg); err != nil {
		return err
	}
	fmt.Printf("registry valid: %d tools\n", len(reg.Tools))
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
	includeLocal := fs.Bool("include-local", false, "include local overlay tools")
	if err := fs.Parse(args); err != nil {
		return err
	}

	tools, err := filteredTools(ctx, *registryPath, *includeLocal, *stage, *tag, "")
	if err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("use.stage", *stage),
		attribute.String("use.tag", *tag),
		attribute.String("use.format", *format),
		attribute.Int("use.results.count", len(tools)),
	)
	return writeToolList(tools, *format)
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
	includeLocal := fs.Bool("include-local", false, "include local overlay tools")
	stage := fs.String("stage", "", "optional stage filter")
	if err := fs.Parse(args); err != nil {
		return err
	}
	query := strings.Join(fs.Args(), " ")
	if strings.TrimSpace(query) == "" {
		return fmt.Errorf("discover query is required")
	}
	tools, err := filteredTools(ctx, *registryPath, *includeLocal, *stage, "", query)
	if err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("use.query", query),
		attribute.String("use.stage", *stage),
		attribute.String("use.format", *format),
		attribute.Int("use.results.count", len(tools)),
	)
	return writeToolList(tools, *format)
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
	includeLocal := fs.Bool("include-local", false, "include local overlay tools")
	includeScoped := fs.Bool("include-scoped", false, "include scoped tools even when no query is provided")
	includePlanned := fs.Bool("include-planned", false, "include planned tools in emitted context")
	if err := fs.Parse(args); err != nil {
		return err
	}
	tools, err := filteredTools(ctx, *registryPath, *includeLocal, *stage, "", *query)
	if err != nil {
		return err
	}
	tools = registry.FilterForContext(tools, *includeScoped, *includePlanned, *query)
	commandSpan.SetAttributes(
		attribute.String("use.stage", *stage),
		attribute.String("use.query", *query),
		attribute.String("use.format", *format),
		attribute.Bool("use.include_scoped", *includeScoped),
		attribute.Bool("use.include_planned", *includePlanned),
		attribute.Int("use.results.count", len(tools)),
	)
	switch strings.ToLower(strings.TrimSpace(*format)) {
	case "json":
		payload := map[string]any{
			"tool_context": map[string]any{
				"stage": *stage,
				"tools": summarizeContextTools(tools),
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
		fmt.Println("  tools:")
		for _, tool := range summarizeContextTools(tools) {
			fmt.Printf("    - id: %s\n", tool["id"])
			fmt.Printf("      name: %s\n", tool["name"])
			fmt.Printf("      description: %s\n", tool["description"])
			fmt.Printf("      status: %s\n", tool["status"])
			fmt.Printf("      availability: %s\n", tool["availability"])
			fmt.Printf("      guidance: %s\n", tool["guidance"])
			fmt.Printf("      call_type: %s\n", tool["call_type"])
			fmt.Printf("      credential_ref: %s\n", tool["credential_ref"])
			fmt.Printf("      support_tier: %s\n", tool["support_tier"])
			fmt.Printf("      schema_summary: %s\n", tool["schema_summary"])
			fields := tool["parameters"].([]map[string]any)
			if len(fields) == 0 {
				fmt.Println("      parameters: []")
				continue
			}
			fmt.Println("      parameters:")
			for _, field := range fields {
				fmt.Printf("        - name: %s\n", field["name"])
				fmt.Printf("          type: %s\n", field["type"])
				fmt.Printf("          required: %t\n", field["required"])
				if description, ok := field["description"].(string); ok && description != "" {
					fmt.Printf("          description: %s\n", description)
				}
				if enumValues, ok := field["enum"].([]string); ok && len(enumValues) > 0 {
					fmt.Printf("          enum: [%s]\n", strings.Join(enumValues, ", "))
				}
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
	stage := fs.String("stage", "", "optional stage to evaluate fit against")
	format := fs.String("format", "text", "output format: text|json|yaml")
	includeLocal := fs.Bool("include-local", false, "include local overlay tools")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("inspect requires exactly one tool id")
	}
	toolID := strings.TrimSpace(fs.Args()[0])
	reg, _, _, err := loadRegistry(ctx, *registryPath, *includeLocal)
	if err != nil {
		return err
	}
	if err := registry.Validate(reg); err != nil {
		return err
	}
	tool, ok := registry.FindByID(reg.Tools, toolID)
	if !ok {
		return fmt.Errorf("unknown tool id: %s", toolID)
	}
	summary := summarizeInspectTool(tool, *stage)
	commandSpan.SetAttributes(
		attribute.String("use.tool.id", toolID),
		attribute.String("use.stage", *stage),
		attribute.String("use.format", *format),
		attribute.String("use.tool.status", tool.Status),
		attribute.String("use.tool.availability", tool.Availability),
	)
	return writeInspectSummary(summary, *format)
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

func filteredTools(ctx context.Context, registryPath string, includeLocal bool, stage, tag, query string) ([]types.Tool, error) {
	reg, _, _, err := loadRegistry(ctx, registryPath, includeLocal)
	if err != nil {
		return nil, err
	}
	if err := registry.Validate(reg); err != nil {
		return nil, err
	}
	tools := reg.Tools
	tools = registry.FilterByStage(tools, stage)
	tools = registry.FilterByTag(tools, tag)
	tools = registry.Discover(tools, query)
	return tools, nil
}

func writeToolList(tools []types.Tool, format string) error {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		data, err := json.MarshalIndent(summarizeTools(tools), "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	case "yaml":
		for _, tool := range summarizeTools(tools) {
			fmt.Printf("- id: %s\n", tool["id"])
			fmt.Printf("  name: %s\n", tool["name"])
			fmt.Printf("  description: %s\n", tool["description"])
			fmt.Printf("  call_type: %s\n", tool["call_type"])
			fmt.Printf("  credential_ref: %s\n", tool["credential_ref"])
			fmt.Printf("  support_tier: %s\n", tool["support_tier"])
		}
		return nil
	case "text":
		for _, tool := range tools {
			fmt.Printf("%s  %s\n", tool.ID, tool.Name)
			fmt.Printf("  %s\n", tool.Description)
			fmt.Printf("  Status: %s\n", tool.Status)
			fmt.Printf("  Availability: %s\n", tool.Availability)
			fmt.Printf("  Credential: %s\n", tool.CredentialRef)
			fmt.Printf("  Schema: %s\n", schemaSummary(tool.Schema))
			fmt.Printf("  Tier: %s\n", tool.SupportTier)
		}
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func summarizeTools(tools []types.Tool) []map[string]any {
	summary := make([]map[string]any, 0, len(tools))
	for _, tool := range tools {
		summary = append(summary, map[string]any{
			"id":             tool.ID,
			"name":           tool.Name,
			"description":    tool.Description,
			"status":         tool.Status,
			"availability":   tool.Availability,
			"guidance":       tool.Guidance,
			"complements":    append([]string{}, tool.Complements...),
			"schema":         schemaSummary(tool.Schema),
			"credential_ref": tool.CredentialRef,
			"call_type":      tool.Call.Type,
			"support_tier":   tool.SupportTier,
		})
	}
	return summary
}

func summarizeContextTools(tools []types.Tool) []map[string]any {
	summary := make([]map[string]any, 0, len(tools))
	for _, tool := range tools {
		summary = append(summary, map[string]any{
			"id":             tool.ID,
			"name":           tool.Name,
			"description":    tool.Description,
			"status":         tool.Status,
			"availability":   tool.Availability,
			"guidance":       tool.Guidance,
			"complements":    append([]string{}, tool.Complements...),
			"schema":         schemaSummary(tool.Schema),
			"schema_summary": schemaSummary(tool.Schema),
			"parameters":     summarizeSchemaFields(tool.Schema),
			"credential_ref": tool.CredentialRef,
			"call_type":      tool.Call.Type,
			"support_tier":   tool.SupportTier,
		})
	}
	return summary
}

func summarizeInspectTool(tool types.Tool, stage string) map[string]any {
	stage = strings.TrimSpace(stage)
	stageMatch := stage == "" || matchesStage(tool, stage)
	statusActive := tool.Status != registry.StatusPlanned
	contextEligible := tool.Availability != registry.AvailabilityScoped || stage == "" || stageMatch
	if tool.Availability == registry.AvailabilityScoped && stage == "" {
		contextEligible = false
	}
	if tool.Availability == registry.AvailabilityScoped && stage != "" {
		contextEligible = false
	}
	if !statusActive {
		contextEligible = false
	}
	return map[string]any{
		"id":                 tool.ID,
		"name":               tool.Name,
		"description":        tool.Description,
		"status":             tool.Status,
		"availability":       tool.Availability,
		"guidance":           tool.Guidance,
		"complements":        append([]string{}, tool.Complements...),
		"support_tier":       tool.SupportTier,
		"stage_affinity":     append([]string{}, tool.StageAffinity...),
		"stage":              stage,
		"stage_match":        stageMatch,
		"default_in_context": statusActive && tool.Availability != registry.AvailabilityScoped && stageMatch,
		"context_note":       inspectContextNote(tool, stage, stageMatch),
		"call_type":          tool.Call.Type,
		"call_method":        tool.Call.Method,
		"call_url":           tool.Call.URL,
		"call_command":       tool.Call.Command,
		"credential_ref":     tool.CredentialRef,
		"schema_summary":     schemaSummary(tool.Schema),
		"parameters":         summarizeSchemaFields(tool.Schema),
		"context_eligible":   contextEligible,
	}
}

func writeInspectSummary(summary map[string]any, format string) error {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		data, err := json.MarshalIndent(map[string]any{"tool": summary}, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	case "yaml":
		fmt.Println("tool:")
		writeInspectYAML(summary, "  ")
		return nil
	case "text":
		fmt.Printf("%s  %s\n", summary["id"], summary["name"])
		fmt.Printf("  %s\n", summary["description"])
		fmt.Printf("  Status: %s\n", summary["status"])
		fmt.Printf("  Availability: %s\n", summary["availability"])
		fmt.Printf("  Context: %s\n", summary["context_note"])
		if guidance, ok := summary["guidance"].(string); ok && guidance != "" {
			fmt.Printf("  Guidance: %s\n", guidance)
		}
		if stage, ok := summary["stage"].(string); ok && stage != "" {
			fmt.Printf("  Stage: %s\n", stage)
			fmt.Printf("  Stage Match: %t\n", summary["stage_match"])
		}
		fmt.Printf("  Call Type: %s\n", summary["call_type"])
		fmt.Printf("  Credential: %s\n", summary["credential_ref"])
		fmt.Printf("  Schema: %s\n", summary["schema_summary"])
		if complements, ok := summary["complements"].([]string); ok && len(complements) > 0 {
			fmt.Printf("  Complements: %s\n", strings.Join(complements, ", "))
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
	fmt.Printf("%savailability: %s\n", prefix, summary["availability"])
	fmt.Printf("%scontext_note: %s\n", prefix, summary["context_note"])
	if guidance, ok := summary["guidance"].(string); ok && guidance != "" {
		fmt.Printf("%sguidance: %s\n", prefix, guidance)
	}
	fmt.Printf("%ssupport_tier: %s\n", prefix, summary["support_tier"])
	fmt.Printf("%scall_type: %s\n", prefix, summary["call_type"])
	if method, ok := summary["call_method"].(string); ok && method != "" {
		fmt.Printf("%scall_method: %s\n", prefix, method)
	}
	if url, ok := summary["call_url"].(string); ok && url != "" {
		fmt.Printf("%scall_url: %s\n", prefix, url)
	}
	if command, ok := summary["call_command"].(string); ok && command != "" {
		fmt.Printf("%scall_command: %s\n", prefix, command)
	}
	fmt.Printf("%scredential_ref: %s\n", prefix, summary["credential_ref"])
	fmt.Printf("%sschema_summary: %s\n", prefix, summary["schema_summary"])
	fmt.Printf("%sdefault_in_context: %t\n", prefix, summary["default_in_context"])
	if stage, ok := summary["stage"].(string); ok && stage != "" {
		fmt.Printf("%sstage: %s\n", prefix, stage)
		fmt.Printf("%sstage_match: %t\n", prefix, summary["stage_match"])
	}
	if stages, ok := summary["stage_affinity"].([]string); ok && len(stages) > 0 {
		fmt.Printf("%sstage_affinity: [%s]\n", prefix, strings.Join(stages, ", "))
	}
	if complements, ok := summary["complements"].([]string); ok && len(complements) > 0 {
		fmt.Printf("%scomplements: [%s]\n", prefix, strings.Join(complements, ", "))
	}
	fields := summary["parameters"].([]map[string]any)
	if len(fields) == 0 {
		fmt.Printf("%sparameters: []\n", prefix)
		return
	}
	fmt.Printf("%sparameters:\n", prefix)
	for _, field := range fields {
		fmt.Printf("%s  - name: %s\n", prefix, field["name"])
		fmt.Printf("%s    type: %s\n", prefix, field["type"])
		fmt.Printf("%s    required: %t\n", prefix, field["required"])
		if description, ok := field["description"].(string); ok && description != "" {
			fmt.Printf("%s    description: %s\n", prefix, description)
		}
		if enumValues, ok := field["enum"].([]string); ok && len(enumValues) > 0 {
			fmt.Printf("%s    enum: [%s]\n", prefix, strings.Join(enumValues, ", "))
		}
	}
}

func inspectContextNote(tool types.Tool, stage string, stageMatch bool) string {
	if tool.Status == registry.StatusPlanned {
		if stage != "" && !stageMatch {
			return "planned tool; not active and not a fit for this stage without explicit future implementation"
		}
		return "planned tool; discoverable for planning, but not included in active context by default"
	}
	switch tool.Availability {
	case registry.AvailabilityRequired:
		if stage != "" && !stageMatch {
			return "required tool, but current stage does not match its affinity"
		}
		return "required tool; appropriate to include by default"
	case registry.AvailabilityScoped:
		if stage != "" && !stageMatch {
			return "scoped tool; not a fit for this stage without an explicit override"
		}
		return "scoped tool; discoverable and inspectable, but not included by default"
	default:
		if stage != "" && !stageMatch {
			return "default tool, but current stage does not match its affinity"
		}
		return "default tool; appropriate for normal context injection"
	}
}

func matchesStage(tool types.Tool, stage string) bool {
	stage = strings.TrimSpace(stage)
	if stage == "" || len(tool.StageAffinity) == 0 {
		return true
	}
	for _, candidate := range tool.StageAffinity {
		if candidate == stage {
			return true
		}
	}
	return false
}

func summarizeSchemaFields(fields []types.SchemaField) []map[string]any {
	if len(fields) == 0 {
		return []map[string]any{}
	}
	summary := make([]map[string]any, 0, len(fields))
	for _, field := range fields {
		summary = append(summary, map[string]any{
			"name":        field.Name,
			"type":        field.Type,
			"required":    field.Required,
			"description": field.Description,
			"enum":        append([]string{}, field.Enum...),
		})
	}
	return summary
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
