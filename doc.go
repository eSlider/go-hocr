// Package hocr parses hOCR 1.2 files produced by Tesseract and other OCR engines.
//
// hOCR is an HTML-based format for OCR results. This package provides a structured
// document model with YAML and positioned HTML export.
//
// See https://kba.github.io/hocr-spec/1.2 for the specification.
//
// # Quick start
//
//	doc, err := hocr.ReadFile("page.hocr")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	yaml, err := doc.ToYaml()
//	html, err := doc.ToHtml(1)
//
// For low-level XML types, see the v1_2 subpackage.
package hocr
