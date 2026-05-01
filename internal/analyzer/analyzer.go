package analyzer

import (
	"fmt"

	"github.com/cheryl-chun/confgen/internal/parser"
)

// Analyzer performs static type inference on a ConfigNode tree to derive 
// idiomatic Go structural representations. It maintains internal state 
// to ensure unique identifier generation and handle potential naming collisions.
type Analyzer struct {
	result        *AnalyzeResult
	structCounter map[string]int // Registry for tracking occurrences to resolve naming collisions.
}

// NewAnalyzer initializes a new inference engine with an empty result set 
// and an active naming registry for structural tracking.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		result:        NewAnalyzeResult(),
		structCounter: make(map[string]int),
	}
}

// Analyze traverses the provided ConfigNode tree and generates a collection 
// of Go type definitions. It requires the root node to be of type 'Object'.
// It returns a comprehensive AnalyzeResult or an error if the tree structure is invalid.
func (a *Analyzer) Analyze(root *parser.ConfigNode) (*AnalyzeResult, error) {
	if root == nil {
		return nil, fmt.Errorf("root node is nil")
	}

	// The entry point of a configuration must represent a top-level associative object.
	if !root.IsObject() {
		return nil, fmt.Errorf("root node must be an object, got %v", root.Type)
	}

	// Recursively resolve the root object into the primary 'Config' structure.
	rootStruct := a.analyzeObject("Config", root)
	a.result.RootStruct = rootStruct

	return a.result, nil
}

// analyzeObject maps a ConfigNode group to a StructDef. It iterates through 
// all child nodes and triggers field-level type inference for each member.
func (a *Analyzer) analyzeObject(structName string, node *parser.ConfigNode) *StructDef {
	def := &StructDef{
		Name:   structName,
		Fields: make([]*FieldDef, 0, len(node.Children)),
	}

	// Process each child as a potential struct field.
	for key, child := range node.Children {
		field := a.analyzeField(key, child)
		def.Fields = append(def.Fields, field)
	}

	return def
}

// analyzeField evaluates an individual node to determine its corresponding Go field 
// metadata, including PascalCase identifier conversion and serialization tags.
func (a *Analyzer) analyzeField(key string, node *parser.ConfigNode) *FieldDef {
	field := &FieldDef{
		Name:         ToFieldName(key), // Perform casing normalization for exported visibility.
		JSONTag:      key,
		YAMLTag:      key,
		MapStructTag: key,
	}

	// Dispatch type inference logic based on node characteristics.
	switch {
	case node.IsPrimitive():
		// Map scalar types: string, int, float64, bool.
		field.Type = GoType(node.Type)

	case node.IsArray():
		// Delegate sequential collection analysis for slice type derivation.
		field.Type = a.analyzeArrayType(key, node)

	case node.IsObject():
		// Nested structure detected: instantiate a new auxiliary struct definition.
		nestedStructName := ToStructName(key)
		nestedStruct := a.analyzeObject(nestedStructName, node)

		// Register the newly discovered structural dependency.
		a.result.AddStruct(nestedStruct)

		// Reference the nested struct as the field's data type.
		field.Type = nestedStructName

	default:
		// Fallback to empty interface for indeterminate or unsupported types.
		field.Type = "interface{}"
	}

	return field
}

// analyzeArrayType employs a sampling heuristic to determine the element type 
// for sequential collections. It assumes type homogeneity across the collection.
func (a *Analyzer) analyzeArrayType(key string, node *parser.ConfigNode) string {
	if len(node.Items) == 0 {
		// Insufficient metadata for inference; defaulting to a generic interface slice.
		return "[]interface{}"
	}

	// Sample the first element to derive the representative type for the entire slice.
	firstItem := node.Items[0]

	var elemType string
	switch {
	case firstItem.IsPrimitive():
		elemType = GoType(firstItem.Type)

	case firstItem.IsObject():
		// For object arrays, generate a structural definition using singularized keys.
		// e.g., "servers[0]" -> "ServerConfig".
		itemStructName := ToStructName(key)
		
		// Apply a basic singularization heuristic by stripping the "s" suffix before the "Config" identifier.
		if len(itemStructName) > 7 && itemStructName[len(itemStructName)-7:] == "sConfig" {
			itemStructName = itemStructName[:len(itemStructName)-7] + "Config"
		}

		itemStruct := a.analyzeObject(itemStructName, firstItem)
		a.result.AddStruct(itemStruct)
		elemType = itemStructName

	case firstItem.IsArray():
		// Handle multi-dimensional collections via recursive analysis.
		elemType = a.analyzeArrayType(key+"_item", firstItem)

	default:
		// Fallback for complex or indeterminate element types.
		elemType = "interface{}"
	}

	return "[]" + elemType
}

// Analyze is a high-level package helper that abstracts the inference engine's 
// instantiation and execution for standard configuration trees.
func Analyze(root *parser.ConfigNode) (*AnalyzeResult, error) {
	analyzer := NewAnalyzer()
	return analyzer.Analyze(root)
}