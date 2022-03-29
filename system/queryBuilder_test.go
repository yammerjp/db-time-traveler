package system

import (
	"fmt"
	"testing"
)

func TestBuildStmtToUpdate(t *testing.T) {

	expected := "UPDATE accounts SET trial_end_date = (trial_end_date - INTERVAL 1 MONTH), registered_campaign_end_datetime = (registered_campaign_end_datetime - INTERVAL 1 MONTH), created_at = (created_at - INTERVAL 1 MONTH), updated_at = (updated_at - INTERVAL 1 MONTH) WHERE (id) IN ( 3 )"
	ret, err := QueryBuilderSourceToUpdate{
		QueryBuilderSourceForColumnValues: QueryBuilderSourceForColumnValues{
			QueryBuilderSourceForSchemaInformation: QueryBuilderSourceForSchemaInformation{
				targetTable: "accounts",
			},
			columns:       []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"},
			primaryKeys:   []string{"id"},
			stmtInWhereIn: "3",
		},
		QueryBuilderSourcePartOfInterval: QueryBuilderSourcePartOfInterval{
			isPast: true,
			num:    1,
			term:   "MONTH",
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
	expected := "SELECT trial_end_date, registered_campaign_end_datetime, created_at, updated_at FROM accounts WHERE (id) IN ( 3 )"
	ret, err := QueryBuilderSourceForColumnValues{
		QueryBuilderSourceForSchemaInformation: QueryBuilderSourceForSchemaInformation{
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

func TestBuildStmtToSelectUpdatingColumnValues(t *testing.T) {
	expected := "SELECT trial_end_date - INTERVAL 1 MONTH, registered_campaign_end_datetime - INTERVAL 1 MONTH, created_at - INTERVAL 1 MONTH, updated_at - INTERVAL 1 MONTH FROM accounts WHERE (id) IN ( 3 )"
	ret, err := QueryBuilderSourceToUpdate{
		QueryBuilderSourceForColumnValues: QueryBuilderSourceForColumnValues{
			QueryBuilderSourceForSchemaInformation: QueryBuilderSourceForSchemaInformation{
				targetTable: "accounts",
			},
			columns:       []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"},
			primaryKeys:   []string{"id"},
			stmtInWhereIn: "3",
		},
		QueryBuilderSourcePartOfInterval: QueryBuilderSourcePartOfInterval{
			isPast: true,
			num:    1,
			term:   "MONTH",
		},
	}.buildStmtToSelect()
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectUpdatingColumnValuesQueryBuilder must be return a expected statement")
	}
}

func TestBuildStmtToSelectColumnNamesDateRelated(t *testing.T) {
	expected := "SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = \"accounts\" AND DATA_TYPE IN (\"date\", \"datetime\", \"timestamp\")"
	ret, err := QueryBuilderSourceForSchemaInformation{
		targetTable: "accounts",
	}.buildStmtToSelectColumnNamesDateRelated()
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
	ret, err := QueryBuilderSourceForSchemaInformation{
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

func TestBuildStmtToSelectBeforeAndAfter(t *testing.T) {
	expected := "SELECT id, trial_end_date, trial_end_date - INTERVAL 1 MONTH, registered_campaign_end_datetime, registered_campaign_end_datetime - INTERVAL 1 MONTH, created_at, created_at - INTERVAL 1 MONTH, updated_at, updated_at - INTERVAL 1 MONTH FROM accounts WHERE (id) IN ( 3 )"
	ret, err := QueryBuilderSourceToUpdate{
		QueryBuilderSourceForColumnValues: QueryBuilderSourceForColumnValues{
			QueryBuilderSourceForSchemaInformation: QueryBuilderSourceForSchemaInformation{
				targetTable: "accounts",
			},
			columns:       []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"},
			primaryKeys:   []string{"id"},
			stmtInWhereIn: "3",
		},
		QueryBuilderSourcePartOfInterval: QueryBuilderSourcePartOfInterval{
			isPast: true,
			num:    1,
			term:   "MONTH",
		},
	}.buildStmtToSelectBeforeAndAfter()
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectUpdatingColumnValuesBeforeAndAfterQueryBuilder must be return a expected statement")
	}
}
