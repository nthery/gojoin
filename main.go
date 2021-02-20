// Gojoin is a clone of the unix join(1) utility.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nthery/gojoin/join"
)

func main() {
	log.SetPrefix("gojoin: ")
	log.SetFlags(0)

	sep := flag.String("t", " ", "field separator")
	flag.Parse()

	if len(*sep) == 0 {
		log.Fatal("separator must not be empty")
	}

	args := flag.Args()
	if len(args) < 2 {
		usage()
	}

	inputs := make([]join.Input, 0, len(args))
	for _, path := range args {
		f, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		defer f.Close()
		inputs = append(inputs, join.NewInput(f, path))
	}

	err := join.Join(inputs, *sep, os.Stdout)
	if err != nil {
		log.Fatalf("error while joining files: %v", err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: gojoin [-t sep] file...")
	os.Exit(1)
}
