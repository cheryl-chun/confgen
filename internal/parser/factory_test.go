package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewParserFactory(t *testing.T) {
	factory := NewParserFactory()

	if factory == nil {
		t.Fatal("Factory should not be nil")
	}

	if factory.parsers == nil {
		t.Error("Parsers map should be initialized")
	}

	if len(factory.parsers) != 0 {
		t.Error("New factory should have no parsers registered")
	}
}

func TestGetFactory_Singleton(t *testing.T) {
	factory1 := GetFactory()
	factory2 := GetFactory()

	if factory1 != factory2 {
		t.Error("GetFactory should return the same instance (singleton)")
	}
}

func TestParserFactory_Register(t *testing.T) {
	factory := NewParserFactory()

	yamlParser := NewYAMLParser()
	factory.Register(yamlParser)

	// 验证 .yaml 和 .yml 都被注册
	if len(factory.parsers) != 2 {
		t.Errorf("Expected 2 extensions registered, got %d", len(factory.parsers))
	}

	if factory.parsers["yaml"] == nil {
		t.Error("yaml extension not registered")
	}

	if factory.parsers["yml"] == nil {
		t.Error("yml extension not registered")
	}
}

func TestParserFactory_RegisterDefaultParsers(t *testing.T) {
	factory := NewParserFactory()
	factory.RegisterDefaultParsers()

	// 应该有 YAML (.yaml, .yml) 和 JSON (.json) = 3 个扩展名
	if len(factory.parsers) < 3 {
		t.Errorf("Expected at least 3 extensions, got %d", len(factory.parsers))
	}

	// 验证必须的解析器
	requiredExtensions := []string{"yaml", "yml", "json"}
	for _, ext := range requiredExtensions {
		if factory.parsers[ext] == nil {
			t.Errorf("Extension %s not registered", ext)
		}
	}
}

func TestParserFactory_GetParser(t *testing.T) {
	factory := NewParserFactory()
	factory.RegisterDefaultParsers()

	tests := []struct {
		ext         string
		shouldExist bool
	}{
		{".yaml", true},
		{"yaml", true},
		{".yml", true},
		{"yml", true},
		{".json", true},
		{"json", true},
		{".toml", false},
		{".xml", false},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			parser, err := factory.GetParser(tt.ext)

			if tt.shouldExist {
				if err != nil {
					t.Errorf("Expected parser for %s, got error: %v", tt.ext, err)
				}
				if parser == nil {
					t.Errorf("Parser for %s should not be nil", tt.ext)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for unsupported extension %s", tt.ext)
				}
				if parser != nil {
					t.Errorf("Parser for %s should be nil", tt.ext)
				}
			}
		})
	}
}

func TestParserFactory_GetParserByFilePath(t *testing.T) {
	factory := NewParserFactory()
	factory.RegisterDefaultParsers()

	tests := []struct {
		path        string
		shouldExist bool
		parserName  string
	}{
		{"config.yaml", true, "YAML"},
		{"config.yml", true, "YAML"},
		{"config.json", true, "JSON"},
		{"/path/to/config.yaml", true, "YAML"},
		{"config", false, ""}, // no extension
		{"config.toml", false, ""},
		{"config.xml", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			parser, err := factory.GetParserByFilePath(tt.path)

			if tt.shouldExist {
				if err != nil {
					t.Errorf("Expected parser for %s, got error: %v", tt.path, err)
				}
				if parser == nil {
					t.Errorf("Parser for %s should not be nil", tt.path)
				} else if parser.Name() != tt.parserName {
					t.Errorf("Expected parser %s, got %s", tt.parserName, parser.Name())
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for unsupported path %s", tt.path)
				}
			}
		})
	}
}

func TestParserFactory_SupportedFormats(t *testing.T) {
	factory := NewParserFactory()
	factory.RegisterDefaultParsers()

	formats := factory.SupportedFormats()

	if len(formats) < 3 {
		t.Errorf("Expected at least 3 formats, got %d", len(formats))
	}

	// 检查是否包含必需的格式
	requiredFormats := map[string]bool{
		"yaml": false,
		"yml":  false,
		"json": false,
	}

	for _, format := range formats {
		if _, exists := requiredFormats[format]; exists {
			requiredFormats[format] = true
		}
	}

	for format, found := range requiredFormats {
		if !found {
			t.Errorf("Format %s not found in supported formats", format)
		}
	}
}

