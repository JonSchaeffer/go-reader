package rss

import (
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// Content processor struct
type ContentProcessor struct {
	sanitizer *bluemonday.Policy
}

func NewContentProcessor() *ContentProcessor {
	// Create a sanitization policy
	p := bluemonday.UGCPolicy()

	// Allow specific tags and attributes for articles
	p.AllowElements("p", "br", "strong", "em", "b", "i", "u", "h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowElements("ul", "ol", "li", "blockquote", "pre", "code")
	p.AllowElements("a").AllowAttrs("href", "title").OnElements("a")
	p.AllowElements("img").AllowAttrs("src", "alt", "title", "width", "height").OnElements("img")

	return &ContentProcessor{sanitizer: p}
}

func (cp *ContentProcessor) ProcessContent(rawContent string) string {
	// 1. Clean HTML
	cleaned := cp.sanitizer.Sanitize(rawContent)

	// 2. Fix common issues
	cleaned = cp.fixCommonIssues(cleaned)

	// 3. Normalize structure
	cleaned = cp.normalizeStructure(cleaned)

	return cleaned
}

func (cp *ContentProcessor) fixCommonIssues(content string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	content = re.ReplaceAllString(content, " ")

	// Fix multiple line breaks
	re = regexp.MustCompile(`<br\s*/?>\s*<br\s*/?>`)
	content = re.ReplaceAllString(content, "</p><p>")

	// Wrap orphaned text in paragraphs
	if !strings.Contains(content, "<p>") && len(content) > 0 {
		content = "<p>" + content + "</p>"
	}

	return strings.TrimSpace(content)
}

func (cp *ContentProcessor) normalizeStructure(content string) string {
	// Convert <b> to <strong>, <i> to <em>, etc.
	replacements := map[string]string{
		"<b>":  "<strong>",
		"</b>": "</strong>",
		"<i>":  "<em>",
		"</i>": "</em>",
	}

	for old, new := range replacements {
		content = strings.ReplaceAll(content, old, new)
	}

	return content
}
