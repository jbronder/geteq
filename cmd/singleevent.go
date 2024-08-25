package cmd

import (
	"fmt"

	"github.com/jbronder/geteq/logic"
	"github.com/spf13/cobra"
)

func init() {
	queryCmd.AddCommand(singleEventCmd)
}

var singleEventCmd = &cobra.Command{
	Use:     "event",
	Aliases: []string{"se", "e", "s"},
	Short:   "Detailed information about a single event given an eventid",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		endpoint, err := logic.ExtractId("query", FDSNFormatFlag, args[0])
		if err != nil {
			return err
		}

		content, err := logic.RequestContent(endpoint)
		if err != nil {
			return err
		}

		switch FDSNFormatFlag {
		case "table":
			feature, err := logic.ExtractSingleFeature(content)
			if err != nil {
				return err
			}
			logic.StdoutSingleEvent(feature)
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
