package hocr

import (
	"fmt"
	"github.com/eslider/go-hocr/v1_2"
	"strings"
)

const (
	TextFloatLine string = "textfloat"
	TextLine             = "line"
	HeaderLine           = "header"
	CaptionLine          = "caption"
)

type Line struct {
	Element    `yaml:",inline"`
	BaseLine   []float32 `yaml:"baseline,flow"`
	Size       float32
	Ascenders  float32
	Descenders float32
	Words      []*Word `yaml:"words"`

	// Textangle - The angle in degrees by which textual content has been rotated relative to the rest of the page
	// (if not present, the angle is assumed to be zero); rotations are counter-clockwise,
	// so an angle of 90 degrees is vertical text running from bottom to top in Latin script;
	// note that this is different from reading order, which should be indicated using standard HTML properties
	Textangle float32
}

func (l *Line) GetHtml() string {
	var ws []string

	for _, w := range l.Words {
		ws = append(ws, w.GetText())
	}

	if ws == nil {
		return "<br/>"
	}

	tags := strings.Join(ws, " ")
	if l.IsHeader() {
		return "<h1>" + tags + "</h1>"
	}
	if l.IsCaption() {
		return "<b>" + tags + "</b>"
	}
	if !l.IsLine() && !l.IsTextFloat() {
		return "<h1>" + tags + "</h1>"
	}

	return tags
}

func (l *Line) IsTextFloat() bool {
	return l.Class == TextFloatLine
}

func (l *Line) IsLine() bool {
	return l.Class == TextLine

}

func (l *Line) IsHeader() bool {
	return l.Class == HeaderLine

}

func (l *Line) IsCaption() bool {
	return l.Class == CaptionLine
}

func NewLine(el *v1_2.Element) *Line {
	// Check for new fields
	checkForNewProperties(el, []string{
		"baseline",
		"x_size",
		"x_descenders",
		"x_ascenders",
		"bbox",
		"textangle",
	})

	// Get words
	var words []*Word
	for _, sub := range el.GetElements() {
		if sub.IsWord() {
			words = append(words, NewWord(sub))
		} else {
			fmt.Println("Found new line element:", sub)
		}

	}
	return &Line{
		Element: Element{
			Id:          el.GetId(),
			Class:       el.GetClass(),
			BoundingBox: el.GetBoundingBox(),
		},
		Textangle:  el.GetPropF32("textangle"),
		BaseLine:   el.GetPropF32Slice("baseline"),
		Descenders: el.GetPropF32("x_descenders"),
		Ascenders:  el.GetPropF32("x_ascenders"),
		Size:       el.GetPropF32("x_size"),
		Words:      words,
	}
}
