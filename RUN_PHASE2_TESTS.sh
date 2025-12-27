#!/bin/bash
# Test Implementation - Phase 2 Parallel Delegations
# Created: 2025-12-27
# Purpose: Delegate 3 dependent test implementation tasks (after Phase 1 completes)
# Note: These tasks can run in parallel AFTER Phase 1 is complete

echo "🚀 PHASE 2: TEST IMPLEMENTATION TASKS (After Phase 1)"
echo "═════════════════════════════════════════════════════════════════"
echo ""
echo "ℹ️  Only run this AFTER Phase 1 tasks (T1, T2, T5, T7) are complete"
echo ""
echo "Starting 3 independent agents to implement 100+ tests in parallel..."
echo ""

# TASK 3: Graph Builder Tests
echo "1️⃣  Delegating T3: Graph Builder Tests"
/delegate Implement comprehensive tests for graph construction. Create internal/graph/builder_test.go with 30+ test cases covering: NewBuilder (4 tests), addNode (6 tests for deduplication), addEdge (8 tests for O(1) deduplication with edgeKey map), AddFromLocalAnalysis (16 tests for inbound/outbound/kafka/external). Test type-safe constants (NodeTypeApp, EdgeTypeCall, etc). Verify edge deduplication uses map, not linear scan. Test nil handling and namespace preservation. Follow specification in docs/tasks/T3-graph-builder-tests.md.

# TASK 4: Command Handler Tests
echo "2️⃣  Delegating T4: Command Handler Tests"
/delegate Implement comprehensive tests for CLI command handlers. Create cmd/commands_test.go and cmd/flags_test.go with 50+ test cases covering: mustGetString/mustGetBool helpers (4 tests), runAnalyze (8 tests), runBatch (10 tests for topic/config flags), runLocal (10 tests for org parameter), runLocalBatch (10 tests), runGraph (12 tests). Use t.TempDir() for test repos. Mock filesystem. Test flag handling, error cases, output generation. Follow specification in docs/tasks/T4-command-tests.md.

# TASK 6: GitHub Client Tests
echo "3️⃣  Delegating T6: GitHub Client Tests"
/delegate Implement comprehensive tests for GitHub API client (all mocked). Create internal/github/client_test.go with 20+ test cases covering: NewClient (2 tests), GetRepository (5 tests), GetRepositoriesByTopic (6 tests), GetFileContent (6 tests for base64 decoding), FileExists (4 tests), GetDirectoryContents (5 tests). Mock all gh CLI command execution. Include error scenarios (404, 403, 500). Test pagination. NO REAL API CALLS. Follow specification in docs/tasks/T6-github-client-tests.md.

echo ""
echo "═════════════════════════════════════════════════════════════════"
echo "✅ Phase 2 delegations submitted!"
echo ""
echo "These 3 agents will work in parallel and create separate PRs."
echo "Total tests: 100+ (30 + 50 + 20)"
echo "Expected duration: ~2.5 hours"
echo ""
echo "Prerequisites:"
echo "✓ Phase 1 must be complete (T1, T2, T5, T7)"
echo "✓ All Phase 1 PRs merged"
echo ""
echo "Next steps:"
echo "1. Wait for all Phase 2 tasks to complete"
echo "2. Review and merge the 3 PRs"
echo "3. Run: go test ./... -cover"
echo "4. Start Phase 3 when Phase 2 is done"
echo ""
