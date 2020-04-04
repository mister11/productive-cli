package utils

import "testing"

func TestParseTime(t *testing.T) {
	ans := ParseTime("7:30")

	if ans != 450 {
		t.Errorf("Wrong")
	}
}
