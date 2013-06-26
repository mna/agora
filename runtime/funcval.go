package runtime

import (
	"errors"
)

var (
	ErrInvalidConvFuncToInt    = errors.New("cannot convert Func to Int")
	ErrInvalidConvFuncToFloat  = errors.New("cannot convert Func to Float")
	ErrInvalidConvFuncToString = errors.New("cannot convert Func to String")

	ErrInvalidOpAddOnFunc = errors.New("cannot apply Add on a Func value")
	ErrInvalidOpSubOnFunc = errors.New("cannot apply Sub on a Func value")
	ErrInvalidOpMulOnFunc = errors.New("cannot apply Mul on a Func value")
	ErrInvalidOpDivOnFunc = errors.New("cannot apply Div on a Func value")
	ErrInvalidOpPowOnFunc = errors.New("cannot apply Pow on a Func value")
	ErrInvalidOpModOnFunc = errors.New("cannot apply Mod on a Func value")
	ErrInvalidOpUnmOnFunc = errors.New("cannot apply Unm on a Func value")
)

type funcVal struct{}

// Int is an invalid conversion.
func (ø funcVal) Int() int {
	panic(ErrInvalidConvFuncToInt)
}

// Float is an invalid conversion.
func (ø funcVal) Float() float64 {
	panic(ErrInvalidConvFuncToFloat)
}

// String is an invalid conversion.
func (ø funcVal) String() string {
	panic(ErrInvalidConvFuncToString)
}

// Bool returns true.
func (ø funcVal) Bool() bool {
	return true
}

func (ø funcVal) Native() interface{} {
	return ø
}

// TODO : Maybe override this one in the real func values
func (ø funcVal) Cmp(v Val) int {
	if ø == v {
		// Point to same function
		return 0
	}
	// Otherwise, always return -1 (no rational way to compare 2 functions)
	return -1
}

// Add is an invalid operation.
func (ø funcVal) Add(v Val) Val {
	panic(ErrInvalidOpAddOnFunc)
}

// Sub is an invalid operation.
func (ø funcVal) Sub(v Val) Val {
	panic(ErrInvalidOpSubOnFunc)
}

// Mul is an invalid operation.
func (ø funcVal) Mul(v Val) Val {
	panic(ErrInvalidOpMulOnFunc)
}

// Div is an invalid operation.
func (ø funcVal) Div(v Val) Val {
	panic(ErrInvalidOpDivOnFunc)
}

// Mod is an invalid operation.
func (ø funcVal) Mod(v Val) Val {
	panic(ErrInvalidOpModOnFunc)
}

// Pow is an invalid operation.
func (ø funcVal) Pow(v Val) Val {
	panic(ErrInvalidOpPowOnFunc)
}

// Unm is an invalid operation.
func (ø funcVal) Unm() Val {
	panic(ErrInvalidOpUnmOnFunc)
}
