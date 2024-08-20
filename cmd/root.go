package cmd

import (
	_"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "geteq",
	Short: "geteq returns real-time earthquake data from USGS",
	Long:  `geteq: A CLI tool to obtain real-time earthquake data from USGS in
	multiple formats including the terminal in a human readable format.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// throws the error at the end after printing the Usage prompt
		// fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
