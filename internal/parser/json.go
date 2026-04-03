package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// JSONParser JSON 文件解析器
type JSONParser struct {
	*BaseParser
}

// NewJSONParser 创建 JSON 解析器
func NewJSONParser() *JSONParser {
	return &JSONParser{
		BaseParser: NewBaseParser("JSON", []string{".json"}),
	}
}

// Parse 从 io.Reader 解析 JSON
func (p *JSONParser) Parse(reader io.Reader) (*ParseResult, error) {
	decoder := json.NewDecoder(reader)

	var raw interface{}
	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// 复用 YAML 的树构建逻辑（JSON 和 YAML 的数据结构相同）
	root := buildConfigTree("root", raw)

	return &ParseResult{
		Root: root,
		Raw:  raw,
	}, nil
}

// ParseFile 从文件解析 JSON
func (p *JSONParser) ParseFile(path string) (*ParseResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return p.Parse(file)
}