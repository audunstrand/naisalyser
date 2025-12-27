#!/bin/bash
# Phase 1: Mermaid Output Tests Implementation

cd /Users/audunfauchaldstrand/naisalyser

echo "🚀 DELEGATING T7: Mermaid Output Tests (15+ tests)"
echo "══════════════════════════════════════════════════════════════"
echo ""
echo "Task: Implement comprehensive tests for graph visualization"
echo "Files: internal/graph/mermaid_test.go"
echo "Tests: 15+ (ToMermaid, sanitizeMermaidID)"
echo "Spec: docs/tasks/T7-mermaid-tests.md"
echo ""
echo "Delegating to agent..."
echo ""

/delegate "Task T7: Mermaid Output Tests - Implement 15+ tests in internal/graph/mermaid_test.go for sanitizeMermaidID and ToMermaid functions. Test node rendering (app/kafka/external), edge rendering (call/kafka/external), stateful vs stateless styling, ID sanitization (hyphens, dots, slashes, colons). Verify type-safe constants (NodeTypeApp, EdgeTypeCall, etc) are used and output is valid Mermaid syntax. Reference docs/tasks/T7-mermaid-tests.md for complete specifications."
