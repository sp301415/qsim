package algorithms_test

import (
	"testing"

	"github.com/sp301415/qsim/algorithms"
)

func BenchmarkDeutchJozsa2B(t *testing.B) {
	algorithms.DeutchJozsa(2, algorithms.BalancedOracle)
}

func BenchmarkDeutchJozsa2C(t *testing.B) {
	algorithms.DeutchJozsa(2, algorithms.ConstantOracle)
}

func BenchmarkDeutchJozsa8B(t *testing.B) {
	algorithms.DeutchJozsa(8, algorithms.BalancedOracle)
}

func BenchmarkDeutchJozsa8C(t *testing.B) {
	algorithms.DeutchJozsa(8, algorithms.ConstantOracle)
}

func BenchmarkDeutchJozsa16B(t *testing.B) {
	algorithms.DeutchJozsa(16, algorithms.BalancedOracle)
}

func BenchmarkDeutchJozsa16C(t *testing.B) {
	algorithms.DeutchJozsa(16, algorithms.ConstantOracle)
}

func BenchmarkDeutchJozsaClassical16B(t *testing.B) {
	algorithms.DeutchJozsaClassical(16, algorithms.BalancedFunc)
}

func BenchmarkDeutchJozsaClassical16C(t *testing.B) {
	algorithms.DeutchJozsaClassical(16, algorithms.ConstantFunc)
}
