package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/yammerjp/db-time-traveler/configuration"
)

var initConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "Create sample config files to ~/.db-time-traveler.yaml",
	Long:  "Create sample config files to ~/.db-time-traveler.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		err := configuration.CreateConfigFile(configPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initConfigCmd)
	initShowConfigFlg(initConfigCmd.Flags())
}
