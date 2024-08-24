package cmd

import "github.com/spf13/cobra"

var FDSNDateTimeFlag string
var FDSNMagFlag string
var FDSNFormatFlag string

func init() {
	rootCmd.AddCommand(fdsnCmd)
	fdsnCmd.PersistentFlags().StringVarP(&FDSNMagFlag, "magnitude", "m", "", `magnitude or magnitude range (e.g. low[,high] "2.3,4.5")`)
	fdsnCmd.PersistentFlags().StringVarP(&FDSNDateTimeFlag, "time", "t", "", `UTC datetime range (e.g. startdate,enddate "2024-09-20,2024-09-21")`)
	fdsnCmd.PersistentFlags().StringVarP(&FDSNFormatFlag, "output", "o", "table", "output format options: {csv, json, table, text}")
}

var fdsnCmd = &cobra.Command{
	Use:   "fdsn",
	Short: "query historical earthquake records from the FDSN",
	Long: `Retrieve historical earthquake records from the International
	Federation of Digital Seismograph Networks (FDSN)`,
}
