package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/yammerjp/db-time-traveler/system"
)

var findTimeRelatedColumnsCmd = &cobra.Command{
	Use:   "find-time-related-columns",
	Short: "Show time related column names by the specified table name.",
	Long:  `Show time related column names by the specified table name. (Need --table option)`,
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

		system.SelectAndPrintColumns(dap, table)
	},
}

func init() {
	rootCmd.AddCommand(findTimeRelatedColumnsCmd)
	findTimeRelatedColumnsCmd.Flags().StringVarP(&username, "user", "u", "root", "Username for Database Connection")
	findTimeRelatedColumnsCmd.Flags().StringVarP(&password, "password", "p", "password", "Password for Database Connection")
	findTimeRelatedColumnsCmd.Flags().StringVarP(&host, "host", "", "localhost", "Hostname or IPv4 address for Database Connection")
	findTimeRelatedColumnsCmd.Flags().IntVarP(&port, "port", "", 3306, "Port number for Database Connection")
	findTimeRelatedColumnsCmd.Flags().StringVarP(&schema, "schema", "s", "", "Schema name for Database Connection")
	findTimeRelatedColumnsCmd.Flags().StringVarP(&table, "table", "", "", "Table name")
}
