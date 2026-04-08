package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cheryl-chun/confgen/internal/analyzer"
	"github.com/cheryl-chun/confgen/internal/parser"
)

// TestFullPipeline 完整测试从 YAML → Parser → Analyzer 的流程
//
// 这个测试展示了代码生成的完整数据流转：
//
// 步骤 1: YAML 文本 (用户的配置文件)
//    ↓
// 步骤 2: parser.ConfigNode 树 (临时中间表示)
//    ↓
// 步骤 3: analyzer.AnalyzeResult (Go 类型定义)
//    ↓
// 步骤 4: (下一步) codegen 生成 Go 源码
func TestFullPipeline(t *testing.T) {
	// ========================================
	// 步骤 1: 准备测试用的 YAML 配置文件
	// ========================================
	yamlContent := `
server:
  host: "localhost"
  port: 8080
  timeout: 30.5
  enabled: true
  tags:
    - "production"
    - "web"

database:
  host: "db.example.com"
  port: 5432
  max_connections: 100
  ssl_enabled: true

  # 嵌套对象
  credentials:
    username: "admin"
    password: "secret"

# 对象数组
servers:
  - host: "server1.example.com"
    port: 8080
  - host: "server2.example.com"
    port: 8081
`

	// 创建临时文件
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	t.Logf("✓ 步骤 1: 创建测试配置文件: %s", configPath)

	// ========================================
	// 步骤 2: 使用 Parser 解析 YAML 文件
	// ========================================
	//
	// Parser 的作用：
	// - 读取 YAML/JSON 文件
	// - 将文本解析为 ConfigNode 树
	// - ConfigNode 是一个临时的中间表示，包含：
	//   * Key: 字段名 (如 "host", "port")
	//   * Value: 实际值
	//   * Type: 值类型 (TypeString, TypeInt, TypeObject, TypeArray 等)
	//   * Children: 子节点 (用于 Object)
	//   * Items: 数组元素 (用于 Array)
	//
	t.Log("\n--- 步骤 2: Parser 解析 YAML ---")

	parseResult, err := parser.ParseFile(configPath)
	if err != nil {
		t.Fatalf("Parser.ParseFile failed: %v", err)
	}

	t.Logf("✓ 解析成功，根节点类型: %v", parseResult.Root.Type)
	t.Logf("✓ 根节点包含 %d 个字段", len(parseResult.Root.Children))

	// 查看解析后的树结构
	printConfigNode(t, parseResult.Root, 0)

	// ========================================
	// 步骤 3: 使用 Analyzer 推断 Go 类型
	// ========================================
	//
	// Analyzer 的作用：
	// - 遍历 ConfigNode 树
	// - 推断每个节点对应的 Go 类型
	// - 生成 struct 定义 (StructDef)
	// - 处理嵌套结构、数组、命名规范等
	//
	// Analyzer 的核心逻辑：
	// 1. 对于 Object 节点 → 生成一个 struct
	// 2. 对于 Array 节点 → 推断数组元素类型 (如 []string, []ServerConfig)
	// 3. 对于 Primitive 节点 → 映射为 Go 基础类型 (string, int, bool, float64)
	// 4. 字段命名转换: snake_case → PascalCase (max_connections → MaxConnections)
	// 5. struct 命名: 字段名 + "Config" 后缀 (server → ServerConfig)
	//
	t.Log("\n--- 步骤 3: Analyzer 推断类型 ---")

	analyzeResult, err := analyzer.Analyze(parseResult.Root)
	if err != nil {
		t.Fatalf("Analyzer.Analyze failed: %v", err)
	}

	t.Logf("✓ 类型推断成功")
	t.Logf("✓ 生成了 1 个根 struct: %s", analyzeResult.RootStruct.Name)
	t.Logf("✓ 生成了 %d 个嵌套 struct", len(analyzeResult.SubStructs))

	// 打印生成的类型定义
	printAnalyzeResult(t, analyzeResult)

	// ========================================
	// 验证生成的结构
	// ========================================
	t.Log("\n--- 验证生成的类型定义 ---")

	// 验证根 struct
	if analyzeResult.RootStruct.Name != "Config" {
		t.Errorf("RootStruct.Name = %q, want %q", analyzeResult.RootStruct.Name, "Config")
	}

	// 验证根 struct 的字段数量
	if len(analyzeResult.RootStruct.Fields) != 3 {
		t.Errorf("len(RootStruct.Fields) = %d, want 3 (server, database, servers)",
			len(analyzeResult.RootStruct.Fields))
	}

	// 验证嵌套 struct
	expectedStructs := []string{"ServerConfig", "DatabaseConfig", "CredentialsConfig"}
	for _, name := range expectedStructs {
		if _, ok := analyzeResult.SubStructs[name]; !ok {
			t.Errorf("Missing expected struct: %s", name)
		}
	}

	// 验证具体的字段类型
	serverStruct := analyzeResult.SubStructs["ServerConfig"]
	if serverStruct != nil {
		// 验证字段
		expectedFields := map[string]string{
			"Host":    "string",
			"Port":    "int",
			"Timeout": "float64",
			"Enabled": "bool",
			"Tags":    "[]string",
		}

		for _, field := range serverStruct.Fields {
			expectedType, ok := expectedFields[field.Name]
			if !ok {
				t.Logf("Warning: Unexpected field in ServerConfig: %s", field.Name)
				continue
			}
			if field.Type != expectedType {
				t.Errorf("ServerConfig.%s: Type = %q, want %q",
					field.Name, field.Type, expectedType)
			}
		}
	}

	// 验证数组对象
	for _, field := range analyzeResult.RootStruct.Fields {
		if field.Name == "Servers" {
			if field.Type != "[]ServerConfig" {
				t.Errorf("Servers field: Type = %q, want %q",
					field.Type, "[]ServerConfig")
			}
			t.Logf("✓ 正确识别了对象数组类型: %s", field.Type)
		}
	}

	t.Log("\n✓ 所有验证通过！")
	t.Log("\n下一步: codegen 模块将生成最终的 Go 源码")
}

