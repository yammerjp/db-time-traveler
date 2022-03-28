package system

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
)

type ConnectionConfig struct {
	Name     string
	Driver   string
	Hostname string
	Port     string
	Username string
	Password string
	Database string
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

func (c *Config) findConnection(targetConnectionName string) (*ConnectionConfig, error) {
	for _, v := range c.Connections {
		if targetConnectionName != "" && targetConnectionName == v.Name {
			return &v, nil
		} else if targetConnectionName == "" && c.DefaultConnection == v.Name {
			return &v, nil
		}
	}
	return nil, errors.New("Default Connection is not found")
}

func (c *Config) ToDap(selectedConnection string) (*DatabaseAccessPoint, error) {
	connection, err := c.findConnection(selectedConnection)
	if err != nil {
		return nil, err
	}
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

func LoadConfig(specifiedPath string, selectedConnection string) (*DatabaseAccessPoint, error) {
	path, found, err := findConfig(specifiedPath)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("A config file is not found")
	}

	config, err := loadConfigFromYaml(path)
	if err != nil {
		return nil, err
	}
	return config.ToDap(selectedConnection)
}
