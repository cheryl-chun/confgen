package parser

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type YAMLParser struct {
	*BaseParser
}

var _ Parser = (*YAMLParser)(nil) // Ensure YAMLParser implements Parser interface

func NewYAMLParser() *YAMLParser {
	return &YAMLParser{
		BaseParser: NewBaseParser("YAML", []string{".yaml", ".yml"}),
	}
}

// Parse parse YAML from io.Reader
func (p *YAMLParser) Parse(reader io.Reader) (*ParseResult, error) {
	decoder := yaml.NewDecoder(reader)

	var raw interface{}
	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}

	root := buildConfigTree("root", raw)

	return &ParseResult{
		Root: root,
		Raw:  raw,
	}, nil
}

// ParseFile parse YAML from file
func (p *YAMLParser) ParseFile(path string) (*ParseResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return p.Parse(file)
}

// buildConfigTree 递归构建配置树
func buildConfigTree(key string, value interface{}) *ConfigNode {
	node := NewConfigNode(key)

	switch v := value.(type) {
	case map[string]interface{}:
		// 对象类型
		node.Type = TypeObject
		node.Value = v
		for k, val := range v {
			child := buildConfigTree(k, val)
			node.AddChild(child)
		}

	case []interface{}:
		// 数组类型
		node.Type = TypeArray
		node.Value = v
		for i, item := range v {
			itemNode := buildConfigTree(fmt.Sprintf("[%d]", i), item)
			node.AddItem(itemNode)
		}

	case string:
		node.Type = TypeString
		node.Value = v

	case int, int8, int16, int32, int64:
		node.Type = TypeInt
		node.Value = v

	case float32, float64:
		node.Type = TypeFloat
		node.Value = v

	case bool:
		node.Type = TypeBool
		node.Value = v

	case nil:
		node.Type = TypeNull
		node.Value = nil

	default:
		// 未知类型，尝试转换为字符串
		node.Type = TypeString
		node.Value = fmt.Sprint(v)
	}

	return node
}