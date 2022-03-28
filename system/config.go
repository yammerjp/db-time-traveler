package system

import (
	"errors"
	"fmt"
	"io/ioutil"
	"ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type ConnectionConfig struct {
	Database DB
}

type RootConfig struct {
	defaultConnection ConnectionConfig
}

func loadConfigYaml(path string) (*ConnectionConfig, error) {
	rootConfig := RootConfig{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	rootConfig, err = yaml.Unmarshal(file)
	fmt.Printf("%s", rootConfig)
	return nil, errors.New("not implemented")
}
