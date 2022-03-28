package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const DbTimeTravelerVersion = "v0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show db-time-traveler version",
	Long:  "Show db-time-traveler version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(DbTimeTravelerVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
