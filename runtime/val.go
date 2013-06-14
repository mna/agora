package runtime

// Converter declares the required methods to convert a value
// to any one of the supported types (except Object).
type Converter interface {
	Int() int
	Float() float64
	String() string
	Bool() bool
}

// Arithmetic defines the methods required to compute all
// the supported arithmetic operations.
type Arithmetic interface {
	Add(Val) Val
	Sub(Val) Val
	Mul(Val) Val
	Div(Val) Val
	Mod(Val) Val
	Pow(Val) Val
	Not() Val
	Unm() Val
}

// Val is the representation of a value, any value, in the language.
// The supported value types are the following:
// * Integer (Int)
// * Float
// * String
// * Boolean (Bool)
// * Nil
// * Object
type Val interface {
	Converter
	Arithmetic
}
