package system

import (
	"fmt"
	"os"
	"testing"
)

func TestLoadConfigYaml(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}
	expected := ConnectionConfig{
		Database: DB{
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "password",
			DBName:   "sampleschema",
		},
	}
	ret, err := loadConfigYaml(home + "/go/src/github.com/yammerjp/db-time-traveler/.db-time-traveler.yml")
	if err != nil {
		t.Error(err)
	}
	if *ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("loadConnectionConfig must be return a expected statement")
	}
}
