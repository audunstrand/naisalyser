package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps the gh CLI for GitHub API access
type Client struct {
	verbose bool
}

// Repository represents a GitHub repository
type Repository struct {
	Name          string   `json:"name"`
	FullName      string   `json:"full_name"`
	Description   string   `json:"description"`
	URL           string   `json:"html_url"`
	DefaultBranch string   `json:"default_branch"`
	Topics        []string `json:"topics"`
	Language      string   `json:"language"`
	Archived      bool     `json:"archived"`
}

// FileContent represents file content from GitHub
type FileContent struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
	SHA     string `json:"sha"`
}

// NewClient creates a new GitHub client
func NewClient(verbose bool) *Client {
	return &Client{verbose: verbose}
}

// execGH runs a gh command and returns the output
func (c *Client) execGH(args ...string) ([]byte, error) {
	if c.verbose {
		fmt.Printf("  → gh %s\n", strings.Join(args, " "))
	}

	cmd := exec.Command("gh", args...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("gh command failed: %s", string(exitErr.Stderr))
		}
		return nil, err
	}
	return output, nil
}

// GetRepository fetches repository metadata
func (c *Client) GetRepository(fullName string) (*Repository, error) {
	output, err := c.execGH("api", fmt.Sprintf("repos/%s", fullName))
	if err != nil {
		return nil, err
	}

	var repo Repository
	if err := json.Unmarshal(output, &repo); err != nil {
		return nil, fmt.Errorf("failed to parse repository data: %w", err)
	}

	return &repo, nil
}

// GetRepositoriesByTopic fetches all repositories with a specific topic
func (c *Client) GetRepositoriesByTopic(org, topic string) ([]string, error) {
	query := fmt.Sprintf("org:%s topic:%s", org, topic)
	output, err := c.execGH("api", "search/repositories",
		"-f", fmt.Sprintf("q=%s", query),
		"-f", "per_page=100",
		"--jq", ".items[].full_name")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var repos []string
	for _, line := range lines {
		if line != "" {
			repos = append(repos, line)
		}
	}

	return repos, nil
}

// GetFileContent fetches a file's content from a repository
func (c *Client) GetFileContent(fullName, path string) (string, error) {
	output, err := c.execGH("api", fmt.Sprintf("repos/%s/contents/%s", fullName, path),
		"--jq", ".content")
	if err != nil {
		return "", err
	}

	// Content is base64 encoded
	encoded := strings.TrimSpace(string(output))
	if encoded == "" || encoded == "null" {
		return "", fmt.Errorf("file not found: %s", path)
	}

	// Decode base64
	decoded, err := decodeBase64(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode file content: %w", err)
	}

	return decoded, nil
}

// GetDirectoryContents lists files in a directory
func (c *Client) GetDirectoryContents(fullName, path string) ([]string, error) {
	output, err := c.execGH("api", fmt.Sprintf("repos/%s/contents/%s", fullName, path),
		"--jq", ".[].name")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []string
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}

	return files, nil
}

// FileExists checks if a file exists in the repository
func (c *Client) FileExists(fullName, path string) bool {
	_, err := c.execGH("api", fmt.Sprintf("repos/%s/contents/%s", fullName, path),
		"--jq", ".name")
	return err == nil
}
