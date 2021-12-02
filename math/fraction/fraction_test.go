package fraction_test

import (
	"fmt"
	"testing"

	"github.com/sp301415/qsim/math/fraction"
)

func TestContinuedFraction(t *testing.T) {
	f := fraction.New(649, 200)

	c := f.ContinuedFraction()
	r := []int{3, 4, 12, 4}

	for i, v := range r {
		if v != c[i] {
			t.Fail()
		}
	}
}

func TestApprox(t *testing.T) {
	f := fraction.New(84375, 100000)
	fmt.Println(f.FractionalApprox(0))
	if f.FractionalApprox(0) != fraction.New(27, 32) {
		t.Fail()
	}
}
