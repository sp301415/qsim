// Package matrix provides various functions for matrix operations.
package matrix

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

// Returns the dimension(row, column) of matrix.
func (m Matrix) Dims() (int, int) {
	return len(m), len(m[0])
}

// Checks if two matrices are equal
func (m Matrix) Equals(n Matrix) bool {
	r1, c1 := m.Dims()
	r2, c2 := n.Dims()

	if (r1 != r2) || (c1 != c2) {
		return false
	}

	for i := 0; i < r1; i++ {
		for j := 0; j < c1; j++ {
			if m[i][j] != n[i][j] {
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

// Tensor product of two matrices.
func (m Matrix) Tensor(n Matrix) Matrix {
	r1, c1 := m.Dims()
	r2, c2 := n.Dims()

	res := make(Matrix, r1*r2)

	for i1 := 0; i1 < r1; i1++ {
		for i2 := 0; i2 < r2; i2++ {
			i := i1*r2 + i2
			m[i] = make([]complex128, c1*c2)

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
func Tensor(m ...Matrix) Matrix {
	if len(m) == 1 {
		return m[0]
	}

	res := m[0]
	for i := 1; i < len(m); i++ {
		res = res.Tensor(m[i])
	}

	return res
}
