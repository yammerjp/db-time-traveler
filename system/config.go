package system

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type ConnectionConfig struct {
	Name          string
	Driver        string
	Hostname      string
	Port          string
	Username      string
	Password      string
	Database      string
	SSHKeyPath    string
	SSHHost       string
	SSHPort       string
	SSHUser       string
	SSHPassphrase string
}

type Config struct {
	DefaultConnection string `yaml:"default_connection"`
	Connections       []ConnectionConfig
}

func findConfig(specifiedPath string) (string, bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", false, err
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", false, err
	}
	searchPaths := []string{
		specifiedPath,
		homeDir + "/.db-time-traveler.yaml",
		homeDir + "/.db-time-traveler.yml",
		configDir + "/db-time-traveler.yaml",
		configDir + "/db-time-traveler.yml",
	}
	for _, v := range searchPaths {
		if _, err := os.Stat(v); err == nil {
			return v, true, nil
		}
	}
	return "", false, nil
}

func loadConfigFromYaml(path string) (*Config, error) {
	var config Config
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(buf, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) FindConnection(targetConnectionName string) (*ConnectionConfig, error) {
	for _, v := range c.Connections {
		if targetConnectionName != "" && targetConnectionName == v.Name {
			return &v, nil
		} else if targetConnectionName == "" && c.DefaultConnection == v.Name {
			return &v, nil
		}
	}
	return nil, errors.New("Default Connection is not found")
}

func (connection *ConnectionConfig) CreateDatabaseConnection() (*DatabaseConnection, error) {
	var dap DatabaseAccessPointHub
	var err error
	if connection.SSHHost != "" {
		dap, err = connection.toDapOnSSH()
	} else {
		dap, err = connection.toDapWithoutSSH()
	}
	if err != nil {
		return nil, err
	}
	return dap.CreateDatabaseConnection()
}

func (connection *ConnectionConfig) toDapWithoutSSH() (*DatabaseAccessPoint, error) {
	port, err := strconv.Atoi(connection.Port)
	if err != nil {
		return nil, err
	}
	return &DatabaseAccessPoint{
		Username: connection.Username,
		Password: connection.Password,
		Host:     connection.Hostname,
		Port:     port,
		Schema:   connection.Database,
	}, nil
}

func (connection *ConnectionConfig) toDapOnSSH() (*DatabaseAccessPointOnSSH, error) {
	return &DatabaseAccessPointOnSSH{
		DB: &DB{
			Host:     connection.Hostname,
			Port:     connection.Port,
			User:     connection.Username,
			Password: connection.Password,
			DBName:   connection.Database,
		},
		SSH: &SSH{
			Key:        connection.SSHKeyPath,
			Host:       connection.SSHHost,
			User:       connection.SSHUser,
			Port:       connection.SSHPort,
			Passphrase: connection.SSHPassphrase,
		},
	}, nil
}

func LoadConfig(specifiedPath string) (*Config, error) {
	path, found, err := findConfig(specifiedPath)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("A config file is not found")
	}
	return loadConfigFromYaml(path)
}

func (connection *ConnectionConfig) ToString() string {
	ret := fmt.Sprintf("  Host: %s\n  Port: %s\n  Database: %s\n  Username: %s\n  Password: %s", connection.Hostname, connection.Port, connection.Database, connection.Username, connection.Password)
	if connection.SSHHost != "" {
		ret += fmt.Sprintf("\n  SSHHost: %s\n  SSHPort: %s\n  SSHUser: %s\n  SSHKeyPath: %s\n  SSHPathphrase: %s", connection.SSHHost, connection.SSHPort, connection.SSHUser, connection.SSHKeyPath, connection.SSHPassphrase)
	}
	return ret
}

func (c *Config) ToString() string {
	ret := ""
	for i, v := range c.Connections {
		if i != 0 {
			ret += "\n"
		}
		ret += v.Name
		if c.DefaultConnection == v.Name {
			ret += " (default)"
		}
		ret += "\n" + v.ToString()
	}
	return ret
}

func (connection *ConnectionConfig) Override(prioritizedConnection *ConnectionConfig) (*ConnectionConfig, error) {
	var ret ConnectionConfig
	if prioritizedConnection.Name != "" {
		ret.Name = prioritizedConnection.Name
	} else {
		ret.Name = connection.Name
	}
	if prioritizedConnection.Driver != "" {
		ret.Driver = prioritizedConnection.Driver
	} else {
		ret.Driver = connection.Driver
	}
	if prioritizedConnection.Hostname != "" {
		ret.Hostname = prioritizedConnection.Hostname
	} else {
		ret.Hostname = connection.Hostname
	}
	if prioritizedConnection.Port != "3306" && false {
		ret.Port = prioritizedConnection.Port
	} else {
		ret.Port = connection.Port
	}
	if prioritizedConnection.Username != "" {
		ret.Username = prioritizedConnection.Username
	} else {
		ret.Username = connection.Username
	}
	if prioritizedConnection.Password != "" {
		ret.Password = prioritizedConnection.Password
	} else {
		ret.Password = connection.Password
	}
	if prioritizedConnection.Database != "" {
		ret.Database = prioritizedConnection.Database
	} else {
		ret.Database = connection.Database
	}
	if prioritizedConnection.SSHKeyPath != "" {
		ret.SSHKeyPath = prioritizedConnection.SSHKeyPath
	} else {
		ret.SSHKeyPath = connection.SSHKeyPath
	}
	if prioritizedConnection.SSHHost != "" {
		ret.SSHHost = prioritizedConnection.SSHHost
	} else {
		ret.SSHHost = connection.SSHHost
	}
	if prioritizedConnection.SSHPort != "22" && false {
		ret.SSHPort = prioritizedConnection.SSHPort
	} else {
		ret.SSHPort = connection.SSHPort
	}
	if prioritizedConnection.SSHUser != "" {
		ret.SSHUser = prioritizedConnection.SSHUser
	} else {
		ret.SSHUser = connection.SSHUser
	}
	if prioritizedConnection.SSHPassphrase != "" {
		ret.SSHPassphrase = prioritizedConnection.SSHPassphrase
	} else {
		ret.SSHPassphrase = connection.SSHPassphrase
	}
	return &ret, nil
}

func CreateConfigFile(specifiedPath string) error {
	content := []byte(
		`default_connection: local
connections:
  -
    name: local
    driver: mysql
    hostname: localhost
    username: root
    password: password
    port: 3306
    database: sampleschema
  -
    name: sshconnection
    driver: mysql
    hostname: localhost
    username: root
    password: password
    port: 3306
    database: sampleschema
    sshhost: bastion.example.com
    sshport: 22
    sshuser: yammer
    sshkeypath: /home/username/.ssh/id_rsa
    sshpassphrase: helloworld
`)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	searchPaths := []string{
		specifiedPath,
		homeDir + "/.db-time-traveler.yaml",
		homeDir + "/.db-time-traveler.yml",
		configDir + "/db-time-traveler.yaml",
		configDir + "/db-time-traveler.yml",
	}
	targetFile := ""
	for _, v := range searchPaths {
		if v == "" {
			continue
		}
		if _, err := os.Stat(v); err != nil {
			targetFile = v
			break
		}
		fmt.Println(v)
	}
	if targetFile == "" {
		return errors.New("all search path are already exist config files")
	}
	err = os.WriteFile(targetFile, content, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("created %s", targetFile)
	return nil
}
