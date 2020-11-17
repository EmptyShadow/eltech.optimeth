package optimize

import (
	"math/rand"
	"sort"

	"github.com/EmptyShadow/eltech.optimize/internal/functions"
	"github.com/EmptyShadow/eltech.optimize/internal/random"
	"gonum.org/v1/gonum/optimize"
)

const (
	HSConcurrent                      = 1
	DefaultHSMemorySize               = 50
	DefaultProbToTakeFromMemory       = 0.5
	DefaultProbToApplyPitchAdjustment = 0.5
	DefaultMaxStep                    = 1
)

var DefaultHSFD = functions.NewSingleFuncDomain(functions.VarDomain{
	Bottom: -100, //nolint
	Top:    100,  //nolint
})

type HSConfig struct {
	FD                         functions.FuncDomain
	MemorySize                 int
	ProbToTakeFromMemory       float64 // вероятность взять значение из памяти, иначе возьмем рандомное значение.
	ProbToApplyPitchAdjustment float64 // вероятность сделать шаг, иначе возьмем из памяти.
	MaxStep                    float64 // область определения шага.
}

func DefaultHSConfig() *HSConfig {
	return &HSConfig{
		FD:                         DefaultHSFD,
		MemorySize:                 DefaultHSMemorySize,
		ProbToTakeFromMemory:       DefaultProbToTakeFromMemory,
		ProbToApplyPitchAdjustment: DefaultProbToApplyPitchAdjustment,
		MaxStep:                    DefaultMaxStep,
	}
}

var _ optimize.Method = (*HS)(nil)

type memoryComponent struct {
	F float64
	X []float64
}

type hsState struct {
	memory []*memoryComponent
	dim    int

	status optimize.Status
	err    error
}

func newHSState(dim int, conf *HSConfig) *hsState {
	s := &hsState{}

	s.memory = make([]*memoryComponent, conf.MemorySize)
	for i := 0; i < len(s.memory); i++ {
		x := make([]float64, dim)

		for j := 0; j < dim; j++ {
			x[j] = random.RandValueVarInDomain(j, conf.FD)
		}

		s.memory[i] = &memoryComponent{X: x}
	}

	s.dim = dim
	s.status = optimize.NotTerminated
	s.err = nil

	return s
}

// HS метод гармонического поиска (HarmonySearch).
type HS struct {
	conf  *HSConfig
	state *hsState
}

// NewHS создать экземпляр метода гармонического поиска.
func NewHS(conf *HSConfig) *HS {
	return &HS{
		conf: conf,
	}
}

func (h *HS) Init(dim, _ int) int {
	h.state = newHSState(dim, h.conf)

	return HSConcurrent
}

func (h *HS) Run(operation chan<- optimize.Task, result <-chan optimize.Task, tasks []optimize.Task) {
	defer close(operation)

	h.evaluateMemory(operation, result)
	h.sortMemory()

	x := tasks[0].X

	funcEvaluation(operation, x) // Стартовое вычисление функции в стартовой точке.

	for res := range result {
		switch res.Op {
		default:
			panic("unknown operation")
		case optimize.PostIteration:
			break
		case optimize.MajorIteration:
			funcEvaluation(operation, res.X)
		case optimize.FuncEvaluation: // вычисление следующей импровизации и если она лучше какой то в памяти,то супер.
			if h.updateMemory(res.X, res.F) {
				h.sortMemory()

				x = res.X

				majorIteration(operation, res.X, res.F)

				continue
			}

			improvised := h.improvisation(x)
			funcEvaluation(operation, improvised)
		}
	}
}

func (h *HS) evaluateMemory(operation chan<- optimize.Task, result <-chan optimize.Task) {
	// расчитываем функцию в памяти.
	for i := 0; i < h.conf.MemorySize; i++ {
		m := h.state.memory[i]

		funcEvaluation(operation, m.X)

		res := <-result

		m.F = res.F
	}
}

func (h *HS) sortMemory() {
	sort.Slice(h.state.memory, func(i, j int) bool {
		l := h.state.memory[i]
		r := h.state.memory[j]

		return l.F > r.F
	})
}

func (h *HS) improvisation(improvised []float64) []float64 {
	dim := h.state.dim

	for varIndex := 0; varIndex < dim; varIndex++ {
		prob1 := rand.Float64()

		if prob1 >= h.conf.ProbToTakeFromMemory {
			improvised[varIndex] = random.RandValueVarInDomain(varIndex, h.conf.FD)

			continue
		}

		prod2 := rand.Float64()

		if prod2 >= h.conf.ProbToApplyPitchAdjustment {
			m := h.state.memory[rand.Intn(h.conf.MemorySize)]

			improvised[varIndex] = m.X[varIndex]

			continue
		}

		d := h.conf.FD.VarDomain(varIndex)
		step := random.RandFloatInRange(-1.0, 1.0) * h.conf.MaxStep
		improvised[varIndex] = d.Normalize(improvised[varIndex] + step)
	}

	return improvised
}

func (h *HS) updateMemory(x []float64, f float64) bool {
	for i := 0; i < h.conf.MemorySize; i++ {
		if f < h.state.memory[i].F {
			h.state.memory[i].X = x
			h.state.memory[i].F = f

			return true
		}
	}

	return false
}

func (h *HS) Uses(_ optimize.Available) (uses optimize.Available, err error) {
	return optimize.Available{
		Grad: false,
		Hess: false,
	}, nil
}

func (h *HS) Status() (optimize.Status, error) {
	return h.state.status, h.state.err
}

type HSStep struct {
	D functions.FuncDomain
}
