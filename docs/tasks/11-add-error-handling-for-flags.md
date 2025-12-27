# Task 11: Add Error Handling for Cobra Flag Retrieval

## Deliverable
Replace ignored errors from `cmd.Flags().GetString()` with proper error handling.

## Context
Throughout the cmd package, flag values are retrieved with errors being discarded:

```go
outputDir, _ := cmd.Flags().GetString("output")
verbose, _ := cmd.Flags().GetBool("verbose")
```

While these errors are unlikely (Cobra validates flags before RunE), ignoring them is bad practice and could hide bugs during development.

## Key Decisions and Principles
- Check errors explicitly and return if non-nil
- Use helper function to reduce repetition
- Fail fast on configuration errors

## Delivers
Proper error handling for all flag retrieval calls.

## Acceptance Criteria
- All `cmd.Flags().Get*()` calls check and handle errors
- Clear error messages indicating which flag failed
- No behavior change for normal operation
- Code compiles and runs correctly

## Dependencies
None.

## Related Code

**Files to modify:** All files in `cmd/` package

### Option A: Inline Error Handling (Simple)

**Modify `cmd/analyze.go` (lines 28-31):**

Current:
```go
func runAnalyze(cmd *cobra.Command, args []string) error {
    repo := args[0]
    outputDir, _ := cmd.Flags().GetString("output")
    verbose, _ := cmd.Flags().GetBool("verbose")
```

New:
```go
func runAnalyze(cmd *cobra.Command, args []string) error {
    repo := args[0]
    
    outputDir, err := cmd.Flags().GetString("output")
    if err != nil {
        return fmt.Errorf("failed to get output flag: %w", err)
    }
    
    verbose, err := cmd.Flags().GetBool("verbose")
    if err != nil {
        return fmt.Errorf("failed to get verbose flag: %w", err)
    }
```

**Apply same pattern to:**
- `cmd/batch.go` (lines 36-41)
- `cmd/local.go` (lines 42-45 and 73-76)
- `cmd/graph.go` (lines 36-39)

### Option B: Helper Function (DRY)

**Create `cmd/flags.go`:**

```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

// mustGetString retrieves a string flag or returns an error
func mustGetString(cmd *cobra.Command, name string) (string, error) {
    val, err := cmd.Flags().GetString(name)
    if err != nil {
        return "", fmt.Errorf("failed to get flag %q: %w", name, err)
    }
    return val, nil
}

// mustGetBool retrieves a bool flag or returns an error
func mustGetBool(cmd *cobra.Command, name string) (bool, error) {
    val, err := cmd.Flags().GetBool(name)
    if err != nil {
        return false, fmt.Errorf("failed to get flag %q: %w", name, err)
    }
    return val, nil
}
```

**Then update usages:**

```go
func runAnalyze(cmd *cobra.Command, args []string) error {
    repo := args[0]
    
    outputDir, err := mustGetString(cmd, "output")
    if err != nil {
        return err
    }
    
    verbose, err := mustGetBool(cmd, "verbose")
    if err != nil {
        return err
    }
    // ... rest of function
}
```

## Full List of Changes

**cmd/analyze.go:28-31:**
```go
// Before
outputDir, _ := cmd.Flags().GetString("output")
verbose, _ := cmd.Flags().GetBool("verbose")

// After
outputDir, err := cmd.Flags().GetString("output")
if err != nil {
    return fmt.Errorf("failed to get output flag: %w", err)
}
verbose, err := cmd.Flags().GetBool("verbose")
if err != nil {
    return fmt.Errorf("failed to get verbose flag: %w", err)
}
```

**cmd/batch.go:36-41:**
```go
// Before
topic, _ := cmd.Flags().GetString("topic")
configFile, _ := cmd.Flags().GetString("config")
org, _ := cmd.Flags().GetString("org")
outputDir, _ := cmd.Flags().GetString("output")
verbose, _ := cmd.Flags().GetBool("verbose")

// After - add error checks for each
```

**cmd/local.go:42-45:**
```go
// Before
outputDir, _ := cmd.Flags().GetString("output")
verbose, _ := cmd.Flags().GetBool("verbose")

// After - add error checks
```

**cmd/local.go:73-76 (runLocalBatch):**
```go
// Same pattern
```

**cmd/graph.go:36-39:**
```go
// Same pattern
```

## Verification

```bash
# Build should succeed
go build ./...

# Test that normal operation works
./naisalyser --help
./naisalyser analyze --help
./naisalyser local ./repos/some-repo -v

# Verify error handling works (this would require mocking, so just verify compilation)
go vet ./cmd/...
```

## Files Changed
- `cmd/analyze.go` (modify ~6 lines)
- `cmd/batch.go` (modify ~12 lines)
- `cmd/local.go` (modify ~12 lines)
- `cmd/graph.go` (modify ~8 lines)
- `cmd/flags.go` (optional new file, ~20 lines if using helper approach)
