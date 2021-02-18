package join

import (
	"bufio"
	"io"
	"strings"
)

// Input keeps track of the state of an input source.
// Lines must be sorted in ascending key order.
type Input struct {
	rd *bufio.Scanner

	// The file we are reading from.  Used for error reporting.
	path string

	// Words on the current line
	words []string
}

// Init initializes this input.
func (in *Input) Init(r io.Reader, path string) {
	in.rd = bufio.NewScanner(r)
	in.path = path
}

// key returns the join key for the current line.
// At least one line must have been read.
func (in *Input) key() string {
	return in.words[0]
}

// hasWords returns true if a line has been read and false on EOF or error.
func (in *Input) hasWords() bool {
	return len(in.words) > 0
}

// err returns the first error if any detected on this input.
func (in *Input) err() error {
	return in.rd.Err()
}

// read reads up to the the next non empty input line and splits it in words.
func (in *Input) read() {
	for {
		if !in.rd.Scan() {
			in.words = in.words[0:0]
			return
		}

		line := strings.Trim(in.rd.Text(), " \t")
		if len(line) == 0 {
			continue
		}

		in.words = strings.Split(line, " ")
		return
	}
}
