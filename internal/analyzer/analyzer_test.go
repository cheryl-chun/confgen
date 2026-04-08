package analyzer

import (
	"testing"

	"github.com/cheryl-chun/confgen/internal/parser"
)

func TestAnalyzeSimpleObject(t *testing.T) {
	// 构建测试用的 ConfigNode
	root := parser.NewConfigNode("root")
	root.Type = parser.TypeObject

	// 添加简单字段
	host := parser.NewConfigNode("host")
	host.Type = parser.TypeString
	host.Value = "localhost"
	root.AddChild(host)

	port := parser.NewConfigNode("port")
	port.Type = parser.TypeInt
	port.Value = 8080
	root.AddChild(port)

	enabled := parser.NewConfigNode("enabled")
	enabled.Type = parser.TypeBool
	enabled.Value = true
	root.AddChild(enabled)

	// 分析
	result, err := Analyze(root)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// 验证根 struct
	if result.RootStruct == nil {
		t.Fatal("RootStruct is nil")
	}

	if result.RootStruct.Name != "Config" {
		t.Errorf("RootStruct.Name = %q, want %q", result.RootStruct.Name, "Config")
	}

	if len(result.RootStruct.Fields) != 3 {
		t.Errorf("len(Fields) = %d, want 3", len(result.RootStruct.Fields))
	}

	// 验证字段
	fields := result.RootStruct.Fields
	expectedFields := map[string]string{
		"Host":    "string",
		"Port":    "int",
		"Enabled": "bool",
	}

	for _, field := range fields {
		expectedType, ok := expectedFields[field.Name]
		if !ok {
			t.Errorf("Unexpected field: %s", field.Name)
			continue
		}
		if field.Type != expectedType {
			t.Errorf("Field %s: Type = %q, want %q", field.Name, field.Type, expectedType)
		}
	}
}

func TestAnalyzeNestedObject(t *testing.T) {
	// 构建嵌套对象
	root := parser.NewConfigNode("root")
	root.Type = parser.TypeObject

	// server 节点
	server := parser.NewConfigNode("server")
	server.Type = parser.TypeObject
	root.AddChild(server)

	host := parser.NewConfigNode("host")
	host.Type = parser.TypeString
	host.Value = "localhost"
	server.AddChild(host)

	port := parser.NewConfigNode("port")
	port.Type = parser.TypeInt
	port.Value = 8080
	server.AddChild(port)

	// 分析
	result, err := Analyze(root)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// 验证根 struct
	if len(result.RootStruct.Fields) != 1 {
		t.Errorf("len(RootStruct.Fields) = %d, want 1", len(result.RootStruct.Fields))
	}

	serverField := result.RootStruct.Fields[0]
	if serverField.Name != "Server" {
		t.Errorf("serverField.Name = %q, want %q", serverField.Name, "Server")
	}

	if serverField.Type != "ServerConfig" {
		t.Errorf("serverField.Type = %q, want %q", serverField.Type, "ServerConfig")
	}

	// 验证嵌套 struct
	serverStruct, ok := result.SubStructs["ServerConfig"]
	if !ok {
		t.Fatal("ServerConfig not found in SubStructs")
	}

	if len(serverStruct.Fields) != 2 {
		t.Errorf("len(ServerConfig.Fields) = %d, want 2", len(serverStruct.Fields))
	}
}

func TestAnalyzeArray(t *testing.T) {
	// 构建数组节点
	root := parser.NewConfigNode("root")
	root.Type = parser.TypeObject

	tags := parser.NewConfigNode("tags")
	tags.Type = parser.TypeArray
	root.AddChild(tags)

	// 添加数组元素
	tag1 := parser.NewConfigNode("[0]")
	tag1.Type = parser.TypeString
	tag1.Value = "production"
	tags.AddItem(tag1)

	tag2 := parser.NewConfigNode("[1]")
	tag2.Type = parser.TypeString
	tag2.Value = "web"
	tags.AddItem(tag2)

	// 分析
	result, err := Analyze(root)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// 验证数组字段
	if len(result.RootStruct.Fields) != 1 {
		t.Fatalf("len(Fields) = %d, want 1", len(result.RootStruct.Fields))
	}

	tagsField := result.RootStruct.Fields[0]
	if tagsField.Name != "Tags" {
		t.Errorf("tagsField.Name = %q, want %q", tagsField.Name, "Tags")
	}

	if tagsField.Type != "[]string" {
		t.Errorf("tagsField.Type = %q, want %q", tagsField.Type, "[]string")
	}
}

func TestAnalyzeArrayOfObjects(t *testing.T) {
	// 构建对象数组
	root := parser.NewConfigNode("root")
	root.Type = parser.TypeObject

	servers := parser.NewConfigNode("servers")
	servers.Type = parser.TypeArray
	root.AddChild(servers)

	// 第一个 server
	server1 := parser.NewConfigNode("[0]")
	server1.Type = parser.TypeObject
	servers.AddItem(server1)

	host1 := parser.NewConfigNode("host")
	host1.Type = parser.TypeString
	host1.Value = "server1"
	server1.AddChild(host1)

	// 分析
	result, err := Analyze(root)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// 验证数组字段
	serversField := result.RootStruct.Fields[0]
	if serversField.Name != "Servers" {
		t.Errorf("serversField.Name = %q, want %q", serversField.Name, "Servers")
	}

	// 应该是 []ServerConfig
	if serversField.Type != "[]ServerConfig" {
		t.Errorf("serversField.Type = %q, want %q", serversField.Type, "[]ServerConfig")
	}

	// 验证 ServerConfig struct
	if _, ok := result.SubStructs["ServerConfig"]; !ok {
		t.Error("ServerConfig not found in SubStructs")
	}
}
