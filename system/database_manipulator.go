package system

import (
	"github.com/yammerjp/db-time-traveler/database"
	"github.com/yammerjp/db-time-traveler/query"
)

func selectDateRelatedColumns(c *database.Connection, table string, ignoreColumns []string) ([]string, error) {
	query := query.Table{TargetTable: table}.BuildStmtToSelectColumnNamesDateRelated(ignoreColumns)
	// SelectStatement to string
	return c.QueryExecWithReturningSingleValue(string(query))
}

func selectPrimaryKeyColumns(c *database.Connection, table string) ([]string, error) {
	query := query.Table{TargetTable: table}.BuildStmtToSelectColumnNamesOfPrimaryKey()
	// SelectStatement to string
	return c.QueryExecWithReturningSingleValue(string(query))
}
