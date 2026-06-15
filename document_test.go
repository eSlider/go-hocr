package hocr_test

import (
	"path/filepath"
	"testing"

	hocr "github.com/eslider/go-hocr"
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
