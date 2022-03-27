package system

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DatabaseConnection struct {
	connection *sql.DB
}

func createDatabaseConnection(schema string) (*DatabaseConnection, error) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/"+schema+"?parseTime=true")
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &DatabaseConnection{connection: db}, nil
}

func (dap *DatabaseAccessPoint) createDatabaseConnection() (*DatabaseConnection, error) {
	db, err := sql.Open("mysql", dap.toString())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &DatabaseConnection{connection: db}, nil
}

func (c *DatabaseConnection) Close() {
	c.connection.Close()
}

func (c *DatabaseConnection) ping() error {
	return c.connection.Ping()
}

func ping(schema string) error {
	c, err := createDatabaseConnection(schema)
	if err != nil {
		return err
	}
	defer c.Close()

	return nil
}

func (c *DatabaseConnection) queryProcessorReturnsSingleColumn(query string) ([]string, error) {
	rows, err := c.connection.Query(query)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			return []string{}, err
		}
		ret = append(ret, column)
	}
	return ret, nil
}

func queryProcessorReturnsSingleColumn(schema string, query string) ([]string, error) {
	c, err := createDatabaseConnection(schema)
	if err != nil {
		return []string{}, err
	}
	defer c.Close()

	return c.queryProcessorReturnsSingleColumn(query)
}

func (c *DatabaseConnection) queryProcessorReturnsNothing(query string) error {
	_, err := c.connection.Exec(query)
	return err
}

func queryProcessorReturnsNothing(schema string, query string) error {
	c, err := createDatabaseConnection(schema)
	if err != nil {
		return err
	}
	defer c.Close()

	return c.queryProcessorReturnsNothing(query)
}

func (c *DatabaseConnection) queryProcessorReturn4Columns(query string) ([][]string, error) {
	rows, err := c.connection.Query(query)
	if err != nil {
		return nil, err
	}

	ret := [][]string{}
	for rows.Next() {
		var column0Nullable sql.NullString
		var column1Nullable sql.NullString
		var column2Nullable sql.NullString
		var column3Nullable sql.NullString
		var column0 string
		var column1 string
		var column2 string
		var column3 string

		if err := rows.Scan(&column0Nullable, &column1Nullable, &column2Nullable, &column3Nullable); err != nil {
			return [][]string{}, err
		}
		if column0Nullable.Valid {
			column0 = column0Nullable.String
		} else {
			column0 = "NULL"
		}
		if column1Nullable.Valid {
			column1 = column1Nullable.String
		} else {
			column1 = "NULL"
		}
		if column2Nullable.Valid {
			column2 = column2Nullable.String
		} else {
			column2 = "NULL"
		}
		if column3Nullable.Valid {
			column3 = column3Nullable.String
		} else {
			column3 = "NULL"
		}
		ret = append(ret, []string{column0, column1, column2, column3})
	}
	return ret, nil
}
