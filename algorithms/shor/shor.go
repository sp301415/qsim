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

func Shor(N int) int {
	fmt.Printf("[+] Factoring: %d\n", N)

start:

	// Classical part.
	a := 0
	for {
		a = rand.Intn(N) + 1
		K := numbers.GCD(a, N)

		if K != 1 {
			fmt.Println("[-] Found factor by luck. We start again.")
			goto start
		}
		break
	}

	fmt.Printf("[+] Using a: %d\n", a)

	n := numbers.BitLength(N)

	fmt.Println("[*] Initializing Qbit State...")

	q := qsim.NewCircuit(3 * n)
	q.InitQbit(qbit.NewFromCbit((1<<n)-1, 3*n))

	for i := n; i < 3*n; i++ {
		q.H(i)
	}

	newState := vector.Zeros(q.State.Dim())

	fmt.Println("[*] Applying Shor's Oracle...")

	for qn, qa := range q.State {
		if qa == 0 {
			continue
		}
		// Take the upmost 2n bit...
		x := qn >> n
		// pow it...
		r := numbers.PowMod(a, x, N)
		// xor with output register...
		o := (qn % (1 << n)) ^ r
		// then add it with x!
		newState[(x<<n)+o] += qa
	}

	q.State = newState

	fmt.Println("[*] Applying Inverse QFT...")

	q.InvQFT(n, 3*n)

	fmt.Println("[*] Measuring...")

	y := 0
	for i := n; i < 3*n; i++ {
		y += q.Measure(i) * (1 << (i - n))
	}

	fmt.Printf("[+] Found y: %d\n", y)

	Q := 1 << (2 * n)
	yQ := fraction.New(y, Q)

	r := 0
	found_r := false

	for i := 1; i <= len(yQ.ContinuedFraction()); i++ {
		ds := yQ.FractionalApprox(i)
		r = ds.D

		fmt.Printf("[*] Trying with r: %d...\n", r)

		if r > N {
			break
		}

		if numbers.PowMod(a, r, N) == 1 {
			found_r = true
			break
		}
	}

	if !found_r {
		fmt.Println("[!] Failed to find r on this try :(")
		goto start
	}

	if r%2 != 0 || numbers.PowMod(a, r/2, N) == N-1 {
		fmt.Println("[!] Failed to find factor on this try :(")
		goto start
	}

	fmt.Printf("[+] Found r: %d\n", r)

	return numbers.GCD(N, numbers.PowMod(a, r/2, N)-1)
}
