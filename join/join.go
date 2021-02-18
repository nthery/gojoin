// Package join implements the join algorithm described in the join(1) unix command manual
// generalized to handle more than two input sources.
package join

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Join iterates over all input lines and generates lines containing all words
// for lines with matching key.
func Join(inputs []Input, output io.Writer) error {
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

func allInputsHaveWords(inputs []Input) bool {
	for i := 0; i < len(inputs); i++ {
		if !inputs[i].hasWords() {
			return false
		}
	}
	return true
}

func allInputsHaveSameKey(inputs []Input) bool {
	for i := 1; i < len(inputs); i++ {
		if inputs[0].key() != inputs[i].key() {
			return false
		}
	}
	return true
}

// keySorter stores state needed to sort keys for current lines in descending order.
type keySorter struct {
	inputs  []Input
	indices []int
}

func newKeySorter(inputs []Input) keySorter {
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
