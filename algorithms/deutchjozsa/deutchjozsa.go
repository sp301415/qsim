package deutchjozsa

import (
	"github.com/sp301415/qsim"
	"github.com/sp301415/qsim/math/number"
	"github.com/sp301415/qsim/utils/slice"
)

func BalancedFunc(x int) int {
	r := 0

	for i := 0; i < number.BitLen(x); i++ {
		r ^= (x >> i) % 2
	}

	return r
}

func ConstantFunc(x int) int {
	return 0
}

// Returns true if oracle is constant. false if it is not.
func DeutchJozsa(n int, oracle func(int) int) bool {
	iregs := slice.Range(0, n)

	// Prepare n + 1 registers with |0...01>.
	q := qsim.NewCircuit(n + 1)
	q.X(n)

	// Apply H Gate to every register.
	q.H(iregs...)
	q.H(n)

	// Apply Oracle!
	q.ApplyOracle(oracle, iregs, []int{n})

	// Hadamard, then Measure
	q.H(iregs...)
	res := q.Measure(iregs...)

	if res == 0 {
		return true
	} else if res == (1<<n)-1 {
		return false
	} else {
		panic("INVALID MEASUREMENT")
	}
}

// Returns true if oracle is constant. false if it is not.
func DeutchJozsaClassical(n int, oracle func(int) int) bool {
	// Evaluate 2^(n-1) + 1 times
	cnt := [2]int{0, 0}

	for i := 0; i < (1<<(n-1))+1; i++ {
		cnt[oracle(i)] += 1
	}

	return cnt[0] > (1<<(n-1)) || cnt[1] > (1<<(n-1))
}
