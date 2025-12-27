# Task 08: Implement NPM Dependency Parser

## Deliverable
Parse `package.json` to extract npm dependencies.

## Context
The `parseNpmDeps` function was a stub returning empty slice (removed in Task 04). This task implements actual parsing of package.json files to extract dependencies.

## Key Decisions and Principles
- Use Go's `encoding/json` for proper JSON parsing
- Extract both `dependencies` and `devDependencies`
- Mark dev dependencies with a flag or separate type
- Handle malformed JSON gracefully (return empty slice, don't crash)

## Delivers
Working npm dependency parser that populates the Dependencies field for JavaScript/TypeScript projects.

## Acceptance Criteria
- New `parseNpmDeps(content string) []Dependency` function
- Parses `dependencies` object from package.json
- Parses `devDependencies` object from package.json
- Each dependency has Name, Version, and Type="npm" or "npm-dev"
- Malformed JSON returns empty slice without error
- Integration with `fetchDependencies` and `fetchLocalDependencies`

## Dependencies
Task 04 should be completed first (removes the stub).

## Related Code

**File to modify:** `internal/analyzer/analyzer.go`

**Add the implementation (after `parseGoModDeps` function, around line 276):**

```go
// packageJSON represents the structure of package.json
type packageJSON struct {
    Dependencies    map[string]string `json:"dependencies"`
    DevDependencies map[string]string `json:"devDependencies"`
}

// parseNpmDeps extracts dependencies from package.json content
func parseNpmDeps(content string) []Dependency {
    var pkg packageJSON
    if err := json.Unmarshal([]byte(content), &pkg); err != nil {
        return nil
    }

    var deps []Dependency

    // Parse production dependencies
    for name, version := range pkg.Dependencies {
        deps = append(deps, Dependency{
            Name:    name,
            Version: version,
            Type:    "npm",
        })
    }

    // Parse dev dependencies
    for name, version := range pkg.DevDependencies {
        deps = append(deps, Dependency{
            Name:    name,
            Version: version,
            Type:    "npm-dev",
        })
    }

    return deps
}
```

**Add import for `encoding/json`:**
```go
import (
    "encoding/json"
    "fmt"
    "strings"
    // ... existing imports
)
```

**Update `fetchDependencies` to call the parser (modify the javascript/typescript case):**

```go
case "javascript", "typescript":
    if content, err := gh.GetFileContent(repo.FullName, "package.json"); err == nil {
        deps = append(deps, parseNpmDeps(content)...)
    }
```

**Update `fetchLocalDependencies` in `local_analyzer.go` similarly:**

```go
case "javascript", "typescript":
    if content, err := reader.GetFileContent(repo.Path, "package.json"); err == nil {
        deps = append(deps, parseNpmDeps(content)...)
    }
```

## Verification

```bash
# Build should succeed
go build ./...

# Create test
cat > internal/analyzer/npm_test.go << 'EOF'
package analyzer

import "testing"

func TestParseNpmDeps_Valid(t *testing.T) {
    content := `{
        "name": "my-app",
        "dependencies": {
            "react": "^18.2.0",
            "axios": "1.4.0"
        },
        "devDependencies": {
            "typescript": "^5.0.0",
            "jest": "^29.0.0"
        }
    }`

    deps := parseNpmDeps(content)

    if len(deps) != 4 {
        t.Errorf("expected 4 deps, got %d", len(deps))
    }

    // Check for specific dependency
    found := false
    for _, d := range deps {
        if d.Name == "react" && d.Version == "^18.2.0" && d.Type == "npm" {
            found = true
            break
        }
    }
    if !found {
        t.Error("expected to find react dependency")
    }

    // Check for dev dependency
    foundDev := false
    for _, d := range deps {
        if d.Name == "typescript" && d.Type == "npm-dev" {
            foundDev = true
            break
        }
    }
    if !foundDev {
        t.Error("expected to find typescript as dev dependency")
    }
}

func TestParseNpmDeps_Invalid(t *testing.T) {
    content := `not valid json`
    deps := parseNpmDeps(content)
    if deps != nil && len(deps) != 0 {
        t.Errorf("expected empty slice for invalid JSON, got %d deps", len(deps))
    }
}

func TestParseNpmDeps_Empty(t *testing.T) {
    content := `{"name": "empty-app"}`
    deps := parseNpmDeps(content)
    if len(deps) != 0 {
        t.Errorf("expected 0 deps, got %d", len(deps))
    }
}
EOF

go test ./internal/analyzer/... -v -run TestParseNpmDeps
```

## Files Changed
- `internal/analyzer/analyzer.go` (add ~35 lines)
- `internal/analyzer/local_analyzer.go` (modify ~3 lines)
- `internal/analyzer/npm_test.go` (new file, ~60 lines)
