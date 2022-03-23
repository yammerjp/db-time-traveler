package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
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
		fmt.Println("update called")
		db, err := sql.Open("mysql", "root:password@/sampleschema")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
		rows, err := db.Query("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = 'accounts'")
		if err != nil {
			log.Fatal(err)
		}

		columns := make([]string, 0)
		dataTypes := make([]string, 0)

		for rows.Next() {
			var column string
			var dataType string
			if err := rows.Scan(&column, &dataType); err != nil {
				log.Fatal(err)
			}
			columns = append(columns, column)
			dataTypes = append(dataTypes, dataType)
		}

		for i, column := range columns {
			fmt.Printf("%40s %40s\n", column, dataTypes[i])
		}

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
