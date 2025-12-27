# Test Task T4: Command Handler Tests

## Description
Implement comprehensive tests for CLI command handlers:
- `runAnalyze()` (analyze.go)
- `runBatch()` (batch.go)
- `runLocal()` (local.go)
- `runLocalBatch()` (local.go)
- `runGraph()` (graph.go)
- Flag helper functions: `mustGetString()`, `mustGetBool()`

## File Location
`cmd/commands_test.go` and `cmd/flags_test.go`

## Test Cases Required

### Flag Helper Tests (4 tests)
1. `TestMustGetString_ValidFlag` - Returns string value correctly
2. `TestMustGetString_MissingFlag` - Returns error for invalid flag
3. `TestMustGetBool_ValidFlag` - Returns bool value correctly
4. `TestMustGetBool_MissingFlag` - Returns error for invalid flag

### RunAnalyze Tests (8 tests)
1. `TestRunAnalyze_RequiresRepoArg` - Fails without repository argument
2. `TestRunAnalyze_FlagParsing` - Gets output and verbose flags correctly
3. `TestRunAnalyze_DefaultOutputDir` - Uses "./docs/generated" by default
4. `TestRunAnalyze_CustomOutputDir` - Respects --output flag
5. `TestRunAnalyze_VerboseFlag` - Respects --verbose flag
6. `TestRunAnalyze_ErrorHandling` - Handles repo not found gracefully
7. `TestRunAnalyze_OutputDirectoryCreated` - Creates output directory if missing
8. `TestRunAnalyze_SuccessPath` - Completes successfully with valid input

### RunBatch Tests (10 tests)
1. `TestRunBatch_RequiresTopicOrConfig` - Fails if both missing
2. `TestRunBatch_TopicFlag` - Fetches repos by topic
3. `TestRunBatch_ConfigFlag` - Reads repos from config file
4. `TestRunBatch_ConfigFileNotFound` - Error handling for missing config
5. `TestRunBatch_ConfigFileMalformed` - YAML parsing error handling
6. `TestRunBatch_BothFlagsFails` - Error when both --topic and --config specified
7. `TestRunBatch_OrgFlag` - Respects --org flag
8. `TestRunBatch_OutputDirFlag` - Respects --output flag
9. `TestRunBatch_MultipleRepos` - Iterates through all repos
10. `TestRunBatch_SkipsFailedRepos` - Continues on individual repo errors

### RunLocal Tests (10 tests)
1. `TestRunLocal_RequiresPathArg` - Fails without path argument
2. `TestRunLocal_InvalidPath` - Error for non-existent path
3. `TestRunLocal_NotDirectory` - Error for file path
4. `TestRunLocal_OrgFlagDefault` - Defaults to "navikt"
5. `TestRunLocal_OrgFlagCustom` - Accepts custom org
6. `TestRunLocal_OutputDirFlag` - Respects --output flag
7. `TestRunLocal_VerboseFlag` - Respects --verbose flag
8. `TestRunLocal_OutputDirCreated` - Creates output directory
9. `TestRunLocal_SuccessPath` - Completes successfully
10. `TestRunLocal_MarkdownGenerated` - Creates output markdown file

### RunLocalBatch Tests (10 tests)
1. `TestRunLocalBatch_RequiresDirArg` - Fails without directory argument
2. `TestRunLocalBatch_InvalidDir` - Error for non-existent directory
3. `TestRunLocalBatch_EmptyDir` - Handles empty repo directory
4. `TestRunLocalBatch_SingleRepo` - Processes single repo
5. `TestRunLocalBatch_MultipleRepos` - Iterates through all repos
6. `TestRunLocalBatch_SkipsHiddenDirs` - Ignores .hidden directories
7. `TestRunLocalBatch_OrgFlagApplied` - Org flag passed to all repos
8. `TestRunLocalBatch_OutputDirStructure` - Creates per-repo output dirs
9. `TestRunLocalBatch_ContinuesOnError` - Continues on repo errors
10. `TestRunLocalBatch_AllReposProcessed` - Reports correct count

### RunGraph Tests (12 tests)
1. `TestRunGraph_RequiresDirArg` - Fails without directory argument
2. `TestRunGraph_InvalidDir` - Error for non-existent directory
3. `TestRunGraph_EmptyDir` - Handles empty repo directory
4. `TestRunGraph_OutputDirCreated` - Creates output directory
5. `TestRunGraph_JsonGenerated` - Creates dependencies.json
6. `TestRunGraph_MermaidGenerated` - Creates dependencies.md
7. `TestRunGraph_OutputDirFlag` - Respects --output flag
8. `TestRunGraph_VerboseFlag` - Respects --verbose flag
9. `TestRunGraph_ValidJsonOutput` - Generated JSON is valid
10. `TestRunGraph_ValidMermaidOutput` - Generated Markdown has mermaid block
11. `TestRunGraph_MultipleRepos` - Builds graph from multiple repos
12. `TestRunGraph_ErrorRecovery` - Continues on individual repo errors

## Test Infrastructure

### Mock Filesystem
Use `t.TempDir()` for temporary test directories:
```go
func TestRunLocal_SuccessPath(t *testing.T) {
    tmpDir := t.TempDir()
    // Create test repo structure
    // Run command with temp dir
}
```

### Mock Config Files
```go
type MockBatchConfig struct {
    Repositories []string
}
```

### Cobra Command Mocking
```go
func newTestCmd(flags map[string]interface{}) *cobra.Command {
    cmd := &cobra.Command{RunE: runAnalyze}
    for name, val := range flags {
        switch v := val.(type) {
        case string:
            cmd.Flags().String(name, v, "")
        case bool:
            cmd.Flags().Bool(name, v, "")
        }
    }
    return cmd
}
```

## Test Repository Structure
Create test repo structure in temp dirs:
```
test-repo/
├── README.md
├── nais.yaml (with/without handlebars)
├── build.gradle.kts
├── pom.xml
├── package.json
└── go.mod
```

## Edge Cases to Test
- Unicode in directory names
- Very long paths
- Permission errors
- Disk space errors
- Concurrent repo processing
- Large number of repos

## Implementation Notes
- Use table-driven tests for flag combinations
- Mock external systems (GitHub client)
- Test error messages are helpful
- Verify exit codes
- Test concurrent safety if applicable

## Success Criteria
- ✅ All 50+ tests pass
- ✅ Code coverage for cmd package ≥ 80%
- ✅ No external API calls (all mocked)
- ✅ Proper error handling verified
- ✅ All flag combinations tested
