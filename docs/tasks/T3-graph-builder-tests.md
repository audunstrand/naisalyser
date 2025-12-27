# Test Task T3: Graph Builder Tests

## Description
Implement comprehensive tests for the dependency graph builder:
- `NewBuilder()`
- `AddFromLocalAnalysis()`
- `addNode()`
- `addEdge()`
- `Build()`
- Edge deduplication with type-safe constants

## File Location
`internal/graph/builder_test.go`

## Test Cases Required

### NewBuilder Tests (3 tests)
1. `TestNewBuilder_InitializesEmpty` - Builder starts with no nodes/edges
2. `TestNewBuilder_InitializesNodeMap` - nodes map is empty
3. `TestNewBuilder_InitializesEdgeSet` - edgeSet map is initialized
4. `TestNewBuilder_InitializesEdgeSlice` - edges slice is empty

### addNode Tests (5 tests)
1. `TestAddNode_AddsNewNode` - Single node is added
2. `TestAddNode_DoesNotDuplicate` - Duplicate nodes with same ID rejected
3. `TestAddNode_PreservesType` - Node type (app, kafka, external) preserved
4. `TestAddNode_PreservesStateful` - Stateful flag correctly stored
5. `TestAddNode_PreservesNamespace` - Namespace field preserved
6. `TestAddNode_MultipleNodes` - Multiple different nodes all added

### addEdge Tests (8 tests)
1. `TestAddEdge_AddsNewEdge` - Single edge is added
2. `TestAddEdge_DeduplicatesExactMatch` - Duplicate edges rejected
3. `TestAddEdge_DistinguishesByType` - Same nodes, different edge types allowed
4. `TestAddEdge_DistinguishesByDirection` - A→B differs from B→A
5. `TestAddEdge_PreservesType` - Edge type (call, kafka, external) preserved
6. `TestAddEdge_O1Performance` - Deduplication uses map (O(1) not O(n))
7. `TestAddEdge_MaintainsOrder` - Edges in slice maintain insertion order
8. `TestAddEdge_MultipleEdges` - Multiple different edges all added
9. `TestAddEdge_TypeSafeConstants` - Uses EdgeTypeCall, EdgeTypeKafka, EdgeTypeExternal

### AddFromLocalAnalysis Tests (15 tests)
1. `TestAddFromLocalAnalysis_NilAnalysis` - Handles nil analysis gracefully
2. `TestAddFromLocalAnalysis_NilRepository` - Handles nil repository gracefully
3. `TestAddFromLocalAnalysis_AddsAppNode` - Adds main application node
4. `TestAddFromLocalAnalysis_AppNodeType` - App node has NodeTypeApp
5. `TestAddFromLocalAnalysis_AppNodeName` - App node ID is repo.Name
6. `TestAddFromLocalAnalysis_StatefulFlag` - Uses NaisConfig.IsStateful()
7. `TestAddFromLocalAnalysis_InboundRules` - Adds inbound access policy rules
8. `TestAddFromLocalAnalysis_InboundEdges` - Creates edges from inbound apps
9. `TestAddFromLocalAnalysis_OutboundRules` - Adds outbound access policy rules
10. `TestAddFromLocalAnalysis_OutboundEdges` - Creates edges to outbound apps
11. `TestAddFromLocalAnalysis_ExternalHosts` - Adds external host nodes
12. `TestAddFromLocalAnalysis_ExternalEdges` - Creates edges to external hosts
13. `TestAddFromLocalAnalysis_KafkaTopics` - Adds Kafka topic nodes
14. `TestAddFromLocalAnalysis_KafkaEdges` - Creates edges to Kafka topics
15. `TestAddFromLocalAnalysis_KafkaPool` - Handles Kafka pool fallback
16. `TestAddFromLocalAnalysis_NilNaisConfig` - Handles missing NaisConfig

### Build Tests (3 tests)
1. `TestBuild_ReturnsGraph` - Build() returns *Graph
2. `TestBuild_IncludesAllNodes` - All added nodes in result
3. `TestBuild_IncludesAllEdges` - All deduplicated edges in result

### Integration Tests (3 tests)
1. `TestBuilder_FullWorkflow` - Create → Add → Build workflow
2. `TestBuilder_TypeSafeConstants` - Uses NodeType and EdgeType constants
3. `TestBuilder_DuplicateEdgeDetection` - O(1) deduplication proves faster

## Test Data Fixtures

Create mock analysis data:
```go
// Mock NaisConfig with IsStateful() = true
mockConfig := &NaisConfig{
    Spec: NaisSpec{
        GCP: &NaisGCP{
            SQLInstances: []NaisSQLInstance{{Type: "POSTGRES_14"}},
        },
    },
}

// Mock LocalAnalysis
mockAnalysis := &LocalAnalysis{
    Repository: &local.Repository{Name: "my-app", FullName: "navikt/my-app"},
    NaisConfig: mockConfig,
    KafkaTopics: []KafkaTopic{{Name: "events", ConstName: "EVENTS_TOPIC"}},
}
```

## Edge Cases to Test
- Multiple calls to AddFromLocalAnalysis (accumulates)
- Duplicate nodes from different analyses
- Duplicate edges from different analyses
- Empty access policies
- Empty Kafka topics
- No external hosts
- Namespace metadata on nodes

## Implementation Notes
- Use subtests for organization
- Test the `edgeKey` struct indirectly through deduplication
- Verify type-safe constants are used (not string literals)
- Performance: verify O(1) deduplication (can't really benchmark, but verify map is used)
- No external dependencies needed (test fixtures only)

## Success Criteria
- ✅ All 30+ tests pass
- ✅ Code coverage for builder functions ≥ 90%
- ✅ Type-safe constants verified in output
- ✅ Edge deduplication confirmed working
- ✅ No string literals (all using constants)
