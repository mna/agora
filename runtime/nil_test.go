package runtime

import (
	"testing"
)

func TestNilAsBool(t *testing.T) {
	res := Nil.Bool()
	if res != false {
		t.Errorf("Nil as bool : expected %v, got %v", false, res)
	}
}

func TestNilAsString(t *testing.T) {
	res := Nil.String()
	if res != NilString {
		t.Errorf("Nil as string : expected %s, got %s", NilString, res)
	}
}

func TestInvalidOpNil(t *testing.T) {
	assert := func(exp error) {
		if err := recover(); err != exp {
			t.Errorf("expected panic with error '%s', got '%v'", exp, err)
		}
	}

	func() {
		defer assert(ErrInvalidOpAddOnNil)
		Nil.Add(Bool(false))
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidOpSubOnNil)
		Nil.Sub(Bool(false))
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidOpMulOnNil)
		Nil.Mul(Bool(false))
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidOpDivOnNil)
		Nil.Div(Bool(false))
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidOpModOnNil)
		Nil.Mod(Bool(false))
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidOpUnmOnNil)
		Nil.Unm()
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidConvNilToInt)
		Nil.Int()
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidConvNilToFloat)
		Nil.Float()
		panic(nil)
	}()
}
