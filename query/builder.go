package query

import (
	"errors"
	"fmt"
	"strings"
)

type Statement string
type UpdateStatement Statement
type SelectStatement Statement

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

func (i *Interval) buildInterval() string {
	var symbol string
	if i.IsPast {
		symbol = "-"
	} else {
		symbol = "+"
	}
	return fmt.Sprintf(" %s INTERVAL %d %s", symbol, i.Num, i.Term)
}

func (t Table) buildFrom() string {
	return " FROM " + t.targetTable
}

func (s SelectSource) buildWhereIn() string {
	pks := strings.Join(s.primaryKeys, ", ")

	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(s.stmtInWhereIn)), "SELECT") {
		return fmt.Sprintf(" WHERE %s IN ( %s )", pks, s.stmtInWhereIn)
	}

	// MySQL does not allow specifying a table with the same name in WHERE IN SELECT during UPDATE, so an alias is applied to avoid this problem.
	return fmt.Sprintf(" WHERE %s IN ( SELECT %s FROM ( %s ) as any )", pks, pks, s.stmtInWhereIn)
}

func (u UpdateSource) buildStmtToUpdate() (UpdateStatement, error) {
	if len(u.columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	query := "UPDATE " + u.targetTable + " SET "
	for i, column := range u.columns {
		if i != 0 {
			query += ", "
		}
		query += column + " = (" + column + u.buildInterval() + ")"
	}
	return UpdateStatement(query + u.buildWhereIn()), nil
}

func (q UpdateSource) buildStmtToSelectBeforeAndAfter() (SelectStatement, error) {
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
	return SelectStatement(query + q.buildFrom() + q.buildWhereIn()), nil
}

func (t Table) buildStmtToSelectColumnNames() SelectStatement {
	return SelectStatement(`SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = "` + t.targetTable + `"`)
}

func (t Table) buildStmtToSelectColumnNamesDateRelated(ignoreColumnNames []string) (SelectStatement, error) {
	return SelectStatement(string(t.buildStmtToSelectColumnNames()) + " AND DATA_TYPE IN (\"date\", \"datetime\", \"timestamp\") AND COLUMN_NAME NOT IN (\"" + strings.Join(ignoreColumnNames, "\", \"") + "\")"), nil // + DATA_TYPE = time
}

func (t Table) buildStmtToSelectColumnNamesOfPrimaryKey() (SelectStatement, error) {
	return t.buildStmtToSelectColumnNames() + " AND COLUMN_KEY = \"PRI\"", nil
}