package random

import (
	"math/rand"

	"github.com/EmptyShadow/eltech.optimize/internal/functions"
	"gonum.org/v1/gonum/mat"
)

// RandFloatInRange генерация значения в пределах [b, t].
func RandFloatInRange(b, t float64) float64 {
	return b + rand.Float64()*(t-b)
}

// RandValueVar генерация значения из области определения.
func RandValueVar(d functions.VarDomain) float64 {
	return RandFloatInRange(d.Bottom, d.Top)
}

// RandValueVarInDomain генерация значения из области определения для определенной переменной.
func RandValueVarInDomain(varIndex int, fd functions.FuncDomain) float64 {
	return RandValueVar(fd.VarDomain(varIndex))
}

// RandMatrixInDomain генерация матрицы в пределах определеня функции.
func RandMatrixInDomain(r, c int, fd functions.FuncDomain) *mat.Dense {
	memory := make([]float64, r*c)

	for i := 0; i < r*c; i++ {
		varIndex := i % c
		memory[i] = RandValueVarInDomain(varIndex, fd)
	}

	return mat.NewDense(r, c, memory)
}
