package engine

// equalSet treats an array as a set and determines if two arrays are equivalent
func equalSet(a, b []int) bool {
	aExistsB := func(a, b []int) bool {
		x := make(map[int]bool)
		for _, v := range a {
			x[v] = true
		}
		for _, v := range b {
			if !x[v] {
				return false
			}
		}
		return true
	}

	return aExistsB(a, b) && aExistsB(b, a)
}
