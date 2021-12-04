// Package qsim provides functions for a quantum circuit.
package qsim

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"sort"
	"strings"
	"sync"

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
	if q.Dim() != (1 << circ.N) {
		panic("Invalid qbit length.")
	}
	circ.State = q
}

func (circ *Circuit) Apply(operator matrix.Matrix, iregs ...int) {
	if !operator.IsUnitary() {
		panic("Operator must be unitary.")
	}

	if len(operator) != 1<<len(iregs) {
		panic("Operator size does not match with input qbits.")
	}

	// Special case
	// If operator is pure => trivial parallelization is possible.
	// If operator is not pure, but still single qbit
	// => Actually, almost every gate is single qbit. This case, we can still parallelize somehow.
	if operator.IsPureGate() {
		circ.applyPure(operator, iregs...)
		return
	} else {
		if len(iregs) == 1 && circ.N > 1 {
			circ.applySingle(operator, iregs[0])
			return
		}
	}

	// Generic Fallback.

	// Tensor Product takes too long, we need another method.
	// Idea: decompose state vector to pure states?
	res := vector.Zeros(circ.State.Dim())

	for basis, amp := range circ.State {
		if amp == 0 {
			continue
		}
		// amp * |basis>
		// First, extract input qbits from basis
		// For example, if basis = |0101> and amp = 0, 2 => ibasis = |11>
		ibasis := 0
		for idx, val := range iregs {
			// Extract val-th bit from basis, plug it in to idx-th bit of ibasis.
			ibasis += ((basis >> val) % 2) << idx
		}
		// Generate new qbit from x, apply operator to it.
		newibasis_q := qbit.NewFromCbit(ibasis, len(iregs)).Apply(operator)

		for newibasis, newamp := range newibasis_q {
			// U*|ibasis> = sum newamp * |newibasis>
			if newamp == 0 {
				continue
			}
			// Make newbasis by merging newibasis to basis.
			// Extract idx-th bit from newibasis, plug it in to val-th bit of basis.
			newbasis := basis
			for idx, val := range iregs {
				bit := (newibasis >> idx) % 2
				newbasis = (newbasis | (1 << val)) - ((bit ^ 1) << val)
			}
			res[newbasis] += amp * newamp
		}
	}

	circ.State = res
}

func (circ *Circuit) applySingle(operator matrix.Matrix, ireg int) {
	// checks for operator is already done in Apply()
	// Note that operator is assumed to be non-pure.

	wg := &sync.WaitGroup{}

	// We can still parallelize 2^(n-1) loops.
	// How? Suppose basis = |0101> and ireg = 2.
	// Then, U|0101> = a|0001> + b|0101>
	// So, we can "group" |0101> and |0001>, and parallelize for 0, 1, 3th qbit.

	chunksize := (1 << ((circ.N - 1) / 2))
	res := vector.Zeros(len(circ.State))

	for i := 0; i < (1 << (circ.N - 1)); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			bases := [2]int{0, 0}
			for n := start; n < start+chunksize; n++ {
				bases[0] = ((n >> ireg) << (ireg + 1)) + (n % (1 << ireg))
				bases[1] = bases[0] | (1 << ireg)

				for ibasis, basis := range bases {
					amp := circ.State[basis]
					if amp == 0 {
						continue
					}

					// Now, same as Apply().
					newbasis_q := qbit.NewFromCbit(ibasis, 1).Apply(operator)

					// Another interesting fact:
					// newbasis_q[0] -> bases[0], newbasis_q[1] -> bases[1] by definition!
					res[bases[0]] += amp * newbasis_q[0]
					res[bases[1]] += amp * newbasis_q[1]
				}
			}
		}(i)
	}

	wg.Wait()

	circ.State = res
}

func (circ *Circuit) applyPure(operator matrix.Matrix, iregs ...int) {
	// Pure operators are trivially parallelizable.

	wg := &sync.WaitGroup{}
	chunksize := (1 << (circ.N / 2))
	res := vector.Zeros(len(circ.State))

	for i := 0; i < circ.State.Dim(); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for basis := start; basis < start+chunksize; basis++ {
				amp := circ.State[basis]

				if amp == 0 {
					continue
				}

				ibasis := 0
				for idx, val := range iregs {
					ibasis += ((basis >> val) % 2) << idx
				}

				newibasis_q := qbit.NewFromCbit(ibasis, len(iregs)).Apply(operator)

				for newibasis, newamp := range newibasis_q {
					if newamp == 0 {
						continue
					}

					newbasis := basis
					for idx, val := range iregs {
						bit := (newibasis >> idx) % 2
						newbasis = (newbasis | (1 << val)) - ((bit ^ 1) << val)
					}
					res[newbasis] = amp * newamp

					// newibasis_q is guaranteed to have only one value, so break.
					break
				}
			}
		}(i)
	}

	wg.Wait()

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

