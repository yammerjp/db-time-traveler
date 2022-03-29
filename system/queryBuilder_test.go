package system

import (
	"fmt"
	"testing"
)

func TestUpdateQueryBuilder(t *testing.T) {

	expected := "UPDATE accounts SET trial_end_date = (trial_end_date - INTERVAL 1 MONTH), registered_campaign_end_datetime = (registered_campaign_end_datetime - INTERVAL 1 MONTH), created_at = (created_at - INTERVAL 1 MONTH), updated_at = (updated_at - INTERVAL 1 MONTH) WHERE (id) IN ( 3 )"
	ret, err := updateQueryBuilder("accounts", []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"}, "1 MONTH", []string{"id"}, "3")
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("updateQueryBuilder must be return a expected statement")
	}
}

func TestSelectTargettedColumnsQueryBuilder(t *testing.T) {
	expected := "SELECT trial_end_date, registered_campaign_end_datetime, created_at, updated_at FROM accounts WHERE (id) IN ( 3 )"
	ret, err := selectTargettedColumnsQueryBuilder("accounts", []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"}, []string{"id"}, "3")
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectTargettedColumnsQueryBuilder must be return a expected statement")
	}
}

func TestSelectUpdatingColumnValuesQueryBuilder(t *testing.T) {
	expected := "SELECT trial_end_date - INTERVAL 1 MONTH, registered_campaign_end_datetime - INTERVAL 1 MONTH, created_at - INTERVAL 1 MONTH, updated_at - INTERVAL 1 MONTH FROM accounts WHERE (id) IN ( 3 )"
	ret, err := selectUpdatingColumnValuesQueryBuilder("accounts", []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"}, "1 MONTH", []string{"id"}, "3")
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectUpdatingColumnValuesQueryBuilder must be return a expected statement")
	}
}

func TestSelectDateRelatedColumnsQueryBuilder(t *testing.T) {
	expected := "SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = \"accounts\" AND DATA_TYPE IN (\"date\", \"datetime\", \"timestamp\")"
	ret, err := selectDateRelatedColumnsQueryBuilder("accounts")
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectDateRelatedColumnsQueryBuilder must be return a expected statement")
	}
}

func TestSelectPrimaryKeyColumnsQueryBuilder(t *testing.T) {
	expected := "SELECT DISTINCT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = \"accounts\" AND COLUMN_KEY = \"PRI\""
	ret, err := selectPrimaryKeyColumnsQueryBuilder("accounts")
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectDateRelatedColumnsQueryBuilder must be return a expected statement")
	}
}

func TestSelectUpdatingColumnValuesBeforeAdnAfterQueryBUilder(t *testing.T) {
	expected := "SELECT id, trial_end_date, trial_end_date - INTERVAL 1 MONTH, registered_campaign_end_datetime, registered_campaign_end_datetime - INTERVAL 1 MONTH, created_at, created_at - INTERVAL 1 MONTH, updated_at, updated_at - INTERVAL 1 MONTH FROM accounts WHERE (id) IN ( 3 )"
	ret, err := selectUpdatingColumnValuesBeforeAndAfterQueryBuilder("accounts", []string{"trial_end_date", "registered_campaign_end_datetime", "created_at", "updated_at"}, "1 MONTH", []string{"id"}, "3")
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("selectUpdatingColumnValuesBeforeAndAfterQueryBuilder must be return a expected statement")
	}
}
