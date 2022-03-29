package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/yammerjp/db-time-traveler/system"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update date related columns",
	Long:  `Update date related columns`,
	Run: func(cmd *cobra.Command, args []string) {
		update(false)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	initConnection(updateCmd.Flags())
	initUpdate(updateCmd.Flags())
}

func initUpdate(f *flag.FlagSet) {
	f.StringVarP(&table, "table", "", "", "Table name")
	f.BoolVarP(&printQuery, "print-query", "", false, "Print query")
	f.StringVarP(&past, "past", "", "", "rewind date/time")
	f.StringVarP(&future, "future", "", "", "fast forward date/time")
	f.StringVarP(&primaryKeysWhereIn, "primary-keys-where-in", "", "", "Primary Key to specify WHERE IN")
	f.StringSliceVarP(&ignoreColumns, "ignore", "", []string{}, "Ignore columns from updating")
}
func update(dryRun bool) {
	c, err := connection()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	interval, err := system.ParseInterval(past, future)
	if err != nil {
		log.Fatal(err)
	}

	beforeAndAfter, err := c.SelectToUpdateToString(table, *interval, primaryKeysWhereIn, ignoreColumns)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(beforeAndAfter)

	if printQuery {
		query, err := c.UpdateQueryBuilder(table, *interval, primaryKeysWhereIn, ignoreColumns)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(query)
	}
	if dryRun {
		return
	}

	if err := c.Update(table, *interval, primaryKeysWhereIn, ignoreColumns); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nUpdated successfully!")
}
