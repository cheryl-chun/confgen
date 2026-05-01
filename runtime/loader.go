package runtime

import (
	"fmt"
	"reflect"

	"github.com/cheryl-chun/confgen/internal/tree"
)

// Loader orchestrates the configuration hydration process by aggregating 
// data from multiple heterogeneous sources into a unified configuration tree.
type Loader struct {
	tree    *tree.ConfigTree
	sources []Source
}

// Source defines the interface for configuration providers. 
// It decouples the loading logic from specific formats or protocols.
type Source interface {
	// Load populates the provided ConfigTree with data from the underlying source.
	Load(tree *tree.ConfigTree) error
	// Priority returns the precedence level of the configuration source.
	Priority() tree.SourceType
}

// NewLoader initializes a Loader with an empty configuration tree 
// and an empty registry of sources.
func NewLoader() *Loader {
	return &Loader{
		tree:    tree.NewConfigTree(),
		sources: make([]Source, 0),
	}
}

// AddSource registers a new configuration provider into the loader's pipeline.
func (l *Loader) AddSource(source Source) {
	l.sources = append(l.sources, source)
}

// AddFile is a convenience wrapper to register a local filesystem source.
func (l *Loader) AddFile(path string) *Loader {
	l.AddSource(&FileSource{Path: path})
	return l
}

// AddEnv is a convenience wrapper to register an environment variable source 
// with a specific lookup prefix.
func (l *Loader) AddEnv(prefix string) *Loader {
	l.AddSource(&EnvSource{Prefix: prefix})
	return l
}

// AddRemoteSource allows the injection of custom remote configuration providers.
func (l *Loader) AddRemoteSource(source Source) *Loader {
	l.AddSource(source)
	return l
}

// Fill executes the end-to-end hydration pipeline. 
// It aggregates data from all registered sources and performs reflection-based 
// mapping to populate the provided target struct.
func (l *Loader) Fill(cfg interface{}) error {
	// 1. Validate that the target is a mutable pointer to a struct.
	_, err := l.validateConfigTarget(cfg)
	if err != nil {
		return err
	}

	// 2. Iterate through all sources and merge their data into the internal tree.
	for _, source := range l.sources {
		if err := source.Load(l.tree); err != nil {
			return fmt.Errorf("failed to load source: %w", err)
		}
	}

	// 3. Perform structural hydration and optional tree injection.
	return l.applyTreeToConfig(cfg)
}

// validateConfigTarget ensures the provided interface is a non-nil pointer, 
// which is mandatory for reflection-based mutation.
func (l *Loader) validateConfigTarget(cfg interface{}) (reflect.Value, error) {
	if cfg == nil {
		return reflect.Value{}, fmt.Errorf("config cannot be nil")
	}

	rv := reflect.ValueOf(cfg)
	if rv.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("config must be a pointer for mutation")
	}

	if rv.IsNil() {
		return reflect.Value{}, fmt.Errorf("config pointer is nil")
	}

	return rv, nil
}

// applyTreeToConfig triggers the recursive population of the struct and 
// injects the ConfigTree reference for dynamic runtime queries.
func (l *Loader) applyTreeToConfig(cfg interface{}) error {
	rv, err := l.validateConfigTarget(cfg)
	if err != nil {
		return err
	}

	// Zero out the struct to ensure a clean state before hydration.
	elem := rv.Elem()
	elem.Set(reflect.Zero(elem.Type()))
	
	// Perform recursive field population.
	if err := l.fillStruct(elem, ""); err != nil {
		return fmt.Errorf("failed to fill struct: %w", err)
	}

	// Inject the internal tree metadata into the ConfigTree field if present.
	if err := l.setConfigTreeField(elem); err != nil {
		// This is a non-critical optional feature; failures are suppressed.
	}

	return nil
}

