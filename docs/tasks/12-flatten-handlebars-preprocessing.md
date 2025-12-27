# Task 12: Flatten Handlebars Preprocessing Function

## Deliverable
Refactor `preprocessHandlebars` function to reduce nesting from 4 levels to 2.

## Context
The `preprocessHandlebars` function in `local_analyzer.go` has deeply nested loops making it hard to read and test. The function replaces Handlebars template syntax (like `{{variable}}`) with placeholder values so YAML can be parsed.

**Current structure (simplified):**
```go
func preprocessHandlebars(content string) string {
    // ... remove {{#each}} blocks ...
    
    lines := strings.Split(result, "\n")
    for i, line := range lines {           // Level 1
        for strings.Contains(line, "{{") { // Level 2
            start := strings.Index(...)
            if start == -1 {                // Level 3
                break
            }
            end := strings.Index(...)
            if end == -1 {                  // Level 3
                break
            }
            // determine replacement       // Level 4 logic here
            line = line[:start] + ...
        }
        lines[i] = line
    }
    return strings.Join(lines, "\n")
}
```

## Key Decisions and Principles
- Extract inner logic into separate functions
- Each function does one thing
- Use early returns to reduce nesting
- Function names reveal intent

## Delivers
More readable `preprocessHandlebars` with extracted helper functions.

## Acceptance Criteria
- New `removeEachBlocks(content string) string` function
- New `replaceTemplateVars(content string) string` function
- New `replaceLineTemplateVars(line string) string` function
- Original function delegates to these
- Behavior unchanged (same output for same input)
- Maximum nesting reduced to 2 levels

## Dependencies
None.

## Related Code

**File to modify:** `internal/analyzer/local_analyzer.go`

**Current function (lines 144-201):**
```go
func preprocessHandlebars(content string) string {
    result := content

    // Replace {{#each ...}} ... {{/each}} blocks with empty
    for strings.Contains(result, "{{#each") {
        start := strings.Index(result, "{{#each")
        if start == -1 {
            break
        }
        end := strings.Index(result[start:], "{{/each}}")
        if end == -1 {
            break
        }
        lineStart := strings.LastIndex(result[:start], "\n") + 1
        lineEnd := start + end + len("{{/each}}")
        if nextNewline := strings.Index(result[lineEnd:], "\n"); nextNewline != -1 {
            lineEnd += nextNewline + 1
        }
        result = result[:lineStart] + result[lineEnd:]
    }

    lines := strings.Split(result, "\n")
    for i, line := range lines {
        for strings.Contains(line, "{{") {
            start := strings.Index(line, "{{")
            if start == -1 {
                break
            }
            end := strings.Index(line[start:], "}}")
            if end == -1 {
                break
            }

            replacement := "placeholder"
            lineLower := strings.ToLower(line)
            if strings.Contains(lineLower, "replica") ||
                strings.Contains(lineLower, "min:") ||
                strings.Contains(lineLower, "max:") ||
                strings.Contains(lineLower, "port") ||
                strings.Contains(lineLower, "timeout") ||
                strings.Contains(lineLower, "delay") {
                replacement = "1"
            }

            line = line[:start] + replacement + line[start+end+2:]
        }
        lines[i] = line
    }
    result = strings.Join(lines, "\n")

    return result
}
```

**Replace with:**

```go
// preprocessHandlebars replaces handlebars template syntax with placeholder values
func preprocessHandlebars(content string) string {
    content = removeEachBlocks(content)
    content = replaceTemplateVars(content)
    return content
}

// removeEachBlocks removes {{#each ...}} ... {{/each}} blocks entirely
func removeEachBlocks(content string) string {
    result := content
    for strings.Contains(result, "{{#each") {
        start := strings.Index(result, "{{#each")
        if start == -1 {
            break
        }
        end := strings.Index(result[start:], "{{/each}}")
        if end == -1 {
            break
        }
        
        lineStart := strings.LastIndex(result[:start], "\n") + 1
        lineEnd := start + end + len("{{/each}}")
        if nextNewline := strings.Index(result[lineEnd:], "\n"); nextNewline != -1 {
            lineEnd += nextNewline + 1
        }
        result = result[:lineStart] + result[lineEnd:]
    }
    return result
}

// replaceTemplateVars replaces {{variable}} with appropriate placeholders
func replaceTemplateVars(content string) string {
    lines := strings.Split(content, "\n")
    for i, line := range lines {
        lines[i] = replaceLineTemplateVars(line)
    }
    return strings.Join(lines, "\n")
}

// replaceLineTemplateVars replaces all template vars in a single line
func replaceLineTemplateVars(line string) string {
    for strings.Contains(line, "{{") {
        start := strings.Index(line, "{{")
        if start == -1 {
            break
        }
        end := strings.Index(line[start:], "}}")
        if end == -1 {
            break
        }
        
        replacement := getPlaceholderForContext(line)
        line = line[:start] + replacement + line[start+end+2:]
    }
    return line
}

// getPlaceholderForContext returns "1" for numeric contexts, "placeholder" otherwise
func getPlaceholderForContext(line string) string {
    lineLower := strings.ToLower(line)
    numericContexts := []string{"replica", "min:", "max:", "port", "timeout", "delay"}
    
    for _, ctx := range numericContexts {
        if strings.Contains(lineLower, ctx) {
            return "1"
        }
    }
    return "placeholder"
}
```

## Verification

```bash
# Build should succeed
go build ./...

# Create test to verify behavior unchanged
cat > internal/analyzer/handlebars_test.go << 'EOF'
package analyzer

import "testing"

func TestPreprocessHandlebars_Simple(t *testing.T) {
    input := "name: {{appName}}"
    expected := "name: placeholder"
    
    result := preprocessHandlebars(input)
    if result != expected {
        t.Errorf("expected %q, got %q", expected, result)
    }
}

func TestPreprocessHandlebars_NumericContext(t *testing.T) {
    input := "replicas: {{replicaCount}}"
    expected := "replicas: 1"
    
    result := preprocessHandlebars(input)
    if result != expected {
        t.Errorf("expected %q, got %q", expected, result)
    }
}

func TestPreprocessHandlebars_EachBlock(t *testing.T) {
    input := `before
{{#each items}}
  - {{this}}
{{/each}}
after`
    
    result := preprocessHandlebars(input)
    if strings.Contains(result, "{{#each") {
        t.Error("expected {{#each}} block to be removed")
    }
    if !strings.Contains(result, "before") || !strings.Contains(result, "after") {
        t.Error("expected surrounding content to remain")
    }
}

func TestPreprocessHandlebars_MultipleVars(t *testing.T) {
    input := "{{a}} and {{b}}"
    expected := "placeholder and placeholder"
    
    result := preprocessHandlebars(input)
    if result != expected {
        t.Errorf("expected %q, got %q", expected, result)
    }
}
EOF

go test ./internal/analyzer/... -v -run TestPreprocessHandlebars
```

## Files Changed
- `internal/analyzer/local_analyzer.go` (refactor ~60 lines into ~70 lines with better structure)
- `internal/analyzer/handlebars_test.go` (new file, ~50 lines)
