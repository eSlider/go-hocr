package hocr_test

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	hocr "github.com/eslider/go-hocr"
	"github.com/eslider/go-hocr/v1_2"
)

func lineText(l *hocr.Line) string {
	var ws []string
	for _, w := range l.Words {
		ws = append(ws, w.Text)
	}
	return strings.Join(ws, " ")
}

// TestBareWordsSynthesizeLines is the regression test for the bug where an
// ocr_par containing ocrx_word elements directly (no ocr_line wrapper) lost all
// words. Before the fix the paragraph had zero lines here.
func TestBareWordsSynthesizeLines(t *testing.T) {
	doc, err := hocr.ReadFile(filepath.Join("testdata", "bare_words.hocr"))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Pages) != 1 {
		t.Fatalf("pages = %d, want 1", len(doc.Pages))
	}
	page := doc.Pages[0]
	if len(page.Blocks) != 1 {
		t.Fatalf("blocks = %d, want 1", len(page.Blocks))
	}
	paras := page.Blocks[0].Paragraphs
	if len(paras) != 1 {
		t.Fatalf("paragraphs = %d, want 1", len(paras))
	}
	par := paras[0]
	if len(par.Lines) != 2 {
		t.Fatalf("lines = %d, want 2 (bare words must be grouped by vertical overlap)", len(par.Lines))
	}

	for i, l := range par.Lines {
		if !l.IsLine() {
			t.Fatalf("line %d: synthetic line must report IsLine()", i)
		}
	}

	// Lines are ordered top to bottom.
	if got, want := lineText(par.Lines[0]), "Example Corp Legal Form"; got != want {
		t.Fatalf("line 0 text = %q, want %q", got, want)
	}
	if got, want := lineText(par.Lines[1]), "DE00000000000000000000"; got != want {
		t.Fatalf("line 1 text = %q, want %q", got, want)
	}

	// Bounding boxes are the union of the member word boxes.
	if got, want := par.Lines[0].BoundingBox, []int{100, 10, 320, 21}; !reflect.DeepEqual(got, want) {
		t.Fatalf("line 0 bbox = %v, want %v", got, want)
	}
	if got, want := par.Lines[1].BoundingBox, []int{100, 25, 250, 35}; !reflect.DeepEqual(got, want) {
		t.Fatalf("line 1 bbox = %v, want %v", got, want)
	}

	// Words within a line are ordered left to right.
	if got, want := par.Lines[0].Words[0].Text, "Example"; got != want {
		t.Fatalf("first word = %q, want %q", got, want)
	}
	if got, want := par.Lines[0].Words[len(par.Lines[0].Words)-1].Text, "Form"; got != want {
		t.Fatalf("last word = %q, want %q", got, want)
	}
}

// TestStandardStructureRegression proves the normal carea > par > line > word
// structure still produces exactly the same lines as before the fix.
func TestStandardStructureRegression(t *testing.T) {
	doc, err := hocr.ReadFile(filepath.Join("testdata", "sample.hocr"))
	if err != nil {
		t.Fatal(err)
	}
	paras := doc.Pages[0].Blocks[0].Paragraphs
	if len(paras) != 1 {
		t.Fatalf("paragraphs = %d, want 1", len(paras))
	}
	par := paras[0]
	if len(par.Lines) != 1 {
		t.Fatalf("lines = %d, want 1", len(par.Lines))
	}
	if got, want := lineText(par.Lines[0]), "Hello World"; got != want {
		t.Fatalf("line text = %q, want %q", got, want)
	}
	if !par.Lines[0].IsLine() {
		t.Fatal("expected ocr_line to remain a real line")
	}
	if got, want := par.Lines[0].BoundingBox, []int{36, 96, 559, 120}; !reflect.DeepEqual(got, want) {
		t.Fatalf("line bbox = %v, want %v", got, want)
	}
}

// TestStructureTolerance covers missing wrappers: ocr_line directly under
// carea, and ocrx_word directly under page.
func TestStructureTolerance(t *testing.T) {
	word := &v1_2.Element{Id: "word_1_1", Class: "ocrx_word", Text: "hi", Properties: "bbox 10 10 50 20"}
	line := &v1_2.Element{Id: "line_1_1", Class: "ocr_line", Properties: "bbox 10 10 100 20"}
	line.Spans = []*v1_2.Element{word}

	carea := &v1_2.Element{Id: "block_1_1", Class: "ocr_carea"}
	carea.Spans = []*v1_2.Element{line}
	block := hocr.NewBlock(carea)
	if len(block.Paragraphs) != 1 {
		t.Fatalf("carea: paragraphs = %d, want 1", len(block.Paragraphs))
	}
	if len(block.Paragraphs[0].Lines) != 1 {
		t.Fatalf("carea: lines = %d, want 1", len(block.Paragraphs[0].Lines))
	}
	if got := lineText(block.Paragraphs[0].Lines[0]); got != "hi" {
		t.Fatalf("carea: line text = %q, want %q", got, "hi")
	}

	page := &v1_2.Element{Id: "page_1", Class: "ocr_page"}
	page.Spans = []*v1_2.Element{{Id: "word_1_1", Class: "ocrx_word", Text: "solo", Properties: "bbox 1 2 3 4"}}
	pg := hocr.NewPage(page)
	if len(pg.Blocks) != 1 || !pg.Blocks[0].IsContentArea() {
		t.Fatalf("page: blocks = %d, want 1 content area", len(pg.Blocks))
	}
	if len(pg.Blocks[0].Paragraphs) != 1 || len(pg.Blocks[0].Paragraphs[0].Lines) != 1 {
		t.Fatalf("page: want 1 paragraph with 1 line")
	}
	if got := lineText(pg.Blocks[0].Paragraphs[0].Lines[0]); got != "solo" {
		t.Fatalf("page: line text = %q, want %q", got, "solo")
	}
}

// TestBrokenStructuresNoPanic ensures malformed input (word without bbox,
// unknown elements) does not panic and does not silently drop known content.
func TestBrokenStructuresNoPanic(t *testing.T) {
	// A word with no bbox at all.
	par := &v1_2.Element{Id: "par_1_1", Class: "ocr_par"}
	par.Spans = []*v1_2.Element{{Id: "word_1_1", Class: "ocrx_word", Text: "orphan"}}
	p := hocr.NewParagraph(par)
	if len(p.Lines) != 1 {
		t.Fatalf("no-bbox word: lines = %d, want 1", len(p.Lines))
	}
	if got := lineText(p.Lines[0]); got != "orphan" {
		t.Fatalf("no-bbox word: text = %q, want %q", got, "orphan")
	}
	if got := p.GetText(); !strings.Contains(got, "orphan") {
		t.Fatalf("no-bbox word: GetText lost the text: %q", got)
	}

	// An unknown element under page and under carea must be ignored, not fatal.
	page := &v1_2.Element{Id: "page_1", Class: "ocr_page"}
	page.Spans = []*v1_2.Element{{Id: "weird_1_1", Class: "weird"}}
	pg := hocr.NewPage(page)
	if len(pg.Blocks) != 0 {
		t.Fatalf("unknown page child: blocks = %d, want 0", len(pg.Blocks))
	}

	carea := &v1_2.Element{Id: "block_1_1", Class: "ocr_carea"}
	carea.Spans = []*v1_2.Element{{Id: "weird_1_1", Class: "weird"}}
	block := hocr.NewBlock(carea)
	if len(block.Paragraphs) != 0 {
		t.Fatalf("unknown carea child: paragraphs = %d, want 0", len(block.Paragraphs))
	}
}
