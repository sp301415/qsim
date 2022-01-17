// Package mat implements complex128 matrix.
package mat

import (
	"fmt"
	"math/cmplx"

	"github.com/sp301415/qsim/math/vec"
)

// type Mat represents a matrix.
type Mat [][]complex128

// Initializations.

// NewMat allocates zero matrix of size r * c.
func NewMat(r, c int) Mat {
	if r <= 0 || c <= 0 {
		panic("Invalid matrix size.")
	}

	m := make(Mat, r)

	for i := range m {
		m[i] = make([]complex128, c)
	}

	return m
}

// NewMatVecs allocates new matrix from slice of vectors.
// Note that function copies vectors, unlike NewVecSlice or NewMatSlice.
// This is because it is likely that vectors will be reused.
func NewMatVecs(vs []vec.Vec) Mat {
	if len(vs) == 0 {
		panic("Can't allocate zero sized matrix.")
	}

	r := len(vs)
	c := len(vs[0])
	m := NewMat(r, c)
	for i, v := range vs {
		if v.Dim() != c {
			panic("Every vector should have same length.")
		}
		m[i] = v.Copy()
	}

	if len(m[0]) == 0 {
		panic("Zero length vectors not allowed.")
	}

	return m
}

// NewMatSlice returns new matrix from complex128 double slice.
// This is equivalant to Mat(v), which means that it does not copy v.
func NewMatSlice(v [][]complex128) Mat {
	if len(v) == 0 {
		panic("Can't allocate zero sized matrix.")
	}

	return Mat(v)
}

// NewMatSliceVars returns new matrix from variable arguments of slices.
// This does not copy each slices.
func NewMatSliceVars(vs ...[]complex128) Mat {
	if len(vs) == 0 {
		panic("Can't allocate zero sized matrix.")
	}

	for _, v := range vs {
		if len(v) == 0 {
			panic("Can't allocate zero sized matrix.")
		}
	}

	return Mat(vs)
}

func NewMatVars(nrows int, vs ...complex128) Mat {
	if len(vs) == 0 {
		panic("Can't allocate zero sized matrix.")
	}

	if len(vs)%nrows != 0 || nrows <= 0 {
		panic("Invalid number of rows.")
	}

	ncols := len(vs) / nrows

	m := make([][]complex128, nrows)
	for i := 0; i < nrows; i++ {
		m[i] = append(make([]complex128, 0), vs[i*ncols:(i+1)*ncols]...)
	}

	return Mat(m)
}

// NewSquare allocates the zero square matrix of size n * n.
func NewSquare(n int) Mat {
	if n <= 0 {
		panic("Can't allocate zero sized matrix.")
	}

	return NewMat(n, n)
}

// NewId allocates the identity matrix of size n * n.
func NewId(n int) Mat {
	m := NewSquare(n)

	for i := 0; i < n; i++ {
		m[i][i] = 1
	}

	return m
}

// Checking functions.

// IsSquare checks if m is a square matrix.
func (m Mat) IsSquare() bool {
	return m.NCols() == m.NRows()
}

// IsUnitary checks if m is a unitary matrix.
func (m Mat) IsUnitary() bool {
	if !m.IsSquare() {
		return false
	}

	return m.Mul(m.Dagger()).Equals(NewId(m.NRows()))
}

// Helper functions.

// NRows returns the number of rows. (or, the length of a column.)
func (m Mat) NRows() int {
	return len(m)
}

// GetRow returns the ith row as a vector.
func (m Mat) GetRow(i int) vec.Vec {
	return vec.NewVecSlice(m[i]).Copy()
}

// SetRow sets the ith row to given vector.
func (m *Mat) SetRow(i int, v vec.Vec) {
	(*m)[i] = v.Copy()
}

// NCols returns the number of columns. (or, the length of a row.)
func (m Mat) NCols() int {
	return len(m[0])
}

// GetCol returns the ith column as a vector.
func (m Mat) GetCol(i int) vec.Vec {
	c := vec.NewVec(m.NRows())
	for k := 0; k < m.NRows(); k++ {
		c[k] = m[k][i]
	}

	return c
}

// SetCol sets the ith column to given vector.
func (m *Mat) SetCol(i int, v vec.Vec) {
	for k := 0; k < m.NRows(); k++ {
		(*m)[k][i] = v[i]
	}
}

// Dim returns the dimension of the matrix.
func (m Mat) Dim() (int, int) {
	return m.NRows(), m.NCols()
}

// Equals checks if two matrices are equal, in an error margin of 1e-6.
func (m Mat) Equals(n Mat) bool {
	if m.NRows() != n.NRows() || m.NCols() != n.NCols() {
		return false
	}

	for i, row := range n {
		for j := range row {
			if cmplx.Abs(n[i][j]-m[i][j]) > 1e-6 {
				return false
			}
		}
	}

	return true
}

// Copy returns the copy of m.
func (m Mat) Copy() Mat {
	r := NewMat(m.NRows(), m.NCols())
	for i := range r {
		r[i] = m.GetRow(i)
	}

	return r
}

// Operations.

// Addd adds matrix to m.
func (m Mat) Add(n Mat) Mat {
	if m.NRows() != n.NRows() || m.NCols() != n.NCols() {
		panic("Matrix size doesn't match.")
	}

	r := m.Copy()
	for i, row := range r {
		for j := range row {
			r[i][j] += n[i][j]
		}
	}

	return r
}

// Sub subtracts matrix from m.
func (m Mat) Sub(n Mat) Mat {
	if m.NRows() != n.NRows() || m.NCols() != n.NCols() {
		panic("Matrix size doesn't match.")
	}

	r := m.Copy()
	for i, row := range r {
		for j := range row {
			r[i][j] -= n[i][j]
		}
	}

	return r
}

// ScalarMul multiplies a constant complex128 to m.
func (m Mat) ScalarMul(a complex128) Mat {
	r := m.Copy()
	for i, row := range r {
		for j := range row {
			r[i][j] *= a
		}
	}

	return r
}

// Mul multiplies matrix to m. Size of m is automatically adjusted.
func (m Mat) Mul(n Mat) Mat {
	if m.NCols() != n.NRows() {
		panic("Matrix size doesn't match.")
	}

	r := NewMat(m.NRows(), n.NCols())
	for i := 0; i < m.NRows(); i++ {
		for j := 0; j < n.NCols(); j++ {
			for k := 0; k < m.NCols(); k++ {
				r[i][j] += m[i][k] * n[k][j]
			}
		}
	}

	return r
}

// Tensor tensor products a matrix to m. Size of m is automatically adjusted.
func (m Mat) Tensor(n Mat) Mat {
	r := NewMat(m.NRows()*n.NRows(), m.NCols()*n.NCols())

	for i := 0; i < m.NRows(); i++ {
		for j := 0; j < n.NRows(); j++ {
			for k := 0; k < m.NCols(); k++ {
				for l := 0; l < n.NCols(); l++ {
					r[i*n.NRows()+j][k*n.NCols()+l] = m[i][k] * n[j][l]
				}
			}
		}
	}

	return r
}

// Dagger returns a conjugate transpose of m.
func (m Mat) Dagger() Mat {
	r := NewMat(m.NCols(), m.NRows())

	for i, row := range m {
		for j := range row {
			r[j][i] = cmplx.Conj(m[i][j])
		}
	}

	return r
}

// String implements the Stringer interface.
func (m Mat) String() string {
	r := ""

	for _, row := range m {
		r += fmt.Sprintln(row)
	}

	return r
}
