package system

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLoadConfigYaml(t *testing.T) {
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
	ret, err := loadConfigFromYaml("../.db-time-traveler.yml")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(*ret, expected) {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("loadConnectionConfig must be return a expected statement")
	}
}
