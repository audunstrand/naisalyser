package analyzer

import (
	"strings"
	"testing"
)

// Tests for removeEachBlocks

func TestRemoveEachBlocks_SingleBlock(t *testing.T) {
	input := `before
{{#each items}}
  - {{this}}
{{/each}}
after`
	expected := `before
after`

	result := removeEachBlocks(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRemoveEachBlocks_MultipleBlocks(t *testing.T) {
	input := `first
{{#each list1}}
  - {{name}}
{{/each}}
middle
{{#each list2}}
  - {{value}}
{{/each}}
last`
	expected := `first
middle
last`

	result := removeEachBlocks(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRemoveEachBlocks_NestedBlocks(t *testing.T) {
	input := `outer
{{#each outer}}
  inner
  {{#each inner}}
    - {{value}}
  {{/each}}
{{/each}}
done`

	result := removeEachBlocks(input)
	// The outer block removal will capture everything up to its closing tag,
	// which includes the nested block, so both should be removed
	if strings.Contains(result, "{{#each") {
		t.Error("expected {{#each}} tags to be removed")
	}
	if !strings.Contains(result, "outer") || !strings.Contains(result, "done") {
		t.Error("expected surrounding content to remain")
	}
	// Note: nested blocks within removed blocks are also removed
}

func TestRemoveEachBlocks_NoBlocks(t *testing.T) {
	input := `simple: {{value}}
name: {{appName}}`
	expected := input

	result := removeEachBlocks(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRemoveEachBlocks_MalformedMissingClose(t *testing.T) {
	input := `before
{{#each items}}
  - {{this}}
after`

	result := removeEachBlocks(input)
	// Should not modify if closing tag is missing
	if result != input {
		t.Errorf("expected unchanged input when closing tag missing, got %q", result)
	}
}

// Tests for replaceLineTemplateVars

func TestReplaceLineTemplateVars_SingleVariable(t *testing.T) {
	input := "name: {{appName}}"
	expected := "name: placeholder"

	result := replaceLineTemplateVars(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestReplaceLineTemplateVars_MultipleVariables(t *testing.T) {
	input := "{{a}} and {{b}} and {{c}}"
	expected := "placeholder and placeholder and placeholder"

	result := replaceLineTemplateVars(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestReplaceLineTemplateVars_NoVariables(t *testing.T) {
	input := "plain text line"
	expected := "plain text line"

	result := replaceLineTemplateVars(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestReplaceLineTemplateVars_MalformedMissingClose(t *testing.T) {
	input := "name: {{appName"

	result := replaceLineTemplateVars(input)
	// Should not replace if closing braces are missing
	if result != input {
		t.Errorf("expected unchanged input when closing braces missing, got %q", result)
	}
}

// Tests for context-aware replacement

func TestReplaceLineTemplateVars_ReplicasContext(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"replicas field", "replicas: {{replicaCount}}", "replicas: 1"},
		{"min replicas", "min: {{minReplicas}}", "min: 1"},
		{"max replicas", "max: {{maxReplicas}}", "max: 1"},
		{"uppercase REPLICAS", "REPLICAS: {{count}}", "REPLICAS: 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceLineTemplateVars(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestReplaceLineTemplateVars_PortContext(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"port field", "port: {{servicePort}}", "port: 1"},
		{"uppercase PORT", "PORT: {{value}}", "PORT: 1"},
		{"targetPort", "targetPort: {{port}}", "targetPort: 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceLineTemplateVars(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestReplaceLineTemplateVars_TimeoutContext(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"timeout field", "timeout: {{requestTimeout}}", "timeout: 1"},
		{"delay field", "delay: {{startupDelay}}", "delay: 1"},
		{"uppercase TIMEOUT", "TIMEOUT: {{value}}", "TIMEOUT: 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceLineTemplateVars(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestReplaceLineTemplateVars_NonNumericContext(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"name field", "name: {{appName}}", "name: placeholder"},
		{"image field", "image: {{dockerImage}}", "image: placeholder"},
		{"namespace field", "namespace: {{ns}}", "namespace: placeholder"},
		{"url field", "url: {{serviceUrl}}", "url: placeholder"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceLineTemplateVars(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestReplaceLineTemplateVars_MultipleVariablesMixedContext(t *testing.T) {
	input := "name: {{appName}} replicas: {{count}}"
	// The line contains "replicas", so all variables get "1"
	expected := "name: 1 replicas: 1"

	result := replaceLineTemplateVars(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// Tests for replaceTemplateVars

func TestReplaceTemplateVars_MultipleLines(t *testing.T) {
	input := `name: {{appName}}
replicas: {{count}}
port: {{servicePort}}`
	expected := `name: placeholder
replicas: 1
port: 1`

	result := replaceTemplateVars(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestReplaceTemplateVars_EmptyContent(t *testing.T) {
	input := ""
	expected := ""

	result := replaceTemplateVars(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestReplaceTemplateVars_NoTemplates(t *testing.T) {
	input := `name: myapp
replicas: 3
port: 8080`
	expected := input

	result := replaceTemplateVars(input)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// Tests for preprocessHandlebars (integration)

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

func TestPreprocessHandlebars_RealNaisYAML(t *testing.T) {
	input := `apiVersion: "nais.io/v1alpha1"
kind: "Application"
metadata:
  name: {{appName}}
  namespace: {{namespace}}
spec:
  image: {{image}}
  port: {{port}}
  replicas:
    min: {{minReplicas}}
    max: {{maxReplicas}}
  resources:
    requests:
      cpu: {{cpu}}
      memory: {{memory}}
  liveness:
    path: /health
    timeout: {{timeout}}
  env:
    - name: ENV_VAR
      value: {{envValue}}`

	result := preprocessHandlebars(input)

	// Verify structure is maintained
	if !strings.Contains(result, "apiVersion:") {
		t.Error("expected apiVersion to be preserved")
	}
	if !strings.Contains(result, "kind:") {
		t.Error("expected kind to be preserved")
	}

	// Verify no template syntax remains
	if strings.Contains(result, "{{") || strings.Contains(result, "}}") {
		t.Error("expected all template syntax to be replaced")
	}

	// Verify numeric contexts got "1"
	if !strings.Contains(result, "port: 1") {
		t.Error("expected port to be replaced with '1'")
	}
	if !strings.Contains(result, "min: 1") {
		t.Error("expected min replicas to be replaced with '1'")
	}
	if !strings.Contains(result, "max: 1") {
		t.Error("expected max replicas to be replaced with '1'")
	}
	if !strings.Contains(result, "timeout: 1") {
		t.Error("expected timeout to be replaced with '1'")
	}

	// Verify non-numeric contexts got "placeholder"
	if !strings.Contains(result, "name: placeholder") {
		t.Error("expected name to be replaced with 'placeholder'")
	}
	if !strings.Contains(result, "namespace: placeholder") {
		t.Error("expected namespace to be replaced with 'placeholder'")
	}
	if !strings.Contains(result, "image: placeholder") {
		t.Error("expected image to be replaced with 'placeholder'")
	}
}

func TestPreprocessHandlebars_RealNaisYAMLWithEachBlock(t *testing.T) {
	input := `apiVersion: "nais.io/v1alpha1"
kind: "Application"
metadata:
  name: {{appName}}
spec:
  replicas:
    min: {{minReplicas}}
  env:
{{#each envVars}}
    - name: {{name}}
      value: {{value}}
{{/each}}
  port: {{port}}`

	result := preprocessHandlebars(input)

	// Verify {{#each}} block is removed
	if strings.Contains(result, "{{#each") || strings.Contains(result, "{{/each}}") {
		t.Error("expected {{#each}} block to be removed")
	}

	// Verify remaining templates are replaced
	if strings.Contains(result, "{{") || strings.Contains(result, "}}") {
		t.Error("expected all template syntax to be replaced")
	}

	// Verify structure is maintained
	if !strings.Contains(result, "apiVersion:") {
		t.Error("expected apiVersion to be preserved")
	}
	if !strings.Contains(result, "env:") {
		t.Error("expected env section to be preserved")
	}

	// Verify replacements
	if !strings.Contains(result, "name: placeholder") {
		t.Error("expected name to be replaced with 'placeholder'")
	}
	if !strings.Contains(result, "min: 1") {
		t.Error("expected min replicas to be replaced with '1'")
	}
	if !strings.Contains(result, "port: 1") {
		t.Error("expected port to be replaced with '1'")
	}
}

// Tests for getPlaceholderForContext

func TestGetPlaceholderForContext_NumericContexts(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{"replica context", "replicas: {{value}}", "1"},
		{"min context", "min: {{value}}", "1"},
		{"max context", "max: {{value}}", "1"},
		{"port context", "port: {{value}}", "1"},
		{"timeout context", "timeout: {{value}}", "1"},
		{"delay context", "delay: {{value}}", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPlaceholderForContext(tt.line)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetPlaceholderForContext_NonNumericContext(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{"name context", "name: {{value}}", "placeholder"},
		{"image context", "image: {{value}}", "placeholder"},
		{"url context", "url: {{value}}", "placeholder"},
		{"random context", "something: {{value}}", "placeholder"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPlaceholderForContext(tt.line)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
