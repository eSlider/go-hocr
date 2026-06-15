package v1_2

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// Document specs: https://kba.github.io/hocr-spec/1.2
// Supports: v1.2
type Document struct {
	// Title is normally not set, but maybe needed?
	Title string `xml:"head>title"`

	// Meta info contains pre-pared values to get unmarshalled and late parsed into Properties
	Meta []*struct {
		Name    string `xml:"name,attr"`
		Content string `xml:"content,attr"`
	} `xml:"head>meta"`

	// Pages hold elements as a tree
	Pages []*Element `xml:"body>div"`

	// Properties used to save parsed meta data
	Properties string
}

// Element is a node with sub-nodes represents every type in document body
// Example <div class='ocr_carea' id='block_1_14' title="bbox 773 4305 2451 4802"
type Element struct {
	// Name of the element, such as "p","div" or"span"
	Name xml.Name `xml:"name"`

	// Id numeric, not original, but never need,
	// course using of slices preserve save order
	Id string `xml:"id,attr"`

	// Class has only value without "ocr_" prefix
	Class string `xml:"class,attr"`

	// Language of only lines(p) like "deu","eng","rus"
	Language string `xml:"lang,attr"`

	// Elements holds div-blocks sub elements same as this one but in slice
	Elements []*Element `xml:"div"`

	// Paragraphs holds only paragraphs-lines
	Paragraphs []*Element `xml:"p"`

	// Paragraphs holds only span-words
	Spans []*Element `xml:"span"`

	// Text of span-words
	Text string `xml:",chardata"`

	// Data is same as `Text` but from XML CDATA
	Data string `xml:",cdata"`

	// Properties contents unparsed raw-string of title-attribute
	Properties string `xml:"title,attr"`

	// Need to parse Properties only one time
	meta map[string]string // For caching
}

// SplitMeta content from attributes
func SplitMeta(meta string) map[string]string {
	a := strings.Split(meta, ";")
	r := make(map[string]string)
	for _, s := range a {
		s = strings.Trim(s, " ")
		kv := strings.Split(s, " ")
		r[kv[0]] = strings.Trim(strings.TrimPrefix(s, kv[0]), " ")
	}
	return r
}

// GetProperties parses and caches Properties for often future use
func (h *Element) GetProperties() map[string]string {
	if h.meta == nil {
		h.meta = SplitMeta(h.Properties)
	}
	return h.meta
}

// GetElements of all types
func (h *Element) GetElements() []*Element {
	var r []*Element
	r = append(r, h.Elements...)
	r = append(r, h.Paragraphs...)
	r = append(r, h.Spans...)
	return r
}

// ParseValue of attribute by given key and postion.
func (h *Element) ParseValue(k string, pos int) string {
	count := strings.Count(k, "_")
	if count < 1 {
		return ""
	}
	if pos < 0 {
		pos = count * -pos
	}
	return strings.Split(k, "_")[pos]
}

// GetType in short.
func (h *Element) GetType() string {
	return h.ParseValue(h.Id, 0)
}

// GetClass in short from Class.
func (h *Element) GetClass() string {
	return h.ParseValue(h.Class, 1)
}

// GetText from Text or Data.
func (h *Element) GetText() string {
	if h.Text != "" {
		return h.Text
	}
	if h.Data != "" {
		return h.Data
	}
	return ""
}

// GetId in short.
func (h *Element) GetId() int16 {
	id, _ := strconv.Atoi(h.ParseValue(h.Id, -1))
	return int16(id)
}

// IsBlock the element?
func (h *Element) IsBlock() bool {
	return h.GetType() == "block"
}

// IsContentArea the element?
func (h *Element) IsContentArea() bool {
	return h.GetClass() == "carea"
}

// GetBoundingBox returns the bounding box of the hOCR element.
func (h *Element) GetBoundingBox() []int {
	var r []int
	ps := h.GetProperties()

	if ps["bbox"] == "" {
		return r
	}

	for _, v := range strings.Split(ps["bbox"], " ") {
		f, _ := strconv.ParseInt(v, 10, 32)
		r = append(r, int(f))
	}
	return r
}

// IsParagraph ?
func (h *Element) IsParagraph() bool {
	return h.GetType() == "par"
}

// IsLine ?
func (h *Element) IsLine() bool {
	return h.GetType() == "line"

}

// GetPropF32 - catch property by given name and convert to float32 value.
func (h *Element) GetPropF32(k string) float32 {
	var r float32
	props := h.GetProperties()
	if props[k] != "" {
		f, _ := strconv.ParseFloat(props[k], 32)
		r = float32(f)
	}

	return r
}

// GetPropF32Slice - catch properties by given name and convert to float32 values.
func (h *Element) GetPropF32Slice(k string) []float32 {
	var r []float32
	props := h.GetProperties()

	if props[k] != "" {
		for _, v := range strings.Split(props[k], " ") {
			f, _ := strconv.ParseFloat(v, 32)
			r = append(r, float32(f))
		}
	}

	return r
}

// IsWord the element?
func (h *Element) IsWord() bool {
	return h.GetType() == "word"
}

// IsPhoto (image)?
func (h *Element) IsPhoto() bool {
	return h.GetClass() == "photo"
}

// GetPropString - catch property by given name and convert to string value.
func (h *Element) GetPropString(k string) string {
	props := h.GetProperties()
	if props[k] == "" {
		return ""
	}

	return props[k]
}

// GetPropInt16 - catch property by given name and convert to int value
func (h *Element) GetPropInt16(s string) int16 {
	props := h.GetProperties()
	if props[s] == "" {
		return 0
	}

	// Get the page number
	id, _ := strconv.Atoi(props[s])
	return int16(id)
}

// GetPropInt16Slice - catch property by given name and convert to int value slice.
func (h *Element) GetPropInt16Slice(k string) []int16 {
	var r []int16
	props := h.GetProperties()

	if props[k] != "" {
		for _, v := range strings.Split(props[k], " ") {
			f, _ := strconv.ParseInt(v, 10, 16)
			r = append(r, int16(f))
		}
	}

	return r
}

// IsSeparator the element?
func (h *Element) IsSeparator() bool {
	return h.GetClass() == "separator"
}

// GetCapabilities returns the capabilities of the hOCR document.
// Returns a slice of strings:
//
//	"ocr_page"
//	"ocr_carea"
//	"ocr_par"
//	"ocr_line"
//	"ocrx_word"
//	"ocrp_wconf"
func (h *Document) GetCapabilities() []string {
	for _, m := range h.Meta {
		if m.Name == "ocr-capabilities" {
			return strings.Split(m.Content, " ")
		}
	}
	return nil
}

// GetSystem of OCR as "name version".
func (h *Document) GetSystem() string {
	for _, m := range h.Meta {
		if m.Name == "ocr-system" {
			return m.Content
		}
	}
	return ""
}

// GetProperties of document by parsing. Not cached.
func (h *Document) GetProperties() interface{} {
	return SplitMeta(h.Properties)
}