func (circ *Circuit) Control(operator matrix.Matrix, cs []int, xs []int) {
	if !operator.IsUnitary() {
		panic("Operator must be unitary.")
	}

	if len(operator) != 1<<len(xs) {
		panic("Operator size does not match with input qbits.")
	}

	if operator.IsPureGate() {
		circ.controlPure(operator, cs, xs)
		return
	} else {
		if len(xs) == 1 && circ.N > 1 {
			circ.controlSingle(operator, cs, xs[0])
			return
		}
	}

	// Generic Fallback.
	res := vector.Zeros(circ.State.Dim())

	for basis, amp := range circ.State {
		if amp == 0 {
			continue
		}

		ctrl := 0
		for _, v := range cs {
			ctrl ^= (basis >> v) % 2
		}

		if ctrl == 0 {
			res[basis] = amp
			continue
		}

		ibasis := 0
		for idx, val := range xs {
			ibasis += ((basis >> val) % 2) << idx
		}
		newibasis_q := qbit.NewFromCbit(ibasis, len(xs)).Apply(operator)

		for newibasis, newamp := range newibasis_q {
			if newamp == 0 {
				continue
			}
			newbasis := basis
			for idx, val := range xs {
				bit := (newibasis >> idx) % 2
				newbasis = (newbasis | (1 << val)) - ((bit ^ 1) << val)
			}
			res[newbasis] += amp * newamp
		}
	}

	circ.State = res
}

func (circ *Circuit) controlSingle(operator matrix.Matrix, cs []int, x int) {
	// checks for operator is already done in Apply()
	// Note that operator is assumed to be non-pure.

	wg := &sync.WaitGroup{}

	// We can still parallelize 2^(n-1) loops.
	// How? Suppose basis = |0101> and ireg = 2.
	// Then, U|0101> = a|0001> + b|0101>
	// So, we can "group" |0101> and |0001>, and parallelize for 0, 1, 3th qbit.

	chunksize := (1 << ((circ.N - 1) / 2))
	res := vector.Zeros(len(circ.State))

	for i := 0; i < (1 << (circ.N - 1)); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			bases := [2]int{0, 0}
			for n := start; n < start+chunksize; n++ {
				bases[0] = ((n >> x) << (x + 1)) + (n % (1 << x))
				bases[1] = bases[0] | (1 << x)

				ctrl := 0
				for _, v := range cs {
					ctrl ^= (bases[0] >> v) % 2
				}

				if ctrl == 0 {
					res[bases[0]] = circ.State[bases[0]]
					res[bases[1]] = circ.State[bases[1]]
					continue
				}

				for ibasis, basis := range bases {
					amp := circ.State[basis]
					if amp == 0 {
						continue
					}

					// Now, same as Apply().
					newbasis_q := qbit.NewFromCbit(ibasis, 1).Apply(operator)

					// Another interesting fact:
					// newbasis_q[0] -> bases[0], newbasis_q[1] -> bases[1] by definition!
					res[bases[0]] += amp * newbasis_q[0]
					res[bases[1]] += amp * newbasis_q[1]
				}
			}
		}(i)
	}

	wg.Wait()

	circ.State = res
}

func (circ *Circuit) controlPure(operator matrix.Matrix, cs []int, xs []int) {
	// Pure operators are trivially parallelizable.

	wg := &sync.WaitGroup{}
	chunksize := (1 << (circ.N / 2))
	res := vector.Zeros(len(circ.State))

	for i := 0; i < circ.State.Dim(); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for basis := start; basis < start+chunksize; basis++ {
				amp := circ.State[basis]

				if amp == 0 {
					continue
				}

				ctrl := 0
				for _, v := range cs {
					ctrl ^= (basis >> v) % 2
				}

				if ctrl == 0 {
					res[basis] = amp
					continue
				}

				ibasis := 0
				for idx, val := range xs {
					ibasis += ((basis >> val) % 2) << idx
				}

				newibasis_q := qbit.NewFromCbit(ibasis, len(xs)).Apply(operator)

				for newibasis, newamp := range newibasis_q {
					if newamp == 0 {
						continue
					}

					newbasis := basis
					for idx, val := range xs {
						bit := (newibasis >> idx) % 2
						newbasis = (newbasis | (1 << val)) - ((bit ^ 1) << val)
					}
					res[newbasis] = amp * newamp

					break
				}
			}
		}(i)
	}

	wg.Wait()

	circ.State = res
}

func (circ *Circuit) controlSingleSingle(operator matrix.Matrix, c int, x int) {
	circ.Control(operator, []int{c}, []int{x})
}

func (circ *Circuit) CX(c int, x int) {
	circ.controlSingleSingle(gate.X(), c, x)
}

func (circ *Circuit) Swap(x int, y int) {
	res := vector.Zeros(circ.State.Dim())

	chunksize := 1 << (circ.N / 2)
	wg := &sync.WaitGroup{}

	for i := 0; i < (1 << circ.N); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for n := start; n < start+chunksize; n++ {
				a := circ.State[n]
				if a == 0 {
					continue
				}

				bx := (n >> x) % 2
				by := (n >> y) % 2

				nn := n

				nn = (nn | (1 << x)) - ((by ^ 1) << x)
				nn = (nn | (1 << y)) - ((bx ^ 1) << y)

				res[nn] = a
			}
		}(i)
	}

	wg.Wait()

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
			circ.controlSingleSingle(gate.P(math.Pi/float64(numbers.Pow(2, i-j))), j, i)
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
			circ.controlSingleSingle(gate.P(-math.Pi/float64(numbers.Pow(2, i-j))), j, i)
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
