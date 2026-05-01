package codegen

import (
	"strings"
	"testing"

	"github.com/cheryl-chun/confgen/internal/analyzer"
)

func TestGenerateSimpleStruct(t *testing.T) {
	result := &analyzer.AnalyzeResult{
		RootStruct: &analyzer.StructDef{
			Name: "Config",
			Fields: []*analyzer.FieldDef{
				{
					Name:         "Host",
					Type:         "string",
					JSONTag:      "host",
					YAMLTag:      "host",
					MapStructTag: "host",
				},
				{
					Name:         "Port",
					Type:         "int",
					JSONTag:      "port",
					YAMLTag:      "port",
					MapStructTag: "port",
				},
			},
		},
		SubStructs: make(map[string]*analyzer.StructDef),
	}

	opts := Options{
		PackageName: "main",
		AddComments: false,
	}

	code, err := Generate(result, opts)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	t.Log("Generated code:")
	t.Log(code)

	// 检查关键内容
	mustContain(t, code, "package main")
	mustContain(t, code, "type Config struct")
	mustContain(t, code, "Host")
	mustContain(t, code, "Port")
	mustContain(t, code, `json:"host"`)
	mustContain(t, code, `yaml:"port"`)
	mustContain(t, code, "DO NOT EDIT")
}

func TestGenerateNestedStruct(t *testing.T) {
	result := &analyzer.AnalyzeResult{
		RootStruct: &analyzer.StructDef{
			Name: "Config",
			Fields: []*analyzer.FieldDef{
				{
					Name:         "Server",
					Type:         "ServerConfig",
					JSONTag:      "server",
					YAMLTag:      "server",
					MapStructTag: "server",
				},
			},
		},
		SubStructs: map[string]*analyzer.StructDef{
			"ServerConfig": {
				Name: "ServerConfig",
				Fields: []*analyzer.FieldDef{
					{
						Name:         "Host",
						Type:         "string",
						JSONTag:      "host",
						YAMLTag:      "host",
						MapStructTag: "host",
					},
					{
						Name:         "Port",
						Type:         "int",
						JSONTag:      "port",
						YAMLTag:      "port",
						MapStructTag: "port",
					},
				},
			},
		},
	}

	opts := DefaultOptions()
	code, err := Generate(result, opts)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	t.Log("Generated code:")
	t.Log(code)

	mustContain(t, code, "type Config struct")
	mustContain(t, code, "type ServerConfig struct")
	mustContain(t, code, "Server")
	mustContain(t, code, "ServerConfig")
}

func TestGenerateArrayTypes(t *testing.T) {
	result := &analyzer.AnalyzeResult{
		RootStruct: &analyzer.StructDef{
			Name: "Config",
			Fields: []*analyzer.FieldDef{
				{
					Name:         "Tags",
					Type:         "[]string",
					JSONTag:      "tags",
					YAMLTag:      "tags",
					MapStructTag: "tags",
				},
				{
					Name:         "Servers",
					Type:         "[]ServerConfig",
					JSONTag:      "servers",
					YAMLTag:      "servers",
					MapStructTag: "servers",
				},
			},
		},
		SubStructs: map[string]*analyzer.StructDef{
			"ServerConfig": {
				Name: "ServerConfig",
				Fields: []*analyzer.FieldDef{
					{
						Name:         "Host",
						Type:         "string",
						JSONTag:      "host",
						YAMLTag:      "host",
						MapStructTag: "host",
					},
				},
			},
		},
	}

	opts := DefaultOptions()
	code, err := Generate(result, opts)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	t.Log("Generated code:")
	t.Log(code)

	mustContain(t, code, "[]string")
	mustContain(t, code, "[]ServerConfig")
}

func TestGenerateWithComments(t *testing.T) {
	result := &analyzer.AnalyzeResult{
		RootStruct: &analyzer.StructDef{
			Name: "Config",
			Fields: []*analyzer.FieldDef{
				{
					Name:         "Host",
					Type:         "string",
					JSONTag:      "host",
					YAMLTag:      "host",
					MapStructTag: "host",
					Comment:      "Server hostname",
				},
			},
		},
		SubStructs: make(map[string]*analyzer.StructDef),
	}

	opts := Options{
		PackageName: "config",
		AddComments: true,
	}

	code, err := Generate(result, opts)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	t.Log("Generated code:")
	t.Log(code)

	mustContain(t, code, "package config")
	mustContain(t, code, "// Server hostname")
}

func mustContain(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("Generated code does not contain %q", substr)
	}
}