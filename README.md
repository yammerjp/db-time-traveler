db-time-traveler
---

db-time-traveler rewrites the time of records in a relational database together.
You specify a table name, and db-time-traveler will rewind that record by the specified amount of time, as if it had been created in the past.

For columns of type DATE or DATETIME or TIMESTAMP in the specified table, subtract the specified time from the existing value and overwrite the record with the result.

Example: db-time-traveler build and execute such as the following SQL from command line options

```sql
UPDATE
  accounts
SET
  trial_end_date = (trial_end_date - INTERVAL 1 MONTH),
  registered_campaign_end_datetime = (registered_campaign_end_datetime - INTERVAL 1 MONTH),
  created_at = (created_at - INTERVAL 1 MONTH),
  updated_at = (updated_at - INTERVAL 1 MONTH)
WHERE id IN 1
```

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
  --ignore updated_at
```

The above command example represents the following actions

- Connect to MySQL whose IP address is 127.0.0.1 and port number is 3306
- Manipulate a database named main-db
- The records that are filtered by the primary key "35"
- Exclude the 'updated_at' column of the "contracts" table from being overwritten
- Rewinds the date or time of all columns of the first specified table ("accounts") that are of type DATE or DATETIME or TIMESTAMP by one month.
- Ignore "updated_at" from updating columns

## Setup

```bash
$ go install github.com/yammerjp/db-time-traveler@latest  # download repository and build and install
$ db-time-traveler config-init                            # generate config yaml
$ vim ~/.db-time-traveler.yaml                            # edit config yaml
```

## Example

There is a executed command example of following the table definition

a executed command example

```bash
$ db-time-traveler \
    update-dry-run\
    --future 1month \
    --primary-keys-where-in="1" \
    --print-query \
    --table accounts \
    --ignore updated_at \
    --ignore trial_end_date
id: 1
  registered_campaign_end_datetime:
    before: 2021-04-03T00:10:20Z
    after:  2021-05-03T00:10:20Z
  created_at:
    before: 2022-03-30T15:59:55Z
    after:  2022-03-30T15:59:55Z
UPDATE accounts SET registered_campaign_end_datetime = (registered_campaign_end_datetime + INTERVAL 1 MONTH), created_at = (created_at + INTERVAL 1 MONTH) WHERE (id) IN ( 1 )
```

the table definition

```sql
DROP SCHEMA IF EXISTS sampleschema;
CREATE SCHEMA sampleschema;
USE sampleschema;

DROP TABLE IF EXISTS accounts;

CREATE TABLE accounts
(
  id                            INT(10) auto_increment,
  name                          VARCHAR(40),
  trial_end_date                DATE,
  registered_campaign_end_datetime  DATETIME,
  created_at                    TIMESTAMP default current_timestamp,
  updated_at                    TIMESTAMP default current_timestamp on update current_timestamp,
  PRIMARY KEY (id)
);

INSERT INTO accounts (id, name, trial_end_date, registered_campaign_end_datetime) VALUES (1, "Nagaoka", "2022-05-30", "2022-04-03 00:10:20");
INSERT INTO accounts (id, name) VALUES (2, "Tanaka");
```

## Config

example

```
default_connection: local
connections:
  -
    name: local
    driver: mysql
    hostname: 127.0.0.1
    username: root
    password: password
    port: 3306
    database: sampleschema
  -
    name: sshconnection
    driver: mysql
    hostname: localhost
    username: hogehogeuser
    password: password
    port: 3306
    database: sampleschema
    sshhost: bastion.example.com
    sshport: 22
    sshuser: yammer
    sshkeypath: /home/username/.ssh/id_rsa
    sshpassphrase: helloworld
```

- `default_connection` define default selected connection by name
- write ssh configuration together, If you want to connect with proxy of SSH

## Author

Keisuke Nakayama ([@yammerjp](https://github.com/yammerjp))

## LICENSE

MIT
