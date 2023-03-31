package math

import (
	"golang.org/x/exp/constraints"
)

// Min 取a、b中较小的对象
func Min[T constraints.Ordered](a, b T) T {
	if a <= b {
		return a
	}

	return b
}

// Max 取a、b中较大的对象
func Max[T constraints.Ordered](a, b T) T {
	if a >= b {
		return a
	}

	return b
}
