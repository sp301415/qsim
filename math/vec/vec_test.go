package vec_test

import (
	"testing"

	"github.com/sp301415/qsim/math/vec"
)

func TestAlloc(t *testing.T) {
	d := []complex128{1, 2, 3}

	v := vec.NewVec(3)
	for i := range d {
		v[i] = d[i]
	}

	w := vec.NewVecSlice(d)
	if !v.Equals(w) {
		t.Fail()
	}

	x := vec.NewVecVars(1, 2, 3)
	if !v.Equals(x) {
		t.Fail()
	}
}

func TestEqual(t *testing.T) {
	v1 := vec.NewVecVars(1, 2, 3)
	v2 := vec.NewVecVars(1, 2, 4)
	v3 := vec.NewVecVars(1, 2, 3, 4)

	if v1.Equals(v2) {
		t.Fail()
	}
	if v2.Equals(v3) {
		t.Fail()
	}
	if v3.Equals(v1) {
		t.Fail()
	}
}

func TestCopy(t *testing.T) {
	v := vec.NewVecVars(1, 2, 3)
	w := v.Copy()
	v[0] = 2

	if v.Equals(w) {
		t.Fail()
	}
}

func TestAdd(t *testing.T) {
	v1 := vec.NewVecVars(1, 2, 3)
	v2 := vec.NewVecVars(3, 2, 1)
	v3 := vec.NewVecVars(4, 4, 4)

	if !v1.Add(v2).Equals(v3) {
		t.Fail()
	}
}

func TestSub(t *testing.T) {
	v1 := vec.NewVecVars(1, 2, 3)
	v2 := vec.NewVecVars(3, 2, 1)
	v3 := vec.NewVecVars(4, 4, 4)

	if !v3.Sub(v2).Equals(v1) {
		t.Fail()
	}

}

func TestScalarMul(t *testing.T) {
	v := vec.NewVecVars(1, 2, 3)
	a := 3 + 0i
	av := vec.NewVecVars(3, 6, 9)

	if !v.ScalarMul(a).Equals(av) {
		t.Fail()
	}
}

func TestDot(t *testing.T) {
	v := vec.NewVecVars(1+1i, 2+2i, 3+3i)

	if v.Dot(v) != 28+0i {
		t.Fail()
	}
}

func TestTensor(t *testing.T) {
	v := vec.NewVecVars(1, 2)
	w := vec.NewVecVars(1, 2, 2, 4)

	if !v.Tensor(v).Equals(w) {
		t.Fail()
	}
}

func TestNorm(t *testing.T) {
	a := 1 + 1i

	v := vec.NewVecSlice([]complex128{a, a})

	if v.Norm() != 2 || v.NormSquared() != 4 {
		t.Fail()
	}
}
