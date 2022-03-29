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

func (q QueryBuilderSourceForSchemaInformation) buildFrom() string {
	return " FROM " + q.targetTable
}

func (q QueryBuilderSourceForColumnValues) buildWhereIn() string {
	return " WHERE (" + strings.Join(q.primaryKeys, ", ") + ") IN ( " + q.stmtInWhereIn + " )"
}

func (q QueryBuilderSourceToUpdate) buildStmtToUpdate() (string, error) {
	if len(q.columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	query := "UPDATE " + q.targetTable + " SET "
	for i, column := range q.columns {
		if i != 0 {
			query += ", "
		}
		query += column + " = (" + column + q.buildInterval() + ")"
	}
	return query + q.buildWhereIn(), nil
}

func (q QueryBuilderSourceForColumnValues) buildStmtToSelect() (string, error) {
	if len(q.columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	return "SELECT " + strings.Join(q.columns, ", ") + q.buildFrom() + q.buildWhereIn(), nil
}

func (q QueryBuilderSourceToUpdate) buildStmtToSelect() (string, error) {
	if len(q.columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	query := "SELECT "
	for i, v := range q.columns {
		if i != 0 {
			query += ", "
		}
		query += v + q.buildInterval()
	}
	return query + q.buildFrom() + q.buildWhereIn(), nil
}

func (q QueryBuilderSourceToUpdate) buildStmtToSelectBeforeAndAfter() (string, error) {
	if len(q.columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	query := "SELECT " + strings.Join(q.primaryKeys, ", ") + ", "
	for i, v := range q.columns {
		if i != 0 {
			query += ", "
		}
		query += v + ", " + v + q.buildInterval()
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
