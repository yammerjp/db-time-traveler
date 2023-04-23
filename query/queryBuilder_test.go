package query

import (
	"fmt"
	"testing"
)

func TestBuildStmtToUpdateWithPrimaryKeyValue(t *testing.T) {

	expected := "UPDATE accounts SET trial_end_date = (trial_end_date - INTERVAL 1 MONTH), registered_campaign_end_datetime = (registered_campaign_end_datetime - INTERVAL 1 MONTH), created_at = (created_at - INTERVAL 1 MONTH), updated_at = (updated_at - INTERVAL 1 MONTH) WHERE id IN ( 3 )"
	ret, err := UpdateSource{
		SelectSource: SelectSource{
			Table: Table{
				targetTable: "accounts",
			},
			columns:       []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"},
			primaryKeys:   []string{"id"},
			stmtInWhereIn: "3",
		},
		Interval: Interval{
			IsPast: true,
			Num:    1,
			Term:   "MONTH",
		},
	}.buildStmtToUpdate()
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("updateQueryBuilder must be return a expected statement")
	}
}

func TestBuildStmtToUpdateWithSelectStmt(t *testing.T) {

	expected := "UPDATE accounts SET trial_end_date = (trial_end_date - INTERVAL 1 MONTH), registered_campaign_end_datetime = (registered_campaign_end_datetime - INTERVAL 1 MONTH), created_at = (created_at - INTERVAL 1 MONTH), updated_at = (updated_at - INTERVAL 1 MONTH) WHERE id IN ( SELECT id FROM ( SELECT id FROM accounts ) as any )"
	ret, err := UpdateSource{
		SelectSource: SelectSource{
			Table: Table{
				targetTable: "accounts",
			},
			columns:       []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"},
			primaryKeys:   []string{"id"},
			stmtInWhereIn: "SELECT id FROM accounts",
		},
		Interval: Interval{
			IsPast: true,
			Num:    1,
			Term:   "MONTH",
		},
	}.buildStmtToUpdate()
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("updateQueryBuilder must be return a expected statement")
	}
}

func TestBuildStmtToSelect(t *testing.T) {
	expected := "SELECT trial_end_date, registered_campaign_end_datetime, created_at, updated_at FROM accounts WHERE id IN ( 3 )"
	ret, err := SelectSource{
		Table: Table{
			targetTable: "accounts",
		},
		columns:       []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"},
		primaryKeys:   []string{"id"},
		stmtInWhereIn: "3",
	}.buildStmtToSelect()
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectTargettedColumnsQueryBuilder must be return a expected statement")
	}
}

func TestBuildStmtToSelectColumnNamesDateRelated(t *testing.T) {
	expected := "SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = \"accounts\" AND DATA_TYPE IN (\"date\", \"datetime\", \"timestamp\") AND COLUMN_NAME NOT IN (\"trial_end_date\", \"updated_at\")"
	ret, err := Table{
		targetTable: "accounts",
	}.buildStmtToSelectColumnNamesDateRelated([]string{"trial_end_date", "updated_at"})
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectDateRelatedColumnsQueryBuilder must be return a expected statement")
	}
}

func TestBuildStmtToSelectColumnNamesOfPrimaryKey(t *testing.T) {
	expected := "SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = \"accounts\" AND COLUMN_KEY = \"PRI\""
	ret, err := Table{
		targetTable: "accounts",
	}.buildStmtToSelectColumnNamesOfPrimaryKey()
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectDateRelatedColumnsQueryBuilder must be return a expected statement")
	}
}
