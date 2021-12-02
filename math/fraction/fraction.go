// Fraction provide packages for fraction. (Frankly just for Shor)
package fraction

import (
	"fmt"

	"github.com/sp301415/qsim/math/numbers"
)

// N / D
type Fraction struct {
	D int
	N int
}

func New(N, D int) Fraction {
	g := numbers.GCD(N, D)

	return Fraction{N: N / g, D: D / g}
}

func (f Fraction) String() string {
	return fmt.Sprintf("Fraction{%d, %d}", f.N, f.D)
}

func (f Fraction) Float64() float64 {
	return float64(f.N) / float64(f.D)
}

func (f Fraction) Int() int {
	return f.N / f.D
}

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

// if n == 0, then go to the end.
func (f Fraction) FractionalApprox(n int) Fraction {
	cf := f.ContinuedFraction()

	if n == 0 {
		n = len(cf)
	} else {
		n = numbers.Min(n, len(cf))
	}

	Ns := make([]int, n+2)
	Ds := make([]int, n+2)

	Ns[1] = 1
	Ds[0] = 1

	for i := 2; i < n+2; i++ {
		a := cf[i-2]
		Ns[i] = a*Ns[i-1] + Ns[i-2]
		Ds[i] = a*Ds[i-1] + Ds[i-2]
	}

	return New(Ns[n+1], Ds[n+1])
}
