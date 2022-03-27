package system

import (
	"errors"
	"fmt"
	"log"
)

type TimeRelatedColumnType int

const (
	Date TimeRelatedColumnType = iota
	Time
	DateTime
	Timestamp
)

type TimeRelatedColumn struct {
	ColumnType TimeRelatedColumnType
	Name       string
	IsPrimary  bool
}

func (trc *TimeRelatedColumn) toString() string {
	var columnTypeString string
	switch trc.ColumnType {
	case Date:
		columnTypeString = "Date"
	case Time:
		columnTypeString = "Time"
	case DateTime:
		columnTypeString = "DateTime"
	case Timestamp:
		columnTypeString = "Timestamp"
	default:
		log.Fatal(errors.New(fmt.Sprintf("Unknown ColumnType: %d", trc.ColumnType)))
	}
	if trc.IsPrimary {
		return columnTypeString + "(primary): " + trc.Name
	} else {
		return columnTypeString + ": " + trc.Name
	}
}

func (trc *TimeRelatedColumn) ToString() string {
	return trc.toString()
}
