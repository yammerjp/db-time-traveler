package cmd

import (
	"fmt"
	"log"
	"os"

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

	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&username, "user", "u", "root", "Username for Database Connection")
	updateCmd.Flags().StringVarP(&password, "password", "p", "password", "Password for Database Connection")
	updateCmd.Flags().StringVarP(&host, "host", "", "localhost", "Hostname or IPv4 address for Database Connection")
	updateCmd.Flags().IntVarP(&port, "port", "", 3306, "Port number for Database Connection")
	updateCmd.Flags().StringVarP(&schema, "schema", "s", "", "Schema name for Database Connection")
	updateCmd.Flags().StringVarP(&table, "table", "", "", "Table name")
	updateCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Dry run")
	updateCmd.Flags().StringVarP(&past, "past", "", "", "rewind date/time")
	updateCmd.Flags().StringVarP(&primaryKeyRaw, "primary-key-raw", "", "", "Primary Key to specify WHERE IN")
	updateCmd.Flags().StringVarP(&sshHost, "ssh-host", "", "", "Host name for bastion SSH host")
	updateCmd.Flags().IntVarP(&sshPort, "ssh-port", "", 22, "Host port number for bastion SSH host")
	updateCmd.Flags().StringVarP(&sshUser, "ssh-user", "", "", "Host username for bastion SSH host")
	updateCmd.Flags().StringVarP(&sshKeyPath, "ssh-key-path", "", home+"/.ssh/id_rsa", "Private key path for bastion SSH host")
	updateCmd.Flags().StringVarP(&sshPassphrase, "ssh-passphrase", "", "", "Private key passphrase for bastion SSH host")
}
