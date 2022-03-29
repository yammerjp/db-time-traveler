package system

import (
	"errors"
	"strings"
)

type QueryBuilderSourceForSchemaInformation struct {
	targetTable string
}

type QueryBuilderSourceForColumnValues struct {
	columns     []string
	primaryKeys []string
	whereInStmt string
	QueryBuilderSourceForSchemaInformation
}

type QueryBuilderSourceToUpdate struct {
	past string
	QueryBuilderSourceForColumnValues
}

func (q *QueryBuilderSourceToUpdate) buildToUpdate() (string, error) {
	if q == nil {
		return "", errors.New("QueryBuilderSource is nil on building query")
	}
	return updateQueryBuilder(q.targetTable, q.columns, q.past, q.primaryKeys, q.whereInStmt)
}

func updateQueryBuilder(targetTable string, columns []string, past string, primaryKeys []string, whereInStmt string) (string, error) {
	query := "UPDATE " + targetTable + " SET"

	if len(columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	for i, column := range columns {
		if i == 0 {
			query += " " + column + " = (" + column + " - INTERVAL " + past + ")"
		} else {
			query += ", " + column + " = (" + column + " - INTERVAL " + past + ")"
		}
	}

	query += " WHERE (" + strings.Join(primaryKeys, ", ") + ") IN ( " + whereInStmt + " )"
	return query, nil
}

func (q *QueryBuilderSourceForColumnValues) buildToSelect() (string, error) {
	if q == nil {
		return "", errors.New("QueryBuilderSource is nil on building query")
	}
	return selectTargettedColumnsQueryBuilder(q.targetTable, q.columns, q.primaryKeys, q.whereInStmt)
}

func selectTargettedColumnsQueryBuilder(targetTable string, columns []string, primaryKeys []string, whereInStmt string) (string, error) {
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
	query += " WHERE (" + strings.Join(primaryKeys, ", ") + ") IN ( " + whereInStmt + " )"
	return query, nil
}

func (q *QueryBuilderSourceToUpdate) buildToSelect() (string, error) {
	if q == nil {
		return "", errors.New("QueryBuilderSource is nil on building query")
	}
	return selectUpdatingColumnValuesQueryBuilder(q.targetTable, q.columns, q.past, q.primaryKeys, q.whereInStmt)
}

func selectUpdatingColumnValuesQueryBuilder(targetTable string, columns []string, past string, primaryKeys []string, whereInStmt string) (string, error) {
	query := "SELECT"

	if len(columns) == 0 {
		return "", errors.New("Must be specify any columns")
	}
	for i, column := range columns {
		if i == 0 {
			query += " " + column + " - INTERVAL " + past
		} else {
			query += ", " + column + " - INTERVAL " + past
		}
	}
	query += " FROM " + targetTable
	query += " WHERE (" + strings.Join(primaryKeys, ", ") + ") IN ( " + whereInStmt + " )"
	return query, nil
}

func (q *QueryBuilderSourceForSchemaInformation) buildToSelectDateRelatedColumns() (string, error) {
	if q == nil {
		return "", errors.New("QueryBuilderSource is nil on building query")
	}
	return selectDateRelatedColumnsQueryBuilder(q.targetTable)
}

func selectDateRelatedColumnsQueryBuilder(targetTable string) (string, error) {
	query := "SELECT COLUMN_NAME"
	query += " FROM INFORMATION_SCHEMA.COLUMNS"
	query += " WHERE table_name = \"" + targetTable + "\""
	query += " AND DATA_TYPE IN (\"date\", \"datetime\", \"timestamp\")" // + DATA_TYPE = time
	return query, nil
}

func (q *QueryBuilderSourceForSchemaInformation) buildToSelectPrimaryColumns() (string, error) {
	if q == nil {
		return "", errors.New("QueryBuilderSource is nil on building query")
	}
	return selectPrimaryKeyColumnsQueryBuilder(q.targetTable)
}

func selectPrimaryKeyColumnsQueryBuilder(targetTable string) (string, error) {
	query := "SELECT COLUMN_NAME"
	query += " FROM INFORMATION_SCHEMA.COLUMNS"
	query += " WHERE table_name = \"" + targetTable + "\""
	query += " AND COLUMN_KEY = \"PRI\""
	return query, nil
}

func (q *QueryBuilderSourceToUpdate) buildToSelectBeforeAndAfter() (string, error) {
	return selectUpdatingColumnValuesBeforeAndAfterQueryBuilder(q.targetTable, q.columns, q.past, q.primaryKeys, q.whereInStmt)
}

func selectUpdatingColumnValuesBeforeAndAfterQueryBuilder(targetTable string, columns []string, past string, primaryKeys []string, whereInStmt string) (string, error) {
	query := "SELECT " + strings.Join(primaryKeys, ", ")
	for _, column := range columns {
		query += ", " + column + ", " + column + " - INTERVAL " + past
	}
	query += " FROM " + targetTable
	query += " WHERE (" + strings.Join(primaryKeys, ", ") + ") IN ( " + whereInStmt + " )"
	return query, nil
}
