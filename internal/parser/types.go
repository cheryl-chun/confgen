package parser

// Unify the data structures of different formats
type ValueType int

const (
	TypeString ValueType = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeArray
	TypeObject
	TypeNull
)

// ConfigValue represents the intermediate representation of a configuration value
// Supports basic types, nested objects, and arrays
type ConfigValue struct {
	Type  ValueType
	Value any
}

// ConfigNode represents a node in the 
// configuration tree (used to build the 
// configuration structure)
type ConfigNode struct {
	Key      string                 // Field name
	Value    any                    // Value (basic type, array, or map)
	Type     ValueType              // Value type
	Children map[string]*ConfigNode // Child nodes (used for object types)
	Items    []*ConfigNode          // Array elements (used for array types)
}

// ParseResult represents the result of parsing
type ParseResult struct {
	Root *ConfigNode // Root node of the configuration tree
	Raw  any         // Original parsed map/slice
}

// NewConfigNode creates a new configuration node
func NewConfigNode(key string) *ConfigNode {
	return &ConfigNode{
		Key:      key,
		Children: make(map[string]*ConfigNode),
		Items:    make([]*ConfigNode, 0),
	}
}

func (n *ConfigNode) AddChild(child *ConfigNode) {
	if n.Children == nil {
		n.Children = make(map[string]*ConfigNode)
	}
	n.Children[child.Key] = child
}

func (n *ConfigNode) AddItem(item *ConfigNode) {
	n.Items = append(n.Items, item)
}

func (n *ConfigNode) IsObject() bool {
	return n.Type == TypeObject
}

func (n *ConfigNode) IsArray() bool {
	return n.Type == TypeArray
}

func (n *ConfigNode) IsPrimitive() bool {
	return n.Type == TypeString || n.Type == TypeInt ||
		n.Type == TypeFloat || n.Type == TypeBool || n.Type == TypeNull
}
