package analyzer

import (
	"strings"
	"unicode"
)

// ToStructName transforms a configuration key into a valid Go struct identifier 
// by applying PascalCase and appending a "Config" suffix.
// Example: "database" -> "DatabaseConfig", "max_connections" -> "MaxConnectionsConfig".
func ToStructName(key string) string {
	if key == "" || key == "root" {
		return "Config"
	}
	return ToPascalCase(key) + "Config"
}

// ToFieldName converts a configuration key into an exported Go field name 
// using PascalCase to ensure visibility outside the package.
// Example: "max_connections" -> "MaxConnections", "api_key" -> "ApiKey".
func ToFieldName(key string) string {
	return ToPascalCase(key)
}

// ToPascalCase normalizes various naming conventions (snake_case, kebab-case, camelCase) 
// into a standardized PascalCase format.
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	// Tokenize the string into discrete words based on delimiters or casing boundaries.
	words := splitWords(s)

	var result strings.Builder
	for _, word := range words {
		if word == "" {
			continue
		}
		// Apply casing transformations while respecting Go initialism conventions.
		result.WriteString(capitalize(word))
	}

	return result.String()
}

// splitWords parses the input string into a slice of tokens by identifying 
// delimiters ('_', '-', ' ') and detecting transitions in camelCase.
func splitWords(s string) []string {
	var words []string
	var currentWord strings.Builder

	for i, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			// Trigger word boundary on explicit delimiters.
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		} else if i > 0 && unicode.IsUpper(r) && unicode.IsLower(rune(s[i-1])) {
			// Trigger word boundary on camelCase transitions (e.g., camelCase -> [camel, Case]).
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			currentWord.WriteRune(r)
		} else {
			currentWord.WriteRune(r)
		}
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

// capitalize applies title casing to a single token. It prioritizes 
// industry-standard initialisms (e.g., "url" -> "URL") to align with 
// Go's static analysis recommendations (golint).
func capitalize(s string) string {
	if s == "" {
		return ""
	}

	// Check for common technical acronyms that should remain fully uppercase.
	upper := strings.ToUpper(s)
	if isCommonAcronym(upper) {
		return upper
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// isCommonAcronym validates if a token is a recognized initialism. 
// This follows the Go community's best practices for identifier naming.
func isCommonAcronym(s string) bool {
	acronyms := map[string]bool{
		"ID":    true,
		"API":   true,
		"URL":   true,
		"URI":   true,
		"HTTP":  true,
		"HTTPS": true,
		"JSON":  true,
		"XML":   true,
		"HTML":  true,
		"SQL":   true,
		"DB":    true,
		"TCP":   true,
		"UDP":   true,
		"IP":    true,
		"TLS":   true,
		"SSL":   true,
		"CPU":   true,
		"RAM":   true,
		"UUID":  true,
	}
	return acronyms[s]
}