package qsim_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/sp301415/qsim"
	"github.com/sp301415/qsim/math/vec"
	"github.com/sp301415/qsim/utils/slice"
)

func TestInitC(t *testing.T) {
	c := qsim.NewCircuit(2)
	c.SetBit(3)

	q := qsim.NewQubit(vec.NewVecSlice([]complex128{0, 0, 0, 1}))

	if !c.State().Equals(q) {
		t.Fail()
	}
}

func TestSingleX(t *testing.T) {
	c := qsim.NewCircuit(1)
	c.X(0)

	q := qsim.NewBit(1, 1)

	if !c.State().Equals(q) {
		t.Fail()
	}
}

func TestSingleH(t *testing.T) {
	c := qsim.NewCircuit(1)
	c.H(0)

	q := qsim.NewQubit(vec.NewVecSlice([]complex128{math.Sqrt2 / 2.0, math.Sqrt2 / 2.0}))

	if !c.State().Equals(q) {
		t.Fail()
	}
}

func TestMultiX(t *testing.T) {
	N := 10
	c := qsim.NewCircuit(N)

	for i := 0; i < N; i++ {
		c.X(i)
	}

	q1 := qsim.NewBit(1<<N-1, N)

	if !c.State().Equals(q1) {
		t.Fail()
	}

	for i := N - 1; i >= 0; i-- {
		c.X(i)
	}

	q2 := qsim.NewBit(0, N)

	if !c.State().Equals(q2) {
		t.Fail()
	}
}

func TestMultiH(t *testing.T) {
	N := 10
	c := qsim.NewCircuit(N)

	for i := 0; i < N; i++ {
		c.H(i)
	}

	q1 := qsim.NewBit(0, N).ToVec()
	for i := range q1 {
		q1[i] = complex(1.0/math.Pow(2.0, float64(N)/2.0), 0)
	}

	if !c.State().Equals(qsim.NewQubit(q1)) {
		t.Fail()
	}

	for i := N - 1; i >= 0; i-- {
		c.H(i)
	}

	q2 := qsim.NewBit(0, N)

	if !c.State().Equals(q2) {
		t.Fail()
	}
}

func TestMultiApply(t *testing.T) {
	N := 5
	regs := slice.Sequence(0, N)

	c1 := qsim.NewCircuit(N)
	c1.H(regs...)

	c2 := qsim.NewCircuit(N)
	Hs := qsim.H()
	for i := 1; i < N; i++ {
		Hs = Hs.Tensor(qsim.H())
	}
	c2.Apply(Hs, regs...)

	if !c1.State().Equals(c2.State()) {
		t.Fail()
	}
}

func TestHH(t *testing.T) {
	c1 := qsim.NewCircuit(3)
	c2 := qsim.NewCircuit(3)

	c1.H(0)
	c1.H(1)

	c2.Apply(qsim.H().Tensor(qsim.H()), 0, 1)

	if !c1.State().Equals(c2.State()) {
		t.Fail()
	}
}

func TestOracle(t *testing.T) {
	c := qsim.NewCircuit(2)

	c.H(0)
	c.CX(0, 1)

	c.ApplyOracle(func(_ int) int { return 1 }, []int{1}, []int{0})

	q := vec.NewVec(1 << 2)
	q[0b01] = math.Sqrt2 / 2.0
	q[0b10] = math.Sqrt2 / 2.0

	if !c.State().Equals(qsim.NewQubit(q)) {
		t.Fail()
	}
}

func TestSingleCX(t *testing.T) {
	c := qsim.NewCircuit(2)
	c.X(0)
	c.Control(qsim.X(), []int{0}, []int{1})

	if !c.State().Equals(qsim.NewBit(0b11, 2)) {
		t.Fail()
	}

	c.Control(qsim.X(), []int{0}, []int{1})

	if !c.State().Equals(qsim.NewBit(0b01, 2)) {
		t.Fail()
	}
}

func TestSingleCH(t *testing.T) {
	c1 := qsim.NewCircuit(2)
	c2 := qsim.NewCircuit(2)

	c1.X(0)
	c2.X(0)

	c1.Control(qsim.H(), []int{0}, []int{1})
	c2.H(1)

	if !c1.State().Equals(c2.State()) {
		t.Fail()
	}
}

func TestMultiCCX(t *testing.T) {
	c := qsim.NewCircuit(3)

	c.H(0)
	c.H(1)
	c.CCX(0, 1, 2)

	q := vec.NewVec(1 << 3)
	q[0b000] = 0.5
	q[0b001] = 0.5
	q[0b010] = 0.5
	q[0b111] = 0.5

	if !c.State().Equals(qsim.NewQubit(q)) {
		t.Fail()
	}
}

