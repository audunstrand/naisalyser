package analyzer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/navikt/naisalyser/internal/local"
	"gopkg.in/yaml.v3"
)

// LocalAnalysis represents analysis from local repository
type LocalAnalysis struct {
	Repository   *local.Repository
	NaisConfig   *NaisConfig
	Dependencies []Dependency
	Readme       string
	APIs         []APIEndpoint
	KafkaTopics  []KafkaTopic
}

// AnalyzeLocal performs complete analysis of a local repository
func AnalyzeLocal(reader *local.Reader, repo *local.Repository) (*LocalAnalysis, error) {
	analysis := &LocalAnalysis{
		Repository: repo,
	}

	// Try to fetch README
	readme, err := reader.GetFileContent(repo.Path, "README.md")
	if err == nil {
		analysis.Readme = readme
	}

	// Try to fetch nais.yaml (check common locations)
	naisConfig, _ := fetchLocalNaisConfig(reader, repo.Path)
	if naisConfig != nil {
		analysis.NaisConfig = naisConfig
	}

	// Try to extract dependencies
	deps, err := fetchLocalDependencies(reader, repo)
	if err == nil {
		analysis.Dependencies = deps
	}

	// Extract Kafka topics from source code
	analysis.KafkaTopics = ExtractKafkaTopics(reader, repo.Path)

	return analysis, nil
}

func fetchLocalNaisConfig(reader *local.Reader, repoPath string) (*NaisConfig, error) {
	// Check common nais.yaml locations first
	staticPaths := []string{
		"nais.yaml",
		"nais/naiserator.yml",
		"nais/naiserator.yaml",
		".nais/app.yaml",
		".nais/nais.yaml",
		"nais/nais.yaml",
		".nais/dev.yaml",
		".nais/naiserator.yaml",
		"nais/app.yaml",
	}

	for _, path := range staticPaths {
		if config := tryParseNaisConfig(reader, repoPath, path); config != nil {
			return config, nil
		}
	}

	// Scan nais/ directory for any yml files, preferring prod over dev
	naisDir := filepath.Join(repoPath, "nais")
	if config := scanNaisDirectory(reader, repoPath, naisDir, "nais"); config != nil {
		return config, nil
	}

	// Also scan .nais/ directory (used by familie-* repos)
	dotNaisDir := filepath.Join(repoPath, ".nais")
	if config := scanNaisDirectory(reader, repoPath, dotNaisDir, ".nais"); config != nil {
		return config, nil
	}

	return nil, nil // No config found, not an error
}

// scanNaisDirectory scans a directory for nais configs, preferring prod over dev
func scanNaisDirectory(reader *local.Reader, repoPath, dirPath, prefix string) *NaisConfig {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil
	}

	// Sort files: prod first, then others
	var prodFiles, otherFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".yml") && !strings.HasSuffix(name, ".yaml") {
			continue
		}
		if strings.Contains(strings.ToLower(name), "prod") {
			prodFiles = append(prodFiles, name)
		} else if !strings.Contains(strings.ToLower(name), "dev") && !strings.Contains(strings.ToLower(name), "lokal") {
			otherFiles = append(otherFiles, name)
		}
	}

	// Try prod files first
	for _, name := range prodFiles {
		path := filepath.Join(prefix, name)
		if config := tryParseNaisConfig(reader, repoPath, path); config != nil {
			return config
		}
	}
	// Then try other non-dev files
	for _, name := range otherFiles {
		path := filepath.Join(prefix, name)
		if config := tryParseNaisConfig(reader, repoPath, path); config != nil {
			return config
		}
	}
	return nil
}

func tryParseNaisConfig(reader *local.Reader, repoPath, path string) *NaisConfig {
	if !reader.FileExists(repoPath, path) {
		return nil
	}
	content, err := reader.GetFileContent(repoPath, path)
	if err != nil {
		return nil
	}
	processedContent := preprocessHandlebars(content)
	var config NaisConfig
	if err := yaml.Unmarshal([]byte(processedContent), &config); err != nil {
		return nil
	}
	return &config
}

// preprocessHandlebars replaces handlebars template syntax with placeholder values
func preprocessHandlebars(content string) string {
	result := removeEachBlocks(content)
	result = replaceTemplateVars(result)
	return result
}

// removeEachBlocks removes {{#each ...}} ... {{/each}} blocks from content
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
		// Find the line start and end to remove entire lines
		lineStart := strings.LastIndex(result[:start], "\n") + 1
		lineEnd := start + end + len("{{/each}}")
		if nextNewline := strings.Index(result[lineEnd:], "\n"); nextNewline != -1 {
			lineEnd += nextNewline + 1
		}
		result = result[:lineStart] + result[lineEnd:]
	}
	return result
}

// replaceTemplateVars replaces simple {{variable}} with context-appropriate placeholders
func replaceTemplateVars(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = replaceLineTemplateVars(line)
	}
	return strings.Join(lines, "\n")
}

// replaceLineTemplateVars replaces template vars in a single line
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

		// Determine replacement based on context
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
	return line
}

func fetchLocalDependencies(reader *local.Reader, repo *local.Repository) ([]Dependency, error) {
	var deps []Dependency

	// Check language and try to parse dependencies accordingly
	switch strings.ToLower(repo.Language) {
	case "kotlin", "java":
		// Try build.gradle.kts
		if content, err := reader.GetFileContent(repo.Path, "build.gradle.kts"); err == nil {
			deps = append(deps, parseGradleDeps(content)...)
		}
		// Try pom.xml
		if content, err := reader.GetFileContent(repo.Path, "pom.xml"); err == nil {
			deps = append(deps, parseMavenDeps(content)...)
		}
	case "javascript", "typescript":
		if content, err := reader.GetFileContent(repo.Path, "package.json"); err == nil {
			deps = append(deps, parseNpmDeps(content)...)
		}
	case "go":
		if content, err := reader.GetFileContent(repo.Path, "go.mod"); err == nil {
			deps = append(deps, parseGoModDeps(content)...)
		}
	}

	return deps, nil
}
