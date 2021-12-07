// Package qsim provides functions for a quantum circuit.
package qsim

import (
	"fmt"
	"math"
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

const PARALLEL_THRESHOLD int = 10

type Circuit struct {
	// Why save gates when we don't need diagram? :)
	N     int
	State vector.Vector
	temp  vector.Vector
}

func (circ *Circuit) cleartemp() {
	for i := range circ.temp {
		circ.temp[i] = 0
	}
}

// Generates new circuit with n qbits. Initializes with |0...0>.
func NewCircuit(n int) Circuit {
	if n < 1 {
		panic("Invalid qbit length.")
	}

	if n > 63 {
		panic("This simulator currently supports up to 63 qbits.")
	}

	q := qbit.Zeros(n)
	temp := vector.Zeros(1 << n)

	return Circuit{N: n, State: q, temp: temp}
}

// Sets state to q.
func (circ *Circuit) InitQbit(q vector.Vector) {
	if q.Dim() != (1 << circ.N) {
		panic("Invalid qbit length.")
	}
	circ.State = q
}

// Sets state to |n>.
func (circ *Circuit) InitCbit(n int) {
	if numbers.BitLength(n) > circ.N {
		panic("Invalid cbit length.")
	}

	for i := range circ.State {
		circ.State[i] = 0
	}
	circ.State[n] = 1
}

// Applies the oracle f to circuit. Maps |x>_{iregs}|y>_{oregs} -> |x>_{iregs}|y^f(x)>_{oregs}.
// NOTE: This function DOES NOT check if oracle is unitary. Use at your own risk.
func (circ *Circuit) ApplyOracle(oracle func(int) int, iregs []int, oregs []int) {
	if len(iregs) == 0 || len(oregs) == 0 {
		panic("Invalid input/output registers.")
	}

	if numbers.Min(iregs...) < 0 || numbers.Max(iregs...) >= circ.N {
		panic("Register index out of range.")
	}

	if circ.N <= PARALLEL_THRESHOLD {
		circ.applyOracleFallback(oracle, iregs, oregs)
		return
	}

	wg := &sync.WaitGroup{}
	chunksize := 1 << (circ.N / 2)
	circ.cleartemp()

	for i := 0; i < len(circ.State); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for basis := start; basis < start+chunksize; basis++ {
				amp := circ.State[basis]
				if amp == 0 {
					continue
				}

				input := 0
				for idx, val := range iregs {
					input += ((basis >> val) & 1) << idx
				}

				output := oracle(input)

				newbasis := basis
				for idx, val := range oregs {
					bit := (output >> idx) & 1
					newbasis ^= bit << val
				}

				circ.temp[newbasis] = amp
			}
		}(i)
	}

	wg.Wait()

	copy(circ.State, circ.temp)
}

// Fallback of ApplyOracle(). No optimization.
func (circ *Circuit) applyOracleFallback(oracle func(int) int, iregs []int, oregs []int) {
	circ.cleartemp()

	for basis, amp := range circ.State {
		if amp == 0 {
			continue
		}

		input := 0
		for idx, val := range iregs {
			input += ((basis >> val) & 1) << idx
		}

		output := oracle(input)

		newbasis := basis
		for idx, val := range oregs {
			bit := (output >> idx) & 1
			newbasis ^= bit << val
		}

		circ.temp[newbasis] = amp
	}

	copy(circ.State, circ.temp)
}

// Applies the operator to iregs.
func (circ *Circuit) Apply(operator matrix.Matrix, iregs ...int) {
	if !operator.IsUnitary() {
		panic("Operator must be unitary.")
	}

	if len(operator) != 1<<len(iregs) {
		panic("Operator size does not match with input qbits.")
	}

	if numbers.Min(iregs...) < 0 || numbers.Max(iregs...) >= circ.N {
		panic("Register index out of range.")
	}

	if circ.N > PARALLEL_THRESHOLD {
		if len(iregs) == 1 {
			circ.applyOneQbit(operator, iregs[0])
			return
		}

		if len(iregs) == 2 {
			circ.applyTwoQbit(operator, iregs[0], iregs[1])
			return
		}
	} else {
		if len(iregs) == 1 {
			circ.applyOneQbitFallback(operator, iregs[0])
			return
		}

		if len(iregs) == 2 {
			circ.applyTwoQbitFallback(operator, iregs[0], iregs[1])
			return
		}
	}

	// Generic Fallback.
	circ.applyFallback(operator, iregs...)
}

