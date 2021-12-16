# QSim

QSim is a quantum computing simulator written in pure go. Currently it supports up to 32 qubits, offering optimizations for one and two qubit gates using dedicated functions and parallelization. Applying one qubit gate to `n`-qubit state takes around `O(2^n)`.

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

## Benchmark
All tests are done in Mac mini with M1.
```
goos: darwin
goarch: arm64
pkg: github.com/sp301415/qsim
BenchmarkTensorApply-8    	1000000000	         0.02787 ns/op	       0 B/op	       0 allocs/op
BenchmarkApply-8           	1000000000	         0.0000753 ns/op	       0 B/op	       0 allocs/op
BenchmarkApplyParallel-8   	1000000000	         0.0001435 ns/op	       0 B/op	       0 allocs/op
```
Benchmark applies bunch of gates to 10 qubit circuit. `BenchmarkTensorApply` tensor products gates first, then applies it to the circuit. `BenchmarkApply` and `BenchmarkApplyParallel` repeatedly applies one qubit gate to each qubit. We can see that QSim's `Apply` is more than 3500 times faster than the naive method. 

Interestingly, using parallel computation takes litte more time, possibly because of the overhead. You can change this behavior by setting `Circuit.Option.PARALLEL_THRESHOLD`. The default value is 10.