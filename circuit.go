package qsim

import (
	"math"
	"math/rand"
	"runtime"
	"sort"
	"sync"

	"github.com/sp301415/qsim/math/numbers"
	"github.com/sp301415/qsim/math/vec"
)

// Options for a circuit.
type Options struct {
	GOROUTINE_CNT      int // Number of goroutines to execute. Defaults to GOMAXPROCS.
	PARALLEL_THRESHOLD int // Size threshold to use parallelization. Defaults to 8.
}

type Circuit struct {
	state  Qubit   // State qubit of this circuit.
	temp   Qubit   // Used for some apply functions.
	Option Options // Options for this circuit.
}

// Clears temp qubit.
func (c *Circuit) cleartemp() {
	for i := range c.temp.data {
		c.temp.data[i] = 0
	}
}

// NewCircuit initializes circuit with nbits size.
func NewCircuit(nbits int) *Circuit {
	if nbits < 0 || nbits > 20 {
		panic("Unsupported amount of qubits. Currently qsim supports up to 20 qubits.")
	}

	return &Circuit{
		state:  NewBit(0, nbits),
		temp:   NewQubit(vec.NewVec(1 << nbits)),
		Option: Options{GOROUTINE_CNT: runtime.GOMAXPROCS(0), PARALLEL_THRESHOLD: 10},
	}
}

// SetBit sets the state qubit to given number.
func (c *Circuit) SetBit(n int) {
	c.state = NewBit(n, c.Size())
}

// Size returns the qubit length of this circuit.
func (c Circuit) Size() int {
	return c.state.size
}

// State returns the copy of this circuit's state.
func (c Circuit) State() Qubit {
	return c.state.Copy()
}

// Gates.

// Applies the I gate.
func (c *Circuit) I(i int) {
	// Just Do Nothing. lol.
}

// Applies the X gate.
func (c *Circuit) X(i int) {
	c.Apply(X(), i)
}

// Applies the Y gate.
func (c *Circuit) Y(i int) {
	c.Apply(Y(), i)
}

// Applies the Z gate.
func (c *Circuit) Z(i int) {
	c.Apply(Z(), i)
}

// Applies the H gate.
func (c *Circuit) H(i int) {
	c.Apply(H(), i)
}

// Applies the P gate.
func (c *Circuit) P(phi float64, i int) {
	c.Apply(P(phi), i)
}

// Applies the S gate.
func (c *Circuit) S(i int) {
	c.Apply(S(), i)
}

// Applies the T gate.
func (c *Circuit) T(i int) {
	c.Apply(T(), i)
}

// Applies the CX gate.
func (circ *Circuit) CX(c0, i int) {
	circ.Control(X(), []int{c0}, []int{i})
}

// Applies the CCX gate.
func (circ *Circuit) CCX(c0, c1, i int) {
	circ.Control(X(), []int{c0, c1}, []int{i})
}

// Apply.

// Apply applies the given gates.
func (c *Circuit) Apply(op Gate, iregs ...int) {
	if len(iregs) == 0 {
		panic("At least one input registers required.")
	}

	if len(iregs) != op.Size() {
		panic("Operator size does not match input registers.")
	}

	if numbers.Min(iregs...) < 0 || numbers.Max(iregs...) > c.Size() {
		panic("Registers out of range.")
	}

	// Special treatment for one and two qubit gates.
	if c.Size() > c.Option.PARALLEL_THRESHOLD {
		switch len(iregs) {
		case 1:
			c.applyOneParallel(op, iregs[0])
			return
		case 2:
			c.applyTwoParallel(op, iregs[0], iregs[1])
			return
		}
	} else {
		switch len(iregs) {
		case 1:
			c.applyOne(op, iregs[0])
			return
		case 2:
			c.applyTwo(op, iregs[0], iregs[1])
			return
		}
	}

	c.applyGeneral(op, iregs...)
}

