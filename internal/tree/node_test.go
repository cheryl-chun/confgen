package tree

import (
	"testing"
)

func TestConfigNode_SetValue_SingleSource(t *testing.T) {
	cases := []struct {
		name     string
		source   SourceType
		value    any
		expected any
	}{
		{"set from file", SourceFile, "localhost", "localhost"},
		{"set from env", SourceSystemEnv, "prod.com", "prod.com"},
		{"set from default", SourceDefault, 8080, 8080},
		{"set nil value", SourceFile, nil, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := NewConfigNode("test")
			node.SetValue(tc.value, tc.source)

			got := node.GetValue()
			if got != tc.expected {
				t.Errorf("GetValue() = %v, want %v", got, tc.expected)
			}

			if !node.HasValue() {
				t.Error("HasValue() should return true")
			}
		})
	}
}

func TestConfigNode_SetValue_MultiSource_Priority(t *testing.T) {
	cases := []struct {
		name     string
		sources  []struct {
			source SourceType
			value  any
		}
		expectedValue  any
		expectedSource SourceType
	}{
		{
			name: "env overrides file",
			sources: []struct {
				source SourceType
				value  any
			}{
				{SourceFile, "localhost"},
				{SourceSystemEnv, "prod.com"},
			},
			expectedValue:  "prod.com",
			expectedSource: SourceSystemEnv,
		},
		{
			name: "code override has highest priority",
			sources: []struct {
				source SourceType
				value  any
			}{
				{SourceFile, "file.com"},
				{SourceSystemEnv, "env.com"},
				{SourceRemote, "remote.com"},
				{SourceCodeOverride, "code.com"},
			},
			expectedValue:  "code.com",
			expectedSource: SourceCodeOverride,
		},
		{
			name: "system env > session env",
			sources: []struct {
				source SourceType
				value  any
			}{
				{SourceSessionEnv, "session.com"},
				{SourceSystemEnv, "system.com"},
			},
			expectedValue:  "system.com",
			expectedSource: SourceSystemEnv,
		},
		{
			name: "file > remote",
			sources: []struct {
				source SourceType
				value  any
			}{
				{SourceRemote, "remote.com"},
				{SourceFile, "local.com"},
			},
			expectedValue:  "local.com",
			expectedSource: SourceFile,
		},
		{
			name: "all sources present",
			sources: []struct {
				source SourceType
				value  any
			}{
				{SourceDefault, "default"},
				{SourceRemote, "remote"},
				{SourceFile, "file"},
				{SourceRuntimeOverride, "runtime"},
				{SourceSessionEnv, "session"},
				{SourceSystemEnv, "system"},
				{SourceCodeOverride, "code"},
			},
			expectedValue:  "code",
			expectedSource: SourceCodeOverride,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := NewConfigNode("host")

			// Set values from different sources
			for _, src := range tc.sources {
				node.SetValue(src.value, src.source)
			}

			// Check highest priority value
			got := node.GetValue()
			if got != tc.expectedValue {
				t.Errorf("GetValue() = %v, want %v", got, tc.expectedValue)
			}

			// Check that the highest priority source is correct
			values := node.GetAllValues()
			if len(values) != len(tc.sources) {
				t.Errorf("GetAllValues() length = %d, want %d", len(values), len(tc.sources))
			}

			if values[0].Source != tc.expectedSource {
				t.Errorf("Highest priority source = %v, want %v", values[0].Source, tc.expectedSource)
			}

			// Verify priority order (descending)
			for i := 0; i < len(values)-1; i++ {
				if values[i].Priority < values[i+1].Priority {
					t.Errorf("Values not sorted by priority: values[%d].Priority=%d < values[%d].Priority=%d",
						i, values[i].Priority, i+1, values[i+1].Priority)
				}
			}
		})
	}
}

