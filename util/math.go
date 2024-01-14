package util

import "golang.org/x/exp/constraints"

// Max returns the largest argument provided.
// https://stackoverflow.com/a/73243983
func Max[T constraints.Ordered](args ...T) T {
	if len(args) == 0 {
		return *new(T) // zero value of T
	}

	if isNan(args[0]) {
		return args[0]
	}

	max := args[0]
	for _, arg := range args[1:] {

		if isNan(arg) {
			return arg
		}

		if arg > max {
			max = arg
		}
	}
	return max
}

func isNan[T constraints.Ordered](arg T) bool {
	return arg != arg
}
