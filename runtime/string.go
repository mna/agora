package runtime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// Predefined errors
	ErrInvalidOpSubOnString = errors.New("cannot apply Sub on a String value")
	ErrInvalidOpDivOnString = errors.New("cannot apply Div on a String value")
	ErrInvalidOpModOnString = errors.New("cannot apply Mod on a String value")
	ErrInvalidOpUnmOnString = errors.New("cannot apply Unm on a String value")
)

// String is the representation of the String type. It is equivalent
// to Go's string type.
type String string

// Pretty-prints the string value.
func (s String) dump() string {
	return fmt.Sprintf("\"%s\" (String)", string(s))
}

// Int converts the string representation of an integer to an integer value.
// If the string doesn't hold a valid integer representation,
// it panics.
func (s String) Int() int64 {
	i, err := strconv.ParseInt(string(s), 10, 0)
	if err != nil {
		panic(err)
	}
	return int64(i)
}

// Float converts the string representation of a float to a float value.
// If the string doesn't hold a valid float representation,
// it panics.
func (s String) Float() float64 {
	f, err := strconv.ParseFloat(string(s), 64)
	if err != nil {
		panic(err)
	}
	return f
}

// String returns itself.
func (s String) String() string {
	return string(s)
}

// Bool returns true if the string value is not empty, false otherwise.
func (s String) Bool() bool {
	return len(string(s)) > 0
}

// Native returns the Go native representation of the value.
func (s String) Native() interface{} {
	return string(s)
}

// Cmp compares the value with another value.
func (s String) Cmp(v Val) int {
	if vs := v.String(); string(s) > vs {
		return 1
	} else if string(s) < vs {
		return -1
	} else {
		return 0
	}
}

// Add performs the concatenation of the string with the supplied value,
// converted to a string.
func (s String) Add(v Val) Val {
	return String(string(s) + v.String())
}

// Sub is an invalid operation.
func (s String) Sub(v Val) Val {
	panic(ErrInvalidOpSubOnString)
}

// Mul repeats n number of times the string, n being the
// value converted to an integer.
// TODO : Is this a *good idea*?
func (s String) Mul(v Val) Val {
	return String(strings.Repeat(string(s), int(v.Int())))
}

// Div is an invalid operation.
func (s String) Div(v Val) Val {
	panic(ErrInvalidOpDivOnString)
}

// Mod is an invalid operation.
func (s String) Mod(v Val) Val {
	panic(ErrInvalidOpModOnString)
}

// Unm is an invalid operation.
func (s String) Unm() Val {
	panic(ErrInvalidOpUnmOnString)
}
