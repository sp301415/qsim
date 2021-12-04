package shor_test

import (
	"testing"

	"github.com/sp301415/qsim/algorithms/shor"
)

func BenchmarkShor15(t *testing.B) {
	shor.Shor(15)
}

func BenchmarkShor35(t *testing.B) {
	shor.Shor(35)
}
