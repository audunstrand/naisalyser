# 🎯 Test Coverage Implementation Plan

**Status:** ✅ READY FOR EXECUTION
**Created:** December 27, 2025
**Coverage Target:** 1.7% → 85%+

---

## 📌 Executive Summary

The project currently has **1.7% test coverage**. This plan provides a complete, parallelizable strategy to achieve **85%+ coverage** through **10 test implementation tasks**.

All tasks are **independent and can run in parallel**, estimated to take **2-3 hours** if executed concurrently, or **8-10 hours** if sequential.

---

## 🎯 The Plan

### Location
- **Master Plan:** [`docs/test-plan.md`](docs/test-plan.md)
- **Task Specs:** [`docs/tasks/T1-10/`](docs/tasks/)
- **This File:** [`TEST_COVERAGE_PLAN.md`](TEST_COVERAGE_PLAN.md)

### Three Execution Phases

#### Phase 1️⃣ - Independent (Run in Parallel) ⏱️ ~1.5 hours
4 tasks with no dependencies - can ALL start immediately:

| Task | File | Tests | Component |
|------|------|-------|-----------|
| **T1** | [T1-parser-tests.md](docs/tasks/T1-parser-tests.md) | 40+ | Gradle, Maven, NPM, Go parsers |
| **T2** | [T2-handlebars-tests.md](docs/tasks/T2-handlebars-tests.md) | 15+ | Template preprocessing |
| **T5** | [T5-local-reader-tests.md](docs/tasks/T5-local-reader-tests.md) | 25+ | File I/O, language detection |
| **T7** | [T7-mermaid-tests.md](docs/tasks/T7-mermaid-tests.md) | 15+ | Graph visualization |

#### Phase 2️⃣ - Medium Dependencies (Start After Phase 1) ⏱️ ~2.5 hours
3 tasks that build on Phase 1 results:

| Task | File | Tests | Component | Dependency |
|------|------|-------|-----------|-----------|
| **T3** | [T3-graph-builder-tests.md](docs/tasks/T3-graph-builder-tests.md) | 30+ | Graph construction | T2 (optional) |
| **T4** | [T4-command-tests.md](docs/tasks/T4-command-tests.md) | 50+ | CLI handlers | None |
| **T6** | [T6-github-client-tests.md](docs/tasks/T6-github-client-tests.md) | 20+ | GitHub API (mocked) | None |

#### Phase 3️⃣ - Integration (Start After Phase 2) ⏱️ ~2 hours
2 tasks that verify all previous work:

| Task | File | Tests | Component |
|------|------|-------|-----------|
| **T8** | [T8-output-tests.md](docs/tasks/T8-output-tests.md) | 30+ | Documentation generation |
| **T10** | [T10-integration-tests.md](docs/tasks/T10-integration-tests.md) | 18+ | End-to-end workflows |

#### ✅ Already Complete
- **T9** - NaisConfig methods (3 tests) - **PASSING**

---

## 📊 By The Numbers

```
Current State      Target          Effort
────────────────────────────────────────────
Coverage:  1.7%  → 85%+          ~280 new tests
Tests:     3      → 280+          1,900+ lines of code
Packages:  1/6    → 6/6           10 task files
Files:     1      → 10            ~50KB documentation
```

---

## 🚀 How to Execute

### Option A: Full Parallelization (Recommended) 🏃 ~2-3 hours
```bash
# Phase 1: Start all 4 tasks simultaneously
/delegate Implement Task T1: Parser Tests
/delegate Implement Task T2: Handlebars Preprocessing  
/delegate Implement Task T5: Local Reader Tests
/delegate Implement Task T7: Graph Mermaid Output

# Wait for Phase 1 to complete (~1.5 hours)...

# Phase 2: Start 3 tasks simultaneously  
/delegate Implement Task T3: Graph Builder Tests
/delegate Implement Task T4: Command Handler Tests
/delegate Implement Task T6: GitHub Client Tests

# Wait for Phase 2 to complete (~2.5 hours)...

# Phase 3: Start 2 tasks simultaneously
/delegate Implement Task T8: Output Generation Tests
/delegate Implement Task T10: Integration Tests

# Wait for Phase 3 to complete (~2 hours)...
```

### Option B: Sequential (If limited resources) 🐢 ~8-10 hours
Run all tasks in order: T1 → T2 → T5 → T7 → T3 → T4 → T6 → T8 → T10

