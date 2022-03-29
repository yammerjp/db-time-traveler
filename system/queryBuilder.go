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
	columns       []string
	primaryKeys   []string
	stmtInWhereIn string
}

type QueryBuilderSourceToUpdate struct {
	QueryBuilderSourceForColumnValues
	QueryBuilderSourcePartOfInterval
}

func (q *QueryBuilderSourcePartOfInterval) buildInterval() string {
	var symbol string
	if q.isPast {
		symbol = "-"
	} else {
		symbol = "+"
	}
	return fmt.Sprintf(" %s INTERVAL %d %s", symbol, q.num, q.term)
}

func (q *QueryBuilderSourceToUpdate) buildToUpdate() (string, error) {
	if q == nil {
		return "", errors.New("QueryBuilderSource is nil on building query")
	}
	return updateQueryBuilder(q.targetTable, q.columns, q.QueryBuilderSourcePartOfInterval, q.primaryKeys, q.stmtInWhereIn)
}

func updateQueryBuilder(targetTable string, columns []string, interval QueryBuilderSourcePartOfInterval, primaryKeys []string, stmtInWhereIn string) (string, error) {
	query := "UPDATE " + targetTable + " SET"

	intervalStr := interval.buildInterval()
	if len(columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	for i, column := range columns {
		if i == 0 {
			query += " " + column + " = (" + column + intervalStr + ")"
		} else {
			query += ", " + column + " = (" + column + intervalStr + ")"
		}
	}

	query += " WHERE (" + strings.Join(primaryKeys, ", ") + ") IN ( " + stmtInWhereIn + " )"
	return query, nil
}

func (q *QueryBuilderSourceForColumnValues) buildToSelect() (string, error) {
	if q == nil {
		return "", errors.New("QueryBuilderSource is nil on building query")
	}
	return selectTargettedColumnsQueryBuilder(q.targetTable, q.columns, q.primaryKeys, q.stmtInWhereIn)
}

func selectTargettedColumnsQueryBuilder(targetTable string, columns []string, primaryKeys []string, stmtInWhereIn string) (string, error) {
	query := "SELECT"

	if len(columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	for i, column := range columns {
		if i == 0 {
			query += " " + column
		} else {
			query += ", " + column
		}
	}
	query += " FROM " + targetTable
	query += " WHERE (" + strings.Join(primaryKeys, ", ") + ") IN ( " + stmtInWhereIn + " )"
	return query, nil
}

func (q *QueryBuilderSourceToUpdate) buildToSelect() (string, error) {
	if q == nil {
		return "", errors.New("QueryBuilderSource is nil on building query")
	}
	return selectUpdatingColumnValuesQueryBuilder(q.targetTable, q.columns, q.QueryBuilderSourcePartOfInterval, q.primaryKeys, q.stmtInWhereIn)
}

func selectUpdatingColumnValuesQueryBuilder(targetTable string, columns []string, interval QueryBuilderSourcePartOfInterval, primaryKeys []string, stmtInWhereIn string) (string, error) {
	query := "SELECT"

	intervalStr := interval.buildInterval()
	if len(columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	for i, column := range columns {
		if i == 0 {
			query += " " + column + intervalStr
		} else {
			query += ", " + column + intervalStr
		}
	}
	query += " FROM " + targetTable
	query += " WHERE (" + strings.Join(primaryKeys, ", ") + ") IN ( " + stmtInWhereIn + " )"
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
	return selectUpdatingColumnValuesBeforeAndAfterQueryBuilder(q.targetTable, q.columns, q.QueryBuilderSourcePartOfInterval, q.primaryKeys, q.stmtInWhereIn)
}

func selectUpdatingColumnValuesBeforeAndAfterQueryBuilder(targetTable string, columns []string, interval QueryBuilderSourcePartOfInterval, primaryKeys []string, stmtInWhereIn string) (string, error) {
	query := "SELECT " + strings.Join(primaryKeys, ", ")
	intervalStr := interval.buildInterval()
	for _, column := range columns {
		query += ", " + column + ", " + column + intervalStr
	}
	query += " FROM " + targetTable
	query += " WHERE (" + strings.Join(primaryKeys, ", ") + ") IN ( " + stmtInWhereIn + " )"
	return query, nil
}
