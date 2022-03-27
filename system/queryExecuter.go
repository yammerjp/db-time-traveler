package system

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DatabaseConnection struct {
	connection *sql.DB
}

func (dap *DatabaseAccessPoint) CreateDatabaseConnection() (*DatabaseConnection, error) {
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

func (c *DatabaseConnection) queryExecWithNothingReturningValues(query string) error {
	_, err := c.connection.Exec(query)
	return err
}

func (c *DatabaseConnection) queryExecWithReturningSingleValue(query string) ([]string, error) {
	rows, err := c.queryExec(query)
	if err != nil {
		return []string{}, err
	}
	columns := make([]string, 0, len(rows))
	for _, v := range rows {
		columns = append(columns, v[0])
	}
	return columns, nil
}

func (c *DatabaseConnection) queryExec(query string) ([][]string, error) {
	rows, err := c.connection.Query(query)
	if err != nil {
		return nil, err
	}
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnNullables := make([]sql.NullString, len(columnNames))
	columns := make([]string, len(columnNames))
	columnRefs := make([]interface{}, len(columnNames))
	for i := range columnNullables {
		columnRefs[i] = &columnNullables[i]
	}

	ret := [][]string{}
	for rows.Next() {
		if err := rows.Scan(columnRefs...); err != nil {
			return [][]string{}, err
		}
		for i := range columnNullables {
			if columnNullables[i].Valid {
				columns[i] = columnNullables[i].String
			} else {
				columns[i] = "NULL"
			}
		}
		ret = append(ret, append([]string{}, columns...))
	}
	return ret, nil
}
