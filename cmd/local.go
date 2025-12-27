package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/navikt/naisalyser/internal/analyzer"
	"github.com/navikt/naisalyser/internal/local"
	"github.com/navikt/naisalyser/internal/output"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local [path]",
	Short: "Analyze a local repository",
	Long: `Analyze a locally cloned repository and generate documentation.

Example:
  naisalyser local ./repos/my-app
  naisalyser local ./repos/my-app --output ./docs`,
	Args: cobra.ExactArgs(1),
	RunE: runLocal,
}

var localBatchCmd = &cobra.Command{
	Use:   "local-batch [repos-dir]",
	Short: "Analyze all repositories in a directory",
	Long: `Analyze all repositories in a directory and generate documentation.

Example:
  naisalyser local-batch ./repos --output ./pleiepenger`,
	Args: cobra.ExactArgs(1),
	RunE: runLocalBatch,
}

func init() {
	rootCmd.AddCommand(localCmd)
	rootCmd.AddCommand(localBatchCmd)
}

func runLocal(cmd *cobra.Command, args []string) error {
	repoPath := args[0]
	outputDir, _ := cmd.Flags().GetString("output")
	verbose, _ := cmd.Flags().GetBool("verbose")

	if verbose {
		fmt.Printf("Analyzing local repository: %s\n", repoPath)
	}

	// Read local repository
	reader := local.NewReader(verbose)
	repoData, err := reader.ReadRepository(repoPath)
	if err != nil {
		return fmt.Errorf("failed to read repository: %w", err)
	}

	// Analyze the repository
	analysis, err := analyzer.AnalyzeLocal(reader, repoData)
	if err != nil {
		return fmt.Errorf("failed to analyze repository: %w", err)
	}

	// Generate output
	if err := output.GenerateMarkdown(analysis, outputDir); err != nil {
		return fmt.Errorf("failed to generate output: %w", err)
	}

	fmt.Printf("✓ Documentation generated in %s\n", outputDir)
	return nil
}

func runLocalBatch(cmd *cobra.Command, args []string) error {
	reposDir := args[0]
	outputDir, _ := cmd.Flags().GetString("output")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// List subdirectories
	entries, err := os.ReadDir(reposDir)
	if err != nil {
		return fmt.Errorf("failed to read repos directory: %w", err)
	}

	var repos []string
	for _, entry := range entries {
		if entry.IsDir() && !isHiddenDir(entry.Name()) {
			repos = append(repos, filepath.Join(reposDir, entry.Name()))
		}
	}

	fmt.Printf("Found %d repositories to analyze\n", len(repos))

	reader := local.NewReader(verbose)

	for i, repoPath := range repos {
		repoName := filepath.Base(repoPath)
		fmt.Printf("[%d/%d] Analyzing %s...\n", i+1, len(repos), repoName)

		repoData, err := reader.ReadRepository(repoPath)
		if err != nil {
			fmt.Printf("  ⚠ Failed to read: %v\n", err)
			continue
		}

		analysis, err := analyzer.AnalyzeLocal(reader, repoData)
		if err != nil {
			fmt.Printf("  ⚠ Failed to analyze: %v\n", err)
			continue
		}

		repoOutputDir := filepath.Join(outputDir, repoName)
		if err := output.GenerateMarkdown(analysis, repoOutputDir); err != nil {
			fmt.Printf("  ⚠ Failed to generate output: %v\n", err)
			continue
		}

		fmt.Printf("  ✓ Done\n")
	}

	return nil
}

func isHiddenDir(name string) bool {
	return len(name) > 0 && name[0] == '.'
}
