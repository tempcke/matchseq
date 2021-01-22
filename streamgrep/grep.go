package streamgrep

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strings"
)

type byteWindow struct {
	bytes []byte
}

func newByteWindow(cap int) *byteWindow {
	return &byteWindow{
		bytes: make([]byte, cap, cap),
	}
}

func (w *byteWindow) push(b byte) {
	w.bytes = append(w.bytes[1:], b)
}

func (w *byteWindow) String() string {
	return trimmedString(w.bytes)
}

func trimmedString(byteList []byte) string {
	return string(bytes.Trim(byteList, "\x00"))
}

// StreamGrep allows you to grep a stream of bytes
type StreamGrep struct {
	target string
	b, a   int // before and after context length
}

// NewStreamGrep creates and returns a streamGrep
func NewStreamGrep(target string, before, after int) StreamGrep {
	return StreamGrep{target, before, after}
}

// Grep works similar to *nix grep but on an io.Reader such as os.Stdin
// It sends matches out on a string channel
func (g StreamGrep) Grep(stream io.Reader, c chan<- string, eos rune) {
	defer close(c)

	start, stop := g.b, g.b+len(g.target)

	in := bufio.NewReader(stream)
	w := newByteWindow(g.b + len(g.target) + g.a)
	for {
		r, _, err := in.ReadRune()
		if err == io.EOF || r == eos {
			for i := 0; i < g.a; i++ {
				w.push(byte(0))
				if string(w.bytes[start:stop]) == g.target {
					c <- format(w.bytes, start, stop)
				}
			}
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		w.push(byte(r))

		if string(w.bytes[start:stop]) == g.target {
			c <- format(w.bytes, start, stop)
		}
	}
}

func format(bytes []byte, left, right int) string {
	var sb strings.Builder

	prefix := trimmedString(bytes[:left])
	if len(prefix) > 0 {
		sb.WriteString(prefix + " ")
	}

	target := string(bytes[left:right])
	sb.WriteString(target)

	suffix := trimmedString(bytes[right:])
	if len(suffix) > 0 {
		sb.WriteString(" " + suffix)
	}

	return sb.String()
}
