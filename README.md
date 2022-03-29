db-time-traveler
---

db-time-traveler rewrites the time of records in a relational database together.
You specify a table name, and db-time-traveler will rewind that record by the specified amount of time, as if it had been created in the past.

For columns of type DATE or DATETIME or TIMESTAMP in the specified table, subtract the specified time from the existing value and overwrite the record with the result.

## Requirement

To develop a function that requires operations to be performed chronologically, we would like to pretend that the records were created in the past for the purpose of debugging in the development environment.

I want to run an UPDATE statement to shift the values of DATE or DATETIME and TIMESTAMP type columns to the past by the time specified by the input, after narrowing down the records by specifying the table name.

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
  --primary-keys-where-in 35 \
```

The above command example represents the following actions

- Connect to MySQL whose IP address is 127.0.0.1 and port number is 3306
- Manipulate a database named main-db
- The records that are filtered by the primary key "35"
- Exclude the 'updated_at' column of the "contracts" table from being overwritten
- Rewinds the date or time of all columns of the first specified table ("accounts") that are of type DATE or DATETIME or TIMESTAMP by one month.

## Author

Keisuke Nakayama ([@yammerjp](https://github.com/yammerjp))

## LICENSE

MIT
