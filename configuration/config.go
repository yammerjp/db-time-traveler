package configuration

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/yammerjp/db-time-traveler/database"
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

func (conn *ConnectionConfig) CreateDatabaseConnection() (*database.DatabaseConnection, error) {
	var dap database.DatabaseAccessPointHub
	var err error
	if conn.SSHHost != "" {
		dap, err = conn.toDapOnSSH()
	} else {
		dap, err = conn.toDapWithoutSSH()
	}
	if err != nil {
		return nil, err
	}
	return dap.CreateDatabaseConnection()
}

func (conn *ConnectionConfig) toDapWithoutSSH() (*database.DatabaseAccessPoint, error) {
	port, err := strconv.Atoi(conn.Port)
	if err != nil {
		return nil, err
	}
	return &database.DatabaseAccessPoint{
		Username: conn.Username,
		Password: conn.Password,
		Host:     conn.Hostname,
		Port:     port,
		Schema:   conn.Database,
	}, nil
}

func (conn *ConnectionConfig) toDapOnSSH() (*database.DatabaseAccessPointOnSSH, error) {
	return &database.DatabaseAccessPointOnSSH{
		DB: &database.DB{
			Host:     conn.Hostname,
			Port:     conn.Port,
			User:     conn.Username,
			Password: conn.Password,
			DBName:   conn.Database,
		},
		SSH: &database.SSH{
			Key:        conn.SSHKeyPath,
			Host:       conn.SSHHost,
			User:       conn.SSHUser,
			Port:       conn.SSHPort,
			Passphrase: conn.SSHPassphrase,
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

func (conn *ConnectionConfig) ToString() string {
	ret := fmt.Sprintf("  Host: %s\n  Port: %s\n  Database: %s\n  Username: %s\n  Password: %s", conn.Hostname, conn.Port, conn.Database, conn.Username, conn.Password)
	if conn.SSHHost != "" {
		ret += fmt.Sprintf("\n  SSHHost: %s\n  SSHPort: %s\n  SSHUser: %s\n  SSHKeyPath: %s\n  SSHPathphrase: %s", conn.SSHHost, conn.SSHPort, conn.SSHUser, conn.SSHKeyPath, conn.SSHPassphrase)
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

func (conn *ConnectionConfig) Override(prioritizedConnection *ConnectionConfig, overridePort bool, overrideSSHPort bool) (*ConnectionConfig, error) {
	var ret ConnectionConfig
	if prioritizedConnection.Name != "" {
		ret.Name = prioritizedConnection.Name
	} else {
		ret.Name = conn.Name
	}
	if prioritizedConnection.Driver != "" {
		ret.Driver = prioritizedConnection.Driver
	} else {
		ret.Driver = conn.Driver
	}
	if prioritizedConnection.Hostname != "" {
		ret.Hostname = prioritizedConnection.Hostname
	} else {
		ret.Hostname = conn.Hostname
	}
	if overridePort {
		ret.Port = prioritizedConnection.Port
	} else {
		ret.Port = conn.Port
	}
	if prioritizedConnection.Username != "" {
		ret.Username = prioritizedConnection.Username
	} else {
		ret.Username = conn.Username
	}
	if prioritizedConnection.Password != "" {
		ret.Password = prioritizedConnection.Password
	} else {
		ret.Password = conn.Password
	}
	if prioritizedConnection.Database != "" {
		ret.Database = prioritizedConnection.Database
	} else {
		ret.Database = conn.Database
	}
	if prioritizedConnection.SSHKeyPath != "" {
		ret.SSHKeyPath = prioritizedConnection.SSHKeyPath
	} else {
		ret.SSHKeyPath = conn.SSHKeyPath
	}
	if prioritizedConnection.SSHHost != "" {
		ret.SSHHost = prioritizedConnection.SSHHost
	} else {
		ret.SSHHost = conn.SSHHost
	}
	if overrideSSHPort {
		ret.SSHPort = prioritizedConnection.SSHPort
	} else {
		ret.SSHPort = conn.SSHPort
	}
	if prioritizedConnection.SSHUser != "" {
		ret.SSHUser = prioritizedConnection.SSHUser
	} else {
		ret.SSHUser = conn.SSHUser
	}
	if prioritizedConnection.SSHPassphrase != "" {
		ret.SSHPassphrase = prioritizedConnection.SSHPassphrase
	} else {
		ret.SSHPassphrase = conn.SSHPassphrase
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
