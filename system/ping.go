package system

import (
	"log"
)

func Ping(dap *DatabaseAccessPoint) {
	c, err := dap.createDatabaseConnection()
	if err != nil {
		panic(err)
	}
	defer c.Close()
	if err := c.ping(); err != nil {
		log.Fatal("PingError: ", err)
	} else {
		log.Println("Ping Success!")
	}
}
