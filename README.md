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
  --database main-db \
  --table accounts \
  --join contracts \
  -on contracts.account_id=accounts.id \
  --where accounts.id=35 \
  --ignore contracts.updated_at
```

The above command example represents the following actions

- Connect to MySQL whose IP address is 127.0.0.1 and port number is 3306
- Manipulate a database named main-db
- JOIN the "contracts" table to the "accounts" table with the condition "contracts.account_id = accounts.id
- The records that are filtered by the condition "accounts.id=35" are the target of the operation.
- Exclude the 'updated_at' column of the "contracts" table from being overwritten
- Rewinds the date or time of all columns of the first specified table ("accounts") and the JOINed table ("contracts") that are of type DATETIME or TIMESTAMP by one month.

## Commandline Options

### past (unimplemented)

The rewind time

Acceptable: `[0-9]+(month|day|hour|minute|second)s?`

Example:

```
  --past 30minutes
```

### host (unimplemented)

Host name or IPv4 address to connect to

Acceptable: `[0-9]+\.[0-9]+\.[0-9]\.[0-9]`
Acceptable: `[a-z][a-z0-9\-.]*`

Example:

```
  --host example.com
```

### port (unimplemented)

Destination port number

Acceptable: `[0-9]+`

Example:

```
  --port 3306
```

### database (unimplemented)

Database name to connect to

Acceptable: `[a-z0-9\-]+`

Example:

```
  --database maindb
```

### table (unimplemented)

Central table name (the table specified in the FROM clause of a SELECT statement or the UPDATE clause of an UPDATE statement)

Acceptable: `<table name>`

Example:

```
  --table accounts
```

### join (unimplemented)

Table name to be joined
(Must be used with the --on option.)

Acceptable: `<table name>`

Example:

```
  --join contracts --on contracts.account_id = accounts.id
```

### on (unimplemented)

Conditions to join

Acceptable: `<SQL expression>`

Example:

```
  --join contracts --on contracts.account_id = accounts.id
```

### where (unimplemented)

Conditions to refine

Acceptable: `<SQL expression>`

Example:

```
  --where accounts.id=35
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

### dry-run (unimplemented)

Check the UPDATE statement to be executed

## Author

Keisuke Nakayama ([@yammerjp](https://github.com/yammerjp))

## LICENSE

MIT
