package cmd

import (
	"github.com/jonstjohn/crdb-settings/pkg/dbpgx"
	"github.com/jonstjohn/crdb-settings/pkg/releases"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "save",
	Short: "Releases save command",
	Run: func(cmd *cobra.Command, args []string) {
		pool, err := dbpgx.NewPoolFromUrl(urlArg)
		if err != nil {
			panic(err)
		}
		err = releases.UpdateReleases(pool)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
