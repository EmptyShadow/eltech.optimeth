package harmonysearch

import (
	"context"
	"errors"
	"math/rand"
)

const (
	defaultMemorySize                        = 50
	defaultNumberOfImprovisations            = 100
	defaultMinObjectValue                    = -100.0
	defaultMaxObjectValue                    = 100.0
	defaultProbabilityToTakeFromMemory       = 0.5
	defaultProbabilityToApplyPitchAdjustment = 0.5
	defaultPitchAdjustingRateWidth           = 1.0
)

var (
	defaultDomainOfDefinitionFunc, _             = SingleDomainOfDefinition(defaultMinObjectValue, defaultMaxObjectValue)
	defaultProbabilityToTakeFromMemoryFunc       = StaticProbability(defaultProbabilityToTakeFromMemory)
	defaultProbabilityToApplyPitchAdjustmentFunc = StaticProbability(defaultProbabilityToApplyPitchAdjustment)
	defaultPitchAdjustingRateWidthFunc           = StaticWidth(defaultPitchAdjustingRateWidth)
)

var (
	ErrMinMoreMax = errors.New("min value more max value")
	ErrOutOfRange = errors.New("out off range")
)

// DomainOfDefinitionFunc функция возвращает допустимый диапазон значений для объекта.
type DomainOfDefinitionFunc func(objIndex int) (min, max float64, err error)

// SingleDomainOfDefinition область определения для всех одинакова.
func SingleDomainOfDefinition(min, max float64) (DomainOfDefinitionFunc, error) {
	if min >= max {
		return nil, ErrMinMoreMax
	}

	return func(_ int) (float64, float64, error) {
		return min, max, nil
	}, nil
}

// DifferentDomainOfDefinition область определения для всех своя.
func DifferentDomainOfDefinition(minmaxValues [][2]float64) (DomainOfDefinitionFunc, error) {
	for i := 0; i < len(minmaxValues); i++ {
		if minmaxValues[i][0] >= minmaxValues[i][1] {
			return nil, ErrMinMoreMax
		}
	}

	return func(objIndex int) (float64, float64, error) {
		if objIndex >= len(minmaxValues) {
			return 0.0, 0.0, ErrOutOfRange
		}

		return minmaxValues[objIndex][0], minmaxValues[objIndex][1], nil
	}, nil
}

// ProbabilityFunc функция должна возвращать вероятность [0,1] выполнения какого то действия по объекту.
type ProbabilityFunc func(objIndex int) float64

// StaticProbability вероятность всегда одна и таже.
func StaticProbability(v float64) ProbabilityFunc {
	return func(objIndex int) float64 {
		return v
	}
}

// RandProbability вероятность вегда рандомная.
func RandProbability() ProbabilityFunc {
	return func(objIndex int) float64 {
		return rand.Float64()
	}
}

// PitchAdjustingRateWidthFunc функция должна возвращать ширину шага [0,1] для объекта.
type PitchAdjustingRateWidthFunc func(objIndex int) float64

func StaticWidth(v float64) PitchAdjustingRateWidthFunc {
	return func(_ int) float64 {
		return v
	}
}

func RandWidth() PitchAdjustingRateWidthFunc {
	return func(_ int) float64 {
		return rand.Float64()
	}
}

type Opt func(opts *opts)

type opts struct {
	ctx context.Context

	// функция для получения области определения переменной.
	domainOfDefinition DomainOfDefinitionFunc
	// функция для получения вероятности выбора гармоники из memory.
	probabilityToTakeFromMemory ProbabilityFunc
	// функция для получения вероятности выполнения шага гармоники.
	probabilityToApplyPitchAdjustment ProbabilityFunc
	// функция для получения ширины с которой может меняться гармоника.
	pitchAdjustingRateWidth PitchAdjustingRateWidthFunc

	memorySize             int  // размер памяти.
	numberOfObjects        int  // количество объектов.
	numberOfImprovisations int  // количество импровизаций.
	isFindingMin           bool // флаг поиска минимума, а не максимума по умолчанию.
}

func Context(ctx context.Context) Opt {
	return func(opts *opts) {
		opts.ctx = ctx
	}
}

// DomainOfDefinition функция для получения области определения переменной.
func DomainOfDefinition(r DomainOfDefinitionFunc) Opt {
	return func(opts *opts) {
		opts.domainOfDefinition = r
	}
}

// ProbabilityToTakeFromMemory функция для получения вероятности выбора гармоники из memory.
func ProbabilityToTakeFromMemory(cr ProbabilityFunc) Opt {
	return func(opts *opts) {
		opts.probabilityToTakeFromMemory = cr
	}
}

// ProbabilityToApplyPitchAdjustment функция для получения вероятности выполнения шага гармоники.
func ProbabilityToApplyPitchAdjustment(par ProbabilityFunc) Opt {
	return func(opts *opts) {
		opts.probabilityToApplyPitchAdjustment = par
	}
}

// PitchAdjustingRateWidth функция для получения ширины с которой может меняться гармоника.
func PitchAdjustingRateWidth(parw PitchAdjustingRateWidthFunc) Opt {
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

// FindMin флаг поиска минимума, а не максимума по умолчанию.
func FindMin() Opt {
	return func(opts *opts) {
		opts.isFindingMin = true
	}
}

func defaultOpts(numberOfObjects int) *opts {
	return &opts{
		ctx:                               context.Background(),
		domainOfDefinition:                defaultDomainOfDefinitionFunc,
		probabilityToTakeFromMemory:       defaultProbabilityToTakeFromMemoryFunc,
		probabilityToApplyPitchAdjustment: defaultProbabilityToApplyPitchAdjustmentFunc,
		pitchAdjustingRateWidth:           defaultPitchAdjustingRateWidthFunc,
		memorySize:                        defaultMemorySize,
		numberOfObjects:                   numberOfObjects,
		numberOfImprovisations:            defaultNumberOfImprovisations,
	}
}
