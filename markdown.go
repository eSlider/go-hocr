package hocr

import (
	"fmt"
	"strings"
)

// ToMarkdown renders the document as readable Markdown (pages / paragraphs / lines).
// Word confidence below minConf is dropped (use 0 to keep all). Headers become ##.
func (d *Document) ToMarkdown(minConf float32) string {
	var b strings.Builder
	if t := strings.TrimSpace(d.Title); t != "" {
		b.WriteString("# ")
		b.WriteString(t)
		b.WriteString("\n\n")
	}
	for i, page := range d.Pages {
		if len(d.Pages) > 1 {
			if i > 0 {
				b.WriteString("\n---\n\n")
			}
			fmt.Fprintf(&b, "## Page %d\n\n", i+1)
		}
		b.WriteString(page.ToMarkdown(minConf))
	}
	return strings.TrimSpace(b.String()) + "\n"
}

// ToMarkdown renders page blocks as Markdown paragraphs.
func (p *Page) ToMarkdown(minConf float32) string {
	var parts []string
	for _, block := range p.Blocks {
		if s := block.ToMarkdown(minConf); s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, "\n\n") + "\n"
}

// ToMarkdown joins paragraphs; separators become a thematic break.
func (b *Block) ToMarkdown(minConf float32) string {
	if b.IsSeparator() {
		return "---"
	}
	if b.IsPhoto() {
		return "![photo]()"
	}
	var parts []string
	for _, p := range b.Paragraphs {
		if s := p.ToMarkdown(minConf); s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, "\n\n")
}

// ToMarkdown joins lines; header lines become ##.
func (p *Paragraph) ToMarkdown(minConf float32) string {
	var lines []string
	for _, l := range p.Lines {
		if s := l.ToMarkdown(minConf); s != "" {
			lines = append(lines, s)
		}
	}
	if len(lines) == 0 {
		return ""
	}
	// Single header line → heading; otherwise keep line breaks within paragraph.
	if len(lines) == 1 && strings.HasPrefix(lines[0], "## ") {
		return lines[0]
	}
	return strings.Join(lines, "\n")
}

// ToMarkdown returns plain words; header/caption lines get Markdown markers.
func (l *Line) ToMarkdown(minConf float32) string {
	var words []string
	for _, w := range l.Words {
		if minConf > 0 && w.Confidence > 0 && w.Confidence < minConf {
			continue
		}
		t := strings.TrimSpace(w.Text)
		if t != "" {
			words = append(words, t)
		}
	}
	if len(words) == 0 {
		return ""
	}
	text := strings.Join(words, " ")
	switch {
	case l.IsHeader():
		return "## " + text
	case l.IsCaption():
		return "**" + text + "**"
	default:
		return text
	}
}
