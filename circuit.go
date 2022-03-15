package qsim

import (
	"math"
	"math/rand"
	"runtime"
	"sync"

	"github.com/sp301415/qsim/math/number"
	"github.com/sp301415/qsim/math/vec"
	"github.com/sp301415/qsim/utils/slice"
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
	if nbits < 0 || nbits > 24 {
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
func (c *Circuit) I(iregs ...int) {
	// Just Do Nothing. lol.
}

// Applies the X gate.
func (c *Circuit) X(iregs ...int) {
	if len(iregs) == 0 {
		panic("At least one input registers required.")
	}

	for _, i := range iregs {
		c.Apply(X(), i)
	}
}

// Applies the Y gate.
func (c *Circuit) Y(iregs ...int) {
	if len(iregs) == 0 {
		panic("At least one input registers required.")
	}

	for _, i := range iregs {
		c.Apply(Y(), i)
	}
}

// Applies the Z gate.
func (c *Circuit) Z(iregs ...int) {
	if len(iregs) == 0 {
		panic("At least one input registers required.")
	}

	for _, i := range iregs {
		c.Apply(Z(), i)
	}
}

// Applies the H gate.
func (c *Circuit) H(iregs ...int) {
	if len(iregs) == 0 {
		panic("At least one input registers required.")
	}

	for _, i := range iregs {
		c.Apply(H(), i)
	}
}

// Applies the P gate.
func (c *Circuit) P(phi float64, iregs ...int) {
	if len(iregs) == 0 {
		panic("At least one input registers required.")
	}

	for _, i := range iregs {
		c.Apply(P(phi), i)
	}
}

// Applies the S gate.
func (c *Circuit) S(iregs ...int) {
	if len(iregs) == 0 {
		panic("At least one input registers required.")
	}

	for _, i := range iregs {
		c.Apply(S(), i)
	}
}

// Applies the T gate.
func (c *Circuit) T(iregs ...int) {
	if len(iregs) == 0 {
		panic("At least one input registers required.")
	}

	for _, i := range iregs {
		c.Apply(T(), i)
	}
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

	if number.Min(iregs...) < 0 || number.Max(iregs...) > c.Size() {
		panic("Registers out of range.")
	}

	if slice.HasDuplicate(iregs) {
		panic("Duplicate registers.")
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
	mask := (1 << i) - 1

	for n := 0; n < c.state.Dim()/2; n++ {
		// n0 = XXX0XXX, n1 = XXX1XXX
		n0 := ((n & ^mask) << 1) + (n & mask)
		n1 := n0 | (mask + 1)

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

	chunkidx := make([]int, c.Option.GOROUTINE_CNT+1)
	for i := 0; i < c.Option.GOROUTINE_CNT; i++ {
		chunkidx[i] = i * chunksize
	}
	chunkidx[c.Option.GOROUTINE_CNT] = jobsize

	var wg sync.WaitGroup
	wg.Add(c.Option.GOROUTINE_CNT)

	mask := (1 << i) - 1

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
		go func(start, end int) {
			defer wg.Done()

			for n := start; n < end; n++ {
				n0 := ((n & ^mask) << 1) + (n & mask)
				n1 := n0 | (mask + 1)

				a0 := c.state.data[n0]
				a1 := c.state.data[n1]

				c.state.data[n0] = a0*op.data[0][0] + a1*op.data[0][1]
				c.state.data[n1] = a0*op.data[1][0] + a1*op.data[1][1]
			}
		}(chunkidx[n], chunkidx[n+1])
	}

	wg.Wait()
}

// applyTwo applies two qubit gate.
func (c *Circuit) applyTwo(op Gate, i0, i1 int) {
	if i0 == i1 {
		panic("Cannot apply gate to same registers.")
	}

	if i0 > i1 {
		i0, i1 = i1, i0
		op.data[1], op.data[2] = op.data[2], op.data[1]
	}

	mask0 := (1 << i0) - 1
	mask1 := (1 << i1) - 1

	for n := 0; n < c.state.Dim()/4; n++ {
		// n00 = XXX0(i1)XXX0(i0)XXX
		// n01 = XXX0(i1)XXX1(i0)XXX
		// ...

		t := ((n & ^mask0) << 1) + (n & mask0)

		n00 := ((t & ^mask1) << 1) + (t & mask1)
		n01 := n00 | (mask0 + 1)
		n10 := n00 | (mask1 + 1)
		n11 := n10 | (mask0 + 1)

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

	chunkidx := make([]int, c.Option.GOROUTINE_CNT+1)
	for i := 0; i < c.Option.GOROUTINE_CNT; i++ {
		chunkidx[i] = i * chunksize
	}
	chunkidx[c.Option.GOROUTINE_CNT] = jobsize

	var wg sync.WaitGroup
	wg.Add(c.Option.GOROUTINE_CNT)

	if i0 > i1 {
		i0, i1 = i1, i0
		op.data[1], op.data[2] = op.data[2], op.data[1]
	}

	mask0 := (1 << i0) - 1
	mask1 := (1 << i1) - 1

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
		go func(start, end int) {
			defer wg.Done()

			for n := start; n < end; n++ {
				t := ((n & ^mask0) << 1) + (n & mask0)

				n00 := ((t & ^mask1) << 1) + (t & mask1)
				n01 := n00 | (mask0 + 1)
				n10 := n00 | (mask1 + 1)
				n11 := n10 | (mask0 + 1)

				a00 := c.state.data[n00]
				a01 := c.state.data[n01]
				a10 := c.state.data[n10]
				a11 := c.state.data[n11]

				c.state.data[n00] = a00*op.data[0][0] + a01*op.data[0][1] + a10*op.data[0][2] + a11*op.data[0][3]
				c.state.data[n01] = a00*op.data[1][0] + a01*op.data[1][1] + a10*op.data[1][2] + a11*op.data[1][3]
				c.state.data[n10] = a00*op.data[2][0] + a01*op.data[2][1] + a10*op.data[2][2] + a11*op.data[2][3]
				c.state.data[n11] = a00*op.data[3][0] + a01*op.data[3][1] + a10*op.data[3][2] + a11*op.data[3][3]
			}
		}(chunkidx[n], chunkidx[n+1])
	}

	wg.Wait()
}

// applyGeneral applies gate to this circuit.
func (c *Circuit) applyGeneral(op Gate, iregs ...int) {
	c.cleartemp()
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
}

// ApplyOracle applies the oracle f to circuit. Maps |x>_{iregs}|y>_{oregs} -> |x>_{iregs}|y^f(x)>_{oregs}.
// NOTE: This function DOES NOT check if oracle is unitary. Use at your own risk.
func (c *Circuit) ApplyOracle(oracle func(int) int, iregs []int, oregs []int) {
	if len(iregs) == 0 || len(oregs) == 0 {
		panic("Invalid input/output registers.")
	}

	if number.Min(iregs...) < 0 || number.Max(iregs...) >= c.Size() {
		panic("Register index out of range.")
	}

	if number.Min(oregs...) < 0 || number.Max(oregs...) >= c.Size() {
		panic("Register index out of range.")
	}

	if slice.HasDuplicate(iregs) || slice.HasDuplicate(oregs) {
		panic("Duplicate registers.")
	}

	if slice.HasCommon(iregs, oregs) {
		panic("Duplicate registers.")
	}

	if c.state.Dim() > c.Option.PARALLEL_THRESHOLD {
		c.applyOracleParallel(oracle, iregs, oregs)
		return
	}

	c.cleartemp()
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
			newbasis ^= ((output >> idx) & 1) << val
		}

		c.temp.data[newbasis] = amp
	}

	c.state, c.temp = c.temp, c.state
}

// applyOracleParallel applies oracle with parallelizaiton.
func (c *Circuit) applyOracleParallel(oracle func(int) int, iregs, oregs []int) {
	c.cleartemp()

	jobsize := c.state.Dim()
	chunksize := jobsize / c.Option.GOROUTINE_CNT

	chunkidx := make([]int, c.Option.GOROUTINE_CNT+1)
	for i := 0; i < c.Option.GOROUTINE_CNT; i++ {
		chunkidx[i] = i * chunksize
	}
	chunkidx[c.Option.GOROUTINE_CNT] = jobsize

	var wg sync.WaitGroup
	wg.Add(c.Option.GOROUTINE_CNT)

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
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
					newbasis ^= ((output >> idx) & 1) << val
				}

				c.temp.data[newbasis] = amp
			}
		}(chunkidx[n], chunkidx[n+1])
	}

	wg.Wait()

	c.state, c.temp = c.temp, c.state
}

