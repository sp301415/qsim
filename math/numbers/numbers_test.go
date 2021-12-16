package numbers_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/sp301415/qsim/math/numbers"
)

func TestPow(t *testing.T) {
	if numbers.Pow(2, 3) != 8 {
		t.Fail()
	}
}

func TestPowMod(t *testing.T) {
	if numbers.PowMod(2, 4, 7) != 2 {
		t.Fail()
	}
}

func TestGCD(t *testing.T) {
	if numbers.GCD(12, 16) != 4 {
		t.Fail()
	}
}

func TestBinLen(t *testing.T) {
	n := 1545
	l := len(fmt.Sprintf("%b", n))

	if l != numbers.BitLen(n) {
		t.Fail()
	}
}

func TestMinMax(t *testing.T) {
	N := 100
	data := make([]int, N)

	for i := range data {
		data[i] = i
	}

	rand.Shuffle(N, func(i, j int) { data[i], data[j] = data[j], data[i] })

	if numbers.Min(data...) != 0 || numbers.Max(data...) != N-1 {
		t.Fail()
	}
}
