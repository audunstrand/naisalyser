package cmd

import (
	"fmt"
	"os"

	"github.com/navikt/naisalyser/internal/analyzer"
	"github.com/navikt/naisalyser/internal/github"
	"github.com/navikt/naisalyser/internal/output"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Analyze multiple repositories",
	Long: `Analyze multiple repositories based on a configuration file or topic tag.

Examples:
  naisalyser batch --topic my-product-area
  naisalyser batch --config apps.yaml`,
	RunE: runBatch,
}

func init() {
	batchCmd.Flags().String("topic", "", "GitHub topic tag to filter repositories")
	batchCmd.Flags().String("config", "", "YAML config file with list of repositories")
	batchCmd.Flags().String("org", "navikt", "GitHub organization to search in")
	rootCmd.AddCommand(batchCmd)
}

type BatchConfig struct {
	Repositories []string `yaml:"repositories"`
}

func runBatch(cmd *cobra.Command, args []string) error {
	topic, err := mustGetString(cmd, "topic")
	if err != nil {
		return err
	}
	
	configFile, err := mustGetString(cmd, "config")
	if err != nil {
		return err
	}
	
	org, err := mustGetString(cmd, "org")
	if err != nil {
		return err
	}
	
	outputDir, err := mustGetString(cmd, "output")
	if err != nil {
		return err
	}
	
	verbose, err := mustGetBool(cmd, "verbose")
	if err != nil {
		return err
	}

	gh := github.NewClient(verbose)

	var repos []string

	if topic != "" {
		// Fetch repos by topic
		repoList, err := gh.GetRepositoriesByTopic(org, topic)
		if err != nil {
			return fmt.Errorf("failed to fetch repositories by topic: %w", err)
		}
		repos = repoList
	} else if configFile != "" {
		// Read repos from config file
		data, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		var config BatchConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}
		repos = config.Repositories
	} else {
		return fmt.Errorf("either --topic or --config must be specified")
	}

	fmt.Printf("Found %d repositories to analyze\n", len(repos))

	for i, repo := range repos {
		fmt.Printf("[%d/%d] Analyzing %s...\n", i+1, len(repos), repo)

		repoData, err := gh.GetRepository(repo)
		if err != nil {
			fmt.Printf("  ⚠ Failed to fetch: %v\n", err)
			continue
		}

		analysis, err := analyzer.Analyze(gh, repoData)
		if err != nil {
			fmt.Printf("  ⚠ Failed to analyze: %v\n", err)
			continue
		}

		repoOutputDir := fmt.Sprintf("%s/%s", outputDir, repo)
		if err := output.GenerateMarkdown(analysis, repoOutputDir); err != nil {
			fmt.Printf("  ⚠ Failed to generate output: %v\n", err)
			continue
		}

		fmt.Printf("  ✓ Done\n")
	}

	return nil
}
