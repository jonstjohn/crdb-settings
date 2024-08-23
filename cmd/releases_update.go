package cmd

import (
	"github.com/jonstjohn/crdb-settings/pkg/releases"
	"github.com/spf13/cobra"
)

var releasesUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update db releases from remote yaml",
	Run: func(cmd *cobra.Command, args []string) {
		rm, err := releases.NewReleasesManager(urlArg)
		if err != nil {
			panic(err)
		}
		err = rm.UpdateReleases()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	releasesCmd.AddCommand(releasesUpdateCmd)
}
