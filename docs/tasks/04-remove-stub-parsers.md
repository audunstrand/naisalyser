# Task 04: Remove Stub Dependency Parsers

## Deliverable
Remove stub `parseMavenDeps` and `parseNpmDeps` functions that return empty slices, or mark them clearly as unimplemented.

## Context
In `internal/analyzer/analyzer.go`, two functions exist that do nothing:

```go
func parseMavenDeps(content string) []Dependency {
    // Very basic XML parsing - in production use proper XML parser
    return []Dependency{}
}

func parseNpmDeps(content string) []Dependency {
    // Would parse package.json
    return []Dependency{}
}
```

These are called in `fetchDependencies()` but always return empty slices, giving users the false impression that dependencies are being parsed for Java/JS projects.

## Key Decisions and Principles
- Option A: Remove the functions entirely (chosen for this task - YAGNI)
- Option B: Keep but add TODO comments and log warnings when called
- We choose Option A because the functions are misleading in their current state

## Delivers
Removed dead code that pretends to parse dependencies but doesn't.

## Acceptance Criteria
- `parseMavenDeps` function removed from `analyzer.go`
- `parseNpmDeps` function removed from `analyzer.go`
- Calls to these functions removed from `fetchDependencies()`
- Comments added indicating these parsers are not yet implemented
- Code compiles

## Dependencies
None.

## Related Code

**File to modify:** `internal/analyzer/analyzer.go`

**Current `fetchDependencies` function (lines 194-219):**
```go
func fetchDependencies(gh *github.Client, repo *github.Repository) ([]Dependency, error) {
    var deps []Dependency

    switch strings.ToLower(repo.Language) {
    case "kotlin", "java":
        if content, err := gh.GetFileContent(repo.FullName, "build.gradle.kts"); err == nil {
            deps = append(deps, parseGradleDeps(content)...)
        }
        if content, err := gh.GetFileContent(repo.FullName, "pom.xml"); err == nil {
            deps = append(deps, parseMavenDeps(content)...)  // ← RETURNS EMPTY
        }
    case "javascript", "typescript":
        if content, err := gh.GetFileContent(repo.FullName, "package.json"); err == nil {
            deps = append(deps, parseNpmDeps(content)...)  // ← RETURNS EMPTY
        }
    case "go":
        if content, err := gh.GetFileContent(repo.FullName, "go.mod"); err == nil {
            deps = append(deps, parseGoModDeps(content)...)
        }
    }

    return deps, nil
}
```

**Replace with:**
```go
func fetchDependencies(gh *github.Client, repo *github.Repository) ([]Dependency, error) {
    var deps []Dependency

    switch strings.ToLower(repo.Language) {
    case "kotlin", "java":
        // Gradle dependencies are parsed
        if content, err := gh.GetFileContent(repo.FullName, "build.gradle.kts"); err == nil {
            deps = append(deps, parseGradleDeps(content)...)
        }
        // TODO: Maven pom.xml parsing not implemented
    case "javascript", "typescript":
        // TODO: package.json parsing not implemented
    case "go":
        if content, err := gh.GetFileContent(repo.FullName, "go.mod"); err == nil {
            deps = append(deps, parseGoModDeps(content)...)
        }
    }

    return deps, nil
}
```

**Remove these functions entirely (lines 239-247):**
```go
// DELETE THIS:
func parseMavenDeps(content string) []Dependency {
    // Very basic XML parsing - in production use proper XML parser
    return []Dependency{}
}

// DELETE THIS:
func parseNpmDeps(content string) []Dependency {
    // Would parse package.json
    return []Dependency{}
}
```

**Also update `fetchLocalDependencies` in `local_analyzer.go` (lines 203-228):**

```go
func fetchLocalDependencies(reader *local.Reader, repo *local.Repository) ([]Dependency, error) {
    var deps []Dependency

    switch strings.ToLower(repo.Language) {
    case "kotlin", "java":
        if content, err := reader.GetFileContent(repo.Path, "build.gradle.kts"); err == nil {
            deps = append(deps, parseGradleDeps(content)...)
        }
        // TODO: Maven pom.xml parsing not implemented
    case "javascript", "typescript":
        // TODO: package.json parsing not implemented
    case "go":
        if content, err := reader.GetFileContent(repo.Path, "go.mod"); err == nil {
            deps = append(deps, parseGoModDeps(content)...)
        }
    }

    return deps, nil
}
```

## Verification

```bash
# Build should succeed
go build ./...

# Verify no references to removed functions
grep -r "parseMavenDeps\|parseNpmDeps" internal/

# Should return empty (no matches)
```

## Files Changed
- `internal/analyzer/analyzer.go` (remove ~12 lines, modify ~8 lines)
- `internal/analyzer/local_analyzer.go` (modify ~6 lines)
