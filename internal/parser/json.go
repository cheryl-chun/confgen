package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type JSONParser struct {
	*BaseParser
}

func NewJSONParser() *JSONParser {
	return &JSONParser{
		BaseParser: NewBaseParser("JSON", []string{".json"}),
	}
}

func (p *JSONParser) Parse(reader io.Reader) (*ParseResult, error) {
	decoder := json.NewDecoder(reader)

	var raw interface{}
	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	root := buildConfigTree("root", raw)

	return &ParseResult{
		Root: root,
		Raw:  raw,
	}, nil
}

func (p *JSONParser) ParseFile(path string) (*ParseResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return p.Parse(file)
}