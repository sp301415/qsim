// Package vector provides various functions for vectors (i.e. n*1 matrix).
package vector

import (
	"math"
	"math/cmplx"

	"github.com/sp301415/qsim/math/matrix"
)

type Vector []complex128

// Generates new vector.
func New(data []complex128) Vector {
	res := make(Vector, len(data))
	copy(res, data)

	return res
}

// Generates zero vector.
func Zeros(n int) Vector {
	res := make(Vector, n)

	for i := range res {
		res[i] = 0
	}

	return res
}

// Returns the length of vector.
func (v Vector) Dim() int {
	return len(v)
}

// Checks if two vectors are equal.
func (v Vector) Equals(w Vector) bool {
	if v.Dim() != w.Dim() {
		return false
	}

	for i := range v {
		if cmplx.Abs(v[i]-w[i]) > 1e-6 {
			return false
		}
	}

	return true
}

// Returns the dual vector.
func (v Vector) Dual() Vector {
	res := make(Vector, len(v))

	for i := range res {
		res[i] = cmplx.Conj(v[i])
	}

	return res
}

// Adds two vectors.
func (v Vector) Add(w Vector) Vector {
	if v.Dim() != w.Dim() {
		panic("Vector size doesn't match.")
	}

	res := make(Vector, len(v))

	for i := range res {
		res[i] = v[i] + w[i]
	}

	return res
}

// Alternate version of v.Add(w).
func Add(v, w Vector) Vector {
	return v.Add(w)
}

// Returns inner product of two vectors.
func (v Vector) InnerProduct(w Vector) complex128 {
	if v.Dim() != w.Dim() {
		panic("Vector size doesn't match.")
	}

	res := 0 + 0i

	for i := range v {
		res += cmplx.Conj(v[i]) * w[i]
	}

	return res
}

// Alternate version of v.InnerProduct(w)
func InnerProduct(v, w Vector) complex128 {
	return v.InnerProduct(w)
}

// Returns the outer product of two vectors.
func (v Vector) OuterProduct(w Vector) matrix.Matrix {
	if v.Dim() != w.Dim() {
		panic("Vector size doesn't match.")
	}

	d := w.Dual()
	n := v.Dim()

	res := make([][]complex128, n*n)

	for i := 0; i < n; i++ {
		res[i] = make([]complex128, n)

		for j := 0; j < n; j++ {
			res[i][j] = v[i] * d[j]
		}
	}

	return matrix.New(res)
}

// Alternate version of v.OuterProduct(w)
func OuterProduct(v, w Vector) matrix.Matrix {
	return v.OuterProduct(w)
}

// Returns tensor product of two vectors
func (v Vector) Tensor(w Vector) Vector {
	res := make(Vector, v.Dim()*w.Dim())
	idx := 0

	for _, x := range v {
		for _, y := range w {
			res[idx] = x * y
			idx++
		}
	}

	return res
}

// Alternate version of v.Tensor(w)
// Note that for convinience, this function has variable arguments.
func Tensor(vs ...Vector) Vector {
	if len(vs) == 1 {
		return vs[0]
	}

	res := vs[0]

	for _, v := range vs {
		res = res.Tensor(v)
	}

	return res
}

// Returns M*v.
func (v Vector) Apply(m matrix.Matrix) Vector {
	r, c := m.Dims()

	if c != v.Dim() {
		panic("Matrix and Vector size doesn't fit.")
	}

	res := make(Vector, r)

	for i := 0; i < r; i++ {
		res[i] = 0
		for j := 0; j < c; j++ {
			res[i] += m[i][j] * v[j]
		}
	}

	return res
}

// Returns the normalized vector.
func (v Vector) Normalize() Vector {
	norm := 0.0
	for _, n := range v {
		norm += math.Pow(cmplx.Abs(n), 2)
	}
	norm = math.Sqrt(norm)

	res := make(Vector, v.Dim())

	for i := range res {
		res[i] = v[i] / complex(norm, 0)
	}

	return res
}
