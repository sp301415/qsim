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
func Contains(s []int, a int) bool {
	for _, v := range s {
		if v == a {
			return true
		}
	}

	return false
}

// HasCommon returns if two array has common element.
func HasCommon(s1, s2 []int) bool {
	for _, v := range s1 {
		if Contains(s2, v) {
			return true
		}
	}

	return false
}

// HasDuplicate returns if given slice contains duplicate elements.
func HasDuplicate(s []int) bool {
	m := make(map[int]bool)

	for _, v := range s {
		if _, e := m[v]; e {
			return true
		}
		m[v] = true
	}

	return false
}
