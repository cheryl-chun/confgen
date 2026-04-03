package parser

import "io"

/**
* The goal of the Parser is to uniformly convert
* configuration files of different formats
* (such as YAML/JSON/TOML, etc.) into a standard tree structure,
* which can be used for subsequent type inference and
* code generation.
 */

// Parser configuration file parser interface
// All parsers for various formats need to implement this interface
type Parser interface {
	// Parse the configuration file from io.Reader
	Parse(reader io.Reader) (*ParseResult, error)

	// ParseFile parses the configuration from a file path
	ParseFile(path string) (*ParseResult, error)

	// SupportedExtensions return the supported file extensions
	SupportedExtensions() []string

	// Name returns the parser name
	Name() string
}

// Provide general functions to avoid 
// each parser having to repeatedly implement the Name() 
// and SupportedExtensions() methods.
type BaseParser struct {
	name       string
	extensions []string
}

func NewBaseParser(name string, extensions []string) *BaseParser {
	return &BaseParser{
		name:       name,
		extensions: extensions,
	}
}

func (b *BaseParser) SupportedExtensions() []string {
	return b.extensions
}

func (b *BaseParser) Name() string {
	return b.name
}
