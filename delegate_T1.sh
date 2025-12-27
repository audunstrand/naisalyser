#!/bin/bash
# Phase 1: Parser Tests Implementation
# This will delegate the parser tests task to an agent

cd /Users/audunfauchaldstrand/naisalyser

echo "🚀 DELEGATING T1: Parser Tests (40+ tests)"
echo "══════════════════════════════════════════════════════════════"
echo ""
echo "Task: Implement comprehensive tests for dependency parsers"
echo "Files: internal/analyzer/parsers_test.go"
echo "Tests: 40+ (parseGradleDeps, parseNpmDeps, parseMavenDeps, parseGoModDeps)"
echo "Spec: docs/tasks/T1-parser-tests.md"
echo ""
echo "Delegating to agent..."
echo ""

/delegate "Task T1: Parser Tests - Implement 40+ comprehensive tests for Gradle, Maven, NPM, and Go dependency parsers in internal/analyzer/parsers_test.go. Test parseGradleDeps, parseNpmDeps, parseMavenDeps, parseGoModDeps with valid inputs, malformed data, edge cases, and version handling. Reference docs/tasks/T1-parser-tests.md for complete specifications. Include table-driven tests and realistic test data fixtures."
