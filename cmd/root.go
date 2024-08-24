package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "geteq",
	Short: "geteq returns real-time and historical earthquake records from USGS",
	Long: `geteq: A CLI tool to obtain real-time and historical earthquake records from USGS in
	multiple formats including the terminal in a tabular format.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Errors are caught from individual cobra.Command RunE functions
		os.Exit(1)
	}
}
