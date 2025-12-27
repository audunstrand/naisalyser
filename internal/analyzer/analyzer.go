package analyzer

import (
	"fmt"
	"strings"

	"github.com/navikt/naisalyser/internal/github"
	"gopkg.in/yaml.v3"
)

// Analysis represents the complete analysis of a repository
type Analysis struct {
	Repository   *github.Repository
	NaisConfig   *NaisConfig
	Dependencies []Dependency
	Readme       string
	APIs         []APIEndpoint
}

// NaisConfig represents parsed nais.yaml configuration
type NaisConfig struct {
	APIVersion string       `yaml:"apiVersion"`
	Kind       string       `yaml:"kind"`
	Metadata   NaisMetadata `yaml:"metadata"`
	Spec       NaisSpec     `yaml:"spec"`
}

type NaisMetadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
}

type NaisSpec struct {
	Image        string            `yaml:"image"`
	Port         int               `yaml:"port"`
	Ingresses    []string          `yaml:"ingresses"`
	Replicas     *NaisReplicas     `yaml:"replicas"`
	Resources    *NaisResources    `yaml:"resources"`
	Liveness     *NaisProbe        `yaml:"liveness"`
	Readiness    *NaisProbe        `yaml:"readiness"`
	Env          []NaisEnvVar      `yaml:"env"`
	EnvFrom      []NaisEnvFrom     `yaml:"envFrom"`
	AccessPolicy *NaisAccessPolicy `yaml:"accessPolicy"`
	Azure        *NaisAzure        `yaml:"azure"`
	GCP          *NaisGCP          `yaml:"gcp"`
	Kafka        *NaisKafka        `yaml:"kafka"`
}

type NaisReplicas struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type NaisResources struct {
	Requests *NaisResourceSpec `yaml:"requests"`
	Limits   *NaisResourceSpec `yaml:"limits"`
}

type NaisResourceSpec struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type NaisProbe struct {
	Path string `yaml:"path"`
	Port int    `yaml:"port"`
}

type NaisEnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type NaisEnvFrom struct {
	Secret string `yaml:"secret"`
}

type NaisAccessPolicy struct {
	Inbound  *NaisAccessRules `yaml:"inbound"`
	Outbound *NaisAccessRules `yaml:"outbound"`
}

type NaisAccessRules struct {
	Rules    []NaisAccessRule   `yaml:"rules"`
	External []NaisExternalHost `yaml:"external"`
}

type NaisAccessRule struct {
	Application string `yaml:"application"`
	Namespace   string `yaml:"namespace"`
	Cluster     string `yaml:"cluster"`
}

type NaisExternalHost struct {
	Host string `yaml:"host"`
}

type NaisKafka struct {
	Pool string `yaml:"pool"`
}

type NaisAzure struct {
	Application *NaisAzureApp `yaml:"application"`
}

type NaisAzureApp struct {
	Enabled bool `yaml:"enabled"`
}

type NaisGCP struct {
	SQLInstances []NaisSQLInstance `yaml:"sqlInstances"`
	Buckets      []NaisBucket      `yaml:"buckets"`
}

type NaisSQLInstance struct {
	Type      string         `yaml:"type"`
	Databases []NaisDatabase `yaml:"databases"`
}

type NaisDatabase struct {
	Name string `yaml:"name"`
}

type NaisBucket struct {
	Name string `yaml:"name"`
}

// Dependency represents an application dependency
type Dependency struct {
	Name    string
	Version string
	Type    string // "maven", "npm", "go"
}

// APIEndpoint represents an exposed API endpoint
type APIEndpoint struct {
	Method string
	Path   string
	Source string
}

// Analyze performs complete analysis of a repository
func Analyze(gh *github.Client, repo *github.Repository) (*Analysis, error) {
	analysis := &Analysis{
		Repository: repo,
	}

	// Try to fetch README
	readme, err := gh.GetFileContent(repo.FullName, "README.md")
	if err == nil {
		analysis.Readme = readme
	}

	// Try to fetch nais.yaml (check common locations)
	naisConfig, err := fetchNaisConfig(gh, repo.FullName)
	if err == nil {
		analysis.NaisConfig = naisConfig
	}

	// Try to extract dependencies
	deps, err := fetchDependencies(gh, repo)
	if err == nil {
		analysis.Dependencies = deps
	}

	return analysis, nil
}

func fetchNaisConfig(gh *github.Client, fullName string) (*NaisConfig, error) {
	// Check common nais.yaml locations
	paths := []string{
		"nais.yaml",
		".nais/app.yaml",
		".nais/nais.yaml",
		"nais/nais.yaml",
		".nais/dev.yaml",
	}

	for _, path := range paths {
		content, err := gh.GetFileContent(fullName, path)
		if err == nil {
			var config NaisConfig
			if err := yaml.Unmarshal([]byte(content), &config); err != nil {
				continue
			}
			return &config, nil
		}
	}

	return nil, fmt.Errorf("no nais.yaml found")
}

func fetchDependencies(gh *github.Client, repo *github.Repository) ([]Dependency, error) {
	var deps []Dependency

	// Check language and try to parse dependencies accordingly
	switch strings.ToLower(repo.Language) {
	case "kotlin", "java":
		// Try build.gradle.kts
		if content, err := gh.GetFileContent(repo.FullName, "build.gradle.kts"); err == nil {
			deps = append(deps, parseGradleDeps(content)...)
		}
		// Try pom.xml
		if content, err := gh.GetFileContent(repo.FullName, "pom.xml"); err == nil {
			deps = append(deps, parseMavenDeps(content)...)
		}
	case "javascript", "typescript":
		if content, err := gh.GetFileContent(repo.FullName, "package.json"); err == nil {
			deps = append(deps, parseNpmDeps(content)...)
		}
	case "go":
		if content, err := gh.GetFileContent(repo.FullName, "go.mod"); err == nil {
			deps = append(deps, parseGoModDeps(content)...)
		}
	}

	return deps, nil
}

// Placeholder parsers - these would need more sophisticated implementation
func parseGradleDeps(content string) []Dependency {
	// Simple heuristic: look for implementation/api lines
	var deps []Dependency
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "implementation(") || strings.HasPrefix(line, "api(") {
			// Extract dependency string
			deps = append(deps, Dependency{
				Name: extractQuotedString(line),
				Type: "gradle",
			})
		}
	}
	return deps
}

func parseMavenDeps(content string) []Dependency {
	// Very basic XML parsing - in production use proper XML parser
	return []Dependency{}
}

func parseNpmDeps(content string) []Dependency {
	// Would parse package.json
	return []Dependency{}
}

func parseGoModDeps(content string) []Dependency {
	var deps []Dependency
	lines := strings.Split(content, "\n")
	inRequire := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "require (" {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}
		if inRequire && line != "" {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps = append(deps, Dependency{
					Name:    parts[0],
					Version: parts[1],
					Type:    "go",
				})
			}
		}
	}
	return deps
}

func extractQuotedString(s string) string {
	start := strings.Index(s, "\"")
	if start == -1 {
		return s
	}
	end := strings.LastIndex(s, "\"")
	if end <= start {
		return s
	}
	return s[start+1 : end]
}
