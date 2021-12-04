package shor_test

import (
	"testing"

	"github.com/sp301415/qsim/algorithms/shor"
)

func BenchmarkShor15(t *testing.B) {
	N := 15
	factor := shor.Shor(N)

	if N%factor != 0 {
		t.Fail()
	}
}

func BenchmarkShor35(t *testing.B) {
	N := 35
	factor := shor.Shor(N)

	if N%factor != 0 {
		t.Fail()
	}
}

func BenchmarkShor55(t *testing.B) {
	N := 55
	factor := shor.Shor(N)

	if N%factor != 0 {
		t.Fail()
	}
}
