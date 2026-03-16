package intake

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/net/html"
	"rsc.io/pdf"
)

const DefaultInbox = "01 Inbox"

type SourceType string

const (
	SourceURL    SourceType = "url"
	SourceFile   SourceType = "file"
	SourceFolder SourceType = "folder"
)

type Artifact struct {
	Title      string
	Source     string
	SourceType SourceType
	Markdown   string
	Links      []string
	Hash       string
	CreatedAt  time.Time
}

func Inspect(target string) (SourceType, string, error) {
	trimmed := strings.TrimSpace(target)
	if trimmed == "" {
		return "", "", errors.New("target is required")
	}
	if parsed, err := url.Parse(trimmed); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		return SourceURL, trimmed, nil
	}
	info, err := os.Stat(trimmed)
	if err != nil {
		return "", "", err
	}
	if info.IsDir() {
		return SourceFolder, trimmed, nil
	}
	return SourceFile, trimmed, nil
}

func NormalizeURL(ctx context.Context, target string) (Artifact, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return Artifact{}, err
	}
	req.Header.Set("User-Agent", "Cairn-Intake/0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Artifact{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Artifact{}, err
	}
	mediaType, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	artifact := Artifact{
		Source:     target,
		SourceType: SourceURL,
		Hash:       digest(body),
		CreatedAt:  time.Now().UTC(),
	}
	switch {
	case strings.Contains(mediaType, "html"), mediaType == "":
		title, markdown, links := htmlToMarkdown(body, target)
		artifact.Title = firstNonEmpty(title, fallbackTitle(target))
		artifact.Markdown = markdown
		artifact.Links = links
	default:
		text := strings.TrimSpace(string(body))
		artifact.Title = fallbackTitle(target)
		artifact.Markdown = "# " + artifact.Title + "\n\n" + text + "\n"
	}
	return artifact, nil
}

func NormalizeFile(path string) (Artifact, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Artifact{}, err
	}
	ext := strings.ToLower(filepath.Ext(path))
	artifact := Artifact{
		Source:     path,
		SourceType: SourceFile,
		Hash:       digest(data),
		CreatedAt:  time.Now().UTC(),
	}
	switch ext {
	case ".md", ".markdown":
		artifact.Title = firstMarkdownHeading(data, fallbackTitle(path))
		artifact.Markdown = string(data)
	case ".txt", ".log":
		artifact.Title = fallbackTitle(path)
		artifact.Markdown = "# " + artifact.Title + "\n\n" + strings.TrimSpace(string(data)) + "\n"
	case ".html", ".htm":
		title, markdown, links := htmlToMarkdown(data, "")
		artifact.Title = firstNonEmpty(title, fallbackTitle(path))
		artifact.Markdown = markdown
		artifact.Links = links
	case ".pdf":
		title, markdown, err := pdfToMarkdown(path)
		if err != nil {
			return Artifact{}, err
		}
		artifact.Title = firstNonEmpty(title, fallbackTitle(path))
		artifact.Markdown = markdown
	default:
		return Artifact{}, fmt.Errorf("unsupported file type: %s", ext)
	}
	return artifact, nil
}

