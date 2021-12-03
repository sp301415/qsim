package shor

import (
	"fmt"
	"math/rand"

	"github.com/sp301415/qsim"
	"github.com/sp301415/qsim/math/fraction"
	"github.com/sp301415/qsim/math/numbers"
	"github.com/sp301415/qsim/math/vector"
	"github.com/sp301415/qsim/quantum/qbit"
)

func shorInstance(N int, verbose bool) int {
	// Classical Part.
	a := 0
	for {
		a = rand.Intn(N) + 1
		K := numbers.GCD(a, N)

		if K != 1 {
			if verbose {
				fmt.Println("[-] Found factor by luck. We start again.")
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
	q.InitQbit(qbit.NewFromCbit((1<<n)-1, 3*n))

	for i := n; i < 3*n; i++ {
		q.H(i)
	}

	newState := vector.Zeros(q.State.Dim())

	if verbose {
		fmt.Println("[*] Applying Shor's Oracle...")
	}

	for qn, qa := range q.State {
		if qa == 0 {
			continue
		}
		x := qn >> n
		r := numbers.PowMod(a, x, N)
		o := (qn % (1 << n)) ^ r
		newState[(x<<n)+o] += qa
	}
	q.State = newState

	if verbose {
		fmt.Println("[*] Applying Inverse QFT...")
	}

	q.InvQFT(n, 3*n)

	if verbose {
		fmt.Println("[*] Measuring...")
	}

	mslice := make([]int, 2*n)
	for i := range mslice {
		mslice[i] = i + n
	}
	y := q.Measure(mslice...)

	if verbose {
		fmt.Printf("[+] Found y: %d\n", y)
	}

	// Again, Classical Part.
	Q := 1 << (2 * n)
	yQ := fraction.New(y, Q)

	r := 0
	found_r := false

	for i := 1; i <= len(yQ.ContinuedFraction()); i++ {
		ds := yQ.FractionalApprox(i)
		r = ds.D

		if r == 1 {
			continue
		}

		if verbose {
			fmt.Printf("[*] Trying with r: %d...\n", r)
		}

		if r > N {
			break
		}

		if numbers.PowMod(a, r, N) == 1 {
			found_r = true
			break
		}
	}

	if !found_r {
		if verbose {
			fmt.Println("[!] Failed to find r on this try :(")
		}
		return 0
	}

	if r%2 != 0 || numbers.PowMod(a, r/2, N) == N-1 {
		if verbose {
			fmt.Println("[!] Failed to find factor on this try :(")
		}
		return 0
	}

	if verbose {
		fmt.Printf("[+] Found r: %d\n", r)
	}

	return numbers.GCD(N, numbers.PowMod(a, r/2, N)-1)
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
