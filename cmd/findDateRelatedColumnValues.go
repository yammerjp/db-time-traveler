package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/yammerjp/db-time-traveler/system"
)

var findDateRelatedColumnValuesCmd = &cobra.Command{
	Use:   "find-date-related-column-values",
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

		columns, columnValues, err := dap.SelectDateRelatedColumnValues(table, primaryKeyRaw)
		if err != nil {
			log.Fatal("returned error:" + err.Error())
		}
		for _, v := range columnValues {
			for j, r := range columns {
				if j != 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s: %s", r, v[j])
			}
			fmt.Print("\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(findDateRelatedColumnValuesCmd)
	findDateRelatedColumnValuesCmd.Flags().StringVarP(&username, "user", "u", "root", "Username for Database Connection")
	findDateRelatedColumnValuesCmd.Flags().StringVarP(&password, "password", "p", "password", "Password for Database Connection")
	findDateRelatedColumnValuesCmd.Flags().StringVarP(&host, "host", "", "localhost", "Hostname or IPv4 address for Database Connection")
	findDateRelatedColumnValuesCmd.Flags().IntVarP(&port, "port", "", 3306, "Port number for Database Connection")
	findDateRelatedColumnValuesCmd.Flags().StringVarP(&schema, "schema", "s", "", "Schema name for Database Connection")
	findDateRelatedColumnValuesCmd.Flags().StringVarP(&table, "table", "", "", "Table name")
	findDateRelatedColumnValuesCmd.Flags().StringVarP(&primaryKeyRaw, "primary-key-raw", "", "", "Primary Key to specify WHERE IN")
}
