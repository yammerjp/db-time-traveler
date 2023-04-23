package system

import (
	"fmt"

	"github.com/yammerjp/db-time-traveler/database"
	"github.com/yammerjp/db-time-traveler/query"
)

func BeforeAndAfter(c *database.Connection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) (string, error) {
	ret := ""
	pks, stmtInWhereIns, columns, beforeValsArr, afterValsArr, err := beforeAndAfter(c, table, interval, stmtInWhereIn, ignoreColumns)
	if err != nil {
		return "", nil
	}
	for i, v := range stmtInWhereIns {
		if i != 0 {
			ret += "\n"
		}
		for j, w := range pks {
			if j != 0 {
				ret += ", "
			}
			ret += fmt.Sprintf("%s: %s", w, v[j])
		}
		for j := range columns {
			ret += fmt.Sprintf("\n  %s:\n    before: %s\n    after:  %s", columns[j], beforeValsArr[i][j], afterValsArr[i][j])
		}
	}
	return ret, nil
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

func Update(c *database.Connection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) error {
	query, err := BuildQueryUpdate(c, table, interval, stmtInWhereIn, ignoreColumns)
	if err != nil {
		return err
	}
	// UpdateStatement to string
	return c.QueryExecWithNothingReturningValues(string(query))
}

func selectDateRelatedColumns(c *database.Connection, table string, ignoreColumns []string) ([]string, error) {
	query, err := query.Table{TargetTable: table}.BuildStmtToSelectColumnNamesDateRelated(ignoreColumns)
	if err != nil {
		return []string{}, err
	}
	// SelectStatement to string
	return c.QueryExecWithReturningSingleValue(string(query))
}

func selectPrimaryKeyColumns(c *database.Connection, table string) ([]string, error) {
	query, err := query.Table{TargetTable: table}.BuildStmtToSelectColumnNamesOfPrimaryKey()
	if err != nil {
		return []string{}, err
	}
	// SelectStatement to string
	primaryKeys, err := c.QueryExecWithReturningSingleValue(string(query))
	if err != nil {
		return []string{}, err
	}
	return primaryKeys, nil
}

func BuildStmtToSelectBeforeAndAfter(c *database.Connection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) (query.SelectStatement, []string, error) {
	columns, err := selectDateRelatedColumns(c, table, ignoreColumns)
	if err != nil {
		return "", []string{}, err
	}

	pks, err := selectPrimaryKeyColumns(c, table)
	if err != nil {
		return "", columns, err
	}
	stmt, err := query.UpdateSource{
		SelectSource: query.SelectSource{
			Table: query.Table{
				TargetTable: table,
			},
			Columns:       columns,
			PrimaryKeys:   pks,
			StmtInWhereIn: stmtInWhereIn,
		},
		Interval: interval,
	}.BuildStmtToSelectBeforeAndAfter()

	return stmt, columns, err
}

func selectBeforeAndAfter(c *database.Connection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) ([][]string, []string, error) {
	query, columns, err := BuildStmtToSelectBeforeAndAfter(c, table, interval, stmtInWhereIn, ignoreColumns)
	if err != nil {
		return [][]string{}, columns, err
	}
	// SelectStatement to string
	result, err := c.QueryExec(string(query))
	if err != nil {
		return [][]string{}, columns, err
	}
	return result, columns, nil
}

func beforeAndAfter(c *database.Connection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) ([]string, [][]string, []string, [][]string, [][]string, error) {
	pkNames, err := selectPrimaryKeyColumns(c, table)
	if err != nil {
		return []string{}, [][]string{}, []string{}, [][]string{}, [][]string{}, err
	}
	result, columnsToUpdate, err := selectBeforeAndAfter(c, table, interval, stmtInWhereIn, ignoreColumns)
	if err != nil {
		return pkNames, [][]string{}, columnsToUpdate, [][]string{}, [][]string{}, err
	}
	pkValsArr := [][]string{}
	beforeValsArr := [][]string{}
	afterValsArr := [][]string{}
	/*
		ex)
			SELECT
				id,
				trial_end_date, trial_end_date - INTERVAL 1 MONTH,
				registered_campaign_end_datetime, registered_campaign_end_datetime - INTERVAL 1 MONTH,
				created_at, created_at - INTERVAL 1 MONTH,
				updated_at, updated_at - INTERVAL 1 MONTH
			FROM accounts
			WHERE id IN ( 3 )
	*/
	for _, v := range result {
		pkVals := []string{}
		beforeVals := []string{}
		afterVals := []string{}
		for j, r := range v {
			if j < len(pkNames) {
				// ex) id
				pkVals = append(pkVals, r)
			} else if (j-len(pkNames))%2 == 0 {
				// ex) trial_end_date, registered_campaign_end_datetime, created_at, updated_at
				beforeVals = append(beforeVals, r)
			} else {
				// ex) trial_end_date - INTERVAL 1 MONTH, registered_campaign_end_datetime - INTERVAL 1 MONTH, created_at - INTERVAL 1 MONTH, updated_at - INTERVAL 1 MONTH
				afterVals = append(afterVals, r)
			}
		}
		beforeValsArr = append(beforeValsArr, append([]string{}, beforeVals...))
		afterValsArr = append(afterValsArr, append([]string{}, afterVals...))
		pkValsArr = append(pkValsArr, append([]string{}, pkVals...))
	}
	return pkNames, pkValsArr, columnsToUpdate, beforeValsArr, afterValsArr, nil
}
