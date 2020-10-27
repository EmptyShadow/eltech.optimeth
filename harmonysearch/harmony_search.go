// Пакет harmonysearch реализует алгоритм гармонического поиска.

package harmonysearch

import (
	ai "github.com/EmptyShadow/eltech.ai"
	"math/rand"
)

type Compositor struct {
	*opts
}

func NewCompositor(numberOfObjects int, opts ...Opt) *Compositor {
	_opts := defaultOpts(numberOfObjects)

	for _, opt := range opts {
		opt(_opts)
	}

	return &Compositor{opts: _opts}
}

func (c *Compositor) Improvisation(f ai.OptiFunc, start ai.Vector) (bestImprovised ai.Vector, err error) {
	m, err := c.initialization()
	if err != nil {
		return nil, err
	}

	currentImprovisation := 0

	bestImprovised = start
	bestValue := f(bestImprovised)

	for {
		select {
		case <-c.ctx.Done():
			return bestImprovised, c.ctx.Err()
		default:
		}

		improvised, err := c.improvisation(m)
		if err != nil {
			return nil, err
		}

		currentValue := f(improvised)
		norma := bestImprovised.Diff(improvised).Norma()

		if c.isSolutionBetter(currentValue, bestValue) {
			bestImprovised = improvised
			bestValue = currentValue
		}

		currentImprovisation++

		if currentImprovisation >= c.numberOfImprovisations || norma <= c.eps {
			return bestImprovised, nil
		}
	}
}

func (c *Compositor) initialization() (ai.Matrix, error) {
	return ai.NewMatrixWithInitFunc(c.numberOfObjects, c.memorySize, func(i, _ int) (float64, error) {
		return ai.VectorWithDomainOfDefinitionInitFunc(c.domainOfDefinition)(i)
	})
}

func (c *Compositor) improvisation(m ai.Matrix) (improvised ai.Vector, err error) {
	improvised = ai.NewVector(c.numberOfObjects)

	for j := 0; j < c.numberOfObjects; j++ {
		prob1 := rand.Float64()

		if prob1 >= c.probabilityToTakeFromMemory(j) {
			if improvised[j], err = c.randImprovisedElement(j); err != nil {
				return nil, err
			}

			continue
		}

		prod2 := rand.Float64()
		randDimensionIndex := rand.Intn(c.memorySize)

		if prod2 >= c.probabilityToApplyPitchAdjustment(j) {
			randDimension := m.Row(randDimensionIndex)
			improvised[j] = randDimension[j]

			continue
		}

		min, max, err := c.domainOfDefinition(j)
		if err != nil {
			return nil, err
		}

		pitchAdjusting := rand.Float64() * ai.RandInRangeFloat64(min, max) * c.pitchAdjustingRateWidth(j)

		if improvised[j]+pitchAdjusting <= max {
			improvised[j] += pitchAdjusting
		} else if improvised[j]-pitchAdjusting >= min {
			improvised[j] -= pitchAdjusting
		}
	}

	return improvised, nil
}

func (c *Compositor) randImprovisedElement(i int) (element float64, err error) {
	min, max, err := c.domainOfDefinition(i)
	if err != nil {
		return 0.0, err
	}

	return ai.RandInRangeFloat64(min, max), nil
}

func (c *Compositor) isSolutionBetter(currentValue, bestValue float64) bool {
	if c.isFindingMin {
		return currentValue < bestValue
	}

	return currentValue > bestValue
}
