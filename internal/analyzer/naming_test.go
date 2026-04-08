package analyzer

import "testing"

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"max_connections", "MaxConnections"},
		{"api_key", "APIKey"},          // Go 规范: 缩写词全大写
		{"server-port", "ServerPort"},
		{"camelCase", "CamelCase"},
		{"snake_case_name", "SnakeCaseName"},
		{"simple", "Simple"},
		{"", ""},
		{"db_host", "DBHost"},          // Go 规范: DB 全大写
		{"api_url", "APIURL"},          // Go 规范: API 和 URL 全大写
		{"id", "ID"},                   // Go 规范: ID 全大写
		{"user_id", "UserID"},          // Go 规范: ID 全大写
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToStructName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"server", "ServerConfig"},
		{"database", "DatabaseConfig"},
		{"max_connections", "MaxConnectionsConfig"},
		{"root", "Config"},
		{"", "Config"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToStructName(tt.input)
			if result != tt.expected {
				t.Errorf("ToStructName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToFieldName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"max_connections", "MaxConnections"},
		{"host", "Host"},
		{"port", "Port"},
		{"api_key", "APIKey"}, // Go 规范: API 全大写
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ToFieldName(tt.input)
			if result != tt.expected {
				t.Errorf("ToFieldName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
