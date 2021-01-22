package main

import (
	"fmt"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tt := []struct {
		args []string
		t    string
		x, y int
	}{
		{[]string{"-T:ACGT"}, "ACGT", 0, 0},
		{[]string{"-T:ACGT", "-x:1", "-y:2"}, "ACGT", 1, 2},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprint(tc.args), func(t *testing.T) {
			c := newConfFromArgs(tc.args)
			assertEqual(t, nil, c.validate())
			assertEqual(t, tc.t, c.target)
			assertEqual(t, tc.x, c.x)
			assertEqual(t, tc.y, c.y)
		})
	}
}

func TestConfValidation(t *testing.T) {
	// so long as c.t != "" it should be valid
	var c conf
	c = newConfFromArgs([]string{"-T:A"})
	if c.validate() != nil {
		t.Error("when target is set it should be valid")
	}
	c = newConfFromArgs([]string{"-x:1"})
	if c.validate() == nil {
		t.Error("when target is not set it should not be valid")
	}
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Fatalf(
			"Not Equal\nWant: %v\t%T\nGot:  %v\t%T",
			expected, expected,
			actual, actual)
	}
}
