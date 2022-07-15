package shor_test

import (
	"testing"

	"github.com/sp301415/qsim/algorithms/shor"
)

func BenchmarkShor15(b *testing.B) {
	N := 15
	for i := 0; i < b.N; i++ {
		factor := shor.Shor(N)

		if N%factor != 0 {
			b.Fail()
		}
	}
}

func BenchmarkShor35(b *testing.B) {
	N := 35

	for i := 0; i < b.N; i++ {
		factor := shor.Shor(N)

		if N%factor != 0 {
			b.Fail()
		}
	}
}

func BenchmarkShor55(b *testing.B) {
	N := 55

	for i := 0; i < b.N; i++ {
		factor := shor.Shor(N)

		if N%factor != 0 {
			b.Fail()
		}
	}
}

func BenchmarkShor85(b *testing.B) {
	N := 85

	for i := 0; i < b.N; i++ {
		factor := shor.Shor(N)

		if N%factor != 0 {
			b.Fail()
		}
	}
}
