// Package qsim provides functions for a quantum circuit.
package qsim

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"strings"
	"time"

	"github.com/sp301415/qsim/math/matrix"
	"github.com/sp301415/qsim/math/numbers"
	"github.com/sp301415/qsim/math/vector"
	"github.com/sp301415/qsim/quantum/gate"
	"github.com/sp301415/qsim/quantum/qbit"
)

type Circuit struct {
	// Why save gates when we don't need diagram? :)
	N     int
	State vector.Vector
}

func NewCircuit(n int) Circuit {
	q := qbit.Zeros(n)

	return Circuit{N: n, State: q}
}

func (circ *Circuit) InitQbit(q vector.Vector) {
	circ.State = q
}

func (circ *Circuit) Apply(operator matrix.Matrix, qbits ...int) {
	if !operator.IsUnitary() {
		panic("Operator must be unitary.")
	}

	if len(operator) != 1<<len(qbits) {
		panic("Operator size does not match with input qbits.")
	}

	// Tensor Product takes too long, we need another method.
	// Idea: decompose state vector to pure states?

	res := make(vector.Vector, circ.State.Dim())

	for n, a := range circ.State {
		// a * |n>
		// First, extract input qbits from n
		// For example, if n = |0101> and qbit = 0, 2, x = |11>
		x := 0
		for i, v := range qbits {
			// Extract vth bit from n, plug it in to ith bit of x.
			x += ((n >> v) % 2) * (1 << i)
		}
		// Generate new qbit from x, apply operator to it.
		q := qbit.NewFromCbit(x, len(qbits))
		q = q.Apply(operator)

		for qx, qa := range q {
			// Extract ith bit from x, plug it in to vth bit of n.
			nn := 0
			for i, v := range qbits {
				qi := (qx >> i) % 2
				nn = n | (1 << v)
				if qi == 0 {
					nn = nn - (1 << v)
				}
			}
			res[nn] += a * qa
		}
	}

	circ.State = res
}

func (circ *Circuit) I(n int) {
	// Just Do Nothing. lol.
}

func (circ *Circuit) X(n int) {
	circ.Apply(gate.X(), n)
}

func (circ *Circuit) Y(n int) {
	circ.Apply(gate.Y(), n)
}

func (circ *Circuit) Z(n int) {
	circ.Apply(gate.Z(), n)
}

func (circ *Circuit) H(n int) {
	circ.Apply(gate.H(), n)
}

func (circ *Circuit) ControlOperator(operator matrix.Matrix, c int, x int) {
	if !operator.IsUnitary() {
		panic("Operator not unitary.")
	}

	if len(operator) != 2 {
		panic("Currently only one-qbit operator allowed.")
	}

	res := make(vector.Vector, circ.State.Dim())

	for n, a := range circ.State {
		cc := (n >> c) % 2

		if cc == 0 {
			res[n] += a
			continue
		}

		qx := qbit.NewFromCbit((n>>x)%2, 1)
		qx = qx.Apply(operator)
		nn := n | (1 << x)

		res[nn] += qx[1] * a
		res[nn-(1<<x)] += qx[0] * a
	}

	circ.State = res

}

func (circ *Circuit) CX(c int, x int) {
	circ.ControlOperator(gate.X(), c, x)
}

func (circ *Circuit) Measure(x int) int {
	prob := [2]float64{0, 0}

	for n, a := range circ.State {
		prob[(n>>x)%2] += math.Pow(cmplx.Abs(a), 2)
	}

	// Wait, Golang does not have weighted sampling? WTF.
	rand.Seed(time.Now().UnixNano())
	rand := rand.Float64()

	m := 0
	if rand < prob[0] {
		m = 0
	} else {
		m = 1
	}

	for n := range circ.State {
		if (n>>x)%2 != m {
			circ.State[n] = 0
		}
	}

	circ.State = circ.State.Normalize()

	return m
}

func (circ Circuit) StateToString() string {
	q := circ.State

	qs := make([]string, 0)
	d := numbers.Log2(q.Dim())

	for i, v := range q {
		if v == 0 {
			continue
		}

		res := ""

		if v == 1.0 {
			// Do Nothing
		} else if v == 1.0i {
			res += "i"
		} else if real(v) == 0.0 {
			res += fmt.Sprint(imag(v)) + "i"
		} else if imag(v) == 0.0 {
			res += fmt.Sprint(real(v))
		} else {
			res += fmt.Sprint(v)
		}

		qs = append(qs, res+fmt.Sprintf("|%0*b>", d, i))
	}

	return strings.Join(qs, " + ")
}
