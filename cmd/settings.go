package cmd

import "github.com/spf13/cobra"

var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Settings commands",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)
}
