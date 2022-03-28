package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/yammerjp/db-time-traveler/system"
)

var updateDryRunCmd = &cobra.Command{
	Use:   "update-dry-run",
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

		beforeAndAfter, err := c.SelectToUpdateToString(table, parsedPast, primaryKeyRaw)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(beforeAndAfter)

		if printQuery {
			query, err := c.UpdateQueryBuilder(table, parsedPast, primaryKeyRaw)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(query)
		}
	},
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.AddCommand(updateDryRunCmd)
	updateDryRunCmd.Flags().StringVarP(&username, "user", "u", "root", "Username for Database Connection")
	updateDryRunCmd.Flags().StringVarP(&password, "password", "p", "password", "Password for Database Connection")
	updateDryRunCmd.Flags().StringVarP(&host, "host", "", "localhost", "Hostname or IPv4 address for Database Connection")
	updateDryRunCmd.Flags().IntVarP(&port, "port", "", 3306, "Port number for Database Connection")
	updateDryRunCmd.Flags().StringVarP(&schema, "schema", "s", "", "Schema name for Database Connection")
	updateDryRunCmd.Flags().StringVarP(&table, "table", "", "", "Table name")
	updateDryRunCmd.Flags().BoolVarP(&printQuery, "print-query", "", false, "Print query")
	updateDryRunCmd.Flags().StringVarP(&past, "past", "", "", "rewind date/time")
	updateDryRunCmd.Flags().StringVarP(&primaryKeyRaw, "primary-key-raw", "", "", "Primary Key to specify WHERE IN")
	updateDryRunCmd.Flags().StringVarP(&sshHost, "ssh-host", "", "", "Host name for bastion SSH host")
	updateDryRunCmd.Flags().IntVarP(&sshPort, "ssh-port", "", 22, "Host port number for bastion SSH host")
	updateDryRunCmd.Flags().StringVarP(&sshUser, "ssh-user", "", "", "Host username for bastion SSH host")
	updateDryRunCmd.Flags().StringVarP(&sshKeyPath, "ssh-key-path", "", home+"/.ssh/id_rsa", "Private key path for bastion SSH host")
	updateDryRunCmd.Flags().StringVarP(&sshPassphrase, "ssh-passphrase", "", "", "Private key passphrase for bastion SSH host")
}
