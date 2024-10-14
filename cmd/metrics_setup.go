package cmd

import (
	"github.com/jonstjohn/crdb-settings/pkg/metrics"
	"github.com/spf13/cobra"
)

var metricsSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup metrics database and tables",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := metrics.NewManager(urlArg)
		if err != nil {
			panic(err)
		}
		if err = m.InitializeDatabase(); err != nil {
			panic(err)
		}
	},
}

func init() {
	metricsCmd.AddCommand(metricsSetupCmd)
}
