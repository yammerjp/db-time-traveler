package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/yammerjp/db-time-traveler/system"
)

var findDateRelatedColumnsCmd = &cobra.Command{
	Use:   "find-date-related-columns",
	Short: "Show date related column names by the specified table name.",
	Long:  `Show date related column names by the specified table name. (Need --table option)`,
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

		columns, err := dap.SelectDateRelatedColumns(table)
		if err != nil {
			log.Fatal("returned error:" + err.Error())
		}
		for _, v := range columns {
			fmt.Println(v)
		}
	},
}

func init() {
	rootCmd.AddCommand(findDateRelatedColumnsCmd)
	findDateRelatedColumnsCmd.Flags().StringVarP(&username, "user", "u", "root", "Username for Database Connection")
	findDateRelatedColumnsCmd.Flags().StringVarP(&password, "password", "p", "password", "Password for Database Connection")
	findDateRelatedColumnsCmd.Flags().StringVarP(&host, "host", "", "localhost", "Hostname or IPv4 address for Database Connection")
	findDateRelatedColumnsCmd.Flags().IntVarP(&port, "port", "", 3306, "Port number for Database Connection")
	findDateRelatedColumnsCmd.Flags().StringVarP(&schema, "schema", "s", "", "Schema name for Database Connection")
	findDateRelatedColumnsCmd.Flags().StringVarP(&table, "table", "", "", "Table name")
}
