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

type Table struct {
	targetTable string
}

type SelectSource struct {
	Table
	columns       []string
	primaryKeys   []string
	stmtInWhereIn string
}

type UpdateSource struct {
	SelectSource
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

func (q Table) buildFrom() string {
	return " FROM " + q.targetTable
}

func (q SelectSource) buildWhereIn() string {
	pks := strings.Join(q.primaryKeys, ", ")

	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(q.stmtInWhereIn)), "SELECT") {
		return fmt.Sprintf(" WHERE %s IN ( %s )", pks, q.stmtInWhereIn)
	}

	// MySQL does not allow specifying a table with the same name in WHERE IN SELECT during UPDATE, so an alias is applied to avoid this problem.
	return fmt.Sprintf(" WHERE %s IN ( SELECT %s FROM ( %s ) as any )", pks, pks, q.stmtInWhereIn)
}

func (q UpdateSource) buildStmtToUpdate() (string, error) {
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

func (q SelectSource) buildStmtToSelect() (string, error) {
	if len(q.columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	return "SELECT " + strings.Join(q.columns, ", ") + q.buildFrom() + q.buildWhereIn(), nil
}

func (q UpdateSource) buildStmtToSelectBeforeAndAfter() (string, error) {
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

func (q Table) buildStmtToSelectColumnNames() string {
	return `SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = "` + q.targetTable + `"`
}

func (q Table) buildStmtToSelectColumnNamesDateRelated(ignoreColumnNames []string) (string, error) {
	return q.buildStmtToSelectColumnNames() + " AND DATA_TYPE IN (\"date\", \"datetime\", \"timestamp\") AND COLUMN_NAME NOT IN (\"" + strings.Join(ignoreColumnNames, "\", \"") + "\")", nil // + DATA_TYPE = time
}

func (q Table) buildStmtToSelectColumnNamesOfPrimaryKey() (string, error) {
	return q.buildStmtToSelectColumnNames() + " AND COLUMN_KEY = \"PRI\"", nil
}
