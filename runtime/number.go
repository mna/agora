package runtime

import (
	"fmt"
	"math"
	"strconv"
)

// Number is the representation of the Number type. It is equivalent
// to Go's float64 type.
type Number float64

func (f Number) dump() string {
	return fmt.Sprintf("%s (Number)", strconv.FormatFloat(float64(f), 'f', -1, 64))
}

// Int returns the integer part of the float value.
func (f Number) Int() int64 {
	return int64(math.Trunc(float64(f)))
}

// Float returns the float value itself.
func (f Number) Float() float64 {
	return float64(f)
}

// String returns a string representation of the float value.
func (f Number) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

// Bool returns true if the float value is non-zero, false otherwise.
func (f Number) Bool() bool {
	return float64(f) != 0
}

// Native returns the Go native representation of the value.
func (f Number) Native() interface{} {
	return float64(f)
}

// Cmp compares the Number value to the provided value.
func (f Number) Cmp(v Val) int {
	if vf := v.Float(); float64(f) > vf {
		return 1
	} else if float64(f) < vf {
		return -1
	} else {
		return 0
	}
}

// Add performs the addition of the float value to another Val value, converted
// to a float.
func (f Number) Add(v Val) Val {
	return Number(float64(f) + v.Float())
}

// Sub performs the subtraction of another Val value, converted
// to a float, from the float value.
func (f Number) Sub(v Val) Val {
	return Number(float64(f) - v.Float())
}

// Mul performs the multiplication of the float value with another Val value,
// converted to a float.
func (f Number) Mul(v Val) Val {
	return Number(float64(f) * v.Float())
}

// Div performs the division of the float value by another Val value,
// converted to a float.
func (f Number) Div(v Val) Val {
	return Number(float64(f) / v.Float())
}

// Mod returns the modulo (remainder) of the division of the float value by
// another Val value, converted to a float.
func (f Number) Mod(v Val) Val {
	return Number(math.Mod(float64(f), v.Float()))
}

// Unm returns the unary minus operation applied to the float value.
// It switches the sign of the value.
func (f Number) Unm() Val {
	return Number(-float64(f))
}
