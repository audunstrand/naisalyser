package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper function to load test fixtures
func loadFixture(t *testing.T, filename string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("failed to load fixture %s: %v", filename, err)
	}
	return string(content)
}

// ========================================
// ParseGradleDeps Tests (10 tests)
// ========================================

func TestParseGradleDeps_Empty(t *testing.T) {
	deps := parseGradleDeps("")
	if len(deps) != 0 {
		t.Errorf("expected empty slice, got %d dependencies", len(deps))
	}
}

func TestParseGradleDeps_SingleDependency(t *testing.T) {
	content := `implementation("org.example:lib:1.0.0")`
	deps := parseGradleDeps(content)
	
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}
	
	if deps[0].Name != "org.example:lib:1.0.0" {
		t.Errorf("expected name 'org.example:lib:1.0.0', got '%s'", deps[0].Name)
	}
	if deps[0].Type != "gradle" {
		t.Errorf("expected type 'gradle', got '%s'", deps[0].Type)
	}
}

func TestParseGradleDeps_ApiDependency(t *testing.T) {
	content := `api("com.example:api-lib:2.0.0")`
	deps := parseGradleDeps(content)
	
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}
	
	if deps[0].Name != "com.example:api-lib:2.0.0" {
		t.Errorf("expected name 'com.example:api-lib:2.0.0', got '%s'", deps[0].Name)
	}
}

func TestParseGradleDeps_MultipleTypes(t *testing.T) {
	content := loadFixture(t, "gradle_basic.gradle.kts")
	deps := parseGradleDeps(content)
	
	if len(deps) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(deps))
	}
	
	// Check that we have both implementation and api types
	hasImplementation := false
	hasApi := false
	for _, dep := range deps {
		if dep.Name == "org.springframework.boot:spring-boot-starter-web:3.0.0" {
			hasImplementation = true
		}
		if dep.Name == "com.fasterxml.jackson.core:jackson-databind:2.14.0" {
			hasApi = true
		}
	}
	
	if !hasImplementation || !hasApi {
		t.Errorf("expected both implementation and api dependencies")
	}
}

func TestParseGradleDeps_WithVersions(t *testing.T) {
	content := `implementation("org.example:lib:1.2.3")`
	deps := parseGradleDeps(content)
	
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}
	
	// Note: Current implementation doesn't parse version separately
	// It includes the full string in Name
	if deps[0].Name != "org.example:lib:1.2.3" {
		t.Errorf("expected full name with version, got '%s'", deps[0].Name)
	}
}

func TestParseGradleDeps_MalformedSyntax(t *testing.T) {
	content := loadFixture(t, "gradle_malformed.gradle.kts")
	deps := parseGradleDeps(content)
	
	// Should handle gracefully - malformed lines are ignored
	// Only valid lines are parsed
	if len(deps) > 2 {
		t.Errorf("expected malformed syntax to be ignored, got %d deps", len(deps))
	}
}

func TestParseGradleDeps_IgnoresComments(t *testing.T) {
	content := loadFixture(t, "gradle_with_comments.gradle.kts")
	deps := parseGradleDeps(content)
	
	// Should only find 2 dependencies (commented ones should be ignored)
	if len(deps) != 2 {
		t.Errorf("expected 2 dependencies (ignoring comments), got %d", len(deps))
	}
	
	// Verify we don't have commented dependencies
	for _, dep := range deps {
		if dep.Name == "org.commented:commented-lib:1.0.0" {
			t.Error("commented dependency should not be parsed")
		}
	}
}

