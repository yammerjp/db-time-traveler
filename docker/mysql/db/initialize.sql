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