// Special case of Apply(), when there is only one input registers.
// Implemented as in-place swapping.
func (circ *Circuit) applyOneQbit(operator matrix.Matrix, ireg int) {
	wg := &sync.WaitGroup{}
	chunksize := 1 << ((circ.N - 1) / 2)

	for i := 0; i < (len(circ.State) >> 1); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for n := start; n < start+chunksize; n++ {
				n0 := ((n >> ireg) << (ireg + 1)) + (n & ((1 << ireg) - 1))
				n1 := n0 | (1 << ireg)

				amp0 := circ.State[n0]
				amp1 := circ.State[n1]

				circ.State[n0] = amp0*operator[0][0] + amp1*operator[0][1]
				circ.State[n1] = amp0*operator[1][0] + amp1*operator[1][1]
			}
		}(i)
	}

	wg.Wait()
}

// applyOneQbit() with no parallelization.
func (circ *Circuit) applyOneQbitFallback(operator matrix.Matrix, ireg int) {
	for n := 0; n < (len(circ.State) >> 1); n++ {
		n0 := ((n >> ireg) << (ireg + 1)) + (n & ((1 << ireg) - 1))
		n1 := n0 | (1 << ireg)

		amp0 := circ.State[n0]
		amp1 := circ.State[n1]

		circ.State[n0] = amp0*operator[0][0] + amp1*operator[0][1]
		circ.State[n1] = amp0*operator[1][0] + amp1*operator[1][1]
	}
}

// Special case of Apply(), when there are two input registers.
// Implemented as in-place swapping.
func (circ *Circuit) applyTwoQbit(operator matrix.Matrix, ireg0, ireg1 int) {
	if ireg0 == ireg1 {
		panic("Same input registers.")
	}

	if circ.N == 2 {
		// No need for further optimization, it is fast enough. (I think.)
		circ.State = circ.State.Apply(operator)
		return
	}

	wg := &sync.WaitGroup{}
	chunksize := 1 << ((circ.N - 2) / 2)

	for i := 0; i < (len(circ.State) >> 2); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for n := start; n < start+chunksize; n++ {
				n0 := ((n >> ireg0) << (ireg0 + 1)) + (n & ((1 << ireg0) - 1))
				n1 := n0 | (1 << ireg0)

				n00 := ((n0 >> ireg1) << (ireg1 + 1)) + (n0 & ((1 << ireg1) - 1))
				n01 := n00 | (1 << ireg1)
				n10 := ((n1 >> ireg1) << (ireg1 + 1)) + (n1 & ((1 << ireg1) - 1))
				n11 := n10 | (1 << ireg1)

				amp00 := circ.State[n00]
				amp01 := circ.State[n01]
				amp10 := circ.State[n10]
				amp11 := circ.State[n11]

				circ.State[n00] = amp00*operator[0][0] + amp01*operator[0][1] + amp10*operator[0][2] + amp11*operator[0][1]
				circ.State[n01] = amp00*operator[1][0] + amp01*operator[1][1] + amp10*operator[1][2] + amp11*operator[1][1]
				circ.State[n10] = amp00*operator[2][0] + amp01*operator[2][1] + amp10*operator[2][2] + amp11*operator[2][1]
				circ.State[n11] = amp00*operator[3][0] + amp01*operator[3][1] + amp10*operator[3][2] + amp11*operator[3][1]
			}
		}(i)
	}

	wg.Wait()
}