func NormalizeFolder(path string) (Artifact, error) {
	var entries []Artifact
	err := filepath.WalkDir(path, func(current string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && current != path {
				return filepath.SkipDir
			}
			return nil
		}
		artifact, err := NormalizeFile(current)
		if err != nil {
			return nil
		}
		entries = append(entries, artifact)
		return nil
	})
	if err != nil {
		return Artifact{}, err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Source < entries[j].Source
	})
	var builder strings.Builder
	title := fallbackTitle(path)
	builder.WriteString("## Intake Summary\n\n")
	builder.WriteString("- Source folder: `" + path + "`\n")
	builder.WriteString("- Supported documents: " + fmt.Sprintf("%d", len(entries)) + "\n\n")
	if len(entries) == 0 {
		builder.WriteString("No supported files found.\n")
	} else {
		builder.WriteString("## Included Documents\n\n")
		for _, entry := range entries {
			rel, relErr := filepath.Rel(path, entry.Source)
			if relErr != nil {
				rel = entry.Source
			}
			builder.WriteString("- `" + rel + "`\n")
		}
		builder.WriteString("\n")
		builder.WriteString("## Document Contents\n\n")
		for _, entry := range entries {
			rel, relErr := filepath.Rel(path, entry.Source)
			if relErr != nil {
				rel = entry.Source
			}
			builder.WriteString("### " + rel + "\n\n")
			builder.WriteString(indentHeadings(trimTopHeading(strings.TrimSpace(entry.Markdown)), 1))
			builder.WriteString("\n\n")
		}
	}
	hashInput := make([]string, 0, len(entries))
	for _, entry := range entries {
		hashInput = append(hashInput, entry.Hash)
	}
	return Artifact{
		Title:      title,
		Source:     path,
		SourceType: SourceFolder,
		Markdown:   builder.String(),
		Hash:       digest([]byte(strings.Join(hashInput, "|"))),
		CreatedAt:  time.Now().UTC(),
	}, nil
}

func StageArtifact(vaultRoot, into string, artifact Artifact) (string, error) {
	vaultRoot = strings.TrimSpace(vaultRoot)
	if vaultRoot == "" {
		return "", errors.New("vault root is required")
	}
	into = strings.TrimSpace(into)
	if into == "" {
		into = DefaultInbox
	}
	createdAt := artifact.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	targetDir := filepath.Join(vaultRoot, filepath.FromSlash(into))
	if err := os.MkdirAll(targetDir, 0o775); err != nil {
		return "", err
	}
	filename := createdAt.Format("20060102-150405") + " " + sanitizeFilename(artifact.Title) + ".md"
	fullPath := filepath.Join(targetDir, filename)
	content := RenderArtifact(artifact)
	if err := os.WriteFile(fullPath, []byte(content), 0o664); err != nil {
		return "", err
	}
	return fullPath, nil
}

func RenderArtifact(artifact Artifact) string {
	var builder strings.Builder
	builder.WriteString("---\n")
	builder.WriteString("id: ath-artifact-intake-" + artifact.CreatedAt.Format("20060102-150405") + "\n")
	builder.WriteString("type: artifact\n")
	builder.WriteString("status: active\n")
	builder.WriteString("domain: research\n")
	builder.WriteString("updated: " + artifact.CreatedAt.Format("2006-01-02") + "\n")
	builder.WriteString("owner: matt\n")
	builder.WriteString("source_of_truth: agent\n")
	builder.WriteString("sensitivity: internal\n")
	builder.WriteString("artifact_kind: intake\n")
	builder.WriteString("supports_claims: []\n")
	builder.WriteString("refutes_claims: []\n")
	builder.WriteString("agent_last_touch: " + artifact.CreatedAt.Format(time.RFC3339) + "\n")
	builder.WriteString("agent_intent: ingest\n")
	builder.WriteString("review_state: pending\n")
	builder.WriteString("source_type: " + string(artifact.SourceType) + "\n")
	builder.WriteString("source_ref: \"" + escapeYAML(artifact.Source) + "\"\n")
	builder.WriteString("source_hash: " + artifact.Hash + "\n")
	builder.WriteString("---\n\n")
	builder.WriteString("# " + artifact.Title + "\n\n")
	builder.WriteString("## Context\n\n")
	builder.WriteString("- Source: `" + artifact.Source + "`\n")
	builder.WriteString("- Source type: `" + string(artifact.SourceType) + "`\n")
	if len(artifact.Links) > 0 {
		builder.WriteString("- Discovered links: " + fmt.Sprintf("%d", len(artifact.Links)) + "\n")
	}
	builder.WriteString("\n## Method\n\n")
	builder.WriteString("- Normalized into markdown for Athena review.\n")
	builder.WriteString("- Requires human triage before AthenaMind promotion.\n")
	builder.WriteString("\n## Review Scaffold\n\n")
	builder.WriteString("### Highlights\n\n")
	builder.WriteString("- Pending review\n")
	builder.WriteString("\n### Keep or Discard\n\n")
	builder.WriteString("- Decision: pending\n")
	builder.WriteString("- Reason: pending\n")
	builder.WriteString("\n### Promotion Notes\n\n")
	builder.WriteString("- AthenaMind candidate: no\n")
	builder.WriteString("- Notes: pending\n")
	builder.WriteString("\n## Result\n\n")
	builder.WriteString(trimTopHeading(strings.TrimSpace(artifact.Markdown)))
	builder.WriteString("\n\n## Links\n\n")
	if len(artifact.Links) == 0 {
		builder.WriteString("- None\n")
	} else {
		for _, link := range artifact.Links {
			builder.WriteString("- " + link + "\n")
		}
	}
	return builder.String()
}

