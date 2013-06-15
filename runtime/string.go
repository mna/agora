package runtime

import (
	"bytes"
	"strconv"
	"strings"
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

// Sub removes Int characters from the string if the value is not a string,
// or it removes the character set of the value if it is a string.
func (ø String) Sub(v Val) Val {
	switch x := v.(type) {
	case String:
		return String(strings.Map(func(r rune) rune {
			if ix := strings.IndexRune(string(x), r); ix >= 0 {
				return rune(-1)
			}
			return r
		}, string(ø)))
	default:
		s := string(ø)
		return String(s[:len(s)-v.Int()])
	}
}

// Mul repeats n number of times the string, n being the
// value converted to an integer.
func (ø String) Mul(v Val) Val {
	return String(strings.Repeat(string(ø), v.Int()))
}

// Div splits the string in n number of substrings, n being the value
// converted to an integer. TODO : Div by a string, splits at each chararacter in the set?
func (ø String) Div(v Val) Val {
	return String("")
}

// Mod returns the modulo (remainder) of the division of the float value by
// another Val value, converted to a float.
func (ø String) Mod(v Val) Val {
	return String("")
}

// Pow returns the float raised at the power of the other Val value, converted
// to a float.
func (ø String) Pow(v Val) Val {
	return String("")
}

// Not switches the boolean value of the string, and returns a Boolean.
func (ø String) Not() Val {
	return Bool(!ø.Bool())
}

// Unm reverses the string.
func (ø String) Unm() Val {
	s := string(ø)
	l := len(s)
	if l == 0 {
		return ø
	}
	// TODO : Need to use a Reader.ReadRune
	b := make([]byte, l)
	buf := bytes.NewBuffer(b)
	for i := l; i > 0; i-- {
		buf.WriteByte(s[i-1])
	}
	return String(buf.String())
}
