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

func BenchmarkPowMod(t *testing.B) {
	a := rand.Intn(1 << 20)
	b := rand.Intn(1 << 20)
	c := rand.Intn(1 << 20)

	numbers.PowMod(a, b, c)
}

func TestBinLen(t *testing.T) {
	n := 1545
	l := len(fmt.Sprintf("%b", n))

	if l != numbers.BitLength(n) {
		t.Fail()
	}
}
