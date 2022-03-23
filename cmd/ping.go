package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yammerjp/db-time-traveler/system"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Try to connect to a relational database",
	Long:  `Try to connect to a relational database`,
	Run: func(cmd *cobra.Command, args []string) {
		system.Ping()
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