func TestConfigNode_SetValue_UpdateExisting(t *testing.T) {
	node := NewConfigNode("test")

	// First set
	node.SetValue("initial", SourceFile)
	if node.GetValue() != "initial" {
		t.Errorf("Initial value = %v, want 'initial'", node.GetValue())
	}

	// Update same source
	node.SetValue("updated", SourceFile)
	if node.GetValue() != "updated" {
		t.Errorf("Updated value = %v, want 'updated'", node.GetValue())
	}

	// Should still have only 1 value
	if len(node.Values) != 1 {
		t.Errorf("Values length = %d, want 1", len(node.Values))
	}
}

func TestConfigNode_GetValueFromSource(t *testing.T) {
	node := NewConfigNode("test")
	node.SetValue("file-value", SourceFile)
	node.SetValue("env-value", SourceSystemEnv)
	node.SetValue("remote-value", SourceRemote)

	cases := []struct {
		name        string
		source      SourceType
		wantValue   any
		wantExists  bool
	}{
		{"get file source", SourceFile, "file-value", true},
		{"get env source", SourceSystemEnv, "env-value", true},
		{"get remote source", SourceRemote, "remote-value", true},
		{"get non-existent source", SourceCodeOverride, nil, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, exists := node.GetValueFromSource(tc.source)
			if exists != tc.wantExists {
				t.Errorf("exists = %v, want %v", exists, tc.wantExists)
			}
			if got != tc.wantValue {
				t.Errorf("value = %v, want %v", got, tc.wantValue)
			}
		})
	}
}

func TestConfigNode_RemoveSource(t *testing.T) {
	cases := []struct {
		name               string
		initialSources     []SourceType
		removeSource       SourceType
		expectRemoved      bool
		expectedRemaining  int
		expectedHighestVal any
	}{
		{
			name:               "remove existing source",
			initialSources:     []SourceType{SourceFile, SourceSystemEnv},
			removeSource:       SourceFile,
			expectRemoved:      true,
			expectedRemaining:  1,
			expectedHighestVal: "env",
		},
		{
			name:               "remove non-existent source",
			initialSources:     []SourceType{SourceFile},
			removeSource:       SourceSystemEnv,
			expectRemoved:      false,
			expectedRemaining:  1,
			expectedHighestVal: "file",
		},
		{
			name:               "remove highest priority source",
			initialSources:     []SourceType{SourceFile, SourceSystemEnv, SourceCodeOverride},
			removeSource:       SourceCodeOverride,
			expectRemoved:      true,
			expectedRemaining:  2,
			expectedHighestVal: "env",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := NewConfigNode("test")

			// Set initial values
			for _, src := range tc.initialSources {
				var val string
				switch src {
				case SourceFile:
					val = "file"
				case SourceSystemEnv:
					val = "env"
				case SourceCodeOverride:
					val = "code"
				}
				node.SetValue(val, src)
			}

			// Remove source
			removed := node.RemoveSource(tc.removeSource)
			if removed != tc.expectRemoved {
				t.Errorf("RemoveSource() = %v, want %v", removed, tc.expectRemoved)
			}

			// Check remaining count
			if len(node.Values) != tc.expectedRemaining {
				t.Errorf("Remaining values = %d, want %d", len(node.Values), tc.expectedRemaining)
			}

			// Check highest priority value after removal
			if node.HasValue() {
				got := node.GetValue()
				if got != tc.expectedHighestVal {
					t.Errorf("After removal, GetValue() = %v, want %v", got, tc.expectedHighestVal)
				}
			}
		})
	}
}

func TestConfigNode_AddChild(t *testing.T) {
	parent := NewConfigNode("parent")
	child1 := NewConfigNode("child1")
	child2 := NewConfigNode("child2")

	parent.AddChild(child1)
	parent.AddChild(child2)

	if len(parent.Children) != 2 {
		t.Errorf("Children count = %d, want 2", len(parent.Children))
	}

	if got, ok := parent.GetChild("child1"); !ok || got != child1 {
		t.Error("child1 not found or incorrect")
	}

	if got, ok := parent.GetChild("child2"); !ok || got != child2 {
		t.Error("child2 not found or incorrect")
	}

	if _, ok := parent.GetChild("nonexistent"); ok {
		t.Error("GetChild should return false for non-existent child")
	}
}

