package cmd

import (
	"github.com/spf13/cobra"
)

// NewCLI creates a new instance of the root CLI
func NewCLI() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tracee-polr-adapter",
		Short: "Generates PolicyReports for Tracee Events",
	}

	rootCmd.AddCommand(newRunCMD())

	return rootCmd
}
