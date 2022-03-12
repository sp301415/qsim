// Package number provides various functions for integers.
package number

// Pow returns a ** b.
func Pow(a, b int) int {
	if a == 2 {
		return 1 << b
	}

	res := 1

	for b > 0 {
		if b&1 != 0 {
			res *= a
		}
		b >>= 1
		a *= a
	}

	return res
}

// PowMod returns a ** b mod c.
// Panics if c < 0.
func PowMod(a, b, c int) int {
	if c < 0 {
		panic("Negative modulo not allowed.")
	}

	res := 1
	a %= c

	for b > 0 {
		if b&1 != 0 {
			res = (res * a) % c
		}
		b >>= 1
		a = (a * a) % c
	}

	return res
}

// GCD returns greatest common divisor of a and b.
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// BitLen returns binary length of n.
// If n == 0, it returns 1. If n < 0, it panics.
func BitLen(n int) int {
	if n < 0 {
		panic("Negative integer not allowed.")
	}

	if n == 0 {
		return 1
	}

	len := 0
	for n > 0 {
		len += 1
		n >>= 1
	}

	return len
}

// Min returns the smallest integer.
func Min(ns ...int) int {
	if len(ns) == 0 {
		panic("At least one argument is required.")
	}

	if len(ns) == 1 {
		return ns[0]
	}

	min := ns[0]

	for i := 1; i < len(ns); i++ {
		if ns[i] < min {
			min = ns[i]
		}
	}

	return min
}

// Max returns the largest integer.
func Max(ns ...int) int {
	if len(ns) == 0 {
		panic("At least one argument is required.")
	}

	if len(ns) == 1 {
		return ns[0]
	}

	max := ns[0]

	for i := 1; i < len(ns); i++ {
		if ns[i] > max {
			max = ns[i]
		}
	}

	return max
}
