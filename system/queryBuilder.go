package system

import (
	"errors"
	"fmt"
	"strings"
)

type Interval struct {
	IsPast bool
	Num    int
	Term   string
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
	Interval
}

func (q *Interval) buildInterval() string {
	var symbol string
	if q.IsPast {
		symbol = "-"
	} else {
		symbol = "+"
	}
	return fmt.Sprintf(" %s INTERVAL %d %s", symbol, q.Num, q.Term)
}

func (q QueryBuilderSourceForSchemaInformation) buildFrom() string {
	return " FROM " + q.targetTable
}

func (q QueryBuilderSourceForColumnValues) buildWhereIn() string {
	pks := strings.Join(q.primaryKeys, ", ")

	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(q.stmtInWhereIn)), "SELECT") {
		return fmt.Sprintf(" WHERE %s IN ( %s )", pks, q.stmtInWhereIn)
	}

	// MySQL does not allow specifying a table with the same name in WHERE IN SELECT during UPDATE, so an alias is applied to avoid this problem.
	return fmt.Sprintf(" WHERE %s IN ( SELECT %s FROM ( %s ) as any )", pks, pks, q.stmtInWhereIn)
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

func (q QueryBuilderSourceForSchemaInformation) buildStmtToSelectColumnNamesDateRelated(ignoreColumnNames []string) (string, error) {
	return q.buildStmtToSelectColumnNames() + " AND DATA_TYPE IN (\"date\", \"datetime\", \"timestamp\") AND COLUMN_NAME NOT IN (\"" + strings.Join(ignoreColumnNames, "\", \"") + "\")", nil // + DATA_TYPE = time
}

func (q QueryBuilderSourceForSchemaInformation) buildStmtToSelectColumnNamesOfPrimaryKey() (string, error) {
	return q.buildStmtToSelectColumnNames() + " AND COLUMN_KEY = \"PRI\"", nil
}
