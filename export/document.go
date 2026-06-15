package export

import (
	"github.com/eslider/go-hocr/legacy"
)

// Document is a flattened representation of a legacy hOCR file.
type Document struct {
	Pages []Page
}

// Page holds exported lines from a single OCR page.
type Page struct {
	Lines []interface{}
}

// NewDocument converts a parsed legacy hOCR document into an export-friendly structure.
func NewDocument(h legacy.Hocr, meta bool, withWords bool) (Document, error) {
	var pages []Page

	for _, p := range h.Pages {
		var lines []interface{}

		for _, l := range p.Lines {
			if withWords {
				var words []interface{}
				for _, w := range l.Words {
					if meta {
						words = append(words, struct {
							Text string
							Meta *map[string]string
						}{
							Text: w.GetCleanText(),
							Meta: w.GetMeta(),
						})
					} else {
						words = append(words, w.GetCleanText())
					}
				}

				if meta {
					lines = append(lines, struct {
						Text  string
						Meta  *map[string]string
						Words []interface{}
					}{
						Text:  legacy.LineText(l),
						Meta:  l.GetMeta(),
						Words: words,
					})
				} else {
					lines = append(lines, struct {
						Text  string
						Words []interface{}
					}{
						Text:  legacy.LineText(l),
						Words: words,
					})
				}
			} else {
				if meta {
					lines = append(lines, struct {
						Text string
						Meta *map[string]string
					}{
						Text: legacy.LineText(l),
						Meta: l.GetMeta(),
					})
				} else {
					lines = append(lines, legacy.LineText(l))
				}
			}
		}

		pages = append(pages, Page{Lines: lines})
	}

	return Document{Pages: pages}, nil
}
