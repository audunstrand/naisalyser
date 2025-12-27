package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// mustGetString retrieves a string flag or returns an error
func mustGetString(cmd *cobra.Command, name string) (string, error) {
	val, err := cmd.Flags().GetString(name)
	if err != nil {
		return "", fmt.Errorf("failed to get flag %q: %w", name, err)
	}
	return val, nil
}

// mustGetBool retrieves a bool flag or returns an error
func mustGetBool(cmd *cobra.Command, name string) (bool, error) {
	val, err := cmd.Flags().GetBool(name)
	if err != nil {
		return false, fmt.Errorf("failed to get flag %q: %w", name, err)
	}
	return val, nil
}
