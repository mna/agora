package runtime

// Converter declares the required methods to convert a value
// to any one of the supported types (except Object and Func).
type Converter interface {
	Int() int
	Float() float64
	String() string
	Bool() bool
	Native() interface{}
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

// Comparer defines the method required to compare two Values.
// Cmp() returns 1 if the method receiver value is greater, 0 if
// it is equal, and -1 if it is lower.
type Comparer interface {
	Cmp(Val) int
}

// Val is the representation of a value, any value, in the language.
// The supported value types are the following:
// * Integer (Int)
// * Float
// * String
// * Boolean (Bool)
// * Nil
// * Object
// * Func
type Val interface {
	Converter
	Comparer
	Arithmetic
}
