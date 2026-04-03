package tree

import (
	"fmt"
	"strings"
)

// ConfigTree is a Trie tree structure for configuration data
type ConfigTree struct {
	Root *ConfigNode
}

func NewConfigTree() *ConfigTree {
	root := NewConfigNode("root")
	root.Type = TypeObject // Root is always an object
	return &ConfigTree{
		Root: root,
	}
}

// Get gets the node by path
// path format: "server.host"
func (t *ConfigTree) Get(path string) *ConfigNode {
	if path == "" {
		return nil
	}
	return t.GetByPath(strings.Split(path, "."))
}

// GetByPath gets the node by path
// path format: []string{"server", "host"}
func (t *ConfigTree) GetByPath(path []string) *ConfigNode {
	node := t.Root
	for _, key := range path {
		if key == "" {
			continue
		}

		child, ok := node.GetChild(key)
		if !ok {
			return nil
		}
		node = child
	}
	return node
}

// GetValue gets the value of the node by path
func (t *ConfigTree) GetValue(path string) (any, bool) {
	node := t.Get(path)
	if node == nil || !node.HasValue() {
		return nil, false
	}
	return node.GetValue(), true
}

// Set sets the value at the specified path and source
// If the path does not exist, intermediate nodes will be created
func (t *ConfigTree) Set(path string, value any, source SourceType, valueType ValueType) error {
	return t.SetByPath(strings.Split(path, "."), value, source, valueType)
}

// SetByPath sets the value using a path array
func (t *ConfigTree) SetByPath(path []string, value any, source SourceType, valueType ValueType) error {
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	// Traverse the path and create missing nodes
	node := t.Root
	for i, key := range path {
		if key == "" {
			continue
		}

		child, ok := node.GetChild(key)
		if !ok {
			child = NewConfigNode(key)

			// if child is the last key, set it to the actual type; 
			// otherwise, set it to object
			if i == len(path)-1 {
				child.Type = valueType
			} else {
				child.Type = TypeObject
			}

			node.AddChild(child)
		}
		node = child
	}

	node.SetValue(value, source)
	return nil
}

// Merge is used to merge another ConfigTree into this one
// with a specified source type for the new values
func (t *ConfigTree) Merge(other *ConfigTree, source SourceType) {
	t.mergeNode(t.Root, other.Root, source)
}

func (t *ConfigTree) mergeNode(target, source *ConfigNode, sourceType SourceType) {
	if source.HasValue() {
		target.SetValue(source.GetValue(), sourceType)
		target.Type = source.Type
	}

	// recursively merge child nodes
	for key, sourceChild := range source.Children {
		targetChild, ok := target.GetChild(key)
		if !ok {
			// target node does not exist, copy directly
			targetChild = copyNode(sourceChild, sourceType)
			target.AddChild(targetChild)
		} else {
			// recursively merge
			t.mergeNode(targetChild, sourceChild, sourceType)
		}
	}

	// merge array elements
	if source.IsArray() {
		target.Items = make([]*ConfigNode, len(source.Items))
		for i, item := range source.Items {
			target.Items[i] = copyNode(item, sourceType)
		}
	}
}

// copyNode (deepcopy)
func copyNode(source *ConfigNode, sourceType SourceType) *ConfigNode {
	node := NewConfigNode(source.Key)
	node.Type = source.Type

	if source.HasValue() {
		node.SetValue(source.GetValue(), sourceType)
	}

	for _, child := range source.Children {
		node.AddChild(copyNode(child, sourceType))
	}

	for _, item := range source.Items {
		node.AddItem(copyNode(item, sourceType))
	}

	return node
}

// GetAllWithPrefix gets all nodes with the specified prefix
func (t *ConfigTree) GetAllWithPrefix(prefix string) map[string]*ConfigNode {
	node := t.Get(prefix)
	if node == nil {
		return nil
	}

	result := make(map[string]*ConfigNode)
	for key, child := range node.Children {
		fullPath := prefix + "." + key
		result[fullPath] = child
	}
	return result
}

// Walk traverses the entire tree
func (t *ConfigTree) Walk(fn func(path string, node *ConfigNode)) {
	t.walkNode(t.Root, "", fn)
}

// walkNode recursively traverses nodes
func (t *ConfigTree) walkNode(node *ConfigNode, path string, fn func(string, *ConfigNode)) {
	if path != "" { // skip root node
		fn(path, node)
	}

	for key, child := range node.Children {
		childPath := path
		if childPath != "" {
			childPath += "."
		}
		childPath += key
		t.walkNode(child, childPath, fn)
	}
}

// ToMap converts the tree to a nested map (for serialization)
// Only returns the highest priority values
func (t *ConfigTree) ToMap() map[string]any {
	return t.nodeToMap(t.Root)
}

func (t *ConfigTree) nodeToMap(node *ConfigNode) map[string]any {
	if node.IsPrimitive() {
		return nil // 基本类型直接返回值
	}

	result := make(map[string]any)

	if node.IsObject() {
		for key, child := range node.Children {
			if child.IsPrimitive() {
				result[key] = child.GetValue()
			} else if child.IsObject() {
				result[key] = t.nodeToMap(child)
			} else if child.IsArray() {
				result[key] = t.nodeToArray(child)
			}
		}
	}

	return result
}

// nodeToArray converts an array node
func (t *ConfigTree) nodeToArray(node *ConfigNode) []any {
	result := make([]any, len(node.Items))
	for i, item := range node.Items {
		if item.IsPrimitive() {
			result[i] = item.GetValue()
		} else if item.IsObject() {
			result[i] = t.nodeToMap(item)
		} else if item.IsArray() {
			result[i] = t.nodeToArray(item)
		}
	}
	return result
}

func (t *ConfigTree) Print() {
	t.printNode(t.Root, "", true)
}

func (t *ConfigTree) printNode(node *ConfigNode, prefix string, isLast bool) {
	if node.Key != "root" {
		marker := "├─"
		if isLast {
			marker = "└─"
		}
		fmt.Printf("%s%s %s\n", prefix, marker, node.String())

		if isLast {
			prefix += "  "
		} else {
			prefix += "│ "
		}
	}

	children := make([]*ConfigNode, 0, len(node.Children))
	for _, child := range node.Children {
		children = append(children, child)
	}

	for i, child := range children {
		t.printNode(child, prefix, i == len(children)-1)
	}
}
