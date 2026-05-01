package analyzer

import "github.com/cheryl-chun/confgen/internal/parser"

// StructDef encapsulates the metadata required to define a Go structure.
// It serves as a blueprint for the code generation engine to render 
// valid Go source code.
type StructDef struct {
	Name   string      // The unique identifier for the struct (e.g., "ServerConfig").
	Fields []*FieldDef // A collection of field definitions contained within the struct.
}

// FieldDef represents an individual struct field, encompassing its 
// identifier, Go data type, and multi-format serialization tags.
type FieldDef struct {
	Name         string    // The exported field name in PascalCase (e.g., "MaxConnections").
	Type         string    // The resolved Go type (e.g., "string", "*DatabaseConfig").
	JSONTag      string    // The value for the 'json' struct tag used for marshalling.
	YAMLTag      string    // The value for the 'yaml' struct tag.
	MapStructTag string    // The value for the 'mapstructure' tag (standard for Viper integration).
	Comment      string    // An optional inline comment for the field (for generated documentation).
}

// AnalyzeResult acts as a central registry for all inferred Go types.
// It separates the entry-point structure from its nested dependencies.
type AnalyzeResult struct {
	RootStruct *StructDef            // The primary configuration entry point (typically "Config").
	SubStructs map[string]*StructDef // A registry of auxiliary struct definitions, keyed by their type names.
}

// NewAnalyzeResult initializes an empty AnalyzeResult with an active 
// map for structural dependency tracking.
func NewAnalyzeResult() *AnalyzeResult {
	return &AnalyzeResult{
		SubStructs: make(map[string]*StructDef),
	}
}

// AddStruct registers a new structural definition into the internal registry.
// This ensures that all nested components are available during the code rendering phase.
func (r *AnalyzeResult) AddStruct(def *StructDef) {
	r.SubStructs[def.Name] = def
}

// GoType performs a lookup to resolve internal ValueType abstractions 
// into concrete Go primitive type strings.
func GoType(vt parser.ValueType) string {
	switch vt {
	case parser.TypeString:
		return "string"
	case parser.TypeInt:
		return "int"
	case parser.TypeFloat:
		return "float64"
	case parser.TypeBool:
		return "bool"
	case parser.TypeNull:
		// Fallback to interface{} for null values to provide maximum flexibility.
		return "interface{}"
	default:
		// Unknown or indeterminate types default to an empty interface.
		return "interface{}"
	}
}