// applyOne applies one qubit gate.
func (c *Circuit) applyOne(op Gate, i int) {
	lo := 1 << i

	for n := 0; n < c.state.Dim()/2; n++ {
		// n0 = XXX0XXX, n1 = XXX1XXX
		n0 := ((n >> i) << (i + 1)) + (n % lo)
		n1 := n0 | lo

		a0 := c.state.data[n0]
		a1 := c.state.data[n1]

		c.state.data[n0] = a0*op.data[0][0] + a1*op.data[0][1]
		c.state.data[n1] = a0*op.data[1][0] + a1*op.data[1][1]
	}
}

// applyOneParallel applies one qubit gate with parallelization.
func (c *Circuit) applyOneParallel(op Gate, i int) {
	jobsize := c.state.Dim() / 2
	chunksize := jobsize / c.Option.GOROUTINE_CNT

	wg := &sync.WaitGroup{}
	wg.Add(c.Option.GOROUTINE_CNT)

	lo := 1 << i

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
		start := n * chunksize
		end := start + chunksize

		if n == c.Option.GOROUTINE_CNT-1 {
			end += jobsize % c.Option.GOROUTINE_CNT
		}

		go func(start, end int) {
			defer wg.Done()

			for n := start; n < end; n++ {
				n0 := ((n >> i) << (i + 1)) + (n % lo)
				n1 := n0 | lo

				a0 := c.state.data[n0]
				a1 := c.state.data[n1]

				c.state.data[n0] = a0*op.data[0][0] + a1*op.data[0][1]
				c.state.data[n1] = a0*op.data[1][0] + a1*op.data[1][1]
			}
		}(start, end)
	}

	wg.Wait()
}

// applyTwo applies two qubit gate.
func (c *Circuit) applyTwo(op Gate, i0, i1 int) {
	if i0 == i1 {
		panic("Cannot apply gate to same registers.")
	}

	lo0 := 1 << i0
	lo1 := 1 << i1

	for n := 0; n < c.state.Dim()/4; n++ {
		// n00 = XXX0(i1)XXX0(i0)XXX
		// n01 = XXX0(i1)XXX1(i0)XXX
		// ...

		n0 := ((n >> i0) << (i0 + 1)) + (n % lo0)

		n00 := ((n0 >> i1) << (i1 + 1)) + (n % lo1)
		n01 := n00 | lo0
		n10 := n00 | lo1
		n11 := n10 | lo0

		a00 := c.state.data[n00]
		a01 := c.state.data[n01]
		a10 := c.state.data[n10]
		a11 := c.state.data[n11]

		c.state.data[n00] = a00*op.data[0][0] + a01*op.data[0][1] + a10*op.data[0][2] + a11*op.data[0][3]
		c.state.data[n01] = a00*op.data[1][0] + a01*op.data[1][1] + a10*op.data[1][2] + a11*op.data[1][3]
		c.state.data[n10] = a00*op.data[2][0] + a01*op.data[2][1] + a10*op.data[2][2] + a11*op.data[2][3]
		c.state.data[n11] = a00*op.data[3][0] + a01*op.data[3][1] + a10*op.data[3][2] + a11*op.data[3][3]
	}
}

// applyTwoParallel applies two qubit gate with parallelizaition.
func (c *Circuit) applyTwoParallel(op Gate, i0, i1 int) {
	jobsize := c.state.Dim() / 4
	chunksize := jobsize / c.Option.GOROUTINE_CNT

	wg := &sync.WaitGroup{}
	wg.Add(c.Option.GOROUTINE_CNT)

	lo0 := 1 << i0
	lo1 := 1 << i1

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
		start := n * chunksize
		end := start + chunksize

		if n == c.Option.GOROUTINE_CNT-1 {
			end += jobsize % c.Option.GOROUTINE_CNT
		}

		go func(start, end int) {
			defer wg.Done()

			for n := start; n < end; n++ {
				n0 := ((n >> i0) << (i0 + 1)) + (n % lo0)

				n00 := ((n0 >> i1) << (i1 + 1)) + (n % lo1)
				n01 := n00 | lo0
				n10 := n00 | lo1
				n11 := n10 | lo0

				a00 := c.state.data[n00]
				a01 := c.state.data[n01]
				a10 := c.state.data[n10]
				a11 := c.state.data[n11]

				c.state.data[n00] = a00*op.data[0][0] + a01*op.data[0][1] + a10*op.data[0][2] + a11*op.data[0][3]
				c.state.data[n01] = a00*op.data[1][0] + a01*op.data[1][1] + a10*op.data[1][2] + a11*op.data[1][3]
				c.state.data[n10] = a00*op.data[2][0] + a01*op.data[2][1] + a10*op.data[2][2] + a11*op.data[2][3]
				c.state.data[n11] = a00*op.data[3][0] + a01*op.data[3][1] + a10*op.data[3][2] + a11*op.data[3][3]
			}
		}(start, end)
	}

	wg.Wait()
}

