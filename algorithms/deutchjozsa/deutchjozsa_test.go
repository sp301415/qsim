package deutchjozsa_test

import (
	"testing"

	"github.com/sp301415/qsim/algorithms/deutchjozsa"
)

var IS_BALANCED bool = false
var IS_CONSTANT bool = true

func BenchmarkDeutchJozsa2B(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if deutchjozsa.DeutchJozsa(2, deutchjozsa.BalancedFunc) != IS_BALANCED {
			b.Fail()
		}
	}
}

func BenchmarkDeutchJozsa2C(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if deutchjozsa.DeutchJozsa(2, deutchjozsa.ConstantFunc) != IS_CONSTANT {
			b.Fail()
		}
	}
}

func BenchmarkDeutchJozsa8B(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if deutchjozsa.DeutchJozsa(8, deutchjozsa.BalancedFunc) != IS_BALANCED {
			b.Fail()
		}
	}
}

func BenchmarkDeutchJozsa8C(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if deutchjozsa.DeutchJozsa(8, deutchjozsa.ConstantFunc) != IS_CONSTANT {
			b.Fail()
		}
	}
}

func BenchmarkDeutchJozsa16B(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if deutchjozsa.DeutchJozsa(16, deutchjozsa.BalancedFunc) != IS_BALANCED {
			b.Fail()
		}
	}
}

func BenchmarkDeutchJozsa16C(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if deutchjozsa.DeutchJozsa(16, deutchjozsa.ConstantFunc) != IS_CONSTANT {
			b.Fail()
		}
	}
}

func BenchmarkDeutchJozsaClassical16B(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if deutchjozsa.DeutchJozsaClassical(16, deutchjozsa.BalancedFunc) != IS_BALANCED {
			b.Fail()
		}
	}
}

func BenchmarkDeutchJozsaClassical16C(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if deutchjozsa.DeutchJozsaClassical(16, deutchjozsa.ConstantFunc) != IS_CONSTANT {
			b.Fail()
		}
	}
}
