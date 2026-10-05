package hocr

import (
	"fmt"
	"sort"
	"strings"

	"github.com/eslider/go-hocr/v1_2"
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
//
// Normally an ocr_par contains ocr_line elements. Some Tesseract PSM modes
// (e.g. PSM 3) emit ocrx_word elements directly under ocr_par without an
// ocr_line wrapper. Those bare words are grouped into synthesised lines by
// vertical bbox overlap so that no content is lost.
func NewParagraph(el *v1_2.Element) *Paragraph {
	// Check for new fields
	checkForNewProperties(el, []string{
		"lang",
		"bbox",
	})

	// Collect lines and bare words (ocrx_word directly under ocr_par).
	var lines []*Line
	var words []*Word
	for _, sub := range el.GetElements() {
		switch {
		case sub.IsLine():
			lines = append(lines, NewLine(sub))
		case sub.IsWord():
			words = append(words, NewWord(sub))
		default:
			fmt.Println("Found new paragraph element:", sub)
		}
	}
	if len(words) > 0 {
		lines = append(lines, synthesizeLines(words)...)
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

// newSyntheticParagraph builds a paragraph from already-parsed lines and bare
// words. It is used when the hOCR structure skips the expected ocr_par wrapper:
// a direct ocr_line under a carea/page becomes a paragraph, and bare ocrx_word
// elements are wrapped into synthesised lines.
func newSyntheticParagraph(lines []*Line, words []*Word) *Paragraph {
	if len(words) > 0 {
		lines = append(lines, synthesizeLines(words)...)
	}
	return &Paragraph{
		Element: Element{Class: "par"},
		Lines:   lines,
	}
}

// synthesizeLines groups words into visual lines by vertical bbox overlap.
// Words whose vertical intervals [y0,y1] intersect belong to the same line;
// words from different visual lines of one paragraph become separate Lines.
// Words within a line are ordered left to right (by x0). The resulting line has
// Class "line" so that Line.IsLine reports true, and its BoundingBox is the
// union of the word boxes.
func synthesizeLines(words []*Word) []*Line {
	sorted := make([]*Word, len(words))
	copy(sorted, words)
	sort.SliceStable(sorted, func(i, j int) bool {
		xi, yi, _, _, oki := wordBBox(sorted[i])
		xj, yj, _, _, okj := wordBBox(sorted[j])
		if !oki || !okj {
			return oki && !okj // words without a usable bbox go last
		}
		if yi != yj {
			return yi < yj
		}
		return xi < xj
	})

	type group struct {
		y0, y1 int
		valid  bool
		words  []*Word
	}
	var groups []*group
	for _, w := range sorted {
		_, y0, _, y1, ok := wordBBox(w)
		if !ok {
			groups = append(groups, &group{words: []*Word{w}})
			continue
		}
		placed := false
		for _, g := range groups {
			if !g.valid {
				continue
			}
			if y0 < g.y1 && g.y0 < y1 { // vertical intervals overlap
				g.words = append(g.words, w)
				if y0 < g.y0 {
					g.y0 = y0
				}
				if y1 > g.y1 {
					g.y1 = y1
				}
				placed = true
				break
			}
		}
		if !placed {
			groups = append(groups, &group{y0: y0, y1: y1, valid: true, words: []*Word{w}})
		}
	}

	lines := make([]*Line, 0, len(groups))
	for _, g := range groups {
		// order words left to right
		sort.SliceStable(g.words, func(i, j int) bool {
			xi, _, _, _, oki := wordBBox(g.words[i])
			xj, _, _, _, okj := wordBBox(g.words[j])
			if !oki || !okj {
				return oki && !okj
			}
			return xi < xj
		})
		lines = append(lines, &Line{
			Element: Element{
				Class:       TextLine,
				BoundingBox: unionBBox(g.words),
			},
			Words: g.words,
		})
	}
	return lines
}

// wordBBox returns the four bbox values of a word and whether they are usable.
func wordBBox(w *Word) (x0, y0, x1, y1 int, ok bool) {
	b := w.BoundingBox
	if len(b) < 4 {
		return 0, 0, 0, 0, false
	}
	return b[0], b[1], b[2], b[3], true
}

// unionBBox returns the smallest bbox enclosing all words that have a bbox.
func unionBBox(words []*Word) []int {
	var box []int
	for _, w := range words {
		b := w.BoundingBox
		if len(b) < 4 {
			continue
		}
		if box == nil {
			box = []int{b[0], b[1], b[2], b[3]}
			continue
		}
		if b[0] < box[0] {
			box[0] = b[0]
		}
		if b[1] < box[1] {
			box[1] = b[1]
		}
		if b[2] > box[2] {
			box[2] = b[2]
		}
		if b[3] > box[3] {
			box[3] = b[3]
		}
	}
	return box
}
