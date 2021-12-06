package shor

import (
	"fmt"
	"math/rand"

	"github.com/sp301415/qsim"
	"github.com/sp301415/qsim/math/fraction"
	"github.com/sp301415/qsim/math/numbers"
)

func shorInstance(N int, verbose bool) int {
	// Classical Part.
	a := 0
	for {
		a = rand.Intn(N) + 1
		K := numbers.GCD(a, N)

		if K != 1 {
			if verbose {
				fmt.Println("[!] Found factor by luck. We start again.")
			}
			return 0
		} else {
			break
		}
	}

	if verbose {
		fmt.Printf("[+] Using a: %d\n", a)
	}

	n := numbers.BitLength(N)

	if verbose {
		fmt.Println("[*] Initializing Qbit State...")
	}

	// Quantum Part.
	q := qsim.NewCircuit(3 * n)
	q.InitCbit((1 << n) - 1)

	iregs := make([]int, 2*n)
	oregs := make([]int, n)

	for i := range iregs {
		iregs[i] = i + n
	}
	for i := range oregs {
		oregs[i] = i
	}

	for _, v := range iregs {
		q.H(v)
	}

	if verbose {
		fmt.Println("[*] Applying Shor's Oracle...")
	}

	oracle := func(x int) int { return numbers.PowMod(a, x, N) }
	q.ApplyOracle(oracle, iregs, oregs)

	if verbose {
		fmt.Println("[*] Applying Inverse QFT...")
	}

	q.InvQFT(n, 3*n)

	if verbose {
		fmt.Println("[*] Measuring...")
	}

	y := q.Measure(iregs...)

	if verbose {
		fmt.Printf("[+] Measured output: %d\n", y)
	}

	// Again, Classical Part.
	Q := 1 << (2 * n)
	approxes := fraction.New(y, Q).FractionalApprox()

	r := 0
	// Try from reverse
	for i := len(approxes) - 1; i >= 0; i-- {
		r = approxes[i].D

		// r should be smaller than N.
		if r < N {
			break
		}
	}

	if verbose {
		fmt.Printf("[*] Trying with r: %d...\n", r)
	}

	factor := 0
	for v := -1; v <= 1; v += 2 {
		factor = numbers.GCD(numbers.PowMod(a, r/2, N)+v, N)

		if verbose {
			fmt.Printf("[*] Checking factor: %d...\n", factor)
		}

		if factor != 1 && factor != N && N%factor == 0 {
			return factor
		}
	}

	if verbose {
		fmt.Println("[!] Failed to find factor. :(")
	}

	return 0
}

func Shor(N int) int {
	factor := 0
	for {
		factor = shorInstance(N, false)
		if factor != 0 {
			break
		}
	}

	return factor
}

func ShorVerbose(N int) int {
	factor := 0
	for {
		factor = shorInstance(N, true)
		if factor != 0 {
			break
		}
	}

	return factor
}
