package streamgrep

import (
	"fmt"
	"strings"
	"testing"
)

func TestWindow(t *testing.T) {
	w := newWindow(3)
	var (
		a = 'A'
		b = 'B'
		c = '了'
		d = 'D'
	)
	w.push(a)
	w.push(b)
	assertEqual(t, a, w.runes[1])
	assertEqual(t, b, w.runes[2])
	assertEqual(t, "AB", w.String()) // ensure first zero bit is trimmed
	assertEqual(t, "AB", string(w.runes[1:]))
	assertEqual(t, 3, len(w.runes))

	w.push(c)
	w.push(d)
	assertEqual(t, "B了D", w.String())
}

func TestStreamGrep(t *testing.T) {
	tt := []struct {
		input    string
		target   string
		b, a     int // before and after context length
		expected []string
	}{
		{"", "A", 0, 0, []string{}},
		{"A", "A", 0, 0, []string{"A"}},
		{"ABC", "A", 0, 0, []string{"A"}},
		{"ABC", "B", 1, 1, []string{"A B C"}},
		{"ABC", "A", 1, 1, []string{"A B"}},
		{"ABC", "C", 1, 1, []string{"B C"}},
		{"ABCD", "BC", 1, 1, []string{"A BC D"}},
		// shorter than expected prefix
		{"ABC", "BC", 2, 1, []string{"A BC"}},
		{"ABCDABD", "BC", 3, 3, []string{"A BC DAB"}},

		{"ABCBCDE", "BC", 3, 3, []string{"A BC BCD", "ABC BC DE"}},
		{"ABCDABCD", "BC", 3, 3, []string{"A BC DAB", "CDA BC D"}},
		{"ABCDABCDεABCDABCD", "BC", 3, 3, []string{"A BC DAB", "CDA BC D"}},
		{
			"AAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTAε",
			"AGTA", 5, 7,
			[]string{
				"A AGTA CGTGCAG",
				"CAGTG AGTA GTAGACC",
				"TGAGT AGTA GACCTGA",
				"ATATA AGTA GCTA",
			},
		},
		// rune support?
		{"的是了见", "是了", 1, 1, []string{"的 是了 见"}},
	}

	eos := 'ε' // end of stream rune
	for _, tc := range tt {
		testName := fmt.Sprintf("%v: %v", tc.input, tc.expected)
		t.Run(testName, func(t *testing.T) {
			c := make(chan string)
			g := NewStreamGrep(tc.target, tc.b, tc.a)
			go g.Grep(strings.NewReader(tc.input), c, eos)
			i := 0
			for s := range c { // should block until chan is closed
				assertEqual(t, tc.expected[i], s)
				i++
			}
			assertEqual(t, len(tc.expected), i)
		})
	}
}

// assertions

func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Fatalf(
			"Not Equal\nWant: %v\t%T\nGot:  %v\t%T",
			expected, expected,
			actual, actual)
	}
}
