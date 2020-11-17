package optimize_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/EmptyShadow/eltech.optimize/internal/functions"
	internaloptimize "github.com/EmptyShadow/eltech.optimize/internal/optimize"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/optimize"
)

type test struct {
	enabled bool
	name    string

	prob     optimize.Problem
	minimum  float64
	initX    [][]float64
	settings *optimize.Settings

	conf *internaloptimize.HSConfig
}

func TestHS_Run(t *testing.T) {
	tests := []test{
		{
			enabled: true,
			name:    "Levi13",
			prob:    functions.MustProblem(functions.Levi13, nil, nil),
			minimum: 0.0,
			initX: [][]float64{
				{10.0, 10.0}, {9.0, 9.0}, {2.0, 7.8},
			},
			settings: &optimize.Settings{
				Runtime: time.Minute,
				Recorder: optimize.NewPrinter(),
				Converger: &optimize.FunctionConverge{
					Absolute:   1e-2,
					Iterations: 100,
				},
				//MajorIterations: 1_000_000, // максимальное количество найденных лучших значений.
				//FuncEvaluations: 1_000_000, // максимально количество вычислений кондидатов.
			},
			conf: &internaloptimize.HSConfig{
				FD: functions.NewSingleFuncDomain(functions.VarDomain{
					Bottom: -10,
					Top:    10,
				}),
				MemorySize:                 internaloptimize.DefaultHSMemorySize,
				ProbToTakeFromMemory:       internaloptimize.DefaultProbToTakeFromMemory,
				ProbToApplyPitchAdjustment: internaloptimize.DefaultProbToApplyPitchAdjustment,
				MaxStep:                    internaloptimize.DefaultMaxStep,
			},
		},
		{
			enabled: true,
			name:    "Matias",
			prob:    functions.MustProblem(functions.Matias, nil, nil),
			minimum: 0.0,
			initX: [][]float64{
				{10.0, 10.0}, {9.0, 9.0}, {2.0, 7.8},
			},
			settings: &optimize.Settings{
				Runtime: time.Minute,
				Recorder: optimize.NewPrinter(),
				Converger: &optimize.FunctionConverge{
					Absolute:   1e-2,
					Iterations: 150,
				},
				//MajorIterations: 1_000_000, // максимальное количество найденных лучших значений.
				//FuncEvaluations: 1_000_000, // максимально количество вычислений кондидатов.
			},
			conf: &internaloptimize.HSConfig{
				FD: functions.NewSingleFuncDomain(functions.VarDomain{
					Bottom: -10,
					Top:    10,
				}),
				MemorySize:                 internaloptimize.DefaultHSMemorySize,
				ProbToTakeFromMemory:       internaloptimize.DefaultProbToTakeFromMemory,
				ProbToApplyPitchAdjustment: internaloptimize.DefaultProbToApplyPitchAdjustment,
				MaxStep:                    internaloptimize.DefaultMaxStep,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.enabled {
				t.Skip()

				return
			}

			for _, initX := range test.initX {
				t.Run(fmt.Sprint(initX), func(t *testing.T) {
					hs := internaloptimize.NewHS(test.conf)

					asserting := assert.New(t)

					result, err := optimize.Minimize(test.prob, initX, test.settings, hs)
					asserting.NoError(err)
					asserting.Contains([]optimize.Status{optimize.Success, optimize.FunctionConvergence,
						optimize.FunctionEvaluationLimit, optimize.RuntimeLimit},
						result.Status)
					diff := math.Abs(math.Abs(result.F)-math.Abs(test.minimum))
					t.Log(diff)
					asserting.True(diff <= test.settings.Converger.(*optimize.FunctionConverge).Absolute)
				})
			}
		})
	}
}