// applyGeneral applies gate to this circuit.
func (c *Circuit) applyGeneral(op Gate, iregs ...int) {
	for basis, amp := range c.state.data {
		if amp == 0 {
			continue
		}
		// amp * |basis>
		// First, extract input qubits from basis
		// For example, if basis = |0101> and amp = 0, 2 => ibasis = |11>
		ibasis := 0
		for idx, val := range iregs {
			// Extract val-th bit from basis, plug it in to idx-th bit of ibasis.
			ibasis += ((basis >> val) & 1) << idx
		}
		// Apply gate to ibasis.
		// Note that ibasis is just a basis state.
		// This means that applying is taking columns from gate.
		for newibasis := 0; newibasis < (1 << len(iregs)); newibasis++ {
			newamp := op.data[newibasis][ibasis]
			if newamp == 0 {
				continue
			}

			newbasis := basis
			for idx, val := range iregs {
				bit := (newibasis >> idx) & 1
				newbasis = (newbasis | (1 << val)) - ((bit ^ 1) << val)
			}
			c.temp.data[newbasis] += amp * newamp
		}
	}

	c.state, c.temp = c.temp, c.state
	c.cleartemp()
}

// ApplyOracle applies the oracle f to circuit. Maps |x>_{iregs}|y>_{oregs} -> |x>_{iregs}|y^f(x)>_{oregs}.
// NOTE: This function DOES NOT check if oracle is unitary. Use at your own risk.
func (c *Circuit) ApplyOracle(oracle func(int) int, iregs []int, oregs []int) {
	if len(iregs) == 0 || len(oregs) == 0 {
		panic("Invalid input/output registers.")
	}

	if numbers.Min(iregs...) < 0 || numbers.Max(iregs...) >= c.Size() {
		panic("Register index out of range.")
	}

	if numbers.Min(oregs...) < 0 || numbers.Max(oregs...) >= c.Size() {
		panic("Register index out of range.")
	}

	if c.state.Dim() > c.Option.PARALLEL_THRESHOLD {
		c.applyOracleParallel(oracle, iregs, oregs)
		return
	}

	for basis, amp := range c.state.data {
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

		c.temp.data[newbasis] = amp
	}

	c.state, c.temp = c.temp, c.state
	c.cleartemp()
}

// applyOracleGeneral applies oracle as ApplyOracle with no parallelization.
func (c *Circuit) applyOracleParallel(oracle func(int) int, iregs, oregs []int) {
	jobsize := c.state.Dim()
	chunksize := jobsize / c.Option.GOROUTINE_CNT

	wg := &sync.WaitGroup{}
	wg.Add(c.Option.GOROUTINE_CNT)

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
		start := n * chunksize
		end := start + chunksize

		if n == c.Option.GOROUTINE_CNT-1 {
			end += jobsize % c.Option.GOROUTINE_CNT
		}

		go func(start, end int) {
			defer wg.Done()

			for basis := start; basis < end; basis++ {
				amp := c.state.data[basis]
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

				c.temp.data[newbasis] = amp
			}
		}(start, end)
	}

	wg.Wait()

	c.state, c.temp = c.temp, c.state
	c.cleartemp()
}

