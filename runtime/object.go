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
	ErrInvalidOpModOnObj = errors.New("cannot apply Mod on this Object")
	ErrInvalidOpUnmOnObj = errors.New("cannot apply Unm on this Object")

	ErrNoSuchMethod     = errors.New("method does not exist")
	ErrFieldNotFunction = errors.New("field is not a function")
)

type Object struct {
	m map[Val]Val
}

func NewObject() *Object {
	return &Object{
		make(map[Val]Val),
	}
}

func (ø *Object) dump() string {
	buf := bytes.NewBuffer(nil)
	for k, v := range ø.m {
		buf.WriteString(fmt.Sprintf(" %s: %s, ", k.dump(), v.dump()))
	}
	return fmt.Sprintf("{%s} (Object)", buf)
}

func (ø *Object) Int() int {
	if i, ok := ø.m[String("__toInt")]; ok {
		if f, ok := i.(Func); ok {
			return f.Call(ø).Int()
		}
	}
	panic(ErrInvalidConvObjToInt)
}

func (ø *Object) Float() float64 {
	if l, ok := ø.m[String("__toFloat")]; ok {
		if f, ok := l.(Func); ok {
			return f.Call(ø).Float()
		}
	}
	panic(ErrInvalidConvObjToFloat)
}

func (ø *Object) String() string {
	if s, ok := ø.m[String("__toString")]; ok {
		if f, ok := s.(Func); ok {
			return f.Call(ø).String()
		}
	}
	panic(ErrInvalidConvObjToString)
}

func (ø *Object) Bool() bool {
	if b, ok := ø.m[String("__toBool")]; ok {
		if f, ok := b.(Func); ok {
			return f.Call(ø).Bool()
		}
	}
	// If __toBool is not defined, object returns true (since it is not nil)
	return true
}

func (ø *Object) Native() interface{} {
	if o, ok := ø.m[String("__toNative")]; ok {
		if f, ok := o.(Func); ok {
			return f.Call(ø).Native()
		}
	}
	panic(ErrInvalidConvObjToNative)
}

func (ø *Object) Cmp(v Val) int {
	// First check for a custom Cmp method
	if c, ok := ø.m[String("__cmp")]; ok {
		if f, ok := c.(Func); ok {
			return f.Call(ø, v).Int()
		}
	}
	// Else, default Cmp - if same reference as v, return 0 (equal)
	if ø == v {
		return 0
	}
	return -1
}

func (ø *Object) callBinaryMethod(nm String, err error, v Val) Val {
	if m, ok := ø.m[nm]; ok {
		if f, ok := m.(Func); ok {
			return f.Call(ø, v)
		}
	}
	panic(err)
}

func (ø *Object) Add(v Val) Val {
	return ø.callBinaryMethod(String("__add"), ErrInvalidOpAddOnObj, v)
}

func (ø *Object) Sub(v Val) Val {
	return ø.callBinaryMethod(String("__sub"), ErrInvalidOpSubOnObj, v)
}

func (ø *Object) Mul(v Val) Val {
	return ø.callBinaryMethod(String("__mul"), ErrInvalidOpMulOnObj, v)
}

func (ø *Object) Div(v Val) Val {
	return ø.callBinaryMethod(String("__div"), ErrInvalidOpDivOnObj, v)
}

func (ø *Object) Mod(v Val) Val {
	return ø.callBinaryMethod(String("__mod"), ErrInvalidOpModOnObj, v)
}

func (ø *Object) Unm() Val {
	if m, ok := ø.m[String("__unm")]; ok {
		if f, ok := m.(Func); ok {
			return f.Call(ø)
		}
	}
	panic(ErrInvalidOpUnmOnObj)
}

func (ø *Object) Get(key Val) Val {
	if v, ok := ø.m[key]; ok {
		return v
	}
	return Nil
}

func (ø *Object) Set(key Val, v Val) {
	ø.m[key] = v
}

func (ø *Object) callMethod(nm Val, args ...Val) Val {
	v, ok := ø.m[nm]
	if ok {
		if f, ok := v.(Func); ok {
			return f.Call(ø, args...)
		} else {
			panic(ErrFieldNotFunction)
		}
	} else {
		// Method not found - call __noSuchMethod if it exists, otherwise panic
		if m, ok := ø.m[String("__noSuchMethod")]; ok {
			if f, ok := m.(Func); ok {
				args = append([]Val{nm}, args...)
				return f.Call(ø, args...)
			}
		}
		panic(ErrNoSuchMethod)
	}
}
