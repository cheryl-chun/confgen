package tree

import (
	"fmt"
	"sort"
)

// SourceValue represents a configuration value 
// from a specific source
type SourceValue struct {
	Value    any        // Actual value
	Source   SourceType // Source
	Priority int        // Priority (from SourcePriority)
}

// ConfigNode represents a node in the configuration Trie tree
// This is a multi-source configuration node, each node can store multiple values from different sources
type ConfigNode struct {
	Key  string    // Field name (Trie edge label)
	Type ValueType // Node type

	// Multi-source value storage (sorted by priority, highest priority first)
	// This allows configuration overrides: Env > File > Default
	Values []SourceValue

	// Trie tree child nodes
	Children map[string]*ConfigNode // Object type child nodes
	Items    []*ConfigNode          // Array type elements
}

// NewConfigNode creates a new configuration node
func NewConfigNode(key string) *ConfigNode {
	return &ConfigNode{
		Key:      key,
		Values:   make([]SourceValue, 0, 4), // Preallocate space for 4 sources
		Children: make(map[string]*ConfigNode),
		Items:    make([]*ConfigNode, 0),
	}
}

// SetValue sets the value from a specific source
// If a value from the same source already exists, it updates it; otherwise, it inserts and sorts by priority
func (n *ConfigNode) SetValue(value any, source SourceType) {
	priority := SourcePriority[source]
	sourceValue := SourceValue{
		Value:    value,
		Source:   source,
		Priority: priority,
	}

	// Check if a value from the same source already exists
	for i, sv := range n.Values {
		if sv.Source == source {
			// Update existing value
			n.Values[i] = sourceValue
			return
		}
	}

	// Binary insert to maintain sorted order (descending by priority)
	// Use sort.Search to find insertion position
	// Values are sorted in descending order (highest priority first)
	insertPos := sort.Search(len(n.Values), func(i int) bool {
		// Find first position where priority is less than or equal to new priority
		return n.Values[i].Priority <= priority
	})

	// Insert at the correct position
	n.Values = append(n.Values, SourceValue{})
	copy(n.Values[insertPos+1:], n.Values[insertPos:]) 
	n.Values[insertPos] = sourceValue
}

// GetValue gets the highest priority value
func (n *ConfigNode) GetValue() any {
	if len(n.Values) == 0 {
		return nil
	}
	// The value with the highest priority is always at index 0
	return n.Values[0].Value 
}

// GetValueFromSource gets the value from a specific source
func (n *ConfigNode) GetValueFromSource(source SourceType) (any, bool) {
	for _, sv := range n.Values {
		if sv.Source == source {
			return sv.Value, true
		}
	}
	return nil, false
}

func (n *ConfigNode) GetAllValues() []SourceValue {
	return n.Values
}

// RemoveSource removes the value from a specific source
func (n *ConfigNode) RemoveSource(source SourceType) bool {
	for i, sv := range n.Values {
		if sv.Source == source {
			n.Values = append(n.Values[:i], n.Values[i+1:]...)
			return true
		}
	}
	return false
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

func (n *ConfigNode) GetChild(key string) (*ConfigNode, bool) {
	child, ok := n.Children[key]
	return child, ok
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

func (n *ConfigNode) HasValue() bool {
	return len(n.Values) > 0
}

func (n *ConfigNode) String() string {
	if len(n.Values) == 0 {
		return fmt.Sprintf("%s (%s): <no value>", n.Key, n.Type)
	}
	highestValue := n.Values[0]
	return fmt.Sprintf("%s (%s): %v [from %s]",
		n.Key, n.Type, highestValue.Value, highestValue.Source)
}
