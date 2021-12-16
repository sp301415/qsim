// Package vec implements complex128 vector.
package vec

import (
	"math"
	"math/cmplx"
)

// type Vec represents a vector.
type Vec []complex128

// Initializations.

// NewVec allocates zero vector of size n.
func NewVec(n int) Vec {
	if n <= 0 {
		panic("Invalid vector size.")
	}

	return make(Vec, n)
}

// NewVecSlice returns a new vector from complex128 slice.
// This is equivalent to Vec(v), which means it does not copy v.
func NewVecSlice(v []complex128) Vec {
	if len(v) == 0 {
		panic("Can't allocate zero size vector.")
	}

	return Vec(v)
}

// NewVecVars returns a new vector from variable arguments.
func NewVecVars(vs ...complex128) Vec {
	if len(vs) == 0 {
		panic("Can't allocate zero size vector.")
	}

	return Vec(vs)
}

// Helper functions.

// Dim returns the length of v.
func (v Vec) Dim() int {
	return len(v)
}

// Equals checks if two vectors are equal, in an error margin of 1e-6.
func (v Vec) Equals(w Vec) bool {
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

// Copy returns the copy of v.
func (v Vec) Copy() Vec {
	r := NewVec(v.Dim())
	copy(r, v)

	return r
}

// Operations.

// Add adds vector to v.
func (v Vec) Add(w Vec) Vec {
	if v.Dim() != w.Dim() {
		panic("Vector size doesn't match.")
	}

	r := v.Copy()
	for i := range r {
		r[i] += w[i]
	}

	return r
}

// Sub subtracts a vector from v.
func (v Vec) Sub(w Vec) Vec {
	if v.Dim() != w.Dim() {
		panic("Vector size doesn't match.")
	}

	r := v.Copy()
	for i := range r {
		r[i] -= w[i]
	}

	return r
}

// ScalarMul multiplies a constant complex128 to v.
func (v Vec) ScalarMul(a complex128) Vec {
	r := v.Copy()
	for i := range r {
		r[i] *= a
	}

	return r
}

// Dot returns the inner product of v and w.
func (v Vec) Dot(w Vec) complex128 {
	if v.Dim() != w.Dim() {
		panic("Vector size doesn't match.")
	}

	r := 0i
	for i := range v {
		r += v[i] * cmplx.Conj(w[i])
	}

	return r
}

// Tensor tensor products a vector to v.
func (v Vec) Tensor(w Vec) Vec {
	r := NewVec(v.Dim() * w.Dim())

	for i := 0; i < v.Dim(); i++ {
		for j := 0; j < w.Dim(); j++ {
			r[i*v.Dim()+j] = v[i] * w[j]
		}
	}

	return r
}

// NormSquared returns the squared eucledian norm of v.
// This may be useful sometimes because of numerical stability.
func (v Vec) NormSquared() float64 {
	r := 0.0

	for _, x := range v {
		r += real(x)*real(x) + imag(x)*imag(x)
	}

	return r
}

// Norm returns the eucledian norm of v.
func (v Vec) Norm() float64 {
	return math.Sqrt(v.NormSquared())
}
