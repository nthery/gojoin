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

// NewInput constructs an input for the specified source.
// path is the name of the file being read.  It is used for error reporting.
func NewInput(r io.Reader, path string) Input {
	return Input{rd: bufio.NewScanner(r), path: path}
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

// read reads up to the the next non empty input line and splits it in words
// separarated by sep.
func (in *Input) read(sep string) {
	for {
		if !in.rd.Scan() {
			in.words = in.words[0:0]
			return
		}

		line := strings.Trim(in.rd.Text(), " \t")
		if len(line) == 0 {
			continue
		}

		in.words = strings.Split(line, sep)
		return
	}
}
