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

var ZEROVEC vector.Vector

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

				circ.temp[newbasis] += amp
			}
		}(i)
	}

	wg.Wait()

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

	// Special case
	// If operator is generalized permutation matrix => trivial parallelization is possible.
	// If operator is not generalized permutation matrix, but still single qbit
	// => Actually, almost every gate is single qbit. This case, we can still parallelize somehow.
	if operator.IsGenPermutMatrix() {
		circ.applyGenPermut(operator, iregs...)
		return
	}

	if len(iregs) == 1 {
		circ.applySingle(operator, iregs[0])
		return
	}

	// Generic Fallback.

	// Tensor Product takes too long, we need another method.
	// Idea: decompose state vector to basis states?
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

// Special case of apply, when operator is a single-qbit gate.
// This case, we use parallelization for 2^(n-1) loops.
func (circ *Circuit) applySingle(operator matrix.Matrix, ireg int) {
	// checks for operator is already done in Apply()
	// Note that operator is assumed to be non-generalized permutation matrix.

	wg := &sync.WaitGroup{}

	// We can still parallelize 2^(n-1) loops.
	// How? Suppose basis = |0101> and ireg = 2.
	// Then, U|0101> = a|0001> + b|0101>
	// So, we can "group" |0101> and |0001>, and parallelize for 0, 1, 3th qbit.

	chunksize := (1 << ((circ.N - 1) / 2))
	circ.cleartemp()

	memo := [2][2]complex128{{operator[0][0], operator[1][0]}, {operator[0][1], operator[1][1]}}

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

					circ.temp[bases[0]] += amp * memo[0][ibasis]
					circ.temp[bases[1]] += amp * memo[1][ibasis]
				}
			}
		}(i)
	}

	wg.Wait()

	copy(circ.State, circ.temp)
}

// Special case of apply, when operator is generalized permutation matrix.
// This case, trivial parallelization is possible.
func (circ *Circuit) applyGenPermut(operator matrix.Matrix, iregs ...int) {
	// Generalized permutation matrix operators are trivially parallelizable.

	// Precomputate maps.
	memo_basis := make([]int, 1<<len(iregs))
	memo_amp := make([]complex128, 1<<len(iregs))

	wg_memo := &sync.WaitGroup{}

	for i := 0; i < len(operator); i++ {
		wg_memo.Add(1)
		go func(col int) {
			defer wg_memo.Done()

			for j := 0; j < len(operator); j++ {
				if operator[col][j] != 0 {
					memo_basis[col] = j
					memo_amp[col] = operator[j][col]
				}
			}
		}(i)
	}

	wg_memo.Wait()

	wg := &sync.WaitGroup{}
	chunksize := 1 << (circ.N / 2)
	circ.cleartemp()

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
					ibasis += ((basis >> val) & 1) << idx
				}

				newibasis := memo_basis[ibasis]
				newamp := memo_amp[ibasis]

				newbasis := basis
				for idx, val := range iregs {
					bit := (newibasis >> idx) & 1
					newbasis = (newbasis | (1 << val)) - ((bit ^ 1) << val)
				}
				circ.temp[newbasis] = amp * newamp
			}
		}(i)
	}

	wg.Wait()

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

	if operator.IsGenPermutMatrix() {
		circ.controlGenPermut(operator, cs, xs)
		return
	}

	if len(xs) == 1 {
		circ.controlSingle(operator, cs, xs[0])
		return
	}

	// Generic Fallback.
	circ.cleartemp()

	for basis, amp := range circ.State {
		if amp == 0 {
			continue
		}

		ctrl := 0
		for _, v := range cs {
			ctrl ^= (basis >> v) & 1
		}

		if ctrl == 0 {
			circ.temp[basis] = amp
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

// Special case of control-operator, when the operator is single qbit.
// Similar to applySingle.
func (circ *Circuit) controlSingle(operator matrix.Matrix, cs []int, x int) {
	wg := &sync.WaitGroup{}

	chunksize := (1 << ((circ.N - 1) / 2))
	circ.cleartemp()

	memo := [2][2]complex128{{operator[0][0], operator[1][0]}, {operator[0][1], operator[1][1]}}

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
					ctrl ^= (bases[0] >> v) & 1
				}

				if ctrl == 0 {
					circ.temp[bases[0]] = circ.State[bases[0]]
					circ.temp[bases[1]] = circ.State[bases[1]]
					continue
				}

				for ibasis, basis := range bases {
					amp := circ.State[basis]
					if amp == 0 {
						continue
					}

					circ.temp[bases[0]] += amp * memo[0][ibasis]
					circ.temp[bases[1]] += amp * memo[1][ibasis]
				}
			}
		}(i)
	}

	wg.Wait()

	copy(circ.State, circ.temp)
}

