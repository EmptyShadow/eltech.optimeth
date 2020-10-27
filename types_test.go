package ai_test

import (
	ai "github.com/EmptyShadow/eltech.ai"
	"testing"
)

func TestRandInRangeFloat64(t *testing.T) {
	min := -3.0
	max := 3.0

	v := ai.RandInRangeFloat64(min, max)

	if max < v || v < min {
		t.Fatalf("failed %f < %f < %f", max, v, min)
	}
}
