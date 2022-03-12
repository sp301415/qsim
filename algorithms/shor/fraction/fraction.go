// Fraction provide packages for fraction.
package fraction

import (
	"fmt"

	"github.com/sp301415/qsim/math/number"
)

// N / D
type Fraction struct {
	D int
	N int
}

// Generates new fraction.
func New(N, D int) Fraction {
	if D == 0 {
		panic("Division by zero.")
	}

	g := number.GCD(N, D)

	return Fraction{N: N / g, D: D / g}
}

// Prints fraction to string.
func (f Fraction) String() string {
	return fmt.Sprintf("Fraction{%d, %d}", f.N, f.D)
}

// Returns fraction as float.
func (f Fraction) Float64() float64 {
	return float64(f.N) / float64(f.D)
}

// Returns fraction as int.
func (f Fraction) Int() int {
	return f.N / f.D
}

// Returns the continued fraction expression of this fraction.
func (f Fraction) ContinuedFraction() []int {
	res := make([]int, 0)
	ff := f
	a := 0

	for i := 0; ; i++ {
		a = ff.Int()
		res = append(res, a)
		ND := ff.N - a*ff.D

		if ND == 0 {
			break
		}

		ff = New(ff.D, ND)
	}

	return res
}

// Returns all possible fractional approximation of this fraction.
func (f Fraction) FractionalApprox() []Fraction {
	cf := f.ContinuedFraction()

	n := len(cf)
	Ns := make([]int, n+2)
	Ds := make([]int, n+2)

	res := make([]Fraction, n)

	Ns[1] = 1
	Ds[0] = 1

	for i := 2; i < n+2; i++ {
		a := cf[i-2]
		Ns[i] = a*Ns[i-1] + Ns[i-2]
		Ds[i] = a*Ds[i-1] + Ds[i-2]
		res[i-2] = New(Ns[i], Ds[i])
	}

	return res
}