// applyTwoQbit() with no parallelization.
func (circ *Circuit) applyTwoQbitFallback(operator matrix.Matrix, ireg0, ireg1 int) {
	for n := 0; n < (len(circ.State) >> 2); n++ {
		n0 := ((n >> ireg0) << (ireg0 + 1)) + (n & ((1 << ireg0) - 1))
		n1 := n0 | (1 << ireg0)

		n00 := ((n0 >> ireg1) << (ireg1 + 1)) + (n0 & ((1 << ireg1) - 1))
		n01 := n00 | (1 << ireg1)
		n10 := ((n1 >> ireg1) << (ireg1 + 1)) + (n1 & ((1 << ireg1) - 1))
		n11 := n10 | (1 << ireg1)

		amp00 := circ.State[n00]
		amp01 := circ.State[n01]
		amp10 := circ.State[n10]
		amp11 := circ.State[n11]

		circ.State[n00] = amp00*operator[0][0] + amp01*operator[0][1] + amp10*operator[0][2] + amp11*operator[0][1]
		circ.State[n01] = amp00*operator[1][0] + amp01*operator[1][1] + amp10*operator[1][2] + amp11*operator[1][1]
		circ.State[n10] = amp00*operator[2][0] + amp01*operator[2][1] + amp10*operator[2][2] + amp11*operator[2][1]
		circ.State[n11] = amp00*operator[3][0] + amp01*operator[3][1] + amp10*operator[3][2] + amp11*operator[3][1]
	}
}

// General Fallback for Apply(). No optimization.
func (circ *Circuit) applyFallback(operator matrix.Matrix, iregs ...int) {
	circ.cleartemp()

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
			ibasis += ((basis >> val) & 1) << idx
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
				bit := (newibasis >> idx) & 1
				newbasis = (newbasis | (1 << val)) - ((bit ^ 1) << val)
			}
			circ.temp[newbasis] += amp * newamp
		}
	}

	copy(circ.State, circ.temp)
}

// Applies the I gate to ireg.
func (circ *Circuit) I(ireg int) {
	// Just Do Nothing. lol.
}

// Applies the X gate to ireg.
func (circ *Circuit) X(ireg int) {
	circ.Apply(gate.X(), ireg)
}

// Applies the Y gate to ireg.
func (circ *Circuit) Y(ireg int) {
	circ.Apply(gate.Y(), ireg)
}

// Applies the Z gate to ireg.
func (circ *Circuit) Z(ireg int) {
	circ.Apply(gate.Z(), ireg)
}

// Applies the H gate to ireg.
func (circ *Circuit) H(ireg int) {
	circ.Apply(gate.H(), ireg)
}

// Applies the P gate to ireg.
func (circ *Circuit) P(phi float64, ireg int) {
	circ.Apply(gate.P(phi), ireg)
}

// Applies the S gate to ireg.
func (circ *Circuit) S(ireg int) {
	circ.Apply(gate.S(), ireg)
}

// Applies the T gate to ireg.
func (circ *Circuit) T(ireg int) {
	circ.Apply(gate.T(), ireg)
}

// Used for calculating control bits in control-functions.
func checkControlBit(n int, cs []int) bool {
	res := 0

	for _, idx := range cs {
		res ^= (n >> idx) & 1
	}

	return res == 1
}

// Applies the control-version of operator to circuit. cs is the control qbits, xs is the input qbits.
func (circ *Circuit) Control(operator matrix.Matrix, cs []int, xs []int) {
	if !operator.IsUnitary() {
		panic("Operator must be unitary.")
	}

	if len(operator) != 1<<len(xs) {
		panic("Operator size does not match with input qbits.")
	}

	if numbers.Min(cs...) < 0 || numbers.Max(cs...) >= circ.N || numbers.Min(xs...) < 0 || numbers.Max(xs...) >= circ.N {
		panic("Register index out of range.")
	}

	if len(cs)+len(xs) > circ.N {
		panic("Too many registers.")
	}

	if circ.N > PARALLEL_THRESHOLD {
		if len(xs) == 1 {
			circ.controlOneQbit(operator, cs, xs[0])
			return
		}

		if len(xs) == 2 {
			circ.controlTwoQubit(operator, cs, xs[0], xs[1])
			return
		}
	} else {
		if len(xs) == 1 {
			circ.controlOneQbitFallback(operator, cs, xs[0])
			return
		}

		if len(xs) == 2 {
			circ.controlTwoQbitFallback(operator, cs, xs[0], xs[1])
			return
		}
	}

	// Generic Fallback.
	circ.controlFallback(operator, cs, xs)
}