func TestParseGradleDeps_ComplexVersions(t *testing.T) {
	content := loadFixture(t, "gradle_complex_versions.gradle.kts")
	deps := parseGradleDeps(content)
	
	if len(deps) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(deps))
	}
	
	// Verify complex version formats are captured
	expectedVersions := map[string]bool{
		"org.example:lib-with-range:1.0.+":           false,
		"com.example:snapshot-version:2.0.0-SNAPSHOT": false,
		"io.example:semantic-version:1.2.3-beta.4":   false,
	}
	
	for _, dep := range deps {
		if _, exists := expectedVersions[dep.Name]; exists {
			expectedVersions[dep.Name] = true
		}
	}
	
	for name, found := range expectedVersions {
		if !found {
			t.Errorf("expected to find dependency with version: %s", name)
		}
	}
}

func TestParseGradleDeps_NoImportStatements(t *testing.T) {
	content := loadFixture(t, "gradle_empty.gradle.kts")
	deps := parseGradleDeps(content)
	
	if len(deps) != 0 {
		t.Errorf("expected no dependencies in empty gradle file, got %d", len(deps))
	}
}

func TestParseGradleDeps_ExtractsName(t *testing.T) {
	content := `implementation("com.example:my-library:1.0.0")`
	deps := parseGradleDeps(content)
	
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}
	
	// Verify the full dependency string is extracted
	if deps[0].Name != "com.example:my-library:1.0.0" {
		t.Errorf("dependency name not correctly extracted: got '%s'", deps[0].Name)
	}
}

// ========================================
// ParseNpmDeps Tests (10 tests)
// ========================================

func TestParseNpmDeps_ValidDependencies(t *testing.T) {
	content := loadFixture(t, "package_basic.json")
	deps := parseNpmDeps(content)
	
	if len(deps) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(deps))
	}
	
	// Verify types are correct
	for _, dep := range deps {
		if dep.Type != "npm" {
			t.Errorf("expected type 'npm', got '%s'", dep.Type)
		}
	}
	
	// Verify specific dependencies
	found := make(map[string]string)
	for _, dep := range deps {
		found[dep.Name] = dep.Version
	}
	
	if found["express"] != "^4.18.2" {
		t.Errorf("expected express version '^4.18.2', got '%s'", found["express"])
	}
}

func TestParseNpmDeps_ValidDevDependencies(t *testing.T) {
	content := loadFixture(t, "package_dev.json")
	deps := parseNpmDeps(content)
	
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(deps))
	}
	
	// Verify types are correct for dev dependencies
	for _, dep := range deps {
		if dep.Type != "npm-dev" {
			t.Errorf("expected type 'npm-dev' for dev dependencies, got '%s'", dep.Type)
		}
	}
}

func TestParseNpmDeps_BothDependenciesAndDev(t *testing.T) {
	content := loadFixture(t, "package_with_dev.json")
	deps := parseNpmDeps(content)
	
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependencies (1 prod + 1 dev), got %d", len(deps))
	}
	
	// Verify we have both types
	hasNpm := false
	hasNpmDev := false
	for _, dep := range deps {
		if dep.Type == "npm" {
			hasNpm = true
		}
		if dep.Type == "npm-dev" {
			hasNpmDev = true
		}
	}
	
	if !hasNpm || !hasNpmDev {
		t.Error("expected both 'npm' and 'npm-dev' types")
	}
}

func TestParseNpmDeps_MalformedJSON(t *testing.T) {
	content := loadFixture(t, "package_malformed.json")
	deps := parseNpmDeps(content)
	
	if deps != nil {
		t.Errorf("expected nil for malformed JSON, got %d dependencies", len(deps))
	}
}

func TestParseNpmDeps_EmptyObject(t *testing.T) {
	content := loadFixture(t, "package_empty.json")
	deps := parseNpmDeps(content)
	
	// Valid JSON but no dependencies - parser returns empty slice (length 0)
	if len(deps) != 0 {
		t.Errorf("expected no dependencies, got %d", len(deps))
	}
}

