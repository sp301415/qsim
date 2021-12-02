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

// Returns integer log base 2.
func Log2(n int) int {
	res := 0
	n >>= 1

	for n > 0 {
		res++
		n >>= 1
	}

	return res
}

// Returns binary length of number
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

func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
