// Package gate provides collection of popular gates.
package gate

import (
	"math"
	"math/cmplx"

	"github.com/sp301415/qsim/math/matrix"
)

func I() matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{1, 0},
			{0, 1},
		},
	)
}

func X() matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{0, 1},
			{1, 0},
		},
	)
}

func Y() matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{0, -1i},
			{1i, 0},
		},
	)
}

func Z() matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{1, 0},
			{0, -1},
		},
	)
}

func H() matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{complex(math.Sqrt(0.5), 0), complex(math.Sqrt(0.5), 0)},
			{complex(math.Sqrt(0.5), 0), complex(-math.Sqrt(0.5), 0)},
		},
	)
}

func P(phi float64) matrix.Matrix {
	return matrix.New(
		[][]complex128{
			{1, 0},
			{0, cmplx.Rect(1, phi)},
		},
	)
}
