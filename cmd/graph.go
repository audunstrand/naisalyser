package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/navikt/naisalyser/internal/analyzer"
	"github.com/navikt/naisalyser/internal/graph"
	"github.com/navikt/naisalyser/internal/local"
	"github.com/spf13/cobra"
)

var graphCmd = &cobra.Command{
	Use:   "graph [repos-dir]",
	Short: "Generate dependency graph from local repositories",
	Long: `Generate a dependency graph showing relationships between applications.

The graph includes:
- App-to-app communication (from accessPolicy)
- Database dependencies
- Kafka dependencies
- External service dependencies

Example:
  naisalyser graph ./repos --output ./graph`,
	Args: cobra.ExactArgs(1),
	RunE: runGraph,
}

func init() {
	rootCmd.AddCommand(graphCmd)
}

func runGraph(cmd *cobra.Command, args []string) error {
	reposDir := args[0]
	
	outputDir, err := mustGetString(cmd, "output")
	if err != nil {
		return err
	}
	
	verbose, err := mustGetBool(cmd, "verbose")
	if err != nil {
		return err
	}

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

	fmt.Printf("Building dependency graph from %d repositories\n", len(repos))

	reader := local.NewReader(verbose)
	builder := graph.NewBuilder()

	for _, repoPath := range repos {
		repoName := filepath.Base(repoPath)
		if verbose {
			fmt.Printf("  Processing %s...\n", repoName)
		}

		repoData, err := reader.ReadRepository(repoPath, "navikt")
		if err != nil {
			if verbose {
				fmt.Printf("    ⚠ Failed to read: %v\n", err)
			}
			continue
		}

		analysis, err := analyzer.AnalyzeLocal(reader, repoData)
		if err != nil {
			if verbose {
				fmt.Printf("    ⚠ Failed to analyze: %v\n", err)
			}
			continue
		}

		builder.AddFromLocalAnalysis(analysis)
	}

	g := builder.Build()
	fmt.Printf("Found %d nodes and %d edges\n", len(g.Nodes), len(g.Edges))

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write JSON
	jsonPath := filepath.Join(outputDir, "dependencies.json")
	jsonData, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}
	fmt.Printf("✓ Written %s\n", jsonPath)

	// Write Mermaid markdown
	mermaidPath := filepath.Join(outputDir, "dependencies.md")
	mermaidContent := fmt.Sprintf("# Dependency Graph\n\n```mermaid\n%s```\n", g.ToMermaid())
	if err := os.WriteFile(mermaidPath, []byte(mermaidContent), 0644); err != nil {
		return fmt.Errorf("failed to write Mermaid: %w", err)
	}
	fmt.Printf("✓ Written %s\n", mermaidPath)

	return nil
}