func TestParseNpmDeps_ComplexVersions(t *testing.T) {
	content := loadFixture(t, "package_complex_versions.json")
	deps := parseNpmDeps(content)
	
	if len(deps) != 4 {
		t.Fatalf("expected 4 dependencies, got %d", len(deps))
	}
	
	// Verify complex version formats are preserved
	found := make(map[string]string)
	for _, dep := range deps {
		found[dep.Name] = dep.Version
	}
	
	if found["git-package"] != "git+https://github.com/user/repo.git#v1.0.0" {
		t.Errorf("git version not preserved correctly")
	}
	if found["wildcard"] != "*" {
		t.Errorf("wildcard version not preserved correctly")
	}
}

func TestParseNpmDeps_NoGivenDependencies(t *testing.T) {
	content := `{"name": "test", "version": "1.0.0"}`
	deps := parseNpmDeps(content)
	
	// Parser returns empty slice when no dependencies section
	if len(deps) != 0 {
		t.Errorf("expected 0 dependencies, got %d", len(deps))
	}
}

func TestParseNpmDeps_NestedObjects(t *testing.T) {
	content := `{
		"name": "test",
		"config": {
			"dependencies": {
				"nested": "1.0.0"
			}
		},
		"dependencies": {
			"real-dep": "2.0.0"
		}
	}`
	deps := parseNpmDeps(content)
	
	// Should only parse the root dependencies
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}
	if deps[0].Name != "real-dep" {
		t.Errorf("expected 'real-dep', got '%s'", deps[0].Name)
	}
}

func TestParseNpmDeps_TypesCorrect(t *testing.T) {
	content := `{
		"dependencies": {
			"prod-dep": "1.0.0"
		},
		"devDependencies": {
			"dev-dep": "2.0.0"
		}
	}`
	deps := parseNpmDeps(content)
	
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(deps))
	}
	
	typeMap := make(map[string]string)
	for _, dep := range deps {
		typeMap[dep.Name] = dep.Type
	}
	
	if typeMap["prod-dep"] != "npm" {
		t.Errorf("expected 'npm' type for prod-dep, got '%s'", typeMap["prod-dep"])
	}
	if typeMap["dev-dep"] != "npm-dev" {
		t.Errorf("expected 'npm-dev' type for dev-dep, got '%s'", typeMap["dev-dep"])
	}
}

func TestParseNpmDeps_LargePackageJson(t *testing.T) {
	// Create a large package.json programmatically
	largeJSON := `{
		"dependencies": {
			"dep1": "1.0.0",
			"dep2": "2.0.0",
			"dep3": "3.0.0",
			"dep4": "4.0.0",
			"dep5": "5.0.0"
		},
		"devDependencies": {
			"devdep1": "1.0.0",
			"devdep2": "2.0.0",
			"devdep3": "3.0.0"
		}
	}`
	
	deps := parseNpmDeps(largeJSON)
	
	if len(deps) != 8 {
		t.Errorf("expected 8 dependencies, got %d", len(deps))
	}
}

// ========================================
// ParseMavenDeps Tests (10 tests)
// ========================================

func TestParseMavenDeps_ValidDependencies(t *testing.T) {
	content := loadFixture(t, "pom_basic.xml")
	deps := parseMavenDeps(content)
	
	if len(deps) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(deps))
	}
	
	// Verify type is correct
	for _, dep := range deps {
		if dep.Type != "maven" {
			t.Errorf("expected type 'maven', got '%s'", dep.Type)
		}
	}
}

func TestParseMavenDeps_FilterTestScope(t *testing.T) {
	content := loadFixture(t, "pom_with_test_scope.xml")
	deps := parseMavenDeps(content)
	
	// Should only have 1 dependency (spring-core), test-scoped ones should be filtered
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency (test scope filtered), got %d", len(deps))
	}
	
	if deps[0].Name != "org.springframework:spring-core" {
		t.Errorf("expected spring-core, got '%s'", deps[0].Name)
	}
}

func TestParseMavenDeps_MissingVersion(t *testing.T) {
	content := loadFixture(t, "pom_missing_version.xml")
	deps := parseMavenDeps(content)
	
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}
	
	// Version should be empty string
	if deps[0].Version != "" {
		t.Errorf("expected empty version, got '%s'", deps[0].Version)
	}
}

