package algorithms

import (
	"github.com/sp301415/qsim"
	"github.com/sp301415/qsim/math/numbers"
	"github.com/sp301415/qsim/math/vector"
	"github.com/sp301415/qsim/quantum/qbit"
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

func BalancedOracle(regs vector.Vector, n int) vector.Vector {
	// Input: n, Output: 1
	if regs.Dim() != 1<<(n+1) {
		panic("Invalid size of registers.")
	}

	res := make(vector.Vector, regs.Dim())

	for x, a := range regs {
		res[x^BalancedFunc(x>>1)] += a
	}

	return res
}

func ConstantOracle(regs vector.Vector, n int) vector.Vector {
	// Input: n, Output: 1
	if regs.Dim() != 1<<(n+1) {
		panic("Invalid size of registers.")
	}

	res := make(vector.Vector, regs.Dim())

	for x, a := range regs {
		res[x^ConstantFunc(x>>1)] += a
	}

	return res
}

// Returns true if oracle is constant. false if it is not.
func DeutchJozsa(n int, oracle func(vector.Vector, int) vector.Vector) bool {
	// Prepare n + 1 registers with |0...01>.
	q := qsim.NewCircuit(n + 1)
	q.InitQbit(qbit.NewFromCbit(1, n+1))

	// Apply H Gate to every register.
	for i := 0; i < n+1; i++ {
		q.H(i)
	}

	// Apply Oracle!
	q.State = oracle(q.State, n)

	// Hadamard, then Measure
	for i := 1; i < n+1; i++ {
		q.H(i)
	}

	res := 0
	for i := 1; i < n+1; i++ {
		m := q.Measure(i)
		res += m * (1 << (i - 1))
	}

	return res == 0
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
