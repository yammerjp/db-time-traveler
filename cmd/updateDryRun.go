package cmd

import (
	"github.com/spf13/cobra"
)

var updateDryRunCmd = &cobra.Command{
	Use:   "update-dry-run",
	Short: "Dry run to update date related columns",
	Long:  `Dry run to update date related columns`,
	Run: func(cmd *cobra.Command, args []string) {
		update(true)
	},
}

func init() {
	rootCmd.AddCommand(updateDryRunCmd)
	initConnection(updateDryRunCmd.Flags())
	initUpdate(updateDryRunCmd.Flags())
}
