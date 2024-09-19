package cmd

import (
	"github.com/jonstjohn/crdb-settings/pkg/gh"
	"github.com/spf13/cobra"
)

var githubSettingFlag string
var githubCmdAccessTokenFlag string

var settingsGithubCmd = &cobra.Command{
	Use:   "github",
	Short: "Settings github command",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := gh.NewManager(&githubCmdAccessTokenFlag, urlArg)
		if err != nil {
			panic(err)
		}
		err = m.UpdateIssuesForSetting(githubSettingFlag)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	settingsCmd.AddCommand(settingsGithubCmd)
	settingsGithubCmd.Flags().StringVar(&githubSettingFlag, "setting", "all", "Setting to search github for")
	settingsGithubCmd.Flags().StringVar(&githubCmdAccessTokenFlag, "token", "", "Github access token, optional")
}
