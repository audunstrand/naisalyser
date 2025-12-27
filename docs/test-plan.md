# Test Coverage Implementation Plan

## Executive Summary
Current test coverage: **1.7%**
Target test coverage: **≥80%**

This plan breaks down test implementation into 10 independent, parallelizable tasks that can be worked on simultaneously.

---

## Task Overview

| Task | Component | Priority | Est. Lines | Tests | Dependencies |
|------|-----------|----------|-----------|-------|--------------|
| T1 | Parser Tests | HIGH | 300+ | 40+ | None |
| T2 | Handlebars Preprocessing | HIGH | 150 | 15+ | None |
| T3 | Graph Builder Tests | HIGH | 250+ | 30+ | Task T2 (optional) |
| T4 | Command Handler Tests | MEDIUM | 400+ | 50+ | None |
| T5 | Local Reader Tests | HIGH | 200+ | 25+ | None |
| T6 | GitHub Client Tests | MEDIUM | 200+ | 20+ | Mocking |
| T7 | Graph Mermaid Output | MEDIUM | 150 | 15+ | None |
| T8 | Output Generation Tests | MEDIUM | 300+ | 30+ | None |
| T9 | NaisConfig Methods | COMPLETE | 42 | 3 | ✅ DONE |
| T10 | Integration Tests | LOW | 200+ | 15+ | All others |

---

## Detailed Task Breakdown

### Task T1: Parser Tests
**File:** `internal/analyzer/parsers_test.go`
**Coverage:** parseGradleDeps, parseNpmDeps, parseMavenDeps, parseGoModDeps

**Test Cases:**
- ParseGradleDeps:
  - Empty content
  - Single dependency (implementation)
  - Multiple dependencies mixed (api + implementation)
  - Malformed gradle syntax
  - Dependencies with versions
  
- ParseNpmDeps:
  - Valid package.json with dependencies
  - Valid package.json with devDependencies
  - Both dependencies and devDependencies
  - Malformed JSON (returns nil)
  - Empty dependencies object
  - Complex semver versions
  
- ParseMavenDeps:
  - Valid pom.xml with dependencies
  - Filters test-scoped dependencies
  - Missing version (optional)
  - Malformed XML (returns nil)
  - Multiple dependencies with groupId:artifactId format
  
- ParseGoModDeps:
  - Valid go.mod require block
  - Multiple dependencies with versions
  - Malformed go.mod
  - Empty require block

**Estimated Tests:** 40+
**Can start:** Immediately (no dependencies)

---

### Task T2: Handlebars Preprocessing Tests
**File:** `internal/analyzer/handlebars_test.go`
**Coverage:** preprocessHandlebars, removeEachBlocks, replaceTemplateVars, replaceLineTemplateVars

