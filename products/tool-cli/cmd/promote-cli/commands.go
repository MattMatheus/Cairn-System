package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"toolcli/internal/promote"
	"toolcli/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
)

func runInspect(args []string) (err error) {
	_, span := telemetry.StartCommandSpan(context.Background(), "inspect")
	defer func() { telemetry.EndSpan(span, err) }()

	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	format := fs.String("format", "json", "output format: text|json|yaml")
	id := fs.String("id", "", "optional override for the memory-cli entry id")
	title := fs.String("title", "", "optional override for the memory-cli entry title")
	domain := fs.String("domain", "", "optional override for the memory-cli entry domain")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("inspect requires exactly one note path")
	}
	note, err := promote.LoadVaultNote(fs.Args()[0])
	if err != nil {
		return err
	}
	plan, err := promote.BuildPlan(note, promote.Options{ID: *id, Title: *title, Domain: *domain})
	if err != nil {
		return err
	}
	span.SetAttributes(
		attribute.String("promote.note_path", note.Path),
		attribute.String("promote.memory_id", plan.MemoryID),
	)
	return writeOutput(map[string]any{
		"note_path":       note.Path,
		"note_type":       note.NoteType,
		"note_status":     note.Status,
		"source_of_truth": note.SourceOfTruth,
		"plan":            plan,
	}, *format)
}

func runNote(args []string) (err error) {
	_, span := telemetry.StartCommandSpan(context.Background(), "note")
	defer func() { telemetry.EndSpan(span, err) }()

	fs := flag.NewFlagSet("note", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", defaultMemoryRoot(), "memory-cli memory root")
	stage := fs.String("stage", "pm", "workflow stage: planning|architect|pm")
	reviewer := fs.String("reviewer", "", "reviewer identity")
	decision := fs.String("decision", "approved", "review decision: approved|rejected")
	reason := fs.String("reason", "", "reason for promotion")
	risk := fs.String("risk", "", "risk note")
	notes := fs.String("notes", "", "decision notes")
	id := fs.String("id", "", "optional override for the memory-cli entry id")
	title := fs.String("title", "", "optional override for the memory-cli entry title")
	domain := fs.String("domain", "", "optional override for the memory-cli entry domain")
	format := fs.String("format", "json", "output format: text|json|yaml")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 1 {
		return fmt.Errorf("note requires exactly one note path")
	}
	note, err := promote.LoadVaultNote(fs.Args()[0])
	if err != nil {
		return err
	}
	plan, err := promote.BuildPlan(note, promote.Options{ID: *id, Title: *title, Domain: *domain})
	if err != nil {
		return err
	}
	span.SetAttributes(
		attribute.String("promote.note_path", note.Path),
		attribute.String("promote.memory_root", *root),
		attribute.String("promote.memory_id", plan.MemoryID),
	)
	lock, err := promote.AcquireMemoryRootLock(*root)
	if err != nil {
		return err
	}
	defer lock.Close()

	tmp, err := os.CreateTemp("", "cairn-promote-body-*.md")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := tmp.WriteString(plan.Body); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	cmd := exec.Command(
		goCommand(),
		"run",
		"./cmd/memory-cli",
		"write",
		"--root", *root,
		"--id", plan.MemoryID,
		"--title", plan.Title,
		"--type", plan.MemoryType,
		"--domain", plan.Domain,
		"--body-file", tmpPath,
		"--stage", *stage,
		"--reviewer", *reviewer,
		"--decision", *decision,
		"--reason", *reason,
		"--risk", *risk,
		"--notes", *notes,
		"--source-ref", plan.SourceRef,
		"--source-kind", plan.SourceKind,
		"--source-type", plan.SourceType,
	)
	cmd.Dir = filepath.Dir(filepath.Dir(memoryCLIPath()))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}
	return writeOutput(map[string]any{
		"note_path":    note.Path,
		"memory_root":  *root,
		"memory_id":    plan.MemoryID,
		"title":        plan.Title,
		"domain":       plan.Domain,
		"source_ref":   plan.SourceRef,
		"source_type":  plan.SourceType,
		"review_stage": *stage,
		"memory_cli":   string(out),
	}, *format)
}
