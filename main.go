// Gojoin is a clone of the unix join(1) utility.
package main

import (
	"fmt"
	"os"

	"github.com/nthery/gojoin/join"
)

// Number of files expected on command-line.
const nfiles = 2

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		usage()
		os.Exit(1)
	}

	var inputs [nfiles]join.Input
	for i := 0; i < nfiles; i++ {
		f, err := os.OpenFile(args[i], os.O_RDONLY, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gojoin: cannot open %v: %v", args[i], err)
			os.Exit(1)
		}
		defer f.Close()
		inputs[i].Init(f, args[i])
	}

	err := join.Join(inputs[:], os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gojoin: error while joining files: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	out := "usage: gojoin"
	for i := 1; i <= nfiles; i++ {
		out += fmt.Sprintf(" file%d", i)
	}
	fmt.Fprintln(os.Stderr, out)
}