func TestParseMavenDeps_MalformedXML(t *testing.T) {
	content := loadFixture(t, "pom_malformed.xml")
	deps := parseMavenDeps(content)
	
	if deps != nil {
		t.Errorf("expected nil for malformed XML, got %d dependencies", len(deps))
	}
}

func TestParseMavenDeps_GroupIdArtifactFormat(t *testing.T) {
	content := loadFixture(t, "pom_basic.xml")
	deps := parseMavenDeps(content)
	
	if len(deps) < 1 {
		t.Fatal("expected at least 1 dependency")
	}
	
	// Verify format is groupId:artifactId
	if deps[0].Name != "org.springframework.boot:spring-boot-starter-web" {
		t.Errorf("expected groupId:artifactId format, got '%s'", deps[0].Name)
	}
}

func TestParseMavenDeps_MultipleScopes(t *testing.T) {
	content := loadFixture(t, "pom_with_scopes.xml")
	deps := parseMavenDeps(content)
	
	// Should have 3 dependencies (all non-test scopes)
	if len(deps) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(deps))
	}
	
	// Verify all have maven type
	for _, dep := range deps {
		if dep.Type != "maven" {
			t.Errorf("expected type 'maven', got '%s'", dep.Type)
		}
	}
}

func TestParseMavenDeps_NoDependencies(t *testing.T) {
	content := loadFixture(t, "pom_no_deps.xml")
	deps := parseMavenDeps(content)
	
	// Parser returns empty slice when no dependencies
	if len(deps) != 0 {
		t.Errorf("expected 0 dependencies, got %d", len(deps))
	}
}

func TestParseMavenDeps_NestedDependencies(t *testing.T) {
	// POM with nested structure
	content := `<?xml version="1.0"?>
<project>
	<parent>
		<groupId>org.parent</groupId>
		<artifactId>parent-pom</artifactId>
	</parent>
	<dependencies>
		<dependency>
			<groupId>com.example</groupId>
			<artifactId>lib</artifactId>
			<version>1.0.0</version>
		</dependency>
	</dependencies>
</project>`
	
	deps := parseMavenDeps(content)
	
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}
	if deps[0].Name != "com.example:lib" {
		t.Errorf("expected 'com.example:lib', got '%s'", deps[0].Name)
	}
}

func TestParseMavenDeps_EmptyDependenciesTag(t *testing.T) {
	content := loadFixture(t, "pom_empty_deps.xml")
	deps := parseMavenDeps(content)
	
	// Empty dependencies tag returns empty slice
	if len(deps) != 0 {
		t.Errorf("expected 0 dependencies, got %d", len(deps))
	}
}

func TestParseMavenDeps_TypeIsCorrect(t *testing.T) {
	content := loadFixture(t, "pom_basic.xml")
	deps := parseMavenDeps(content)
	
	for _, dep := range deps {
		if dep.Type != "maven" {
			t.Errorf("expected all dependencies to have type 'maven', got '%s'", dep.Type)
		}
	}
}

// ========================================
// ParseGoModDeps Tests (10 tests)
// ========================================

func TestParseGoModDeps_ValidRequire(t *testing.T) {
	content := loadFixture(t, "go_mod_basic.mod")
	deps := parseGoModDeps(content)
	
	if len(deps) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(deps))
	}
	
	// Verify types are correct
	for _, dep := range deps {
		if dep.Type != "go" {
			t.Errorf("expected type 'go', got '%s'", dep.Type)
		}
	}
}

