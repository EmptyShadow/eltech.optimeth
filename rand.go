package ai

import "math/rand"

type RandFloat64Func func() (float64, error)

func RangeRandFloat64Func(min, max float64) RandFloat64Func {
	return func() (float64, error) {
		return min + rand.Float64()*(max-min), nil
	}
}
