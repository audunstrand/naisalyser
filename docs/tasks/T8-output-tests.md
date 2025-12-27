# Test Task T8: Output Generation Tests

## Description
Implement comprehensive tests for markdown documentation generation:
- `GenerateMarkdown()` - main function
- `toViewModel()` - converts analysis to view model
- `generateOverviewVM()` - generates overview section
- `generateDependenciesVM()` - generates dependencies section
- `generateNaisConfigVM()` - generates configuration section
- `executeTemplate()` - template rendering

## File Location
`internal/output/markdown_test.go`

## Test Cases Required

### GenerateMarkdown Tests (8 tests)
1. `TestGenerateMarkdown_CreatesOutputDir` - Creates directory if missing
2. `TestGenerateMarkdown_WritesReadmeFile` - Creates README.md
3. `TestGenerateMarkdown_InvalidOutputPath` - Error for invalid path
4. `TestGenerateMarkdown_PermissionError` - Error when directory not writable
5. `TestGenerateMarkdown_CompleteOutput` - All sections in output
6. `TestGenerateMarkdown_ValidMarkdownSyntax` - Output is valid markdown
7. `TestGenerateMarkdown_ExistingDir` - Works with existing directory
8. `TestGenerateMarkdown_FilesCreated` - Correct number of files created

### ToViewModel Tests (8 tests)
1. `TestToViewModel_NilAnalysis` - Handles nil analysis gracefully
2. `TestToViewModel_NilRepository` - Handles nil repository
3. `TestToViewModel_NilNaisConfig` - Handles nil NAIS config
4. `TestToViewModel_NilDependencies` - Handles nil dependencies slice
5. `TestToViewModel_FullData` - Converts all fields correctly
6. `TestToViewModel_DataTypes` - All types are correct in output
7. `TestToViewModel_EmptyStrings` - Handles empty string fields
8. `TestToViewModel_SpecialCharacters` - Escapes markdown special chars

### GenerateOverviewVM Tests (8 tests)
1. `TestGenerateOverviewVM_RepoMetadata` - Includes repo name, description
2. `TestGenerateOverviewVM_Language` - Shows detected language
3. `TestGenerateOverviewVM_Namespace` - Shows repository namespace
4. `TestGenerateOverviewVM_MissingFields` - Handles missing metadata
5. `TestGenerateOverviewVM_EmptyDescription` - Handles no description
6. `TestGenerateOverviewVM_LongDescription` - Handles very long descriptions
7. `TestGenerateOverviewVM_UrlFormatting` - Formats URLs correctly
8. `TestGenerateOverviewVM_DateFormatting` - Formats dates if applicable

### GenerateNaisConfigVM Tests (12 tests)

#### Access Policy (4 tests)
1. `TestGenerateNaisConfigVM_InboundRules` - Lists inbound applications
2. `TestGenerateNaisConfigVM_OutboundRules` - Lists outbound applications
3. `TestGenerateNaisConfigVM_ExternalHosts` - Lists external hosts
4. `TestGenerateNaisConfigVM_NoAccessPolicy` - Handles missing access policy

#### Database & Storage (4 tests)
5. `TestGenerateNaisConfigVM_SQLInstances` - Lists databases
6. `TestGenerateNaisConfigVM_Buckets` - Lists GCS buckets
7. `TestGenerateNaisConfigVM_NoStorage` - Handles no storage configured
8. `TestGenerateNaisConfigVM_StorageMetadata` - Shows size, type, etc.

#### Kafka (3 tests)
9. `TestGenerateNaisConfigVM_KafkaConfig` - Shows Kafka pool
10. `TestGenerateNaisConfigVM_KafkaTopics` - Lists Kafka topics
11. `TestGenerateNaisConfigVM_NoKafka` - Handles no Kafka configured

#### General (1 test)
12. `TestGenerateNaisConfigVM_AllSections` - All config sections present

### GenerateDependenciesVM Tests (10 tests)
1. `TestGenerateDependenciesVM_GradleDeps` - Groups Gradle dependencies
2. `TestGenerateDependenciesVM_MavenDeps` - Groups Maven dependencies
3. `TestGenerateDependenciesVM_NpmDeps` - Groups NPM dependencies
4. `TestGenerateDependenciesVM_NpmDevDeps` - Groups dev dependencies
5. `TestGenerateDependenciesVM_GoDeps` - Groups Go dependencies
6. `TestGenerateDependenciesVM_NoDependencies` - Handles empty list
7. `TestGenerateDependenciesVM_MixedTypes` - Multiple dependency types
8. `TestGenerateDependenciesVM_DependencyCount` - Shows correct counts
9. `TestGenerateDependenciesVM_Formatting` - Table format is correct
10. `TestGenerateDependenciesVM_Sorting` - Dependencies are sorted

