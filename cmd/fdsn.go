package cmd

import "github.com/spf13/cobra"

var FDSNDateTimeFlag string
var FDSNMagFlag string
var FDSNFormatFlag string

func init() {
	rootCmd.AddCommand(fdsnCmd)
	fdsnCmd.PersistentFlags().StringVarP(&FDSNMagFlag, "magnitude", "m", "", "a minimum magnitude")
	fdsnCmd.PersistentFlags().StringVarP(&FDSNDateTimeFlag, "time", "t", "", "a UTC datetime")
	fdsnCmd.PersistentFlags().StringVarP(&FDSNFormatFlag, "output", "o", "table", "output format options: {json, table, text}")
}

var fdsnCmd = &cobra.Command{
	Use:   "fdsn",
	Short: "query historical earthquake records from the FDSN",
	Long: `Retreive earthquake records from the International Federation
	of Digital Seismograph Networks (FDSN)`,
}
