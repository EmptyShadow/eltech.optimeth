package optimeth_test

import (
	"testing"

	optimeth "github.com/EmptyShadow/eltech.optimeth"
)

func TestRandInRangeFloat64(t *testing.T) {
	min := -3.0
	max := 3.0

	v := optimeth.RandInRangeFloat64(min, max)

	if max < v || v < min {
		t.Fatalf("failed %f < %f < %f", max, v, min)
	}
}
