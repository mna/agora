package runtime

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidOpSubOnString = errors.New("cannot apply Sub on a String value")
	ErrInvalidOpDivOnString = errors.New("cannot apply Div on a String value")
	ErrInvalidOpPowOnString = errors.New("cannot apply Pow on a String value")
	ErrInvalidOpModOnString = errors.New("cannot apply Mod on a String value")
	ErrInvalidOpUnmOnString = errors.New("cannot apply Unm on a String value")
)

// String is the representation of the String type. It is equivalent
// to Go's string type.
type String string

// Int converts the string representation of an integer to an integer value.
// If the string doesn't hold a valid integer representation,
// it panics.
func (ø String) Int() int {
	i, err := strconv.ParseInt(string(ø), 10, 0)
	if err != nil {
		// TODO : Custom error, or return the native error?
		panic(err)
	}
	return int(i)
}

// Float converts the string representation of a float to a float value.
// If the string doesn't hold a valid float representation,
// it panics.
func (ø String) Float() float64 {
	f, err := strconv.ParseFloat(string(ø), 64)
	if err != nil {
		// TODO : Custom error, or return the native error?
		panic(err)
	}
	return f
}

// String returns itself.
func (ø String) String() string {
	return string(ø)
}

// Bool returns true if the string value is not empty, false otherwise.
func (ø String) Bool() bool {
	return len(string(ø)) > 0
}

// Add performs the concatenation of the string with the supplied value,
// converted to a string.
func (ø String) Add(v Val) Val {
	// TODO : Problem if Int + String + Int, the String + Int will concatenate,
	// then Int + String will return bogus result. The leftmost operand should
	// dictate the expression's type, and force conversion of other Vals before
	// the evaluations. (in the compiler?) Or, always apply this simple rule:
	// the type result of an expression is always the type of the leftmost of
	// the two operands, so Int + String + Int results in Int + String. Use
	// parenthesis if you don't want this, i.e. (Int + String) + Int.
	return String(string(ø) + v.String())
}

// Sub is an invalid operation.
func (ø String) Sub(v Val) Val {
	panic(ErrInvalidOpSubOnString)
}

// Mul repeats n number of times the string, n being the
// value converted to an integer.
func (ø String) Mul(v Val) Val {
	return String(strings.Repeat(string(ø), v.Int()))
}

// Div is an invalid operation.
func (ø String) Div(v Val) Val {
	panic(ErrInvalidOpDivOnString)
}

// Mod is an invalid operation.
func (ø String) Mod(v Val) Val {
	panic(ErrInvalidOpModOnString)
}

// Pow is an invalid operation.
func (ø String) Pow(v Val) Val {
	panic(ErrInvalidOpPowOnString)
}

// Not switches the boolean value of the string, and returns a Boolean.
func (ø String) Not() Val {
	return Bool(!ø.Bool())
}

// Unm is an invalid operation.
func (ø String) Unm() Val {
	panic(ErrInvalidOpUnmOnString)
}
