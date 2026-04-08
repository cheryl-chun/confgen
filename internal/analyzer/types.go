package analyzer

import "github.com/cheryl-chun/confgen/internal/parser"

// StructDef 表示一个 Go struct 的定义
type StructDef struct {
	Name   string      // struct 名称 (e.g., "ServerConfig")
	Fields []*FieldDef // 字段列表
}

// FieldDef 表示 struct 中的一个字段
type FieldDef struct {
	Name         string    // 字段名 (PascalCase, e.g., "MaxConnections")
	Type         string    // Go 类型 (e.g., "string", "int", "[]string", "*DatabaseConfig")
	JSONTag      string    // JSON tag 值 (e.g., "max_connections")
	YAMLTag      string    // YAML tag 值
	MapStructTag string    // mapstructure tag 值
	Comment      string    // 字段注释 (可选)
}

// AnalyzeResult 表示分析结果
type AnalyzeResult struct {
	RootStruct *StructDef            // 根配置结构体 (通常命名为 "Config")
	SubStructs map[string]*StructDef // 所有嵌套的 struct 定义，key 是 struct 名称
}

// NewAnalyzeResult 创建分析结果
func NewAnalyzeResult() *AnalyzeResult {
	return &AnalyzeResult{
		SubStructs: make(map[string]*StructDef),
	}
}

// AddStruct 添加一个 struct 定义
func (r *AnalyzeResult) AddStruct(def *StructDef) {
	r.SubStructs[def.Name] = def
}

// GoType 从 parser.ValueType 推断 Go 类型
func GoType(vt parser.ValueType) string {
	switch vt {
	case parser.TypeString:
		return "string"
	case parser.TypeInt:
		return "int"
	case parser.TypeFloat:
		return "float64"
	case parser.TypeBool:
		return "bool"
	case parser.TypeNull:
		return "interface{}" // null 值默认使用 interface{}
	default:
		return "interface{}"
	}
}
