# NAISalyser Improvement Tasks

This directory contains detailed implementation tasks for improving the NAISalyser codebase.

## Task Priority Order

Execute tasks in this order due to dependencies:

### Phase 1: Foundation (No Dependencies)
1. `01-add-naisconfig-behavior-methods.md` - Add domain methods to NaisConfig
2. `02-add-type-safe-node-edge-types.md` - Replace string types with constants
3. `03-remove-unused-kafka-patterns.md` - Clean up dead code
4. `04-remove-stub-parsers.md` - Remove or implement stub functions

### Phase 2: Refactoring (Depends on Phase 1)
5. `05-extract-repository-reader-interface.md` - Unify GitHub/local readers
6. `06-refactor-graph-builder-use-domain-methods.md` - Use new NaisConfig methods
7. `07-optimize-edge-deduplication.md` - O(1) edge lookup

### Phase 3: Features (Depends on Phase 2)
8. `08-implement-npm-dependency-parser.md` - Parse package.json
9. `09-implement-maven-dependency-parser.md` - Parse pom.xml
10. `10-add-configurable-organization.md` - Remove hardcoded "navikt"
11. `11-add-error-handling-for-flags.md` - Handle cobra flag errors
12. `12-flatten-handlebars-preprocessing.md` - Reduce nesting

## Running Tests

After each task, verify with:
```bash
go build ./...
go test ./...
```

## Code Style

- Keep functions under 30 lines
- One level of indentation preferred
- Use early returns
- Add comments only for non-obvious logic