// Control.

// Used for calculating control bits in control-functions.
func checkControlBit(n int, cregs []int) bool {
	res := 0

	for _, idx := range cregs {
		res ^= (n >> idx) & 1
	}

	return res == 1
}

// Control applies controlled gate.
func (c *Circuit) Control(op Gate, cregs, iregs []int) {
	if len(iregs) == 0 {
		panic("At least one input registers required.")
	}

	if len(iregs) != op.Size() {
		panic("Operator size does not match input registers.")
	}

	if numbers.Min(cregs...) < 0 || numbers.Max(cregs...) > c.Size() {
		panic("Registers out of range.")
	}

	if numbers.Min(iregs...) < 0 || numbers.Max(iregs...) > c.Size() {
		panic("Registers out of range.")
	}

	// Special treatment for one and two qubit gates.
	if c.Size() > c.Option.PARALLEL_THRESHOLD {
		switch len(iregs) {
		case 1:
			c.controlOneParallel(op, cregs, iregs[0])
			return
		case 2:
			c.controlTwoParallel(op, cregs, iregs[0], iregs[1])
			return
		}
	} else {
		switch len(iregs) {
		case 1:
			c.controlOne(op, cregs, iregs[0])
			return
		case 2:
			c.controlTwo(op, cregs, iregs[0], iregs[1])
			return
		}
	}

	c.controlGeneral(op, cregs, iregs)
}

// controlOne applies one qubit controlled gate.
func (c *Circuit) controlOne(op Gate, cregs []int, i int) {
	lo := 1 << i

	for n := 0; n < c.state.Dim()/2; n++ {
		n0 := ((n >> i) << (i + 1)) + (n % lo)
		n1 := n0 | lo

		if !checkControlBit(n0, cregs) || !checkControlBit(n1, cregs) {
			continue
		}

		a0 := c.state.data[n0]
		a1 := c.state.data[n1]

		c.state.data[n0] = a0*op.data[0][0] + a1*op.data[0][1]
		c.state.data[n1] = a0*op.data[1][0] + a1*op.data[1][1]
	}
}

// controlOneParallel applies one qubit controlled gate with parallelization.
func (c *Circuit) controlOneParallel(op Gate, cregs []int, i int) {
	jobsize := c.state.Dim() / 2
	chunksize := jobsize / c.Option.GOROUTINE_CNT

	wg := &sync.WaitGroup{}
	wg.Add(c.Option.GOROUTINE_CNT)

	lo := 1 << i

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
		start := n * chunksize
		end := start + chunksize

		if n == c.Option.GOROUTINE_CNT-1 {
			end += jobsize % c.Option.GOROUTINE_CNT
		}

		go func(start, end int) {
			defer wg.Done()

			for n := start; n < end; n++ {
				n0 := ((n >> i) << (i + 1)) + (n % lo)
				n1 := n0 | lo

				if !checkControlBit(n0, cregs) || !checkControlBit(n1, cregs) {
					continue
				}

				a0 := c.state.data[n0]
				a1 := c.state.data[n1]

				c.state.data[n0] = a0*op.data[0][0] + a1*op.data[0][1]
				c.state.data[n1] = a0*op.data[1][0] + a1*op.data[1][1]
			}
		}(start, end)
	}

	wg.Wait()
}