// printConfigNode 递归打印 ConfigNode 树结构（用于调试）
func printConfigNode(t *testing.T, node *parser.ConfigNode, indent int) {
	if node == nil {
		return
	}

	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	switch {
	case node.IsObject():
		t.Logf("%s%s (Object) {", prefix, node.Key)
		for _, child := range node.Children {
			printConfigNode(t, child, indent+1)
		}
		t.Logf("%s}", prefix)

	case node.IsArray():
		t.Logf("%s%s (Array) [%d items]", prefix, node.Key, len(node.Items))
		if len(node.Items) > 0 {
			// 只打印第一个元素
			printConfigNode(t, node.Items[0], indent+1)
			if len(node.Items) > 1 {
				t.Logf("%s  ... (%d more items)", prefix, len(node.Items)-1)
			}
		}

	case node.IsPrimitive():
		t.Logf("%s%s (%v) = %v", prefix, node.Key, node.Type, node.Value)
	}
}

// printAnalyzeResult 打印分析结果（用于调试）
func printAnalyzeResult(t *testing.T, result *analyzer.AnalyzeResult) {
	t.Log("\n生成的 Go 类型定义：")
	t.Log("==================")

	// 打印根 struct
	printStructDef(t, result.RootStruct, 0)

	// 打印所有嵌套 struct
	for _, structDef := range result.SubStructs {
		t.Log("")
		printStructDef(t, structDef, 0)
	}
}

// printStructDef 打印 struct 定义
func printStructDef(t *testing.T, def *analyzer.StructDef, indent int) {
	if def == nil {
		return
	}

	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	t.Logf("%stype %s struct {", prefix, def.Name)
	for _, field := range def.Fields {
		t.Logf("%s  %-20s %-20s `json:\"%s\" yaml:\"%s\"`",
			prefix,
			field.Name,
			field.Type,
			field.JSONTag,
			field.YAMLTag,
		)
	}
	t.Logf("%s}", prefix)
}