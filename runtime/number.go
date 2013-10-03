package runtime

import (
	"fmt"
	"math"
	"strconv"
)

// Number is the representation of the Number type. It is equivalent
// to Go's float64 type.
type Number float64

func (f Number) Dump() string {
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
