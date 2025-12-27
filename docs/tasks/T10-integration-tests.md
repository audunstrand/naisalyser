# Test Task T10: Integration Tests

## Description
Implement end-to-end integration tests that verify complete workflows:
- Full local analysis pipeline
- Full graph building pipeline  
- Parser integration with real files
- Command-to-output complete flows

## File Location
`internal/integration_test.go` (or separate package)

## Test Cases Required

### Local Analysis Workflow Tests (6 tests)
1. `TestIntegration_LocalAnalysisWorkflow` - Complete local analysis
   - Create temp repo
   - Add NAIS config
   - Analyze with local reader
   - Verify output structure
   
2. `TestIntegration_LocalAnalysisWithDependencies` - Analysis with deps
   - Create temp repo with build.gradle.kts
   - Analyze with dependency parsing
   - Verify deps are extracted
   
3. `TestIntegration_LocalAnalysisWithKafka` - Analysis with Kafka
   - Create temp repo with Kafka config
   - Analyze and extract topics
   - Verify topics identified
   
4. `TestIntegration_LocalAnalysisWithAccessPolicy` - Access policy extraction
   - NAIS config with inbound/outbound rules
   - Verify rules are extracted
   - Verify nodes and edges created
   
5. `TestIntegration_LocalAnalysisWithStorage` - Database and bucket detection
   - NAIS config with GCP resources
   - Verify stateful flag set correctly
   - Verify resources listed in output
   
6. `TestIntegration_LocalAnalysisWithHandlebars` - Template preprocessing
   - NAIS config with Handlebars templates
   - Preprocess and parse
   - Verify templates replaced correctly

### Graph Building Workflow Tests (5 tests)
1. `TestIntegration_GraphBuildingMultipleRepos` - Multi-repo graph
   - Create 3 test repos with dependencies
   - Build graph from all
   - Verify all nodes and edges present
   
2. `TestIntegration_GraphBuildingEdgeDeduplication` - Edge dedup works end-to-end
   - Create repos with duplicate edges
   - Build graph
   - Verify duplicates removed
   
3. `TestIntegration_GraphBuildingMermaidOutput` - Complete Mermaid generation
   - Build graph from repos
   - Generate Mermaid diagram
   - Verify valid output
   
4. `TestIntegration_GraphBuildingJsonOutput` - JSON serialization works
   - Build graph
   - Serialize to JSON
   - Verify valid JSON structure
   
5. `TestIntegration_GraphBuildingComplexDependencies` - Complex graph structure
   - Multi-type nodes (app, kafka, external)
   - Multi-type edges (call, kafka, external)
   - Verify graph represents all correctly

### Parser Integration Tests (6 tests)
1. `TestIntegration_GradleParserWithRealFile` - Real build.gradle.kts
   - Create valid gradle file
   - Parse dependencies
   - Verify extraction works
   
2. `TestIntegration_MavenParserWithRealFile` - Real pom.xml
   - Create valid pom.xml
   - Parse dependencies
   - Verify extraction and filtering
   
3. `TestIntegration_NpmParserWithRealFile` - Real package.json
   - Create valid package.json
   - Parse dependencies and dev deps
   - Verify both sections parsed
   
4. `TestIntegration_GoModParserWithRealFile` - Real go.mod
   - Create valid go.mod
   - Parse dependencies
   - Verify extraction works
   
5. `TestIntegration_MultiParserProject` - Project with multiple build systems
   - Create repo with gradle + npm
   - Parse both
   - Verify all deps extracted
   
6. `TestIntegration_ParserErrorHandling` - Malformed files handled
   - Malformed gradle, pom, package.json
   - Verify parsers handle gracefully
   - Verify no crashes

### Language Detection Workflow Tests (4 tests)
1. `TestIntegration_LanguageDetectionGo` - Go project detection
   - Create Go project structure
   - Read repo metadata
   - Verify language detected
   
2. `TestIntegration_LanguageDetectionKotlin` - Kotlin project detection
   - Create Kotlin/Spring project
   - Verify language detected
   
3. `TestIntegration_LanguageDetectionNode` - Node.js project detection
   - Create package.json
   - Verify language detected
   
