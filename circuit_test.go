package qsim_test

import (
	"math/rand"
	"testing"

	"github.com/sp301415/qsim"
	"github.com/sp301415/qsim/math/matrix"
	"github.com/sp301415/qsim/math/vector"
	"github.com/sp301415/qsim/quantum/gate"
	"github.com/sp301415/qsim/quantum/qbit"
)

func TestInitC(t *testing.T) {
	c := qsim.NewCircuit(2)
	c.InitCbit(3)

	q := vector.New([]complex128{0, 0, 0, 1})

	if !c.State.Equals(q) {
		t.Fail()
	}
}

func TestSingleX(t *testing.T) {
	c := qsim.NewCircuit(1)
	c.X(0)

	q := qbit.NewFromCbit(1, 1)

	if !c.State.Equals(q) {
		t.Fail()
	}
}

func TestSingleH(t *testing.T) {
	c := qsim.NewCircuit(1)
	c.H(0)

	q := vector.New([]complex128{1, 1}).Normalize()

	if !c.State.Equals(q) {
		t.Fail()
	}
}

func TestMultiX(t *testing.T) {
	N := 10
	c := qsim.NewCircuit(N)

	for i := 0; i < N; i++ {
		c.X(i)
	}

	q1 := qbit.Ones(N)

	if !c.State.Equals(q1) {
		t.Fail()
	}

	for i := N - 1; i >= 0; i-- {
		c.X(i)
	}

	q2 := qbit.Zeros(N)

	if !c.State.Equals(q2) {
		t.Fail()
	}
}

func TestMultiH(t *testing.T) {
	N := 10
	c := qsim.NewCircuit(N)

	for i := 0; i < N; i++ {
		c.H(i)
	}

	q1 := qbit.Zeros(N)
	for i := range q1 {
		q1[i] = 1
	}
	q1 = q1.Normalize()

	if !c.State.Equals(q1) {
		t.Fail()
	}

	for i := N - 1; i >= 0; i-- {
		c.H(i)
	}

	q2 := qbit.Zeros(N)

	if !c.State.Equals(q2) {
		t.Fail()
	}
}

func TestMultiApply(t *testing.T) {
	N := 5
	c1 := qsim.NewCircuit(N)
	c2 := qsim.NewCircuit(N)
	Hs := gate.H()
	regs := make([]int, N)

	for i := 0; i < N; i++ {
		c1.H(i)
		regs[i] = i

		if i > 0 {
			Hs = Hs.Tensor(gate.H())
		}
	}

	c2.Apply(Hs, regs...)

	if !c1.State.Equals(c2.State) {
		t.Fail()
	}
}

func TestHH(t *testing.T) {
	c1 := qsim.NewCircuit(3)
	c2 := qsim.NewCircuit(3)

	c1.H(0)
	c1.H(1)

	c2.Apply(gate.H().Tensor(gate.H()), 0, 1)

	if !c1.State.Equals(c2.State) {
		t.Fail()
	}
}

func TestOracle(t *testing.T) {
	c := qsim.NewCircuit(2)

	c.H(0)
	c.CX(0, 1)
	c.ApplyOracle(func(_ int) int { return 1 }, []int{1}, []int{0})

	q := vector.Zeros(1 << 2)
	q[0b01] = 1
	q[0b10] = 1
	q = q.Normalize()

	if !c.State.Equals(q) {
		t.Fail()
	}
}

func TestSingleCX(t *testing.T) {
	c := qsim.NewCircuit(2)
	c.X(0)
	c.Control(gate.X(), []int{0}, []int{1})

	if !c.State.Equals(qbit.NewFromCbit(0b11, 2)) {
		t.Fail()
	}

	c.Control(gate.X(), []int{0}, []int{1})

	if !c.State.Equals(qbit.NewFromCbit(0b01, 2)) {
		t.Fail()
	}
}

func TestSingleCH(t *testing.T) {
	c1 := qsim.NewCircuit(2)
	c2 := qsim.NewCircuit(2)

	c1.X(0)
	c2.X(0)

	c1.Control(gate.H(), []int{0}, []int{1})
	c2.H(1)

	if !c1.State.Equals(c2.State) {
		t.Fail()
	}
}

func TestMultiCCX(t *testing.T) {
	c := qsim.NewCircuit(3)

	c.H(0)
	c.CCX(0, 1, 2)

	q := vector.Zeros(1 << 3)
	q[0b000] = 1
	q[0b101] = 1
	q = q.Normalize()

	if !c.State.Equals(q) {
		t.Fail()
	}
}

func TestCHH(t *testing.T) {
	c := qsim.NewCircuit(3)

	c.X(0)
	c.Control(matrix.Tensor(gate.H(), gate.H()), []int{0}, []int{1, 2})

	q := vector.Zeros(1 << 3)
	q[0b001] = 1
	q[0b011] = 1
	q[0b101] = 1
	q[0b111] = 1

	q = q.Normalize()

	if !c.State.Equals(q) {
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
		if c.State[0] == 0 {
			t.Fail()
		}
	} else {
		if c.State[len(c.State)-1] == 0 {
			t.Fail()
		}
	}
}

func TestQFT(t *testing.T) {
	c := qsim.NewCircuit(2)
	c.X(0)
	c.QFT(0, 2)

	q := vector.New([]complex128{0.5, 0.5i, -0.5, -0.5i})

	if !c.State.Equals(q) {
		t.Fail()
	}
}

func TestInvQFT(t *testing.T) {
	N := 10
	c := qsim.NewCircuit(N)

	c.QFT(0, N)
	c.InvQFT(0, N)

	if !c.State.Equals(qbit.Zeros(N)) {
		t.Fail()
	}
}

func TestMeasure(t *testing.T) {
	N := 10
	M := rand.Intn(1 << N)

	c := qsim.NewCircuit(N)
	c.InitCbit(M)

	regs := make([]int, N)
	for i := range regs {
		regs[i] = i
	}

	m := c.Measure(regs...)

	if m != M {
		t.Fail()
	}
}

func BenchmarkApplywithOptimization(t *testing.B) {
	// Hadamard gate for every qbit.

	N := 8
	c := qsim.NewCircuit(N)

	for i := 0; i < N; i++ {
		c.H(i)
	}

	for i := N - 1; i >= 0; i-- {
		c.H(i)
	}
}

func BenchmarkApplywithoutOptimization(t *testing.B) {
	// Hadamard gate for every qbit.
	// But with Tensor Product.

	N := 8
	c := qsim.NewCircuit(N)
	Hs := gate.H()
	regs := make([]int, N)

	for i := 0; i < N; i++ {
		regs[i] = i

		if i > 0 {
			Hs = Hs.Tensor(gate.H())
		}
	}

	c.Apply(Hs, regs...)
	c.Apply(Hs, regs...)
}
