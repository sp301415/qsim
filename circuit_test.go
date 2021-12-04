package qsim_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/sp301415/qsim"
	"github.com/sp301415/qsim/math/matrix"
	"github.com/sp301415/qsim/math/vector"
	"github.com/sp301415/qsim/quantum/gate"
	"github.com/sp301415/qsim/quantum/qbit"
)

func TestInitQbits(t *testing.T) {
	c := qsim.NewCircuit(10)
	q := qbit.NewFromCbit(37, 10)
	c.InitQbit(q)

	fmt.Println(c.StateToString())

	if !c.State.Equals(q) {
		t.Fail()
	}
}

func TestX(t *testing.T) {
	c := qsim.NewCircuit(10)

	for i := 0; i < 10; i++ {
		c.X(i)
	}

	if !c.State.Equals(qbit.NewFromCbit((1<<10)-1, 10)) {
		t.Fail()
	}
}

func BenchmarkApply(t *testing.B) {
	N := 16
	c := qsim.NewCircuit(N)

	for i := 0; i < N; i++ {
		c.Apply(gate.X(), i)
	}

	for i := 0; i < N; i++ {
		c.Apply(gate.X(), i)
	}

	if !c.State.Equals(qbit.NewFromCbit(0, N)) {
		t.Fail()
	}
}

func BenchmarkApplyNonPurewOp(t *testing.B) {
	N := 8
	c := qsim.NewCircuit(N)

	t.StartTimer()
	for i := 0; i < N; i++ {
		c.H(i)
	}
}

func BenchmarkApplyNonPurewoOp(t *testing.B) {
	N := 8
	c := qsim.NewCircuit(N)
	Hs := make([]matrix.Matrix, N)
	idx := make([]int, N)

	for i := range Hs {
		Hs[i] = gate.H()
		idx[i] = i
	}

	t.StartTimer()

	H := matrix.Tensor(Hs...)
	c.Apply(H, idx...)
}

func TestApplyNonEntagled(t *testing.T) {
	c := qsim.NewCircuit(2)

	c.Apply(gate.H(), 0)
	c.Apply(gate.H(), 1)

	q := vector.New([]complex128{0.5, 0.5, 0.5, 0.5})

	if !q.Equals(c.State) {
		t.Fail()
	}
}

func TestApplyEntangled(t *testing.T) {
	c := qsim.NewCircuit(2)
	c.H(0)
	c.CX(0, 1)
	q := vector.New([]complex128{complex(1/math.Sqrt(2), 0), 0, 0, complex(1/math.Sqrt(2), 0)})

	if !q.Equals(c.State) {
		t.Fail()
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
func TestQFTInvQFT(t *testing.T) {
	c := qsim.NewCircuit(3)
	c.InitQbit(qbit.NewFromCbit(3, 3))

	c.QFT(0, 3)
	c.InvQFT(0, 3)

	if !c.State.Equals(qbit.NewFromCbit(3, 3)) {
		t.Fail()
	}
}

func BenchmarkQFT(t *testing.B) {
	c := qsim.NewCircuit(16)
	c.QFT(0, 16)
}

func TestMeasure(t *testing.T) {
	c := qsim.NewCircuit(2)

	c.H(0)
	c.CX(0, 1)

	r1 := qbit.NewFromCbit(0b11, 2)
	r2 := qbit.NewFromCbit(0b00, 2)

	c.Measure(0)

	if !c.State.Equals(r1) && !c.State.Equals(r2) {
		t.Fail()
	}
}

func TestMultiMeasure(t *testing.T) {
	N := 3
	I := 5
	c := qsim.NewCircuit(N)
	c.InitQbit(qbit.NewFromCbit(I, N))

	x := c.Measure(0, 2)

	if x != 3 {
		t.Fail()
	}
}
