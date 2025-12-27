package local

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestReadRepository_ValidPath tests reading a valid repository path
func TestReadRepository_ValidPath(t *testing.T) {
	tmpDir := t.TempDir()
	
	reader := NewReader(false)
	repo, err := reader.ReadRepository(tmpDir, "testorg")
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo == nil {
		t.Fatal("expected repo to not be nil")
	}
	if repo.Path != tmpDir {
		t.Errorf("expected Path=%s, got %s", tmpDir, repo.Path)
	}
	if repo.Name != filepath.Base(tmpDir) {
		t.Errorf("expected Name=%s, got %s", filepath.Base(tmpDir), repo.Name)
	}
}

// TestReadRepository_WithOrgParameter tests org parameter is used correctly
func TestReadRepository_WithOrgParameter(t *testing.T) {
	tmpDir := t.TempDir()
	
	reader := NewReader(false)
	repo, err := reader.ReadRepository(tmpDir, "custom-org")
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	expectedFullName := "custom-org/" + filepath.Base(tmpDir)
	if repo.FullName != expectedFullName {
		t.Errorf("expected FullName=%s, got %s", expectedFullName, repo.FullName)
	}
}

// TestReadRepository_NaviktOrg tests default navikt org
func TestReadRepository_NaviktOrg(t *testing.T) {
	tmpDir := t.TempDir()
	
	reader := NewReader(false)
	repo, err := reader.ReadRepository(tmpDir, "navikt")
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	expectedFullName := "navikt/" + filepath.Base(tmpDir)
	if repo.FullName != expectedFullName {
		t.Errorf("expected FullName=%s, got %s", expectedFullName, repo.FullName)
	}
}

// TestReadRepository_NonExistentPath tests error handling for missing path
func TestReadRepository_NonExistentPath(t *testing.T) {
	reader := NewReader(false)
	_, err := reader.ReadRepository("/non/existent/path", "testorg")
	
	if err == nil {
		t.Fatal("expected error for non-existent path")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("expected error to mention 'does not exist', got: %v", err)
	}
}

// TestReadRepository_FileNotDirectory tests error when path is a file
func TestReadRepository_FileNotDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	_, err := reader.ReadRepository(filePath, "testorg")
	
	if err == nil {
		t.Fatal("expected error for file path")
	}
	if !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("expected error to mention 'not a directory', got: %v", err)
	}
}

