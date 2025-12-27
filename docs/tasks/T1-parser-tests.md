# Test Task T1: Parser Tests

## Description
Implement comprehensive tests for all dependency parsers:
- `parseGradleDeps()`
- `parseNpmDeps()`
- `parseMavenDeps()`
- `parseGoModDeps()`

## File Location
`internal/analyzer/parsers_test.go`

## Test Cases Required

### ParseGradleDeps Tests (10 tests)
1. `TestParseGradleDeps_Empty` - Empty content returns empty slice
2. `TestParseGradleDeps_SingleDependency` - Parses implementation() correctly
3. `TestParseGradleDeps_ApiDependency` - Parses api() correctly
4. `TestParseGradleDeps_MultipleTypes` - Mixed implementation and api
5. `TestParseGradleDeps_WithVersions` - Extracts versions correctly
6. `TestParseGradleDeps_MalformedSyntax` - Handles invalid syntax gracefully
7. `TestParseGradleDeps_IgnoresComments` - Ignores commented lines
8. `TestParseGradleDeps_ComplexVersions` - Handles version ranges
9. `TestParseGradleDeps_NoImportStatements` - Empty gradle file
10. `TestParseGradleDeps_ExtractsName` - Correctly extracts dependency names

### ParseNpmDeps Tests (10 tests)
1. `TestParseNpmDeps_ValidDependencies` - Parses dependencies object
2. `TestParseNpmDeps_ValidDevDependencies` - Parses devDependencies object
3. `TestParseNpmDeps_BothDependenciesAndDev` - Parses both sections
4. `TestParseNpmDeps_MalformedJSON` - Returns nil on invalid JSON
5. `TestParseNpmDeps_EmptyObject` - Handles empty dependencies
6. `TestParseNpmDeps_ComplexVersions` - Handles semver, ranges, git refs
7. `TestParseNpmDeps_NoGivenDependencies` - No dependencies section returns nil
8. `TestParseNpmDeps_NestedObjects` - Handles nested JSON gracefully
9. `TestParseNpmDeps_TypesCorrect` - Type is "npm" or "npm-dev"
10. `TestParseNpmDeps_LargePackageJson` - Handles big package.json files

### ParseMavenDeps Tests (10 tests)
1. `TestParseMavenDeps_ValidDependencies` - Parses pom.xml correctly
2. `TestParseMavenDeps_FilterTestScope` - Skips test-scoped dependencies
3. `TestParseMavenDeps_MissingVersion` - Handles missing version tag
4. `TestParseMavenDeps_MalformedXML` - Returns nil on invalid XML
5. `TestParseMavenDeps_GroupIdArtifactFormat` - Correctly formats as groupId:artifactId
6. `TestParseMavenDeps_MultipleScopes` - Handles provided, runtime, compile scopes
7. `TestParseMavenDeps_NoDependencies` - No dependencies section returns nil
8. `TestParseMavenDeps_NestedDependencies` - Handles complex POM structure
9. `TestParseMavenDeps_EmptyDependenciesTag` - Empty dependencies element returns nil
10. `TestParseMavenDeps_TypeIsCorrect` - Type is always "maven"

### ParseGoModDeps Tests (10 tests)
1. `TestParseGoModDeps_ValidRequire` - Parses require block correctly
2. `TestParseGoModDeps_MultipleRequirements` - Multiple dependencies with versions
3. `TestParseGoModDeps_MissingVersion` - Handles missing version gracefully
4. `TestParseGoModDeps_NoRequireBlock` - Empty go.mod file
5. `TestParseGoModDeps_MalformedVersion` - Handles invalid semver
6. `TestParseGoModDeps_ComplexModuleNames` - Handles nested module paths
7. `TestParseGoModDeps_ReplaceDirectives` - Ignores replace directives
8. `TestParseGoModDeps_Whitespace` - Handles various whitespace
9. `TestParseGoModDeps_EmptyRequireBlock` - Empty require() block
10. `TestParseGoModDeps_TypeIsCorrect` - Type is always "go"

## Test Data Fixtures
Create test fixtures in `internal/analyzer/testdata/`:
- `gradle_basic.gradle.kts`
- `package_with_dev.json`
- `pom_with_deps.xml`
- `go_mod_valid.mod`

## Implementation Notes
- Use table-driven tests where appropriate
- Include edge cases and error scenarios
- Test both valid and invalid input
- Verify correct field population (Name, Version, Type)

## Success Criteria
- ✅ All 40+ tests pass
- ✅ Code coverage for parser functions ≥ 95%
- ✅ No external dependencies (use test fixtures)
