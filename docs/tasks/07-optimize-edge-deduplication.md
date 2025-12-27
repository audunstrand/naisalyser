# Task 07: Optimize Edge Deduplication with Map Lookup

## Deliverable
Replace O(n) linear scan for duplicate edges with O(1) map-based lookup.

## Context
The `addEdge` function in `graph/builder.go` checks for duplicate edges by iterating through all existing edges. For a graph with 1000 edges, adding one edge requires 1000 comparisons. Adding n edges is O(n²).

**Current code (builder.go:153-160):**
```go
func (b *Builder) addEdge(edge Edge) {
    // Check for duplicate edges
    for _, e := range b.edges {
        if e.From == edge.From && e.To == edge.To && e.Type == edge.Type {
            return
        }
    }
    b.edges = append(b.edges, edge)
}
```

## Key Decisions and Principles
- Use composite key struct for map lookup
- Maintain insertion order with slice (for deterministic output)
- Map is private implementation detail

## Delivers
O(1) edge insertion with deduplication.

## Acceptance Criteria
- New `edgeKey` struct type for composite map key
- `Builder` has `edgeSet map[edgeKey]struct{}` field
- `addEdge` checks map before adding to slice
- `NewBuilder` initializes the map
- Behavior unchanged (same edges in output)
- Performance improved for large graphs

## Dependencies
Recommended after Task 02 (type-safe edge types) but not required.

## Related Code

**File to modify:** `internal/graph/builder.go`

**Add new type after Edge struct (around line 27):**
```go
// edgeKey is used for O(1) duplicate detection
type edgeKey struct {
    From string
    To   string
    Type EdgeType // or string if Task 02 not done yet
}
```

**Modify Builder struct (lines 28-32):**
```go
// Builder builds a dependency graph from multiple analysis results
type Builder struct {
    nodes   map[string]Node
    edges   []Edge
    edgeSet map[edgeKey]struct{} // For O(1) duplicate detection
}
```

**Modify NewBuilder function (lines 35-40):**
```go
// NewBuilder creates a new graph builder
func NewBuilder() *Builder {
    return &Builder{
        nodes:   make(map[string]Node),
        edges:   []Edge{},
        edgeSet: make(map[edgeKey]struct{}),
    }
}
```

**Replace addEdge function (lines 153-160):**

Current:
```go
func (b *Builder) addEdge(edge Edge) {
    // Check for duplicate edges
    for _, e := range b.edges {
        if e.From == edge.From && e.To == edge.To && e.Type == edge.Type {
            return
        }
    }
    b.edges = append(b.edges, edge)
}
```

New:
```go
func (b *Builder) addEdge(edge Edge) {
    key := edgeKey{From: edge.From, To: edge.To, Type: edge.Type}
    if _, exists := b.edgeSet[key]; exists {
        return
    }
    b.edgeSet[key] = struct{}{}
    b.edges = append(b.edges, edge)
}
```

## Verification

```bash
# Build should succeed
go build ./...

# Create a benchmark test
cat > internal/graph/builder_bench_test.go << 'EOF'
package graph

import "testing"

func BenchmarkAddEdge(b *testing.B) {
    builder := NewBuilder()
    
    // Pre-populate with some edges
    for i := 0; i < 100; i++ {
        builder.addEdge(Edge{
            From: "app-" + string(rune('a'+i%26)),
            To:   "app-" + string(rune('a'+(i+1)%26)),
            Type: "call",
        })
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Try to add duplicate (worst case for old implementation)
        builder.addEdge(Edge{From: "app-a", To: "app-b", Type: "call"})
    }
}

func TestAddEdge_NoDuplicates(t *testing.T) {
    builder := NewBuilder()
    
    edge := Edge{From: "a", To: "b", Type: "call"}
    builder.addEdge(edge)
    builder.addEdge(edge) // duplicate
    builder.addEdge(edge) // duplicate
    
    g := builder.Build()
    if len(g.Edges) != 1 {
        t.Errorf("expected 1 edge, got %d", len(g.Edges))
    }
}
EOF

# Run tests
go test ./internal/graph/... -v

# Run benchmark
go test ./internal/graph/... -bench=. -benchmem
```

## Files Changed
- `internal/graph/builder.go` (modify ~15 lines)
- `internal/graph/builder_bench_test.go` (new file, ~40 lines)
