package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jonstjohn/crdb-settings/pkg/releases"
	"github.com/spf13/cobra"
)

var releasesListCmdSourceArg string

var releasesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Releases list command",
	Run: func(cmd *cobra.Command, args []string) {
		if releasesListCmdSourceArg == "db" {
			rp, err := releases.NewDbDatasource(urlArg)
			if err != nil {
				panic(err)
			}
			releases, err := rp.GetReleases()
			if err != nil {
				panic(err)
			}
			b, err := json.MarshalIndent(releases, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		} else {
			rp := releases.NewRemoteDataSource()
			releases, err := rp.GetReleases()
			if err != nil {
				panic(err)
			}
			b, err := json.MarshalIndent(releases, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		}
	},
}

func init() {
	releasesListCmd.Flags().StringVar(&releasesListCmdSourceArg, "source", "db", "Source for releases list command - 'yaml' or 'db'")
	releasesCmd.AddCommand(releasesListCmd)

}
