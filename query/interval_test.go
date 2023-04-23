package query

import (
	"fmt"
	"testing"
)

func TestParseInterval(t *testing.T) {
	expected := Interval{
		IsPast: false,
		Num:    3,
		Term:   "MONTH",
	}
	ret, err := ParseInterval("", "3month")
	if err != nil {
		t.Error(err)
	}
	if *ret != expected {
		fmt.Printf("expected: %v\nreturned: %v\n", expected, *ret)
		t.Error("parseInterval must be return a expected statement")
	}
}