func TestParserFactory_ParseFile_Integration(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	// 创建测试文件
	yamlFile := filepath.Join(tmpDir, "test.yaml")
	yamlContent := `
server:
  host: localhost
  port: 8080
`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	jsonFile := filepath.Join(tmpDir, "test.json")
	jsonContent := `{
		"server": {
			"host": "localhost",
			"port": 8080
		}
	}`
	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	factory := NewParserFactory()
	factory.RegisterDefaultParsers()

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{"YAML file", yamlFile, false},
		{"JSON file", jsonFile, false},
		{"Non-existent file", filepath.Join(tmpDir, "notexist.yaml"), true},
		{"Unsupported format", filepath.Join(tmpDir, "test.toml"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := factory.ParseFile(tt.filePath)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Result should not be nil")
				}
				if result.Root == nil {
					t.Error("Root node should not be nil")
				}

				// 验证解析结果
				server := result.Root.Children["server"]
				if server == nil {
					t.Fatal("server field not found")
				}
				if !server.IsObject() {
					t.Error("server should be an object")
				}
			}
		})
	}
}

func TestParseFile_PackageLevel(t *testing.T) {
	// 测试包级别的便捷函数
	tmpDir := t.TempDir()
	yamlFile := filepath.Join(tmpDir, "config.yaml")
	yamlContent := `name: test`
	if err := os.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := ParseFile(yamlFile)

	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if result == nil || result.Root == nil {
		t.Fatal("Result should not be nil")
	}

	if result.Root.Children["name"] == nil {
		t.Error("name field not found")
	}
}

func TestRegister_PackageLevel(t *testing.T) {
	// 测试包级别的注册函数
	// 注意：这会影响全局工厂，所以要小心
	customParser := NewYAMLParser()
	Register(customParser)

	// 验证可以获取到解析器
	factory := GetFactory()
	parser, err := factory.GetParser("yaml")

	if err != nil {
		t.Errorf("Failed to get registered parser: %v", err)
	}

	if parser == nil {
		t.Error("Parser should not be nil")
	}
}

func TestSupportedFormats_PackageLevel(t *testing.T) {
	formats := SupportedFormats()

	if len(formats) == 0 {
		t.Error("Should have at least one supported format")
	}

	// 验证包含基本格式
	hasYAML := false
	hasJSON := false

	for _, format := range formats {
		if format == "yaml" || format == "yml" {
			hasYAML = true
		}
		if format == "json" {
			hasJSON = true
		}
	}

	if !hasYAML {
		t.Error("YAML format should be supported")
	}

	if !hasJSON {
		t.Error("JSON format should be supported")
	}
}

func TestParserFactory_ExtensionNormalization(t *testing.T) {
	factory := NewParserFactory()
	yamlParser := NewYAMLParser()
	factory.Register(yamlParser)

	// 测试各种格式的扩展名都能正确匹配
	tests := []string{
		".yaml",
		"yaml",
		".YAML",
		"YAML",
		".Yaml",
	}

	for _, ext := range tests {
		t.Run(ext, func(t *testing.T) {
			parser, err := factory.GetParser(ext)
			if err != nil {
				t.Errorf("Failed to get parser for %s: %v", ext, err)
			}
			if parser == nil {
				t.Errorf("Parser should not be nil for %s", ext)
			}
		})
	}
}

func TestParserFactory_ErrorMessages(t *testing.T) {
	factory := NewParserFactory()
	factory.RegisterDefaultParsers()

	tests := []struct {
		name        string
		path        string
		expectedErr error
	}{
		{
			"unsupported extension",
			"config.toml",
			ErrParserNotFound,
		},
		{
			"no extension",
			"config",
			ErrUnsupportedFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := factory.GetParserByFilePath(tt.path)
			if err == nil {
				t.Error("Expected error, got nil")
			}

			// 验证错误类型
			if !strings.Contains(err.Error(), tt.expectedErr.Error()) {
				t.Errorf("Expected error containing '%v', got '%v'", tt.expectedErr, err)
			}
		})
	}
}