func TestCHH(t *testing.T) {
	c := qsim.NewCircuit(3)

	c.X(0)
	c.Control(qsim.H().Tensor(qsim.H()), []int{0}, []int{1, 2})

	q := vec.NewVec(1 << 3)
	q[0b001] = 0.5
	q[0b011] = 0.5
	q[0b101] = 0.5
	q[0b111] = 0.5

	if !c.State().Equals(qsim.NewQubit(q)) {
		t.Fail()
	}
}

func TestEntangle(t *testing.T) {
	N := 10
	c := qsim.NewCircuit(N)

	c.H(0)

	for i := 0; i < N-1; i++ {
		c.CX(i, i+1)
	}

	m := c.Measure(0)

	if m == 0 {
		if c.State().ToVec()[0] == 0 {
			t.Fail()
		}
	} else {
		if c.State().ToVec()[c.State().Dim()-1] == 0 {
			t.Fail()
		}
	}
}

func TestQFT(t *testing.T) {
	c := qsim.NewCircuit(2)
	c.X(0)
	c.QFT(0, 1)

	q := qsim.NewQubit([]complex128{0.5, 0.5i, -0.5, -0.5i})

	if !c.State().Equals(q) {
		t.Fail()
	}
}

func TestInvQFT(t *testing.T) {
	N := 10
	c := qsim.NewCircuit(N)
	iregs := slice.Sequence(0, N)

	c.QFT(iregs...)
	c.InvQFT(iregs...)

	if !c.State().Equals(qsim.NewBit(0, N)) {
		t.Fail()
	}
}

func TestMeasure(t *testing.T) {
	N := 10
	M := rand.Intn(1 << N)

	c := qsim.NewCircuit(N)
	c.SetBit(M)

	m := c.Measure(slice.Sequence(0, N)...)

	if m != M {
		t.Fail()
	}
}

func BenchmarkTensorApply(t *testing.B) {
	N := 10
	regs := slice.Sequence(0, N)

	H := qsim.H()
	X := qsim.X()
	Z := qsim.Z()
	T := qsim.T()

	c := qsim.NewCircuit(N)
	for i := 1; i < N; i++ {
		H = H.Tensor(qsim.H())
		X = X.Tensor(qsim.X())
		Z = Z.Tensor(qsim.Z())
		T = T.Tensor(qsim.T())
	}

	c.Apply(H, regs...)
	c.Apply(X, regs...)
	c.Apply(Z, regs...)
	c.Apply(T, regs...)
}

func BenchmarkApply(t *testing.B) {
	N := 10

	c := qsim.NewCircuit(N)
	c.Option.PARALLEL_THRESHOLD = 20

	for i := 0; i < N; i++ {
		c.H(i)
	}
	for i := 0; i < N; i++ {
		c.X(i)
	}
	for i := 0; i < N; i++ {
		c.Z(i)
	}
	for i := 0; i < N; i++ {
		c.T(i)
	}
}

func BenchmarkApplyParallel(t *testing.B) {
	N := 10

	c := qsim.NewCircuit(N)
	c.Option.PARALLEL_THRESHOLD = 5

	for i := 0; i < N; i++ {
		c.H(i)
	}
	for i := 0; i < N; i++ {
		c.X(i)
	}
	for i := 0; i < N; i++ {
		c.Z(i)
	}
	for i := 0; i < N; i++ {
		c.T(i)
	}
}

func BenchmarkGrover4(t *testing.B) {
	q := qsim.NewCircuit(4)

	// superposition
	for i := 0; i < q.Size(); i++ {
		q.H(i)
	}

	// iteration
	N := 1 << q.Size()
	r := math.Floor(math.Pi / 4 * math.Sqrt(float64(N)))
	for i := 0; i < int(r); i++ {
		q.X(0, 1)
		q.H(0)
		q.Control(qsim.X(), []int{1, 2, 3}, []int{0})
		q.H(0)
		q.X(0, 1)

		q.H(0, 1, 2, 3)
		q.X(0, 1, 2, 3)
		q.H(0)
		q.Control(qsim.X(), []int{1, 2, 3}, []int{0})
		q.H(0)
		q.X(0, 1, 2, 3)
		q.H(0, 1, 2, 3)
	}

	if q.Measure(0, 1, 2, 3) != 12 {
		t.Fail()
	}
}
