package runtime

import (
	"errors"
)

var (
	ErrInvalidConvNilToInt    = errors.New("cannot convert Nil to Int")
	ErrInvalidConvNilToFloat  = errors.New("cannot convert Nil to Float")
	ErrInvalidConvNilToString = errors.New("cannot convert Nil to String")

	ErrInvalidOpAddOnNil = errors.New("cannot apply Add on a Nil value")
	ErrInvalidOpSubOnNil = errors.New("cannot apply Sub on a Nil value")
	ErrInvalidOpMulOnNil = errors.New("cannot apply Mul on a Nil value")
	ErrInvalidOpDivOnNil = errors.New("cannot apply Div on a Nil value")
	ErrInvalidOpPowOnNil = errors.New("cannot apply Pow on a Nil value")
	ErrInvalidOpModOnNil = errors.New("cannot apply Mod on a Nil value")
	ErrInvalidOpNotOnNil = errors.New("cannot apply Not on a Nil value")
	ErrInvalidOpUnmOnNil = errors.New("cannot apply Unm on a Nil value")
)

// Nil is the representation of the null type. It is semantically equivalent
// to Go's nil value, but it is represented as a struct.
type Nil struct{}

// Int is an invalid conversion.
func (ø Nil) Int() int {
	panic(ErrInvalidConvNilToInt)
}

// Float is an invalid conversion.
func (ø Nil) Float() float64 {
	panic(ErrInvalidConvNilToFloat)
}

// String returns the string "nil".
func (ø Nil) String() string {
	return "nil"
}

// Bool returns false.
func (ø Nil) Bool() bool {
	return false
}

// Add is an invalid operation.
func (ø Nil) Add(v Val) Val {
	panic(ErrInvalidOpAddOnNil)
}

// Sub is an invalid operation.
func (ø Nil) Sub(v Val) Val {
	panic(ErrInvalidOpSubOnNil)
}

// Mul is an invalid operation.
func (ø Nil) Mul(v Val) Val {
	panic(ErrInvalidOpMulOnNil)
}

// Div is an invalid operation.
func (ø Nil) Div(v Val) Val {
	panic(ErrInvalidOpDivOnNil)
}

// Mod is an invalid operation.
func (ø Nil) Mod(v Val) Val {
	panic(ErrInvalidOpModOnNil)
}

// Pow is an invalid operation.
func (ø Nil) Pow(v Val) Val {
	panic(ErrInvalidOpPowOnNil)
}

// Not is an invalid operation.
func (ø Nil) Not() Val {
	panic(ErrInvalidOpNotOnNil)
}

// Unm is an invalid operation.
func (ø Nil) Unm() Val {
	panic(ErrInvalidOpUnmOnNil)
}
