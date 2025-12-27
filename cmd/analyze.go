package cmd

import (
	"fmt"

	"github.com/navikt/naisalyser/internal/analyzer"
	"github.com/navikt/naisalyser/internal/github"
	"github.com/navikt/naisalyser/internal/output"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [repo]",
	Short: "Analyze a single repository",
	Long: `Analyze a GitHub repository and generate documentation.

Example:
  naisalyser analyze navikt/my-app
  naisalyser analyze navikt/my-app --output ./docs`,
	Args: cobra.ExactArgs(1),
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	repo := args[0]
	
	outputDir, err := mustGetString(cmd, "output")
	if err != nil {
		return err
	}
	
	verbose, err := mustGetBool(cmd, "verbose")
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Analyzing repository: %s\n", repo)
	}

	// Fetch repository data using gh CLI
	gh := github.NewClient(verbose)
	repoData, err := gh.GetRepository(repo)
	if err != nil {
		return fmt.Errorf("failed to fetch repository: %w", err)
	}

	// Analyze the repository
	analysis, err := analyzer.Analyze(gh, repoData)
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
