package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/yammerjp/db-time-traveler/system"
)

var showConfigCmd = &cobra.Command{
	Use:   "update",
	Short: "show configs",
	Long:  `show the inner of config yaml file on ~/.db-time-traveler.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		if schema == "" {
			log.Fatal("Need --schema option")
		}
		c, err := connection()
		if err != nil {
			panic(err)
		}
		defer c.Close()

		parsedPast, err := system.ParsePast(past)
		if err != nil {
			log.Fatal(err)
		}

		if dryRun {
			primaryKeys, err := c.SelectPrimaryKeyColumns(table)
			if err != nil {
				log.Fatal(err)
			}
			if len(primaryKeys) != 1 {
				log.Fatal("Multiple column primary keys is not supported")
			}
			primaryKeyValues, columns, columnValuesBefore, columnValuesAfter, err := c.SelectToUpdate(table, parsedPast, primaryKeyRaw)
			for i := range primaryKeyValues {
				for j := range columns {
					fmt.Printf("%s: %s\n  %s:\n    before: %s\n    after:  %s\n", primaryKeys[0], primaryKeyValues[i], columns[j], columnValuesBefore[i][j], columnValuesAfter[i][j])
				}
			}
		} else {
			if err := c.Update(table, "1 MONTH", primaryKeyRaw); err != nil {
				log.Fatal(err)
			}
			fmt.Println("updated successfully")
		}
	},
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.AddCommand(showConfigCmd)
	showConfigCmd.Flags().StringVarP(&configPath, "config-path", "", home+"/.db-time-traveler.yml", "config path")
}
