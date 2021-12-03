package main

import (
	"flag"
	"fmt"

	"github.com/sp301415/qsim/algorithms/shor"
)

func main() {
	lenPtr := flag.Int("n", 0, "Number to factorize.")
	verbPtr := flag.Bool("verbose", false, "Prints messages when on.")

	flag.Parse()

	if *lenPtr == 0 {
		panic("Invalid argument")
	}

	n := *lenPtr
	verb := *verbPtr

	if verb {
		fmt.Printf("[+] Found factor of %d: %d\n", n, shor.ShorVerbose(n))
	} else {
		fmt.Printf("[+] Found factor of %d: %d\n", n, shor.Shor(n))
	}
}
