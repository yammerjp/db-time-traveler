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

func (q QueryBuilderSourceForColumnValues) buildWhereIn() string {
	return " WHERE (" + strings.Join(q.primaryKeys, ", ") + ") IN ( " + q.stmtInWhereIn + " )"
}

func (q QueryBuilderSourceForSchemaInformation) buildFrom() string {
	return " FROM " + q.targetTable
}

func (q QueryBuilderSourceToUpdate) buildStmtToUpdate() (string, error) {
	query := "UPDATE " + q.targetTable + " SET"

	intervalStr := q.buildInterval()
	if len(q.columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	for i, column := range q.columns {
		if i == 0 {
			query += " " + column + " = (" + column + intervalStr + ")"
		} else {
			query += ", " + column + " = (" + column + intervalStr + ")"
		}
	}
	return query + q.buildWhereIn(), nil
}

func (q QueryBuilderSourceForColumnValues) buildStmtToSelect() (string, error) {
	query := "SELECT"

	if len(q.columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	for i, column := range q.columns {
		if i == 0 {
			query += " " + column
		} else {
			query += ", " + column
		}
	}
	return query + q.buildFrom() + q.buildWhereIn(), nil
}

func (q QueryBuilderSourceToUpdate) buildStmtToSelect() (string, error) {
	query := "SELECT"

	intervalStr := q.buildInterval()
	if len(q.columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	for i, column := range q.columns {
		if i != 0 {
			query += ","
		}
		query += " " + column + intervalStr
	}
	return query + q.buildFrom() + q.buildWhereIn(), nil
}

func (q QueryBuilderSourceForSchemaInformation) buildStmtToSelectColumnNames() string {
	return `SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = "` + q.targetTable + `"`
}

func (q QueryBuilderSourceForSchemaInformation) buildStmtToSelectColumnNamesDateRelated() (string, error) {
	return q.buildStmtToSelectColumnNames() + " AND DATA_TYPE IN (\"date\", \"datetime\", \"timestamp\")", nil // + DATA_TYPE = time
}

func (q QueryBuilderSourceForSchemaInformation) buildStmtToSelectColumnNamesOfPrimaryKey() (string, error) {
	return q.buildStmtToSelectColumnNames() + " AND COLUMN_KEY = \"PRI\"", nil
}

func (q QueryBuilderSourceToUpdate) buildStmtToSelectBeforeAndAfter() (string, error) {
	query := "SELECT " + strings.Join(q.primaryKeys, ", ")
	intervalStr := q.buildInterval()
	for _, column := range q.columns {
		query += ", " + column + ", " + column + intervalStr
	}
	return query + q.buildFrom() + q.buildWhereIn(), nil
}
