package harmonysearch

import optimeth "github.com/EmptyShadow/eltech.optimeth"

const (
	DefaultMemorySize                        = 50
	DefaultNumberOfImprovisations            = 100
	DefaultMinObjectValue                    = -100.0
	DefaultMaxObjectValue                    = 100.0
	DefaultProbabilityToTakeFromMemory       = 0.5
	DefaultProbabilityToApplyPitchAdjustment = 0.5
	DefaultMinStep                           = 0.2
	DefaultMaxStep                           = 1.0
	DefaultEps                               = 1e-6
)

var (
	DefaultDomainOfDefinitionFunc = optimeth.MustSingleDomainOfDefinition(DefaultMinObjectValue,
		DefaultMaxObjectValue)
	DefaultProbabilityToTakeFromMemoryFunc       = optimeth.StaticProbability(DefaultProbabilityToTakeFromMemory)
	DefaultProbabilityToApplyPitchAdjustmentFunc = optimeth.StaticProbability(DefaultProbabilityToApplyPitchAdjustment)
	DefaultDomainOfDefinitionStep                = optimeth.MustSingleDomainOfDefinition(DefaultMinStep, DefaultMaxStep)
)

type Opt func(opts *opts)

type opts struct {
	// функция для получения области определения переменной.
	domainOfDefinition optimeth.DomainOfDefinitionFunc
	// функция для получения вероятности выбора гармоники из memory.
	probabilityToTakeFromMemory optimeth.ProbabilityFunc
	// функция для получения вероятности выполнения шага гармоники.
	probabilityToApplyPitchAdjustment optimeth.ProbabilityFunc
	// функция для получения области определения для шага.
	domainOfDefinitionStep optimeth.DomainOfDefinitionFunc

	memorySize             int     // размер памяти.
	numberOfObjects        int     // количество объектов.
	numberOfImprovisations int     // количество импровизаций.
	eps                    float64 // приближение.
	isFindingMin           bool    // флаг поиска минимума, а не максимума по умолчанию.
}

// DomainOfDefinition функция для получения области определения переменной.
func DomainOfDefinition(r optimeth.DomainOfDefinitionFunc) Opt {
	return func(opts *opts) {
		opts.domainOfDefinition = r
	}
}

// ProbabilityToTakeFromMemory функция для получения вероятности выбора гармоники из memory.
func ProbabilityToTakeFromMemory(cr optimeth.ProbabilityFunc) Opt {
	return func(opts *opts) {
		opts.probabilityToTakeFromMemory = cr
	}
}

// ProbabilityToApplyPitchAdjustment функция для получения вероятности выполнения шага гармоники.
func ProbabilityToApplyPitchAdjustment(par optimeth.ProbabilityFunc) Opt {
	return func(opts *opts) {
		opts.probabilityToApplyPitchAdjustment = par
	}
}

// DomainOfDefinitionStep функция для получения области определения для шага.
func DomainOfDefinitionStep(dds optimeth.DomainOfDefinitionFunc) Opt {
	return func(opts *opts) {
		opts.domainOfDefinitionStep = dds
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
		domainOfDefinition:                DefaultDomainOfDefinitionFunc,
		probabilityToTakeFromMemory:       DefaultProbabilityToTakeFromMemoryFunc,
		probabilityToApplyPitchAdjustment: DefaultProbabilityToApplyPitchAdjustmentFunc,
		domainOfDefinitionStep:            DefaultDomainOfDefinitionStep,
		memorySize:                        DefaultMemorySize,
		numberOfObjects:                   numberOfObjects,
		numberOfImprovisations:            DefaultNumberOfImprovisations,
		eps:                               DefaultEps,
	}
}
