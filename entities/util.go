package entities

import neonMath "Neon/engine/math"

// simple utility functions

func unsetVec(slice []neonMath.Vector2D, s int) []neonMath.Vector2D {
	return append(slice[:s], slice[s+1:]...)
}

func unsetInt(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func del(slice []int, s int) []int {
	id := -1
	for k, v := range slice {
		if v == s {
			id = k
			break
		}
	}
	if id == -1 {
		return slice
	}

	return unsetInt(slice, id)
}

// swaps the target value with the result in the slice
func swap(slice []int, target int, result int) []int {
	for i, x := range slice {
		if x == target {
			slice[i] = result
			continue
		}
	}
	return slice
}
