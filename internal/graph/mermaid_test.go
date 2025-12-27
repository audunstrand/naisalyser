package graph

import (
	"strings"
	"testing"
)

// TestSanitizeMermaidID_Hyphens tests that hyphens are replaced with underscores
func TestSanitizeMermaidID_Hyphens(t *testing.T) {
	input := "my-app-name"
	expected := "my_app_name"
	result := sanitizeMermaidID(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestSanitizeMermaidID_Dots tests that dots are replaced with underscores
func TestSanitizeMermaidID_Dots(t *testing.T) {
	input := "my.app.name"
	expected := "my_app_name"
	result := sanitizeMermaidID(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestSanitizeMermaidID_Slashes tests that slashes are replaced with underscores
func TestSanitizeMermaidID_Slashes(t *testing.T) {
	input := "namespace/app"
	expected := "namespace_app"
	result := sanitizeMermaidID(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestSanitizeMermaidID_Colons tests that colons are replaced with underscores
func TestSanitizeMermaidID_Colons(t *testing.T) {
	input := "app:8080"
	expected := "app_8080"
	result := sanitizeMermaidID(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestSanitizeMermaidID_Mixed tests multiple special characters together
func TestSanitizeMermaidID_Mixed(t *testing.T) {
	input := "my-app.service/api:v1"
	expected := "my_app_service_api_v1"
	result := sanitizeMermaidID(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestSanitizeMermaidID_Empty tests empty string handling
func TestSanitizeMermaidID_Empty(t *testing.T) {
	input := ""
	expected := ""
	result := sanitizeMermaidID(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestSanitizeMermaidID_NoSpecialChars tests that normal strings pass through unchanged
func TestSanitizeMermaidID_NoSpecialChars(t *testing.T) {
	input := "myapp123"
	expected := "myapp123"
	result := sanitizeMermaidID(input)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestToMermaid_AppNode tests rendering of a basic app node
func TestToMermaid_AppNode(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "test-app", Type: "app", Stateful: false},
		},
		Edges: []Edge{},
	}
	result := graph.ToMermaid()
	
	if !strings.Contains(result, "graph LR") {
		t.Error("Expected output to contain 'graph LR'")
	}
	if !strings.Contains(result, "test_app[test-app]") {
		t.Error("Expected output to contain app node definition")
	}
	if !strings.Contains(result, "class test_app stateless") {
		t.Error("Expected stateless app to have stateless class")
	}
}

// TestToMermaid_StatefulApp tests rendering of a stateful app node
func TestToMermaid_StatefulApp(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "db-app", Type: "app", Stateful: true},
		},
		Edges: []Edge{},
	}
	result := graph.ToMermaid()
	
	if !strings.Contains(result, "db_app[(db-app)]") {
		t.Error("Expected stateful app to use cylinder shape [(name)]")
	}
	if !strings.Contains(result, "class db_app stateful") {
		t.Error("Expected stateful app to have stateful class")
	}
}

// TestToMermaid_StatelessApp tests rendering of a stateless app node
func TestToMermaid_StatelessApp(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "api-app", Type: "app", Stateful: false},
		},
		Edges: []Edge{},
	}
	result := graph.ToMermaid()
	
	if !strings.Contains(result, "api_app[api-app]") {
		t.Error("Expected stateless app to use box shape [name]")
	}
	if !strings.Contains(result, "class api_app stateless") {
		t.Error("Expected stateless app to have stateless class")
	}
}

// TestToMermaid_KafkaNode tests rendering of a Kafka topic node
func TestToMermaid_KafkaNode(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "my-topic", Type: "kafka"},
		},
		Edges: []Edge{},
	}
	result := graph.ToMermaid()
	
	if !strings.Contains(result, "my_topic{{my-topic}}") {
		t.Error("Expected Kafka node to use hexagon shape {{name}}")
	}
	if !strings.Contains(result, "class my_topic kafka") {
		t.Error("Expected Kafka node to have kafka class")
	}
}

// TestToMermaid_ExternalNode tests rendering of an external service node
func TestToMermaid_ExternalNode(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "api.example.com", Type: "external"},
		},
		Edges: []Edge{},
	}
	result := graph.ToMermaid()
	
	if !strings.Contains(result, "api_example_com>api.example.com]") {
		t.Error("Expected external node to use flag shape >name]")
	}
	if !strings.Contains(result, "class api_example_com external") {
		t.Error("Expected external node to have external class")
	}
}

// TestToMermaid_CallEdge tests rendering of a call edge
func TestToMermaid_CallEdge(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "app1", Type: "app"},
			{ID: "app2", Type: "app"},
		},
		Edges: []Edge{
			{From: "app1", To: "app2", Type: "call"},
		},
	}
	result := graph.ToMermaid()
	
	if !strings.Contains(result, "app1 --> app2") {
		t.Error("Expected call edge to use solid arrow -->")
	}
}

// TestToMermaid_KafkaEdge tests rendering of a Kafka edge
func TestToMermaid_KafkaEdge(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "producer", Type: "app"},
			{ID: "topic", Type: "kafka"},
		},
		Edges: []Edge{
			{From: "producer", To: "topic", Type: "kafka"},
		},
	}
	result := graph.ToMermaid()
	
	if !strings.Contains(result, "producer -.->|kafka| topic") {
		t.Error("Expected Kafka edge to use dashed arrow with label -.->|kafka|")
	}
}

// TestToMermaid_ExternalEdge tests rendering of an external edge
func TestToMermaid_ExternalEdge(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "app", Type: "app"},
			{ID: "external-api", Type: "external"},
		},
		Edges: []Edge{
			{From: "app", To: "external-api", Type: "external"},
		},
	}
	result := graph.ToMermaid()
	
	if !strings.Contains(result, "app -->|ext| external_api") {
		t.Error("Expected external edge to use solid arrow with label -->|ext|")
	}
}

