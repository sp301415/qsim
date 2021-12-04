package deutchjozsa_test

import (
	"testing"

	"github.com/sp301415/qsim/algorithms/deutchjozsa"
)

var IS_BALANCED bool = false
var IS_CONSTANT bool = true

func BenchmarkDeutchJozsa2B(t *testing.B) {
	if deutchjozsa.DeutchJozsa(2, deutchjozsa.BalancedFunc) != IS_BALANCED {
		t.Fail()
	}
}

func BenchmarkDeutchJozsa2C(t *testing.B) {
	if deutchjozsa.DeutchJozsa(2, deutchjozsa.ConstantFunc) != IS_CONSTANT {
		t.Fail()
	}
}

func BenchmarkDeutchJozsa8B(t *testing.B) {
	if deutchjozsa.DeutchJozsa(8, deutchjozsa.BalancedFunc) != IS_BALANCED {
		t.Fail()
	}
}

func BenchmarkDeutchJozsa8C(t *testing.B) {
	if deutchjozsa.DeutchJozsa(8, deutchjozsa.ConstantFunc) != IS_CONSTANT {
		t.Fail()
	}
}

func BenchmarkDeutchJozsa16B(t *testing.B) {
	if deutchjozsa.DeutchJozsa(16, deutchjozsa.BalancedFunc) != IS_BALANCED {
		t.Fail()
	}
}

func BenchmarkDeutchJozsa16C(t *testing.B) {
	if deutchjozsa.DeutchJozsa(16, deutchjozsa.ConstantFunc) != IS_CONSTANT {
		t.Fail()
	}
}

func BenchmarkDeutchJozsaClassical16B(t *testing.B) {
	if deutchjozsa.DeutchJozsaClassical(16, deutchjozsa.BalancedFunc) != IS_BALANCED {
		t.Fail()
	}
}

func BenchmarkDeutchJozsaClassical16C(t *testing.B) {
	if deutchjozsa.DeutchJozsaClassical(16, deutchjozsa.ConstantFunc) != IS_CONSTANT {
		t.Fail()
	}
}
