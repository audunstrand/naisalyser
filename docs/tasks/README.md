# Test Implementation Tasks - Complete Plan

## 📋 Overview

This directory contains the complete test implementation plan to increase code coverage from **1.7% → 85%+**.

All tasks are designed to be **executed in parallel** using `/delegate` commands.

## 🗂️ Files

### Master Plan
- **[../test-plan.md](../test-plan.md)** - Overall strategy, timeline, success criteria

### Individual Task Specifications
Each file details a specific testing phase with:
- Test cases to implement
- Test data fixtures needed  
- Edge cases to cover
- Success criteria

#### Phase 1: Independent Tasks (Run in Parallel)
1. **[T1-parser-tests.md](T1-parser-tests.md)** - Dependency parsers (40+ tests)
   - `parseGradleDeps`, `parseNpmDeps`, `parseMavenDeps`, `parseGoModDeps`
   - File: `internal/analyzer/parsers_test.go`

2. **[T2-handlebars-tests.md](T2-handlebars-tests.md)** - Template preprocessing (15+ tests)
   - `preprocessHandlebars`, `removeEachBlocks`, `replaceTemplateVars`
   - File: `internal/analyzer/handlebars_test.go`

3. **[T5-local-reader-tests.md](T5-local-reader-tests.md)** - File I/O operations (25+ tests)
   - `ReadRepository`, `GetFileContent`, `FileExists`, `detectLanguage`
   - File: `internal/local/reader_test.go`

4. **[T7-mermaid-tests.md](T7-mermaid-tests.md)** - Graph visualization (15+ tests)
   - `ToMermaid`, `sanitizeMermaidID`
   - File: `internal/graph/mermaid_test.go`

#### Phase 2: Medium Dependencies (Start After Phase 1)
5. **[T3-graph-builder-tests.md](T3-graph-builder-tests.md)** - Graph construction (30+ tests)
   - `NewBuilder`, `AddFromLocalAnalysis`, edge deduplication
   - File: `internal/graph/builder_test.go`
   - Optional dependency: T2

6. **[T4-command-tests.md](T4-command-tests.md)** - CLI handlers (50+ tests)
   - `runAnalyze`, `runBatch`, `runLocal`, `runLocalBatch`, `runGraph`
   - Files: `cmd/commands_test.go`, `cmd/flags_test.go`

7. **[T6-github-client-tests.md](T6-github-client-tests.md)** - GitHub API integration (20+ tests)
   - `GetRepository`, `GetRepositoriesByTopic`, `GetFileContent`
   - File: `internal/github/client_test.go`
   - Note: All tests use mocks, no real API calls

#### Phase 3: Integration (Start After Phase 2)
8. **[T8-output-tests.md](T8-output-tests.md)** - Documentation generation (30+ tests)
   - `GenerateMarkdown`, `toViewModel`, `executeTemplate`
   - File: `internal/output/markdown_test.go`

9. **[T10-integration-tests.md](T10-integration-tests.md)** - End-to-end workflows (18+ tests)
   - Complete analysis pipelines
   - Real file parsing
   - Multi-repo scenarios
   - File: `internal/integration_test.go`

10. **T9 - COMPLETE** ✅
    - NaisConfig behavior methods (3 tests)
    - File: `internal/analyzer/naisconfig_test.go`
    - Status: Already passing

## 📊 Statistics

| Metric | Value |
|--------|-------|
| Total Test Cases | 280+ |
| Total Test Code | 1,900+ lines |
| Test Files | 10 |
| Coverage Target | 85%+ |
| Estimated Time (Parallel) | 2-3 hours |
| Estimated Time (Sequential) | 8-10 hours |

## 🚀 Quick Start

### Step 1: Review the Plan
```bash
cat docs/test-plan.md
```

### Step 2: Start Phase 1 (4 parallel tasks)
Each can be delegated to a separate agent:
```bash
# In real /delegate usage (these are placeholders)
/delegate Implement Task T1: Parser Tests (40+ test cases)
/delegate Implement Task T2: Handlebars Preprocessing (15+ test cases)
/delegate Implement Task T5: Local Reader Tests (25+ test cases)
/delegate Implement Task T7: Graph Mermaid Output (15+ test cases)
```

### Step 3: Monitor Progress
```bash
cd docs
# Check coverage as you go
go test ./internal/analyzer/... -cover
go test ./... -cover
```

### Step 4: Complete Remaining Phases
After Phase 1 completes, start Phase 2 (T3, T4, T6), then Phase 3 (T8, T10)

## ✅ Success Checklist

- [ ] Phase 1 tasks completed (T1, T2, T5, T7)
- [ ] Phase 2 tasks completed (T3, T4, T6)  
- [ ] Phase 3 tasks completed (T8, T10)
- [ ] All 280+ tests passing
- [ ] Coverage ≥ 85%
- [ ] No regressions in existing code
- [ ] Test execution time < 10 seconds

## 📈 Coverage Targets by Package

```
internal/analyzer   90%+  (core analysis logic)
internal/graph      85%+  (dependency graph)
internal/local      90%+  (file I/O)
internal/github     80%+  (API client - mocked)
internal/output     80%+  (documentation)
cmd                 75%+  (CLI handlers)
───────────────────────────
Overall             85%+
```

## 🔍 Key Testing Principles

1. **Unit Tests** - Test individual functions in isolation
2. **Integration Tests** - Verify complete workflows
3. **Edge Cases** - Test boundary conditions and error scenarios
4. **Mock Data** - Use realistic test fixtures
5. **No External Dependencies** - All tests are self-contained
6. **Performance** - Verify operations complete in reasonable time

## 📝 Test Data Management

Each task includes:
- Minimal test fixtures (files needed for tests)
- Mock data structures (realistic example data)
- Edge case scenarios (unusual but valid inputs)

All test fixtures are stored in `internal/*/testdata/` directories.

## 🔗 Related Documentation

- [Main Test Plan](../test-plan.md) - Strategic overview
- [Project Tasks](../) - All implementation tasks
- [Completed Tasks](../../docs/tasks/) - Previously completed tasks

---

**Last Updated:** December 27, 2025
**Status:** Ready for parallel execution
**Total Effort:** 2-3 hours (parallel) | 8-10 hours (sequential)