func htmlToMarkdown(data []byte, base string) (string, string, []string) {
	root, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		text := strings.TrimSpace(string(data))
		return "", text, nil
	}
	baseURL, _ := url.Parse(strings.TrimSpace(base))
	title := extractDocumentTitle(root)
	var links []string
	var builder strings.Builder
	contentRoot := findPreferredContentRoot(root)
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "noscript":
				return
			case "h1", "h2", "h3":
				text := collectText(n)
				if text != "" {
					builder.WriteString(strings.Repeat("#", headingLevel(n.Data)) + " " + text + "\n\n")
				}
				return
			case "p", "li", "blockquote", "pre":
				text := collectText(n)
				if text != "" {
					prefix := ""
					if n.Data == "li" {
						prefix = "- "
					}
					builder.WriteString(prefix + text + "\n\n")
				}
				return
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						link := resolveLink(baseURL, strings.TrimSpace(attr.Val))
						if link != "" {
							links = append(links, link)
						}
						break
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(contentRoot)
	markdown := strings.TrimSpace(builder.String())
	if markdown == "" {
		markdown = strings.TrimSpace(collectText(root))
	}
	if title != "" && !strings.HasPrefix(markdown, "# ") {
		markdown = "# " + title + "\n\n" + markdown
	}
	return title, markdown + "\n", dedupeStrings(links)
}

func ApplyTitle(artifact Artifact, title string) Artifact {
	if strings.TrimSpace(title) != "" {
		artifact.Title = strings.TrimSpace(title)
		if strings.TrimSpace(artifact.Markdown) != "" {
			artifact.Markdown = replaceTopHeading(artifact.Markdown, artifact.Title)
		}
	}
	return artifact
}

func pdfToMarkdown(path string) (string, string, error) {
	reader, err := pdf.Open(path)
	if err != nil {
		return "", "", err
	}
	title := strings.TrimSpace(reader.Trailer().Key("Info").Key("Title").Text())
	var builder strings.Builder
	for i := 1; i <= reader.NumPage(); i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		text := extractPageText(page.Content().Text)
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		builder.WriteString("## Page ")
		builder.WriteString(fmt.Sprintf("%d", i))
		builder.WriteString("\n\n")
		builder.WriteString(text)
		builder.WriteString("\n\n")
	}
	markdown := strings.TrimSpace(builder.String())
	if title != "" {
		markdown = "# " + title + "\n\n" + markdown
	}
	if markdown == "" {
		markdown = "# " + fallbackTitle(path) + "\n\n(No extractable text found)\n"
	}
	return title, markdown + "\n", nil
}

