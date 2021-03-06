package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/yammerjp/db-time-traveler/system"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Try to connect to a relational database",
	Long:  `Try to connect to a relational database`,
	Run: func(cmd *cobra.Command, args []string) {
		ping()
	},
}

func ping() {
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
}

func connection() (*system.DatabaseConnection, error) {
	return connect()
}

func init() {
	rootCmd.AddCommand(pingCmd)
	initConnection(pingCmd.Flags())
}

func initConnection(f *flag.FlagSet) {
	f.StringVarP(&username, "user", "u", "", "Username for Database Connection")
	f.StringVarP(&password, "password", "p", "", "Password for Database Connection")
	f.StringVarP(&host, "host", "", "", "Hostname or IPv4 address for Database Connection")
	f.IntVarP(&port, "port", "", 3306, "Port number for Database Connection")
	f.StringVarP(&schema, "schema", "s", "", "Schema name for Database Connection")
	f.StringVarP(&sshHost, "ssh-host", "", "", "Host name for bastion SSH host")
	f.IntVarP(&sshPort, "ssh-port", "", 22, "Host port number for bastion SSH host")
	f.StringVarP(&sshUser, "ssh-user", "", "", "Host username for bastion SSH host")
	f.StringVarP(&sshKeyPath, "ssh-key-path", "", "", "Private key path for bastion SSH host")
	f.StringVarP(&sshPassphrase, "ssh-passphrase", "", "", "Private key passphrase for bastion SSH host")
	initShowConfigFlg(f)
	f.StringVarP(&selectedConnection, "connection", "", "", "connection name")
}

func loadConnectionConfigFromCommandlineArguments() *system.ConnectionConfig {
	return &system.ConnectionConfig{
		Hostname:      host,
		Port:          fmt.Sprintf("%d", port),
		Username:      username,
		Password:      password,
		Database:      schema,
		SSHKeyPath:    sshKeyPath,
		SSHHost:       sshHost,
		SSHUser:       sshUser,
		SSHPort:       fmt.Sprintf("%d", sshPort),
		SSHPassphrase: sshPassphrase,
	}
}

func loadConnectionConfig() (*system.ConnectionConfig, error) {
	config, err := system.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	connection, err := config.FindConnection(selectedConnection)
	if err != nil {
		return nil, err
	}
	fromCmdArgs := loadConnectionConfigFromCommandlineArguments()

	overridePort := false
	overrideSSHPort := false
	for _, v := range os.Args {
		if v == "--port" {
			overridePort = true
		}
		if v == "--ssh-port" {
			overrideSSHPort = true
		}
	}
	return connection.Override(fromCmdArgs, overridePort, overrideSSHPort)
}

func connect() (*system.DatabaseConnection, error) {
	connection, err := loadConnectionConfig()
	if err != nil {
		return nil, err
	}
	return connection.CreateDatabaseConnection()
}