// Special case of Apply(), when there is only one input registers.
// Implemented as in-place swapping.
func (circ *Circuit) controlOneQbit(operator matrix.Matrix, cs []int, ireg int) {
	cs_shifted := make([]int, len(cs))

	for i := range cs_shifted {
		cs_shifted[i] = cs[i]
		if cs[i] > ireg {
			cs_shifted[i]--
		}
	}

	wg := &sync.WaitGroup{}
	chunksize := 1 << ((circ.N - 1) / 2)

	for i := 0; i < (len(circ.State) >> 1); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for n := start; n < start+chunksize; n++ {
				if !checkControlBit(n, cs_shifted) {
					continue
				}

				n0 := ((n >> ireg) << (ireg + 1)) + (n & ((1 << ireg) - 1))
				n1 := n0 | (1 << ireg)

				amp0 := circ.State[n0]
				amp1 := circ.State[n1]

				circ.State[n0] = amp0*operator[0][0] + amp1*operator[0][1]
				circ.State[n1] = amp0*operator[1][0] + amp1*operator[1][1]
			}
		}(i)
	}

	wg.Wait()
}

// Fallback of controlOneQbit(), with no parallelization.
func (circ *Circuit) controlOneQbitFallback(operator matrix.Matrix, cs []int, ireg int) {
	cs_shifted := make([]int, len(cs))

	for i := range cs_shifted {
		cs_shifted[i] = cs[i]
		if cs[i] > ireg {
			cs_shifted[i]--
		}
	}

	for n := 0; n < (len(circ.State) >> 1); n++ {
		if !checkControlBit(n, cs_shifted) {
			continue
		}

		n0 := ((n >> ireg) << (ireg + 1)) + (n & ((1 << ireg) - 1))
		n1 := n0 | (1 << ireg)

		amp0 := circ.State[n0]
		amp1 := circ.State[n1]

		circ.State[n0] = amp0*operator[0][0] + amp1*operator[0][1]
		circ.State[n1] = amp0*operator[1][0] + amp1*operator[1][1]
	}

}

// Special case of Control(), when there are two input registers.
// Implemented as in-place swapping.
func (circ *Circuit) controlTwoQubit(operator matrix.Matrix, cs []int, ireg0, ireg1 int) {
	if ireg0 == ireg1 {
		panic("Same input registers.")
	}

	cs_shifted := make([]int, len(cs))

	for i := range cs_shifted {
		cs_shifted[i] = cs[i]
		if cs[i] > ireg0 {
			cs_shifted[i]--
		}
		if cs[i] > ireg1 {
			cs_shifted[i]--
		}
	}

	wg := &sync.WaitGroup{}
	chunksize := 1 << ((circ.N - 2) / 2)

	for i := 0; i < (len(circ.State) >> 2); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for n := start; n < start+chunksize; n++ {
				if !checkControlBit(n, cs_shifted) {
					continue
				}

				n0 := ((n >> ireg0) << (ireg0 + 1)) + (n & ((1 << ireg0) - 1))
				n1 := n0 | (1 << ireg0)

				n00 := ((n0 >> ireg1) << (ireg1 + 1)) + (n0 & ((1 << ireg1) - 1))
				n01 := n00 | (1 << ireg1)
				n10 := ((n1 >> ireg1) << (ireg1 + 1)) + (n1 & ((1 << ireg1) - 1))
				n11 := n10 | (1 << ireg1)

				amp00 := circ.State[n00]
				amp01 := circ.State[n01]
				amp10 := circ.State[n10]
				amp11 := circ.State[n11]

				circ.State[n00] = amp00*operator[0][0] + amp01*operator[0][1] + amp10*operator[0][2] + amp11*operator[0][1]
				circ.State[n01] = amp00*operator[1][0] + amp01*operator[1][1] + amp10*operator[1][2] + amp11*operator[1][1]
				circ.State[n10] = amp00*operator[2][0] + amp01*operator[2][1] + amp10*operator[2][2] + amp11*operator[2][1]
				circ.State[n11] = amp00*operator[3][0] + amp01*operator[3][1] + amp10*operator[3][2] + amp11*operator[3][1]
			}
		}(i)
	}

	wg.Wait()
}

