package runtime

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	// Predefined errors
	ErrInvalidConvObjToInt    = errors.New("cannot convert Object to Int")
	ErrInvalidConvObjToFloat  = errors.New("cannot convert Object to Float")
	ErrInvalidConvObjToString = errors.New("cannot convert Object to String")

	ErrInvalidOpAddOnObj = errors.New("cannot apply Add on this Object")
	ErrInvalidOpSubOnObj = errors.New("cannot apply Sub on this Object")
	ErrInvalidOpMulOnObj = errors.New("cannot apply Mul on this Object")
	ErrInvalidOpDivOnObj = errors.New("cannot apply Div on this Object")
	ErrInvalidOpModOnObj = errors.New("cannot apply Mod on this Object")
	ErrInvalidOpUnmOnObj = errors.New("cannot apply Unm on this Object")

	ErrNoSuchMethod     = errors.New("method does not exist")
	ErrFieldNotFunction = errors.New("field is not a function")
)

// An Object is a map of values, an associative array.
type Object struct {
	m map[Val]Val
}

// NewObject returns a new instance of an object.
func NewObject() *Object {
	return &Object{
		make(map[Val]Val),
	}
}

// dump pretty-prints the content of the object.
func (o *Object) dump() string {
	buf := bytes.NewBuffer(nil)
	for k, v := range o.m {
		buf.WriteString(fmt.Sprintf(" %s: %s, ", k.dump(), v.dump()))
	}
	return fmt.Sprintf("{%s} (Object)", buf)
}

// Int returns the integer value of the object. Such behaviour can be defined
// if a `__toInt` method is available on the object.
func (o *Object) Int() int {
	if i, ok := o.m[String("__toInt")]; ok {
		if f, ok := i.(Func); ok {
			return f.Call(o).Int()
		}
	}
	panic(ErrInvalidConvObjToInt)
}

// Float returns the float value of the object. Such behaviour can be defined
// if a `__toFloat` method is available on the object.
func (o *Object) Float() float64 {
	if l, ok := o.m[String("__toFloat")]; ok {
		if f, ok := l.(Func); ok {
			return f.Call(o).Float()
		}
	}
	panic(ErrInvalidConvObjToFloat)
}

// String returns the string value of the object. Such behaviour can be defined
// if a `__toString` method is available on the object.
func (o *Object) String() string {
	if s, ok := o.m[String("__toString")]; ok {
		if f, ok := s.(Func); ok {
			return f.Call(o).String()
		}
	}
	panic(ErrInvalidConvObjToString)
}

// Bool returns the boolean value of the object. Such behaviour can be defined
// if a `__toBool` method is available on the object. Otherwise it returns true.
func (o *Object) Bool() bool {
	if b, ok := o.m[String("__toBool")]; ok {
		if f, ok := b.(Func); ok {
			return f.Call(o).Bool()
		}
	}
	// If __toBool is not defined, object returns true (since it is not nil)
	return true
}

// Native returns the Go native value of the object. Such behaviour can be defined
// if a `__toNative` method is available on the object.
func (o *Object) Native() interface{} {
	if n, ok := o.m[String("__toNative")]; ok {
		if f, ok := n.(Func); ok {
			return f.Call(o).Native()
		}
	}
	return o.m
}

// Cmp compares the object to another value. Such behaviour can be defined
// if a `__cmp` method is available on the object. Otherwise it returns 0 if
// the compared value is the object, or -1 otherwise.
func (o *Object) Cmp(v Val) int {
	// First check for a custom Cmp method
	if c, ok := o.m[String("__cmp")]; ok {
		if f, ok := c.(Func); ok {
			return f.Call(o, v).Int()
		}
	}
	// Else, default Cmp - if same reference as v, return 0 (equal)
	if o == v {
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

// Add performs addition. Such behaviour can be defined
// if a `__add` method is available on the object.
func (o *Object) Add(v Val) Val {
	return o.callBinaryMethod(String("__add"), ErrInvalidOpAddOnObj, v)
}

// Sub performs subtraction. Such behaviour can be defined
// if a `__sub` method is available on the object.
func (o *Object) Sub(v Val) Val {
	return o.callBinaryMethod(String("__sub"), ErrInvalidOpSubOnObj, v)
}

// Mul performs multiplication. Such behaviour can be defined
// if a `__mul` method is available on the object.
func (o *Object) Mul(v Val) Val {
	return o.callBinaryMethod(String("__mul"), ErrInvalidOpMulOnObj, v)
}

// Div performs division. Such behaviour can be defined
// if a `__div` method is available on the object.
func (o *Object) Div(v Val) Val {
	return o.callBinaryMethod(String("__div"), ErrInvalidOpDivOnObj, v)
}

// Mod computes the modulo. Such behaviour can be defined
// if a `__mod` method is available on the object.
func (o *Object) Mod(v Val) Val {
	return o.callBinaryMethod(String("__mod"), ErrInvalidOpModOnObj, v)
}

// Unm computes the unary minus. Such behaviour can be defined
// if a `__unm` method is available on the object.
func (o *Object) Unm() Val {
	if m, ok := o.m[String("__unm")]; ok {
		if f, ok := m.(Func); ok {
			return f.Call(o)
		}
	}
	panic(ErrInvalidOpUnmOnObj)
}

// Get returns the value of the field identified by key. It returns Nil
// if the field does not exist.
func (o *Object) Get(key Val) Val {
	if v, ok := o.m[key]; ok {
		return v
	}
	return Nil
}

// Set assigns the value v to the field identified by key.
func (o *Object) Set(key Val, v Val) {
	o.m[key] = v
}

// callMethod calls the method identified by nm with the provided arguments.
// It panics if the field does not hold a function. If the field does not
// exist and a method named `__noSuchMethod` is defined, it is called instead.
func (o *Object) callMethod(nm Val, args ...Val) Val {
	v, ok := o.m[nm]
	if ok {
		if f, ok := v.(Func); ok {
			return f.Call(o, args...)
		} else {
			panic(ErrFieldNotFunction)
		}
	} else {
		// Method not found - call __noSuchMethod if it exists, otherwise panic
		if m, ok := o.m[String("__noSuchMethod")]; ok {
			if f, ok := m.(Func); ok {
				args = append([]Val{nm}, args...)
				return f.Call(o, args...)
			}
		}
		panic(ErrNoSuchMethod)
	}
}
