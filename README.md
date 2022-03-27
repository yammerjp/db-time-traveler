db-time-traveler
---

db-time-traveler rewrites the time of records in a relational database together.
You specify a table name, and db-time-traveler will rewind that record by the specified amount of time, as if it had been created in the past.

For columns of type DATETIME or TIMESTAMP in the specified table, subtract the specified time from the existing value and overwrite the record with the result.

## Requirement

To develop a function that requires operations to be performed chronologically, we would like to pretend that the records were created in the past for the purpose of debugging in the development environment.

I want to run an UPDATE statement to shift the values of DATETIME and TIMESTAMP type columns to the past by the time specified by the input, after narrowing down the records by specifying the table name.

In addition to the above, it would be nice to have the ability to check the SQL before execution and to rewind the contents of a record after SQL execution as auxiliary functions.

## Specification

It is a CLI that can execute the following commands

```bash
db-time-traveler \
  update
  --past 1month \
  --host 127.0.0.1 \
  --port 3306 \
  --user username \
  --password password \
  --schema main-db \
  --table accounts \
  --primary-key-raw 35 \
```

The above command example represents the following actions

- Connect to MySQL whose IP address is 127.0.0.1 and port number is 3306
- Manipulate a database named main-db
- The records that are filtered by the primary key "35"
- Exclude the 'updated_at' column of the "contracts" table from being overwritten
- Rewinds the date or time of all columns of the first specified table ("accounts") that are of type DATE or DATETIME or TIMESTAMP by one month.

## Commandline Options

### past

The rewind time

Acceptable: `[0-9]+(month|week|day)`

Example:

```
  --past 2week
```

### host

Host name or IPv4 address to connect to

Acceptable: `[0-9]+\.[0-9]+\.[0-9]\.[0-9]`
Acceptable: `[a-z][a-z0-9\-.]*`

Example:

```
  --host example.com
```

### port

Destination port number

Acceptable: `[0-9]+`

Example:

```
  --port 3306
```

### user

Username for Database Connection

Acceptable: `any string`

Example:

```
  --user root
```

### password

Password for Database Connection

Acceptable: `any string`

Example:

```
  --password pa55w0rd
```

### schema

Database schema name to connect to

Acceptable: `[a-z0-9\-]+`

Example:

```
  --schema maindb
```

### table

Central table name (the table specified in the FROM clause of a SELECT statement or the UPDATE clause of an UPDATE statement)

Acceptable: `<table name>`

Example:

```
  --table accounts
```

### ignore (unimplemented)

Tables and columns to exclude from updates

Acceptable: `<table name>`
Acceptable: `<table name>.<column name>`

Example:

```
  --ignore accounts.updated_at
```

### rollback (unimplemented)

Revert to before UPDATE

### logfile (unimplemented)

Log output destination

### verbose (unimplemented)

Verbose output of execution steps and results

### dry-run

Check the UPDATE statement to be executed

## Author

Keisuke Nakayama ([@yammerjp](https://github.com/yammerjp))

## LICENSE

MIT
