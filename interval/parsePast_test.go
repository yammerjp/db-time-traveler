package interval

import (
	"fmt"
	"testing"

	"github.com/yammerjp/db-time-traveler/system"
)

func TestParsePast(t *testing.T) {
	expected := "1 MONTH"
	ret, err := ParsePast("1month")
	if err != nil {
		t.Error(err)
	}
	if ret != expected {
		fmt.Printf("expected: %s\nreturned: %s\n", expected, ret)
		t.Error("parsePast must be return a expected statement")
	}
}

func TestParseInterval(t *testing.T) {
	expected := system.Interval{
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
