package matrix_test

import (
	"testing"

	"github.com/sp301415/qsim/math/matrix"
)

func TestTensor(t *testing.T) {
	m := matrix.New(
		[][]complex128{
			{1, 0},
			{0, 1},
		},
	)
	n := matrix.New(
		[][]complex128{
			{0, 1},
			{1, 0},
		},
	)
	mn := matrix.New(
		[][]complex128{
			{0, 1, 0, 0},
			{1, 0, 0, 0},
			{0, 0, 0, 1},
			{0, 0, 1, 0},
		},
	)

	if !mn.Equals(m.Tensor(n)) {
		t.Fail()
	}
}
