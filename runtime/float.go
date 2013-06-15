package runtime

import (
	"math"
	"strconv"
)

// Float is the representation of the Float type. It is equivalent
// to Go's float64 type.
type Float float64

// Int returns the integer part of the float value.
func (ø Float) Int() int {
	return int(math.Trunc(float64(ø)))
}

// Float returns the float value itself.
func (ø Float) Float() float64 {
	return float64(ø)
}

// String returns a string representation of the float value.
func (ø Float) String() string {
	return strconv.FormatFloat(float64(ø), 'f', -1, 64)
}

// Bool returns true if the float value is non-zero, false otherwise.
func (ø Float) Bool() bool {
	return float64(ø) != 0
}

// Add performs the addition of the float value to another Val value, converted
// to a float.
func (ø Float) Add(v Val) Val {
	return Float(float64(ø) + v.Float())
}

// Sub performs the subtraction of another Val value, converted
// to a float, from the float value.
func (ø Float) Sub(v Val) Val {
	return Float(float64(ø) - v.Float())
}

// Mul performs the multiplication of the float value with another Val value,
// converted to a float.
func (ø Float) Mul(v Val) Val {
	return Float(float64(ø) * v.Float())
}

// Div performs the division of the float value by another Val value,
// converted to a float.
func (ø Float) Div(v Val) Val {
	return Float(float64(ø) / v.Float())
}

// Mod returns the modulo (remainder) of the division of the float value by
// another Val value, converted to a float.
func (ø Float) Mod(v Val) Val {
	return Float(math.Mod(float64(ø), v.Float()))
}

// Pow returns the float raised at the power of the other Val value, converted
// to a float.
func (ø Float) Pow(v Val) Val {
	return Float(math.Pow(float64(ø), v.Float()))
}

// Not switches the boolean value of the float and returns a boolean.
func (ø Float) Not() Val {
	return Bool(!ø.Bool())
}

// Unm returns the unary minus operation applied to the float value.
// It switches the sign of the value.
func (ø Float) Unm() Val {
	return Float(-float64(ø))
}
