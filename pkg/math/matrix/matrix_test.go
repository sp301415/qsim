package matrix_test

import (
	"testing"

	"github.com/sp301415/qsim/pkg/math/matrix"
)

func TestAdd(t *testing.T) {
	m1 := matrix.New(
		[][]complex128{
			{1, 2},
			{3, 4},
		},
	)
	m2 := matrix.New(
		[][]complex128{
			{2, 4},
			{6, 8},
		},
	)
	m3 := matrix.New(
		[][]complex128{
			{3, 6},
			{9, 12},
		},
	)

	if !m3.Equals(matrix.Add(m1, m2)) {
		t.Fail()
	}
}

func TestMul(t *testing.T) {
	m1 := matrix.New(
		[][]complex128{
			{1, 2},
			{3, 4},
		},
	)
	m2 := matrix.New(
		[][]complex128{
			{2, 4},
			{6, 8},
		},
	)
	m3 := matrix.New(
		[][]complex128{
			{14, 20},
			{30, 44},
		},
	)

	if !m3.Equals(matrix.Mul(m1, m2)) {
		t.Fail()
	}
}