// Fallback of controlTwoQbit(), with no parallelization.
func (circ *Circuit) controlTwoQbitFallback(operator matrix.Matrix, cs []int, ireg0, ireg1 int) {
	if ireg0 == ireg1 {
		panic("Same input registers.")
	}

	cs_shifted := make([]int, len(cs))

	for i := range cs_shifted {
		cs_shifted[i] = cs[i]
		if cs[i] > ireg0 {
			cs_shifted[i]--
		}
		if cs[i] > ireg1 {
			cs_shifted[i]--
		}
	}

	for n := 0; n < (len(circ.State) >> 2); n++ {
		if !checkControlBit(n, cs_shifted) {
			continue
		}

		n0 := ((n >> ireg0) << (ireg0 + 1)) + (n & ((1 << ireg0) - 1))
		n1 := n0 | (1 << ireg0)

		n00 := ((n0 >> ireg1) << (ireg1 + 1)) + (n0 & ((1 << ireg1) - 1))
		n01 := n00 | (1 << ireg1)
		n10 := ((n1 >> ireg1) << (ireg1 + 1)) + (n1 & ((1 << ireg1) - 1))
		n11 := n10 | (1 << ireg1)

		amp00 := circ.State[n00]
		amp01 := circ.State[n01]
		amp10 := circ.State[n10]
		amp11 := circ.State[n11]

		circ.State[n00] = amp00*operator[0][0] + amp01*operator[0][1] + amp10*operator[0][2] + amp11*operator[0][1]
		circ.State[n01] = amp00*operator[1][0] + amp01*operator[1][1] + amp10*operator[1][2] + amp11*operator[1][1]
		circ.State[n10] = amp00*operator[2][0] + amp01*operator[2][1] + amp10*operator[2][2] + amp11*operator[2][1]
		circ.State[n11] = amp00*operator[3][0] + amp01*operator[3][1] + amp10*operator[3][2] + amp11*operator[3][1]
	}
}

// General Fallback for Control(). Similar to applyFallback(), No optimization.
func (circ *Circuit) controlFallback(operator matrix.Matrix, cs []int, xs []int) {
	circ.cleartemp()

	for basis, amp := range circ.State {
		if amp == 0 {
			continue
		}

		if !checkControlBit(basis, cs) {
			continue
		}

		ibasis := 0
		for idx, val := range xs {
			ibasis += ((basis >> val) & 1) << idx
		}
		newibasis_q := qbit.NewFromCbit(ibasis, len(xs)).Apply(operator)

		for newibasis, newamp := range newibasis_q {
			if newamp == 0 {
				continue
			}
			newbasis := basis
			for idx, val := range xs {
				bit := (newibasis >> idx) & 1
				newbasis = (newbasis | (1 << val)) - ((bit ^ 1) << val)
			}
			circ.temp[newbasis] += amp * newamp
		}
	}

	copy(circ.State, circ.temp)
}

// Alias for Control() when control qbit and input qbit are all single.
func (circ *Circuit) controlSingleSingle(operator matrix.Matrix, c int, x int) {
	circ.Control(operator, []int{c}, []int{x})
}

// Applies Control-X gate to circuit. c is control qbit, and x is input qbit.
func (circ *Circuit) CX(c int, x int) {
	circ.controlSingleSingle(gate.X(), c, x)
}

// Applies Tofolli gate(CCX gate) to circuit. c1, c2 are control qbits, and x is input qbit.
func (circ *Circuit) CCX(c1, c2, x int) {
	circ.Control(gate.X(), []int{c1, c2}, []int{x})
}

