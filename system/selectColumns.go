package system

import (
	"fmt"
	"log"
)

func SelectColumns(dap *DatabaseAccessPoint) {
	db, err := dap.connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = 'accounts'")
	if err != nil {
		log.Fatal(err)
	}

	columns := make([]string, 0)
	dataTypes := make([]string, 0)

	for rows.Next() {
		var column string
		var dataType string
		if err := rows.Scan(&column, &dataType); err != nil {
			log.Fatal(err)
		}
		columns = append(columns, column)
		dataTypes = append(dataTypes, dataType)
	}

	for i, column := range columns {
		fmt.Printf("%40s %40s\n", column, dataTypes[i])
	}
}
