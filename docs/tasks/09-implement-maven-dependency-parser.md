# Task 09: Implement Maven Dependency Parser

## Deliverable
Parse `pom.xml` to extract Maven dependencies.

## Context
The `parseMavenDeps` function was a stub returning empty slice (removed in Task 04). This task implements actual parsing of pom.xml files to extract dependencies.

## Key Decisions and Principles
- Use Go's `encoding/xml` for proper XML parsing
- Extract groupId, artifactId, and version from `<dependency>` elements
- Format name as `groupId:artifactId` (Maven coordinate style)
- Handle missing version (often inherited from parent POM)
- Handle malformed XML gracefully

## Delivers
Working Maven dependency parser that populates the Dependencies field for Java projects using Maven.

## Acceptance Criteria
- New `parseMavenDeps(content string) []Dependency` function
- Parses `<dependency>` elements from pom.xml
- Each dependency has Name=`groupId:artifactId`, Version, Type="maven"
- Version can be empty string if not specified
- Malformed XML returns empty slice without error
- Excludes test-scoped dependencies (scope=test)
- Integration with `fetchDependencies` and `fetchLocalDependencies`

## Dependencies
Task 04 should be completed first (removes the stub).

## Related Code

**File to modify:** `internal/analyzer/analyzer.go`

**Add the implementation (after `parseNpmDeps` or `parseGoModDeps`):**

```go
// pomXML represents the structure of pom.xml
type pomXML struct {
    XMLName      xml.Name        `xml:"project"`
    Dependencies pomDependencies `xml:"dependencies"`
}

type pomDependencies struct {
    Dependency []pomDependency `xml:"dependency"`
}

type pomDependency struct {
    GroupID    string `xml:"groupId"`
    ArtifactID string `xml:"artifactId"`
    Version    string `xml:"version"`
    Scope      string `xml:"scope"`
}

// parseMavenDeps extracts dependencies from pom.xml content
func parseMavenDeps(content string) []Dependency {
    var pom pomXML
    if err := xml.Unmarshal([]byte(content), &pom); err != nil {
        return nil
    }

    var deps []Dependency
    for _, d := range pom.Dependencies.Dependency {
        // Skip test dependencies
        if d.Scope == "test" {
            continue
        }

        name := d.GroupID + ":" + d.ArtifactID
        deps = append(deps, Dependency{
            Name:    name,
            Version: d.Version,
            Type:    "maven",
        })
    }

    return deps
}
```

**Add import for `encoding/xml`:**
```go
import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "strings"
    // ... existing imports
)
```

**Update `fetchDependencies` to call the parser (modify the kotlin/java case):**

```go
case "kotlin", "java":
    if content, err := gh.GetFileContent(repo.FullName, "build.gradle.kts"); err == nil {
        deps = append(deps, parseGradleDeps(content)...)
    }
    if content, err := gh.GetFileContent(repo.FullName, "pom.xml"); err == nil {
        deps = append(deps, parseMavenDeps(content)...)
    }
```

**Update `fetchLocalDependencies` in `local_analyzer.go` similarly:**

```go
case "kotlin", "java":
    if content, err := reader.GetFileContent(repo.Path, "build.gradle.kts"); err == nil {
        deps = append(deps, parseGradleDeps(content)...)
    }
    if content, err := reader.GetFileContent(repo.Path, "pom.xml"); err == nil {
        deps = append(deps, parseMavenDeps(content)...)
    }
```

## Verification

```bash
# Build should succeed
go build ./...

# Create test
cat > internal/analyzer/maven_test.go << 'EOF'
package analyzer

import "testing"

func TestParseMavenDeps_Valid(t *testing.T) {
    content := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
            <version>3.1.0</version>
        </dependency>
        <dependency>
            <groupId>org.postgresql</groupId>
            <artifactId>postgresql</artifactId>
            <version>42.6.0</version>
        </dependency>
        <dependency>
            <groupId>org.junit.jupiter</groupId>
            <artifactId>junit-jupiter</artifactId>
            <version>5.9.3</version>
            <scope>test</scope>
        </dependency>
    </dependencies>
</project>`

    deps := parseMavenDeps(content)

    // Should have 2 deps (test dependency excluded)
    if len(deps) != 2 {
        t.Errorf("expected 2 deps, got %d", len(deps))
    }

    // Check for specific dependency
    found := false
    for _, d := range deps {
        if d.Name == "org.springframework.boot:spring-boot-starter-web" && 
           d.Version == "3.1.0" && 
           d.Type == "maven" {
            found = true
            break
        }
    }
    if !found {
        t.Error("expected to find spring-boot-starter-web dependency")
    }
}

func TestParseMavenDeps_NoVersion(t *testing.T) {
    content := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <dependencies>
        <dependency>
            <groupId>org.example</groupId>
            <artifactId>my-lib</artifactId>
        </dependency>
    </dependencies>
</project>`

    deps := parseMavenDeps(content)
    if len(deps) != 1 {
        t.Fatalf("expected 1 dep, got %d", len(deps))
    }
    if deps[0].Version != "" {
        t.Errorf("expected empty version, got %s", deps[0].Version)
    }
}

func TestParseMavenDeps_Invalid(t *testing.T) {
    content := `not valid xml`
    deps := parseMavenDeps(content)
    if deps != nil && len(deps) != 0 {
        t.Errorf("expected empty slice for invalid XML, got %d deps", len(deps))
    }
}

func TestParseMavenDeps_Empty(t *testing.T) {
    content := `<?xml version="1.0"?><project></project>`
    deps := parseMavenDeps(content)
    if len(deps) != 0 {
        t.Errorf("expected 0 deps, got %d", len(deps))
    }
}
EOF

go test ./internal/analyzer/... -v -run TestParseMavenDeps
```

## Files Changed
- `internal/analyzer/analyzer.go` (add ~45 lines)
- `internal/analyzer/local_analyzer.go` (modify ~3 lines)
- `internal/analyzer/maven_test.go` (new file, ~90 lines)