// TestGetFileContent_Success tests reading an existing file
func TestGetFileContent_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := "test.txt"
	expectedContent := "hello world"
	
	fullPath := filepath.Join(tmpDir, filePath)
	if err := os.WriteFile(fullPath, []byte(expectedContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	content, err := reader.GetFileContent(tmpDir, filePath)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if content != expectedContent {
		t.Errorf("expected content=%q, got %q", expectedContent, content)
	}
}

// TestGetFileContent_FileNotFound tests error for missing file
func TestGetFileContent_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	
	reader := NewReader(false)
	_, err := reader.GetFileContent(tmpDir, "nonexistent.txt")
	
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

// TestGetFileContent_Subdirectory tests reading file in subdirectory
func TestGetFileContent_Subdirectory(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	filePath := "subdir/test.txt"
	expectedContent := "nested content"
	fullPath := filepath.Join(tmpDir, filePath)
	if err := os.WriteFile(fullPath, []byte(expectedContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	content, err := reader.GetFileContent(tmpDir, filePath)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if content != expectedContent {
		t.Errorf("expected content=%q, got %q", expectedContent, content)
	}
}

// TestGetFileContent_EmptyFile tests reading an empty file
func TestGetFileContent_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := "empty.txt"
	
	fullPath := filepath.Join(tmpDir, filePath)
	if err := os.WriteFile(fullPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	content, err := reader.GetFileContent(tmpDir, filePath)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if content != "" {
		t.Errorf("expected empty content, got %q", content)
	}
}

// TestFileExists_ExistingFile tests FileExists returns true for existing file
func TestFileExists_ExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := "test.txt"
	
	fullPath := filepath.Join(tmpDir, filePath)
	if err := os.WriteFile(fullPath, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	exists := reader.FileExists(tmpDir, filePath)
	
	if !exists {
		t.Error("expected FileExists to return true for existing file")
	}
}

// TestFileExists_NonExistentFile tests FileExists returns false for missing file
func TestFileExists_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	
	reader := NewReader(false)
	exists := reader.FileExists(tmpDir, "nonexistent.txt")
	
	if exists {
		t.Error("expected FileExists to return false for non-existent file")
	}
}

// TestFileExists_Directory tests FileExists works with directories
func TestFileExists_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	exists := reader.FileExists(tmpDir, "subdir")
	
	if !exists {
		t.Error("expected FileExists to return true for existing directory")
	}
}

// TestDetectLanguage_Go tests Go language detection
func TestDetectLanguage_Go(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create go.mod file
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	lang := reader.detectLanguage(tmpDir)
	
	if lang != "Go" {
		t.Errorf("expected language=Go, got %s", lang)
	}
}

// TestDetectLanguage_Java tests Java language detection via pom.xml
func TestDetectLanguage_Java(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create pom.xml file
	pomPath := filepath.Join(tmpDir, "pom.xml")
	if err := os.WriteFile(pomPath, []byte("<project></project>"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	lang := reader.detectLanguage(tmpDir)
	
	if lang != "Java" {
		t.Errorf("expected language=Java, got %s", lang)
	}
}

// TestDetectLanguage_JavaGradle tests Java language detection via build.gradle
func TestDetectLanguage_JavaGradle(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create build.gradle file
	gradlePath := filepath.Join(tmpDir, "build.gradle")
	if err := os.WriteFile(gradlePath, []byte("plugins {}"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	lang := reader.detectLanguage(tmpDir)
	
	if lang != "Java" {
		t.Errorf("expected language=Java, got %s", lang)
	}
}

// TestDetectLanguage_Kotlin tests Kotlin language detection
func TestDetectLanguage_Kotlin(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create build.gradle.kts file
	ktsPath := filepath.Join(tmpDir, "build.gradle.kts")
	if err := os.WriteFile(ktsPath, []byte("plugins {}"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	lang := reader.detectLanguage(tmpDir)
	
	if lang != "Kotlin" {
		t.Errorf("expected language=Kotlin, got %s", lang)
	}
}

// TestDetectLanguage_TypeScript tests TypeScript language detection
func TestDetectLanguage_TypeScript(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create package.json and tsconfig.json files
	pkgPath := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(pkgPath, []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	tsconfigPath := filepath.Join(tmpDir, "tsconfig.json")
	if err := os.WriteFile(tsconfigPath, []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	lang := reader.detectLanguage(tmpDir)
	
	if lang != "TypeScript" {
		t.Errorf("expected language=TypeScript, got %s", lang)
	}
}

// TestDetectLanguage_Node tests Node language detection (no TypeScript)
func TestDetectLanguage_Node(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create package.json but no tsconfig.json
	pkgPath := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(pkgPath, []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	lang := reader.detectLanguage(tmpDir)
	
	if lang != "TypeScript" {
		t.Errorf("expected language=TypeScript, got %s", lang)
	}
}

// TestDetectLanguage_Python tests Python language detection
func TestDetectLanguage_Python(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create requirements.txt file
	reqPath := filepath.Join(tmpDir, "requirements.txt")
	if err := os.WriteFile(reqPath, []byte("pytest==7.0.0"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	lang := reader.detectLanguage(tmpDir)
	
	if lang != "Python" {
		t.Errorf("expected language=Python, got %s", lang)
	}
}

// TestDetectLanguage_Rust tests Rust language detection
func TestDetectLanguage_Rust(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create Cargo.toml file
	cargoPath := filepath.Join(tmpDir, "Cargo.toml")
	if err := os.WriteFile(cargoPath, []byte("[package]"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	lang := reader.detectLanguage(tmpDir)
	
	if lang != "Rust" {
		t.Errorf("expected language=Rust, got %s", lang)
	}
}

// TestDetectLanguage_Unknown tests Unknown language when no marker files
func TestDetectLanguage_Unknown(t *testing.T) {
	tmpDir := t.TempDir()
	
	reader := NewReader(false)
	lang := reader.detectLanguage(tmpDir)
	
	if lang != "Unknown" {
		t.Errorf("expected language=Unknown, got %s", lang)
	}
}

// TestDetectLanguage_Priority tests language priority (Kotlin over Java)
func TestDetectLanguage_Priority(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create both build.gradle.kts and build.gradle
	// Should prioritize Kotlin
	ktsPath := filepath.Join(tmpDir, "build.gradle.kts")
	if err := os.WriteFile(ktsPath, []byte("plugins {}"), 0644); err != nil {
		t.Fatal(err)
	}
	gradlePath := filepath.Join(tmpDir, "build.gradle")
	if err := os.WriteFile(gradlePath, []byte("plugins {}"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	lang := reader.detectLanguage(tmpDir)
	
	if lang != "Kotlin" {
		t.Errorf("expected language=Kotlin (priority over Java), got %s", lang)
	}
}

// TestReadFirstLineOfReadme_Success tests reading first line from README
func TestReadFirstLineOfReadme_Success(t *testing.T) {
	tmpDir := t.TempDir()
	
	readmeContent := `# Title

This is the first content line.
More content here.`
	
	readmePath := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	desc := reader.readFirstLineOfReadme(tmpDir)
	
	expected := "This is the first content line."
	if desc != expected {
		t.Errorf("expected description=%q, got %q", expected, desc)
	}
}

// TestReadFirstLineOfReadme_NoReadme tests handling of missing README
func TestReadFirstLineOfReadme_NoReadme(t *testing.T) {
	tmpDir := t.TempDir()
	
	reader := NewReader(false)
	desc := reader.readFirstLineOfReadme(tmpDir)
	
	if desc != "" {
		t.Errorf("expected empty description, got %q", desc)
	}
}

// TestReadFirstLineOfReadme_EmptyReadme tests handling of empty README
func TestReadFirstLineOfReadme_EmptyReadme(t *testing.T) {
	tmpDir := t.TempDir()
	
	readmePath := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	desc := reader.readFirstLineOfReadme(tmpDir)
	
	if desc != "" {
		t.Errorf("expected empty description, got %q", desc)
	}
}

// TestReadFirstLineOfReadme_HeaderOnly tests README with only headers
func TestReadFirstLineOfReadme_HeaderOnly(t *testing.T) {
	tmpDir := t.TempDir()
	
	readmeContent := `# Title
## Subtitle
### Section`
	
	readmePath := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	desc := reader.readFirstLineOfReadme(tmpDir)
	
	if desc != "" {
		t.Errorf("expected empty description for header-only README, got %q", desc)
	}
}

// TestReadFirstLineOfReadme_Truncation tests long lines are truncated
func TestReadFirstLineOfReadme_Truncation(t *testing.T) {
	tmpDir := t.TempDir()
	
	longLine := strings.Repeat("a", 150)
	readmeContent := "# Title\n\n" + longLine
	
	readmePath := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	desc := reader.readFirstLineOfReadme(tmpDir)
	
	if len(desc) > 103 {
		t.Errorf("expected truncated description (max 103 chars with '...'), got %d chars", len(desc))
	}
	if !strings.HasSuffix(desc, "...") {
		t.Error("expected truncated description to end with '...'")
	}
}

// TestReadFirstLineOfReadme_SkipEmptyLines tests skipping empty lines
func TestReadFirstLineOfReadme_SkipEmptyLines(t *testing.T) {
	tmpDir := t.TempDir()
	
	readmeContent := `# Title


First content line after empty lines.`
	
	readmePath := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	desc := reader.readFirstLineOfReadme(tmpDir)
	
	expected := "First content line after empty lines."
	if desc != expected {
		t.Errorf("expected description=%q, got %q", expected, desc)
	}
}

// TestReadRepository_WithDescription tests repository with README description
func TestReadRepository_WithDescription(t *testing.T) {
	tmpDir := t.TempDir()
	
	readmeContent := `# My Project

This is my awesome project.`
	
	readmePath := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	repo, err := reader.ReadRepository(tmpDir, "testorg")
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	expected := "This is my awesome project."
	if repo.Description != expected {
		t.Errorf("expected description=%q, got %q", expected, repo.Description)
	}
}

// TestReadRepository_WithLanguage tests repository with detected language
func TestReadRepository_WithLanguage(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create go.mod file
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}
	
	reader := NewReader(false)
	repo, err := reader.ReadRepository(tmpDir, "testorg")
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	if repo.Language != "Go" {
		t.Errorf("expected language=Go, got %s", repo.Language)
	}
}

// TestNewReader_VerboseFlag tests creating reader with verbose flag
func TestNewReader_VerboseFlag(t *testing.T) {
	reader := NewReader(true)
	
	if !reader.verbose {
		t.Error("expected verbose to be true")
	}
	
	reader = NewReader(false)
	if reader.verbose {
		t.Error("expected verbose to be false")
	}
}
