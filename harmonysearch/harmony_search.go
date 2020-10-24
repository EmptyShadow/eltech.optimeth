package harmonysearch

import (
	"context"
	"fmt"

	"github.com/EmptyShadow/eltech.ai"
)

const (
	DefaultNumberOfSolutionCandidates = 50
	DefaultPitchAdjustingRate         = 1.5
	DefaultMaxNumberOfSearches        = 1000
)

var DefaultRandFloatFunc = ai.RangeRandFloat64Func(-100.0, 100.0)

// matrix матрица гармоник.
type matrix struct {
	hm [][]float64
}

// newMatrix функция инициализирует матрицу гармоник.
func newMatrix(r ai.RandFloat64Func, ndp, nsc int) (*matrix, error) {
	hm := make([][]float64, ndp)

	for i := 0; i < ndp; i++ {
		hm[i] = make([]float64, nsc)

		for j := 0; j < nsc; j++ {
			rv, err := r()
			if err != nil {
				return nil, fmt.Errorf("failed generate rand value: %w", err)
			}

			hm[i][j] = rv
		}
	}

	return &matrix{hm: hm}, nil
}

type Opt func(opts *opts)

type opts struct {
	ctx                 context.Context
	ndp                 int     // количество размерных задач.
	nsc                 int     // количество кандидатов на решение, по умолчанию, 50.
	par                 float64 // скорость регулирования шага.
	maxNumberOfSearches int     // максимальное количество итераций поиска.
	r                   ai.RandFloat64Func
}

func Context(ctx context.Context) Opt {
	return func(opts *opts) {
		opts.ctx = ctx
	}
}

func NumberOfSolutionCandidates(n int) Opt {
	return func(opts *opts) {
		opts.nsc = n
	}
}

func PitchAdjustingRate(par float64) Opt {
	return func(opts *opts) {
		opts.par = par
	}
}

func Randomizer(r ai.RandFloat64Func) Opt {
	return func(opts *opts) {
		opts.r = r
	}
}

func defaultOpts(ndp int) *opts {
	return &opts{
		ctx:                 context.Background(),
		ndp:                 ndp,
		nsc:                 DefaultNumberOfSolutionCandidates,
		par:                 DefaultPitchAdjustingRate,
		maxNumberOfSearches: DefaultMaxNumberOfSearches,
		r:                   DefaultRandFloatFunc,
	}
}

type state struct {
	m        *matrix
	solution []float64
}

type Compositor struct {
	*opts
	state *state
}

func NewCompositor(ndp int, opts ...Opt) *Compositor {
	_opts := defaultOpts(ndp)

	for _, opt := range opts {
		opt(_opts)
	}

	s := &state{}

	return &Compositor{opts: _opts, state: s}
}

func (c *Compositor) Optimize() ([]float64, error) {
	if err := c.randSolutions(); err != nil {
		return nil, err
	}

	currentStep := 0

	for {
		select {
		case <-c.ctx.Done():
			// TODO: возвращать промежуточное решение
			return nil, c.ctx.Err()
		default:
		}

		c.selection()
		c.evaluation()
		c.comparison()

		if c.isBetterSolution() {
			c.replacement()
		} else {
			c.elimination()
		}

		if currentStep >= c.maxNumberOfSearches {
			return nil, nil
		}
	}
}

func (c *Compositor) randSolutions() (err error) {
	c.state.m, err = newMatrix(c.r, c.ndp, c.nsc)
	if err != nil {
		return err
	}

	c.state.solution = make([]float64, c.ndp)

	return nil
}

func (c *Compositor) selection() {

}

func (c *Compositor) evaluation() {

}

func (c *Compositor) comparison() {

}

func (c *Compositor) isBetterSolution() bool {

}

func (c *Compositor) elimination() {

}

func (c *Compositor) replacement() {

}
