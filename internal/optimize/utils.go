package optimize

import (
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/optimize"
)

func funcEvaluation(opr chan<- optimize.Task, x []float64) {
	opr <- optimize.Task{
		Op:       optimize.FuncEvaluation,
		Location: &optimize.Location{X: x},
	}
}

func majorIteration(opr chan<- optimize.Task, x []float64, f float64) {
	opr <- optimize.Task{
		ID:       0,
		Op:       optimize.MajorIteration,
		Location: &optimize.Location{X: x, F: f},
	}
}

type FunctionConverge struct {
	Absolute float64

	first bool
	bestF float64
	bestX []float64
}

func (c *FunctionConverge) Init(_ int) {
	c.first = true
}

func (c *FunctionConverge) Converged(loc *optimize.Location) optimize.Status {
	f := loc.F
	x := loc.X

	if c.first {
		c.bestF = f
		c.bestX = x
		c.first = false

		return optimize.NotTerminated
	}

	diffX := floats.Distance(c.bestX, x, 2)

	if f < c.bestF {
		c.bestF = f
		c.bestX = x
	}

	if diffX >= c.Absolute {
		return optimize.NotTerminated
	}

	return optimize.FunctionConvergence
}