// Control.

// Used for calculating control bits in control-functions.
func checkControlBit(n int, cregs []int) bool {
	for _, idx := range cregs {
		if (n>>idx)&1 == 0 {
			return false
		}
	}

	return true
}

// Control applies controlled gate.
func (c *Circuit) Control(op Gate, cregs, iregs []int) {
	if len(iregs) == 0 {
		panic("At least one input registers required.")
	}

	if len(iregs) != op.Size() {
		panic("Operator size does not match input registers.")
	}

	if number.Min(cregs...) < 0 || number.Max(cregs...) > c.Size() {
		panic("Registers out of range.")
	}

	if number.Min(iregs...) < 0 || number.Max(iregs...) > c.Size() {
		panic("Registers out of range.")
	}

	if slice.HasDuplicate(cregs) || slice.HasDuplicate(iregs) {
		panic("Duplicate registers.")
	}

	if slice.HasCommon(cregs, iregs) {
		panic("Duplicate registers.")
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
	mask := (1 << i) - 1

	for n := 0; n < c.state.Dim()/2; n++ {
		n0 := ((n & ^mask) << 1) + (n & mask)
		n1 := n0 | (mask + 1)

		if !checkControlBit(n0, cregs) {
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

	chunkidx := make([]int, c.Option.GOROUTINE_CNT+1)
	for i := 0; i < c.Option.GOROUTINE_CNT; i++ {
		chunkidx[i] = i * chunksize
	}
	chunkidx[c.Option.GOROUTINE_CNT] = jobsize

	var wg sync.WaitGroup
	wg.Add(c.Option.GOROUTINE_CNT)

	mask := (1 << i) - 1

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
		go func(start, end int) {
			defer wg.Done()

			for n := start; n < end; n++ {
				n0 := ((n & ^mask) << 1) + (n & mask)
				n1 := n0 | (mask + 1)

				if !checkControlBit(n0, cregs) {
					continue
				}

				a0 := c.state.data[n0]
				a1 := c.state.data[n1]

				c.state.data[n0] = a0*op.data[0][0] + a1*op.data[0][1]
				c.state.data[n1] = a0*op.data[1][0] + a1*op.data[1][1]
			}
		}(chunkidx[n], chunkidx[n+1])
	}

	wg.Wait()
}

// controlTwo applies two qubit controlled gate.
func (c *Circuit) controlTwo(op Gate, cregs []int, i0, i1 int) {
	if i0 == i1 {
		panic("Cannot apply gate to same registers.")
	}

	if i0 > i1 {
		i0, i1 = i1, i0
		op.data[1], op.data[2] = op.data[2], op.data[1]
	}

	mask0 := (1 << i0) - 1
	mask1 := (1 << i1) - 1

	for n := 0; n < c.state.Dim()/4; n++ {
		t := ((n & ^mask0) << 1) + (n & mask0)

		n00 := ((t & ^mask1) << 1) + (t & mask1)
		n01 := n00 | (mask0 + 1)
		n10 := n00 | (mask1 + 1)
		n11 := n10 | (mask0 + 1)

		if !checkControlBit(n00, cregs) {
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

	chunkidx := make([]int, c.Option.GOROUTINE_CNT+1)
	for i := 0; i < c.Option.GOROUTINE_CNT; i++ {
		chunkidx[i] = i * chunksize
	}
	chunkidx[c.Option.GOROUTINE_CNT] = jobsize

	var wg sync.WaitGroup
	wg.Add(c.Option.GOROUTINE_CNT)

	if i0 > i1 {
		i0, i1 = i1, i0
		op.data[1], op.data[2] = op.data[2], op.data[1]
	}

	mask0 := (1 << i0) - 1
	mask1 := (1 << i1) - 1

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
		go func(start, end int) {
			defer wg.Done()

			for n := start; n < end; n++ {
				t := ((n & ^mask0) << 1) + (n & mask0)

				n00 := ((t & ^mask1) << 1) + (t & mask1)
				n01 := n00 | (mask0 + 1)
				n10 := n00 | (mask1 + 1)
				n11 := n10 | (mask0 + 1)

				if !checkControlBit(n00, cregs) {
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
		}(chunkidx[n], chunkidx[n+1])
	}

	wg.Wait()
}

// controlGeneral applies controlled gate to this circuit.
func (c *Circuit) controlGeneral(op Gate, cregs, iregs []int) {
	c.cleartemp()
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
				newbasis = (newbasis | (1 << val)) - ((((newibasis >> idx) & 1) ^ 1) << val)
			}
			c.temp.data[newbasis] += amp * newamp
		}
	}

	c.state, c.temp = c.temp, c.state
}

// Misc Gates.

// Swap swaps two qubit.
func (c *Circuit) Swap(i0, i1 int) {
	if i0 < 0 || i0 >= c.Size() || i1 < 0 || i1 >= c.Size() {
		panic("Register index out of range.")
	}

	if i0 == i1 {
		panic("Duplicate registers.")
	}

	if c.Size() == 2 {
		c.state.data[0b01], c.state.data[0b10] = c.state.data[0b10], c.state.data[0b01]
		return
	}

	if c.Size() > c.Option.PARALLEL_THRESHOLD {
		c.swapParallel(i0, i1)
		return
	}

	if i0 > i1 {
		i0, i1 = i1, i0
	}

	mask0 := (1 << i0) - 1
	mask1 := (1 << i1) - 1

	for n := 0; n < c.state.Dim()/4; n++ {
		t := ((n & ^mask0) << 1) + (n & mask0)

		n00 := ((t & ^mask1) << 1) + (t & mask1)
		n01 := n00 | (mask0 + 1)
		n10 := n00 | (mask1 + 1)

		c.state.data[n01], c.state.data[n10] = c.state.data[n10], c.state.data[n01]
	}
}

// swapParallel swaps to qubit with parallelization.
func (c *Circuit) swapParallel(i0, i1 int) {
	jobsize := c.state.Dim() / 4
	chunksize := jobsize / c.Option.GOROUTINE_CNT

	chunkidx := make([]int, c.Option.GOROUTINE_CNT+1)
	for i := 0; i < c.Option.GOROUTINE_CNT; i++ {
		chunkidx[i] = i * chunksize
	}
	chunkidx[c.Option.GOROUTINE_CNT] = jobsize

	var wg sync.WaitGroup
	wg.Add(c.Option.GOROUTINE_CNT)

	if i0 > i1 {
		i0, i1 = i1, i0
	}

	mask0 := (1 << i0) - 1
	mask1 := (1 << i1) - 1

	for n := 0; n < c.Option.GOROUTINE_CNT; n++ {
		go func(start, end int) {
			defer wg.Done()

			for n := start; n < end; n++ {
				t := ((n & ^mask0) << 1) + (n & mask0)

				n00 := ((t & ^mask1) << 1) + (t & mask1)
				n01 := n00 | (mask0 + 1)
				n10 := n00 | (mask1 + 1)

				c.state.data[n01], c.state.data[n10] = c.state.data[n10], c.state.data[n01]
			}
		}(chunkidx[n], chunkidx[n+1])
	}

	wg.Wait()
}

// QFT applies QFT.
func (c *Circuit) QFT(iregs ...int) {
	phis := make([]float64, len(iregs))
	for i := range phis {
		phis[i] = math.Pi / float64(number.Pow(2, i))
	}

	for i := len(iregs) - 1; i >= 0; i-- {
		c.H(iregs[i])
		for j := 0; j < i; j++ {
			c.Control(P(phis[i-j]), []int{iregs[j]}, []int{iregs[i]})
		}
	}

	for i, j := 0, len(iregs)-1; i < j; i, j = i+1, j-1 {
		c.Swap(iregs[i], iregs[j])
	}
}

// InvQFT applies Inverse QFT.
func (c *Circuit) InvQFT(iregs ...int) {
	for i, j := 0, len(iregs)-1; i < j; i, j = i+1, j-1 {
		c.Swap(iregs[i], iregs[j])
	}

	phis := make([]float64, len(iregs))

	for i := range phis {
		phis[i] = -math.Pi / float64(number.Pow(2, i))
	}

	for i := 0; i < len(iregs); i++ {
		for j := 0; j < i; j++ {
			c.Control(P(phis[i-j]), []int{iregs[j]}, []int{iregs[i]})
		}
		c.H(iregs[i])
	}
}

// Measure measures qubits.
func (c *Circuit) Measure(iregs ...int) int {
	if number.Min(iregs...) < 0 || number.Max(iregs...) > c.Size() {
		panic("Register index out of range.")
	}

	if slice.HasDuplicate(iregs) {
		panic("Duplicate registers.")
	}

	probs := make([]float64, 1<<len(iregs))

	for n, amp := range c.state.data {
		if amp == 0 {
			continue
		}
		o := 0
		for i, q := range iregs {
			o += ((n >> q) & 1) << i
		}
		probs[o] += real(amp)*real(amp) + imag(amp)*imag(amp)
	}

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
		for i, q := range iregs {
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