func TestParseGoModDeps_MultipleRequirements(t *testing.T) {
	content := loadFixture(t, "go_mod_basic.mod")
	deps := parseGoModDeps(content)
	
	if len(deps) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(deps))
	}
	
	// Verify specific dependencies and versions
	found := make(map[string]string)
	for _, dep := range deps {
		found[dep.Name] = dep.Version
	}
	
	if found["github.com/gorilla/mux"] != "v1.8.0" {
		t.Errorf("expected version 'v1.8.0', got '%s'", found["github.com/gorilla/mux"])
	}
}

func TestParseGoModDeps_MissingVersion(t *testing.T) {
	content := loadFixture(t, "go_mod_missing_version.mod")
	deps := parseGoModDeps(content)
	
	// Should have 1 with version, 1 without (ignored or empty version)
	hasWithVersion := false
	for _, dep := range deps {
		if dep.Name == "github.com/with/version" && dep.Version == "v1.2.3" {
			hasWithVersion = true
		}
	}
	
	if !hasWithVersion {
		t.Error("expected to find dependency with version")
	}
}

func TestParseGoModDeps_NoRequireBlock(t *testing.T) {
	content := loadFixture(t, "go_mod_no_require.mod")
	deps := parseGoModDeps(content)
	
	if len(deps) != 0 {
		t.Errorf("expected no dependencies, got %d", len(deps))
	}
}

func TestParseGoModDeps_MalformedVersion(t *testing.T) {
	content := `module test
require (
	github.com/example/module invalid-version
)`
	deps := parseGoModDeps(content)
	
	// Should still parse, version is just treated as string
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}
	if deps[0].Version != "invalid-version" {
		t.Errorf("expected version 'invalid-version', got '%s'", deps[0].Version)
	}
}

func TestParseGoModDeps_ComplexModuleNames(t *testing.T) {
	content := loadFixture(t, "go_mod_complex_names.mod")
	deps := parseGoModDeps(content)
	
	if len(deps) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(deps))
	}
	
	// Verify complex module paths are preserved
	found := make(map[string]bool)
	for _, dep := range deps {
		found[dep.Name] = true
	}
	
	if !found["github.com/deep/nested/module/path/example"] {
		t.Error("expected to find deep nested module path")
	}
	if !found["k8s.io/api"] {
		t.Error("expected to find k8s.io/api")
	}
}

func TestParseGoModDeps_ReplaceDirectives(t *testing.T) {
	content := loadFixture(t, "go_mod_with_replace.mod")
	deps := parseGoModDeps(content)
	
	// Should only parse require block, ignore replace
	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency (replace ignored), got %d", len(deps))
	}
	
	if deps[0].Name != "github.com/original/module" {
		t.Errorf("expected original module name, got '%s'", deps[0].Name)
	}
}

func TestParseGoModDeps_Whitespace(t *testing.T) {
	content := loadFixture(t, "go_mod_whitespace.mod")
	deps := parseGoModDeps(content)
	
	if len(deps) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(deps))
	}
	
	// Verify all dependencies are parsed despite varying whitespace
	expectedModules := map[string]bool{
		"github.com/standard/module": false,
		"github.com/extra/spaces":    false,
		"github.com/tabs/module":     false,
	}
	
	for _, dep := range deps {
		if _, exists := expectedModules[dep.Name]; exists {
			expectedModules[dep.Name] = true
		}
	}
	
	for module, found := range expectedModules {
		if !found {
			t.Errorf("expected to find module: %s", module)
		}
	}
}

func TestParseGoModDeps_EmptyRequireBlock(t *testing.T) {
	content := loadFixture(t, "go_mod_empty_require.mod")
	deps := parseGoModDeps(content)
	
	if len(deps) != 0 {
		t.Errorf("expected no dependencies in empty require block, got %d", len(deps))
	}
}

func TestParseGoModDeps_TypeIsCorrect(t *testing.T) {
	content := loadFixture(t, "go_mod_basic.mod")
	deps := parseGoModDeps(content)
	
	for _, dep := range deps {
		if dep.Type != "go" {
			t.Errorf("expected all dependencies to have type 'go', got '%s'", dep.Type)
		}
	}
}
