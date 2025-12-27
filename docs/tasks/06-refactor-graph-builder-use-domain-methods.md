# Task 06: Refactor Graph Builder to Use NaisConfig Methods

## Deliverable
Replace inline business logic in `graph/builder.go` with calls to `NaisConfig` methods created in Task 01.

## Context
The `AddFromLocalAnalysis` function in `graph/builder.go` reaches deep into `NaisConfig` internals to determine if an app is stateful. This violates encapsulation (Feature Envy anti-pattern).

**Current code (builder.go:50-56):**
```go
stateful := false
if analysis.NaisConfig != nil && analysis.NaisConfig.Spec.GCP != nil {
    if len(analysis.NaisConfig.Spec.GCP.SQLInstances) > 0 || len(analysis.NaisConfig.Spec.GCP.Buckets) > 0 {
        stateful = true
    }
}
```

After Task 01, `NaisConfig` has an `IsStateful()` method that encapsulates this logic.

## Key Decisions and Principles
- Use domain methods instead of reaching into object internals
- Keep graph builder focused on graph construction, not config interpretation
- One line should replace 6 lines of logic

## Delivers
Cleaner graph builder code that delegates to domain methods.

## Acceptance Criteria
- Stateful check replaced with `analysis.NaisConfig.IsStateful()`
- Code is more readable and shorter
- Behavior is identical
- Tests pass

## Dependencies
**Requires Task 01 completed first** (NaisConfig behavior methods)

## Related Code

**File to modify:** `internal/graph/builder.go`

**Current code (lines 42-65):**
```go
func (b *Builder) AddFromLocalAnalysis(analysis *analyzer.LocalAnalysis) {
    if analysis == nil || analysis.Repository == nil {
        return
    }

    appName := analysis.Repository.Name

    // Determine if app is stateful (has database or storage)
    stateful := false
    if analysis.NaisConfig != nil && analysis.NaisConfig.Spec.GCP != nil {
        if len(analysis.NaisConfig.Spec.GCP.SQLInstances) > 0 || len(analysis.NaisConfig.Spec.GCP.Buckets) > 0 {
            stateful = true
        }
    }

    // Add app node
    b.addNode(Node{
        ID:       appName,
        Type:     "app",
        Stateful: stateful,
    })
    // ... rest of function
```

**Replace with:**
```go
func (b *Builder) AddFromLocalAnalysis(analysis *analyzer.LocalAnalysis) {
    if analysis == nil || analysis.Repository == nil {
        return
    }

    appName := analysis.Repository.Name

    // Add app node - NaisConfig.IsStateful() handles nil checks internally
    b.addNode(Node{
        ID:       appName,
        Type:     "app",
        Stateful: analysis.NaisConfig.IsStateful(),
    })
    // ... rest of function
```

This change:
- Removes 6 lines of inline logic
- Replaces with 1 method call
- `IsStateful()` handles nil checks internally (from Task 01)

## Verification

```bash
# Build should succeed
go build ./...

# Run graph command on test repos to verify behavior unchanged
./naisalyser graph ./repos --output /tmp/graph-test -v

# Check output is valid
cat /tmp/graph-test/dependencies.json | head -20

# Clean up
rm -rf /tmp/graph-test
```

## Files Changed
- `internal/graph/builder.go` (modify ~6 lines to ~2 lines)
