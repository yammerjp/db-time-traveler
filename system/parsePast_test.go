package system

import (
	"fmt"
	"testing"
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