func TestConfigNode_AddItem(t *testing.T) {
	arrayNode := NewConfigNode("array")
	arrayNode.Type = TypeArray

	item1 := NewConfigNode("[0]")
	item2 := NewConfigNode("[1]")

	arrayNode.AddItem(item1)
	arrayNode.AddItem(item2)

	if len(arrayNode.Items) != 2 {
		t.Errorf("Items count = %d, want 2", len(arrayNode.Items))
	}

	if arrayNode.Items[0] != item1 {
		t.Error("Item[0] incorrect")
	}

	if arrayNode.Items[1] != item2 {
		t.Error("Item[1] incorrect")
	}
}

func TestConfigNode_TypeChecks(t *testing.T) {
	cases := []struct {
		name          string
		nodeType      ValueType
		isObject      bool
		isArray       bool
		isPrimitive   bool
	}{
		{"object type", TypeObject, true, false, false},
		{"array type", TypeArray, false, true, false},
		{"string type", TypeString, false, false, true},
		{"int type", TypeInt, false, false, true},
		{"float type", TypeFloat, false, false, true},
		{"bool type", TypeBool, false, false, true},
		{"null type", TypeNull, false, false, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := NewConfigNode("test")
			node.Type = tc.nodeType

			if node.IsObject() != tc.isObject {
				t.Errorf("IsObject() = %v, want %v", node.IsObject(), tc.isObject)
			}
			if node.IsArray() != tc.isArray {
				t.Errorf("IsArray() = %v, want %v", node.IsArray(), tc.isArray)
			}
			if node.IsPrimitive() != tc.isPrimitive {
				t.Errorf("IsPrimitive() = %v, want %v", node.IsPrimitive(), tc.isPrimitive)
			}
		})
	}
}

func TestConfigNode_String(t *testing.T) {
	cases := []struct {
		name         string
		setupNode    func() *ConfigNode
		expectSubstr string
	}{
		{
			name: "node with value",
			setupNode: func() *ConfigNode {
				node := NewConfigNode("host")
				node.Type = TypeString
				node.SetValue("localhost", SourceFile)
				return node
			},
			expectSubstr: "localhost",
		},
		{
			name: "node without value",
			setupNode: func() *ConfigNode {
				node := NewConfigNode("empty")
				node.Type = TypeObject
				return node
			},
			expectSubstr: "<no value>",
		},
		{
			name: "node with source info",
			setupNode: func() *ConfigNode {
				node := NewConfigNode("port")
				node.Type = TypeInt
				node.SetValue(8080, SourceSystemEnv)
				return node
			},
			expectSubstr: "SystemEnv",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.setupNode()
			str := node.String()

			if str == "" {
				t.Error("String() should not be empty")
			}

			// Just check that the expected substring is present
			// (full string matching would be too brittle)
			if tc.expectSubstr != "" {
				// Simple substring check
				found := false
				for i := 0; i <= len(str)-len(tc.expectSubstr); i++ {
					if str[i:i+len(tc.expectSubstr)] == tc.expectSubstr {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("String() = %q, should contain %q", str, tc.expectSubstr)
				}
			}
		})
	}
}

func TestConfigNode_HasValue(t *testing.T) {
	cases := []struct {
		name      string
		setupNode func() *ConfigNode
		want      bool
	}{
		{
			name: "node with value",
			setupNode: func() *ConfigNode {
				node := NewConfigNode("test")
				node.SetValue("value", SourceFile)
				return node
			},
			want: true,
		},
		{
			name: "new node without value",
			setupNode: func() *ConfigNode {
				return NewConfigNode("test")
			},
			want: false,
		},
		{
			name: "node after removing all sources",
			setupNode: func() *ConfigNode {
				node := NewConfigNode("test")
				node.SetValue("value", SourceFile)
				node.RemoveSource(SourceFile)
				return node
			},
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.setupNode()
			if got := node.HasValue(); got != tc.want {
				t.Errorf("HasValue() = %v, want %v", got, tc.want)
			}
		})
	}
}
