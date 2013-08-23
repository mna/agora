package runtime

import (
	"testing"
)

func TestBoolAsInt(t *testing.T) {
	cases := []struct {
		x   bool
		exp int
	}{
		{x: true, exp: 1},
		{x: false, exp: 0},
	}

	for _, c := range cases {
		vx := Bool(c.x)
		res := vx.Int()
		if c.exp != res {
			t.Errorf("%v as int : expected %d, got %d", c.x, c.exp, res)
		}
	}
}

func TestBoolAsFloat(t *testing.T) {
	cases := []struct {
		x   bool
		exp float64
	}{
		{x: true, exp: 1.0},
		{x: false, exp: 0.0},
	}

	for _, c := range cases {
		vx := Bool(c.x)
		res := vx.Float()
		if c.exp != res {
			t.Errorf("%v as float : expected %f, got %f", c.x, c.exp, res)
		}
	}
}

func TestBoolAsBool(t *testing.T) {
	cases := []struct {
		x   bool
		exp bool
	}{
		{x: true, exp: true},
		{x: false, exp: false},
	}

	for _, c := range cases {
		vx := Bool(c.x)
		res := vx.Bool()
		if c.exp != res {
			t.Errorf("%v as bool : expected %v, got %v", c.x, c.exp, res)
		}
	}
}

func TestBoolAsString(t *testing.T) {
	cases := []struct {
		x   bool
		exp string
	}{
		{x: true, exp: "true"},
		{x: false, exp: "false"},
	}

	for _, c := range cases {
		vx := Bool(c.x)
		res := vx.String()
		if c.exp != res {
			t.Errorf("%v as string : expected %s, got %s", c.x, c.exp, res)
		}
	}
}

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
