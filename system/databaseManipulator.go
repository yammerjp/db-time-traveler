package system

import (
	"fmt"

	"github.com/yammerjp/db-time-traveler/query"
)

func SelectToUpdateToString(c *DatabaseConnection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) (string, error) {
	ret := ""
	primaryKeys, stmtInWhereIns, columns, columnValuesBefore, columnValuesAfter, err := selectToUpdate(c, table, interval, stmtInWhereIn, ignoreColumns)
	if err != nil {
		return "", nil
	}
	for i, v := range stmtInWhereIns {
		if i != 0 {
			ret += "\n"
		}
		for j, w := range primaryKeys {
			if j != 0 {
				ret += ", "
			}
			ret += fmt.Sprintf("%s: %s", w, v[j])
		}
		for j := range columns {
			ret += fmt.Sprintf("\n  %s:\n    before: %s\n    after:  %s", columns[j], columnValuesBefore[i][j], columnValuesAfter[i][j])
		}
	}
	return ret, nil
}

func UpdateQueryBuilder(c *DatabaseConnection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) (string, error) {
	columns, err := selectDateRelatedColumns(c, table, ignoreColumns)
	if err != nil {
		return "", err
	}

	pks, err := selectPrimaryKeyColumns(c, table)
	if err != nil {
		return "", err
	}
	return query.BuildStmtToUpdate(table, columns, pks, stmtInWhereIn, interval)
}

func Update(c *DatabaseConnection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) error {
	query, err := UpdateQueryBuilder(c, table, interval, stmtInWhereIn, ignoreColumns)
	if err != nil {
		return err
	}
	return c.queryExecWithNothingReturningValues(query)
}

func selectDateRelatedColumns(c *DatabaseConnection, table string, ignoreColumns []string) ([]string, error) {
	query, err := query.BuildStmtToSelectColumnNamesDateRelated(table, ignoreColumns)
	if err != nil {
		return []string{}, err
	}
	return c.queryExecWithReturningSingleValue(query)
}

func selectPrimaryKeyColumns(c *DatabaseConnection, table string) ([]string, error) {
	query, err := query.BuildStmtToSelectColumnNamesOfPrimaryKey(table)
	if err != nil {
		return []string{}, err
	}
	primaryKeys, err := c.queryExecWithReturningSingleValue(query)
	if err != nil {
		return []string{}, err
	}
	return primaryKeys, nil
}

func selectToUpdateQueryBuilder(c *DatabaseConnection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) (string, []string, error) {
	columns, err := selectDateRelatedColumns(c, table, ignoreColumns)
	if err != nil {
		return "", []string{}, err
	}

	pks, err := selectPrimaryKeyColumns(c, table)
	if err != nil {
		return "", columns, err
	}
	query, err := query.BuildStmtToUpdate(table, columns, pks, stmtInWhereIn, interval)
	return query, columns, err
}

func selectToUpdate(c *DatabaseConnection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) ([]string, [][]string, []string, [][]string, [][]string, error) {
	primaryKeys, err := selectPrimaryKeyColumns(c, table)
	if err != nil {
		return []string{}, [][]string{}, []string{}, [][]string{}, [][]string{}, err
	}
	query, columns, err := selectToUpdateQueryBuilder(c, table, interval, stmtInWhereIn, ignoreColumns)
	if err != nil {
		return primaryKeys, [][]string{}, columns, [][]string{}, [][]string{}, err
	}
	selectStmtReturnColumnValues, err := c.queryExec(query)
	if err != nil {
		return primaryKeys, [][]string{}, columns, [][]string{}, [][]string{}, err
	}
	retPrimaryKeys := [][]string{}
	retColumnValuesBefore := [][]string{}
	retColumnValuesAfter := [][]string{}
	for _, v := range selectStmtReturnColumnValues {
		columnPrimaryKeys := []string{}
		columnValuesBefore := []string{}
		columnValuesAfter := []string{}
		for j, r := range v {
			if j < len(primaryKeys) {
				columnPrimaryKeys = append(columnPrimaryKeys, r)
			} else if j%2 == 1 {
				columnValuesBefore = append(columnValuesBefore, r)
			} else {
				columnValuesAfter = append(columnValuesAfter, r)
			}
		}
		retColumnValuesBefore = append(retColumnValuesBefore, append([]string{}, columnValuesBefore...))
		retColumnValuesAfter = append(retColumnValuesAfter, append([]string{}, columnValuesAfter...))
		retPrimaryKeys = append(retPrimaryKeys, append([]string{}, columnPrimaryKeys...))
	}
	return primaryKeys, retPrimaryKeys, columns, retColumnValuesBefore, retColumnValuesAfter, nil
}
