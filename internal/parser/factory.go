package parser

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ErrUnsupportedFormat = fmt.Errorf("unsupported config file format")
	ErrParserNotFound = fmt.Errorf("parser not found for file type")
)

type ParserFactory struct {
	mu      sync.RWMutex
	// key: file extension name
	// value: Parser instance
	parsers map[string]Parser 
}

var (
	globalFactory     *ParserFactory
	globalFactoryOnce sync.Once
)

// ParserFactory singleton
func GetFactory() *ParserFactory {
	globalFactoryOnce.Do(func() {
		globalFactory = NewParserFactory()
		// Register default parsers (YAML, JSON, etc.)
		globalFactory.RegisterDefaultParsers()
	})
	return globalFactory
}

func NewParserFactory() *ParserFactory {
	return &ParserFactory{
		parsers: make(map[string]Parser),
	}
}

func (f *ParserFactory) Register(parser Parser) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, ext := range parser.SupportedExtensions() {
		ext = strings.ToLower(strings.TrimPrefix(ext, "."))
		f.parsers[ext] = parser
	}
}

func (f *ParserFactory) RegisterDefaultParsers() {
	f.Register(NewYAMLParser())
	f.Register(NewJSONParser())

	// TODO: You can add your own parsers here
	// TODO: such as TOML, INI, etc.
}

// GetParser according to file extension
func (f *ParserFactory) GetParser(ext string) (Parser, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	parser, ok := f.parsers[ext]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrParserNotFound, ext)
	}
	return parser, nil
}

// GetParserByFilePath according to file path
// extract the extension and get the corresponding parser
func (f *ParserFactory) GetParserByFilePath(path string) (Parser, error) {
	ext := filepath.Ext(path)
	if ext == "" {
		return nil, fmt.Errorf("%w: no file extension in %s", ErrUnsupportedFormat, path)
	}
	return f.GetParser(ext)
}

// SupportedFormats returns all supported formats
func (f *ParserFactory) SupportedFormats() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	formats := make([]string, 0, len(f.parsers))
	seen := make(map[string]bool)

	for ext := range f.parsers {
		if !seen[ext] {
			formats = append(formats, ext)
			seen[ext] = true
		}
	}
	return formats
}

// ParseFile convenient method: automatically select the parser and parse the file
func (f *ParserFactory) ParseFile(path string) (*ParseResult, error) {
	parser, err := f.GetParserByFilePath(path)
	if err != nil {
		return nil, err
	}
	return parser.ParseFile(path)
}

func ParseFile(path string) (*ParseResult, error) {
	return GetFactory().ParseFile(path)
}

func Register(parser Parser) {
	GetFactory().Register(parser)
}

func SupportedFormats() []string {
	return GetFactory().SupportedFormats()
}