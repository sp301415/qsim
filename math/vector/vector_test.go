package vector_test

import (
	"testing"

	"github.com/sp301415/qsim/math/matrix"
	"github.com/sp301415/qsim/math/vector"
)

func TestApply(t *testing.T) {
	m := matrix.New(
		[][]complex128{
			{1, 1, 1},
			{2, 2, 2},
		},
	)
	v := vector.New(
		[]complex128{
			2, 2, 2,
		},
	)
	w := vector.New(
		[]complex128{
			6, 12,
		},
	)

	if !w.Equals(v.Apply(m)) {
		t.Fail()
	}
}

func TestTensor(t *testing.T) {
	v := vector.New(
		[]complex128{
			2, 2, 2,
		},
	)
	w := vector.New(
		[]complex128{
			4, 4, 4, 4, 4, 4, 4, 4, 4,
		},
	)

	if !w.Equals(v.Tensor(v)) {
		t.Fail()
	}
}
