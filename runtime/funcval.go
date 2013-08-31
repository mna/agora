package runtime

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidConvFuncToInt    = errors.New("cannot convert Func to Int")
	ErrInvalidConvFuncToFloat  = errors.New("cannot convert Func to Float")
	ErrInvalidConvFuncToString = errors.New("cannot convert Func to String")

	ErrInvalidOpAddOnFunc = errors.New("cannot apply Add on a Func value")
	ErrInvalidOpSubOnFunc = errors.New("cannot apply Sub on a Func value")
	ErrInvalidOpMulOnFunc = errors.New("cannot apply Mul on a Func value")
	ErrInvalidOpDivOnFunc = errors.New("cannot apply Div on a Func value")
	ErrInvalidOpModOnFunc = errors.New("cannot apply Mod on a Func value")
	ErrInvalidOpUnmOnFunc = errors.New("cannot apply Unm on a Func value")
)

// funcVal implements most of the Val interface's methods, except
// Native() and Cmp() which must be done on the actual type.
type funcVal struct {
	ctx  *Ctx
	name string // this field is set only for Native funcs, because it uses the funcVal.dump()
	// call. Agora functions override dump() so that the agora name is used
	// automatically.
}

func (ø *funcVal) dump() string {
	return fmt.Sprintf("%s (Func)", ø.name)
}

// Int is an invalid conversion.
func (ø *funcVal) Int() int {
	panic(ErrInvalidConvFuncToInt)
}

// Float is an invalid conversion.
func (ø *funcVal) Float() float64 {
	panic(ErrInvalidConvFuncToFloat)
}

// String is an invalid conversion.
func (ø *funcVal) String() string {
	panic(ErrInvalidConvFuncToString)
}

// Bool returns true.
func (ø *funcVal) Bool() bool {
	return true
}

// Add is an invalid operation.
func (ø *funcVal) Add(v Val) Val {
	panic(ErrInvalidOpAddOnFunc)
}

// Sub is an invalid operation.
func (ø *funcVal) Sub(v Val) Val {
	panic(ErrInvalidOpSubOnFunc)
}

// Mul is an invalid operation.
func (ø *funcVal) Mul(v Val) Val {
	panic(ErrInvalidOpMulOnFunc)
}

// Div is an invalid operation.
func (ø *funcVal) Div(v Val) Val {
	panic(ErrInvalidOpDivOnFunc)
}

// Mod is an invalid operation.
func (ø *funcVal) Mod(v Val) Val {
	panic(ErrInvalidOpModOnFunc)
}

// Unm is an invalid operation.
func (ø *funcVal) Unm() Val {
	panic(ErrInvalidOpUnmOnFunc)
}
