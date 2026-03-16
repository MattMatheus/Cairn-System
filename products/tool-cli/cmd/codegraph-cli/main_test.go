package main

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"toolcli/internal/codegraph"
)

func captureStdout(t *testing.T, fn func() error) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	if err := fn(); err != nil {
		t.Fatal(err)
	}
	_ = w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestRunAnalyzePassesForceAndRepo(t *testing.T) {
	old := executeCodegraph
	defer func() { executeCodegraph = old }()
	executeCodegraph = func(_ context.Context, opts codegraph.Options) ([]byte, error) {
		if opts.Command != codegraph.CommandAnalyze || opts.Repo != "/tmp/repo" || !opts.Force {
			t.Fatalf("unexpected options: %+v", opts)
		}
		return []byte("ok"), nil
	}
	out := captureStdout(t, func() error {
		return runAnalyze([]string{"--repo", "/tmp/repo", "--force"})
	})
	if !strings.Contains(out, "ok") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestRunContextRequiresSymbol(t *testing.T) {
	if err := runContext(nil); err == nil {
		t.Fatal("expected symbol validation")
	}
}

func TestRunImpactPassesDirection(t *testing.T) {
	old := executeCodegraph
	defer func() { executeCodegraph = old }()
	executeCodegraph = func(_ context.Context, opts codegraph.Options) ([]byte, error) {
		if opts.Command != codegraph.CommandImpact || opts.Target != "BoundField" || opts.Direction != "downstream" {
			t.Fatalf("unexpected options: %+v", opts)
		}
		return []byte("impact"), nil
	}
	out := captureStdout(t, func() error {
		return runImpact([]string{"--direction", "downstream", "BoundField"})
	})
	if !strings.Contains(out, "impact") {
		t.Fatalf("unexpected output: %s", out)
	}
}
