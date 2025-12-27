package analyzer

// RepositoryReader abstracts file reading from GitHub or local filesystem
type RepositoryReader interface {
	// GetFileContent reads a file's content. repoPath is the repo identifier
	// (full name for GitHub, directory path for local).
	GetFileContent(repoPath, filePath string) (string, error)

	// FileExists checks if a file exists in the repository
	FileExists(repoPath, filePath string) bool
}
