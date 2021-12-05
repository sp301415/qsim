// Package gate provides collection of popular gates.
package gate

import (
	"math"
	"math/cmplx"

	"github.com/sp301415/qsim/math/matrix"
)

// Returns single-qbit I gate.
func I() matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{1, 0},
			{0, 1},
		},
	)
}

// Returns single-qbit X gate.
func X() matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{0, 1},
			{1, 0},
		},
	)
}

// Returns single-qbit Y gate.
func Y() matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{0, -1i},
			{1i, 0},
		},
	)
}

// Returns single-qbit Z gate.
func Z() matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{1, 0},
			{0, -1},
		},
	)
}

// Returns single-qbit H gate.
func H() matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{complex(math.Sqrt(0.5), 0), complex(math.Sqrt(0.5), 0)},
			{complex(math.Sqrt(0.5), 0), complex(-math.Sqrt(0.5), 0)},
		},
	)
}

// Returns single-qbit P(phi) gate.
func P(phi float64) matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{1, 0},
			{0, cmplx.Rect(1, phi)},
		},
	)
}

// Returns single-qbit S gate. Same as P(pi/2).
func S() matrix.Matrix {
	return P(math.Pi / 2.0)
}

// Returns single-qbit T gate. Same as P(pi/4).
func T() matrix.Matrix {
	return P(math.Pi / 4.0)
}
