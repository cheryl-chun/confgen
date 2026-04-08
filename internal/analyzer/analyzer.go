package analyzer

import (
	"fmt"

	"github.com/cheryl-chun/confgen/internal/parser"
)

// Analyzer 负责从 ConfigNode 树推断 Go 类型结构
type Analyzer struct {
	result        *AnalyzeResult
	structCounter map[string]int // 用于生成唯一的 struct 名称
}

// NewAnalyzer 创建分析器
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		result:        NewAnalyzeResult(),
		structCounter: make(map[string]int),
	}
}

// Analyze 分析 ConfigNode 树，生成类型定义
func (a *Analyzer) Analyze(root *parser.ConfigNode) (*AnalyzeResult, error) {
	if root == nil {
		return nil, fmt.Errorf("root node is nil")
	}

	// 根节点必须是 Object 类型
	if !root.IsObject() {
		return nil, fmt.Errorf("root node must be an object, got %v", root.Type)
	}

	// 分析根节点，生成 Config struct
	rootStruct := a.analyzeObject("Config", root)
	a.result.RootStruct = rootStruct

	return a.result, nil
}

// analyzeObject 分析一个 Object 节点，返回 StructDef
func (a *Analyzer) analyzeObject(structName string, node *parser.ConfigNode) *StructDef {
	def := &StructDef{
		Name:   structName,
		Fields: make([]*FieldDef, 0, len(node.Children)),
	}

	// 遍历所有子节点
	for key, child := range node.Children {
		field := a.analyzeField(key, child)
		def.Fields = append(def.Fields, field)
	}

	return def
}

// analyzeField 分析一个字段节点
func (a *Analyzer) analyzeField(key string, node *parser.ConfigNode) *FieldDef {
	field := &FieldDef{
		Name:         ToFieldName(key),
		JSONTag:      key,
		YAMLTag:      key,
		MapStructTag: key,
	}

	// 根据节点类型推断字段类型
	switch {
	case node.IsPrimitive():
		// 基础类型: string, int, float64, bool
		field.Type = GoType(node.Type)

	case node.IsArray():
		// 数组类型
		field.Type = a.analyzeArrayType(key, node)

	case node.IsObject():
		// 嵌套对象，创建新的 struct
		nestedStructName := ToStructName(key)
		nestedStruct := a.analyzeObject(nestedStructName, node)

		// 添加到结果中
		a.result.AddStruct(nestedStruct)

		// 字段类型为嵌套 struct
		field.Type = nestedStructName

	default:
		// 未知类型，使用 interface{}
		field.Type = "interface{}"
	}

	return field
}

// analyzeArrayType 分析数组类型，返回 Go 数组类型字符串
func (a *Analyzer) analyzeArrayType(key string, node *parser.ConfigNode) string {
	if len(node.Items) == 0 {
		// 空数组，无法推断类型，使用 []interface{}
		return "[]interface{}"
	}

	// 取第一个元素推断类型（假设数组元素类型一致）
	firstItem := node.Items[0]

	var elemType string
	switch {
	case firstItem.IsPrimitive():
		elemType = GoType(firstItem.Type)

	case firstItem.IsObject():
		// 数组元素是对象，创建嵌套 struct
		// 例如: servers[0] -> ServerConfig
		itemStructName := ToStructName(key) // 使用字段名的单数形式
		// 如果字段名已经是复数，去掉 's' (简单处理)
		if len(itemStructName) > 7 && itemStructName[len(itemStructName)-7:] == "sConfig" {
			itemStructName = itemStructName[:len(itemStructName)-7] + "Config"
		}

		itemStruct := a.analyzeObject(itemStructName, firstItem)
		a.result.AddStruct(itemStruct)
		elemType = itemStructName

	case firstItem.IsArray():
		// 嵌套数组 (较少见)
		elemType = a.analyzeArrayType(key+"_item", firstItem)

	default:
		elemType = "interface{}"
	}

	return "[]" + elemType
}

// Analyze 是包级别的便捷函数
func Analyze(root *parser.ConfigNode) (*AnalyzeResult, error) {
	analyzer := NewAnalyzer()
	return analyzer.Analyze(root)
}
