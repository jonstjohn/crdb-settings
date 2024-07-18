package cmd

import (
	"github.com/jonstjohn/crdb-settings/pkg/settings"
	"github.com/spf13/cobra"
)

var saveSettingsVersionFlag string

var settingsUpdateCmd = &cobra.Command{
	Use:   "settings update",
	Short: "Settings update command",
	Run: func(cmd *cobra.Command, args []string) {
		err := settings.SaveClusterSettingsForVersion(saveSettingsVersionFlag, urlArg)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(settingsUpdateCmd)
	settingsUpdateCmd.Flags().StringVar(&saveSettingsVersionFlag, "version", "v23.2.1", "Specify a single CRDB version, starting with 'v'")
}
