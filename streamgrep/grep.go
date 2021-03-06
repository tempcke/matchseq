package streamgrep

import (
	"bufio"
	"io"
	"log"
	"strings"
)

// Match is a single found instance of target within the stream
type Match struct {
	before string
	target string
	after  string
}

func (m Match) String() string {
	var sb strings.Builder
	if len(m.before) > 0 {
		sb.WriteString(m.before + " ")
	}
	sb.WriteString(m.target)
	if len(m.after) > 0 {
		sb.WriteString(" " + m.after)
	}
	return sb.String()
}

func newMatch(runes []rune, left, right int) Match {
	return Match{
		before: trimmedString(runes[:left]),
		target: string(runes[left:right]),
		after:  trimmedString(runes[right:]),
	}
}

// StreamGrep allows you to grep a stream of runes
type StreamGrep struct {
	target []rune
	b, a   int // before and after context length
}

// NewStreamGrep creates and returns a streamGrep
func NewStreamGrep(target string, before, after int) StreamGrep {
	return StreamGrep{[]rune(target), before, after}
}

// Grep works similar to *nix grep but on an io.Reader such as os.Stdin
// It sends matches out on a string channel
func (g StreamGrep) Grep(stream io.Reader, c chan<- Match, eos rune) {
	defer close(c)

	start, stop := g.b, g.b+len(g.target)

	in := bufio.NewReader(stream)
	w := newWindow(g.b + len(g.target) + g.a)

	checkAndHandleMatch := func() {
		s := w.runes[start:stop]
		for i := range s {
			if s[i] != g.target[i] {
				return
			}
		}
		c <- newMatch(w.runes, start, stop)
	}

	for {
		r, _, err := in.ReadRune()
		if err == io.EOF || r == eos {
			for i := 0; i < g.a; i++ {
				w.push(rune(0))
				checkAndHandleMatch()
			}
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		w.push(r)
		checkAndHandleMatch()
	}
}

type window struct {
	runes []rune
}

func newWindow(cap int) *window {
	return &window{
		runes: make([]rune, cap, cap),
	}
}

func (w *window) push(b rune) {
	w.runes = append(w.runes[1:], b)
}

func (w *window) String() string {
	return trimmedString(w.runes)
}

func trimmedString(runeList []rune) string {
	return strings.Trim(string(runeList), "\x00")
}
