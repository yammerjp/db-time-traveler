package system

import (
	"errors"
	"fmt"
	"strings"
)

type QueryBuilderSourcePartOfInterval struct {
	isPast bool
	num    int
	term   string
}

type QueryBuilderSourceForSchemaInformation struct {
	targetTable string
}

type QueryBuilderSourceForColumnValues struct {
	QueryBuilderSourceForSchemaInformation
	columns     []string
	primaryKeys []string
	whereInStmt string
}

type QueryBuilderSourceToUpdate struct {
	QueryBuilderSourceForColumnValues
  QueryBuilderSourcePartOfInterval
}

func (q *QueryBuilderSourcePartOfInterval) buildInterval() (string, error) {
  var symbol string
  if q.isPast {
    symbol = "-"
  } else {
    symbol = "+"
  }

  return fmt.Sprintf(" %s INTERVAL %d %s", symbol, q.num, q.term)  , nil
}

func (q *QueryBuilderSourceToUpdate) buildToUpdate() (string, error) {
	if q == nil {
		return "", errors.New("QueryBuilderSource is nil on building query")
	}
  interval, err := q.buildInterval()
  if err != nil {
    return "", err
  }
	return updateQueryBuilder(q.targetTable, q.columns, interval, q.primaryKeys, q.whereInStmt)
}

func updateQueryBuilder(targetTable string, columns []string, interval string, primaryKeys []string, whereInStmt string) (string, error) {
	query := "UPDATE " + targetTable + " SET"

	if len(columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	for i, column := range columns {
		if i == 0 {
			query += " " + column + " = (" + column + interval + ")"
		} else {
			query += ", " + column + " = (" + column + interval + ")"
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
  interval, err := q.buildInterval()
  if err != nil {
    return "", err
  }
	return selectUpdatingColumnValuesQueryBuilder(q.targetTable, q.columns, interval, q.primaryKeys, q.whereInStmt)
}

func selectUpdatingColumnValuesQueryBuilder(targetTable string, columns []string, interval string, primaryKeys []string, whereInStmt string) (string, error) {
	query := "SELECT"

	if len(columns) == 0 {
		return "", errors.New("Must be specify any columns")
	}
	for i, column := range columns {
		if i == 0 {
			query += " " + column + interval
		} else {
			query += ", " + column + interval
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
	query := "SELECT DISTINCT COLUMN_NAME"
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
	query := "SELECT DISTINCT COLUMN_NAME"
	query += " FROM INFORMATION_SCHEMA.COLUMNS"
	query += " WHERE table_name = \"" + targetTable + "\""
	query += " AND COLUMN_KEY = \"PRI\""
	return query, nil
}

func (q *QueryBuilderSourceToUpdate) buildToSelectBeforeAndAfter() (string, error) {
  interval, err := q.buildInterval()
  if err != nil {
    return "", err
  }
	return selectUpdatingColumnValuesBeforeAndAfterQueryBuilder(q.targetTable, q.columns, interval, q.primaryKeys, q.whereInStmt)
}

func selectUpdatingColumnValuesBeforeAndAfterQueryBuilder(targetTable string, columns []string, interval string, primaryKeys []string, whereInStmt string) (string, error) {
	query := "SELECT " + strings.Join(primaryKeys, ", ")
	for _, column := range columns {
		query += ", " + column + ", " + column + interval
	}
	query += " FROM " + targetTable
	query += " WHERE (" + strings.Join(primaryKeys, ", ") + ") IN ( " + whereInStmt + " )"
	return query, nil
}
