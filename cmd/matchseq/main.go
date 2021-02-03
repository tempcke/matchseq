package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tempcke/matchseq/streamgrep"
)

const eos = 'ε' // end of stream rune

var help = []string{
	"help",
	"--help",
	"-h",
}

type conf struct {
	target string
	x, y   int
}

func (c conf) validate() error {
	if c.target == "" {
		return errors.New("Target is required, set with -T:value")
	}
	return nil
}

func newConfFromArgs(args []string) conf {
	c := conf{}
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "-T:"):
			c.target = strings.TrimPrefix(arg, "-T:")
		case strings.HasPrefix(arg, "-x:"):
			x, err := strconv.Atoi(strings.TrimPrefix(arg, "-x:"))
			if err == nil {
				c.x = x
			}
		case strings.HasPrefix(arg, "-y:"):
			y, err := strconv.Atoi(strings.TrimPrefix(arg, "-y:"))
			if err == nil {
				c.y = y
			}
		}
	}
	return c
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 || isHelp(args[0]) {
		outputUsage()
		os.Exit(0)
	}

	cfg := newConfFromArgs(os.Args[1:])
	if err := cfg.validate(); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}

	if err := run(os.Stdin, cfg); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}

func isHelp(h string) bool {
	for _, s := range help {
		if h == s {
			return true
		}
	}
	return false
}

func run(stream io.Reader, cfg conf) error {
	c := make(chan streamgrep.Match)
	recorder := NewRecorder(stream)
	g := streamgrep.NewStreamGrep(cfg.target, cfg.x, cfg.y)
	go g.Grep(recorder, c, eos)
	fmt.Print("\nFinding Matches....\n")
	for m := range c { // should block until chan is closed
		fmt.Println(m.String())
	}
	fmt.Println("\nNucleotide Histogram:\n Type\tCount")
	for b, n := range recorder.History() {
		fmt.Printf("   %s\t%4v\n", string(b), n)
	}
	return nil
}

func outputUsage() {
	fmt.Print(`
matchseq Usage
	pipe a string into: matchseq -T:TARGET [options]

Required Options
	-T:TARGET
		TARGET is the substring we are searching for in the stream

Other Options
	-x:NUM
		Print NUM bytes of leading context before matching target

	-y:NUM
		Print NUM bytes of trailing context after matching target

	-h, --help, help
		Print this usage information

Example 1
	echo "ACACGTCAε" | matchseq -T:ACGT -x:1 -y:2

	Will Output:
		C ACGT CA

Example 2
	echo "AAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTAε" | matchseq -T:AGTA -x:5 -y:7

	Will Output:
		A AGTA CGTGCAG
		CAGTG AGTA GTAGACC
		TGAGT AGTA GACCTGA
		ATATA AGTA GCTA
` + "\n")
}
