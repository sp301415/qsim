package deutchjozsa

import (
	"github.com/sp301415/qsim"
	"github.com/sp301415/qsim/math/numbers"
	"github.com/sp301415/qsim/quantum/qubit"
)

func BalancedFunc(x int) int {
	r := 0

	for i := 0; i < numbers.BitLength(x); i++ {
		r ^= (x >> i) % 2
	}

	return r
}

func ConstantFunc(x int) int {
	return 0
}

// Returns true if oracle is constant. false if it is not.
func DeutchJozsa(n int, oracle func(int) int) bool {
	iregs := make([]int, n)
	for i := range iregs {
		iregs[i] = i + 1
	}

	// Prepare n + 1 registers with |0...01>.
	q := qsim.NewCircuit(n + 1)
	q.InitQubit(qubit.NewFromCbit(1, n+1))

	// Apply H Gate to every register.
	for i := 0; i < n+1; i++ {
		q.H(i)
	}

	// Apply Oracle!
	q.ApplyOracle(oracle, iregs, []int{0})

	// Hadamard, then Measure
	for _, v := range iregs {
		q.H(v)
	}
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
