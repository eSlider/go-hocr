package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/eslider/go-hocr/export"
	"github.com/eslider/go-hocr/legacy"
	"gopkg.in/yaml.v3"
)

func main() {
	filePath := "testdata/sample.hocr"
	format := "yaml"
	includeMeta := false
	includeWords := false

	flag.StringVar(&filePath, "input", filePath, "Input hOCR file")
	flag.StringVar(&format, "format", format, "Output format: json or yaml")
	flag.BoolVar(&includeMeta, "meta", includeMeta, "Include metadata")
	flag.BoolVar(&includeWords, "words", includeWords, "Include words")
	flag.Parse()

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("read input: %v", err)
	}

	h, err := legacy.Parse(file)
	if err != nil {
		log.Fatalf("parse hocr: %v", err)
	}

	document, err := export.NewDocument(h, includeMeta, includeWords)
	if err != nil {
		log.Fatalf("export document: %v", err)
	}

	var out []byte
	switch format {
	case "json":
		out, err = json.Marshal(document)
	default:
		out, err = yaml.Marshal(document)
	}
	if err != nil {
		log.Fatalf("marshal: %v", err)
	}

	fmt.Print(string(out))
}
