package parser

import (
	"testing"
)

func TestNewConfigNode(t *testing.T) {
	node := NewConfigNode("test_key")

	if node.Key != "test_key" {
		t.Errorf("Expected key 'test_key', got '%s'", node.Key)
	}

	if node.Children == nil {
		t.Error("Children map should be initialized")
	}

	if node.Items == nil {
		t.Error("Items slice should be initialized")
	}

	if len(node.Children) != 0 {
		t.Errorf("Children should be empty, got %d items", len(node.Children))
	}

	if len(node.Items) != 0 {
		t.Errorf("Items should be empty, got %d items", len(node.Items))
	}
}

func TestConfigNode_AddChild(t *testing.T) {
	parent := NewConfigNode("parent")
	child1 := NewConfigNode("child1")
	child2 := NewConfigNode("child2")

	parent.AddChild(child1)
	parent.AddChild(child2)

	if len(parent.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(parent.Children))
	}

	if parent.Children["child1"] != child1 {
		t.Error("child1 not found in Children map")
	}

	if parent.Children["child2"] != child2 {
		t.Error("child2 not found in Children map")
	}
}

func TestConfigNode_AddItem(t *testing.T) {
	arrayNode := NewConfigNode("array")
	item1 := NewConfigNode("[0]")
	item1.Type = TypeString
	item1.Value = "first"

	item2 := NewConfigNode("[1]")
	item2.Type = TypeString
	item2.Value = "second"

	arrayNode.AddItem(item1)
	arrayNode.AddItem(item2)

	if len(arrayNode.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(arrayNode.Items))
	}

	if arrayNode.Items[0] != item1 {
		t.Error("item1 not at index 0")
	}

	if arrayNode.Items[1] != item2 {
		t.Error("item2 not at index 1")
	}
}

func TestConfigNode_IsObject(t *testing.T) {
	tests := []struct {
		name     string
		nodeType ValueType
		expected bool
	}{
		{"Object type", TypeObject, true},
		{"String type", TypeString, false},
		{"Array type", TypeArray, false},
		{"Int type", TypeInt, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewConfigNode("test")
			node.Type = tt.nodeType

			if got := node.IsObject(); got != tt.expected {
				t.Errorf("IsObject() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConfigNode_IsArray(t *testing.T) {
	tests := []struct {
		name     string
		nodeType ValueType
		expected bool
	}{
		{"Array type", TypeArray, true},
		{"Object type", TypeObject, false},
		{"String type", TypeString, false},
		{"Int type", TypeInt, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewConfigNode("test")
			node.Type = tt.nodeType

			if got := node.IsArray(); got != tt.expected {
				t.Errorf("IsArray() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConfigNode_IsPrimitive(t *testing.T) {
	tests := []struct {
		name     string
		nodeType ValueType
		expected bool
	}{
		{"String type", TypeString, true},
		{"Int type", TypeInt, true},
		{"Float type", TypeFloat, true},
		{"Bool type", TypeBool, true},
		{"Null type", TypeNull, true},
		{"Object type", TypeObject, false},
		{"Array type", TypeArray, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewConfigNode("test")
			node.Type = tt.nodeType

			if got := node.IsPrimitive(); got != tt.expected {
				t.Errorf("IsPrimitive() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConfigNode_ComplexTree(t *testing.T) {
	// 构建一个复杂的树结构
	// root
	//   ├── server (object)
	//   │   ├── host (string)
	//   │   └── port (int)
	//   └── features (array)
	//       ├── [0] (string)
	//       └── [1] (string)

	root := NewConfigNode("root")
	root.Type = TypeObject

	// 添加 server 对象
	server := NewConfigNode("server")
	server.Type = TypeObject

	host := NewConfigNode("host")
	host.Type = TypeString
	host.Value = "localhost"

	port := NewConfigNode("port")
	port.Type = TypeInt
	port.Value = 8080

	server.AddChild(host)
	server.AddChild(port)
	root.AddChild(server)

	// 添加 features 数组
	features := NewConfigNode("features")
	features.Type = TypeArray

	item1 := NewConfigNode("[0]")
	item1.Type = TypeString
	item1.Value = "cache"

	item2 := NewConfigNode("[1]")
	item2.Type = TypeString
	item2.Value = "metrics"

	features.AddItem(item1)
	features.AddItem(item2)
	root.AddChild(features)

	// 验证树结构
	if len(root.Children) != 2 {
		t.Errorf("Root should have 2 children, got %d", len(root.Children))
	}

	if !root.Children["server"].IsObject() {
		t.Error("server should be an object")
	}

	if len(root.Children["server"].Children) != 2 {
		t.Errorf("server should have 2 children, got %d", len(root.Children["server"].Children))
	}

	if root.Children["server"].Children["host"].Value != "localhost" {
		t.Errorf("host value should be 'localhost', got '%v'", root.Children["server"].Children["host"].Value)
	}

	if root.Children["server"].Children["port"].Value != 8080 {
		t.Errorf("port value should be 8080, got %v", root.Children["server"].Children["port"].Value)
	}

	if !root.Children["features"].IsArray() {
		t.Error("features should be an array")
	}

	if len(root.Children["features"].Items) != 2 {
		t.Errorf("features should have 2 items, got %d", len(root.Children["features"].Items))
	}
}
