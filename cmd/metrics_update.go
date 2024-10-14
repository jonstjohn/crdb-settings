package cmd

import (
	"github.com/jonstjohn/crdb-settings/pkg/metrics"
	"github.com/spf13/cobra"
)

var updateMetricsCmdReleaseFlag string

var metricsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update metrics",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := metrics.NewManager(urlArg)
		if err != nil {
			panic(err)
		}
		if err = m.SaveMetricsForRelease(updateMetricsCmdReleaseFlag); err != nil {
			panic(err)
		}
	},
}

func init() {
	metricsCmd.AddCommand(metricsUpdateCmd)
	metricsUpdateCmd.Flags().StringVarP(&updateMetricsCmdReleaseFlag, "release", "r", "", "Release name")
}
