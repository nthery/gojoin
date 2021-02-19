// Gojoin is a clone of the unix join(1) utility.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nthery/gojoin/join"
)

var sep = flag.String("t", " ", "field separator")

func main() {
	flag.Parse()

	if len(*sep) == 0 {
		fmt.Fprintln(os.Stderr, "gojoin: separator must not be empty")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) < 2 {
		usage()
		os.Exit(1)
	}

	inputs := make([]join.Input, 0, len(args))
	for _, path := range args {
		f, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gojoin: cannot open %v: %v\n", path, err)
			os.Exit(1)
		}
		defer f.Close()
		inputs = append(inputs, join.NewInput(f, path))
	}

	err := join.Join(inputs, *sep, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gojoin: error while joining files: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: gojoin [-t sep] file...")
}
