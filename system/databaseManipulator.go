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