// TestToMermaid_MultipleNodes tests rendering of multiple nodes of different types
func TestToMermaid_MultipleNodes(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "app1", Type: "app", Stateful: false},
			{ID: "app2", Type: "app", Stateful: true},
			{ID: "topic1", Type: "kafka"},
			{ID: "api.external.com", Type: "external"},
		},
		Edges: []Edge{},
	}
	result := graph.ToMermaid()
	
	// Check all nodes are present
	if !strings.Contains(result, "app1[app1]") {
		t.Error("Expected stateless app node")
	}
	if !strings.Contains(result, "app2[(app2)]") {
		t.Error("Expected stateful app node")
	}
	if !strings.Contains(result, "topic1{{topic1}}") {
		t.Error("Expected Kafka node")
	}
	if !strings.Contains(result, "api_external_com>api.external.com]") {
		t.Error("Expected external node")
	}
}

// TestToMermaid_ComplexGraph tests a complete graph with nodes and edges
func TestToMermaid_ComplexGraph(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "frontend", Type: "app", Stateful: false},
			{ID: "backend", Type: "app", Stateful: true},
			{ID: "events", Type: "kafka"},
			{ID: "api.github.com", Type: "external"},
		},
		Edges: []Edge{
			{From: "frontend", To: "backend", Type: "call"},
			{From: "backend", To: "events", Type: "kafka"},
			{From: "backend", To: "api.github.com", Type: "external"},
		},
	}
	result := graph.ToMermaid()
	
	// Verify structure
	if !strings.Contains(result, "graph LR") {
		t.Error("Expected graph to start with 'graph LR'")
	}
	
	// Verify nodes section
	if !strings.Contains(result, "%% Node definitions") {
		t.Error("Expected node definitions section")
	}
	
	// Verify edges section
	if !strings.Contains(result, "%% Edges") {
		t.Error("Expected edges section")
	}
	
	// Verify styling section
	if !strings.Contains(result, "%% Styling") {
		t.Error("Expected styling section")
	}
	
	// Verify all edges are present
	if !strings.Contains(result, "frontend --> backend") {
		t.Error("Expected call edge from frontend to backend")
	}
	if !strings.Contains(result, "backend -.->|kafka| events") {
		t.Error("Expected Kafka edge from backend to events")
	}
	if !strings.Contains(result, "backend -->|ext| api_github_com") {
		t.Error("Expected external edge from backend to api.github.com")
	}
}

// TestToMermaid_StylingClasses tests that all styling classes are defined
func TestToMermaid_StylingClasses(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "app", Type: "app"},
		},
		Edges: []Edge{},
	}
	result := graph.ToMermaid()
	
	// Check all class definitions exist
	if !strings.Contains(result, "classDef stateless fill:#4A90D9") {
		t.Error("Expected stateless class definition")
	}
	if !strings.Contains(result, "classDef stateful fill:#F5A623") {
		t.Error("Expected stateful class definition")
	}
	if !strings.Contains(result, "classDef kafka fill:#7B68EE") {
		t.Error("Expected kafka class definition")
	}
	if !strings.Contains(result, "classDef external fill:#808080") {
		t.Error("Expected external class definition")
	}
}

// TestToMermaid_ValidMermaidSyntax tests that output follows basic Mermaid syntax rules
func TestToMermaid_ValidMermaidSyntax(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "app-1", Type: "app", Stateful: false},
			{ID: "app-2", Type: "app", Stateful: true},
		},
		Edges: []Edge{
			{From: "app-1", To: "app-2", Type: "call"},
		},
	}
	result := graph.ToMermaid()
	
	// Must start with graph declaration
	if !strings.HasPrefix(result, "graph LR") {
		t.Error("Mermaid diagram must start with 'graph LR'")
	}
	
	// Check for proper indentation (4 spaces)
	lines := strings.Split(result, "\n")
	hasIndentedLines := false
	for _, line := range lines {
		if strings.HasPrefix(line, "    ") && len(line) > 4 {
			hasIndentedLines = true
			break
		}
	}
	if !hasIndentedLines {
		t.Error("Expected properly indented lines in Mermaid output")
	}
	
	// Verify sanitized IDs are used (no hyphens in node IDs)
	if strings.Contains(result, "app-1 -->") || strings.Contains(result, "--> app-2") {
		t.Error("Node IDs in edges should be sanitized (no hyphens)")
	}
	
	// Verify original IDs are preserved in labels
	if !strings.Contains(result, "[app-1]") || !strings.Contains(result, "[(app-2)]") {
		t.Error("Original node IDs should be preserved in node labels")
	}
}

// TestToMermaid_EmptyGraph tests rendering of an empty graph
func TestToMermaid_EmptyGraph(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{},
		Edges: []Edge{},
	}
	result := graph.ToMermaid()
	
	// Should still have basic structure
	if !strings.Contains(result, "graph LR") {
		t.Error("Expected graph declaration even for empty graph")
	}
	if !strings.Contains(result, "%% Styling") {
		t.Error("Expected styling section even for empty graph")
	}
}

// TestToMermaid_NodeWithoutEdges tests a node without any edges
func TestToMermaid_NodeWithoutEdges(t *testing.T) {
	graph := &Graph{
		Nodes: []Node{
			{ID: "isolated-app", Type: "app", Stateful: false},
		},
		Edges: []Edge{},
	}
	result := graph.ToMermaid()
	
	if !strings.Contains(result, "isolated_app[isolated-app]") {
		t.Error("Expected isolated node to be rendered")
	}
	if !strings.Contains(result, "class isolated_app stateless") {
		t.Error("Expected isolated node to have correct styling")
	}
}
