package system

import (
	"log"
)

func Ping(dap *DatabaseAccessPoint) {
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
