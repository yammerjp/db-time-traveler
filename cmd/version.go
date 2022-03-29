package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show db-time-traveler version",
	Long:  "Show db-time-traveler version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(builtInformations)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
