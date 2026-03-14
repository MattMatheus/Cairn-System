package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestRunInspect(t *testing.T) {
	out := captureStdout(t, func() error {
		return runInspect([]string{"--format", "json", "https://example.com"})
	})
	if !strings.Contains(out, "\"source_type\": \"url\"") {
		t.Fatalf("unexpected inspect output: %s", out)
	}
}

func TestRunFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.md")
	if err := os.WriteFile(path, []byte("# Sample\n\nBody\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() error {
		return runFile([]string{"--format", "json", "--title", "Custom Sample", path})
	})
	if !strings.Contains(out, "\"title\": \"Custom Sample\"") {
		t.Fatalf("unexpected file output: %s", out)
	}
}

func TestRunStage(t *testing.T) {
	vault := t.TempDir()
	filePath := filepath.Join(t.TempDir(), "sample.txt")
	if err := os.WriteFile(filePath, []byte("Body"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() error {
		return runStage([]string{"--vault", vault, "--into", "01 Inbox", "--format", "json", filePath})
	})
	if !strings.Contains(out, "\"staged_path\":") {
		t.Fatalf("unexpected stage output: %s", out)
	}
	files, err := os.ReadDir(filepath.Join(vault, "01 Inbox"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected one staged file, got %d", len(files))
	}
}

func TestDefaultVaultRootUsesEnv(t *testing.T) {
	t.Setenv("ATHENA_VAULT", "/tmp/athena-vault")
	if got := defaultVaultRoot(); got != "/tmp/athena-vault" {
		t.Fatalf("unexpected vault root: %s", got)
	}
}
