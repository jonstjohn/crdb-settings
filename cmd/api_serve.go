package cmd

import (
	"github.com/jonstjohn/crdb-settings/pkg/api"
	"github.com/spf13/cobra"
)

var apiServeCmd = &cobra.Command{
	Use:   "api serve",
	Short: "Run a local test server and output the settings",
	Run: func(cmd *cobra.Command, args []string) {
		api.Serve(urlArg)
	},
}

func init() {
	rootCmd.AddCommand(apiServeCmd)
}
