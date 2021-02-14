// Gojoin is a clone of the unix join(1) utility.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Number of files expected on command-line.
const nfiles = 2

// input keeps track of the state of a single input file.
type input struct {
	rd *bufio.Scanner

	// The file we are reading from.  Used for error reporting.
	path string

	// Words on the current line
	words []string
}

func (in *input) init(r io.Reader, path string) {
	in.rd = bufio.NewScanner(r)
	in.path = path
}

// key returns the join key for the current line.
// At least one line must have been read.
func (in *input) key() string {
	return in.words[0]
}

// hasWords returns true if a line has been read and false on EOF or error.
func (in *input) hasWords() bool {
	return len(in.words) > 0
}

// err returns the first error if any detected on this input.
func (in *input) err() error {
	return in.rd.Err()
}

// read reads up to the the next non empty input line and splits it in words.
func (in *input) read() {
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

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		usage()
		os.Exit(1)
	}

	var inputs [nfiles]input
	for i := 0; i < nfiles; i++ {
		f, err := os.OpenFile(args[i], os.O_RDONLY, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gojoin: cannot open %v: %v", args[i], err)
			os.Exit(1)
		}
		defer f.Close()
		inputs[i].init(f, args[i])
	}

	err := join(inputs[:], os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gojoin: error while joining files: %v\n", err)
		os.Exit(1)
	}
}

// join iterates over all input lines and generates lines containing all words
// for lines with matching key.
func join(inputs []input, output io.Writer) error {
	out := bufio.NewWriter(output)

	for i := 0; i < len(inputs); i++ {
		inputs[i].read()
	}

	ks := newKeySorter(inputs)

	for allInputsHaveWords(inputs) {
		if !allInputsHaveSameKey(inputs) {
			highest := ks.highestInput()
			for i := 0; i < len(inputs); i++ {
				if i != highest {
					inputs[i].read()
				}
			}
			continue
		}

		fmt.Fprintf(out, "%s", inputs[0].words[0])
		for i := 0; i < len(inputs); i++ {
			ws := strings.Join(inputs[i].words[1:], " ")
			if len(ws) > 0 {
				fmt.Fprintf(out, " %s", ws)
			}
		}
		fmt.Fprintln(out)

		for i := 0; i < len(inputs); i++ {
			inputs[i].read()
		}
	}

	for i := 0; i < len(inputs); i++ {
		if inputs[i].err() != nil {
			return fmt.Errorf("read error in %s: %v", inputs[i].path, inputs[i].err())
		}
	}

	return out.Flush()
}

func allInputsHaveWords(inputs []input) bool {
	for i := 0; i < len(inputs); i++ {
		if !inputs[i].hasWords() {
			return false
		}
	}
	return true
}

func allInputsHaveSameKey(inputs []input) bool {
	for i := 1; i < len(inputs); i++ {
		if inputs[0].key() != inputs[i].key() {
			return false
		}
	}
	return true
}

// keySorter stores state needed to sort keys for current lines in descending order.
type keySorter struct {
	inputs  []input
	indices []int
}

func newKeySorter(inputs []input) keySorter {
	return keySorter{
		inputs:  inputs,
		indices: make([]int, len(inputs), len(inputs)),
	}
}

func (ks *keySorter) Len() int {
	return len(ks.indices)
}

func (ks *keySorter) Less(i, j int) bool {
	// descending order
	return ks.inputs[ks.indices[i]].key() > ks.inputs[ks.indices[j]].key()
}

func (ks *keySorter) Swap(i, j int) {
	ks.indices[i], ks.indices[j] = ks.indices[j], ks.indices[i]
}

// highestInput returns the index of the input with the highest key.
func (ks *keySorter) highestInput() int {
	for i := 0; i < len(ks.indices); i++ {
		ks.indices[i] = i
	}
	sort.Sort(ks)
	return ks.indices[0]
}

func usage() {
	out := "usage: gojoin"
	for i := 1; i <= nfiles; i++ {
		out += fmt.Sprintf(" file%d", i)
	}
	fmt.Fprintln(os.Stderr, out)
}
