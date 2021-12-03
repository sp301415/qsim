// Package qsim provides functions for a quantum circuit.
package qsim

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"sort"
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

	res := vector.Zeros(circ.State.Dim())

	for n, a := range circ.State {
		if a == 0 {
			continue
		}
		// a * |n>
		// First, extract input qbits from n
		// For example, if n = |0101> and qbit = 0, 2, x = |11>
		x := 0
		for i, v := range qbits {
			// Extract vth bit from n, plug it in to ith bit of x.
			x += ((n >> v) % 2) << i
		}
		// Generate new qbit from x, apply operator to it.
		q := qbit.NewFromCbit(x, len(qbits))
		q = q.Apply(operator)

		for qx, qa := range q {
			// Extract ith bit from x, plug it in to vth bit of n.
			nn := 0
			for i, v := range qbits {
				qi := (qx >> i) % 2
				nn = (n | (1 << v)) - ((qi ^ 1) << v)
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

func (circ *Circuit) ControlOperator(operator matrix.Matrix, cs []int, xs []int) {
	if !operator.IsUnitary() {
		panic("Operator not unitary.")
	}

	if len(operator) != len(xs) {
		panic("Operator size does not match.")
	}

	if len(cs) == 1 && len(xs) == 1 {
		circ.ControlSingle(operator, cs[0], xs[0])
		return
	}

	// This is almost same to Apply(), except it checks the control qbit.

	res := vector.Zeros(len(circ.State))

	for n, a := range circ.State {
		if a == 0 {
			continue
		}

		ctrl := 0
		for _, c := range cs {
			ctrl ^= ((n >> c) % 2)
		}

		if ctrl == 0 {
			res[n] += a
			continue
		}

		x := 0
		for i, v := range xs {
			x += ((n >> v) % 2) << i
		}
		q := qbit.NewFromCbit(x, len(xs))
		q = q.Apply(operator)

		for qx, qa := range q {
			nn := 0
			for i, v := range xs {
				qi := (qx >> i) % 2
				nn = (n | (1 << v)) - ((qi ^ 1) << v)
			}
			res[nn] += a * qa
		}
	}

	circ.State = res
}

// A more efficient implementation, when control/output qbit is single.
func (circ *Circuit) ControlSingle(operator matrix.Matrix, c int, x int) {
	res := vector.Zeros(circ.State.Dim())

	for n, a := range circ.State {
		if a == 0 {
			continue
		}

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
	circ.ControlSingle(gate.X(), c, x)
}

func (circ *Circuit) Swap(x int, y int) {
	res := vector.Zeros(len(circ.State))

	for n, a := range circ.State {
		bx := (n >> x) % 2
		by := (n >> y) % 2

		nn := n

		nn = (nn | (1 << x)) - ((by ^ 1) << x)
		nn = (nn | (1 << y)) - ((bx ^ 1) << y)

		res[nn] = a
	}

	circ.State = res
}

// QFT from [start, end)
func (circ *Circuit) QFT(start, end int) {
	if start < 0 || end > circ.N {
		panic("Index out of range.")
	}

	if start >= end {
		panic("Invalid start / end parameters.")
	}

	for i := end - 1; i >= start; i-- {
		circ.H(i)

		for j := start; j < i; j++ {
			circ.ControlSingle(gate.P(math.Pi/float64(numbers.Pow(2, i-j))), j, i)
		}
	}

	for i, j := start, end-1; i < j; i, j = i+1, j-1 {
		circ.Swap(i, j)
	}

}

func (circ *Circuit) InvQFT(start, end int) {
	if start < 0 || end > circ.N {
		panic("Index out of range.")
	}

	if start >= end {
		panic("Invalid start / end parameters.")
	}

	for i, j := start, end-1; i < j; i, j = i+1, j-1 {
		circ.Swap(i, j)
	}

	for i := start; i < end; i++ {
		for j := start; j < i; j++ {
			circ.ControlSingle(gate.P(-math.Pi/float64(numbers.Pow(2, i-j))), j, i)
		}
		circ.H(i)
	}
}

func (circ *Circuit) Measure(qbits ...int) int {
	qbits = sort.IntSlice(qbits)

	if qbits[0] < 0 || qbits[len(qbits)-1] > circ.N-1 {
		panic("Invalid registers.")
	}

	prob := make([]float64, 1<<len(qbits))

	for n, a := range circ.State {
		if a == 0 {
			continue
		}
		o := 0
		for i, q := range qbits {
			o += ((n >> q) % 2) << i
		}
		prob[o] += math.Pow(cmplx.Abs(a), 2.0)
	}

	// Wait, Golang does not have weighted sampling? WTF.
	rand.Seed(time.Now().UnixNano())
	rand := rand.Float64()

	output := 0
	accsum := 0.0

	for i, p := range prob {
		accsum += p
		if accsum >= rand {
			output = i
			break
		}
	}

	for n := range circ.State {
		for i, q := range qbits {
			if (output>>i)%2 != (n>>q)%2 {
				circ.State[n] = 0
				continue
			}
		}
	}

	circ.State = circ.State.Normalize()

	return output
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
