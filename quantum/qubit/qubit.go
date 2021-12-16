package qubit

import (
	"fmt"
	"math/cmplx"
	"strings"

	"github.com/sp301415/qsim/math/numbers"
	"github.com/sp301415/qsim/math/vec"
)

type Qubit vec.Vec

// NewQubit allocates new qubit from given vector. This does not copy vector.
// NOTE: This checks whether length of v is power of 2. Direct casting might lead to unexpected behavior.
func NewQubit(v vec.Vec) Qubit {
	if (v.Dim()&(v.Dim()-1) != 0) && (v.Dim() <= 0) {
		panic("Vector size must be power of two.")
	}

	return Qubit(v)
}

// NewBit allocates new qubit from classical bit.
// If size == 0, then it automatically adjust sizes.
func NewBit(n, size int) Qubit {
	if n < 0 {
		panic("Cannot convert negative number to qubit.")
	}

	if size == 0 {
		size = numbers.BitLen(n)
	}

	if n > 1<<size {
		panic("Size too small.")
	}

	v := vec.NewVec(1 << size)
	v[n] = 1

	return NewQubit(v)
}

// ToVec copies the underlying vec and returns.
func (q Qubit) ToVec() vec.Vec {
	return vec.Vec(q).Copy()
}

// Copy copies q.
func (q Qubit) Copy() Qubit {
	return Qubit(q.ToVec())
}

// Size returns the bit length of a qubit.
func (q Qubit) Size() int {
	return numbers.BitLen(vec.Vec(q).Dim()) - 1
}

// Dim returns the vector length of a qubit.
func (q Qubit) Dim() int {
	return vec.Vec(q).Dim()
}

// Equals checks if two qubits are equal.
func (q Qubit) Equals(p Qubit) bool {
	return vec.Vec(q).Equals(vec.Vec(p))
}

// String implements Stringer interface.
func (q Qubit) String() string {
	r := ""

	for n, a := range q {
		if cmplx.Abs(a) < 1e-6 {
			continue
		}

		r += fmt.Sprintf("|%0*b>: %f\n", q.Size(), n, a)
	}
	return strings.TrimSpace(r)
}
