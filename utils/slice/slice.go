// Package slice provides helper functions for silces.
package slice

// Range returns slices containing [a, b).
func Range(a, b int) []int {
	res := make([]int, b-a)

	for i := range res {
		res[i] = i + a
	}

	return res
}

// Contains returns if slice contains element.
func Contains[T comparable](s []T, a T) bool {
	for _, v := range s {
		if v == a {
			return true
		}
	}

	return false
}

// HasCommon returns if two array has common element.
func HasCommon[T comparable](s1, s2 []T) bool {
	for _, v := range s1 {
		if Contains(s2, v) {
			return true
		}
	}

	return false
}

// HasDuplicate returns if given slice contains duplicate elements.
func HasDuplicate[T comparable](s []T) bool {
	m := make(map[T]struct{})

	for _, v := range s {
		if _, e := m[v]; e {
			return true
		}
		m[v] = struct{}{}
	}

	return false
}
