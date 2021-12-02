// Package qbit provides functions for quantum state.
// Qbit is not much different from a vector, though.
package qbit

import (
	"github.com/sp301415/qsim/math/numbers"
	"github.com/sp301415/qsim/math/vector"
)

// Returns |0>
func Zero() vector.Vector {
	return Zeros(1)
}

// Returns |00...0>
func Zeros(n int) vector.Vector {
	data := make(vector.Vector, numbers.Pow(2, n))
	data[0] = 1

	return vector.New(data)
}

// Returns |1>
func One() vector.Vector {
	return Ones(1)
}

// Returns |11...1>
func Ones(n int) vector.Vector {
	data := make(vector.Vector, numbers.Pow(2, n))
	data[len(data)-1] = 1

	return vector.New(data)
}

// Changes cbit to qbit. Useful for initialization.
// Note: l denotes the size of the qbit.
// If l == 0, then it automatically finds the right size.
func NewFromCbit(n int, l int) vector.Vector {
	if l == 0 {
		l = numbers.Pow(2, numbers.BitLength(n))
	}

	q := make(vector.Vector, numbers.Pow(2, l))
	q[n] = 1

	return vector.New(q)
}

func BitLength(q vector.Vector) int {
	return numbers.BitLength(q.Dim())
}
