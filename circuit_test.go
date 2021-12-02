package qsim_test

import (
	"fmt"
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
