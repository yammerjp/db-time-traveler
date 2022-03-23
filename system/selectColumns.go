package system

import (
	"fmt"
	"log"
	"strings"
)

func SelectAndPrintColumns(dap *DatabaseAccessPoint, tableName string) {
	trcs, err := dap.SelectColumns(tableName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", trcs.toString())
}

type TimeRelatedColumns struct {
	Date        []string
	Time        []string
	DateTime    []string
	Timestamp   []string
	PrimaryKeys []string
}

func (trcs *TimeRelatedColumns) toString() string {
	lines := make([]string, 0)
	for _, column := range trcs.Date {
		lines = append(lines, fmt.Sprintf("Date        %40s", column))
	}
	for _, column := range trcs.Time {
		lines = append(lines, fmt.Sprintf("Time        %40s", column))
	}
	for _, column := range trcs.DateTime {
		lines = append(lines, fmt.Sprintf("DateTime    %40s", column))
	}
	for _, column := range trcs.Timestamp {
		lines = append(lines, fmt.Sprintf("Timestamp   %40s", column))
	}
	for _, column := range trcs.PrimaryKeys {
		lines = append(lines, fmt.Sprintf("PrimaryKeys %40s", column))
	}

	return strings.Join(lines, "\n")
}

func (dap *DatabaseAccessPoint) SelectColumns(tableName string) (*TimeRelatedColumns, error) {
	db, err := dap.connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT COLUMN_NAME, DATA_TYPE, COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = ?", tableName)
	if err != nil {
		return nil, err
	}

	timeRelatedColumns := TimeRelatedColumns{}

	for rows.Next() {
		var column string
		var dataType string
		var columnKey string
		if err := rows.Scan(&column, &dataType, &columnKey); err != nil {
			log.Fatal(err)
		}
		switch dataType {
		case "date":
			timeRelatedColumns.Date = append(timeRelatedColumns.Date, column)
		case "time":
			timeRelatedColumns.Time = append(timeRelatedColumns.Time, column)
		case "datetime":
			timeRelatedColumns.DateTime = append(timeRelatedColumns.DateTime, column)
		case "timestamp":
			timeRelatedColumns.Timestamp = append(timeRelatedColumns.Timestamp, column)
		default:
		}
		if columnKey == "PRI" {
			timeRelatedColumns.PrimaryKeys = append(timeRelatedColumns.PrimaryKeys, column)
		}
	}
	return &timeRelatedColumns, nil
}
