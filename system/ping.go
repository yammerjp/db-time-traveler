package system

import (
	"log"
)

func Ping() {
	dap := &DatabaseAccessPoint{
		username: "root",
		password: "password",
		host:     "127.0.0.1",
		port:     3306,
		schema:   "sampleschema",
	}
	db, err := dap.connect()

	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("PingError: ", err)
	} else {
		log.Println("Ping Success!")
	}
}
