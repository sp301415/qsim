# QSim

QSim is a quantum computing simulator written in pure go. Currently it supports up to 24 qubits, offering optimizations for one and two qubit gates using dedicated functions and parallelization. Applying one qubit gate to `n`-qubit state takes around `O(2^n)`.

NOTE: All measurements are not random for now, for benchmarking purposes.

## Example
### GHZ
```go
// Prepare two qubit circuit.
c := qsim.NewCircuit(2)

// Apply Hadamard and CNOT gates.
c.H(0)
c.CX(0, 1)

fmt.Println(c)
// Output:
// [0]|00>: (0.707107+0.000000i)
// [3]|11>: (0.707107+0.000000i)

// Measure the first qubit.
c.Measure(0)

fmt.Println(c)
// Output:
// [3]|11>: (1.000000+0.000000i)
```

### Grover
```go 
// Grover's Algorithm with 4 qubits.

// Prepare four qubit circuit.
c := qsim.NewCircuit(4)

// Apply Hadamard Gate.
c.H(0, 1, 2, 3)

// Iterate.
n := 1 << c.Size()
r := math.Floor(math.Pi / 4 * math.Sqrt(float64(n)))
for i := 0; i < int(r); i++ {
    c.X(0, 1)
    c.H(0)
    c.Control(qsim.X(), []int{1, 2, 3}, []int{0})
    c.H(0)
    c.X(0, 1)

    c.H(0, 1, 2, 3)
    c.X(0, 1, 2, 3)
    c.H(0)
    c.Control(qsim.X(), []int{1, 2, 3}, []int{0})
    c.H(0)
    c.X(0, 1, 2, 3)
    c.H(0, 1, 2, 3)
}

fmt.Println(c)
// Output:
// [ 0] |0000>: (0.050781+0.000000i)
// [ 1] |0001>: (0.050781+0.000000i)
// [ 2] |0010>: (0.050781+0.000000i)
// [ 3] |0011>: (0.050781+0.000000i)
// [ 4] |0100>: (0.050781+0.000000i)
// [ 5] |0101>: (0.050781+0.000000i)
// [ 6] |0110>: (0.050781+0.000000i)
// [ 7] |0111>: (0.050781+0.000000i)
// [ 8] |1000>: (0.050781+0.000000i)
// [ 9] |1001>: (0.050781+0.000000i)
// [10] |1010>: (0.050781+0.000000i)
// [11] |1011>: (0.050781+0.000000i)
// [12] |1100>: (-0.980469+0.000000i)
// [13] |1101>: (0.050781+0.000000i)
// [14] |1110>: (0.050781+0.000000i)
// [15] |1111>: (0.050781+0.000000i)
```