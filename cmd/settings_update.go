package cmd

import (
	"github.com/jonstjohn/crdb-settings/pkg/settings"
	"github.com/spf13/cobra"
)

var saveSettingsReleaseFlag string

var settingsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Settings update command",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := settings.NewSettingsManager(urlArg)
		if err != nil {
			panic(err)
		}
		err = s.SaveClusterSettingsForVersion(saveSettingsReleaseFlag, urlArg)
		if err != nil {
			panic(err)
		}
		/*
			err := settings.SaveClusterSettingsForVersion(saveSettingsReleaseFlag, urlArg)
			if err != nil {
				panic(err)
			}

		*/
	},
}

func init() {
	settingsCmd.AddCommand(settingsUpdateCmd)
	settingsUpdateCmd.Flags().StringVar(&saveSettingsReleaseFlag, "release", "all", "Update all or specify a single CRDB release, starting with 'v'")
}
