# Test Task T6: GitHub Client Tests

## Description
Implement comprehensive tests for GitHub API interaction:
- `NewClient()`
- `GetRepository()`
- `GetRepositoriesByTopic()`
- `GetFileContent()`
- `FileExists()`
- `GetDirectoryContents()`

## File Location
`internal/github/client_test.go`

## Test Cases Required

### NewClient Tests (2 tests)
1. `TestNewClient_VerboseFalse` - Creates client with verbose=false
2. `TestNewClient_VerboseTrue` - Creates client with verbose=true

### GetRepository Tests (5 tests)
1. `TestGetRepository_ValidRepo` - Fetches valid repository
2. `TestGetRepository_InvalidRepo` - Error for non-existent repo
3. `TestGetRepository_ParsesMetadata` - Correctly parses name, description, language
4. `TestGetRepository_FullNameFormat` - Returns owner/repo format
5. `TestGetRepository_ErrorHandling` - Handles GitHub API errors gracefully

### GetRepositoriesByTopic Tests (6 tests)
1. `TestGetRepositoriesByTopic_ValidTopic` - Returns repos by topic
2. `TestGetRepositoriesByTopic_EmptyResult` - Handles no repos for topic
3. `TestGetRepositoriesByTopic_InvalidTopic` - Handles non-existent topic
4. `TestGetRepositoriesByTopic_Pagination` - Handles paginated results
5. `TestGetRepositoriesByTopic_Organization` - Filters by organization
6. `TestGetRepositoriesByTopic_ErrorHandling` - Handles API errors

### GetFileContent Tests (6 tests)
1. `TestGetFileContent_ValidFile` - Reads file content
2. `TestGetFileContent_NonExistentFile` - Error for missing file
3. `TestGetFileContent_Base64Decoding` - Correctly decodes base64
4. `TestGetFileContent_BinaryContent` - Handles binary files
5. `TestGetFileContent_LargeFile` - Handles large files
6. `TestGetFileContent_ErrorHandling` - Handles API errors

### FileExists Tests (4 tests)
1. `TestFileExists_ExistingFile` - Returns true for existing file
2. `TestFileExists_NonExistentFile` - Returns false for missing file
3. `TestFileExists_Directory` - Returns false for directory path
4. `TestFileExists_ErrorHandling` - Returns false on API error

### GetDirectoryContents Tests (5 tests)
1. `TestGetDirectoryContents_ValidDirectory` - Lists directory entries
2. `TestGetDirectoryContents_EmptyDirectory` - Handles empty directory
3. `TestGetDirectoryContents_NonExistentDirectory` - Error for missing dir
4. `TestGetDirectoryContents_Pagination` - Handles paginated results
5. `TestGetDirectoryContents_ErrorHandling` - Handles API errors

### ExecGH Tests (4 tests)
1. `TestExecGH_ValidCommand` - Executes valid gh command
2. `TestExecGH_InvalidCommand` - Error for bad command
3. `TestExecGH_JsonOutput` - Parses JSON output
4. `TestExecGH_ErrorOutput` - Handles error output

## Mocking Strategy

### Environment Variable Mocking
Mock `gh` CLI output instead of actual API calls:
```bash
export GH_TOKEN=fake-token
```

### Command Output Mocking
Create helper functions to mock `exec.Command`:
```go
func mockExecCommand(command string, args ...string) *exec.Cmd {
    cs := []string{"-test.run=TestHelperProcess", "--"}
    cs = append(cs, command)
    cs = append(cs, args...)
    cmd := exec.Command(os.Args[0], cs...)
    cmd.Env = append(os.Environ(), "GO_TEST_PROCESS=1")
    return cmd
}
```

### Mock Repository Data
```go
var mockRepos = map[string]interface{}{
    "navikt/aap-app": map[string]interface{}{
        "name":        "aap-app",
        "description": "AAP Application",
        "language":    "kotlin",
    },
}

var mockTopic = []string{
    "navikt/aap-app",
    "navikt/aap-api",
    "navikt/aap-scheduler",
}
```

### Mock File Content
```go
var mockFiles = map[string]string{
    "navikt/aap-app:nais.yaml": `
apiVersion: nais.io/v1alpha1
kind: Application
metadata:
  name: aap-app
spec:
  image: navikt/aap-app:latest
`,
    "navikt/aap-app:build.gradle.kts": `
dependencies {
    implementation("org.springframework.boot:spring-boot-starter-web:3.0.0")
}
`,
}
```

## Test Data Fixtures
Create fixtures in `internal/github/testdata/`:
- `repository.json` - Sample GitHub repository response
- `repositories_list.json` - Sample list response
- `file_content.json` - Sample file content response
- `directory_contents.json` - Sample directory listing response

## Error Scenarios to Test
- 404 Not Found
- 403 Forbidden (permission denied)
- 401 Unauthorized (missing token)
- 500 Internal Server Error
- Network timeout
- Rate limit exceeded
- Invalid JSON response
- Truncated response

## Edge Cases to Test
- Very long file paths
- Files with special characters
- Very large repositories (1000+ files)
- Repositories with no description
- Repositories with no language
- Organizations with many repos
- Topics with no repositories

## Implementation Notes
- Do NOT make actual API calls
- Use environment variable substitution
- Mock all external command execution
- Test error handling thoroughly
- Verify JSON parsing works correctly
- Test pagination logic
- Verify organization filtering works

## Success Criteria
- ✅ All 20+ tests pass
- ✅ Code coverage for client functions ≥ 85%
- ✅ No actual GitHub API calls made
- ✅ All error cases handled
- ✅ Mock data is realistic and valid
