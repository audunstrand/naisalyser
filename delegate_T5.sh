#!/bin/bash
# Phase 1: Local Reader Tests Implementation

cd /Users/audunfauchaldstrand/naisalyser

echo "🚀 DELEGATING T5: Local Reader Tests (25+ tests)"
echo "══════════════════════════════════════════════════════════════"
echo ""
echo "Task: Implement comprehensive tests for file I/O operations"
echo "Files: internal/local/reader_test.go"
echo "Tests: 25+ (ReadRepository, GetFileContent, FileExists, detectLanguage, readFirstLineOfReadme)"
echo "Spec: docs/tasks/T5-local-reader-tests.md"
echo ""
echo "Delegating to agent..."
echo ""

/delegate "Task T5: Local Reader Tests - Implement 25+ tests in internal/local/reader_test.go covering ReadRepository with org parameter, GetFileContent, FileExists, detectLanguage (Go, Java, Kotlin, Node, Python), and readFirstLineOfReadme. Use t.TempDir() for test repositories. Include language priority detection, edge cases, and error handling. Reference docs/tasks/T5-local-reader-tests.md for complete specifications."
