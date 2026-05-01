package analyzer

import "testing"

// TestToPascalCase verifies the transformation logic from various string formats 
// (snake_case, kebab-case, camelCase) to PascalCase.
// It specifically validates compliance with Go's official initialism guidelines 
// (e.g., converting "api" to "API" instead of "Api").
func TestToPascalCase(t *testing.T) {
	// Define a suite of test cases covering standard scenarios and edge cases.
	tests := []struct {
		input    string
		expected string
	}{
		{"max_connections", "MaxConnections"},
		{"api_key", "APIKey"},         // Compliance check: API should be fully capitalized.
		{"server-port", "ServerPort"},
		{"camelCase", "CamelCase"},
		{"snake_case_name", "SnakeCaseName"},
		{"simple", "Simple"},
		{"", ""},                      // Edge case: empty input.
		{"db_host", "DBHost"},         // Compliance check: DB should be fully capitalized.
		{"api_url", "APIURL"},         // Multi-acronym check: Both API and URL should be capitalized.
		{"id", "ID"},                  // Singular acronym check.
		{"user_id", "UserID"},         // Common suffix acronym check.
	}

	for _, tt := range tests {
		// Utilize subtests for granular failure reporting and isolated execution.
		t.Run(tt.input, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToStructName ensures that configuration keys are correctly mapped to 
// struct identifiers with the appropriate "Config" suffix.
func TestToStructName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"server", "ServerConfig"},
		{"database", "DatabaseConfig"},
		{"max_connections", "MaxConnectionsConfig"},
		{"root", "Config"},            // Special case: 'root' should map to the primary 'Config' struct.
		{"", "Config"},                // Fallback: empty keys default to 'Config'.
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

// TestToFieldName validates the generation of exported Go field names.
// It ensures that fields are properly capitalized to maintain package-level visibility.
func TestToFieldName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"max_connections", "MaxConnections"},
		{"host", "Host"},
		{"port", "Port"},
		{"api_key", "APIKey"},         // Ensuring acronyms are respected in field naming.
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