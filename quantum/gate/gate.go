package gate

import (
	"math"
	"math/cmplx"

	"github.com/sp301415/qsim/math/mat"
	"github.com/sp301415/qsim/math/numbers"
)

type Gate mat.Mat

// NewGate allocates new gate of given matrix.
func NewGate(m mat.Mat) Gate {
	if !m.IsUnitary() {
		panic("Matrix not unitary.")
	}

	if (m.NRows()&(m.NRows()-1) != 0) && (m.NRows() <= 0) {
		panic("Matrix size should be a power of two.")
	}

	return Gate(m)
}

// ToMat copies the underlying mat and returns.
func (g Gate) ToMat() mat.Mat {
	return mat.Mat(g).Copy()
}

// Copy copies g.
func (g Gate) Copy() Gate {
	return Gate(g.ToMat())
}

// Size returns the size of the gate. Here, size means the qubit length of a gate.
// For example, a 4*4 gate has size 2.
func (g Gate) Size() int {
	return numbers.BitLen(mat.Mat(g).NRows()) - 1
}

// Tensor returns the tensor product of g and given gate.
func (g Gate) Tensor(o Gate) Gate {
	return Gate(mat.Mat(g).Tensor(mat.Mat(o)))
}

// Famous Gates.

// X returns the NOT Gate (X Gate).
func X() Gate {
	return Gate(
		[][]complex128{
			{0, 1},
			{1, 0},
		},
	)
}

// Y returns the Y Gate.
func Y() Gate {
	return Gate(
		[][]complex128{
			{0, -1i},
			{1i, 0},
		},
	)
}

// Z returns the Z gate.
func Z() Gate {
	return Gate(
		[][]complex128{
			{1, 0},
			{0, -1},
		},
	)
}

// H returns the Hadamard Gate.
func H() Gate {
	h := complex(math.Sqrt2/2.0, 0)
	return Gate(
		[][]complex128{
			{h, h},
			{h, -h},
		},
	)
}

// P returns the P(phi) Gate.
func P(phi float64) Gate {
	return Gate(
		[][]complex128{
			{1, 0},
			{0, cmplx.Rect(1, phi)},
		},
	)
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
