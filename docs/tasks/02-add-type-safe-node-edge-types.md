# Task 02: Add Type-Safe Constants for Node and Edge Types

## Deliverable
Replace string literals for `Node.Type` and `Edge.Type` with typed constants.

## Context
Currently `internal/graph/builder.go` uses raw strings like `"app"`, `"kafka"`, `"external"` for node types and `"call"`, `"kafka"`, `"external"` for edge types. This provides no compile-time safety - a typo like `"ap"` would silently create a wrong node type.

**Current code (builder.go:8-20):**
```go
type Node struct {
    Type string `json:"type"` // "app", "kafka", "external"
}
type Edge struct {
    Type string `json:"type"` // "inbound", "outbound", "database", "kafka", "external"
}
```

## Key Decisions and Principles
- Use Go idiom: `type NodeType string` with `const` block
- Keep JSON serialization unchanged (types serialize as strings)
- Update all usages in same PR to avoid mixed state

## Delivers
Type-safe constants that prevent typos in node/edge type assignments.

## Acceptance Criteria
- `NodeType` type with constants: `NodeTypeApp`, `NodeTypeKafka`, `NodeTypeExternal`
- `EdgeType` type with constants: `EdgeTypeCall`, `EdgeTypeKafka`, `EdgeTypeExternal`
- `Node.Type` field changed from `string` to `NodeType`
- `Edge.Type` field changed from `string` to `EdgeType`
- All usages updated to use constants
- JSON output unchanged (still serializes as strings)

## Dependencies
None - this is a foundational task.

## Related Code

**File to modify:** `internal/graph/builder.go`

**Replace lines 7-20 with:**

```go
// NodeType represents the type of a graph node
type NodeType string

const (
    NodeTypeApp      NodeType = "app"
    NodeTypeKafka    NodeType = "kafka"
    NodeTypeExternal NodeType = "external"
)

// EdgeType represents the type of a graph edge
type EdgeType string

const (
    EdgeTypeCall     EdgeType = "call"
    EdgeTypeKafka    EdgeType = "kafka"
    EdgeTypeExternal EdgeType = "external"
)

// Node represents a node in the dependency graph
type Node struct {
    ID        string   `json:"id"`
    Type      NodeType `json:"type"`
    Namespace string   `json:"namespace,omitempty"`
    Stateful  bool     `json:"stateful,omitempty"`
}

// Edge represents a connection between nodes
type Edge struct {
    From string   `json:"from"`
    To   string   `json:"to"`
    Type EdgeType `json:"type"`
}
```

**Update usages in `AddFromLocalAnalysis` (lines 59-144):**

Replace all string literals:
- `"app"` → `NodeTypeApp`
- `"kafka"` → `NodeTypeKafka`
- `"external"` → `NodeTypeExternal`
- `"call"` → `EdgeTypeCall`
- `"kafka"` (for edges) → `EdgeTypeKafka`
- `"external"` (for edges) → `EdgeTypeExternal`

**Also update `internal/graph/mermaid.go` (lines 17-30):**

Replace string comparisons:
- `case "app":` → `case NodeTypeApp:`
- `case "kafka":` → `case NodeTypeKafka:`
- `case "external":` → `case NodeTypeExternal:`
- `case "call":` → `case EdgeTypeCall:`
- etc.

## Verification

```bash
# Build should succeed
go build ./...

# Test that JSON output is unchanged
cat > /tmp/test_types.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "github.com/navikt/naisalyser/internal/graph"
)

func main() {
    g := &graph.Graph{
        Nodes: []graph.Node{
            {ID: "test-app", Type: graph.NodeTypeApp},
            {ID: "my-topic", Type: graph.NodeTypeKafka},
        },
        Edges: []graph.Edge{
            {From: "test-app", To: "my-topic", Type: graph.EdgeTypeKafka},
        },
    }
    data, _ := json.MarshalIndent(g, "", "  ")
    fmt.Println(string(data))
}
EOF

# The output should show "type": "app" and "type": "kafka" as strings
go run /tmp/test_types.go

# Clean up
rm /tmp/test_types.go
```

Expected JSON output:
```json
{
  "nodes": [
    {"id": "test-app", "type": "app"},
    {"id": "my-topic", "type": "kafka"}
  ],
  "edges": [
    {"from": "test-app", "to": "my-topic", "type": "kafka"}
  ]
}
```

## Files Changed
- `internal/graph/builder.go` (modify type definitions and usages)
- `internal/graph/mermaid.go` (update switch statements)
