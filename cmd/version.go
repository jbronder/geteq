package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of geteq",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("geteq CLI v0.1 {Development Version}")
	},
}