### Monitor Progress
```bash
# Check coverage as you go
go test ./... -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

---

## 📋 What Gets Tested

### Parsers (T1) - 40+ tests
✅ Gradle dependencies  
✅ Maven dependencies (XML parsing)  
✅ NPM dependencies (JSON parsing)  
✅ Go modules (go.mod)  

### Preprocessing (T2) - 15+ tests
✅ Handlebars template removal  
✅ Variable replacement  
✅ Context-aware placeholders  

### Local I/O (T5) - 25+ tests
✅ Repository reading  
✅ Language detection (Go, Java, Node, Python, Kotlin)  
✅ File operations  
✅ README extraction  

### Graph Building (T3) - 30+ tests
✅ Node management (deduplication)  
✅ Edge management (O(1) deduplication)  
✅ Type-safe constants  
✅ Access policy conversion  

### Commands (T4) - 50+ tests
✅ `analyze` command  
✅ `batch` command  
✅ `local` command  
✅ `local-batch` command  
✅ `graph` command  
✅ Flag handling  

### GitHub Client (T6) - 20+ tests
✅ Repository fetching (mocked)  
✅ File content reading (mocked)  
✅ Topic-based search (mocked)  
✅ Base64 decoding  

### Mermaid Output (T7) - 15+ tests
✅ Node rendering (app, kafka, external)  
✅ Edge rendering (call, kafka, external)  
✅ ID sanitization  
✅ Styling and classes  

### Documentation (T8) - 30+ tests
✅ Markdown generation  
✅ Template rendering  
✅ View model conversion  
✅ Section generation (overview, deps, config)  

### Integration (T10) - 18+ tests
✅ Full local analysis workflow  
✅ Graph building pipeline  
✅ Parser integration  
✅ Language detection pipeline  
✅ Command-to-output flows  

---

## ✅ Success Criteria

All of the following must be true:

- ✅ All 280+ tests passing
- ✅ Coverage ≥ 85% overall
- ✅ No regressions in existing code
- ✅ No external dependencies for tests
- ✅ Test execution < 10 seconds
- ✅ All error paths covered
- ✅ All edge cases tested
- ✅ Real-world workflows validated

---

## 📈 Coverage by Package (Target)

```
internal/analyzer   90%+  (core analysis logic)
internal/graph      85%+  (dependency graph building)
internal/local      90%+  (file I/O and language detection)
internal/github     80%+  (GitHub API client - all mocked)
internal/output     80%+  (documentation generation)
cmd                 75%+  (CLI command handlers)
───────────────────────────────────────────────
Overall             85%+  ← Target
```

---

## 🎯 Key Implementation Principles

1. **No External APIs** - All tests are self-contained (GitHub client is mocked)
2. **Realistic Data** - Test fixtures mirror real-world data
3. **Error Testing** - All error paths covered (nil checks, parsing errors, etc.)
4. **Performance** - Verify operations complete in acceptable time
5. **Parallel Safe** - Tasks don't interfere with each other
6. **Incremental** - Can verify coverage increasing after each phase

---

## 📁 File Structure After Implementation

```
naisalyser/
├── internal/
│   ├── analyzer/
│   │   ├── parsers_test.go          (T1)
│   │   ├── handlebars_test.go       (T2)
│   │   └── naisconfig_test.go       (T9 - exists)
│   ├── graph/
│   │   ├── builder_test.go          (T3)
│   │   └── mermaid_test.go          (T7)
│   ├── local/
│   │   └── reader_test.go           (T5)
│   ├── github/
│   │   └── client_test.go           (T6)
│   ├── output/
│   │   └── markdown_test.go         (T8)
│   └── integration_test.go          (T10)
├── cmd/
│   ├── commands_test.go             (T4)
│   └── flags_test.go                (T4)
└── docs/
    ├── test-plan.md                 (Master plan)
    └── tasks/
        ├── T1-parser-tests.md
        ├── T2-handlebars-tests.md
        ├── T3-graph-builder-tests.md
        ├── T4-command-tests.md
        ├── T5-local-reader-tests.md
        ├── T6-github-client-tests.md
        ├── T7-mermaid-tests.md
        ├── T8-output-tests.md
        ├── T10-integration-tests.md
        └── README.md
```

---

## 🔗 Quick Links

- 📖 **[Complete Plan](docs/test-plan.md)** - Detailed strategy document
- 📋 **[Task Index](docs/tasks/README.md)** - All tasks and specifications
- ✅ **[Completed Features](docs/tasks/01-*.md, etc)** - Previously completed work

---

## 💡 Tips for Success

1. **Start with Phase 1** - It has no dependencies, can all run in parallel
2. **Use Table-Driven Tests** - Makes tests more maintainable and comprehensive
3. **Mock External Systems** - GitHub client should never make real API calls
4. **Test Data in `testdata/`** - Keep test fixtures organized
5. **Verify Coverage Frequently** - Run `go test ./... -cover` after each task
6. **Read Task Specs First** - Each task file has complete details

---

## ⏱️ Timeline

```
Start → Phase 1 (1.5h) → Phase 2 (2.5h) → Phase 3 (2h) → Finish
        ↓
        4 tasks in parallel
                          ↓
                          3 tasks in parallel
                                            ↓
                                            2 tasks in parallel

Total: ~2-3 hours (parallel) or ~8-10 hours (sequential)
```

---

## 🎓 Learning Resources

Each task file includes:
- ✅ Detailed test case specifications
- ✅ Example test data
- ✅ Mock data structures
- ✅ Edge case scenarios
- ✅ Success criteria

---

## 📞 Questions?

Refer to:
1. **[docs/test-plan.md](docs/test-plan.md)** - Comprehensive strategy
2. **Task-specific files** - Detailed specifications for each task
3. **Existing tests** - `internal/analyzer/naisconfig_test.go` is a good template

---

**Status:** ✅ Ready to execute  
**Created:** 2025-12-27  
**Target Coverage:** 85%+  
**Estimated Effort:** 2-3 hours (parallel)

---

*This is the master document. Consult it to understand the overall strategy, then read individual task files for implementation details.*
