package harmonysearch_test

import (
	"context"
	"fmt"
	ai "github.com/EmptyShadow/eltech.ai"
	"math"
	"testing"

	"github.com/EmptyShadow/eltech.ai/harmonysearch"
	"github.com/EmptyShadow/eltech.ai/testdata"
	"github.com/stretchr/testify/assert"
)

type test struct {
	enabled bool
	f       testdata.TestFunc
	// минимальное приближение, например, лучшее 0.0, а полученное 1e-8, что достаточно близко к 0.0, поэтому следует
	// поставить приближение 1e-7 или чуть больше.
	minApproximation float64
	opts             opts
}

type opts struct {
	starts                            []ai.Vector
	probabilityToTakeFromMemory       ai.ProbabilityFunc
	probabilityToApplyPitchAdjustment ai.ProbabilityFunc
	domainOfDefinitionStep            ai.DomainOfDefinitionFunc
	memorySize                        int
	numberOfImprovisations            int
	eps                               float64
}

func Test_Compositor_Improvisation(t *testing.T) {
	tests := []test{
		{
			enabled:          true,
			f:                testdata.MatiasFunc,
			minApproximation: 1e-6,
			opts: opts{
				starts:                            []ai.Vector{{10.0, 10.0}, {9.0, 9.0}, {2.0, 7.8}},
				probabilityToTakeFromMemory:       ai.StaticProbability(0.3),
				probabilityToApplyPitchAdjustment: ai.StaticProbability(0.9),
				domainOfDefinitionStep:            ai.MustSingleDomainOfDefinition(0.1, 1.0),
				memorySize:                        100,
				numberOfImprovisations:            100_000_000,
				eps:                               1e-3,
			},
		},
		{
			enabled:          true,
			f:                testdata.SpheresFunc,
			minApproximation: 1e-5,
			opts: opts{
				starts:                            []ai.Vector{{100.0, 100.0}, {-100.0, -100.0}, {2.0, -10.0}},
				probabilityToTakeFromMemory:       ai.StaticProbability(0.3),
				probabilityToApplyPitchAdjustment: ai.StaticProbability(0.9),
				domainOfDefinitionStep:            ai.MustSingleDomainOfDefinition(0.1, 2.0),
				memorySize:                        100,
				numberOfImprovisations:            100_000_000,
				eps:                               1e-3,
			},
		},
		{
			enabled:          true,
			f:                testdata.Levi13Func,
			minApproximation: 1e-3,
			opts: opts{
				starts:                            []ai.Vector{{10.0, 10.0}, {9.0, 9.0}, {2.0, 7.8}},
				probabilityToTakeFromMemory:       ai.StaticProbability(0.3),
				probabilityToApplyPitchAdjustment: ai.StaticProbability(0.9),
				domainOfDefinitionStep:            ai.MustSingleDomainOfDefinition(0.1, 2.0),
				memorySize:                        100,
				numberOfImprovisations:            100_000_000,
				eps:                               1e-3,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.f.Name, func(t *testing.T) {
			if !test.enabled {
				t.Skip()

				return
			}

			c := harmonysearch.NewCompositor(test.f.NumberOfVariables,
				harmonysearch.DomainOfDefinition(test.f.DomainOfDefinition),
				harmonysearch.ProbabilityToTakeFromMemory(test.opts.probabilityToTakeFromMemory),
				harmonysearch.ProbabilityToApplyPitchAdjustment(test.opts.probabilityToApplyPitchAdjustment),
				harmonysearch.DomainOfDefinitionStep(test.opts.domainOfDefinitionStep),
				harmonysearch.MemorySize(test.opts.memorySize),
				harmonysearch.NumberOfImprovisations(test.opts.numberOfImprovisations),
				harmonysearch.Eps(test.opts.eps),
				harmonysearch.FindMin(),
			)

			for _, testStart := range test.opts.starts {
				t.Run(fmt.Sprint(testStart), func(t *testing.T) {
					asserting := assert.New(t)
					solution, min, err := c.Improvisation(context.Background(), test.f.F, testStart)
					asserting.NoError(err)
					asserting.True(math.Abs(min-test.f.Min) <= test.minApproximation)

					t.Log(solution, min)
				})
			}
		})
	}
}
