package testdata

import (
	ai "github.com/EmptyShadow/eltech.ai"
	"math"
)

var (
	MatiasFunc = TestFunc{
		Name: "Функция Матьяса",
		F: func(vars ai.Vector) float64 {
			x := vars[0]
			y := vars[1]

			return 0.26*(x*x+y*y) - 0.48*x*y
		},
		NumberOfVariables:  2,
		Solutions:          []ai.Vector{{0.0, 0.0}},
		Min:                0.0,
		DomainOfDefinition: ai.MustSingleDomainOfDefinition(-10.0, 10.0),
	}

	SpheresFunc = TestFunc{
		Name: "Функция сферы",
		F: func(vars ai.Vector) float64 {
			return vars.SumSquares()
		},
		NumberOfVariables:  2,
		Solutions:          []ai.Vector{ai.NewVector(2)},
		Min:                0.0,
		DomainOfDefinition: ai.MustSingleDomainOfDefinition(-100.0, 100.0),
	}

	Levi13Func = TestFunc{
		Name: "Функция Леви N 13",
		F: func(vars ai.Vector) float64 {
			x := vars[0]
			y := vars[1]

			sinx := math.Pow(math.Sin(3.0*math.Pi*x), 2.0)
			sin3y := math.Pow(math.Sin(3.0*math.Pi*y), 2.0)
			sin2y := math.Pow(math.Sin(2.0*math.Pi*y), 2.0)

			return sinx + math.Pow(x-1.0, 2.0)*(1.0+sin3y) + math.Pow(y-1.0, 2.0)*(1.0+sin2y)
		},
		NumberOfVariables:  2,
		Solutions:          []ai.Vector{{1.0, 1.0}},
		Min:                0.0,
		DomainOfDefinition: ai.MustSingleDomainOfDefinition(-10.0, 10.0),
	}

	HimmelblauFunc = TestFunc{
		Name: "Функция Химмельблау",
		F: func(vars ai.Vector) float64 {
			x := vars[0]
			y := vars[1]

			return math.Pow(x*x+y-11, 2.0) + math.Pow(x+y*y-7, 2.0)
		},
		NumberOfVariables: 2,
		Solutions: []ai.Vector{
			{3.0, 2.0}, {-2.805118, 3.131312},
			{-3.779310, -3.283186}, {3.584428, -1.848126}},
		Min:                0,
		DomainOfDefinition: ai.MustSingleDomainOfDefinition(-5.0, 5.0),
	}
)

type Solution struct {
	Point ai.Vector
	Min   float64
}

type TestFunc struct {
	Name               string
	F                  ai.OptiFunc
	NumberOfVariables  int
	Solutions          []ai.Vector
	Min                float64
	DomainOfDefinition ai.DomainOfDefinitionFunc
}
