/*
Copyright © 2024 Jon St John <jon@element128.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var urlArg string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crdb-settings",
	Short: "View and compare CockroachDB cluster settings.",
	Long: `CRDB settings is a CLI library for GO that enables storing and analyzing CockroachDB cluster
settings across versions.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&urlArg, "url", os.Getenv("CRDB_SETTINGS_URL"), "Database URL")
	rootCmd.MarkFlagRequired("url")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
