package main

import (
	"flag"
	"fmt"

	"github.com/sp301415/qsim/algorithms/deutchjozsa"
)

func main() {
	lenPtr := flag.Int("n", 0, "Number of qbits.")
	typePtr := flag.String("type", "", "Type of function. b: balanced, c: constant.")

	flag.Parse()

	if *typePtr != "b" && *typePtr != "c" {
		panic("Invalid argument")
	}

	if *lenPtr == 0 {
		panic("Invalid argument")
	}

	n := *lenPtr
	is_constant := *typePtr == "c"

	if is_constant {
		if deutchjozsa.DeutchJozsa(n, deutchjozsa.ConstantFunc) {
			fmt.Println("Deutch Jozsa Says: This function is CONSTANT!")
		} else {
			fmt.Println("Seems like Deutch Jozsa is wrong :(")
		}
	} else {
		if !deutchjozsa.DeutchJozsa(n, deutchjozsa.BalancedFunc) {
			fmt.Println("Deutch Jozsa Says: This function is BALANCED!")
		} else {
			fmt.Println("Seems like Deutch Jozsa is wrong :(")
		}
	}
}
