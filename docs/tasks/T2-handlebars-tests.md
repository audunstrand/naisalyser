# Test Task T2: Handlebars Preprocessing Tests

## Description
Implement comprehensive tests for Handlebars template preprocessing:
- `preprocessHandlebars()`
- `removeEachBlocks()`
- `replaceTemplateVars()`
- `replaceLineTemplateVars()`

## File Location
`internal/analyzer/handlebars_test.go`

## Test Cases Required

### RemoveEachBlocks Tests (5 tests)
1. `TestRemoveEachBlocks_SingleBlock` - Removes one {{#each}} block
2. `TestRemoveEachBlocks_MultipleBlocks` - Removes multiple blocks
3. `TestRemoveEachBlocks_NoBlocks` - Content unchanged when no blocks
4. `TestRemoveEachBlocks_BlockAtStart` - Block at beginning of content
5. `TestRemoveEachBlocks_BlockAtEnd` - Block at end of content
6. `TestRemoveEachBlocks_NestedBlocks` - Handles sequential nested blocks
7. `TestRemoveEachBlocks_PreservesOtherContent` - Keeps other lines intact

### ReplaceTemplateVars Tests (4 tests)
1. `TestReplaceTemplateVars_SingleVarPerLine` - One variable per line
2. `TestReplaceTemplateVars_MultipleVarsPerLine` - Multiple variables on same line
3. `TestReplaceTemplateVars_NoVars` - Content unchanged when no variables
4. `TestReplaceTemplateVars_MultipleLines` - Processes all lines correctly
5. `TestReplaceTemplateVars_PreserveNewlines` - Line structure maintained

### ReplaceLineTemplateVars Tests (6 tests)
1. `TestReplaceLineTemplateVars_ReplicaContext` - Uses "1" for replica
2. `TestReplaceLineTemplateVars_MinMaxContext` - Uses "1" for min/max
3. `TestReplaceLineTemplateVars_PortContext` - Uses "1" for port
4. `TestReplaceLineTemplateVars_TimeoutContext` - Uses "1" for timeout
5. `TestReplaceLineTemplateVars_DelayContext` - Uses "1" for delay
6. `TestReplaceLineTemplateVars_DefaultContext` - Uses "placeholder" for others
7. `TestReplaceLineTemplateVars_MultipleVars` - Handles multiple on same line
8. `TestReplaceLineTemplateVars_CaseInsensitive` - Detects context case-insensitively
9. `TestReplaceLineTemplateVars_MixedContexts` - Different contexts on same line
10. `TestReplaceLineTemplateVars_NoVars` - Returns unchanged line

### PreprocessHandlebars Integration Tests (3 tests)
1. `TestPreprocessHandlebars_RealNaisYaml` - Real NAIS config with handlebars
2. `TestPreprocessHandlebars_ComplexTemplate` - Multiple types of templates
3. `TestPreprocessHandlebars_EmptyContent` - Empty string handling

## Test Data Fixtures
Create test fixtures in `internal/analyzer/testdata/`:
- `nais_with_handlebars.yaml` - Real NAIS config with templates
- `each_blocks.yaml` - Multiple {{#each}} blocks
- `template_vars.yaml` - Various template variables

## Test Data Examples

```yaml
# Test: ReplaceLineTemplateVars with replica context
Input: "  replicas: {{replicas}}"
Expected: "  replicas: 1"

# Test: ReplaceLineTemplateVars with default context  
Input: "  image: {{image}}"
Expected: "  image: placeholder"

# Test: RemoveEachBlocks
Input: |
  kind: Application
  {{#each items}}
    name: {{name}}
  {{/each}}
  spec: {}

Expected: |
  kind: Application
  spec: {}
```

## Implementation Notes
- Use subtests for clarity (t.Run)
- Test bidirectional data flow (preprocessor -> YAML parser)
- Verify YAML parses successfully after preprocessing
- Include multiline block examples
- Test case sensitivity carefully

## Success Criteria
- ✅ All 15+ tests pass
- ✅ Code coverage for handlebars functions ≥ 95%
- ✅ Real NAIS YAML files preprocess correctly
- ✅ No regressions in existing functionality
