package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/yammerjp/db-time-traveler/system"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Try to connect to a relational database",
	Long:  `Try to connect to a relational database`,
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

		c, err := dap.CreateDatabaseConnection()
		if err != nil {
			panic(err)
		}
		defer c.Close()
		if err := c.Ping(); err != nil {
			log.Fatal("PingError: ", err)
		} else {
			log.Println("Ping Success!")
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringVarP(&username, "user", "u", "root", "Username for Database Connection")
	pingCmd.Flags().StringVarP(&password, "password", "p", "password", "Password for Database Connection")
	pingCmd.Flags().StringVarP(&host, "host", "", "localhost", "Hostname or IPv4 address for Database Connection")
	pingCmd.Flags().IntVarP(&port, "port", "", 3306, "Port number for Database Connection")
	pingCmd.Flags().StringVarP(&schema, "schema", "s", "", "Schema name for Database Connection")
}
