// Пакет harmonysearch реализует алгоритм гармонического поиска.

package harmonysearch

import (
	"math"
	"math/rand"
)

// memory память гармоник.
type memory [][]float64

func (m memory) Row(index int) []float64 {
	c := make([]float64, len(m[index]))

	copy(m[index], c)

	return c
}

func (m memory) Column(index int) []float64 {
	c := make([]float64, len(m))

	for i := 0; i < len(m); i++ {
		c[i] = m[i][index]
	}

	return c
}

func randInRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

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

type OptiFunc func([]float64) float64

func (c *Compositor) Improvisation(f OptiFunc) (bestImprovised []float64, err error) {
	m, err := c.initialization()
	if err != nil {
		return nil, err
	}

	currentImprovisation := 0
	bestValue := c.defaultBestValue()

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

		if c.isSolutionBetter(currentValue, bestValue) {
			bestImprovised = improvised
		}

		currentImprovisation++

		if currentImprovisation >= c.numberOfImprovisations {
			return bestImprovised, nil
		}
	}
}

func (c *Compositor) initialization() (memory, error) {
	hm := make(memory, c.memorySize)

	for i := 0; i < c.memorySize; i++ {
		hm[i] = make([]float64, c.numberOfObjects)

		for j := 0; j < c.numberOfObjects; j++ {
			min, max, err := c.domainOfDefinition(i)
			if err != nil {
				return nil, err
			}

			hm[i][j] = randInRange(min, max)
		}
	}

	return hm, nil
}

func (c *Compositor) defaultBestValue() float64 {
	if c.isFindingMin {
		return math.MaxFloat64
	}

	return -math.MaxFloat64
}

func (c *Compositor) improvisation(m memory) (improvised []float64, err error) {
	improvised = make([]float64, c.numberOfObjects)

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

		pitchAdjusting := rand.Float64() * (max - min) * c.pitchAdjustingRateWidth(j)

		if improvised[j]+pitchAdjusting > max && improvised[j]-pitchAdjusting >= min {
			improvised[j] -= pitchAdjusting
		} else if improvised[j]-pitchAdjusting < min && improvised[j] <= max {
			improvised[j] += pitchAdjusting
		}
	}

	return improvised, nil
}

func (c *Compositor) randImprovisedElement(i int) (element float64, err error) {
	min, max, err := c.domainOfDefinition(i)
	if err != nil {
		return 0.0, err
	}

	return randInRange(min, max), nil
}

func (c *Compositor) isSolutionBetter(currentValue, bestValue float64) bool {
	if c.isFindingMin {
		return currentValue < bestValue
	}

	return currentValue > bestValue
}
