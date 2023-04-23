package system

import (
	"fmt"

	"github.com/yammerjp/db-time-traveler/database"
	"github.com/yammerjp/db-time-traveler/query"
)

type Preview struct {
	pkNames []string
	rows    []Row
}

type Row struct {
	pkVals          []string
	updatingColumns []UpdatingColumn
}

type UpdatingColumn struct {
	name      string
	beforeVal string
	afterVal  string
}

func BeforeAndAfter(c *database.Connection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) (string, error) {
	ret := ""
	preview, err := beforeAndAfter(c, table, interval, stmtInWhereIn, ignoreColumns)
	if err != nil {
		return "", nil
	}
	for i, v := range preview.rows {
		if i != 0 {
			ret += "\n"
		}
		for j, w := range preview.pkNames {
			if j != 0 {
				ret += ", "
			}
			ret += fmt.Sprintf("%s: %s", w, v.pkVals[j])
		}
		for _, x := range v.updatingColumns {
			ret += fmt.Sprintf("\n  %s:\n    before: %s\n    after:  %s", x.name, x.beforeVal, x.afterVal)
		}
	}
	return ret, nil
	/*
			id: 1
			  registered_campaign_end_datetime:
		        before: 2024-12-03T00:10:20Z
		        after:  2025-01-03T00:10:20Z
			  created_at:
		        before: 2025-12-22T23:58:58Z
		        after:  2026-01-22T23:58:58Z
	*/
}

func beforeAndAfter(c *database.Connection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) (*Preview, error) {
	preview := Preview{}

	pkNames, err := selectPrimaryKeyColumns(c, table)
	if err != nil {
		return nil, err
	}
	preview.pkNames = pkNames

	result, columnsToUpdate, err := selectBeforeAndAfter(c, table, interval, stmtInWhereIn, ignoreColumns)
	if err != nil {
		return nil, err
	}
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

		row := &Row{}
		col := &UpdatingColumn{}
		for j, r := range v {
			if j < len(pkNames) {
				// ex) id
				pkVals = append(pkVals, r)
			} else if (j-len(pkNames))%2 == 0 {
				// ex) trial_end_date, registered_campaign_end_datetime, created_at, updated_at
				col = &UpdatingColumn{
					name:      columnsToUpdate[(j-len(pkNames))/2],
					beforeVal: r,
					afterVal:  "",
				}
			} else {
				// ex) trial_end_date - INTERVAL 1 MONTH, registered_campaign_end_datetime - INTERVAL 1 MONTH, created_at - INTERVAL 1 MONTH, updated_at - INTERVAL 1 MONTH
				col.afterVal = r
				row.updatingColumns = append(row.updatingColumns, *col)
			}
		}
		row.pkVals = pkVals
		preview.rows = append(preview.rows, *row)
	}
	return &preview, nil
}

func selectBeforeAndAfter(c *database.Connection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) ([][]string, []string, error) {
	query, columns, err := buildStmtToSelectBeforeAndAfter(c, table, interval, stmtInWhereIn, ignoreColumns)
	if err != nil {
		return [][]string{}, []string{}, err
	}
	// SelectStatement to string
	result, err := c.QueryExec(string(query))
	if err != nil {
		return [][]string{}, []string{}, err
	}
	return result, columns, nil
}

func buildStmtToSelectBeforeAndAfter(c *database.Connection, table string, interval query.Interval, stmtInWhereIn string, ignoreColumns []string) (query.SelectStatement, []string, error) {
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
