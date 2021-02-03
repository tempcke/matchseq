package main

import "io"

// History is a history of runes
type History map[byte]int

var emptyStruct = struct{}{}

var validBytes = map[byte]struct{}{
	'A': emptyStruct,
	'C': emptyStruct,
	'G': emptyStruct,
	'T': emptyStruct,
}

// Recorder spys on the stream to record a history of valid runes
type Recorder struct {
	history History
	reader  io.Reader
}

// NewRecorder constructs and returns a Recorder
func NewRecorder(r io.Reader) Recorder {
	return Recorder{
		reader:  r,
		history: make(History, 4),
	}
}

// Read is a spy to record the Nucleotides as they are read
func (r Recorder) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)

	for _, b := range p {
		if isNucleotide(b) {
			r.history[b]++
		}
	}

	return n, err
}

// History returns a copy of the byte history
// which reveals how many times each nucleotide was seen
func (r Recorder) History() History {
	hCopy := make(History, len(r.history))
	for k, v := range r.history {
		hCopy[k] = v
	}
	return hCopy
}

func isNucleotide(b byte) bool {
	_, ok := validBytes[b]
	return ok
}
