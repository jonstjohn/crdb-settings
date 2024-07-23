package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jonstjohn/crdb-settings/pkg/settings"
	"github.com/spf13/cobra"
)

var listSettingsVersionFlag string

var settingsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cluster settings for a specific release",
	Run: func(cmd *cobra.Command, args []string) {

		s, err := settings.NewSettingsManager(urlArg)
		if err != nil {
			panic(err)
		}

		sts, err := s.GetSettingsForRelease(listSettingsVersionFlag)

		if err != nil {
			panic(err)
		}
		b, err := json.MarshalIndent(sts, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))

		/*

			pool, err := dbpgx.NewPoolFromUrl(urlArg)
			if err != nil {
				panic(err)
			}
			s := settings.NewDbDatasource(pool)
			releases, err := s.GetRawSettingsForVersion(listSettingsVersionFlag)
			if err != nil {
				panic(err)
			}
			b, err := json.MarshalIndent(releases, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))

		*/
	},
}

func init() {
	settingsCmd.AddCommand(settingsListCmd)
	settingsListCmd.Flags().StringVar(&listSettingsVersionFlag, "version", "v23.2.1", "CRDB version, starting with 'v'")
}