**Test Cases:**
- removeEachBlocks:
  - Single {{#each}} block
  - Multiple nested {{#each}} blocks
  - Block at start of content
  - Block at end of content
  - No {{#each}} blocks
  
- replaceTemplateVars:
  - Single variable per line
  - Multiple variables per line
  - No variables
  - Edge cases with nested braces
  
- replaceLineTemplateVars:
  - Replica/min/max context (uses "1")
  - Port/timeout/delay context (uses "1")
  - Other contexts (uses "placeholder")
  - Multiple variables with different contexts
  - Mixed case insensitive detection

**Estimated Tests:** 15+
**Can start:** Immediately (no dependencies)

---

### Task T3: Graph Builder Tests
**File:** `internal/graph/builder_test.go`
**Coverage:** NewBuilder, AddFromLocalAnalysis, addNode, addEdge, Build, edge deduplication with edgeKey

**Test Cases:**
- Builder initialization:
  - NewBuilder creates empty graph
  - Nodes map initialized
  - Edges slice initialized
  - edgeSet map initialized
  
- Node operations:
  - addNode adds new node
  - addNode doesn't duplicate nodes with same ID
  - addNode preserves node properties (type, stateful, namespace)
  
- Edge operations:
  - addEdge adds new edge
  - addEdge rejects duplicates (O(1) using edgeSet)
  - Edge deduplication works with type-safe EdgeType constants
  - Edges maintain insertion order
  
- AddFromLocalAnalysis:
  - Handles nil analysis
  - Handles nil repository
  - Adds app node with correct properties
  - Adds inbound app nodes and edges
  - Adds outbound app nodes and edges
  - Adds external nodes and edges
  - Adds Kafka topic nodes and edges
  - Handles missing NaisConfig gracefully
  - Uses IsStateful() from NaisConfig
  
- Build:
  - Returns Graph with all nodes
  - Returns Graph with all edges
  - Nodes are in consistent order

**Estimated Tests:** 30+
**Can start:** Immediately (optional T2 dependency for integration)

---

### Task T4: Command Handler Tests
**File:** `cmd/commands_test.go`
**Coverage:** runAnalyze, runBatch, runLocal, runLocalBatch, runGraph

**Test Cases:**
- runAnalyze:
  - Requires valid repo argument
  - Flag parsing works
  - Error handling when repo not found
  - Success path with mock data
  
- runBatch:
  - Topic flag required or config file
  - Both flags specified (error)
  - Config file parsing
  - Repository iteration
  - Error handling for failed repos
  
- runLocal:
  - Single repository analysis
  - Org flag defaults to "navikt"
  - Custom org flag works
  - Output directory creation
  - Error handling for invalid paths
  
- runLocalBatch:
  - Directory iteration
  - Org flag applied to all repos
  - Skips hidden directories
  - Batch error handling
  
- runGraph:
  - Graph building from multiple repos
  - Output file creation (JSON + Markdown)
  - Mermaid diagram generation

**Estimated Tests:** 50+
**Can start:** Immediately (uses mock filesystem)

---

### Task T5: Local Reader Tests
**File:** `internal/local/reader_test.go`
**Coverage:** NewReader, ReadRepository, GetFileContent, FileExists, detectLanguage, readFirstLineOfReadme

**Test Cases:**
- ReadRepository:
  - Valid directory path
  - Org parameter in FullName construction
  - Invalid path error handling
  - Non-directory path error handling
  - Language detection (Go, Java, Node, Python, Kotlin)
  - README description extraction
  
- GetFileContent:
  - Existing file returns content
  - Non-existent file returns error
  - Directory path error handling
  
- FileExists:
  - Returns true for existing files
  - Returns false for non-existent files
  
- detectLanguage:
  - Detects Go (go.mod, go.sum)
  - Detects Java (pom.xml, build.gradle.kts)
  - Detects Kotlin (*.kt files)
  - Detects JavaScript/TypeScript (package.json)
  - Detects Python (requirements.txt, setup.py)
  - Returns empty string for unknown
  
- readFirstLineOfReadme:
  - Reads first line of README.md
  - Handles missing README
  - Handles empty README

**Estimated Tests:** 25+
**Can start:** Immediately (uses temp directories)

---

### Task T6: GitHub Client Tests
**File:** `internal/github/client_test.go`
**Coverage:** NewClient, GetRepository, GetRepositoriesByTopic, GetFileContent, FileExists, GetDirectoryContents

**Test Cases:**
- Client initialization:
  - NewClient creates client
  - Verbose flag stored
  
- GetRepository (mock):
  - Valid repo returns data
  - Invalid repo returns error
  - Parsing repository info
  
- GetRepositoriesByTopic (mock):
  - Valid topic returns repos
  - Empty result handling
  - Pagination handling
  
- GetFileContent (mock):
  - File exists returns content
  - File not found returns error
  - Base64 decoding works
  
- FileExists (mock):
  - Returns true for existing files
  - Returns false for missing files
  
- GetDirectoryContents (mock):
  - Lists directory entries
  - Handles missing directory

**Estimated Tests:** 20+
**Can start:** After setting up test mocks

---

### Task T7: Graph Mermaid Output Tests
**File:** `internal/graph/mermaid_test.go`
**Coverage:** ToMermaid, sanitizeMermaidID

**Test Cases:**
- sanitizeMermaidID:
  - Replaces hyphens with underscores
  - Replaces dots with underscores
  - Replaces slashes with underscores
  - Replaces colons with underscores
  - Handles empty strings
  
- ToMermaid:
  - Empty graph
  - Single node (app, kafka, external)
  - Single edge (call, kafka, external)
  - Mixed nodes and edges
  - Stateful vs stateless app styling
  - Correct Mermaid syntax output
  - Node IDs are sanitized
  - Uses type-safe constants (NodeTypeApp, EdgeTypeCall, etc.)

**Estimated Tests:** 15+
**Can start:** Immediately (no dependencies)

---

### Task T8: Output Generation Tests
**File:** `internal/output/markdown_test.go`
**Coverage:** GenerateMarkdown, toViewModel, generateOverviewVM, generateDependenciesVM, generateNaisConfigVM, executeTemplate

**Test Cases:**
- GenerateMarkdown:
  - Creates output directory
  - Writes README.md
  - Handles invalid output path
  
- toViewModel:
  - Converts Analysis to ViewModel
  - Handles nil NaisConfig
  - Handles nil Dependencies
  
- generateOverviewVM:
  - Generates overview view model
  - Includes repository metadata
  - Handles missing data gracefully
  
- generateDependenciesVM:
  - Lists dependencies correctly
  - Groups by type
  - Handles empty dependencies
  
- generateNaisConfigVM:
  - Generates config view model
  - Shows access policies
  - Shows Kafka config
  - Shows GCP resources
  
- executeTemplate:
  - Renders template correctly
  - Handles template errors
  - Output is valid markdown

**Estimated Tests:** 30+
**Can start:** After parsers (T1)

---

### Task T9: NaisConfig Methods Tests
**Status:** ✅ COMPLETE
**File:** `internal/analyzer/naisconfig_test.go`
**Current Tests:** 3

---

### Task T10: Integration Tests
**File:** `internal/integration_test.go`
**Coverage:** End-to-end workflows

**Test Cases:**
- Full local analysis workflow:
  - Create temp repo directory
  - Add NAIS config file
  - Analyze with local reader
  - Verify output
  
- Full graph building workflow:
  - Create multiple temp repos
  - Build graph from batch
  - Verify nodes and edges
  - Generate Mermaid output
  
- Parser integration:
  - Real build.gradle.kts file
  - Real pom.xml file
  - Real package.json file
  - Real go.mod file

**Estimated Tests:** 15+
**Can start:** After T1, T3, T5

---

## Implementation Sequence

### Phase 1 (Parallel, no dependencies):
- T1: Parser Tests
- T2: Handlebars Preprocessing Tests
- T5: Local Reader Tests
- T7: Graph Mermaid Output Tests

### Phase 2 (After Phase 1):
- T3: Graph Builder Tests (uses T1 mock data)
- T4: Command Handler Tests
- T6: GitHub Client Tests (using mocks)

### Phase 3 (After Phase 2):
- T8: Output Generation Tests
- T10: Integration Tests

---

## Success Criteria

✅ All tests pass
✅ Coverage increases from 1.7% → ≥80%
✅ No breaking changes to existing code
✅ Test execution time < 10 seconds total
✅ Clear test naming following convention: `Test<Function><Scenario>`

---

## Test File Structure Template

```go
package <package>

import (
    "testing"
    // test helpers
)

func Test<Function><Scenario>(t *testing.T) {
    // Arrange
    
    // Act
    
    // Assert
}
```

---

## Parallelization Strategy

Tasks can be executed in parallel by different agents using `/delegate`:

```bash
/delegate Test Task T1: Parser Tests (40+ test cases)
/delegate Test Task T2: Handlebars Preprocessing (15+ test cases)
/delegate Test Task T5: Local Reader Tests (25+ test cases)
/delegate Test Task T7: Graph Mermaid Output (15+ test cases)
```

Then after Phase 1 completes:

```bash
/delegate Test Task T3: Graph Builder Tests
/delegate Test Task T4: Command Handler Tests
/delegate Test Task T6: GitHub Client Tests
```

---

## File Count & Estimate

- **Test Files to Create:** 10
- **Total Test Code Lines:** 1,900+
- **Total Test Cases:** 280+
- **Expected Coverage:** 85%+
- **Estimated Time (Sequential):** 8-10 hours
- **Estimated Time (Parallel):** 2-3 hours

