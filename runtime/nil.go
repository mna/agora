package runtime

import (
	"errors"
)

const (
	NilString = "nil"
)

var (
	// The one and only Nil instance
	Nil = null{}

	ErrInvalidConvNilToInt   = errors.New("cannot convert Nil to Int")
	ErrInvalidConvNilToFloat = errors.New("cannot convert Nil to Float")

	ErrInvalidOpAddOnNil = errors.New("cannot apply Add on a Nil value")
	ErrInvalidOpSubOnNil = errors.New("cannot apply Sub on a Nil value")
	ErrInvalidOpMulOnNil = errors.New("cannot apply Mul on a Nil value")
	ErrInvalidOpDivOnNil = errors.New("cannot apply Div on a Nil value")
	ErrInvalidOpPowOnNil = errors.New("cannot apply Pow on a Nil value")
	ErrInvalidOpModOnNil = errors.New("cannot apply Mod on a Nil value")
	ErrInvalidOpUnmOnNil = errors.New("cannot apply Unm on a Nil value")
)

// Null is the representation of the null type. It is semantically equivalent
// to Go's nil value, but it is represented as a struct.
type null struct{}

func (ø null) dump() string {
	return "[Nil]"
}

// Int is an invalid conversion.
func (ø null) Int() int {
	panic(ErrInvalidConvNilToInt)
}

// Float is an invalid conversion.
func (ø null) Float() float64 {
	panic(ErrInvalidConvNilToFloat)
}

// String returns the string "nil".
func (ø null) String() string {
	return NilString
}

// Bool returns false.
func (ø null) Bool() bool {
	return false
}

func (ø null) Native() interface{} {
	return nil
}

func (ø null) Cmp(v Val) int {
	if v == Nil {
		return 0
	}
	return -1
}

// Add is an invalid operation.
func (ø null) Add(v Val) Val {
	panic(ErrInvalidOpAddOnNil)
}

// Sub is an invalid operation.
func (ø null) Sub(v Val) Val {
	panic(ErrInvalidOpSubOnNil)
}

// Mul is an invalid operation.
func (ø null) Mul(v Val) Val {
	panic(ErrInvalidOpMulOnNil)
}

// Div is an invalid operation.
func (ø null) Div(v Val) Val {
	panic(ErrInvalidOpDivOnNil)
}

// Mod is an invalid operation.
func (ø null) Mod(v Val) Val {
	panic(ErrInvalidOpModOnNil)
}

// Pow is an invalid operation.
func (ø null) Pow(v Val) Val {
	panic(ErrInvalidOpPowOnNil)
}

// Not switches the boolean value of nil, and returns a Boolean.
func (ø null) Not() Val {
	return Bool(!ø.Bool())
}

// Unm is an invalid operation.
func (ø null) Unm() Val {
	panic(ErrInvalidOpUnmOnNil)
}
