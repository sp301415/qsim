package qsim_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/sp301415/qsim"
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

	fmt.Println(c.StateToString())

	q := vector.New([]complex128{complex(1/math.Sqrt(2), 0), 0, 0, complex(1/math.Sqrt(2), 0)})

	if !q.Equals(c.State) {
		t.Fail()
	}
}

func TestQFT(t *testing.T) {
	c := qsim.NewCircuit(2)
	c.InitQbit(qbit.NewFromCbit(1, 2))
	c.QFT(0, 2)

	fmt.Println(c.StateToString())
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
