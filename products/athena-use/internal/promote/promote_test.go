package promote

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadVaultNoteParsesFrontmatterAndBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "note.md")
	content := `---
id: ath-decision-test
type: decision
status: active
domain: research
sensitivity: internal
source_of_truth: human
---

# Test Note

Body line.
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	note, err := LoadVaultNote(path)
	if err != nil {
		t.Fatal(err)
	}
	if note.ID != "ath-decision-test" || note.Title != "Test Note" || note.Domain != "research" {
		t.Fatalf("unexpected note: %+v", note)
	}
	if strings.Contains(note.Body, "# Test Note") || !strings.Contains(note.Body, "Body line.") {
		t.Fatalf("unexpected note body: %q", note.Body)
	}
}

func TestBuildPlanUsesVaultMetadata(t *testing.T) {
	plan, err := BuildPlan(VaultNote{
		Path:        "/tmp/note.md",
		ID:          "ath-note",
		Title:       "Test Note",
		Domain:      "research",
		NoteType:    "decision",
		Sensitivity: "internal",
		Body:        "Body",
	}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if plan.MemoryType != "note" || plan.SourceKind != "obsidian-note" || plan.SourceType != "decision" {
		t.Fatalf("unexpected plan: %+v", plan)
	}
}

func TestBuildPlanRejectsPrivatePersonal(t *testing.T) {
	_, err := BuildPlan(VaultNote{
		Path:        "/tmp/note.md",
		ID:          "ath-private",
		Title:       "Private",
		Domain:      "personal",
		Sensitivity: "private_personal",
		Body:        "Body",
	}, Options{})
	if err == nil || !strings.Contains(err.Error(), "private_personal") {
		t.Fatalf("expected private note rejection, got %v", err)
	}
}

func TestAcquireMemoryRootLockSerializesAccess(t *testing.T) {
	root := t.TempDir()
	first, err := AcquireMemoryRootLock(root)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()

	acquired := make(chan struct{})
	release := make(chan struct{})
	errCh := make(chan error, 1)

	go func() {
		second, err := AcquireMemoryRootLock(root)
		if err != nil {
			errCh <- err
			return
		}
		close(acquired)
		<-release
		errCh <- second.Close()
	}()

	select {
	case <-acquired:
		t.Fatal("second lock should block until first lock is released")
	case <-time.After(100 * time.Millisecond):
	}

	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	select {
	case <-acquired:
	case err := <-errCh:
		t.Fatalf("unexpected lock acquisition error: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for second lock acquisition")
	}

	close(release)
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}
}
