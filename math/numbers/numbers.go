// Package numbers provides various integer algorithms.
// Frankly, Go should already have them. :(
package numbers

// Returns power of two int, a ** b.
func Pow(a, b int) int {
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

// Returns modulo power of two int, a ** b % c.
func PowMod(a, b, c int) int {
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

// Returns GCD of two int.
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Returns binary length of number.
func BitLength(n int) int {
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

// Returns the smallest number.
func Min(ns ...int) int {
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

// Returns the largest number.
func Max(ns ...int) int {
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