// controlTwo applies two qubit controlled gate.
func (c *Circuit) controlTwo(op Gate, cregs []int, i0, i1 int) {
	if i0 == i1 {
		panic("Cannot apply gate to same registers.")
	}

	lo0 := 1 << i0
	lo1 := 1 << i1

	for n := 0; n < c.state.Dim()/4; n++ {
		n0 := ((n >> i0) << (i0 + 1)) + (n % lo0)

		n00 := ((n0 >> i1) << (i1 + 1)) + (n % lo1)
		n01 := n00 | lo0
		n10 := n00 | lo1
		n11 := n10 | lo0

		if !checkControlBit(n00, cregs) || !checkControlBit(n01, cregs) || !checkControlBit(n10, cregs) || !checkControlBit(n11, cregs) {
			continue
		}

		a00 := c.state.data[n00]
		a01 := c.state.data[n01]
		a10 := c.state.data[n10]
		a11 := c.state.data[n11]

		c.state.data[n00] = a00*op.data[0][0] + a01*op.data[0][1] + a10*op.data[0][2] + a11*op.data[0][3]
		c.state.data[n01] = a00*op.data[1][0] + a01*op.data[1][1] + a10*op.data[1][2] + a11*op.data[1][3]
		c.state.data[n10] = a00*op.data[2][0] + a01*op.data[2][1] + a10*op.data[2][2] + a11*op.data[2][3]
		c.state.data[n11] = a00*op.data[3][0] + a01*op.data[3][1] + a10*op.data[3][2] + a11*op.data[3][3]
	}
}

// controlTwoParallel applies two qubit controlled gate. with parallelizaition.
func (c *Circuit) controlTwoParallel(op Gate, cregs []int, i0, i1 int) {
	jobsize := c.state.Dim() / 4
	chunksize := jobsize / c.Option.GOROUTINE_CNT

	wg := &sync.WaitGroup{}
	wg.Add(c.Option.GOROUTINE_CNT)

	lo0 := 1 << i0
	lo1 := 1 << i1

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
		start := n * chunksize
		end := start + chunksize

		if n == c.Option.GOROUTINE_CNT-1 {
			end += jobsize % c.Option.GOROUTINE_CNT
		}

		go func(start, end int) {
			defer wg.Done()

			for n := start; n < end; n++ {
				n0 := ((n >> i0) << (i0 + 1)) + (n % lo0)

				n00 := ((n0 >> i1) << (i1 + 1)) + (n % lo1)
				n01 := n00 | lo0
				n10 := n00 | lo1
				n11 := n10 | lo0

				if !checkControlBit(n00, cregs) || !checkControlBit(n01, cregs) || !checkControlBit(n10, cregs) || !checkControlBit(n11, cregs) {
					continue
				}

				a00 := c.state.data[n00]
				a01 := c.state.data[n01]
				a10 := c.state.data[n10]
				a11 := c.state.data[n11]

				c.state.data[n00] = a00*op.data[0][0] + a01*op.data[0][1] + a10*op.data[0][2] + a11*op.data[0][3]
				c.state.data[n01] = a00*op.data[1][0] + a01*op.data[1][1] + a10*op.data[1][2] + a11*op.data[1][3]
				c.state.data[n10] = a00*op.data[2][0] + a01*op.data[2][1] + a10*op.data[2][2] + a11*op.data[2][3]
				c.state.data[n11] = a00*op.data[3][0] + a01*op.data[3][1] + a10*op.data[3][2] + a11*op.data[3][3]
			}
		}(start, end)
	}

	wg.Wait()
}

// controlGeneral applies controlled gate to this circuit.
func (c *Circuit) controlGeneral(op Gate, cregs, iregs []int) {
	for basis, amp := range c.state.data {
		if amp == 0 {
			continue
		}

		if !checkControlBit(basis, cregs) {
			continue
		}

		ibasis := 0
		for idx, val := range iregs {
			ibasis += ((basis >> val) & 1) << idx
		}

		for newibasis := 0; newibasis < (1 << len(iregs)); newibasis++ {
			newamp := op.data[newibasis][ibasis]
			if newamp == 0 {
				continue
			}

			newbasis := basis
			for idx, val := range iregs {
				bit := (newibasis >> idx) & 1
				newbasis = (newbasis | (1 << val)) - ((bit ^ 1) << val)
			}
			c.temp.data[newbasis] += amp * newamp
		}
	}

	c.state, c.temp = c.temp, c.state
	c.cleartemp()
}

// Misc Gates.

