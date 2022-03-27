package system

import (
	"errors"
)

func (c *DatabaseConnection) SelectDateRelatedColumns(table string) ([]string, error) {
	query, err := selectDateRelatedColumnsQueryBuilder(table)
	if err != nil {
		return []string{}, err
	}
	return c.queryExecWithReturningSingleValue(query)
}

func (c *DatabaseConnection) SelectPrimaryKeyColumns(table string) ([]string, error) {
	query, err := selectPrimaryKeyColumnsQueryBuilder(table)
	primaryKeys, err := c.queryExecWithReturningSingleValue(query)
	if err != nil {
		return []string{}, err
	}
	return primaryKeys, nil
}

func (c *DatabaseConnection) SelectDateRelatedColumnValues(table string, primaryKeyValue string) ([]string, [][]string, error) {
	columns, err := c.SelectDateRelatedColumns(table)
	if err != nil {
		return []string{}, [][]string{}, err
	}

	pks, err := c.SelectPrimaryKeyColumns(table)
	if err != nil {
		return columns, [][]string{}, err
	}
	if len(pks) != 1 {
		return columns, [][]string{}, errors.New("Support only single primary key")
	}

	query, err := selectTargettedColumnsQueryBuilder(table, columns, pks[0], primaryKeyValue)
	columnValues, err := c.queryExec(query)
	if err != nil {
		return columns, [][]string{}, err
	}
	return columns, columnValues, err
}

func (c *DatabaseConnection) SelectDateRelatedColumnValuesToBeUpdated(table string, interval string, primaryKeyValue string) ([]string, [][]string, error) {
	columns, err := c.SelectDateRelatedColumns(table)
	if err != nil {
		return []string{}, [][]string{}, err
	}

	pks, err := c.SelectPrimaryKeyColumns(table)
	if err != nil {
		return columns, [][]string{}, err
	}
	if len(pks) != 1 {
		return columns, [][]string{}, errors.New("Support only single primary key")
	}

	query, err := selectUpdatingColumnValuesQueryBuilder(table, columns, interval, pks[0], primaryKeyValue)
	columnValues, err := c.queryExec(query)
	if err != nil {
		return columns, [][]string{}, err
	}
	return columns, columnValues, err
}

func (c *DatabaseConnection) SelectDateRelatedColumnValuesNowAndToBeUpdated(table string, interval string, primaryKeyValue string) ([]string, [][]string, [][]string, error) {
	columns, err := c.SelectDateRelatedColumns(table)
	if err != nil {
		return []string{}, [][]string{}, [][]string{}, err
	}

	pks, err := c.SelectPrimaryKeyColumns(table)
	if err != nil {
		return columns, [][]string{}, [][]string{}, err
	}
	if len(pks) != 1 {
		return columns, [][]string{}, [][]string{}, errors.New("Support only single primary key")
	}

	query, err := selectTargettedColumnsQueryBuilder(table, columns, pks[0], primaryKeyValue)
	columnValues, err := c.queryExec(query)
	if err != nil {
		return columns, [][]string{}, [][]string{}, err
	}

	query, err = selectUpdatingColumnValuesQueryBuilder(table, columns, interval, pks[0], primaryKeyValue)
	columnValuesToBeUpdated, err := c.queryExec(query)
	if err != nil {
		return columns, columnValues, [][]string{}, err
	}
	return columns, columnValues, columnValuesToBeUpdated, nil
}

func (c *DatabaseConnection) SelectToUpdate(table string, interval string, primaryKeyValue string) ([]string, []string, [][]string, [][]string, error) {
	columns, err := c.SelectDateRelatedColumns(table)
	if err != nil {
		return []string{}, []string{}, [][]string{}, [][]string{}, err
	}

	pks, err := c.SelectPrimaryKeyColumns(table)
	if err != nil {
		return []string{}, columns, [][]string{}, [][]string{}, err
	}
	if len(pks) != 1 {
		return []string{}, columns, [][]string{}, [][]string{}, errors.New("Support only single primary key")
	}

	query, err := selectUpdatingColumnValuesBeforeAndAfterQueryBuilder(table, columns, interval, pks[0], primaryKeyValue)
	if err != nil {
		return []string{}, columns, [][]string{}, [][]string{}, err
	}
	selectStmtReturnColumnValues, err := c.queryExec(query)
	if err != nil {
		return []string{}, columns, [][]string{}, [][]string{}, err
	}
	retPrimaryKeys := []string{}
	retColumnValuesBefore := [][]string{}
	retColumnValuesAfter := [][]string{}
	for _, v := range selectStmtReturnColumnValues {
		columnValuesBefore := []string{}
		columnValuesAfter := []string{}
		for j, r := range v {
			if j == 0 {
				retPrimaryKeys = append(retPrimaryKeys, r)
			} else if j%2 == 1 {
				columnValuesBefore = append(columnValuesBefore, r)
			} else {
				columnValuesAfter = append(columnValuesAfter, r)
			}
		}
		retColumnValuesBefore = append(retColumnValuesBefore, append([]string{}, columnValuesBefore...))
		retColumnValuesAfter = append(retColumnValuesAfter, append([]string{}, columnValuesAfter...))
	}
	return retPrimaryKeys, columns, retColumnValuesBefore, retColumnValuesAfter, nil
}

func (c *DatabaseConnection) Update(table string, interval string, primaryKeyValue string) error {
	columns, err := c.SelectDateRelatedColumns(table)
	if err != nil {
		return err
	}

	pks, err := c.SelectPrimaryKeyColumns(table)
	if err != nil {
		return err
	}
	if len(pks) != 1 {
		return errors.New("Support only single primary key")
	}

	query, err := updateQueryBuilder(table, columns, interval, pks[0], primaryKeyValue)
	if err != nil {
		return err
	}
	return c.queryExecWithNothingReturningValues(query)
}
