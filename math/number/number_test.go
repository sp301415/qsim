package number_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/sp301415/qsim/math/number"
	"github.com/sp301415/qsim/utils/slice"
)

func TestPow(t *testing.T) {
	if number.Pow(2, 3) != 8 {
		t.Fail()
	}
}

func TestPowMod(t *testing.T) {
	if number.PowMod(2, 4, 7) != 2 {
		t.Fail()
	}
}

func TestGCD(t *testing.T) {
	if number.GCD(12, 16) != 4 {
		t.Fail()
	}
}

func TestBinLen(t *testing.T) {
	n := 1545
	l := len(fmt.Sprintf("%b", n))

	if l != number.BitLen(n) {
		t.Fail()
	}
}

func TestMinMax(t *testing.T) {
	N := 100
	data := slice.Sequence(0, N)

	rand.Shuffle(N, func(i, j int) { data[i], data[j] = data[j], data[i] })

	if number.Min(data...) != 0 || number.Max(data...) != N-1 {
		t.Fail()
	}
}
