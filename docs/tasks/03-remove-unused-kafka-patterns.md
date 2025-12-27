# Task 03: Remove Unused Kafka Patterns Variable

## Deliverable
Remove the unused `patterns` variable from `ExtractKafkaTopics` function.

## Context
In `internal/analyzer/kafka.go`, a `patterns` variable is declared but never used. The code explicitly suppresses the "unused variable" error with `_ = patterns`. This is dead code that should be removed per YAGNI (You Aren't Gonna Need It).

**Current code (kafka.go:22-54):**
```go
func ExtractKafkaTopics(reader *local.Reader, repoPath string) []KafkaTopic {
    var topics []KafkaTopic

    // Common patterns for Kafka topic definitions
    patterns := []string{
        "src/main/kotlin/**/kafka/Topics.kt",
        "src/main/kotlin/**/Topics.kt",
        "src/main/java/**/Topics.java",
    }

    // ... actual implementation uses knownPaths instead ...

    _ = patterns // Reserved for future use  ← DEAD CODE

    return topics
}
```

## Key Decisions and Principles
- Remove code that isn't used
- If the patterns are needed later, they can be added back from git history
- Don't leave "reserved for future use" code in the codebase

## Delivers
Cleaner code without unused variables and suppression statements.

## Acceptance Criteria
- `patterns` variable declaration removed (lines 23-28)
- `_ = patterns` suppression statement removed (line 54)
- Function still works correctly
- Code compiles without warnings

## Dependencies
None.

## Related Code

**File to modify:** `internal/analyzer/kafka.go`

**Current function (lines 19-57):**
```go
func ExtractKafkaTopics(reader *local.Reader, repoPath string) []KafkaTopic {
    var topics []KafkaTopic

    // Common patterns for Kafka topic definitions
    patterns := []string{
        "src/main/kotlin/**/kafka/Topics.kt",
        "src/main/kotlin/**/Topics.kt",
        "src/main/java/**/Topics.java",
    }

    // Try known paths first
    knownPaths := []string{
        "src/main/kotlin/no/nav/helse/kafka/Topics.kt",
    }

    for _, path := range knownPaths {
        content, err := reader.GetFileContent(repoPath, path)
        if err != nil {
            continue
        }
        topics = append(topics, parseKotlinTopics(content, path)...)
    }

    // Also search in any file containing "Topic" in the kafka directory
    kafkaDir := filepath.Join(repoPath, "src/main/kotlin")
    entries, _ := filepath.Glob(filepath.Join(kafkaDir, "**/kafka/*Topic*.kt"))
    for _, fullPath := range entries {
        relPath, _ := filepath.Rel(repoPath, fullPath)
        content, err := reader.GetFileContent(repoPath, relPath)
        if err != nil {
            continue
        }
        topics = append(topics, parseKotlinTopics(content, relPath)...)
    }

    _ = patterns // Reserved for future use

    return topics
}
```

**Replace with:**
```go
func ExtractKafkaTopics(reader *local.Reader, repoPath string) []KafkaTopic {
    var topics []KafkaTopic

    // Try known paths first
    knownPaths := []string{
        "src/main/kotlin/no/nav/helse/kafka/Topics.kt",
    }

    for _, path := range knownPaths {
        content, err := reader.GetFileContent(repoPath, path)
        if err != nil {
            continue
        }
        topics = append(topics, parseKotlinTopics(content, path)...)
    }

    // Also search in any file containing "Topic" in the kafka directory
    kafkaDir := filepath.Join(repoPath, "src/main/kotlin")
    entries, _ := filepath.Glob(filepath.Join(kafkaDir, "**/kafka/*Topic*.kt"))
    for _, fullPath := range entries {
        relPath, _ := filepath.Rel(repoPath, fullPath)
        content, err := reader.GetFileContent(repoPath, relPath)
        if err != nil {
            continue
        }
        topics = append(topics, parseKotlinTopics(content, relPath)...)
    }

    return topics
}
```

## Verification

```bash
# Build should succeed with no warnings
go build ./...

# Verify no unused variable warnings
go vet ./...

# Run existing functionality (should still work)
# If you have test repos:
# ./naisalyser local ./repos/some-repo -v
```

## Files Changed
- `internal/analyzer/kafka.go` (remove ~7 lines)
