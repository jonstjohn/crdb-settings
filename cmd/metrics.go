package cmd

import "github.com/spf13/cobra"

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Metrics commands",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(metricsCmd)
}
