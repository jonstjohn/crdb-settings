package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jonstjohn/crdb-settings/pkg/settings"
	"github.com/spf13/cobra"
)

var settingDetailSettingFlag string

var settingsDetailCmd = &cobra.Command{
	Use:   "detail",
	Short: "Settings detail command",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := settings.NewSettingsManager(urlArg)
		if err != nil {
			panic(err)
		}
		detail, err := m.GetSettingDetail(settingDetailSettingFlag)
		if err != nil {
			panic(err)
		}
		b, err := json.MarshalIndent(detail, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	},
}

func init() {
	settingsCmd.AddCommand(settingsDetailCmd)
	settingsDetailCmd.Flags().StringVar(&settingDetailSettingFlag, "setting", "changefeed.random_replica_selection.enabled", "Setting to get details for")
}
