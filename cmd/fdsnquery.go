package cmd

import (
	"fmt"

	"github.com/jbronder/geteq/logic"
	"github.com/spf13/cobra"
)

func init() {
	fdsnCmd.AddCommand(queryCmd)
}

var queryCmd = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q"},
	Short:   "run a record query",
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint, err := logic.ExtractFDSNParams("query", FDSNMagFlag, FDSNFormatFlag, FDSNDateTimeFlag)
		if err != nil {
			return err
		}

		content, err := logic.RequestContent(endpoint)
		if err != nil {
			return err
		}

		switch FDSNFormatFlag {
		case "table":
			features, err := logic.ExtractFeatures(content)
			if err != nil {
				return err
			}
			logic.StdoutFeatures(features)
		case "json":
			fallthrough
		case "csv":
			fallthrough
		case "text":
			fmt.Println(string(content))
		}
		return nil
	},
}
