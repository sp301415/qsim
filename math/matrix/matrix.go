// Package matrix provides various functions for matrix operations.
package matrix

import (
	"math/cmplx"
)

type Matrix [][]complex128

// Generates new matrix.
func New(data [][]complex128) Matrix {
	r := len(data)
	c := len(data[0])

	ret := make(Matrix, r)

	for i := 0; i < r; i++ {
		ret[i] = make([]complex128, c)
		for j := 0; j < c; j++ {
			ret[i][j] = data[i][j]
		}
	}

	return ret
}

// Generates zero square matrix.
func Zeros(n int) Matrix {
	ret := make(Matrix, n)

	for i := 0; i < n; i++ {
		ret[i] = make([]complex128, n)
		for j := 0; j < n; j++ {
			ret[i][j] = 0
		}
	}

	return ret
}

// Generates Identity matrix.
func Identity(n int) Matrix {
	m := Zeros(n)

	for i := 0; i < len(m); i++ {
		m[i][i] = 1
	}

	return m
}

// Returns the dimension(row, column) of matrix.
func (m Matrix) Dims() (int, int) {
	return len(m), len(m[0])
}

// Checks if two matrices are equal.
func (m Matrix) Equals(n Matrix) bool {
	r1, c1 := m.Dims()
	r2, c2 := n.Dims()

	if (r1 != r2) || (c1 != c2) {
		return false
	}

	for i := 0; i < r1; i++ {
		for j := 0; j < c1; j++ {
			if cmplx.Abs(m[i][j]-n[i][j]) > 1e-6 {
				return false
			}
		}
	}

	return true
}

// Adds two matrices, returning it.
// Does not change the original matrix.
func (m Matrix) Add(n Matrix) Matrix {
	r, c := m.Dims()
	rr, cc := n.Dims()

	if (r != rr) || (c != cc) {
		panic("Matrix size does not match.")
	}

	ret := make(Matrix, r)

	for i := 0; i < r; i++ {
		ret[i] = make([]complex128, c)
		for j := 0; j < c; j++ {
			ret[i][j] = m[i][j] + n[i][j]
		}
	}

	return ret
}

// Alternative version of m.Add(n).
func Add(m, n Matrix) Matrix {
	return m.Add(n)
}

// Multiplies two matrices, returning it.
// Does not change the original matrix.
func (m Matrix) Mul(n Matrix) Matrix {
	r1, c1 := m.Dims()
	r2, c2 := n.Dims()

	if c1 != r2 {
		panic("Matrix size does not match.")
	}

	res := make(Matrix, r1)

	for i := 0; i < r1; i++ {
		res[i] = make([]complex128, c2)
		for j := 0; j < c1; j++ {
			t := complex128(0)
			for k := 0; k < c1; k++ {
				t += m[i][k] * n[k][j]
			}
			res[i][j] = t
		}
	}

	return res
}

// Alternative version of m.Mul(n).
func Mul(m, n Matrix) Matrix {
	return m.Mul(n)
}

// Multiplies a float constant to a matrix.
func (m Matrix) CMul(k complex128) Matrix {
	r, c := m.Dims()
	res := make(Matrix, r)

	for i := 0; i < r; i++ {
		res[i] = make([]complex128, c)

		for j := 0; j < c; j++ {
			res[i][j] *= k
		}
	}

	return res
}

// Tensor product of two matrices.
func (m Matrix) Tensor(n Matrix) Matrix {
	r1, c1 := m.Dims()
	r2, c2 := n.Dims()

	res := make(Matrix, r1*r2)

	for i1 := 0; i1 < r1; i1++ {
		for i2 := 0; i2 < r2; i2++ {
			i := i1*r2 + i2
			res[i] = make([]complex128, c1*c2)

			for j1 := 0; j1 < c1; j1++ {
				for j2 := 0; j2 < c2; j2++ {
					j := j1*c2 + j2
					res[i][j] = m[i1][j1] * n[i2][j2]
				}
			}
		}
	}

	return res
}

// Alternative version of m.Tensor(n)
// Note that for convinience, this function has variable arguments.
func Tensor(ms ...Matrix) Matrix {
	if len(ms) == 1 {
		return ms[0]
	}

	res := ms[0]
	for i := 1; i < len(ms); i++ {
		res = res.Tensor(ms[i])
	}

	return res
}

// Dagger = transpose conjugate
func (m Matrix) Dagger() Matrix {
	r, c := m.Dims()
	n := make(Matrix, r)

	for i := 0; i < r; i++ {
		n[i] = make([]complex128, c)

		for j := 0; j < c; j++ {
			n[i][j] = cmplx.Conj(m[j][i])
		}
	}

	return n
}

func (m Matrix) IsSquare() bool {
	r, c := m.Dims()

	return r == c
}

func (m Matrix) IsUnitary() bool {
	if !m.IsSquare() {
		return false
	}

	n, _ := m.Dims()
	id := Identity(n)

	return Mul(m, m.Dagger()).Equals(id)
}
