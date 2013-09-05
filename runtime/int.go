package runtime

import (
	"fmt"
	"strconv"
)

// Int is the representation of the Integer type. It is equivalent
// to Go's int64 type
// TODO : Will probably disappear in future versions, all numbers
// will be of type Float.
type Int int64

func (i Int) dump() string {
	return fmt.Sprintf("%d (Int)", int(i))
}

// Int returns the integer value itself.
func (i Int) Int() int {
	return int(i)
}

// Float returns a floating point representation of the integer value.
func (i Int) Float() float64 {
	return float64(i)
}

// String returns a base 10 string representation of the integer value.
func (i Int) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// Bool returns true if the integer value is non-zero, false otherwise.
func (i Int) Bool() bool {
	return int(i) != 0
}

// Native returns the Go native representation of the value.
func (i Int) Native() interface{} {
	return int(i)
}

// Cmp compares the integer value with another value.
func (i Int) Cmp(v Val) int {
	if vi := v.Int(); int(i) > vi {
		return 1
	} else if int(i) < vi {
		return -1
	} else {
		return 0
	}
}

// Add performs the addition of the integer value to another Val value, converted
// to an int.
func (i Int) Add(v Val) Val {
	return Int(int(i) + v.Int())
}

// Sub performs the subtraction of another Val value, converted
// to an int, from the integer value.
func (i Int) Sub(v Val) Val {
	return Int(int(i) - v.Int())
}

// Mul performs the multiplication of the integer value with another Val value,
// converted to an int.
func (i Int) Mul(v Val) Val {
	return Int(int(i) * v.Int())
}

// Div performs the division of the integer value by another Val value,
// converted to an int.
func (i Int) Div(v Val) Val {
	return Int(int(i) / v.Int())
}

// Mod returns the modulo (remainder) of the division of the integer value by
// another Val value, converted to an int.
func (i Int) Mod(v Val) Val {
	return Int(int(i) % v.Int())
}

// Unm returns the unary minus operation applied to the integer value.
// It switches the sign of the value.
func (i Int) Unm() Val {
	return Int(-int(i))
}
