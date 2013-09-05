package runtime

// Converter declares the required methods to convert a value
// to any one of the supported types (except Object and Func).
type Converter interface {
	Int() int64
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
	Unm() Val
}

// Comparer defines the method required to compare two Values.
// Cmp() returns 1 if the method receiver value is greater, 0 if
// it is equal, and -1 if it is lower.
type Comparer interface {
	Cmp(Val) int
}

// The dumper interface defines the required behaviour to pretty-print
// the values.
type dumper interface {
	dump() string
}

// Val is the representation of a value, any value, in the language.
// The supported value types are the following:
// * Integer (Int) - will likely disappear in future versions
// * Float
// * String
// * Boolean (Bool)
// * Nil (null)
// * Object
// * Func
type Val interface {
	Converter
	Comparer
	Arithmetic
	dumper
}

// The LogicProcessor interface defines the method required to implement
// boolean logic. It is defined as an interface for pluggable replacement
// for testing.
type LogicProcessor interface {
	Not(v Val) Bool
	And(x, y Val) Bool
	Or(x, y Val) Bool
}

// The default implementation of the logic processor.
type defaultLogic struct{}

func (d defaultLogic) Not(v Val) Bool {
	return Bool(!v.Bool())
}

func (d defaultLogic) And(x, y Val) Bool {
	return Bool(x.Bool() && y.Bool())
}

func (d defaultLogic) Or(x, y Val) Bool {
	return Bool(x.Bool() || y.Bool())
}
