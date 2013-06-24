package runtime

import (
	"bytes"
	"errors"
	"fmt"
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

	ErrNoSuchMethod     = errors.New("method does not exist")
	ErrFieldNotFunction = errors.New("field is not a function")
)

type Object struct {
	m map[string]Val
}

func newObject() *Object {
	return &Object{
		make(map[string]Val),
	}
}

func (ø *Object) dump() string {
	buf := bytes.NewBuffer(nil)
	for k, v := range ø.m {
		buf.WriteString(fmt.Sprintf(" \"%s\": %s, ", k, v.dump()))
	}
	return fmt.Sprintf("{%s} (Object)", buf)
}

func (ø *Object) Int() int {
	if i, ok := ø.m["__toInt"]; ok {
		if f, ok := i.(*Func); ok {
			f.This = ø
			return f.Call().Int()
		}
	}
	panic(ErrInvalidConvObjToInt)
}

func (ø *Object) Float() float64 {
	if l, ok := ø.m["__toFloat"]; ok {
		if f, ok := l.(*Func); ok {
			f.This = ø
			return f.Call().Float()
		}
	}
	panic(ErrInvalidConvObjToFloat)
}

func (ø *Object) String() string {
	if s, ok := ø.m["__toString"]; ok {
		if f, ok := s.(*Func); ok {
			f.This = ø
			return f.Call().String()
		}
	}
	panic(ErrInvalidConvObjToString)
}

func (ø *Object) Bool() bool {
	if b, ok := ø.m["__toBool"]; ok {
		if f, ok := b.(*Func); ok {
			f.This = ø
			return f.Call().Bool()
		}
	}
	// If __toBool is not defined, object returns true (since it is not nil)
	return true
}

func (ø *Object) Native() interface{} {
	if o, ok := ø.m["__toNative"]; ok {
		if f, ok := o.(*Func); ok {
			f.This = ø
			return f.Call().Native()
		}
	}
	panic(ErrInvalidConvObjToNative)
}

func (ø *Object) Cmp(v Val) int {
	// First check for a custom Cmp method
	if c, ok := ø.m["__cmp"]; ok {
		if f, ok := c.(*Func); ok {
			f.This = ø
			return f.Call(v).Int()
		}
	}
	// Else, default Cmp - if same reference as v, return 0 (equal)
	if ø == v {
		return 0
	}
	// Otherwise, return -1
	return -1
}

func (ø *Object) callBinaryMethod(nm string, err error, v Val) Val {
	if m, ok := ø.m[nm]; ok {
		if f, ok := m.(*Func); ok {
			f.This = ø
			return f.Call(v)
		}
	}
	panic(err)
}

func (ø *Object) Add(v Val) Val {
	return ø.callBinaryMethod("__add", ErrInvalidOpAddOnObj, v)
}

func (ø *Object) Sub(v Val) Val {
	return ø.callBinaryMethod("__sub", ErrInvalidOpSubOnObj, v)
}

func (ø *Object) Mul(v Val) Val {
	return ø.callBinaryMethod("__mul", ErrInvalidOpMulOnObj, v)
}

func (ø *Object) Div(v Val) Val {
	return ø.callBinaryMethod("__div", ErrInvalidOpDivOnObj, v)
}

func (ø *Object) Mod(v Val) Val {
	return ø.callBinaryMethod("__mod", ErrInvalidOpModOnObj, v)
}

func (ø *Object) Pow(v Val) Val {
	return ø.callBinaryMethod("__pow", ErrInvalidOpPowOnObj, v)
}

func (ø *Object) Unm() Val {
	if m, ok := ø.m["__unm"]; ok {
		if f, ok := m.(*Func); ok {
			f.This = ø
			return f.Call()
		}
	}
	panic(ErrInvalidOpUnmOnObj)
}

func (ø *Object) get(key string) Val {
	v, ok := ø.m[key]
	if !ok {
		return Nil
	}
	return v
}

func (ø *Object) set(key string, v Val) {
	ø.m[key] = v
}

func (ø *Object) callMethod(v Val, args ...Val) Val {
	switch f := v.(type) {
	case *Func:
		// Call the method
		f.This = ø
		return f.Call(args...)
	case null:
		// Method not found - call __noSuchMethod if it exists, otherwise panic
		if m, ok := ø.m["__noSuchMethod"]; ok {
			if f, ok := m.(*Func); ok {
				f.This = ø
				return f.Call(args...) // TODO : pass method name as first value
			}
		}
		panic(ErrNoSuchMethod)
	default:
		// Any other case: not a function
		panic(ErrFieldNotFunction)
	}
}
