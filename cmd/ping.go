package cmd

import (
	"fmt"
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
		c, err := connection()
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

func connection() (*system.DatabaseConnection, error) {
	if sshHost == "" {
		dap := &system.DatabaseAccessPoint{
			Username: username,
			Password: password,
			Host:     host,
			Port:     port,
			Schema:   schema,
		}

		return dap.CreateDatabaseConnection()
	} else {
		dapOnSsh := &system.DatabaseAccessPointOnSSH{
			DB: &system.DB{
				Host:     host,
				Port:     fmt.Sprintf("%d", port),
				User:     username,
				Password: password,
				DBName:   schema,
			},
			SSH: &system.SSH{
				Key:        sshKeyPath,
				Host:       sshHost,
				User:       sshUser,
				Port:       fmt.Sprintf("%d", sshPort),
				Passphrase: sshPassphrase,
			},
		}
		return dapOnSsh.CreateDatabaseConnection()
	}

}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringVarP(&username, "user", "u", "root", "Username for Database Connection")
	pingCmd.Flags().StringVarP(&password, "password", "p", "password", "Password for Database Connection")
	pingCmd.Flags().StringVarP(&host, "host", "", "localhost", "Hostname or IPv4 address for Database Connection")
	pingCmd.Flags().IntVarP(&port, "port", "", 3306, "Port number for Database Connection")
	pingCmd.Flags().StringVarP(&schema, "schema", "s", "", "Schema name for Database Connection")
	pingCmd.Flags().StringVarP(&sshHost, "ssh-host", "", "", "Host name for bastion SSH host")
	pingCmd.Flags().IntVarP(&sshPort, "ssh-port", "", 22, "Host port number for bastion SSH host")
	pingCmd.Flags().StringVarP(&sshUser, "ssh-user", "", "", "Host username for bastion SSH host")
	pingCmd.Flags().StringVarP(&sshKeyPath, "ssh-key-path", "", "~/.ssh/id_rsa", "Private key path for bastion SSH host")
	pingCmd.Flags().StringVarP(&sshPassphrase, "ssh-passphrase", "", "", "Private key passphrase for bastion SSH host")
}
