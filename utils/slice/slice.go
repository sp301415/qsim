// Package slice provides helper functions for silces.
package slice

// Sequence returns slices containing [a, b).
func Sequence(a, b int) []int {
	res := make([]int, b-a)
	for i := range res {
		res[i] = i + a
	}

	return res
}