// Swaps two qbits.
func (circ *Circuit) Swap(x int, y int) {
	if x < 0 || x >= circ.N || y < 0 || y >= circ.N {
		panic("Register index out of range.")
	}

	if x == y {
		panic("Swapping same registers.")
	}

	if circ.N == 2 {
		circ.State[0b01], circ.State[0b10] = circ.State[0b10], circ.State[0b01]
		return
	}

	if circ.N <= PARALLEL_THRESHOLD {
		circ.swapFallback(x, y)
		return
	}

	chunksize := 1 << ((circ.N - 2) / 2)
	wg := &sync.WaitGroup{}

	for i := 0; i < (len(circ.State) >> 2); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for n := start; n < start+chunksize; n++ {
				n0 := ((n >> x) << (x + 1)) + (n & ((1 << x) - 1))
				n1 := n0 | (1 << x)

				n01 := ((n0 >> y) << (y + 1)) + (n0 & ((1 << y) - 1)) + (1 << y)
				n10 := ((n1 >> y) << (y + 1)) + (n1 & ((1 << y) - 1))

				circ.State[n01], circ.State[n10] = circ.State[n10], circ.State[n01]
			}
		}(i)
	}

	wg.Wait()
}

// SWAP without goroutine.
func (circ *Circuit) swapFallback(x, y int) {
	for n := 0; n < (len(circ.State) >> 2); n++ {
		n0 := ((n >> x) << (x + 1)) + (n & ((1 << x) - 1))
		n1 := n0 | (1 << x)

		n01 := ((n0 >> y) << (y + 1)) + (n0 & ((1 << y) - 1)) + (1 << y)
		n10 := ((n1 >> y) << (y + 1)) + (n1 & ((1 << y) - 1))

		circ.State[n01], circ.State[n10] = circ.State[n10], circ.State[n01]
	}
}

// Applies QFT to [start, end).
func (circ *Circuit) QFT(start, end int) {
	if start < 0 || end > circ.N {
		panic("Index out of range.")
	}

	if start >= end {
		panic("Invalid start / end parameters.")
	}

	phis := make([]float64, end-start)

	for i := range phis {
		phis[i] = math.Pi / math.Pow(2.0, float64(i))
	}

	for i := end - 1; i >= start; i-- {
		circ.H(i)
		for j := start; j < i; j++ {
			circ.controlSingleSingle(gate.P(phis[i-j]), j, i)
		}
	}

	for i, j := start, end-1; i < j; i, j = i+1, j-1 {
		circ.Swap(i, j)
	}

}

// Applies Inverse QFT to [start, end).
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

	phis := make([]float64, end-start)

	for i := range phis {
		phis[i] = -math.Pi / math.Pow(2.0, float64(i))
	}

	for i := start; i < end; i++ {
		for j := start; j < i; j++ {
			circ.controlSingleSingle(gate.P(phis[i-j]), j, i)
		}
		circ.H(i)
	}
}

// Measure qbits.
func (circ *Circuit) Measure(qbits ...int) int {
	sort.Ints(qbits)

	if qbits[0] < 0 || qbits[len(qbits)-1] > circ.N-1 {
		panic("Register index out of range.")
	}

	probs := make([]float64, 1<<len(qbits))

	for n, amp := range circ.State {
		if amp == 0 {
			continue
		}
		o := 0
		for i, q := range qbits {
			o += ((n >> q) & 1) << i
		}
		probs[o] += real(amp)*real(amp) + imag(amp)*imag(amp)
	}

	// Wait, Golang does not have weighted sampling? WTF.
	rand := rand.Float64()

	output := 0
	accsum := 0.0

	for i, p := range probs {
		accsum += p
		if accsum >= rand {
			output = i
			break
		}
	}

	s := complex(math.Sqrt(probs[output]), 0)

	for n, amp := range circ.State {
		if amp == 0 {
			continue
		}

		for i, q := range qbits {
			if (n>>q)&1 != (output>>i)&1 {
				circ.State[n] = 0
				goto endloop
			}
		}
		circ.State[n] /= s
	endloop:
	}

	return output
}

// Prints current state to string.
func (circ Circuit) StateToString() string {
	q := circ.State

	qs := make([]string, 0)

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

		qs = append(qs, res+fmt.Sprintf("|%0*b>", circ.N, i))
	}

	return strings.Join(qs, " + ")
}
