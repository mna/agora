package runtime

import (
	"errors"
)

const (
	// The string representation of the nil value
	NilString = "nil"
)

var (
	// The one and only Nil instance
	Nil = null{}

	// Predefined errors
	ErrInvalidConvNilToInt   = errors.New("cannot convert Nil to Int")
	ErrInvalidConvNilToFloat = errors.New("cannot convert Nil to Float")

	ErrInvalidOpAddOnNil = errors.New("cannot apply Add on a Nil value")
	ErrInvalidOpSubOnNil = errors.New("cannot apply Sub on a Nil value")
	ErrInvalidOpMulOnNil = errors.New("cannot apply Mul on a Nil value")
	ErrInvalidOpDivOnNil = errors.New("cannot apply Div on a Nil value")
	ErrInvalidOpModOnNil = errors.New("cannot apply Mod on a Nil value")
	ErrInvalidOpUnmOnNil = errors.New("cannot apply Unm on a Nil value")
)

// Null is the representation of the null type. It is semantically equivalent
// to Go's nil value, but it is represented as an empty struct to implement
// the Val interface so that it is a valid agora value.
type null struct{}

func (n null) dump() string {
	return "[Nil]"
}

// Int is an invalid conversion.
func (n null) Int() int64 {
	panic(ErrInvalidConvNilToInt)
}

// Float is an invalid conversion.
func (n null) Float() float64 {
	panic(ErrInvalidConvNilToFloat)
}

// String returns the string "nil".
func (n null) String() string {
	return NilString
}

// Bool returns false.
func (n null) Bool() bool {
	return false
}

// Native returns the Go native representation of the value.
func (n null) Native() interface{} {
	return nil
}

// Cmp compares the value with another value.
func (n null) Cmp(v Val) int {
	if v == Nil {
		return 0
	}
	return -1
}

// Add is an invalid operation.
func (n null) Add(v Val) Val {
	panic(ErrInvalidOpAddOnNil)
}

// Sub is an invalid operation.
func (n null) Sub(v Val) Val {
	panic(ErrInvalidOpSubOnNil)
}

// Mul is an invalid operation.
func (n null) Mul(v Val) Val {
	panic(ErrInvalidOpMulOnNil)
}

// Div is an invalid operation.
func (n null) Div(v Val) Val {
	panic(ErrInvalidOpDivOnNil)
}

// Mod is an invalid operation.
func (n null) Mod(v Val) Val {
	panic(ErrInvalidOpModOnNil)
}

// Unm is an invalid operation.
func (n null) Unm() Val {
	panic(ErrInvalidOpUnmOnNil)
}
