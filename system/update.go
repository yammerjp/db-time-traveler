package system

import (
	"github.com/yammerjp/db-time-traveler/database"
	"github.com/yammerjp/db-time-traveler/query"
)

func Update(c *database.Connection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) error {
	query, err := BuildQueryUpdate(c, table, interval, stmtInWhereIn, ignoreColumns)
	if err != nil {
		return err
	}
	// UpdateStatement to string
	return c.QueryExecWithNothingReturningValues(string(query))
}

func BuildQueryUpdate(c *database.Connection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) (query.UpdateStatement, error) {
	columns, err := selectDateRelatedColumns(c, table, ignoreColumns)
	if err != nil {
		return "", err
	}

	pks, err := selectPrimaryKeyColumns(c, table)
	if err != nil {
		return "", err
	}

	return query.UpdateSource{
		SelectSource: query.SelectSource{
			Table: query.Table{
				TargetTable: table,
			},
			Columns:       columns,
			PrimaryKeys:   pks,
			StmtInWhereIn: stmtInWhereIn,
		},
		Interval: interval,
	}.BuildStmtToUpdate()
}
