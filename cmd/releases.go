package cmd

import "github.com/spf13/cobra"

var releasesCmd = &cobra.Command{
	Use:   "releases",
	Short: "Releases commands",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(releasesCmd)
}
