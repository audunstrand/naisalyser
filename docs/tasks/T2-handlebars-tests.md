# Task T2: Handlebars Preprocessing Tests

## Deliverable
Implement comprehensive test suite (15+ tests) for the handlebars preprocessing functions extracted in Task 12.

## Context
Following Task 12's refactoring of `preprocessHandlebars` into smaller, testable functions, this task creates a comprehensive test suite to validate the behavior of:
- `removeEachBlocks(content string) string`
- `replaceTemplateVars(content string) string`
- `replaceLineTemplateVars(line string) string`
- `getPlaceholderForContext(line string) string`

These functions handle preprocessing of Handlebars template syntax in NAIS YAML files, replacing template variables with appropriate placeholder values so the YAML can be parsed.

## Key Testing Requirements

### Context-Aware Replacement
The preprocessing uses context-aware logic to determine appropriate placeholder values:
- **Numeric contexts** (use `"1"`): replicas, min, max, port, timeout, delay
- **Non-numeric contexts** (use `"placeholder"`): name, image, namespace, url, etc.

### Edge Cases to Test
- Multiple variables on a single line
- Nested `{{#each}}` blocks
- Malformed templates (missing closing tags)
- Empty strings and content with no templates
- Case-insensitive context matching
- Real NAIS YAML files with mixed contexts

## Test Coverage

### Tests for `removeEachBlocks` (5 tests)
1. **Single block removal** - Basic `{{#each}}...{{/each}}` removal
2. **Multiple blocks** - Remove multiple non-nested blocks
3. **Nested blocks** - Handle blocks within blocks
4. **No blocks** - Pass through content without `{{#each}}`
5. **Malformed missing close** - Handle missing `{{/each}}` tag

### Tests for `replaceLineTemplateVars` (13 tests)
1. **Single variable** - Replace one `{{var}}`
2. **Multiple variables** - Replace multiple `{{var}}` on same line
3. **No variables** - Pass through plain text
4. **Malformed missing close** - Handle missing `}}`
5. **Replicas context** (4 subtests) - Test replica/min/max contexts
6. **Port context** (3 subtests) - Test port/targetPort contexts
7. **Timeout context** (3 subtests) - Test timeout/delay contexts
8. **Non-numeric context** (4 subtests) - Test name/image/namespace/url
9. **Mixed context** - Multiple variables where line contains numeric context

### Tests for `replaceTemplateVars` (3 tests)
1. **Multiple lines** - Process multi-line content
2. **Empty content** - Handle empty string
3. **No templates** - Pass through content without templates

### Tests for `preprocessHandlebars` (6 tests)
1. **Simple replacement** - Basic variable replacement
2. **Numeric context** - Verify numeric context gets "1"
3. **Each block** - Integration test with block removal
4. **Multiple vars** - Multiple variables on same line
5. **Real NAIS YAML** - Complete NAIS config with mixed contexts
6. **Real NAIS YAML with each block** - Complete config with `{{#each}}`

### Tests for `getPlaceholderForContext` (2 tests)
1. **Numeric contexts** (6 subtests) - Test all numeric keywords
2. **Non-numeric contexts** (4 subtests) - Test non-numeric scenarios

## Total Test Count
- **25 top-level test functions**
- **49 total test cases** (including subtests)

This exceeds the requirement of 15+ tests and provides comprehensive coverage of:
- All exported and internal functions
- Context-aware replacement logic
- Edge cases and error conditions
- Integration scenarios with real NAIS YAML

## Implementation

**File:** `internal/analyzer/handlebars_test.go`

The test file is organized into sections:
1. Tests for `removeEachBlocks`
2. Tests for `replaceLineTemplateVars` and context-aware replacement
3. Tests for `replaceTemplateVars`
4. Integration tests for `preprocessHandlebars`
5. Tests for `getPlaceholderForContext`

Each test follows Go testing conventions:
- Clear test names describing what is tested
- Table-driven tests for similar scenarios
- Subtests for grouped test cases
- Descriptive error messages

## Verification

```bash
# Run all tests
go test ./internal/analyzer/... -v

# Run specific test
go test ./internal/analyzer/... -v -run TestPreprocessHandlebars_RealNaisYAML

# Check test coverage
go test ./internal/analyzer/... -cover
```

## Related Tasks
- **Task 12** (`12-flatten-handlebars-preprocessing.md`) - Refactored the functions being tested

## Files Changed
- `internal/analyzer/handlebars_test.go` (new file, ~450 lines)
- `docs/tasks/T2-handlebars-tests.md` (this file, new)

## Success Criteria
- ✅ 15+ test functions implemented (25 total)
- ✅ Tests for all handlebars preprocessing functions
- ✅ Context-aware replacement tested (replicas/ports/timeouts use '1')
- ✅ Multiple variables per line tested
- ✅ Edge cases covered (malformed templates, empty strings)
- ✅ Real NAIS YAML preprocessing scenarios tested
- ✅ All tests pass
