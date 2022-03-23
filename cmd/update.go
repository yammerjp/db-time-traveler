package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/yammerjp/db-time-traveler/system"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if schema == "" {
			log.Fatal("Need --schema option")
		}

		dap := &system.DatabaseAccessPoint{
			Username: username,
			Password: password,
			Host:     host,
			Port:     port,
			Schema:   schema,
		}

		system.SelectColumns(dap)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&username, "user", "u", "root", "Username for Database Connection")
	updateCmd.Flags().StringVarP(&password, "password", "p", "password", "Password for Database Connection")
	updateCmd.Flags().StringVarP(&host, "host", "", "localhost", "Hostname or IPv4 address for Database Connection")
	updateCmd.Flags().IntVarP(&port, "port", "", 3306, "Port number for Database Connection")
	updateCmd.Flags().StringVarP(&schema, "schema", "s", "", "Schema name for Database Connection")
}
