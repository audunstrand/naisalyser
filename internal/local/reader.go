package local

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Reader reads repository data from local filesystem
type Reader struct {
	verbose bool
}

// Repository represents a local repository
type Repository struct {
	Name          string
	FullName      string
	Path          string
	Description   string
	DefaultBranch string
	Language      string
}

// NewReader creates a new local reader
func NewReader(verbose bool) *Reader {
	return &Reader{verbose: verbose}
}

// ReadRepository reads repository metadata from local path
func (r *Reader) ReadRepository(path string, org string) (*Repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("path does not exist: %s", absPath)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", absPath)
	}

	name := filepath.Base(absPath)

	// Try to detect language from files
	language := r.detectLanguage(absPath)

	// Try to read description from README
	description := r.readFirstLineOfReadme(absPath)

	repo := &Repository{
		Name:        name,
		FullName:    fmt.Sprintf("%s/%s", org, name),
		Path:        absPath,
		Description: description,
		Language:    language,
	}

	return repo, nil
}

// GetFileContent reads a file's content from local filesystem
func (r *Reader) GetFileContent(repoPath, filePath string) (string, error) {
	fullPath := filepath.Join(repoPath, filePath)
	if r.verbose {
		fmt.Printf("  → reading %s\n", filePath)
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// FileExists checks if a file exists in the repository
func (r *Reader) FileExists(repoPath, filePath string) bool {
	fullPath := filepath.Join(repoPath, filePath)
	_, err := os.Stat(fullPath)
	return err == nil
}

// detectLanguage detects the primary language based on files present
func (r *Reader) detectLanguage(path string) string {
	checks := []struct {
		file     string
		language string
	}{
		{"build.gradle.kts", "Kotlin"},
		{"build.gradle", "Java"},
		{"pom.xml", "Java"},
		{"package.json", "TypeScript"},
		{"go.mod", "Go"},
		{"requirements.txt", "Python"},
		{"Cargo.toml", "Rust"},
	}

	for _, check := range checks {
		if r.FileExists(path, check.file) {
			// Check for TypeScript
			if check.file == "package.json" && r.FileExists(path, "tsconfig.json") {
				return "TypeScript"
			}
			return check.language
		}
	}

	return "Unknown"
}

// readFirstLineOfReadme reads first non-header line from README
func (r *Reader) readFirstLineOfReadme(path string) string {
	content, err := r.GetFileContent(path, "README.md")
	if err != nil {
		return ""
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Return first real content line, truncated
		if len(line) > 100 {
			return line[:100] + "..."
		}
		return line
	}

	return ""
}
