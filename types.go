package ai

import (
	"errors"
	"math"
	"math/rand"
)

var (
	ErrMinMoreMax = errors.New("min value more max value")
	ErrOutOfRange = errors.New("out off range")
)

// DomainOfDefinitionFunc функция возвращает допустимый диапазон значений для объекта.
type DomainOfDefinitionFunc func(objIndex int) (min, max float64, err error)

func MustSingleDomainOfDefinition(min, max float64) DomainOfDefinitionFunc {
	f, err := SingleDomainOfDefinition(min, max)
	if err != nil {
		panic(err)
	}

	return f
}

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

// StepWidth функция должна возвращать ширину шага [0,1] для объекта.
type StepWidth func(objIndex int) float64

func StaticWidth(v float64) StepWidth {
	return func(_ int) float64 {
		return v
	}
}

func RandWidth() StepWidth {
	return func(_ int) float64 {
		return rand.Float64()
	}
}

type Vector []float64

type VectorInitFunc func(i int) (float64, error)

func VectorWithDomainOfDefinitionInitFunc(definitionFunc DomainOfDefinitionFunc) VectorInitFunc {
	return func(i int) (float64, error) {
		min, max, err := definitionFunc(i)
		if err != nil {
			return 0.0, err
		}

		return RandInRangeFloat64(min, max), nil
	}
}

func NewVector(n int) Vector {
	v := make([]float64, n)

	return v
}

func NewVectorWithInitFunc(n int, initFunc VectorInitFunc) (Vector, error) {
	v := NewVector(n)

	for i := 0; i < n; i++ {
		e, err := initFunc(i)
		if err != nil {
			return nil, err
		}

		v[i] = e
	}

	return v, nil
}

func (v Vector) Diff(v2 Vector) Vector {
	v3 := NewVector(len(v2))

	for i := 0; i < len(v); i++ {
		v3[i] = v[i] - v2[i]
	}

	return v3
}

func (v Vector) Norma() float64 {
	sqs := v.SumSquares()

	return math.Sqrt(sqs)
}

func (v Vector) Sum() float64 {
	s := 0.0

	for i := 0; i < len(v); i++ {
		s += v[i]
	}

	return s
}

func (v Vector) SumSquares() float64 {
	s := 0.0

	for i := 0; i < len(v); i++ {
		s += v[i] * v[i]
	}

	return s
}

type Matrix []Vector

func NewMatrix(n, m int, defaultValue ...float64) Matrix {
	defValue := 0.0

	if len(defaultValue) > 0 {
		defValue = defaultValue[0]
	}

	hm, _ := NewMatrixWithInitFunc(n, m, func(i, j int) (float64, error) {
		return defValue, nil
	})

	return hm
}

type MatrixInitFunc func(i, j int) (float64, error)

// NewMatrixWithInitFunc новая функция с инициализирующей функцией.
//
// n - количество столбцов.
// m - количество строк.
func NewMatrixWithInitFunc(n, m int, initFunc MatrixInitFunc) (Matrix, error) {
	hm := make(Matrix, m)

	for i := 0; i < m; i++ {
		hm[i] = NewVector(n)

		for j := 0; j < n; j++ {
			v, err := initFunc(i, j)
			if err != nil {
				return nil, err
			}

			hm[i][j] = v
		}
	}

	return hm, nil
}

func (m Matrix) Row(index int) Vector {
	c := NewVector(len(m[index]))

	copy(c, m[index])

	return c
}

func (m Matrix) Column(index int) Vector {
	c := NewVector(len(m))

	for i := 0; i < len(m); i++ {
		c[i] = m[i][index]
	}

	return c
}

func RandInRangeFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

type OptiFunc func(vector Vector) float64
