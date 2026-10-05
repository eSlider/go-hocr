package hocr

import (
	"fmt"
	"github.com/eslider/go-hocr/v1_2"
	"strings"
)

const (
	PageImageProperty          = "image"    // Image path
	PageNumberProperty         = "ppageno"  // Page number
	PageScanResolutionProperty = "scan_res" // Scan resolution
)

type Page struct {
	// Element holds universal attributes
	Element `yaml:",inline"`

	// Image string path of original file sent to  OCR.
	// * image file name used as input
	// * syntactically, must be a UNIX-like pathname or http URL (no Windows path names)
	// * may be relative
	// * if the hOCR file is present in a directory hierarchy or file archive, should resolve to the corresponding image file
	Image string `yaml:"path"`

	// Number of Page
	Number int16 `yaml:"number"`

	// ScanResolution in DPI
	ScanResolution []int16 `yaml:"scan_res,flow"`

	// Blocks of Page
	Blocks []*Block
}

// NewPage by given specification v1_2.Element
func NewPage(page *v1_2.Element) *Page {
	// A well-formed page holds ocr_carea blocks. Some engines/PSM modes skip
	// the carea wrapper and place ocr_par, ocr_line or ocrx_word directly under
	// ocr_page; wrap those into an equivalent content-area block.
	var blocks []*Block
	var orphanParagraphs []*Paragraph
	var orphanLines []*Line
	var orphanWords []*Word
	for _, sub := range page.GetElements() {
		switch {
		case sub.IsBlock():
			blocks = append(blocks, NewBlock(sub))
		case sub.IsParagraph():
			orphanParagraphs = append(orphanParagraphs, NewParagraph(sub))
		case sub.IsLine():
			orphanLines = append(orphanLines, NewLine(sub))
		case sub.IsWord():
			orphanWords = append(orphanWords, NewWord(sub))
		default:
			fmt.Println("Found not a block in page", sub)
		}
	}
	if len(orphanParagraphs) > 0 || len(orphanLines) > 0 || len(orphanWords) > 0 {
		block := &Block{
			Element:    Element{Class: ContentAreaBlock},
			Paragraphs: orphanParagraphs,
		}
		if len(orphanLines) > 0 || len(orphanWords) > 0 {
			block.Paragraphs = append(block.Paragraphs, newSyntheticParagraph(orphanLines, orphanWords))
		}
		blocks = append(blocks, block)
	}

	return &Page{
		Element: Element{
			Id:          page.GetId(),
			Class:       page.GetClass(),
			BoundingBox: page.GetBoundingBox(),
		},
		Image:          page.GetPropString(PageImageProperty),
		Number:         page.GetPropInt16(PageNumberProperty),
		ScanResolution: page.GetPropInt16Slice(PageScanResolutionProperty),
		Blocks:         blocks,
	}
}

// GetHtml of Page as string
func (p *Page) GetHtml() string {
	var bs []string
	for _, b := range p.Blocks {
		html := b.GetHtml()
		bs = append(bs, html)
	}

	if bs == nil {
		return ""
	}

	return strings.Join(bs, "\n")
}
