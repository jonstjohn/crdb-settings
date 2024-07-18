package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jonstjohn/crdb-settings/pkg/dbpgx"
	"github.com/jonstjohn/crdb-settings/pkg/settings"
	"github.com/spf13/cobra"
)

var listSettingsVersionFlag string

var settingsListCmd = &cobra.Command{
	Use:   "settings list",
	Short: "Run a local test server and output the settings",
	Run: func(cmd *cobra.Command, args []string) {
		pool, err := dbpgx.NewPoolFromUrl(urlArg)
		if err != nil {
			panic(err)
		}
		s := settings.NewDbDatasource(pool)
		releases, err := s.GetSettingsForVersion(listSettingsVersionFlag)
		if err != nil {
			panic(err)
		}
		b, err := json.MarshalIndent(releases, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	},
}

func init() {
	rootCmd.AddCommand(settingsListCmd)
	settingsListCmd.Flags().StringVar(&listSettingsVersionFlag, "version", "v23.2.1", "CRDB version, starting with 'v'")
}
