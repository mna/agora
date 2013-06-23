package runtime

import (
	"errors"
)

var (
	ErrInvalidConvObjToInt    = errors.New("cannot convert Object to Int")
	ErrInvalidConvObjToFloat  = errors.New("cannot convert Object to Float")
	ErrInvalidConvObjToString = errors.New("cannot convert Object to String")
	ErrInvalidConvObjToNative = errors.New("cannot convert Object to Native")

	ErrInvalidOpAddOnObj = errors.New("cannot apply Add on this Object")
	ErrInvalidOpSubOnObj = errors.New("cannot apply Sub on this Object")
	ErrInvalidOpMulOnObj = errors.New("cannot apply Mul on this Object")
	ErrInvalidOpDivOnObj = errors.New("cannot apply Div on this Object")
	ErrInvalidOpPowOnObj = errors.New("cannot apply Pow on this Object")
	ErrInvalidOpModOnObj = errors.New("cannot apply Mod on this Object")
	ErrInvalidOpUnmOnObj = errors.New("cannot apply Unm on this Object")
)

type Object struct {
	m map[string]Val
}

func (ø Object) Int() int {
	if i, ok := ø["__toInt"]; ok {
		return i
	}
	panic(ErrInvalidConvObjToInt)
}

func (ø Object) Float() float64 {
	if f, ok := ø["__toFloat"]; ok {
		return f
	}
	panic(ErrInvalidConvObjToFloat)
}

func (ø Object) String() string {
	if s, ok := ø["__toString"]; ok {
		return s
	}
	panic(ErrInvalidConvObjToString)
}

func (ø Object) Bool() bool {
	if b, ok := ø["__toBool"]; ok {
		return b
	}
	// If __toBool is not defined, object returns true (since it is not nil)
	return true
}

func (ø Object) Native() interface{} {
	if o, ok := ø["__toNative"]; ok {
		return o // TODO : Need to call the method, not return it!
	}
	panic(ErrInvalidConvObjToNative)
}

func (ø Object) Cmp(v Val) int {
	// First check for a custom Cmp method
	if c, ok := ø["__cmp"]; ok {
		return c.(*Func).Call(v) // TODO : Needs to setup `this`?
	}

	// If same reference as v, return 0 (equal)
	if ø == v {
		return 0
	}
	// Otherwise, return -1
	return -1
}

// TODO : NotFound method, how does it work? See js proposition, Lua, others...
