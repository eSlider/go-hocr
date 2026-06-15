package hocr

import (
	"fmt"
	"github.com/eslider/go-hocr/v1_2"
	"strings"
)

// Paragraph ocr_display, ocr_blockquote and ocr_par
type Paragraph struct {
	Element  `yaml:",inline"`
	Language string // "deu","eng","rus", etc.
	Lines    []*Line
}

func (p *Paragraph) GetText() string {
	var ls []string
	for _, l := range p.Lines {

		ls = append(ls, l.GetHtml())

	}
	if ls == nil {
		return ""
	}
	return strings.Join(ls, "<br/>")

}

// NewParagraph creates a new paragraph from a hocr element
// https://kba.github.io/hocr-spec/1.2/#special-paragraphs
func NewParagraph(el *v1_2.Element) *Paragraph {
	// Check for new fields
	checkForNewProperties(el, []string{
		"lang",
		"bbox",
	})

	// Collect lines
	var lines []*Line
	for _, sub := range el.GetElements() {
		// Todo: blockquote,display
		//switch class {
		//case "display": break;
		//case "blockquote": break;
		//case "line": break;
		//}
		if sub.IsLine() {
			lines = append(lines, NewLine(sub))
		} else {
			fmt.Println("Found new paragraph element:", sub)
		}
	}
	return &Paragraph{
		Element: Element{
			Id:          el.GetId(),
			Class:       el.GetClass(),
			BoundingBox: el.GetBoundingBox(),
		},
		Language: el.Language,
		Lines:    lines,
	}
}