4. `TestIntegration_LanguageDetectionMultiple` - Mixed language handling
   - Create repo with multiple languages
   - Verify first one detected

### Command-to-Output Workflows (4 tests)
1. `TestIntegration_LocalCommandFullWorkflow` - CLI command end-to-end
   - Run `naisalyser local` on temp repo
   - Verify output generated
   - Verify markdown is valid
   
2. `TestIntegration_LocalBatchCommandFullWorkflow` - Batch command
   - Create multiple temp repos
   - Run `naisalyser local-batch`
   - Verify all processed
   
3. `TestIntegration_GraphCommandFullWorkflow` - Graph command
   - Create repos
   - Run `naisalyser graph`
   - Verify JSON and Mermaid generated
   
4. `TestIntegration_OutputMarkdownValid` - Output markdown is valid
   - Generate full documentation
   - Verify markdown syntax
   - Verify can be rendered

### Data Transformation Workflows (3 tests)
1. `TestIntegration_NaisConfigToViewModel` - Config → view model
   - Create complex NAIS config
   - Convert to view model
   - Verify all fields present
   
2. `TestIntegration_AnalysisToMarkdown` - Analysis → markdown
   - Complete analysis
   - Generate markdown
   - Verify complete output
   
3. `TestIntegration_RepositoryMetadataExtraction` - Metadata pipeline
   - Create repo with README
   - Extract metadata
   - Verify completeness

## Test Repositories Setup

### Go Project Structure
```
go-test-repo/
├── go.mod (with dependencies)
├── README.md (description)
├── main.go
└── handler/
    └── handler.go
```

### Kotlin Project Structure
```
kotlin-test-repo/
├── build.gradle.kts (with dependencies)
├── nais.yaml (with config)
├── README.md
├── src/main/kotlin/
│   └── App.kt
└── src/main/kotlin/kafka/
    └── Topics.kt
```

### Node Project Structure
```
node-test-repo/
├── package.json (with deps)
├── tsconfig.json
├── README.md
└── src/
    └── index.ts
```

### Complex Project with All Features
```
complex-test-repo/
├── nais.yaml (with access policy, kafka, storage)
├── build.gradle.kts (gradle deps)
├── pom.xml (maven deps)
├── package.json (npm deps)
├── go.mod (go deps)
├── README.md
└── src/
```

## Data Validation

Validate complete data flows:
```go
func validateAnalysisComplete(a *analyzer.Analysis) error {
    if a.Repository == nil {
        return fmt.Errorf("repository nil")
    }
    if a.Repository.Name == "" {
        return fmt.Errorf("repo name empty")
    }
    // ... more validations
    return nil
}
```

## Performance Baselines

Test that operations complete in reasonable time:
```go
func TestIntegration_PerformanceBaseline(t *testing.T) {
    start := time.Now()
    // ... full workflow ...
    duration := time.Since(start)
    
    if duration > 30*time.Second {
        t.Logf("WARNING: Workflow took %v (baseline: <30s)", duration)
    }
}
```

## Error Recovery Workflows

Test that systems handle errors gracefully:
1. Missing NAIS config → analysis completes without it
2. Invalid markdown in README → parsed without error
3. Malformed YAML → preprocessing handles it
4. Parser errors → continue with other parsers
5. Missing dependencies file → continue without deps

## Edge Cases to Test
- Very large repositories (1000+ files)
- Deep directory structures (10+ levels)
- Unicode in all fields
- Special characters in names
- Circular dependencies
- Duplicate dependency definitions
- Mixed valid and invalid files

## Implementation Notes
- Create minimal but complete test repositories
- Use `t.TempDir()` for isolation
- Test both success and error paths
- Verify no data loss in transformations
- Check complete output structure
- Validate against real NAIS schema
- Test with realistic data volumes

## Success Criteria
- ✅ All 18+ integration tests pass
- ✅ Code coverage increases to 85%+
- ✅ Real-world workflows validated
- ✅ No data loss in transformations
- ✅ Error handling verified
- ✅ Performance acceptable
- ✅ Output quality verified
