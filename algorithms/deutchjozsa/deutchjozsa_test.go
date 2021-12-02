package deutchjozsa_test

import (
	"testing"

	"github.com/sp301415/qsim/algorithms/deutchjozsa"
)

func BenchmarkDeutchJozsa2B(t *testing.B) {
	deutchjozsa.DeutchJozsa(2, deutchjozsa.BalancedOracle)
}

func BenchmarkDeutchJozsa2C(t *testing.B) {
	deutchjozsa.DeutchJozsa(2, deutchjozsa.ConstantOracle)
}

func BenchmarkDeutchJozsa8B(t *testing.B) {
	deutchjozsa.DeutchJozsa(8, deutchjozsa.BalancedOracle)
}

func BenchmarkDeutchJozsa8C(t *testing.B) {
	deutchjozsa.DeutchJozsa(8, deutchjozsa.ConstantOracle)
}

func BenchmarkDeutchJozsa16B(t *testing.B) {
	deutchjozsa.DeutchJozsa(16, deutchjozsa.BalancedOracle)
}

func BenchmarkDeutchJozsa16C(t *testing.B) {
	deutchjozsa.DeutchJozsa(16, deutchjozsa.ConstantOracle)
}

func BenchmarkDeutchJozsaClassical16B(t *testing.B) {
	deutchjozsa.DeutchJozsaClassical(16, deutchjozsa.BalancedFunc)
}

func BenchmarkDeutchJozsaClassical16C(t *testing.B) {
	deutchjozsa.DeutchJozsaClassical(16, deutchjozsa.ConstantFunc)
}
