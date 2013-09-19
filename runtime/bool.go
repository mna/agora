package runtime

import (
	"fmt"
)

// Bool is the representation of the Boolean type. It is equivalent
// to Go's bool type.
type Bool bool

func (b Bool) dump() string {
	return fmt.Sprintf("%v (Bool)", bool(b))
}

// Int returns 1 if true, 0 if false.
func (b Bool) Int() int64 {
	if bool(b) {
		return 1
	}
	return 0
}

// Float returns 1 if true, 0 if false.
func (b Bool) Float() float64 {
	if bool(b) {
		return 1.0
	}
	return 0.0
}

// String returns "true" if true, "false" otherwise.
func (b Bool) String() string {
	if bool(b) {
		return "true"
	}
	return "false"
}

// Bool returns the boolean value itself.
func (b Bool) Bool() bool {
	return bool(b)
}

// Native returns the bool native Go representation.
func (b Bool) Native() interface{} {
	return bool(b)
}
