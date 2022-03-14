package qsim

import (
	"fmt"
	"math/cmplx"

	"github.com/sp301415/qsim/math/number"
	"github.com/sp301415/qsim/math/vec"
)

type Qubit struct {
	data vec.Vec
	size int
}

// NewQubit allocates new qubit from given vector. This does not copy vector.
// NOTE: This checks whether length of v is power of 2. Direct casting might lead to unexpected behavior.
func NewQubit(v vec.Vec) Qubit {
	if (v.Dim()&(v.Dim()-1) != 0) && (v.Dim() <= 0) {
		panic("Vector size must be power of two.")
	}

	return Qubit{data: v, size: number.BitLen(v.Dim()) - 1}
}

// NewBit allocates new qubit from classical bit.
// If size == 0, then it automatically adjust sizes.
func NewBit(n, size int) Qubit {
	if n < 0 {
		panic("Cannot convert negative number to qubit.")
	}

	if size == 0 {
		size = number.BitLen(n)
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
	return q.data.Copy()
}

// Copy copies q.
func (q Qubit) Copy() Qubit {
	return Qubit{data: q.data.Copy(), size: q.size}
}

// Size returns the bit length of a qubit.
func (q Qubit) Size() int {
	return q.size
}

// At returns ith element of a qubit.
func (q Qubit) At(i int) complex128 {
	return q.data[i]
}

// Dim returns the vector length of a qubit.
func (q Qubit) Dim() int {
	return q.data.Dim()
}

// Equals checks if two qubits are equal.
func (q Qubit) Equals(p Qubit) bool {
	return q.data.Equals(p.data)
}

// String implements Stringer interface.
func (q Qubit) String() string {
	r := ""
	idxpad := number.BitLen(q.Size())

	for n, a := range q.data {
		if cmplx.Abs(a) < 1e-6 {
			continue
		}

		r += fmt.Sprintf("[%*d] |%0*b>: %f\n", idxpad, n, q.Size(), n, a)
	}
	return r
}
