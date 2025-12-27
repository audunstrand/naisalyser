# Test Task T7: Graph Mermaid Output Tests

## Description
Implement comprehensive tests for Mermaid diagram generation:
- `ToMermaid()` - generates Mermaid diagram syntax
- `sanitizeMermaidID()` - sanitizes node IDs for Mermaid

## File Location
`internal/graph/mermaid_test.go`

## Test Cases Required

### SanitizeMermaidID Tests (8 tests)
1. `TestSanitizeMermaidID_Hyphens` - Replaces hyphens with underscores
2. `TestSanitizeMermaidID_Dots` - Replaces dots with underscores
3. `TestSanitizeMermaidID_Slashes` - Replaces slashes with underscores
4. `TestSanitizeMermaidID_Colons` - Replaces colons with underscores
5. `TestSanitizeMermaidID_MultipleSpecialChars` - Replaces multiple chars
6. `TestSanitizeMermaidID_Empty` - Handles empty string
7. `TestSanitizeMermaidID_AlphanumericOnly` - Leaves valid IDs unchanged
8. `TestSanitizeMermaidID_MixedCharacters` - Complex real-world ID

### ToMermaid_Nodes Tests (6 tests)
1. `TestToMermaid_EmptyGraph` - Empty graph produces valid syntax
2. `TestToMermaid_SingleAppNode` - Single app node renders correctly
3. `TestToMermaid_StatelessAppNode` - Stateless app uses square brackets
4. `TestToMermaid_StatefulAppNode` - Stateful app uses cylinder shape
5. `TestToMermaid_KafkaNode` - Kafka node uses double braces
6. `TestToMermaid_ExternalNode` - External node uses right-pointing shape

### ToMermaid_Edges Tests (6 tests)
1. `TestToMermaid_CallEdge` - Call edge uses solid arrow
2. `TestToMermaid_KafkaEdge` - Kafka edge uses dashed arrow with label
3. `TestToMermaid_ExternalEdge` - External edge uses solid arrow with label
4. `TestToMermaid_MultipleEdges` - Multiple edges render all
5. `TestToMermaid_SelfLoopEdge` - Edge from node to itself
6. `TestToMermaid_NoEdges` - Graph with nodes but no edges

### ToMermaid_Styling Tests (4 tests)
1. `TestToMermaid_StatelessStyling` - Stateless apps get "stateless" class
2. `TestToMermaid_StatefulStyling` - Stateful apps get "stateful" class
3. `TestToMermaid_KafkaStyling` - Kafka nodes get "kafka" class
4. `TestToMermaid_ExternalStyling` - External nodes get "external" class

### ToMermaid_Complex Tests (5 tests)
1. `TestToMermaid_RealWorldGraph` - Complex graph with multiple node types
2. `TestToMermaid_HighConnectivity` - Node with many edges
3. `TestToMermaid_LongNodeIDs` - Very long node IDs are sanitized
4. `TestToMermaid_SpecialCharacters` - Node IDs with special chars
5. `TestToMermaid_ValidMermaidSyntax` - Output parses as valid Mermaid

### ToMermaid_Integration Tests (3 tests)
1. `TestToMermaid_FullWorkflow` - Builder → Build → ToMermaid flow
2. `TestToMermaid_UsesTypeSafeConstants` - Verifies NodeType and EdgeType usage
3. `TestToMermaid_OutputAsMarkdown` - Can be embedded in markdown

## Test Data Examples

### Sample Graph for Testing
```go
func createTestGraph() *Graph {
    return &Graph{
        Nodes: []Node{
            {ID: "my-app", Type: NodeTypeApp, Stateful: true},
            {ID: "event-topic", Type: NodeTypeKafka},
            {ID: "external-api", Type: NodeTypeExternal},
        },
        Edges: []Edge{
            {From: "my-app", To: "event-topic", Type: EdgeTypeKafka},
            {From: "my-app", To: "external-api", Type: EdgeTypeExternal},
        },
    }
}
```

### Expected Mermaid Output Example
```
graph LR
    %% Node definitions
    my_app[(my-app)]
    event_topic{{event-topic}}
    external_api>external-api]

    %% Edges
    my_app -.->|kafka| event_topic
    my_app -->|ext| external_api

    %% Styling
    classDef stateless fill:#4A90D9,stroke:#2E5A8C,color:white
    classDef stateful fill:#F5A623,stroke:#C47A00,color:white
    classDef kafka fill:#7B68EE,stroke:#5A4DB2,color:white
    classDef external fill:#808080,stroke:#404040,color:white
    
    class my_app stateful
    class event_topic kafka
    class external_api external
```

## Validation Tests

Add tests that verify output is valid Mermaid:
```go
func TestToMermaid_ValidMermaidSyntax(t *testing.T) {
    g := createTestGraph()
    output := g.ToMermaid()
    
    // Verify starts with "graph LR"
    // Verify has %% comments for sections
    // Verify has node definitions
    // Verify has edge definitions
    // Verify has styling definitions
}
```

## Edge Cases to Test
- Nodes with IDs matching Mermaid keywords
- Cyclic graphs
- Disconnected components
- Very large graphs (100+ nodes)
- Node IDs that become identical after sanitization
- Unicode characters in node IDs
- Nodes with dashes, dots, colons, slashes mixed

## Implementation Notes
- Verify Mermaid keywords are not used as node IDs
- Check output line endings are consistent
- Test both Linux and Windows line endings if applicable
- Verify comment structure (#, %%)
- Check CSS class definitions are valid
- Ensure styling covers all node types

## Success Criteria
- ✅ All 15+ tests pass
- ✅ Code coverage for mermaid functions ≥ 95%
- ✅ Output is valid Mermaid syntax
- ✅ Special characters handled correctly
- ✅ Type-safe constants verified
- ✅ All node types and edge types rendered correctly
