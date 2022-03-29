package system

import (
	"fmt"
)

func (c *DatabaseConnection) SelectDateRelatedColumns(table string) ([]string, error) {
	q := QueryBuilderSourceForSchemaInformation{
		targetTable: table,
	}
	query, err := q.buildStmtToSelectColumnNamesDateRelated()
	if err != nil {
		return []string{}, err
	}
	return c.queryExecWithReturningSingleValue(query)
}

func (c *DatabaseConnection) SelectPrimaryKeyColumns(table string) ([]string, error) {
	q := QueryBuilderSourceForSchemaInformation{
		targetTable: table,
	}
	query, err := q.buildStmtToSelectColumnNamesOfPrimaryKey()
	if err != nil {
		return []string{}, err
	}
	primaryKeys, err := c.queryExecWithReturningSingleValue(query)
	if err != nil {
		return []string{}, err
	}
	return primaryKeys, nil
}

func (c *DatabaseConnection) SelectDateRelatedColumnValues(table string, stmtInWhereIn string) ([]string, [][]string, error) {
	columns, err := c.SelectDateRelatedColumns(table)
	if err != nil {
		return []string{}, [][]string{}, err
	}

	pks, err := c.SelectPrimaryKeyColumns(table)
	if err != nil {
		return columns, [][]string{}, err
	}
	q := QueryBuilderSourceForColumnValues{
		QueryBuilderSourceForSchemaInformation: QueryBuilderSourceForSchemaInformation{
			targetTable: table,
		},
		columns:       columns,
		primaryKeys:   pks,
		stmtInWhereIn: stmtInWhereIn,
	}
	query, err := q.buildStmtToSelectColumnNamesDateRelated()
	if err != nil {
		return columns, [][]string{}, err
	}
	columnValues, err := c.queryExec(query)
	if err != nil {
		return columns, [][]string{}, err
	}
	return columns, columnValues, err
}

func (c *DatabaseConnection) SelectDateRelatedColumnValuesToBeUpdated(table string, interval QueryBuilderSourcePartOfInterval, stmtInWhereIn string) ([]string, [][]string, error) {
	columns, err := c.SelectDateRelatedColumns(table)
	if err != nil {
		return []string{}, [][]string{}, err
	}

	pks, err := c.SelectPrimaryKeyColumns(table)
	if err != nil {
		return columns, [][]string{}, err
	}
	q := QueryBuilderSourceToUpdate{
		QueryBuilderSourceForColumnValues: QueryBuilderSourceForColumnValues{
			QueryBuilderSourceForSchemaInformation: QueryBuilderSourceForSchemaInformation{
				targetTable: table,
			},
			columns:       columns,
			primaryKeys:   pks,
			stmtInWhereIn: stmtInWhereIn,
		},
		QueryBuilderSourcePartOfInterval: interval,
	}
	query, err := q.buildStmtToSelectBeforeAndAfter()
	if err != nil {
		return columns, [][]string{}, err
	}
	columnValues, err := c.queryExec(query)
	if err != nil {
		return columns, [][]string{}, err
	}
	return columns, columnValues, err
}

func (c *DatabaseConnection) SelectDateRelatedColumnValuesNowAndToBeUpdated(table string, interval QueryBuilderSourcePartOfInterval, stmtInWhereIn string) ([]string, [][]string, [][]string, error) {
	columns, err := c.SelectDateRelatedColumns(table)
	if err != nil {
		return []string{}, [][]string{}, [][]string{}, err
	}

	pks, err := c.SelectPrimaryKeyColumns(table)
	if err != nil {
		return columns, [][]string{}, [][]string{}, err
	}
	q := QueryBuilderSourceForColumnValues{
		QueryBuilderSourceForSchemaInformation: QueryBuilderSourceForSchemaInformation{
			targetTable: table,
		},
		columns:       columns,
		primaryKeys:   pks,
		stmtInWhereIn: stmtInWhereIn,
	}
	query, err := q.buildStmtToSelect()
	if err != nil {
		return columns, [][]string{}, [][]string{}, err
	}
	columnValues, err := c.queryExec(query)
	if err != nil {
		return columns, [][]string{}, [][]string{}, err
	}

	q4u := QueryBuilderSourceToUpdate{
		QueryBuilderSourceForColumnValues: q,
		QueryBuilderSourcePartOfInterval:  interval,
	}
	query, err = q4u.buildStmtToSelect()
	if err != nil {
		return columns, columnValues, [][]string{}, err
	}
	columnValuesToBeUpdated, err := c.queryExec(query)
	if err != nil {
		return columns, columnValues, [][]string{}, err
	}
	return columns, columnValues, columnValuesToBeUpdated, nil
}

func (c *DatabaseConnection) SelectToUpdateQueryBuilder(table string, interval QueryBuilderSourcePartOfInterval, stmtInWhereIn string) (string, []string, error) {
	columns, err := c.SelectDateRelatedColumns(table)
	if err != nil {
		return "", []string{}, err
	}

	pks, err := c.SelectPrimaryKeyColumns(table)
	if err != nil {
		return "", columns, err
	}
	q := QueryBuilderSourceToUpdate{
		QueryBuilderSourceForColumnValues: QueryBuilderSourceForColumnValues{
			QueryBuilderSourceForSchemaInformation: QueryBuilderSourceForSchemaInformation{
				targetTable: table,
			},
			columns:       columns,
			primaryKeys:   pks,
			stmtInWhereIn: stmtInWhereIn,
		},
		QueryBuilderSourcePartOfInterval: interval,
	}
	query, err := q.buildStmtToSelectBeforeAndAfter()
	return query, columns, err
}

func (c *DatabaseConnection) SelectToUpdate(table string, interval QueryBuilderSourcePartOfInterval, stmtInWhereIn string) ([]string, [][]string, []string, [][]string, [][]string, error) {
	primaryKeys, err := c.SelectPrimaryKeyColumns(table)
	if err != nil {
		return []string{}, [][]string{}, []string{}, [][]string{}, [][]string{}, err
	}
	query, columns, err := c.SelectToUpdateQueryBuilder(table, interval, stmtInWhereIn)
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

func (c *DatabaseConnection) SelectToUpdateToString(table string, interval QueryBuilderSourcePartOfInterval, stmtInWhereIn string) (string, error) {
	ret := ""
	primaryKeys, stmtInWhereIns, columns, columnValuesBefore, columnValuesAfter, err := c.SelectToUpdate(table, interval, stmtInWhereIn)
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

func (c *DatabaseConnection) UpdateQueryBuilder(table string, interval QueryBuilderSourcePartOfInterval, stmtInWhereIn string) (string, error) {
	columns, err := c.SelectDateRelatedColumns(table)
	if err != nil {
		return "", err
	}

	pks, err := c.SelectPrimaryKeyColumns(table)
	if err != nil {
		return "", err
	}
	q := QueryBuilderSourceToUpdate{
		QueryBuilderSourceForColumnValues: QueryBuilderSourceForColumnValues{

			QueryBuilderSourceForSchemaInformation: QueryBuilderSourceForSchemaInformation{
				targetTable: table,
			},
			columns:       columns,
			primaryKeys:   pks,
			stmtInWhereIn: stmtInWhereIn,
		},
		QueryBuilderSourcePartOfInterval: interval,
	}
	return q.buildStmtToUpdate()
}

func (c *DatabaseConnection) Update(table string, interval QueryBuilderSourcePartOfInterval, stmtInWhereIn string) error {
	query, err := c.UpdateQueryBuilder(table, interval, stmtInWhereIn)
	if err != nil {
		return err
	}
	return c.queryExecWithNothingReturningValues(query)
}
