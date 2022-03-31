package system

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
)

var dc DatabaseConnection

func initializeDatabase() {
	dc.connection.Query(`CREATE TABLE accounts
(
  id                            INT(10) auto_increment,
  name                          VARCHAR(40),
  trial_end_date                DATE,
  registered_campaign_end_datetime  DATETIME,
  created_at                    TIMESTAMP default current_timestamp,
  updated_at                    TIMESTAMP default current_timestamp on update current_timestamp,
  PRIMARY KEY (id)
)`)
	dc.connection.Query(`
INSERT INTO accounts (id, name, trial_end_date, registered_campaign_end_datetime) VALUES (1, "Nagaoka", "2022-05-30", "2022-04-03 00:10:20", "2022-04-03 00:10:20")
`)
	dc.connection.Query(`
INSERT INTO accounts (id, name) VALUES (2, "Tanaka")
`)
	dc.connection.Query(`CREATE TABLE multiple_primary_keys
(
  id                            INT(10) auto_increment,
  other_primary_column          VARCHAR(40),
  expired_datetime              DATETIME,
  created_at                    TIMESTAMP default current_timestamp,
  updated_at                    TIMESTAMP default current_timestamp on update current_timestamp,
  PRIMARY KEY (id, other_primary_column)
)`)
	dc.connection.Query(`
INSERT INTO accounts (id, other_primary_column, expired_date, created_at, updated_at) VALUES (1, "Nagaoka", "2022-05-30", "2022-04-03 00:10:20", "2022-04-03 00:10:20")
`)
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		db, err := sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		dc = DatabaseConnection{
			connection:    db,
			sshConnection: nil,
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	initializeDatabase()

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestSelectDateRelatedColumns(t *testing.T) {
	expected := []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"}
	got, err := dc.SelectDateRelatedColumns("accounts", []string{})
	if err != nil {
		t.Logf("%s", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Logf("expected %s, got %s", expected, got)
	}

	// ...
}
