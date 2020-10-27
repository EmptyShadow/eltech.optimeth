package harmonysearch_test

import (
	ai "github.com/EmptyShadow/eltech.ai"
	"testing"

	"github.com/EmptyShadow/eltech.ai/harmonysearch"
	"github.com/EmptyShadow/eltech.ai/testdata"
	"github.com/stretchr/testify/assert"
)

type test struct {
	enabled bool
	f       testdata.TestFunc
	opts    opts
}

type opts struct {
	start                             ai.Vector
	probabilityToTakeFromMemory       ai.ProbabilityFunc
	probabilityToApplyPitchAdjustment ai.ProbabilityFunc
	pitchAdjustingRateWidth           ai.StepWidth
	eps                               float64
}

func Test_Compositor_Improvisation(t *testing.T) {
	asserting := assert.New(t)

	tests := []test{
		{
			enabled: false,
			f: testdata.MatiasFunc,
			opts: opts{
				start:                             ai.Vector{10.0, 10.0},
				probabilityToTakeFromMemory:       harmonysearch.DefaultProbabilityToTakeFromMemoryFunc,
				probabilityToApplyPitchAdjustment: harmonysearch.DefaultProbabilityToApplyPitchAdjustmentFunc,
				pitchAdjustingRateWidth:           harmonysearch.DefaultPitchAdjustingRateWidthFunc,
				eps:                               harmonysearch.DefaultEps,
			},
		},
		{
			enabled: true,
			f: testdata.SpheresFunc,
			opts: opts{
				start:                             ai.Vector{100.0, 100.0},
				probabilityToTakeFromMemory:       harmonysearch.DefaultProbabilityToTakeFromMemoryFunc,
				probabilityToApplyPitchAdjustment: harmonysearch.DefaultProbabilityToApplyPitchAdjustmentFunc,
				pitchAdjustingRateWidth:           harmonysearch.DefaultPitchAdjustingRateWidthFunc,
				eps:                               harmonysearch.DefaultEps,
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
				harmonysearch.PitchAdjustingRateWidth(test.opts.pitchAdjustingRateWidth),
				harmonysearch.Eps(test.opts.eps),
				harmonysearch.FindMin())

			solution, err := c.Improvisation(test.f.F, test.opts.start)
			asserting.NoError(err)
			asserting.Contains(test.f.Solutions, solution)
			asserting.Equal(test.f.F(solution), test.f.Min)
		})
	}
}
