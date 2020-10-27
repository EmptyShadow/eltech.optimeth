package harmonysearch

import (
	"context"

	"github.com/EmptyShadow/eltech.ai"
)

const (
	DefaultMemorySize                        = 50
	DefaultNumberOfImprovisations            = 100
	DefaultMinObjectValue                    = -100.0
	DefaultMaxObjectValue                    = 100.0
	DefaultProbabilityToTakeFromMemory       = 0.5
	DefaultProbabilityToApplyPitchAdjustment = 0.5
	DefaultPitchAdjustingRateWidth           = 1.0
	DefaultEps                               = 1e-6
)

var (
	DefaultDomainOfDefinitionFunc, _             = ai.SingleDomainOfDefinition(DefaultMinObjectValue, DefaultMaxObjectValue)
	DefaultProbabilityToTakeFromMemoryFunc       = ai.StaticProbability(DefaultProbabilityToTakeFromMemory)
	DefaultProbabilityToApplyPitchAdjustmentFunc = ai.StaticProbability(DefaultProbabilityToApplyPitchAdjustment)
	DefaultPitchAdjustingRateWidthFunc           = ai.StaticWidth(DefaultPitchAdjustingRateWidth)
)

type Opt func(opts *opts)

type opts struct {
	ctx context.Context

	// функция для получения области определения переменной.
	domainOfDefinition ai.DomainOfDefinitionFunc
	// функция для получения вероятности выбора гармоники из memory.
	probabilityToTakeFromMemory ai.ProbabilityFunc
	// функция для получения вероятности выполнения шага гармоники.
	probabilityToApplyPitchAdjustment ai.ProbabilityFunc
	// функция для получения ширины с которой может меняться гармоника.
	pitchAdjustingRateWidth ai.StepWidth

	memorySize             int     // размер памяти.
	numberOfObjects        int     // количество объектов.
	numberOfImprovisations int     // количество импровизаций.
	eps                    float64 // приближение.
	isFindingMin           bool    // флаг поиска минимума, а не максимума по умолчанию.
}

func Context(ctx context.Context) Opt {
	return func(opts *opts) {
		opts.ctx = ctx
	}
}

// DomainOfDefinition функция для получения области определения переменной.
func DomainOfDefinition(r ai.DomainOfDefinitionFunc) Opt {
	return func(opts *opts) {
		opts.domainOfDefinition = r
	}
}

// ProbabilityToTakeFromMemory функция для получения вероятности выбора гармоники из memory.
func ProbabilityToTakeFromMemory(cr ai.ProbabilityFunc) Opt {
	return func(opts *opts) {
		opts.probabilityToTakeFromMemory = cr
	}
}

// ProbabilityToApplyPitchAdjustment функция для получения вероятности выполнения шага гармоники.
func ProbabilityToApplyPitchAdjustment(par ai.ProbabilityFunc) Opt {
	return func(opts *opts) {
		opts.probabilityToApplyPitchAdjustment = par
	}
}

// PitchAdjustingRateWidth функция для получения ширины с которой может меняться гармоника.
func PitchAdjustingRateWidth(parw ai.StepWidth) Opt {
	return func(opts *opts) {
		opts.pitchAdjustingRateWidth = parw
	}
}

// MemorySize размер памяти.
func MemorySize(ms int) Opt {
	return func(opts *opts) {
		opts.memorySize = ms
	}
}

// NumberOfImprovisations количество импровизаций.
func NumberOfImprovisations(n int) Opt {
	return func(opts *opts) {
		opts.numberOfImprovisations = n
	}
}

// Eps погрешность приближения.
func Eps(eps float64) Opt {
	return func(opts *opts) {
		opts.eps = eps
	}
}

// FindMin флаг поиска минимума, а не максимума по умолчанию.
func FindMin() Opt {
	return func(opts *opts) {
		opts.isFindingMin = true
	}
}

func defaultOpts(numberOfObjects int) *opts {
	return &opts{
		ctx:                               context.Background(),
		domainOfDefinition:                DefaultDomainOfDefinitionFunc,
		probabilityToTakeFromMemory:       DefaultProbabilityToTakeFromMemoryFunc,
		probabilityToApplyPitchAdjustment: DefaultProbabilityToApplyPitchAdjustmentFunc,
		pitchAdjustingRateWidth:           DefaultPitchAdjustingRateWidthFunc,
		memorySize:                        DefaultMemorySize,
		numberOfObjects:                   numberOfObjects,
		numberOfImprovisations:            DefaultNumberOfImprovisations,
		eps:                               DefaultEps,
	}
}
