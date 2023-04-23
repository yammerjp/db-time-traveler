package query

import (
	"fmt"
	"testing"
)

func TestBuildStmtToUpdateWithPrimaryKeyValue(t *testing.T) {

	expected := UpdateStatement("UPDATE accounts SET trial_end_date = (trial_end_date - INTERVAL 1 MONTH), registered_campaign_end_datetime = (registered_campaign_end_datetime - INTERVAL 1 MONTH), created_at = (created_at - INTERVAL 1 MONTH), updated_at = (updated_at - INTERVAL 1 MONTH) WHERE id IN ( 3 )")
	ret, err := UpdateSource{
		SelectSource: SelectSource{
			Table: Table{
				TargetTable: "accounts",
			},
			Columns:       []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"},
			PrimaryKeys:   []string{"id"},
			StmtInWhereIn: "3",
		},
		Interval: Interval{
			IsPast: true,
			Num:    1,
			Term:   "MONTH",
		},
	}.BuildStmtToUpdate()
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("updateQueryBuilder must be return a expected statement")
	}
}

func TestBuildStmtToSelectBeforeAndAfter(t *testing.T) {
	expected := SelectStatement("SELECT id, trial_end_date, trial_end_date - INTERVAL 1 MONTH, registered_campaign_end_datetime, registered_campaign_end_datetime - INTERVAL 1 MONTH, created_at, created_at - INTERVAL 1 MONTH, updated_at, updated_at - INTERVAL 1 MONTH FROM accounts WHERE id IN ( 3 )")
	ret, err := UpdateSource{
		SelectSource: SelectSource{
			Table: Table{
				TargetTable: "accounts",
			},
			Columns:       []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"},
			PrimaryKeys:   []string{"id"},
			StmtInWhereIn: "3",
		},
		Interval: Interval{
			IsPast: true,
			Num:    1,
			Term:   "MONTH",
		},
	}.BuildStmtToSelectBeforeAndAfter()
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectUpdatingColumnValuesBeforeAndAfterQueryBuilder must be return a expected statement")
	}
}

func TestBuildStmtToUpdateWithSelectStmt(t *testing.T) {

	expected := UpdateStatement("UPDATE accounts SET trial_end_date = (trial_end_date - INTERVAL 1 MONTH), registered_campaign_end_datetime = (registered_campaign_end_datetime - INTERVAL 1 MONTH), created_at = (created_at - INTERVAL 1 MONTH), updated_at = (updated_at - INTERVAL 1 MONTH) WHERE id IN ( SELECT id FROM ( SELECT id FROM accounts ) as any )")
	ret, err := UpdateSource{
		SelectSource: SelectSource{
			Table: Table{
				TargetTable: "accounts",
			},
			Columns:       []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"},
			PrimaryKeys:   []string{"id"},
			StmtInWhereIn: "SELECT id FROM accounts",
		},
		Interval: Interval{
			IsPast: true,
			Num:    1,
			Term:   "MONTH",
		},
	}.BuildStmtToUpdate()
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("updateQueryBuilder must be return a expected statement")
	}
}

func TestBuildStmtToSelectColumnNamesDateRelated(t *testing.T) {
	expected := SelectStatement("SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = \"accounts\" AND DATA_TYPE IN (\"date\", \"datetime\", \"timestamp\") AND COLUMN_NAME NOT IN (\"trial_end_date\", \"updated_at\")")
	ret, err := Table{
		TargetTable: "accounts",
	}.BuildStmtToSelectColumnNamesDateRelated([]string{"trial_end_date", "updated_at"})
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectDateRelatedColumnsQueryBuilder must be return a expected statement")
	}
}

func TestBuildStmtToSelectColumnNamesOfPrimaryKey(t *testing.T) {
	expected := SelectStatement("SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = \"accounts\" AND COLUMN_KEY = \"PRI\"")
	ret, err := Table{
		TargetTable: "accounts",
	}.BuildStmtToSelectColumnNamesOfPrimaryKey()
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectDateRelatedColumnsQueryBuilder must be return a expected statement")
	}
}
