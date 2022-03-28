package system

import (
	"os"

	"gopkg.in/yaml.v2"
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

func loadConfigYaml(path string) (*Config, error) {
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
