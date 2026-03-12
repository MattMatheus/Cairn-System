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
	if err := fs.Parse(args); err != nil {
		return err
	}
	tools, err := filteredTools(ctx, *registryPath, *includeLocal, *stage, "", *query)
	if err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("use.stage", *stage),
		attribute.String("use.query", *query),
		attribute.String("use.format", *format),
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
