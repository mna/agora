package runtime

import (
	"errors"
)

var (
	ErrInvalidConvObjToInt    = errors.New("cannot convert Object to Int")
	ErrInvalidConvObjToFloat  = errors.New("cannot convert Object to Float")
	ErrInvalidConvObjToString = errors.New("cannot convert Object to String")
	ErrInvalidConvObjToBool   = errors.New("cannot convert Object to Bool")
	ErrInvalidConvObjToNative = errors.New("cannot convert Object to Native")
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
	panic(ErrInvalidConvObjToBool)
}

func (ø Object) Native() interface{} {
	if o, ok := ø["__toNative"]; ok {
		return o
	}
	panic(ErrInvalidConvObjToNative)
}
