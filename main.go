// Gojoin is a clone of the unix join(1) utility.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nthery/gojoin/join"
)

// Number of files expected on command-line.
const nfiles = 2

var sep = flag.String("t", " ", "field separator")

func main() {
	flag.Parse()

	if len(*sep) == 0 {
		fmt.Fprintln(os.Stderr, "gojoin: separator must not be empty")
		os.Exit(1)
	}

	args := flag.Args()
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
		inputs[i] = join.NewInput(f, args[i])
	}

	err := join.Join(inputs[:], *sep, os.Stdout)
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
