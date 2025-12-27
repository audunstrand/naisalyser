package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "naisalyser",
	Short: "Analyze and document NAV applications",
	Long: `NAISalyser is a CLI tool for analyzing and generating documentation
for applications in the NAV ecosystem.

It fetches data from GitHub repositories and produces markdown documentation.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("output", "o", "./docs/generated", "Output directory for generated docs")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
}
