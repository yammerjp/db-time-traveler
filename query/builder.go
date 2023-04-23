package query

import (
	"errors"
	"strings"
)

func (u UpdateSource) BuildStmtToUpdate() (UpdateStatement, error) {
	if len(u.Columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	query := "UPDATE " + u.TargetTable + " SET "
	for i, column := range u.Columns {
		if i != 0 {
			query += ", "
		}
		query += column + " = (" + column + u.buildInterval() + ")"
	}
	return UpdateStatement(query + u.buildWhereIn()), nil
}

func (q UpdateSource) BuildStmtToSelectBeforeAndAfter() (SelectStatement, error) {
	if len(q.Columns) == 0 {
		return "", errors.New("must be specify any columns")
	}
	query := "SELECT " + strings.Join(q.PrimaryKeys, ", ") + ", "
	for i, v := range q.Columns {
		if i != 0 {
			query += ", "
		}
		query += v + ", " + v + q.buildInterval()
	}
	return SelectStatement(query + q.buildFrom() + q.buildWhereIn()), nil
}

func (t Table) buildStmtToSelectColumnNames() SelectStatement {
	return SelectStatement(`SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = "` + t.TargetTable + `"`)
}

func (t Table) BuildStmtToSelectColumnNamesDateRelated(ignoreColumnNames []string) SelectStatement {
	return SelectStatement(string(t.buildStmtToSelectColumnNames()) + ` AND DATA_TYPE IN ("date", "datetime", "timestamp") AND COLUMN_NAME NOT IN ("` + strings.Join(ignoreColumnNames, `", "`) + `")`) // + DATA_TYPE = time
}

func (t Table) BuildStmtToSelectColumnNamesOfPrimaryKey() SelectStatement {
	return t.buildStmtToSelectColumnNames() + ` AND COLUMN_KEY = "PRI"`
}
