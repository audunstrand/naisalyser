# Task 05: Extract RepositoryReader Interface

## Deliverable
Create a `RepositoryReader` interface that both `github.Client` and `local.Reader` implement, eliminating duplicate analysis code.

## Context
Currently there are two nearly identical analysis flows:
- `Analyze()` in `analyzer.go` uses `github.Client`
- `AnalyzeLocal()` in `local_analyzer.go` uses `local.Reader`

Both do the same thing: read files, parse nais config, extract dependencies. The duplication means every change must be made twice.

**Current duplication:**
```go
// analyzer.go:144-168
func Analyze(gh *github.Client, repo *github.Repository) (*Analysis, error) { ... }

// local_analyzer.go:23-50
func AnalyzeLocal(reader *local.Reader, repo *local.Repository) (*LocalAnalysis, error) { ... }
```

## Key Decisions and Principles
- Interface should be minimal - only methods actually used
- Interface lives in `analyzer` package (consumer defines interface)
- Both readers already have compatible method signatures
- Keep backwards compatibility - existing functions still work

## Delivers
A shared interface that allows one `Analyze` function to work with both GitHub and local sources.

## Acceptance Criteria
- `RepositoryReader` interface defined with `GetFileContent` and `FileExists` methods
- `github.Client` satisfies the interface (already has these methods)
- `local.Reader` satisfies the interface (already has these methods)
- New unified `AnalyzeWithReader()` function that works with either source
- Existing `Analyze()` and `AnalyzeLocal()` still work (call the unified function)
- All existing code compiles without changes

## Dependencies
None, but recommended after Task 01 (NaisConfig methods).

## Related Code

**New file to create:** `internal/analyzer/reader.go`

```go
package analyzer

// RepositoryReader abstracts file reading from GitHub or local filesystem
type RepositoryReader interface {
    // GetFileContent reads a file's content. repoPath is the repo identifier
    // (full name for GitHub, directory path for local).
    GetFileContent(repoPath, filePath string) (string, error)
    
    // FileExists checks if a file exists in the repository
    FileExists(repoPath, filePath string) bool
}
```

**Verify `github.Client` already implements the interface** (`internal/github/client.go`):
```go
// Already exists - matches interface:
func (c *Client) GetFileContent(fullName, path string) (string, error)
func (c *Client) FileExists(fullName, path string) bool
```

**Verify `local.Reader` already implements the interface** (`internal/local/reader.go`):
```go
// Already exists - matches interface:
func (r *Reader) GetFileContent(repoPath, filePath string) (string, error)
func (r *Reader) FileExists(repoPath, filePath string) bool
```

**Modify `internal/analyzer/analyzer.go`:**

Add new unified function after imports:

```go
// Repository represents common repository metadata
type Repository struct {
    Name        string
    FullName    string
    Path        string // Only set for local repos
    Description string
    Language    string
}

// AnalyzeWithReader performs analysis using any RepositoryReader implementation
func AnalyzeWithReader(reader RepositoryReader, repo Repository) (*Analysis, error) {
    repoPath := repo.FullName
    if repo.Path != "" {
        repoPath = repo.Path
    }

    analysis := &Analysis{
        Repository: &github.Repository{
            Name:        repo.Name,
            FullName:    repo.FullName,
            Description: repo.Description,
            Language:    repo.Language,
        },
    }

    // Try to fetch README
    if readme, err := reader.GetFileContent(repoPath, "README.md"); err == nil {
        analysis.Readme = readme
    }

    // Try to fetch nais.yaml
    if naisConfig, err := fetchNaisConfigWithReader(reader, repoPath); err == nil {
        analysis.NaisConfig = naisConfig
    }

    // Try to extract dependencies
    if deps, err := fetchDependenciesWithReader(reader, repoPath, repo.Language); err == nil {
        analysis.Dependencies = deps
    }

    return analysis, nil
}

func fetchNaisConfigWithReader(reader RepositoryReader, repoPath string) (*NaisConfig, error) {
    paths := []string{
        "nais.yaml",
        ".nais/app.yaml",
        ".nais/nais.yaml",
        "nais/nais.yaml",
        ".nais/dev.yaml",
    }

    for _, path := range paths {
        if !reader.FileExists(repoPath, path) {
            continue
        }
        content, err := reader.GetFileContent(repoPath, path)
        if err != nil {
            continue
        }
        var config NaisConfig
        if err := yaml.Unmarshal([]byte(content), &config); err != nil {
            continue
        }
        return &config, nil
    }

    return nil, fmt.Errorf("no nais.yaml found")
}

func fetchDependenciesWithReader(reader RepositoryReader, repoPath, language string) ([]Dependency, error) {
    var deps []Dependency

    switch strings.ToLower(language) {
    case "kotlin", "java":
        if content, err := reader.GetFileContent(repoPath, "build.gradle.kts"); err == nil {
            deps = append(deps, parseGradleDeps(content)...)
        }
    case "go":
        if content, err := reader.GetFileContent(repoPath, "go.mod"); err == nil {
            deps = append(deps, parseGoModDeps(content)...)
        }
    }

    return deps, nil
}
```

**Keep existing `Analyze()` function for backwards compatibility:**
```go
// Analyze performs complete analysis of a repository (GitHub)
func Analyze(gh *github.Client, repo *github.Repository) (*Analysis, error) {
    return AnalyzeWithReader(gh, Repository{
        Name:        repo.Name,
        FullName:    repo.FullName,
        Description: repo.Description,
        Language:    repo.Language,
    })
}
```

## Verification

```bash
# Build should succeed
go build ./...

# Create interface compliance test
cat > internal/analyzer/reader_test.go << 'EOF'
package analyzer

import (
    "testing"
    
    "github.com/navikt/naisalyser/internal/github"
    "github.com/navikt/naisalyser/internal/local"
)

// Compile-time interface compliance checks
var _ RepositoryReader = (*github.Client)(nil)
var _ RepositoryReader = (*local.Reader)(nil)

func TestRepositoryReader_Interface(t *testing.T) {
    // If this compiles, the interfaces are satisfied
}
EOF

go test ./internal/analyzer/... -v
```

## Files Changed
- `internal/analyzer/reader.go` (new file, ~15 lines)
- `internal/analyzer/analyzer.go` (add ~80 lines)
- `internal/analyzer/reader_test.go` (new file, ~15 lines)
