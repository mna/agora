package runtime

const (
	// The string representation of the nil value
	NilString = "nil"
)

var (
	// The one and only Nil instance
	Nil = null{}
)

// Null is the representation of the null type. It is semantically equivalent
// to Go's nil value, but it is represented as an empty struct to implement
// the Val interface so that it is a valid agora value.
type null struct{}

func (n null) Dump() string {
	return "[Nil]"
}

// Int is an invalid conversion.
func (n null) Int() int64 {
	panic(NewTypeError(Type(n), "", "int"))
}

// Float is an invalid conversion.
func (n null) Float() float64 {
	panic(NewTypeError(Type(n), "", "float"))
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
