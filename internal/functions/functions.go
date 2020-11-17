package functions

import (
	"errors"
	"fmt"
	"math"

	"github.com/Knetic/govaluate"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
)

var (
	ErrIncorrectSizeOfArgsFunc    = errors.New("incorrect size of args function")
	ErrIncorrectSizeOfTheGradient = errors.New("incorrect size of the gradient")
	ErrIncorrectSizeOfTheHessian  = errors.New("incorrect size of the hessian")
)

func MustProblem(exp string, grad, hes *fd.Settings) optimize.Problem {
	prob, err := NewProblem(exp, grad, hes)
	if err != nil {
		panic(err)
	}

	return prob
}

func NewProblem(exp string, grad, hes *fd.Settings) (optimize.Problem, error) {
	f, err := NewExpression(exp)
	if err != nil {
		return optimize.Problem{}, err
	}

	return optimize.Problem{
		Func: f.Func,
		Grad: Gradient(f, grad),
		Hess: Hessian(f, hes),
	}, nil
}

type Function interface {
	Func(x []float64) float64
	Dimension() int
}

func CheckDimension(x []float64, f Function) {
	if len(x) != f.Dimension() {
		panic(fmt.Sprintf("dimension of the problem must be %d", f.Dimension()))
	}
}

type GradientFunc func(grad, x []float64)

type HessianFunc func(dst *mat.SymDense, x []float64)

func Gradient(f Function, settings *fd.Settings) GradientFunc {
	return func(grad, x []float64) {
		CheckDimension(x, f)

		if len(x) != len(grad) {
			panic(ErrIncorrectSizeOfTheGradient)
		}

		fd.Gradient(grad, f.Func, x, settings)
	}
}

func Hessian(f Function, settings *fd.Settings) HessianFunc {
	return func(dst *mat.SymDense, x []float64) {
		CheckDimension(x, f)

		if len(x) != dst.Symmetric() {
			panic(ErrIncorrectSizeOfTheHessian)
		}

		fd.Hessian(dst, f.Func, x, settings)
	}
}

const (
	Levi13 = `sin(3 * pi() * x) ** 2 + 
(x - 1) ** 2 * (1 + sin(3 * pi() * y) ** 2) + 
(y - 1) ** 2 * (1 + sin(2 * pi() * y) ** 2)`
	Himmelblau = "(x ** 2 + y - 11) ** 2 + (x + y ** 2 - 7) ** 2"
	Spheres    = "x ** 2 + y ** 2 + z ** 2 + f ** 2"
	Matias     = "0.26 *(x ** 2 + y ** 2) - 0.48 * x * y"
)

var mathFunctions = map[string]govaluate.ExpressionFunction{
	"sin": func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, ErrIncorrectSizeOfArgsFunc
		}

		return math.Sin(args[0].(float64)), nil
	},
	"cos": func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, ErrIncorrectSizeOfArgsFunc
		}

		return math.Cos(args[0].(float64)), nil
	},
	"pi": func(_ ...interface{}) (interface{}, error) {
		return math.Pi, nil
	},
}

type Expression struct {
	exp *govaluate.EvaluableExpression
}

func NewExpression(expression string) (*Expression, error) {
	exp, err := govaluate.NewEvaluableExpressionWithFunctions(expression, mathFunctions)
	if err != nil {
		return nil, err
	}

	return &Expression{exp: exp}, nil
}

func MustExpression(expression string) *Expression {
	e, err := NewExpression(expression)
	if err != nil {
		panic(err)
	}

	return e
}

func (e *Expression) Func(x []float64) float64 {
	CheckDimension(x, e)

	names := e.vars()
	vars := make(map[string]interface{})

	for i, name := range names {
		vars[name] = x[i]
	}

	y, err := e.exp.Evaluate(vars)
	if err != nil {
		panic(err)
	}

	return y.(float64)
}

func (e *Expression) Dimension() int {
	return len(e.vars())
}

func (e *Expression) vars() []string {
	var vars []string

	for _, val := range e.exp.Tokens() {
		if val.Kind != govaluate.VARIABLE {
			continue
		}

		name := val.Value.(string) //nolint

		dublicate := false

		for _, exists := range vars {
			if exists == name {
				dublicate = true

				break
			}
		}

		if !dublicate {
			vars = append(vars, name)
		}
	}

	return vars
}
