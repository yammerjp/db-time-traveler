package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/yammerjp/db-time-traveler/system"
)

var showConfigCmd = &cobra.Command{
	Use:   "show-config",
	Short: "show configs",
	Long:  `show the inner of config yaml such as ~/.db-time-traveler.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := system.LoadConfig(configPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(config.ToString())
	},
}

func init() {
	rootCmd.AddCommand(showConfigCmd)
	showConfigCmd.Flags().StringVarP(&configPath, "config-path", "", "", "config path")
	// showConftgCmd.Flags().StringVarP(&selectedConnection, "connection", "", "", "connection name")
}
