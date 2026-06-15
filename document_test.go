package hocr_test

import (
	"os"
	"path/filepath"
	"testing"

	hocr "github.com/eslider/go-hocr"
	"github.com/eslider/go-hocr/legacy"
)

func TestReadFile(t *testing.T) {
	path := filepath.Join("testdata", "sample.hocr")
	doc, err := hocr.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if doc.System == "" {
		t.Fatal("expected ocr-system meta")
	}
	if len(doc.Pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(doc.Pages))
	}
	if doc.Pages[0].Image != "\"sample.png\"" {
		t.Fatalf("unexpected image path: %q", doc.Pages[0].Image)
	}

	html, err := doc.ToHtml(1)
	if err != nil {
		t.Fatal(err)
	}
	if html == "" {
		t.Fatal("expected non-empty HTML")
	}

	yaml, err := doc.ToYaml()
	if err != nil {
		t.Fatal(err)
	}
	if yaml == "" {
		t.Fatal("expected non-empty YAML")
	}
}

func TestLegacyParse(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("testdata", "sample.hocr"))
	if err != nil {
		t.Fatal(err)
	}

	h, err := legacy.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(h.Pages) == 0 {
		t.Fatal("expected pages")
	}
}
