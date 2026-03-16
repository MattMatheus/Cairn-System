package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"toolcli/internal/codegraph"
	"toolcli/internal/telemetry"

	"go.opentelemetry.io/otel/attribute"
)

func runAnalyze(args []string) (err error) {
	ctx, span := telemetry.StartCommandSpan(context.Background(), "analyze")
	defer func() { telemetry.EndSpan(span, err) }()

	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	repo := fs.String("repo", "", "target repository path")
	force := fs.Bool("force", false, "force full reindex")
	format := fs.String("format", "text", "output format: text|json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return runAndWrite(ctx, span, *format, codegraph.Options{
		Command: codegraph.CommandAnalyze,
		Repo:    strings.TrimSpace(*repo),
		Force:   *force,
	})
}

func runStatus(args []string) (err error) {
	ctx, span := telemetry.StartCommandSpan(context.Background(), "status")
	defer func() { telemetry.EndSpan(span, err) }()

	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	repo := fs.String("repo", "", "target repository path")
	format := fs.String("format", "text", "output format: text|json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return runAndWrite(ctx, span, *format, codegraph.Options{
		Command: codegraph.CommandStatus,
		Repo:    strings.TrimSpace(*repo),
	})
}

func runContext(args []string) (err error) {
	ctx, span := telemetry.StartCommandSpan(context.Background(), "context")
	defer func() { telemetry.EndSpan(span, err) }()

	fs := flag.NewFlagSet("context", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	repo := fs.String("repo", "", "target repository path")
	file := fs.String("file", "", "optional file path to disambiguate symbol context")
	format := fs.String("format", "text", "output format: text|json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("context requires exactly one symbol")
	}
	return runAndWrite(ctx, span, *format, codegraph.Options{
		Command: codegraph.CommandContext,
		Repo:    strings.TrimSpace(*repo),
		Target:  strings.TrimSpace(fs.Args()[0]),
		File:    strings.TrimSpace(*file),
	})
}

func runImpact(args []string) (err error) {
	ctx, span := telemetry.StartCommandSpan(context.Background(), "impact")
	defer func() { telemetry.EndSpan(span, err) }()

	fs := flag.NewFlagSet("impact", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	repo := fs.String("repo", "", "target repository path")
	direction := fs.String("direction", "upstream", "impact direction: upstream|downstream")
	format := fs.String("format", "text", "output format: text|json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("impact requires exactly one symbol")
	}
	return runAndWrite(ctx, span, *format, codegraph.Options{
		Command:   codegraph.CommandImpact,
		Repo:      strings.TrimSpace(*repo),
		Target:    strings.TrimSpace(fs.Args()[0]),
		Direction: strings.TrimSpace(*direction),
	})
}

func runAndWrite(ctx context.Context, span traceLike, format string, opts codegraph.Options) error {
	span.SetAttributes(
		attribute.String("codegraph.command", string(opts.Command)),
		attribute.String("codegraph.repo", opts.Repo),
		attribute.String("codegraph.target", opts.Target),
		attribute.String("codegraph.direction", opts.Direction),
		attribute.Bool("codegraph.force", opts.Force),
	)
	output, err := executeCodegraph(ctx, opts)
	if err != nil {
		return err
	}
	return writeOutput(output, format)
}

var executeCodegraph = codegraph.Run

type traceLike interface {
	SetAttributes(...attribute.KeyValue)
}

func writeOutput(output []byte, format string) error {
	text := strings.TrimSpace(string(output))
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "text":
		fmt.Println(text)
		return nil
	case "json":
		payload := map[string]any{
			"output": text,
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
