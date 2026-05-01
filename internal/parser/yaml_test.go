package parser

import (
	"strings"
	"testing"
)

func TestYAMLParser_Parse_SimpleObject(t *testing.T) {
	yaml := `
name: myapp
version: 1.0.0
enabled: true
`
	parser := NewYAMLParser()
	result, err := parser.Parse(strings.NewReader(yaml))

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Root == nil {
		t.Fatal("Root node is nil")
	}

	if !result.Root.IsObject() {
		t.Error("Root should be an object")
	}

	if len(result.Root.Children) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(result.Root.Children))
	}

	if name := result.Root.Children["name"]; name == nil {
		t.Error("name field not found")
	} else if name.Value != "myapp" {
		t.Errorf("name value should be 'myapp', got '%v'", name.Value)
	}

	if version := result.Root.Children["version"]; version == nil {
		t.Error("version field not found")
	} else if version.Value != "1.0.0" {
		t.Errorf("version value should be '1.0.0', got '%v'", version.Value)
	}

	if enabled := result.Root.Children["enabled"]; enabled == nil {
		t.Error("enabled field not found")
	} else if enabled.Value != true {
		t.Errorf("enabled value should be true, got %v", enabled.Value)
	}
}

func TestYAMLParser_Parse_NestedObject(t *testing.T) {
	yaml := `
server:
  host: localhost
  port: 8080
  timeout: 30
`
	parser := NewYAMLParser()
	result, err := parser.Parse(strings.NewReader(yaml))

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	server := result.Root.Children["server"]
	if server == nil {
		t.Fatal("server field not found")
	}

	if !server.IsObject() {
		t.Error("server should be an object")
	}

	if len(server.Children) != 3 {
		t.Errorf("server should have 3 fields, got %d", len(server.Children))
	}

	if host := server.Children["host"]; host == nil {
		t.Error("host field not found")
	} else if host.Value != "localhost" {
		t.Errorf("host should be 'localhost', got '%v'", host.Value)
	}

	if port := server.Children["port"]; port == nil {
		t.Error("port field not found")
	} else if port.Value != 8080 {
		t.Errorf("port should be 8080, got %v", port.Value)
	}
}

func TestYAMLParser_Parse_Array(t *testing.T) {
	yaml := `
features:
  - cache
  - metrics
  - logging
`
	parser := NewYAMLParser()
	result, err := parser.Parse(strings.NewReader(yaml))

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	features := result.Root.Children["features"]
	if features == nil {
		t.Fatal("features field not found")
	}

	if !features.IsArray() {
		t.Error("features should be an array")
	}

	if len(features.Items) != 3 {
		t.Errorf("features should have 3 items, got %d", len(features.Items))
	}

	expectedValues := []string{"cache", "metrics", "logging"}
	for i, expected := range expectedValues {
		if features.Items[i].Value != expected {
			t.Errorf("features[%d] should be '%s', got '%v'", i, expected, features.Items[i].Value)
		}
	}
}

func TestYAMLParser_Parse_MixedTypes(t *testing.T) {
	yaml := `
string_val: "hello"
int_val: 42
float_val: 3.14
bool_val: true
null_val: null
`
	parser := NewYAMLParser()
	result, err := parser.Parse(strings.NewReader(yaml))

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	tests := []struct {
		field        string
		expectedType ValueType
		expectedVal  any
	}{
		{"string_val", TypeString, "hello"},
		{"int_val", TypeInt, 42},
		{"float_val", TypeFloat, 3.14},
		{"bool_val", TypeBool, true},
		{"null_val", TypeNull, nil},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			node := result.Root.Children[tt.field]
			if node == nil {
				t.Fatalf("%s field not found", tt.field)
			}

			if node.Type != tt.expectedType {
				t.Errorf("%s type should be %v, got %v", tt.field, tt.expectedType, node.Type)
			}

			if node.Value != tt.expectedVal {
				t.Errorf("%s value should be %v, got %v", tt.field, tt.expectedVal, node.Value)
			}
		})
	}
}

func TestYAMLParser_Parse_ComplexStructure(t *testing.T) {
	yaml := `
server:
  host: localhost
  port: 8080
  features:
    - ssl
    - compression

database:
  connections:
    - host: db1.example.com
      port: 5432
    - host: db2.example.com
      port: 5433
`
	parser := NewYAMLParser()
	result, err := parser.Parse(strings.NewReader(yaml))

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	server := result.Root.Children["server"]
	if server == nil {
		t.Fatal("server not found")
	}

	features := server.Children["features"]
	if features == nil {
		t.Fatal("server.features not found")
	}

	if !features.IsArray() {
		t.Error("server.features should be an array")
	}

	database := result.Root.Children["database"]
	if database == nil {
		t.Fatal("database not found")
	}

	connections := database.Children["connections"]
	if connections == nil {
		t.Fatal("database.connections not found")
	}

	if !connections.IsArray() {
		t.Error("database.connections should be an array")
	}

	if len(connections.Items) != 2 {
		t.Errorf("database.connections should have 2 items, got %d", len(connections.Items))
	}

	// 验证数组中的对象结构
	conn1 := connections.Items[0]
	if !conn1.IsObject() {
		t.Error("connection item should be an object")
	}

	if conn1.Children["host"].Value != "db1.example.com" {
		t.Errorf("conn1 host should be 'db1.example.com', got '%v'", conn1.Children["host"].Value)
	}
}

func TestYAMLParser_Parse_InvalidYAML(t *testing.T) {
	yaml := `
invalid:
  - item1
 - item2  # wrong indentation
`
	parser := NewYAMLParser()
	_, err := parser.Parse(strings.NewReader(yaml))

	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestYAMLParser_SupportedExtensions(t *testing.T) {
	parser := NewYAMLParser()
	extensions := parser.SupportedExtensions()

	expectedExtensions := []string{".yaml", ".yml"}
	if len(extensions) != len(expectedExtensions) {
		t.Errorf("Expected %d extensions, got %d", len(expectedExtensions), len(extensions))
	}

	for _, ext := range expectedExtensions {
		found := false
		for _, e := range extensions {
			if e == ext {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Extension %s not found in supported extensions", ext)
		}
	}
}

func TestYAMLParser_Name(t *testing.T) {
	parser := NewYAMLParser()
	if parser.Name() != "YAML" {
		t.Errorf("Expected parser name 'YAML', got '%s'", parser.Name())
	}
}
