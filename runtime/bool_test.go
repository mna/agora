package runtime

import (
	"testing"
)

func TestInvalidOpBool(t *testing.T) {
	assert := func(exp error) {
		if err := recover(); err != exp {
			t.Errorf("expected panic with error '%s', got '%v'", exp, err)
		}
	}

	b := Bool(true)
	func() {
		defer assert(ErrInvalidOpAddOnBool)
		b.Add(Bool(false))
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidOpSubOnBool)
		b.Sub(Bool(false))
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidOpMulOnBool)
		b.Mul(Bool(false))
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidOpDivOnBool)
		b.Div(Bool(false))
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidOpModOnBool)
		b.Mod(Bool(false))
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidOpPowOnBool)
		b.Pow(Bool(false))
		panic(nil)
	}()
	func() {
		defer assert(ErrInvalidOpUnmOnBool)
		b.Unm()
		panic(nil)
	}()
}

func TestNotBool(t *testing.T) {
	cases := []struct {
		x   bool
		exp bool
	}{
		{x: true, exp: false},
		{x: false, exp: true},
	}

	for _, c := range cases {
		vx := Bool(c.x)
		res := vx.Not()
		if bres := bool(res.(Bool)); c.exp != bres {
			t.Errorf("!%v : expected %v, got %v", c.x, c.exp, bres)
		}
	}
}
