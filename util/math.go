package util

import "golang.org/x/exp/constraints"

// Max returns the largest argument provided.
func Max[T constraints.Ordered](args ...T) T {
	if len(args) == 0 {
		var zero T
		return zero
	}

	max := args[0]
	for _, arg := range args[1:] {
		if arg > max {
			max = arg
		}
	}
	return max
}

// Min returns the smallest argument provided.
func Min[T constraints.Ordered](args ...T) T {
	if len(args) == 0 {
		var zero T
		return zero
	}

	min := args[0]
	for _, arg := range args[1:] {
		if arg < min {
			min = arg
		}
	}
	return min
}
