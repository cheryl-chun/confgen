package analyzer

import (
	"testing"

	"github.com/cheryl-chun/confgen/internal/parser"
)

// TestAnalyzeSimpleObject verifies the engine's ability to perform basic 
// scalar type inference (string, int, bool) and map them to PascalCase field names.
func TestAnalyzeSimpleObject(t *testing.T) {
	// Scaffolding: Manually construct a primitive ConfigNode tree.
	root := parser.NewConfigNode("root")
	root.Type = parser.TypeObject

	// Inject scalar fields for type resolution testing.
	host := parser.NewConfigNode("host")
	host.Type = parser.TypeString
	host.Value = "localhost"
	root.AddChild(host)

	port := parser.NewConfigNode("port")
	port.Type = parser.TypeInt
	port.Value = 8080
	root.AddChild(port)

	enabled := parser.NewConfigNode("enabled")
	enabled.Type = parser.TypeBool
	enabled.Value = true
	root.AddChild(enabled)

	// Execute the inference logic.
	result, err := Analyze(root)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Assertion: Ensure the primary entry-point struct is correctly initialized.
	if result.RootStruct == nil {
		t.Fatal("RootStruct is nil")
	}

	if result.RootStruct.Name != "Config" {
		t.Errorf("RootStruct.Name = %q, want %q", result.RootStruct.Name, "Config")
	}

	if len(result.RootStruct.Fields) != 3 {
		t.Errorf("len(Fields) = %d, want 3", len(result.RootStruct.Fields))
	}

	// Validate field mapping and type resolution accuracy.
	fields := result.RootStruct.Fields
	expectedFields := map[string]string{
		"Host":    "string",
		"Port":    "int",
		"Enabled": "bool",
	}

	for _, field := range fields {
		expectedType, ok := expectedFields[field.Name]
		if !ok {
			t.Errorf("Unexpected field discovered in result: %s", field.Name)
			continue
		}
		if field.Type != expectedType {
			t.Errorf("Field %s: Type mismatch. Got %q, want %q", field.Name, field.Type, expectedType)
		}
	}
}

// TestAnalyzeNestedObject validates the recursive descent logic by ensuring 
// that nested objects are flattened into a decoupled struct registry.
func TestAnalyzeNestedObject(t *testing.T) {
	// Orchestrate a nested structural fixture.
	root := parser.NewConfigNode("root")
	root.Type = parser.TypeObject

	// Definition of the nested 'server' component.
	server := parser.NewConfigNode("server")
	server.Type = parser.TypeObject
	root.AddChild(server)

	host := parser.NewConfigNode("host")
	host.Type = parser.TypeString
	host.Value = "localhost"
	server.AddChild(host)

	port := parser.NewConfigNode("port")
	port.Type = parser.TypeInt
	port.Value = 8080
	server.AddChild(port)

	result, err := Analyze(root)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verification: The root should contain a reference to the nested struct.
	if len(result.RootStruct.Fields) != 1 {
		t.Errorf("len(RootStruct.Fields) = %d, want 1", len(result.RootStruct.Fields))
	}

	serverField := result.RootStruct.Fields[0]
	if serverField.Name != "Server" {
		t.Errorf("serverField.Name = %q, want %q", serverField.Name, "Server")
	}

	if serverField.Type != "ServerConfig" {
		t.Errorf("serverField.Type = %q, want %q", serverField.Type, "ServerConfig")
	}

	// Ensure the nested struct was correctly registered in the auxiliary map.
	serverStruct, ok := result.SubStructs["ServerConfig"]
	if !ok {
		t.Fatal("ServerConfig failed to register in SubStructs registry")
	}

	if len(serverStruct.Fields) != 2 {
		t.Errorf("len(ServerConfig.Fields) = %d, want 2", len(serverStruct.Fields))
	}
}

// TestAnalyzeArray ensures that sequential collections of primitive types 
// are correctly identified and represented as Go slices.
func TestAnalyzeArray(t *testing.T) {
	root := parser.NewConfigNode("root")
	root.Type = parser.TypeObject

	// Construct an array of scalar strings.
	tags := parser.NewConfigNode("tags")
	tags.Type = parser.TypeArray
	root.AddChild(tags)

	tag1 := parser.NewConfigNode("[0]")
	tag1.Type = parser.TypeString
	tag1.Value = "production"
	tags.AddItem(tag1)

	tag2 := parser.NewConfigNode("[1]")
	tag2.Type = parser.TypeString
	tag2.Value = "web"
	tags.AddItem(tag2)

	result, err := Analyze(root)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Assert the field is represented as a slice of strings.
	if len(result.RootStruct.Fields) != 1 {
		t.Fatalf("len(Fields) = %d, want 1", len(result.RootStruct.Fields))
	}

	tagsField := result.RootStruct.Fields[0]
	if tagsField.Name != "Tags" {
		t.Errorf("tagsField.Name = %q, want %q", tagsField.Name, "Tags")
	}

	if tagsField.Type != "[]string" {
		t.Errorf("tagsField.Type = %q, want %q", tagsField.Type, "[]string")
	}
}

// TestAnalyzeArrayOfObjects tests the most complex scenario: an array containing 
// objects, requiring singularization of the struct name and slice type inference.
func TestAnalyzeArrayOfObjects(t *testing.T) {
	root := parser.NewConfigNode("root")
	root.Type = parser.TypeObject

	// Construct an array of homogeneous objects.
	servers := parser.NewConfigNode("servers")
	servers.Type = parser.TypeArray
	root.AddChild(servers)

	// First element fixture for type sampling.
	server1 := parser.NewConfigNode("[0]")
	server1.Type = parser.TypeObject
	servers.AddItem(server1)

	host1 := parser.NewConfigNode("host")
	host1.Type = parser.TypeString
	host1.Value = "server1"
	server1.AddChild(host1)

	result, err := Analyze(root)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify the collection field is typed as a slice of the singularized struct.
	serversField := result.RootStruct.Fields[0]
	if serversField.Name != "Servers" {
		t.Errorf("serversField.Name = %q, want %q", serversField.Name, "Servers")
	}

	// Expected: []ServerConfig (singularized from "servers").
	if serversField.Type != "[]ServerConfig" {
		t.Errorf("serversField.Type = %q, want %q", serversField.Type, "[]ServerConfig")
	}

	// Verification of the sampled struct metadata.
	if _, ok := result.SubStructs["ServerConfig"]; !ok {
		t.Error("ServerConfig struct metadata missing from registry")
	}
}