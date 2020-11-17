package functions

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestExpression_Func(t *testing.T) {
	asserting := assert.New(t)

	x := 3.0
	y := 2.0
	z := 1.0
	want := 1230.0

	f, err := NewExpression("x * 10 + y * 100 + z * 1000")
	asserting.NoError(err)

	got := f.Func([]float64{x, y, z})
	asserting.Equal(got, want)
}

func TestExpression_Func2(t *testing.T) {
	asserting := assert.New(t)

	x := 1.0
	y := 2.0
	want := 33.0

	f, err := NewExpression("x + y ** 3 * 2 * y")
	asserting.NoError(err)

	got := f.Func([]float64{x, y})
	asserting.Equal(got, want)
}

func TestExpression_Func3(t *testing.T) {
	asserting := assert.New(t)

	x := 1.0
	y := 2.0
	want := math.Pow(math.Sin(x*y*math.Pi), 2.0) + math.Pow(math.Sin(x+y), 5.0) + x*y

	f, err := NewExpression("sin(x * y * pi()) ** 2 + sin(x + y) ** 5 + x * y")
	asserting.NoError(err)

	got := f.Func([]float64{x, y})
	asserting.Equal(got, want)
}

func TestExpression_Levi13(t *testing.T) {
	asserting := assert.New(t)

	x := 1.0
	y := 1.0
	sinx := math.Pow(math.Sin(3.0*math.Pi*x), 2.0)
	sin3y := math.Pow(math.Sin(3.0*math.Pi*y), 2.0)
	sin2y := math.Pow(math.Sin(2.0*math.Pi*y), 2.0)

	want := sinx + math.Pow(x-1.0, 2.0)*(1.0+sin3y) + math.Pow(y-1.0, 2.0)*(1.0+sin2y)

	f, err := NewExpression(Levi13)
	asserting.NoError(err)

	got := f.Func([]float64{x, y})
	asserting.Equal(got, want)
}

func TestExpression_Himmelblau(t *testing.T) {
	tests := []struct {
		name string
		x, y float64
	}{
		{
			x: 3.0,
			y: 2.0,
		},
		{
			x: -2.805118,
			y: 3.131312,
		},
		{
			x: -3.779310,
			y: -3.283186,
		},
		{
			x: 3.584428,
			y: -1.848126,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asserting := assert.New(t)

			x := tt.x
			y := tt.y
			want := math.Pow(x*x+y-11, 2.0) + math.Pow(x+y*y-7, 2.0)

			f, err := NewExpression(Himmelblau)
			asserting.NoError(err)

			got := f.Func([]float64{x, y})
			asserting.Equal(got, want)
		})
	}
}

func TestExpression_Spheres(t *testing.T) {
	asserting := assert.New(t)

	x := 10.0
	y := 2.0
	z := 5.0
	fa := 9.0
	want := math.Pow(x, 2.0) + math.Pow(y, 2.0) + math.Pow(z, 2.0) + math.Pow(fa, 2.0)

	f, err := NewExpression(Spheres)
	asserting.NoError(err)

	got := f.Func([]float64{x, y, z, fa})
	asserting.Equal(got, want)
}

func TestExpression_Matias(t *testing.T) {
	asserting := assert.New(t)

	x := 0.0
	y := 0.0
	want := 0.26 * (math.Pow(x, 2.0) + math.Pow(y, 2.0)) - 0.48 * x * y
	f, err := NewExpression(Matias)
	asserting.NoError(err)

	got := f.Func([]float64{x, y})
	asserting.Equal(got, want)
}
