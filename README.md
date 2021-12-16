# QSim

QSim is a quantum computing simulator written in pure go. Currently it supports up to 32 qubits, offering optimizations for one and two qubit gates. Applying one qubit gate to `n`-qubit state takes around `O(2^n)`.

NOTE: All measurements are not random for now, for benchmarking purposes. If you want real random, apply some seed to `rand`.

## Example
```go
// Two Qubit GHZ State.

// Prepare two qubit circuit.
c := qsim.NewCircuit(2)

// Apply Hadamard and CNOT gates.
c.H(0)
c.CX(0, 1)

fmt.Println(c.State)
// Output:
// |00>: (0.707107+0.000000i)
// |11>: (0.707107+0.000000i)

// Measure the first qubit.
c.Measure(0)

fmt.Println(c.State)
// Output:
// |11>: (1.000000+0.000000i)
```

