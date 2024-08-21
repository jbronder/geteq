package cmd

import (
	"fmt"

	"github.com/jbronder/geteq/logic"
	"github.com/spf13/cobra"
)

var RtFormatFlag string
var RtMagFlag string
var RtTimeFlag string

func init() {
	rootCmd.AddCommand(realtimeCmd)
	realtimeCmd.Flags().StringVarP(&RtFormatFlag, "output", "o", "table", "output format options: {csv, json, table}")
	realtimeCmd.Flags().StringVarP(&RtMagFlag, "mag", "m", "major", "magnitude options: {all, 1.0, 2.5, 4.5, major}")
	realtimeCmd.Flags().StringVarP(&RtTimeFlag, "time", "t", "month", "time range options: {hour, day, week, month}")
}

var realtimeCmd = &cobra.Command{
	Use:     "realtime",
	Aliases: []string{"real", "rt"},
	Short:   "query real-time earthquake data",
	RunE: func(cmd *cobra.Command, args []string) error {
		fileEndpoint, err := logic.ExtractRTParams(RtFormatFlag, RtMagFlag, RtTimeFlag)
		if err != nil {
			return err
		}
		content, err := logic.RequestContent(fileEndpoint)
		if err != nil {
			return err
		}

		// Standard output format
		switch RtFormatFlag {
		case "table":
			features, err := logic.ExtractFeatures(content)
			if err != nil {
				return err
			}
			logic.StdoutFeatures(features)
		case "csv":
			fmt.Println(string(content))
		case "json":
			fmt.Println(string(content))
		}
		return nil
	},
}
