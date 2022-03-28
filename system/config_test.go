package system

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestLoadConfigYaml(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}
	expected := Config{
		DefaultConnection: "local",
		Connections: []ConnectionConfig{
			{
				Name:     "local",
				Driver:   "mysql",
				Hostname: "localhost",
				Port:     "3306",
				Username: "root",
				Password: "password",
				Database: "sampleschema",
			},
		},
	}
	ret, err := loadConfigYaml(home + "/go/src/github.com/yammerjp/db-time-traveler/.db-time-traveler.yml")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*ret, expected) {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("loadConnectionConfig must be return a expected statement")
	}
}
