package hocr

import "testing"

func TestToMarkdown_Basic(t *testing.T) {
	doc := &Document{
		Title: "Sample",
		Pages: []*Page{
			{
				Blocks: []*Block{
					{
						Element: Element{Class: ContentAreaBlock},
						Paragraphs: []*Paragraph{
							{
								Lines: []*Line{
									{
										Element: Element{Class: HeaderLine},
										Words: []*Word{
											{Text: "INFORMACIÓN", Confidence: 95},
										},
									},
								},
							},
							{
								Lines: []*Line{
									{
										Element: Element{Class: TextLine},
										Words: []*Word{
											{Text: "Hello", Confidence: 90},
											{Text: "world", Confidence: 10}, // filtered
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	md := doc.ToMarkdown(50)
	if !contains(md, "# Sample") || !contains(md, "## INFORMACIÓN") || !contains(md, "Hello") {
		t.Fatalf("unexpected md:\n%s", md)
	}
	if contains(md, "world") {
		t.Fatalf("low-confidence word should be filtered:\n%s", md)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i+len(sub) <= len(s); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