func extractPageText(items []pdf.Text) string {
	if len(items) == 0 {
		return ""
	}
	sort.Slice(items, func(i, j int) bool {
		if abs(items[i].Y-items[j].Y) < 1.5 {
			return items[i].X < items[j].X
		}
		return items[i].Y > items[j].Y
	})
	var lines []string
	var current strings.Builder
	currentY := items[0].Y
	prevRight := items[0].X
	firstInLine := true
	for _, item := range items {
		if abs(item.Y-currentY) > 1.5 {
			if current.Len() > 0 {
				lines = append(lines, current.String())
			}
			current.Reset()
			currentY = item.Y
			prevRight = item.X
			firstInLine = true
		}
		if !firstInLine && item.X-prevRight > item.FontSize*0.05 {
			current.WriteString(" ")
		}
		current.WriteString(item.S)
		prevRight = item.X + item.W
		firstInLine = false
	}
	if current.Len() > 0 {
		lines = append(lines, current.String())
	}
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(compactNonEmpty(lines), "\n")
}

func collectText(n *html.Node) string {
	var parts []string
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			text := strings.TrimSpace(node.Data)
			if text != "" {
				parts = append(parts, text)
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(n)
	return strings.Join(parts, " ")
}

func headingLevel(tag string) int {
	switch tag {
	case "h1":
		return 1
	case "h2":
		return 2
	default:
		return 3
	}
}

func firstMarkdownHeading(data []byte, fallback string) string {
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
		}
	}
	return fallback
}

func fallbackTitle(target string) string {
	base := filepath.Base(target)
	base = strings.TrimSpace(strings.TrimSuffix(base, filepath.Ext(base)))
	if base == "" || base == "." || base == "/" {
		base = target
	}
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")
	return strings.TrimSpace(base)
}

func sanitizeFilename(input string) string {
	input = strings.TrimSpace(strings.ToLower(input))
	replacer := strings.NewReplacer("/", "-", "\\", "-", ":", "-", "*", "", "?", "", "\"", "", "<", "", ">", "", "|", "")
	input = replacer.Replace(input)
	input = strings.Join(strings.Fields(input), "-")
	if input == "" {
		return "intake-artifact"
	}
	return input
}

func escapeYAML(v string) string {
	return strings.ReplaceAll(v, "\"", "\\\"")
}

func digest(data []byte) string {
	sum := sha1.Sum(data)
	return hex.EncodeToString(sum[:])
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func dedupeStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func findPreferredContentRoot(root *html.Node) *html.Node {
	if node := findElement(root, "main"); node != nil {
		return node
	}
	if node := findElement(root, "article"); node != nil {
		return node
	}
	if node := findElement(root, "body"); node != nil {
		return node
	}
	return root
}

func extractDocumentTitle(root *html.Node) string {
	if node := findElement(root, "title"); node != nil && node.FirstChild != nil {
		return strings.TrimSpace(node.FirstChild.Data)
	}
	return ""
}

func findElement(root *html.Node, tag string) *html.Node {
	var result *html.Node
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if result != nil {
			return
		}
		if n.Type == html.ElementNode && n.Data == tag {
			result = n
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(root)
	return result
}

func resolveLink(base *url.URL, raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	link, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if base != nil {
		return base.ResolveReference(link).String()
	}
	return link.String()
}

func replaceTopHeading(markdown, title string) string {
	lines := strings.Split(markdown, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "# ") {
			lines[i] = "# " + title
			return strings.Join(lines, "\n")
		}
		if strings.TrimSpace(line) != "" {
			break
		}
	}
	return "# " + title + "\n\n" + strings.TrimSpace(markdown) + "\n"
}

func trimTopHeading(markdown string) string {
	lines := strings.Split(strings.TrimSpace(markdown), "\n")
	if len(lines) == 0 {
		return ""
	}
	if strings.HasPrefix(strings.TrimSpace(lines[0]), "# ") {
		return strings.TrimSpace(strings.Join(lines[1:], "\n"))
	}
	return strings.TrimSpace(markdown)
}

func indentHeadings(markdown string, levels int) string {
	if levels <= 0 {
		return markdown
	}
	prefix := strings.Repeat("#", levels)
	lines := strings.Split(markdown, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}

func compactNonEmpty(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, value)
		}
	}
	return out
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
