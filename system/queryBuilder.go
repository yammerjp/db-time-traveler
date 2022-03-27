package system

import (
	"errors"
)

// tableName string, columnNames []TimeRelatedColumn, wheres []WhereClause
func updateQueryBuilder(targetTable string, columns []string, interval string, primaryKey string, whereInStmt string) (string, error) {
	query := "UPDATE " + targetTable + " SET"

	if len(columns) == 0 {
		return "", errors.New("Must be specify any columns")
	}
	for i, column := range columns {
		if i == 0 {
			query += " " + column + " = (" + column + " - INTERVAL " + interval + ")"
		} else {
			query += ", " + column + " = (" + column + " - INTERVAL " + interval + ")"
		}
	}

	query += " WHERE " + primaryKey + " IN ( " + whereInStmt + " )"
	return query, nil
}

func selectTargettedColumnsQueryBuilder(targetTable string, columns []string, primaryKey string, whereInStmt string) (string, error) {
	query := "SELECT"

	if len(columns) == 0 {
		return "", errors.New("Must be specify any columns")
	}
	for i, column := range columns {
		if i == 0 {
			query += " " + column
		} else {
			query += ", " + column
		}
	}
	query += " FROM " + targetTable
	query += " WHERE " + primaryKey + " IN ( " + whereInStmt + " )"
	return query, nil
}

func selectUpdatingColumnValiesQueryBuilder(targetTable string, columns []string, interval string, primaryKey string, whereInStmt string) (string, error) {
	query := "SELECT"

	if len(columns) == 0 {
		return "", errors.New("Must be specify any columns")
	}
	for i, column := range columns {
		if i == 0 {
			query += " " + column + " - INTERVAL " + interval
		} else {
			query += ", " + column + " - INTERVAL " + interval
		}
	}
	query += " FROM " + targetTable
	query += " WHERE " + primaryKey + " IN ( " + whereInStmt + " )"
	return query, nil
}

func selectDateRelatedColumnsQueryBuilder(targetTable string) (string, error) {
	query := "SELECT COLUMN_NAME"
	query += " FROM INFORMATION_SCHEMA.COLUMNS"
	query += " WHERE table_name = \"" + targetTable + "\""
	query += " AND DATA_TYPE IN (\"date\", \"datetime\", \"timestamp\")" // + DATA_TYPE = time
	return query, nil
}

func selectPrimaryKeyColumnsQueryBuilder(targetTable string) (string, error) {
	query := "SELECT COLUMN_NAME"
	query += " FROM INFORMATION_SCHEMA.COLUMNS"
	query += " WHERE table_name = \"" + targetTable + "\""
	query += " AND COLUMN_KEY = \"PRI\""
	return query, nil
}
