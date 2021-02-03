[![Go Report Card](https://goreportcard.com/badge/github.com/tempcke/matchseq)](https://goreportcard.com/report/github.com/tempcke/matchseq)

# Phase I - Output Nucleotide Sequences with Context

Consider nucleotide sequences that contain the possible characters: 'A', 'C', 'G', 'T', and the end-of-sequence character: 'ε'.

**Implement a command-line program that accepts a nucleotide stream on `stdin`, scans the stream for any sequence matching a specified target nucleotide, and writes each match to `stdout` along with its surrounding context. The program should end when encountering the end-of-sequence character or the end of the stream.**

More specifically: given a nucleotide stream *S*, a target sequence *T*, and two integers *x* and *y*; for each *T* in *S*: print the *x* preceding nucleotides, the sequence matching *T*, and the *y* succeeding nucleotides.

An example command line

    echo "ACACGTCAε" | matchseq -T:ACGT -x:1 -y:2

would yield

    C ACGT CA

where C is the *x*=1 preceding nucleotide, ACGT is the nucleotide sequence matching *T*, and CA is the *y*=2 succeeding nucleotides.

Be aware that:

- *x* may be unspecified, which indicates that no preceding nucleotides should be printed. Likewise, *y* may be unspecified, which indicates that no succeeding nucleotides should be printed. However, a non-empty value for *T* is required.
- If the stream contains fewer than *x* nucleotides before *T*, or fewer than *y* nucleotides after *T*, print as many as there actually are.
- The end-of-sequence value 'ε' will not appear in *T*.
- Targets may overlap in the sequence *S*, and each should be treated as a distinct occurrence with its own surrounding context.
- Streams are potentially unlimited, so be sure to consider the case where the size of *S* exceeds a system's memory.

## Example

    echo "AAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTAε" | matchseq -T:AGTA -x:5 -y:7

should print the following lines:

    A AGTA CGTGCAG
    CAGTG AGTA GTAGACC
    TGAGT AGTA GACCTGA
    ATATA AGTA GCTA

Notice that the 2nd and 3rd lines display overlapping targets and that the 1st and 4th show fewer than *x* and *y* elements of context, respectively.


# Phase II - Histagram

In addition to the output from Phase I include a histagram showing how many times each nucleotide was seen.

## Example

    echo "AAGTACGTGCAGTGAGTAGTAGACCTGACGTAGACCGATATAAGTAGCTAε" | matchseq -T:AGTA -x:5 -y:7

Should also communicate that `A` was seen 17 times, `C` 8 times, `G` 14 times, and `T` 11 times

# Execution instructions

- `make build` will create the bin/matchseq binary
- `make test` will execute all tests
- `make example` will execute the final example to demonstrate it is working.