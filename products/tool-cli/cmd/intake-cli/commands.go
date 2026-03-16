package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"toolcli/internal/intake"
	"toolcli/internal/telemetry"
)

func runInspect(args []string) (err error) {
	_, span := telemetry.StartCommandSpan(context.Background(), "inspect")
	defer func() { telemetry.EndSpan(span, err) }()

	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	format := fs.String("format", "text", "output format: text|json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("inspect requires exactly one target")
	}
	sourceType, target, err := intake.Inspect(fs.Args()[0])
	if err != nil {
		return err
	}
	span.SetAttributes(
		attribute.String("intake.target", target),
		attribute.String("intake.source_type", string(sourceType)),
	)
	summary := map[string]any{
		"target":       target,
		"source_type":  string(sourceType),
		"writable":     true,
		"default_into": intake.DefaultInbox,
	}
	return writeOutput(summary, *format)
}

func runURL(args []string) (err error) {
	ctx, span := telemetry.StartCommandSpan(context.Background(), "url")
	defer func() { telemetry.EndSpan(span, err) }()
	fs := flag.NewFlagSet("url", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	format := fs.String("format", "text", "output format: text|json")
	out := fs.String("out", "", "optional output file path")
	title := fs.String("title", "", "optional title override")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("url requires exactly one target")
	}
	artifact, err := intake.NormalizeURL(ctx, fs.Args()[0])
	if err != nil {
		return err
	}
	artifact = intake.ApplyTitle(artifact, *title)
	return emitArtifact(*out, *format, artifact, span)
}

func runFile(args []string) (err error) {
	_, span := telemetry.StartCommandSpan(context.Background(), "file")
	defer func() { telemetry.EndSpan(span, err) }()
	fs := flag.NewFlagSet("file", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	format := fs.String("format", "text", "output format: text|json")
	out := fs.String("out", "", "optional output file path")
	title := fs.String("title", "", "optional title override")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("file requires exactly one path")
	}
	artifact, err := intake.NormalizeFile(fs.Args()[0])
	if err != nil {
		return err
	}
	artifact = intake.ApplyTitle(artifact, *title)
	return emitArtifact(*out, *format, artifact, span)
}

func runFolder(args []string) (err error) {
	_, span := telemetry.StartCommandSpan(context.Background(), "folder")
	defer func() { telemetry.EndSpan(span, err) }()
	fs := flag.NewFlagSet("folder", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	format := fs.String("format", "text", "output format: text|json")
	out := fs.String("out", "", "optional output file path")
	title := fs.String("title", "", "optional title override")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("folder requires exactly one path")
	}
	artifact, err := intake.NormalizeFolder(fs.Args()[0])
	if err != nil {
		return err
	}
	artifact = intake.ApplyTitle(artifact, *title)
	return emitArtifact(*out, *format, artifact, span)
}

func runStage(args []string) (err error) {
	ctx, span := telemetry.StartCommandSpan(context.Background(), "stage")
	defer func() { telemetry.EndSpan(span, err) }()
	fs := flag.NewFlagSet("stage", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	into := fs.String("into", intake.DefaultInbox, "vault-relative destination")
	vault := fs.String("vault", defaultVaultRoot(), "Obsidian vault root")
	format := fs.String("format", "text", "output format: text|json")
	title := fs.String("title", "", "optional title override")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("stage requires exactly one target")
	}
	artifact, err := normalizeTarget(ctx, fs.Args()[0])
	if err != nil {
		return err
	}
	artifact = intake.ApplyTitle(artifact, *title)
	path, err := intake.StageArtifact(*vault, *into, artifact)
	if err != nil {
		return err
	}
	span.SetAttributes(
		attribute.String("intake.staged_path", path),
		attribute.String("intake.vault", *vault),
		attribute.String("intake.into", *into),
	)
	return writeOutput(map[string]any{
		"staged_path": path,
		"title":       artifact.Title,
		"source":      artifact.Source,
		"source_type": string(artifact.SourceType),
	}, *format)
}

func normalizeTarget(ctx context.Context, target string) (intake.Artifact, error) {
	sourceType, normalizedTarget, err := intake.Inspect(target)
	if err != nil {
		return intake.Artifact{}, err
	}
	switch sourceType {
	case intake.SourceURL:
		return intake.NormalizeURL(ctx, normalizedTarget)
	case intake.SourceFile:
		return intake.NormalizeFile(normalizedTarget)
	case intake.SourceFolder:
		return intake.NormalizeFolder(normalizedTarget)
	default:
		return intake.Artifact{}, fmt.Errorf("unsupported source type: %s", sourceType)
	}
}

func emitArtifact(outPath, format string, artifact intake.Artifact, span traceLike) error {
	span.SetAttributes(
		attribute.String("intake.title", artifact.Title),
		attribute.String("intake.source", artifact.Source),
		attribute.String("intake.source_type", string(artifact.SourceType)),
	)
	if strings.TrimSpace(outPath) != "" {
		if err := os.MkdirAll(filepath.Dir(outPath), 0o775); err != nil {
			return err
		}
		if err := os.WriteFile(outPath, []byte(artifact.Markdown), 0o664); err != nil {
			return err
		}
	}
	payload := map[string]any{
		"title":       artifact.Title,
		"source":      artifact.Source,
		"source_type": string(artifact.SourceType),
		"hash":        artifact.Hash,
		"links":       artifact.Links,
		"markdown":    artifact.Markdown,
	}
	return writeOutput(payload, format)
}

type traceLike interface {
	SetAttributes(...attribute.KeyValue)
}

func writeOutput(payload map[string]any, format string) error {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	case "text":
		if markdown, ok := payload["markdown"].(string); ok && markdown != "" {
			fmt.Println(strings.TrimSpace(markdown))
			return nil
		}
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func defaultVaultRoot() string {
	if value := strings.TrimSpace(os.Getenv("CAIRN_VAULT")); value != "" {
		return value
	}
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return "."
	}
	return filepath.Join(home, "Workspace", "Cairn Vault")
}
