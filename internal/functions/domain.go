package functions

import (
	"fmt"
)

type ErrDomain struct {
	V float64
	D VarDomain
}

func (e *ErrDomain) Error() string {
	return fmt.Sprintf("value %f is not included in the domain [%f, %f]", e.V, e.D.Bottom, e.D.Top)
}

// VarDomain область определения [VarDomain.Bottom, VarDomain.Top] переменной.
type VarDomain struct {
	Bottom float64 // наименьшее значение.
	Top    float64 // наибольшее значение.
}

func NewVarDomain(b, t float64) *VarDomain {
	return &VarDomain{Bottom: b, Top: t}
}

// Validate валидация попадания в область определения.
func (d *VarDomain) Validate(v float64) error {
	if d.Bottom <= v && v <= d.Top {
		return nil
	}

	return &ErrDomain{
		V: v,
		D: *d,
	}
}

// Normalize нормализуется значение.
//
// Если меньше Bottom, то вернется Bottom.
// Если больше Top, то вернется Top.
// Иначе все хорошо и вернется v.
func (d *VarDomain) Normalize(v float64) float64 {
	if v < d.Bottom {
		return d.Bottom
	}

	if v > d.Top {
		return d.Bottom
	}

	return v
}

type FuncDomain interface {
	VarDomain(varIndex int) VarDomain
}

var _ FuncDomain = (*SingleFuncDomain)(nil)

// SingleFuncDomain область определения одинакова для всех переменных.
type SingleFuncDomain struct {
	d VarDomain
}

func NewSingleFuncDomain(d VarDomain) *SingleFuncDomain {
	return &SingleFuncDomain{d: d}
}

func (s *SingleFuncDomain) VarDomain(_ int) VarDomain {
	return s.d
}

var _ FuncDomain = (*MultipleFuncDomain)(nil)

// MultipleFuncDomain область определения специфичная для каждой переменной.
type MultipleFuncDomain struct {
	ds []VarDomain
}

func NewMultipleFuncDomain(ds ...VarDomain) *MultipleFuncDomain {
	return &MultipleFuncDomain{ds: ds}
}

func (m *MultipleFuncDomain) VarDomain(varIndex int) VarDomain {
	return m.ds[varIndex]
}
