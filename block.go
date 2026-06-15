package hocr

import (
	"fmt"
	"github.com/eslider/go-hocr/v1_2"
	"strings"
)

const (
	SeparatorBlock   string = "separator"
	ContentAreaBlock        = "carea"
	PhotoBlock              = "photo"
)

type Block struct {
	Element    `yaml:",inline"`
	Paragraphs []*Paragraph
}

func (b *Block) IsSeparator() bool {
	return b.Class == SeparatorBlock
}
func (b *Block) IsContentArea() bool {
	return b.Class == ContentAreaBlock
}

func (b *Block) IsPhoto() bool {
	return b.Class == PhotoBlock
}

func (b *Block) GetHtml() string {
	var ps []string
	if b.IsSeparator() {
		l := b.BoundingBox[0] / scaleFactor
		t := b.BoundingBox[1] / scaleFactor
		w := (b.BoundingBox[2] - b.BoundingBox[0]) / scaleFactor
		h := (b.BoundingBox[3] - b.BoundingBox[1]) / scaleFactor
		return fmt.Sprintf(`<hr style="left: %dpx;top: %dpx;width: %dpx;height: %dpx;"/>`, l, t, w, h)
	}
	if b.IsPhoto() {
		return "<img/>"
	}
	if b.IsContentArea() {
		for _, p := range b.Paragraphs {
			ps = append(ps, fmt.Sprintf(`<div>%s</div>`, p.GetText()))
		}
	} else {
		return ""
	}

	if ps == nil {
		return ""
	}

	return strings.Join(ps, "\n")
}

func NewBlock(el *v1_2.Element) *Block {
	// Check for new fields
	checkForNewProperties(el, []string{
		"bbox",
	})

	// Content area block has at least one element and normally an 'p' element
	var paragraphs []*Paragraph

	if el.IsContentArea() {
		for _, sub := range el.GetElements() {
			if sub.IsParagraph() {
				paragraphs = append(paragraphs, NewParagraph(sub))
			} else {
				fmt.Println("Found not a paragraph in block:", sub)
			}
		}
	} else if el.IsPhoto() {
	} else if el.IsSeparator() {
	} else {
		fmt.Println("Found something else in a block:", el)
	}
	// Separator has no content (normally) but interesting to make space between
	//if b.IsSeparator() {
	//	return b
	//}

	return &Block{
		Element: Element{
			Id:          el.GetId(),
			Class:       el.GetClass(),
			BoundingBox: el.GetBoundingBox(),
		},
		Paragraphs: paragraphs,
	}

}
