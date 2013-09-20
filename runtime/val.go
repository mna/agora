package runtime

import (
	"fmt"
)

type TypeError string

func (te TypeError) Error() string {
	return string(te)
}

func NewTypeError(t, op string) TypeError {
	return TypeError(fmt.Sprintf("type error: %s not allowed with type %s", op, t))
}

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
	Add(Val, Val) Val
	Sub(Val, Val) Val
	Mul(Val, Val) Val
	Div(Val, Val) Val
	Mod(Val, Val) Val
	Unm(Val) Val
}

type defaultArithmetic struct{}

func (ar defaultArithmetic) binaryOp(l, r Val, op string, allowStrings bool) Val {
	lt, rt := Type(l), Type(r)
	mm := "__" + op
	if lt == "number" && rt == "number" {
		// Two numbers, standard arithmetic operation
		switch op {
		case "add":
			return Number(l.Float() + r.Float())
		case "sub":
			return Number(l.Float() - r.Float())
		case "mul":
			return Number(l.Float() * r.Float())
		case "div":
			return Number(l.Float() / r.Float())
		case "mod":
			return Number(l.Int() % r.Int())
		}
	} else if allowStrings && lt == "string" && rt == "string" {
		// Two strings
		switch op {
		case "add":
			return String(l.String() + r.String())
		}
	} else if lt == "object" {
		// If left operand is an object with a meta-method
		lo := l.(Object)
		if v, ok := lo.callMetaMethod(mm, r, Bool(true)); ok {
			return v
		}
	}
	// Last chance: if right operand is an object with a meta-method
	if rt == "object" {
		ro := r.(Object)
		if v, ok := ro.callMetaMethod(mm, l, Bool(false)); ok {
			return v
		}
	}
	panic(NewTypeError(lt, op))
}

func (ar defaultArithmetic) Add(l, r Val) Val {
	return ar.binaryOp(l, r, "add", true)
}

func (ar defaultArithmetic) Sub(l, r Val) Val {
	return ar.binaryOp(l, r, "sub", false)
}

func (ar defaultArithmetic) Mul(l, r Val) Val {
	return ar.binaryOp(l, r, "mul", false)
}

func (ar defaultArithmetic) Div(l, r Val) Val {
	return ar.binaryOp(l, r, "div", false)
}

func (ar defaultArithmetic) Mod(l, r Val) Val {
	return ar.binaryOp(l, r, "mod", false)
}

func (ar defaultArithmetic) Unm(l Val) Val {
	lt := Type(l)
	if lt == "number" {
		return Number(-l.Float())
	} else if lt == "object" {
		lo := l.(Object)
		if v, ok := lo.callMetaMethod("__unm"); ok {
			return v
		}
	}
	panic(NewTypeError(lt, "unm"))
}

// Comparer defines the method required to compare two Values.
// Cmp() returns 1 if the method receiver value is greater, 0 if
// it is equal, and -1 if it is lower.
type Comparer interface {
	Cmp(Val, Val) int
}

var (
	// Unmutable, this would be a const if it was possible
	uneqMatrix = map[string]map[string]int{
		"nil": map[string]int{
			"number": -1,
			"string": -1,
			"bool":   -1,
			"object": -1,
			"func":   -1,
			"custom": -1,
		},
		"number": map[string]int{
			"nil":    1,
			"string": -1,
			"bool":   1,
			"object": -1,
			"func":   1,
			"custom": 1,
		},
		"string": map[string]int{
			"number": 1,
			"nil":    1,
			"bool":   1,
			"object": 1,
			"func":   1,
			"custom": 1,
		},
		"bool": map[string]int{
			"number": -1,
			"string": -1,
			"nil":    1,
			"object": -1,
			"func":   -1,
			"custom": -1,
		},
		"object": map[string]int{
			"number": 1,
			"string": -1,
			"bool":   1,
			"nil":    1,
			"func":   1,
			"custom": 1,
		},
		"func": map[string]int{
			"number": -1,
			"string": -1,
			"bool":   1,
			"object": -1,
			"nil":    1,
			"custom": 1,
		},
		"custom": map[string]int{
			"number": -1,
			"string": -1,
			"bool":   1,
			"object": -1,
			"func":   -1,
			"nil":    1,
		},
	}
)

type defaultComparer struct{}

func (dc defaultComparer) Cmp(l, r Val) int {
	lt, rt := Type(l), Type(r)
	if lt == rt {
		// Comparable types
		switch lt {
		case "nil":
			return 0
		case "number":
			lf, rf := l.Float(), r.Float()
			if lf == rf {
				return 0
			} else if lf < rf {
				return -1
			} else {
				return 1
			}
		case "string":
			ls, rs := l.String(), r.String()
			if ls == rs {
				return 0
			} else if ls < rs {
				return -1
			} else {
				return 1
			}
		case "bool":
			lb, rb := l.Bool(), r.Bool()
			if lb == rb {
				return 0
			} else if lb {
				return 1 // true is greater than false (0)
			} else {
				return -1
			}
		case "func":
			lf, rf := l.Native(), r.Native()
			if lf == rf {
				return 0
			} else {
				// "greater" or "lower" has no sense for funcs, return -1
				return -1
			}
		case "object":
			// If left has meta method, use left, otherwise right, else compare
			lo, ro := l.(Object), r.(Object)
			if v, ok := lo.callMetaMethod("__cmp", r, Bool(true)); ok {
				return int(v.Int())
			}
			if v, ok := ro.callMetaMethod("__cmp", l, Bool(false)); ok {
				return int(v.Int())
			}
			if lo == ro {
				return 0
			} else {
				// "greater" or "lower" has no sense for objects, return -1
				return -1
			}
		case "custom":
			if l == r {
				return 0
			} else {
				// "greater" or "lower" has no sense for custom vals, return -1
				return -1
			}
		default:
			panic(NewTypeError(lt, "cmp"))
		}
	} else {
		// Uncomparable types, first check for meta-methods
		var o Object
		var isLeft bool
		var otherv Val
		if lt == "object" {
			o = l.(Object)
			isLeft = true
			otherv = r
		} else if rt == "object" {
			o = r.(Object)
			isLeft = false
			otherv = l
		}
		if o != nil {
			if v, ok := o.callMetaMethod("__cmp", otherv, Bool(isLeft)); ok {
				return int(v.Int())
			}
		}
		// Else, return arbitrary but constant result
		return uneqMatrix[lt][rt]
	}
}

// The dumper interface defines the required behaviour to pretty-print
// the values.
// TODO : Should provide a default impl or make it public, so that custom types
// can actually be built from outside the package.
type dumper interface {
	dump() string
}

// Val is the representation of a value, any value, in the language.
// The supported value types are the following:
// * Number (float64)
// * String
// * Boolean (Bool)
// * Nil (null)
// * Object
// * Func
type Val interface {
	Converter
	dumper
}

func Type(v Val) string {
	switch v.(type) {
	case String:
		return "string"
	case Number:
		return "number"
	case Bool:
		return "bool"
	case Func:
		return "func"
	case Object:
		return "object"
	default:
		if v == Nil {
			return "nil"
		} else {
			return "custom"
		}
	}
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
