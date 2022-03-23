package system

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseAccessPoint struct {
	Username string
	Password string
	Host     string
	Port     int
	Schema   string
}

func (p DatabaseAccessPoint) toString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", p.Username, p.Password, p.Host, p.Port, p.Schema)
}

func (p DatabaseAccessPoint) connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", p.toString())
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}
