package hocr

import (
	"fmt"
	"github.com/eslider/go-hocr/v1_2"
)

type Word struct {
	Element    `yaml:",inline"`
	Text       string
	Confidence float32 // OCR-engine specific confidence for the entire contained substring
}

func (w Word) GetText() string {
	left := w.BoundingBox[0] / scaleFactor
	top := w.BoundingBox[1] / scaleFactor
	return fmt.Sprintf( /** @lang HTML */ `<div style="left: %dpx; top:%dpx">%s</div>`, left, top, w.Text)
}

func NewWord(el *v1_2.Element) *Word {
	// Check for new fields
	checkForNewProperties(el, []string{
		"bbox",
		"x_wconf",
	})
	return &Word{
		Text:       el.GetText(),
		Confidence: el.GetPropF32("x_wconf"),
		Element: Element{
			Id:          el.GetId(),
			Class:       el.GetClass(),
			BoundingBox: el.GetBoundingBox(),
		},
	}
}