// fillStruct performs a recursive traversal of the struct's fields, 
// using reflection to introspect metadata and assign values from the tree.
func (l *Loader) fillStruct(v reflect.Value, prefix string) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields as they cannot be mutated via reflection.
		if !field.CanSet() {
			continue
		}

		// ConfigTree is a reserved field for internal metadata injection.
		if fieldType.Name == "ConfigTree" {
			continue
		}

		// Resolve the canonical configuration path based on struct tags.
		path := l.getFieldPath(prefix, fieldType)

		// Dispatch hydration logic based on field kind.
		if err := l.fillField(field, path); err != nil {
			return fmt.Errorf("field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// fillField resolves the value for a specific path and performs type-safe 
// assignment to the target reflect.Value.
func (l *Loader) fillField(field reflect.Value, path string) error {
	node := l.tree.Get(path)
	if node == nil {
		// For nested structs, attempt recursive traversal even if the specific node is nil 
		// to support partial configuration structures.
		if field.Kind() == reflect.Struct {
			return l.fillStruct(field, path)
		}
		return nil
	}

	if !node.HasValue() && field.Kind() != reflect.Struct {
		return nil
	}

	value := node.GetValue()

	// Type-driven assignment with basic coercion support.
	switch field.Kind() {
	case reflect.String:
		if v, ok := value.(string); ok {
			field.SetString(v)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, ok := value.(int); ok {
			field.SetInt(int64(v))
		}

	case reflect.Float32, reflect.Float64:
		if v, ok := value.(float64); ok {
			field.SetFloat(v)
		}

	case reflect.Bool:
		if v, ok := value.(bool); ok {
			field.SetBool(v)
		}

	case reflect.Slice:
		return l.fillSlice(field, node)

	case reflect.Struct:
		// Recurse into nested configuration objects.
		return l.fillStruct(field, path)

	default:
		// Unsupported kinds are ignored for robustness.
	}

	return nil
}

// fillSlice allocates and populates a Go slice from an array-typed configuration node.
func (l *Loader) fillSlice(field reflect.Value, node *tree.ConfigNode) error {
	if node.Type != tree.TypeArray || !node.HasValue() {
		return nil
	}

	arrValue, ok := node.GetValue().([]interface{})
	if !ok {
		return nil
	}

	elemType := field.Type().Elem()
	slice := reflect.MakeSlice(field.Type(), len(arrValue), len(arrValue))

	for i, item := range arrValue {
		elem := slice.Index(i)
		if err := l.setReflectValue(elem, item, elemType); err != nil {
			return err
		}
	}

	field.Set(slice)
	return nil
}

// setReflectValue performs deep recursive assignment for complex types 
// like nested slices and anonymous maps.
func (l *Loader) setReflectValue(v reflect.Value, val interface{}, t reflect.Type) error {
	switch t.Kind() {
	case reflect.String:
		if s, ok := val.(string); ok {
			v.SetString(s)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, ok := val.(int); ok {
			v.SetInt(int64(i))
		}
	case reflect.Float32, reflect.Float64:
		if f, ok := val.(float64); ok {
			v.SetFloat(f)
		}
	case reflect.Bool:
		if b, ok := val.(bool); ok {
			v.SetBool(b)
		}
	case reflect.Slice:
		arr, ok := val.([]interface{})
		if !ok {
			return nil
		}

		elemType := t.Elem()
		slice := reflect.MakeSlice(t, len(arr), len(arr))
		for i, item := range arr {
			if err := l.setReflectValue(slice.Index(i), item, elemType); err != nil {
				return err
			}
		}
		v.Set(slice)
	case reflect.Struct:
		obj, ok := val.(map[string]interface{})
		if !ok {
			return nil
		}

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := t.Field(i)

			if !field.CanSet() || fieldType.Name == "ConfigTree" {
				continue
			}

			mapKey := l.getMapKey(fieldType)
			raw, exists := obj[mapKey]
			if !exists {
				continue
			}

			if err := l.setReflectValue(field, raw, fieldType.Type); err != nil {
				return err
			}
		}
	}
	return nil
}

// getMapKey resolves the lookup key for a field by prioritizing struct tags 
// (json/yaml) over the literal field name.
func (l *Loader) getMapKey(field reflect.StructField) string {
	if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
		return tag
	}

	if tag := field.Tag.Get("yaml"); tag != "" && tag != "-" {
		return tag
	}

	return field.Name
}

// getFieldPath resolves the canonical hierarchical path for a field by 
// inspecting metadata tags and the current traversal prefix.
func (l *Loader) getFieldPath(prefix string, field reflect.StructField) string {
	var tagValue string
	
	// Prioritize JSON tags for path resolution.
	if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
		tagValue = tag
	} else if tag := field.Tag.Get("yaml"); tag != "" && tag != "-" {
		// Fallback to YAML tags.
		tagValue = tag
	} else {
		// Final fallback to the field's literal identifier.
		tagValue = field.Name
	}

	if prefix == "" {
		return tagValue
	}
	return prefix + "." + tagValue
}

// setConfigTreeField attempts to inject the runtime tree reference into the 
// target struct for dynamic query capabilities. Supporting multiple tree 
// types ensures backward compatibility.
func (l *Loader) setConfigTreeField(v reflect.Value) error {
	treeField := v.FieldByName("ConfigTree")
	if !treeField.IsValid() {
		return fmt.Errorf("ConfigTree field not found")
	}

	if !treeField.CanSet() {
		return fmt.Errorf("ConfigTree field is not settable (check exported visibility)")
	}

	internalTreeType := reflect.TypeOf((*tree.ConfigTree)(nil))
	publicTreeType := reflect.TypeOf((*Tree)(nil))

	// Direct injection for internal tree type.
	if treeField.Type() == internalTreeType {
		treeField.Set(reflect.ValueOf(l.tree))
		return nil
	}

	// Wrapped injection for the public API tree representation.
	if treeField.Type() == publicTreeType {
		treeField.Set(reflect.ValueOf(wrapTree(l.tree)))
		return nil
	}

	return fmt.Errorf("unsupported type for ConfigTree field: %s", treeField.Type())
}

// GetTree exposes the underlying configuration tree for advanced manipulation.
func (l *Loader) GetTree() *tree.ConfigTree {
	return l.tree
}