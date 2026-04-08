package analyzer

import (
	"strings"
	"unicode"
)

// ToStructName 将配置键转换为 Go struct 名称 (PascalCase + "Config" 后缀)
// 例如: "database" -> "DatabaseConfig", "max_connections" -> "MaxConnectionsConfig"
func ToStructName(key string) string {
	if key == "" || key == "root" {
		return "Config"
	}
	return ToPascalCase(key) + "Config"
}

// ToFieldName 将配置键转换为 Go 字段名 (PascalCase)
// 例如: "max_connections" -> "MaxConnections", "api_key" -> "ApiKey"
func ToFieldName(key string) string {
	return ToPascalCase(key)
}

// ToPascalCase 转换为 PascalCase
// 支持: snake_case, kebab-case, camelCase
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	// 分割字符串（按 '_', '-', 或驼峰边界）
	words := splitWords(s)

	var result strings.Builder
	for _, word := range words {
		if word == "" {
			continue
		}
		// 首字母大写，其余小写
		result.WriteString(capitalize(word))
	}

	return result.String()
}

// splitWords 将字符串分割为单词
func splitWords(s string) []string {
	var words []string
	var currentWord strings.Builder

	for i, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			// 分隔符
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		} else if i > 0 && unicode.IsUpper(r) && unicode.IsLower(rune(s[i-1])) {
			// 驼峰边界: camelCase -> [camel, Case]
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

// capitalize 首字母大写，其余小写
func capitalize(s string) string {
	if s == "" {
		return ""
	}

	// 特殊缩写词保持全大写 (可根据需要扩展)
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

// isCommonAcronym 判断是否为常见缩写词
func isCommonAcronym(s string) bool {
	acronyms := map[string]bool{
		"ID":   true,
		"API":  true,
		"URL":  true,
		"URI":  true,
		"HTTP": true,
		"HTTPS": true,
		"JSON": true,
		"XML":  true,
		"HTML": true,
		"SQL":  true,
		"DB":   true,
		"TCP":  true,
		"UDP":  true,
		"IP":   true,
		"TLS":  true,
		"SSL":  true,
		"CPU":  true,
		"RAM":  true,
		"UUID": true,
	}
	return acronyms[s]
}
