package intake

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestInspect(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "note.md")
	if err := os.WriteFile(filePath, []byte("# Note\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	gotType, _, err := Inspect(filePath)
	if err != nil || gotType != SourceFile {
		t.Fatalf("expected file source, got %q err=%v", gotType, err)
	}
	gotType, _, err = Inspect(dir)
	if err != nil || gotType != SourceFolder {
		t.Fatalf("expected folder source, got %q err=%v", gotType, err)
	}
	gotType, _, err = Inspect("https://example.com")
	if err != nil || gotType != SourceURL {
		t.Fatalf("expected url source, got %q err=%v", gotType, err)
	}
}

func TestNormalizeFileMarkdown(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.md")
	if err := os.WriteFile(path, []byte("# Sample\n\nHello world.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact, err := NormalizeFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if artifact.Title != "Sample" {
		t.Fatalf("unexpected title: %s", artifact.Title)
	}
	if !strings.Contains(artifact.Markdown, "Hello world.") {
		t.Fatalf("expected markdown body, got %s", artifact.Markdown)
	}
}

func TestNormalizeURLHTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><head><title>Example Page</title></head><body><nav><a href=\"/nav\">Nav</a></nav><main><h1>Example Page</h1><p>Hello world</p><a href=\"/a\">A</a></main></body></html>"))
	}))
	defer server.Close()

	artifact, err := NormalizeURL(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if artifact.Title != "Example Page" {
		t.Fatalf("unexpected title: %s", artifact.Title)
	}
	if !strings.Contains(artifact.Markdown, "Hello world") {
		t.Fatalf("expected markdown content, got %s", artifact.Markdown)
	}
	if len(artifact.Links) != 1 {
		t.Fatalf("expected one discovered link, got %+v", artifact.Links)
	}
	if artifact.Links[0] != server.URL+"/a" {
		t.Fatalf("expected resolved absolute link, got %+v", artifact.Links)
	}
	if strings.Contains(artifact.Markdown, "Nav") {
		t.Fatalf("expected main content focus to skip nav, got %s", artifact.Markdown)
	}
}

func TestNormalizeFolder(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("# A\n\nOne\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("Two"), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact, err := NormalizeFolder(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(artifact.Markdown, "a.md") || !strings.Contains(artifact.Markdown, "b.txt") {
		t.Fatalf("expected folder summary, got %s", artifact.Markdown)
	}
	if strings.HasPrefix(strings.TrimSpace(artifact.Markdown), "# ") {
		t.Fatalf("folder markdown should not start with a duplicated top-level heading: %s", artifact.Markdown)
	}
	if !strings.Contains(artifact.Markdown, "## Included Documents") || !strings.Contains(artifact.Markdown, "## Document Contents") {
		t.Fatalf("expected stronger folder structure, got %s", artifact.Markdown)
	}
	if !strings.Contains(artifact.Markdown, "### a.md") {
		t.Fatalf("expected per-document subsection heading, got %s", artifact.Markdown)
	}
}

func TestStageArtifact(t *testing.T) {
	vault := t.TempDir()
	artifact := Artifact{
		Title:      "Example Intake",
		Source:     "https://example.com",
		SourceType: SourceURL,
		Markdown:   "# Example Intake\n\nBody\n",
		Hash:       "abc123",
		CreatedAt:  mustTime("2026-03-13T21:00:00Z"),
	}
	path, err := StageArtifact(vault, "", artifact)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, "artifact_kind: intake") || !strings.Contains(text, "source_ref: \"https://example.com\"") {
		t.Fatalf("unexpected staged artifact: %s", text)
	}
	if !strings.Contains(text, "## Review Scaffold") || !strings.Contains(text, "### Keep or Discard") {
		t.Fatalf("expected review scaffold in staged artifact: %s", text)
	}
}

func TestRenderArtifactTrimsDuplicateTopHeading(t *testing.T) {
	artifact := Artifact{
		Title:      "Folder Packet",
		Source:     "/tmp/folder",
		SourceType: SourceFolder,
		Markdown:   "# Folder Packet\n\n## Intake Summary\n\n- Example\n",
		Hash:       "abc123",
		CreatedAt:  mustTime("2026-03-13T21:00:00Z"),
	}
	text := RenderArtifact(artifact)
	if strings.Count(text, "# Folder Packet") != 1 {
		t.Fatalf("expected only one top-level heading, got %s", text)
	}
	if !strings.Contains(text, "### Highlights") || !strings.Contains(text, "### Promotion Notes") {
		t.Fatalf("expected review scaffold sections, got %s", text)
	}
}

func TestIndentHeadings(t *testing.T) {
	got := indentHeadings("## Section\n\nBody\n\n### Subsection", 1)
	if !strings.Contains(got, "### Section") || !strings.Contains(got, "#### Subsection") {
		t.Fatalf("unexpected heading indentation: %s", got)
	}
}

func TestApplyTitle(t *testing.T) {
	artifact := Artifact{
		Title:    "Old",
		Markdown: "# Old\n\nBody\n",
	}
	updated := ApplyTitle(artifact, "New Title")
	if updated.Title != "New Title" {
		t.Fatalf("unexpected title: %s", updated.Title)
	}
	if !strings.HasPrefix(updated.Markdown, "# New Title") {
		t.Fatalf("expected updated heading, got %s", updated.Markdown)
	}
}

func TestNormalizeFilePDF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.pdf")
	if err := os.WriteFile(path, buildSimplePDF("Hello PDF"), 0o644); err != nil {
		t.Fatal(err)
	}
	artifact, err := NormalizeFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.ReplaceAll(artifact.Markdown, " ", ""), "HelloPDF") {
		t.Fatalf("expected extracted pdf text, got %s", artifact.Markdown)
	}
}

func mustTime(raw string) (outTime time.Time) {
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		panic(err)
	}
	return parsed
}

func buildSimplePDF(text string) []byte {
	objects := []string{
		"1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n",
		"2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n",
		"3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 300 144] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>\nendobj\n",
		fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\nBT\n/F1 12 Tf\n72 72 Td\n(%s) Tj\nET\nendstream\nendobj\n", len("BT\n/F1 12 Tf\n72 72 Td\n("+escapePDFText(text)+") Tj\nET\n"), escapePDFText(text)),
		"5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n",
	}
	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n")
	offsets := make([]int, 0, len(objects)+1)
	offsets = append(offsets, 0)
	for _, object := range objects {
		offsets = append(offsets, buf.Len())
		buf.WriteString(object)
	}
	xrefOffset := buf.Len()
	buf.WriteString("xref\n")
	buf.WriteString(fmt.Sprintf("0 %d\n", len(offsets)))
	buf.WriteString("0000000000 65535 f \n")
	for _, offset := range offsets[1:] {
		buf.WriteString(fmt.Sprintf("%010d 00000 n \n", offset))
	}
	buf.WriteString("trailer\n")
	buf.WriteString(fmt.Sprintf("<< /Size %d /Root 1 0 R >>\n", len(offsets)))
	buf.WriteString("startxref\n")
	buf.WriteString(fmt.Sprintf("%d\n", xrefOffset))
	buf.WriteString("%%EOF\n")
	return buf.Bytes()
}

func escapePDFText(v string) string {
	replacer := strings.NewReplacer("\\", "\\\\", "(", "\\(", ")", "\\)")
	return replacer.Replace(v)
}
