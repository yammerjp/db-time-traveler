package system

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

func SelectAndPrintColumns(dap *DatabaseAccessPoint, tableName string) {
	trcs, err := dap.SelectColumns(tableName)
	if err != nil {
		log.Fatal(err)
	}
	lines := []string{}
	for _, v := range trcs {
		lines = append(lines, v.ToString())
	}
	fmt.Print(strings.Join(lines, "\n"))
}

func (dap *DatabaseAccessPoint) SelectColumns(tableName string) ([]TimeRelatedColumn, error) {
	db, err := dap.connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return SelectColumns(db, tableName)
}

func SelectColumns(db *sql.DB, tableName string) ([]TimeRelatedColumn, error) {
	query := "SELECT COLUMN_NAME, DATA_TYPE, COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS"
	params := []interface{}{}

	if tableName != "" {
		query += " WHERE table_name = ?"
		params = append(params, tableName)
	}

	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}

	timeRelatedColumns := []TimeRelatedColumn{}
	for rows.Next() {
		timeRelatedColumn := TimeRelatedColumn{}
		var column string
		var dataType string
		var columnKey string
		if err := rows.Scan(&column, &dataType, &columnKey); err != nil {
			log.Fatal(err)
		}
		switch dataType {
		case "date":
			timeRelatedColumn.ColumnType = Date
		case "time":
			timeRelatedColumn.ColumnType = Time
		case "datetime":
			timeRelatedColumn.ColumnType = DateTime
		case "timestamp":
			timeRelatedColumn.ColumnType = Timestamp
		default:
			continue
		}
		timeRelatedColumn.IsPrimary = columnKey == "PRI"
		timeRelatedColumn.Name = column

		timeRelatedColumns = append(timeRelatedColumns, timeRelatedColumn)
	}
	return timeRelatedColumns, nil
}

func (dap *DatabaseAccessPoint) SelectToUpdate(tableName string, columnNames []TimeRelatedColumn, wheres []WhereClause) error {
	db, err := dap.connect()
	if err != nil {
		return err
	}
	defer db.Close()
	return SelectToUpdate(db, tableName, columnNames, wheres)
}

func queryBuilder(tableName string, columnNames []TimeRelatedColumn, wheres []WhereClause) (string, []interface{}, error) {
	query := "SELECT"
	params := []interface{}{}

	for i, v := range columnNames {
		if i == 0 {
			query += " ?"
		} else {
			query += ", ?"
		}
		params = append(params, v.Name)
	}

	if tableName == "" {
		return "", nil, errors.New("Need table name")
	}
	// CAUTION: STRING CONCATENATION
	query += " FROM " + tableName

	for i, v := range wheres {
		if v.Operator != "=" {
			return "", nil, errors.New("Unknown where operator")
		}
		if i == 0 {
			query += " WHERE ? = ?"
		} else {
			query += " AND ? = ?"
		}
		params = append(params, v.LeftHand)
		params = append(params, v.RightHand)
	}

	fmt.Println(query)
	for i, v := range params {
		if i == 0 {
			fmt.Printf("%s", v)
		} else {
			fmt.Printf(" %s", v)
		}
	}
	fmt.Printf("\n")
	return query, params, nil
}

func SelectToUpdate(db *sql.DB, tableName string, columnNames []TimeRelatedColumn, wheres []WhereClause) error {
	query, params, err := queryBuilder(tableName, columnNames, wheres)
	columnValues := []interface{}{}
	for i := 0; i < len(columnNames); i++ {
		var columnValue sql.NullString
		columnValues = append(columnValues, &columnValue)
	}
	if err != nil {
		return err
	}
	rows, err := db.Query(query, params...)
	if err != nil {
		return err
	}
	/*
	  results := make([]interface{}, len(columnNames))
	  for i := range results {
	    results[i] = new(interface{})
	  }
	  pretty := [][]string{}
	*/
	for rows.Next() {
		// if err := rows.Scan(results[:]...); err != nil {
		if err := rows.Scan(columnValues...); err != nil {
			return err
		}
		for _, v := range columnValues {
			fmt.Printf("%s", v)
		}
		/*
		   cur := make([]string, len(columnNames))
		   for i := range results {

		     val := *results[i].(*interface{})
		     cur[i] = fmt.Sprintf("%s", val)
		     fmt.Printf("%s", *results[i].(*interface{}))
		   }
		   pretty = append(pretty, cur)
		*/
	}
	/*
	  for i, v := range results {
	    columnValues := []string{}
	    for _, w := range v {
	      columnValues = append(columnValues, fmt.printf("%s", w))
	    }

	    fmt.Printf("%4d: %s\n", i, strings.Join(columnValues, ", "))
	  }
	*/
	return nil
}

func updateQueryBuilder(tableName string, columnNames []TimeRelatedColumn, wheres []WhereClause) (string, []interface{}, error) {
	query := "UPDATE " + tableName + " SET"
	params := []interface{}{}
	for i, v := range columnNames {
		if i == 0 {
			query += fmt.Sprintf(" %s = (%s - INTERVAL 1 MONTH)", v.Name, v.Name)
			// query += " ? = (? - INTERVAL 1 MONTH)"
		} else {
			query += fmt.Sprintf(", %s = (%s - INTERVAL 1 MONTH)", v.Name, v.Name)
			//query += ", ? = (? - INTERVAL 1 MONTH)"
		}
		// params = append(params, v.Name, v.Name)
	}

	for i, v := range wheres {
		if v.Operator != "=" {
			return "", nil, errors.New("Unknown where operator")
		}
		if i == 0 {
			query += " WHERE `?` = `?`"
		} else {
			query += " AND ? = ?"
		}
		params = append(params, v.LeftHand, v.RightHand)
	}
	return query, params, nil
}

func Update(db *sql.DB, tableName string, columnNames []TimeRelatedColumn, wheres []WhereClause) error {
	query, params, err := updateQueryBuilder(tableName, columnNames, wheres)
	if err != nil {
		return err
	}
	result, err := db.Exec(query, params...)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Printf("Last Insert Id: %10d", id)
	return nil
}

/*
TODO: implement to run update statement
  UPDATE
    accounts
  SET
    trial_end_date = (trail_end_date - INTERVAL 1 MONTH),
    registered_campaign_end_datetime = (registered_campaign_end_datetime - INTERVAL 1 MONTH),
    created_at = (created_at - INTERVAL 1 MONTH),
    updated_at = (updated_at - INTERVAL 1 MONTH)
  WHERE
    account_id = ...
*/
