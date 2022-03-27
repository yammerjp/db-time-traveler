package cmd

import (
	"errors"
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
    if !dryRun {
      log.Fatal("real run is not supported yet")
    }

    whereClauses := []system.WhereClause{}
    for _, v := range wheres {

      whereClause, err := system.ParseWhereClause(v)
      if err != nil {
        log.Fatal(err)
      }
      if whereClause.Operator != "=" {
        log.Fatal(errors.New("Unsupported where clause's operator"))
      }
      whereClauses = append(whereClauses, *whereClause)
    }
    if past != "1month" {
      log.Fatal("past without 1month is not supported yet")
    }

		dap := &system.DatabaseAccessPoint{
			Username: username,
			Password: password,
			Host:     host,
			Port:     port,
			Schema:   schema,
		}

    db, err := dap.Connect()
    if err != nil {
      log.Fatal(err)
    }
    defer db.Close()

    trcs, err := system.SelectColumns(db, table)
    if err != nil {
      log.Fatal(err)
    }
    // if err := system.SelectToUpdate(db, table, trcs, whereClauses); err != nil {
    if err := system.Update(db, table, trcs, whereClauses); err != nil {
      log.Fatal(err)
    }
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&username, "user", "u", "root", "Username for Database Connection")
	updateCmd.Flags().StringVarP(&password, "password", "p", "password", "Password for Database Connection")
	updateCmd.Flags().StringVarP(&host, "host", "", "localhost", "Hostname or IPv4 address for Database Connection")
	updateCmd.Flags().IntVarP(&port, "port", "", 3306, "Port number for Database Connection")
	updateCmd.Flags().StringVarP(&schema, "schema", "s", "", "Schema name for Database Connection")
	updateCmd.Flags().StringVarP(&table, "table", "", "", "Table name")
  updateCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Dry run")
  updateCmd.Flags().StringSliceVarP(&wheres, "where", "w", []string{}, "Conditions to refine")
  updateCmd.Flags().StringVarP(&past, "past", "", "", "rewind date/time")
}
