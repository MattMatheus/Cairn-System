package codegraph

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildExecCommandBinaryAnalyze(t *testing.T) {
	cmd, err := BuildExecCommand(Install{Mode: "binary", BinaryPath: "/usr/bin/gitnexus"}, Options{
		Command: CommandAnalyze,
		Repo:    "/tmp/repo",
		Force:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
	got := cmd.Args
	want := []string{"/usr/bin/gitnexus", "analyze", "/tmp/repo", "--force"}
	if len(got) != len(want) {
		t.Fatalf("unexpected args: %+v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected args: %+v", got)
		}
	}
}

func TestBuildExecCommandNodeImpact(t *testing.T) {
	cmd, err := BuildExecCommand(Install{Mode: "node", NodePath: "/usr/bin/node", EntryPoint: "/opt/gitnexus/dist/cli/index.js"}, Options{
		Command:   CommandImpact,
		Target:    "BoundField",
		Direction: "downstream",
	})
	if err != nil {
		t.Fatal(err)
	}
	got := cmd.Args
	want := []string{"/usr/bin/node", "/opt/gitnexus/dist/cli/index.js", "impact", "BoundField", "--direction", "downstream"}
	if len(got) != len(want) {
		t.Fatalf("unexpected args: %+v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected args: %+v", got)
		}
	}
}

func TestResolveInstallPrefersExplicitBinary(t *testing.T) {
	install, err := ResolveInstall(
		func() (string, error) { return "/tmp", nil },
		func(key string) (string, bool) {
			if key == "CAIRN_GITNEXUS_BIN" {
				return "/custom/gitnexus", true
			}
			return "", false
		},
		func(string) (string, error) { return "", errors.New("not found") },
	)
	if err != nil {
		t.Fatal(err)
	}
	if install.Mode != "binary" || install.BinaryPath != "/custom/gitnexus" {
		t.Fatalf("unexpected install: %+v", install)
	}
}

func TestResolveInstallFindsBuiltCheckout(t *testing.T) {
	root := t.TempDir()
	repoRoot := filepath.Join(root, "repos", "untrusted", "GitNexus", "gitnexus")
	if err := os.MkdirAll(filepath.Join(repoRoot, "dist", "cli"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "package.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "dist", "cli", "index.js"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	install, err := ResolveInstall(
		func() (string, error) { return filepath.Join(root, "repos", "trusted", "Cairn"), nil },
		func(string) (string, bool) { return "", false },
		func(name string) (string, error) {
			if name == "node" {
				return "/usr/bin/node", nil
			}
			return "", errors.New("not found")
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if install.Mode != "node" || install.EntryPoint == "" {
		t.Fatalf("unexpected install: %+v", install)
	}
}

func TestBuildArgsRequiresTargets(t *testing.T) {
	if _, err := buildArgs(Options{Command: CommandContext}); err == nil {
		t.Fatal("expected context target validation")
	}
	if _, err := buildArgs(Options{Command: CommandImpact}); err == nil {
		t.Fatal("expected impact target validation")
	}
}