// Swap swaps two qubit.
func (c *Circuit) Swap(i0, i1 int) {
	if i0 < 0 || i0 >= c.Size() || i1 < 0 || i1 >= c.Size() {
		panic("Register index out of range.")
	}

	if i0 == i1 {
		panic("Swapping same registers.")
	}

	if c.Size() == 2 {
		c.state.data[0b01], c.state.data[0b10] = c.state.data[0b10], c.state.data[0b01]
		return
	}

	if c.Size() > c.Option.PARALLEL_THRESHOLD {
		c.swapParallel(i0, i1)
		return
	}

	lo0 := 1 << i0
	lo1 := 1 << i1

	for n := 0; n < c.state.Dim()/4; n++ {
		n0 := ((n >> i0) << (i0 + 1)) + (n % lo0)
		n00 := ((n0 >> i1) << (i1 + 1)) + (n % lo1)

		n01 := n00 | lo0
		n10 := n00 | lo1

		c.state.data[n01], c.state.data[n10] = c.state.data[n10], c.state.data[n01]
	}
}

// swapParallel swaps to qubit with parallelization.
func (c *Circuit) swapParallel(i0, i1 int) {
	jobsize := c.state.Dim() / 4
	chunksize := jobsize / c.Option.GOROUTINE_CNT

	wg := &sync.WaitGroup{}
	wg.Add(c.Option.GOROUTINE_CNT)

	lo0 := 1 << i0
	lo1 := 1 << i1

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
		start := n * chunksize
		end := start + chunksize

		if n == c.Option.GOROUTINE_CNT-1 {
			end += jobsize % c.Option.GOROUTINE_CNT
		}

		go func(start, end int) {
			defer wg.Done()

			for n := start; n < end; n++ {
				n0 := ((n >> i0) << (i0 + 1)) + (n % lo0)
				n00 := ((n0 >> i1) << (i1 + 1)) + (n % lo1)

				n01 := n00 | lo0
				n10 := n00 | lo1

				c.state.data[n01], c.state.data[n10] = c.state.data[n10], c.state.data[n01]
			}
		}(start, end)
	}

	wg.Wait()
}

// QFT applies QFT to [start, end).
func (c *Circuit) QFT(start, end int) {
	if start < 0 || end > c.Size() {
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
		c.H(i)
		for j := start; j < i; j++ {
			c.Control(P(phis[i-j]), []int{j}, []int{i})
		}
	}

	for i, j := start, end-1; i < j; i, j = i+1, j-1 {
		c.Swap(i, j)
	}

}

// InvQFT applies Inverse QFT to [start, end).
func (c *Circuit) InvQFT(start, end int) {
	if start < 0 || end > c.Size() {
		panic("Index out of range.")
	}

	if start >= end {
		panic("Invalid start / end parameters.")
	}

	for i, j := start, end-1; i < j; i, j = i+1, j-1 {
		c.Swap(i, j)
	}

	phis := make([]float64, end-start)

	for i := range phis {
		phis[i] = -math.Pi / math.Pow(2.0, float64(i))
	}

	for i := start; i < end; i++ {
		for j := start; j < i; j++ {
			c.Control(P(phis[i-j]), []int{j}, []int{i})
		}
		c.H(i)
	}
}

// Measure measures qubits.
func (c *Circuit) Measure(iregs ...int) int {
	iregs_s := make([]int, len(iregs))
	copy(iregs_s, iregs)
	sort.Ints(iregs_s)

	if iregs_s[0] < 0 || iregs_s[len(iregs_s)-1] > c.Size() {
		panic("Register index out of range.")
	}

	probs := make([]float64, 1<<len(iregs))

	for n, amp := range c.state.data {
		if amp == 0 {
			continue
		}
		o := 0
		for i, q := range iregs_s {
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

	for n, amp := range c.state.data {
		if amp == 0 {
			continue
		}

		has_output := true
		for i, q := range iregs_s {
			if (n>>q)&1 != (output>>i)&1 {
				c.state.data[n] = 0
				has_output = false
			}
		}
		if has_output {
			c.state.data[n] /= s
		}
	}

	return output
}

// String implements the Stringer interface.
func (q Circuit) String() string {
	return q.state.String()
}
