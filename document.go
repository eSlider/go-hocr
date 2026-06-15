package hocr

import (
	"encoding/xml"
	"fmt"
	"github.com/eslider/go-hocr/v1_2"
	"gopkg.in/yaml.v3"
	"os"
)

// Document is uniformed structure to get useful data from specification hOCR standard.
type Document struct {
	Title        string
	System       string
	Capabilities []string `yaml:"capabilities,flow"`
	Pages        []*Page
}

// ToYaml export YAML as string.
func (d *Document) ToYaml() (string, error) {
	jsonBytes, err := yaml.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// ToHtml export HTML as string.
func (d *Document) ToHtml(pageNr int) (string, error) {
	if d.Pages == nil || len(d.Pages) < pageNr-1 {
		return "", fmt.Errorf("no pages found")
	}

	//html := fmt.Sprintf(`<div>%s</div>`, d.Pages[pageNr-1].GetHtml())
	return d.Pages[pageNr-1].GetHtml(), nil
}

func ReadFile(path string) (*Document, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	doc := new(v1_2.Document)
	err = xml.Unmarshal(b, doc)
	if err != nil {
		return nil, err
	}
	return NewDocument(doc), err
}

func NewDocument(doc *v1_2.Document) *Document {
	d := new(Document)
	d.Title = doc.Title
	d.Capabilities = doc.GetCapabilities()
	d.System = doc.GetSystem()
	for _, page := range doc.Pages {
		d.Pages = append(d.Pages, NewPage(page))
	}
	return d
}
