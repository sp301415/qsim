package qsim

import (
	"math"
	"math/cmplx"

	"github.com/sp301415/qsim/math/mat"
	"github.com/sp301415/qsim/math/numbers"
)

type Gate struct {
	data mat.Mat
	size int
}

// NewGate allocates new gate of given matrix.
func NewGate(m mat.Mat) Gate {
	if !m.IsUnitary() {
		panic("Matrix not unitary.")
	}

	if (m.NRows()&(m.NRows()-1) != 0) && (m.NRows() <= 0) {
		panic("Matrix size should be a power of two.")
	}

	return Gate{data: m, size: numbers.BitLen(m.NRows()) - 1}
}

// ToMat copies the underlying mat and returns.
func (g Gate) ToMat() mat.Mat {
	return g.data.Copy()
}

// Copy copies g.
func (g Gate) Copy() Gate {
	return Gate{data: g.data.Copy()}
}

// Size returns the size of the gate. Here, size means the qubit length of a gate.
// For example, a 4*4 gate has size 2.
func (g Gate) Size() int {
	return g.size
}

// At returns the (i, j)th element.
func (g Gate) At(i, j int) complex128 {
	return g.data[i][j]
}

// Tensor returns the tensor product of g and given gate.
func (g Gate) Tensor(o Gate) Gate {
	return Gate{
		data: g.data.Tensor(o.data),
		size: g.size + o.size,
	}
}

// Famous Gates.

// X returns the NOT Gate (X Gate).
func X() Gate {
	return Gate{
		data: [][]complex128{
			{0, 1},
			{1, 0},
		},
		size: 1,
	}
}

// Y returns the Y Gate.
func Y() Gate {
	return Gate{data: [][]complex128{
		{0, -1i},
		{1i, 0},
	},
		size: 1,
	}
}

// Z returns the Z gate.
func Z() Gate {
	return Gate{data: [][]complex128{
		{1, 0},
		{0, -1},
	},
		size: 1,
	}
}

// H returns the Hadamard Gate.
func H() Gate {
	h := complex(math.Sqrt2/2.0, 0)
	return Gate{
		data: [][]complex128{
			{h, h},
			{h, -h},
		},
		size: 1,
	}
}

// P returns the P(phi) Gate.
func P(phi float64) Gate {
	return Gate{data: [][]complex128{
		{1, 0},
		{0, cmplx.Rect(1, phi)},
	},
		size: 1,
	}
}

// S returns the S Gate. Same as P(pi/2).
func S() Gate {
	return P(math.Pi / 2.0)
}

// T returns the T gate. Same as P(pi/4).
func T() Gate {
	return P(math.Pi / 4.0)
}

// String implements Stringer interface.
func (g Gate) String() string {
	return g.ToMat().String()
}
