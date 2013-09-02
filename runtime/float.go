package runtime

import (
	"fmt"
	"math"
	"strconv"
)

// Float is the representation of the Float type. It is equivalent
// to Go's float64 type.
type Float float64

func (f Float) dump() string {
	return fmt.Sprintf("%f (Float)", float64(f))
}

// Int returns the integer part of the float value.
func (f Float) Int() int {
	return int(math.Trunc(float64(f)))
}

// Float returns the float value itself.
func (f Float) Float() float64 {
	return float64(f)
}

// String returns a string representation of the float value.
func (f Float) String() string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

// Bool returns true if the float value is non-zero, false otherwise.
func (f Float) Bool() bool {
	return float64(f) != 0
}

// Native returns the Go native representation of the value.
func (f Float) Native() interface{} {
	return float64(f)
}

// Cmp compares the Float value to the provided value.
func (f Float) Cmp(v Val) int {
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
func (f Float) Add(v Val) Val {
	return Float(float64(f) + v.Float())
}

// Sub performs the subtraction of another Val value, converted
// to a float, from the float value.
func (f Float) Sub(v Val) Val {
	return Float(float64(f) - v.Float())
}

// Mul performs the multiplication of the float value with another Val value,
// converted to a float.
func (f Float) Mul(v Val) Val {
	return Float(float64(f) * v.Float())
}

// Div performs the division of the float value by another Val value,
// converted to a float.
func (f Float) Div(v Val) Val {
	return Float(float64(f) / v.Float())
}

// Mod returns the modulo (remainder) of the division of the float value by
// another Val value, converted to a float.
func (f Float) Mod(v Val) Val {
	return Float(math.Mod(float64(f), v.Float()))
}

// Unm returns the unary minus operation applied to the float value.
// It switches the sign of the value.
func (f Float) Unm() Val {
	return Float(-float64(f))
}
