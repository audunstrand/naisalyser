# Task 10: Add Configurable Organization Parameter

## Deliverable
Remove hardcoded "navikt" organization and make it configurable.

## Context
The `local.Reader.ReadRepository` function hardcodes `"navikt"` as the organization when constructing the FullName:

```go
// internal/local/reader.go:55
FullName: fmt.Sprintf("navikt/%s", name), // Assume navikt org
```

This limits the tool to only navikt repositories. Users analyzing repos from other organizations get incorrect FullName values.

## Key Decisions and Principles
- Add `--org` flag to root command (available to all subcommands)
- Default to "navikt" for backwards compatibility
- Pass org through to Reader where needed
- Don't break existing API if possible

## Delivers
Configurable organization parameter that works with local commands.

## Acceptance Criteria
- Root command has `--org` persistent flag with default "navikt"
- `local` command passes org to reader
- `local-batch` command passes org to reader
- `Reader.ReadRepository` accepts org parameter
- FullName is constructed as `org/name`
- Existing behavior unchanged when `--org` not specified

## Dependencies
None.

## Related Code

**File to modify:** `cmd/root.go`

**Add org flag in init() (around line 21):**

```go
func init() {
    rootCmd.PersistentFlags().StringP("output", "o", "./docs/generated", "Output directory for generated docs")
    rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
    rootCmd.PersistentFlags().String("org", "navikt", "GitHub organization for repository URLs")
}
```

**File to modify:** `internal/local/reader.go`

**Modify ReadRepository signature (line 31):**

Current:
```go
func (r *Reader) ReadRepository(path string) (*Repository, error) {
```

New:
```go
func (r *Reader) ReadRepository(path string, org string) (*Repository, error) {
```

**Update FullName construction (around line 54):**

Current:
```go
FullName: fmt.Sprintf("navikt/%s", name), // Assume navikt org
```

New:
```go
FullName: fmt.Sprintf("%s/%s", org, name),
```

**File to modify:** `cmd/local.go`

**Update runLocal function (around line 42):**

```go
func runLocal(cmd *cobra.Command, args []string) error {
    repoPath := args[0]
    outputDir, _ := cmd.Flags().GetString("output")
    verbose, _ := cmd.Flags().GetBool("verbose")
    org, _ := cmd.Flags().GetString("org")  // ADD THIS

    if verbose {
        fmt.Printf("Analyzing local repository: %s\n", repoPath)
    }

    reader := local.NewReader(verbose)
    repoData, err := reader.ReadRepository(repoPath, org)  // PASS ORG
    if err != nil {
        return fmt.Errorf("failed to read repository: %w", err)
    }
    // ... rest unchanged
}
```

**Update runLocalBatch function (around line 73):**

```go
func runLocalBatch(cmd *cobra.Command, args []string) error {
    reposDir := args[0]
    outputDir, _ := cmd.Flags().GetString("output")
    verbose, _ := cmd.Flags().GetBool("verbose")
    org, _ := cmd.Flags().GetString("org")  // ADD THIS

    // ... existing code ...

    for i, repoPath := range repos {
        // ...
        repoData, err := reader.ReadRepository(repoPath, org)  // PASS ORG
        // ...
    }
}
```

**File to modify:** `cmd/graph.go`

**Update runGraph function (around line 36):**

```go
func runGraph(cmd *cobra.Command, args []string) error {
    reposDir := args[0]
    outputDir, _ := cmd.Flags().GetString("output")
    verbose, _ := cmd.Flags().GetBool("verbose")
    org, _ := cmd.Flags().GetString("org")  // ADD THIS

    // ... existing code ...

    for _, repoPath := range repos {
        // ...
        repoData, err := reader.ReadRepository(repoPath, org)  // PASS ORG
        // ...
    }
}
```

## Verification

```bash
# Build should succeed
go build ./...

# Test default behavior (should still use navikt)
./naisalyser local ./repos/some-repo -v 2>&1 | grep "navikt"

# Test custom org
./naisalyser local ./repos/some-repo --org myorg -v 2>&1 | grep "myorg"

# Verify help shows the flag
./naisalyser --help | grep org

# Create unit test
cat > internal/local/reader_test.go << 'EOF'
package local

import (
    "os"
    "path/filepath"
    "testing"
)

func TestReadRepository_CustomOrg(t *testing.T) {
    // Create temp directory
    tmpDir, err := os.MkdirTemp("", "test-repo")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(tmpDir)

    reader := NewReader(false)
    repo, err := reader.ReadRepository(tmpDir, "custom-org")
    if err != nil {
        t.Fatal(err)
    }

    expectedFullName := "custom-org/" + filepath.Base(tmpDir)
    if repo.FullName != expectedFullName {
        t.Errorf("expected FullName=%s, got %s", expectedFullName, repo.FullName)
    }
}

func TestReadRepository_DefaultOrg(t *testing.T) {
    tmpDir, err := os.MkdirTemp("", "test-repo")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(tmpDir)

    reader := NewReader(false)
    repo, err := reader.ReadRepository(tmpDir, "navikt")
    if err != nil {
        t.Fatal(err)
    }

    expectedFullName := "navikt/" + filepath.Base(tmpDir)
    if repo.FullName != expectedFullName {
        t.Errorf("expected FullName=%s, got %s", expectedFullName, repo.FullName)
    }
}
EOF

go test ./internal/local/... -v
```

## Files Changed
- `cmd/root.go` (add 1 line)
- `cmd/local.go` (modify ~6 lines)
- `cmd/graph.go` (modify ~4 lines)
- `internal/local/reader.go` (modify ~3 lines)
- `internal/local/reader_test.go` (new file, ~45 lines)