### ExecuteTemplate Tests (6 tests)
1. `TestExecuteTemplate_ValidTemplate` - Renders valid template
2. `TestExecuteTemplate_MissingTemplate` - Error for missing template
3. `TestExecuteTemplate_InvalidData` - Handles bad data gracefully
4. `TestExecuteTemplate_TemplateErrors` - Reports template syntax errors
5. `TestExecuteTemplate_OutputIsMarkdown` - Output is valid markdown
6. `TestExecuteTemplate_SpecialCharacters` - Escapes template variables

## Test Data Fixtures

### Sample Analysis Data
```go
func createTestAnalysis() *analyzer.Analysis {
    return &analyzer.Analysis{
        Repository: &github.Repository{
            Name:        "my-app",
            FullName:    "navikt/my-app",
            Description: "My test application",
            Language:    "kotlin",
        },
        NaisConfig: &analyzer.NaisConfig{
            Spec: analyzer.NaisSpec{
                GCP: &analyzer.NaisGCP{
                    SQLInstances: []analyzer.NaisSQLInstance{
                        {Type: "POSTGRES_14", Tier: "db-custom-2-8192"},
                    },
                    Buckets: []analyzer.NaisBucket{
                        {Name: "my-bucket"},
                    },
                },
                AccessPolicy: &analyzer.NaisAccessPolicy{
                    Inbound: &analyzer.NaisAccessRules{
                        Rules: []analyzer.NaisAccessRule{
                            {Application: "other-app", Namespace: "default"},
                        },
                    },
                    Outbound: &analyzer.NaisAccessRules{
                        Rules: []analyzer.NaisAccessRule{
                            {Application: "api-app", Namespace: "default"},
                        },
                        External: []analyzer.ExternalHost{
                            {Host: "api.example.com"},
                        },
                    },
                },
                Kafka: &analyzer.NaisKafkaConfig{
                    Pool: "nav-prod",
                },
            },
        },
        Dependencies: []analyzer.Dependency{
            {Name: "org.springframework.boot:spring-boot", Version: "3.0.0", Type: "maven"},
            {Name: "express", Version: "^4.18.0", Type: "npm"},
        },
        KafkaTopics: []analyzer.KafkaTopic{
            {Name: "my.events", ConstName: "MY_EVENTS_TOPIC", SourceFile: "src/main/kotlin/Topics.kt"},
        },
    }
}
```

### Expected Markdown Structure
```markdown
# my-app

## Overview
- **Repository:** navikt/my-app
- **Description:** My test application
- **Language:** Kotlin

## Dependencies
| Type | Name | Version |
|------|------|---------|
| Maven | org.springframework.boot:spring-boot | 3.0.0 |
| NPM | express | ^4.18.0 |

## Configuration
### Access Policy
**Inbound Rules:**
- other-app

## Database & Storage
- PostgreSQL instance (db-custom-2-8192)
- GCS bucket: my-bucket

## Kafka
Pool: nav-prod
```

## Markdown Validation Tests

Add tests to verify markdown structure:
```go
func TestGenerateMarkdown_ValidMarkdownStructure(t *testing.T) {
    analysis := createTestAnalysis()
    tmpDir := t.TempDir()
    
    err := GenerateMarkdown(analysis, tmpDir)
    assert.NoError(t, err)
    
    // Read generated file
    content, _ := os.ReadFile(filepath.Join(tmpDir, "README.md"))
    
    // Verify structure
    assert.Contains(t, string(content), "# ")      // H1 header
    assert.Contains(t, string(content), "## ")     // H2 headers
    assert.Contains(t, string(content), "- ")      // Lists
    assert.Contains(t, string(content), "| ")      // Tables
}
```

## Edge Cases to Test
- Very long application names
- Special characters in descriptions
- Very large number of dependencies
- Empty access policies
- No Kafka configured
- Unicode in all fields
- HTML-like characters that need escaping
- URLs in descriptions

## Implementation Notes
- Test both positive and negative cases
- Verify markdown headers are properly formatted
- Check tables have proper alignment
- Verify lists are properly indented
- Test template rendering separately from file I/O
- Mock file system for I/O tests
- Verify escaped characters in output

## Success Criteria
- ✅ All 30+ tests pass
- ✅ Code coverage for output functions ≥ 80%
- ✅ Generated markdown is valid and complete
- ✅ All configuration types displayed correctly
- ✅ Special characters properly escaped
- ✅ Directory creation and file I/O working
