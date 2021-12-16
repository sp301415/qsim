package mat_test

import (
	"math"
	"testing"

	"github.com/sp301415/qsim/math/mat"
	"github.com/sp301415/qsim/math/vec"
)

func TestAlloc(t *testing.T) {
	d := [][]complex128{{1, 0}, {0, 1}}
	dflat := []complex128{1, 0, 0, 1}

	v1 := mat.NewMat(2, 2)
	v1[0][0], v1[1][1] = 1, 1
	v2 := mat.NewMatSlice(d)
	v3 := mat.NewMatSliceVars(d...)
	v4 := mat.NewMatVars(2, dflat...)
	v5 := mat.NewMatVecs([]vec.Vec{d[0], d[1]})
	v6 := mat.NewSquare(2)
	v6[0][0], v6[1][1] = 1, 1
	v7 := mat.NewId(2)

	if !v1.Equals(v2) {
		t.Fail()
	}
	if !v2.Equals(v3) {
		t.Fail()
	}
	if !v3.Equals(v4) {
		t.Fail()
	}
	if !v4.Equals(v5) {
		t.Fail()
	}
	if !v5.Equals(v6) {
		t.Fail()
	}
	if !v6.Equals(v7) {
		t.Fail()
	}
	if !v7.Equals(v1) {
		t.Fail()
	}
}

func TestEqual(t *testing.T) {
	m1 := mat.NewMatVars(2, 1, 2, 3, 4)
	m2 := mat.NewMatVars(2, 1, 2, 3, 5)
	m3 := mat.NewMatVars(1, 1, 2, 3, 4)

	if m1.Equals(m2) {
		t.Fail()
	}
	if m3.Equals(m1) {
		t.Fail()
	}
}

func TestCopy(t *testing.T) {
	m := mat.NewMatVars(2, 1, 2, 3, 4)
	n := m.Copy()
	n[0][0] = 1

	if !m.Equals(n) {
		t.Fail()
	}
}

func TestSquareUnitary(t *testing.T) {
	k := complex(math.Sqrt(2)/2.0, 0)
	m := mat.NewMatVars(2, k, k, k, -k)

	if !(m.IsSquare() && m.IsUnitary()) {
		t.Fail()
	}
}

func TestGetSetRowCol(t *testing.T) {
	m1 := mat.NewMatVars(2, 1, 0, 2, 2)
	m2 := mat.NewMatVars(2, 2, 0, 2, 1)
	I1 := mat.NewId(2)
	I2 := mat.NewId(2)

	v := vec.NewVecVars(2, 2)

	if !v.Equals(m1.GetRow(1)) {
		t.Fail()
	}

	if !v.Equals(m2.GetCol(0)) {
		t.Fail()
	}

	I1.SetRow(1, v)
	I2.SetCol(0, v)

	if !I1.Equals(m1) {
		t.Fail()
	}

	if !I2.Equals(m2) {
		t.Fail()
	}
}

func TestDim(t *testing.T) {
	m := mat.NewMat(4, 8)

	if !(m.NRows() == 4 && m.NCols() == 8) {
		t.Fail()
	}

	r, c := m.Dim()
	if !(r == 4 && c == 8) {
		t.Fail()
	}
}

func TestAdd(t *testing.T) {
	m1 := mat.NewMatVars(2, 1, 2, 3, 4)
	m2 := mat.NewMatVars(2, 4, 3, 2, 1)
	m3 := mat.NewMatVars(2, 5, 5, 5, 5)

	if !m1.Add(m2).Equals(m3) {
		t.Fail()
	}
}

func TestSub(t *testing.T) {
	m1 := mat.NewMatVars(2, 1, 2, 3, 4)
	m2 := mat.NewMatVars(2, 4, 3, 2, 1)
	m3 := mat.NewMatVars(2, 5, 5, 5, 5)

	if !m3.Sub(m2).Equals(m1) {
		t.Fail()
	}
}

func TestScalarMul(t *testing.T) {
	m := mat.NewMatVars(2, 1, 2, 3, 4)
	a := 3 + 2i
	am := mat.NewMatVars(2, 3+2i, 6+4i, 9+6i, 12+8i)

	if !m.ScalarMul(a).Equals(am) {
		t.Fail()
	}
}

func TestTensor(t *testing.T) {
	m1 := mat.NewMatSlice(
		[][]complex128{
			{0, 1},
			{1, 0},
		},
	)
	m2 := mat.NewMatSlice(
		[][]complex128{
			{1, 0},
			{0, 1},
		},
	)
	m3 := mat.NewMatSlice(
		[][]complex128{
			{0, 0, 1, 0},
			{0, 0, 0, 1},
			{1, 0, 0, 0},
			{0, 1, 0, 0},
		},
	)

	if !m1.Tensor(m2).Equals(m3) {
		t.Fail()
	}
}

func TestDagger(t *testing.T) {
	m1 := mat.NewMatSlice(
		[][]complex128{
			{1 + 2i, 2 + 3i},
			{3 + 4i, 4 + 5i},
		},
	)
	m2 := mat.NewMatSlice(
		[][]complex128{
			{1 - 2i, 3 - 4i},
			{2 - 3i, 4 - 5i},
		},
	)

	if !m1.Dagger().Equals(m2) {
		t.Fail()
	}
}
