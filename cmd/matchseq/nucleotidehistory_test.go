package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// A, C, G, T

func TestHistory(t *testing.T) {
	tt := []struct {
		input    string
		expected History
	}{
		{"", History{}},
		{"A", History{'A': 1}},
		{"AA", History{'A': 2}},
		{"ACGT", History{'A': 1, 'C': 1, 'G': 1, 'T': 1}},
		{"AB", History{'A': 1}}, // B should be ignored because it is not A,C,G, or T
	}

	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			stream := strings.NewReader(tc.input)
			recorder := NewRecorder(stream)

			// read all data through the recorder
			bytes, err := ioutil.ReadAll(recorder)
			assert.NoError(t, err)
			assert.Equal(t, tc.input, string(bytes))

			recordedHistory := recorder.History()
			assert.Equal(t, tc.expected, recordedHistory)
		})
	}
}