// Special case of control-operator, when the operator is generalized permutation matrix.
// Similar to applyGenPermut.
func (circ *Circuit) controlGenPermut(operator matrix.Matrix, cs []int, xs []int) {
	// Precomputate maps.
	memo_basis := make([]int, 1<<len(xs))
	memo_amp := make([]complex128, 1<<len(xs))

	wg_memo := &sync.WaitGroup{}

	for i := 0; i < len(operator); i++ {
		wg_memo.Add(1)
		go func(col int) {
			defer wg_memo.Done()

			for j := 0; j < len(operator); j++ {
				if operator[col][j] != 0 {
					memo_basis[col] = j
					memo_amp[col] = operator[j][col]
				}
			}
		}(i)
	}

	wg_memo.Wait()

	wg := &sync.WaitGroup{}
	chunksize := (1 << (circ.N / 2))
	circ.cleartemp()

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
					ctrl ^= (basis >> v) & 1
				}

				if ctrl == 0 {
					circ.temp[basis] = amp
					continue
				}

				ibasis := 0
				for idx, val := range xs {
					ibasis += ((basis >> val) & 1) << idx
				}

				newibasis := memo_basis[ibasis]
				newamp := memo_amp[ibasis]

				newbasis := basis
				for idx, val := range xs {
					bit := (newibasis >> idx) & 1
					newbasis = (newbasis | (1 << val)) - ((bit ^ 1) << val)
				}
				circ.temp[newbasis] = amp * newamp
			}
		}(i)
	}

	wg.Wait()

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
	chunksize := 1 << (circ.N / 2)
	wg := &sync.WaitGroup{}
	circ.cleartemp()

	for i := 0; i < len(circ.State); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for n := start; n < start+chunksize; n++ {
				a := circ.State[n]
				if a == 0 {
					continue
				}

				bx := (n >> x) & 1
				by := (n >> y) & 1

				nn := n

				nn = (nn | (1 << x)) - ((by ^ 1) << x)
				nn = (nn | (1 << y)) - ((bx ^ 1) << y)

				circ.temp[nn] = a
			}
		}(i)
	}

	wg.Wait()

	copy(circ.State, circ.temp)
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
	qbits = sort.IntSlice(qbits)

	if qbits[0] < 0 || qbits[len(qbits)-1] > circ.N-1 {
		panic("Register index out of range.")
	}

	prob := make([]float64, 1<<len(qbits))

	for n, a := range circ.State {
		if a == 0 {
			continue
		}
		o := 0
		for i, q := range qbits {
			o += ((n >> q) & 1) << i
		}
		prob[o] += cmplx.Abs(a) * cmplx.Abs(a)
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

	chunksize := 1 << (circ.N / 2)
	wg := &sync.WaitGroup{}

	for i := 0; i < len(circ.State); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for n := start; n < start+chunksize; n++ {
				for i, q := range qbits {
					if (output>>i)&1 != (n>>q)&1 {
						circ.State[n] = 0
						continue
					}
				}
			}
		}(i)
	}

	wg.Wait()

	// In-place normalization with goroutines.
	var sum complex128 = 0
	for _, v := range circ.State {
		sum += complex(cmplx.Abs(v)*cmplx.Abs(v), 0)
	}

	wg = &sync.WaitGroup{}

	for i := 0; i < len(circ.State); i += chunksize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			for n := start; n < start+chunksize; n++ {
				circ.State[n] /= sum
			}
		}(i)
	}

	wg.Wait()

	return output
}

// Prints current state to string.
func (circ Circuit) StateToString() string {
	q := circ.State

	qs := make([]string, 0)
	d := numbers.BitLength(q.Dim()) - 1

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
