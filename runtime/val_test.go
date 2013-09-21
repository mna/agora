package runtime

import (
	"math"
	"testing"
)

const (
	floatCompareBuffer = 1e-6
)

// Test cases for arithmetic
type arithCase struct {
	l, r, exp Val
	err       bool
}

// Define a custom type
type cusType struct{}

func (c cusType) Int() int64 {
	return 1
}
func (c cusType) Float() float64 {
	return 1.0
}
func (c cusType) String() string {
	return "cus!"
}
func (c cusType) Bool() bool {
	return true
}
func (c cusType) Native() interface{} {
	return c
}
func (c cusType) dump() string {
	return c.String()
}

var (
	// Common variables for the tests
	ctx   = NewCtx(nil, nil)
	ari   = defaultArithmetic{}
	o     = NewObject()
	oplus = NewObject()
	fn    = NewNativeFunc(ctx, "", func(_ ...Val) Val { return Nil })
	cus   = cusType{}

	// Common cases, same result regardless of operation
	common = []arithCase{
		{l: Nil, r: Nil, err: true},
		{l: Nil, r: Number(2), err: true},
		{l: Nil, r: String("test"), err: true},
		{l: Nil, r: Bool(true), err: true},
		{l: Nil, r: o, err: true},
		{l: Nil, r: oplus, exp: Nil},
		{l: Nil, r: fn, err: true},
		{l: Nil, r: cusType{}, err: true},
		{l: Number(2), r: Nil, err: true},
		{l: Number(2), r: String("test"), err: true},
		{l: Number(2), r: Bool(true), err: true},
		{l: Number(2), r: o, err: true},
		{l: Number(2), r: oplus, exp: Number(2)},
		{l: Number(2), r: fn, err: true},
		{l: Number(2), r: cusType{}, err: true},
		{l: String("ok"), r: Nil, err: true},
		{l: String("ok"), r: Number(2), err: true},
		{l: String("ok"), r: Bool(true), err: true},
		{l: String("ok"), r: o, err: true},
		{l: String("ok"), r: oplus, exp: String("ok")},
		{l: String("ok"), r: fn, err: true},
		{l: String("ok"), r: cusType{}, err: true},
		{l: Bool(true), r: Nil, err: true},
		{l: Bool(true), r: Number(2), err: true},
		{l: Bool(true), r: String("test"), err: true},
		{l: Bool(true), r: Bool(true), err: true},
		{l: Bool(true), r: o, err: true},
		{l: Bool(true), r: oplus, exp: Bool(true)},
		{l: Bool(true), r: fn, err: true},
		{l: Bool(true), r: cusType{}, err: true},
		{l: oplus, r: Nil, exp: Nil},
		{l: oplus, r: Number(2), exp: Number(2)},
		{l: oplus, r: String("test"), exp: String("test")},
		{l: oplus, r: Bool(true), exp: Bool(true)},
		{l: oplus, r: o, exp: o},
		{l: oplus, r: oplus, exp: oplus},
		{l: oplus, r: fn, exp: fn},
		{l: oplus, r: cus, exp: cus},
		{l: o, r: Nil, err: true},
		{l: o, r: Number(2), err: true},
		{l: o, r: String("test"), err: true},
		{l: o, r: Bool(true), err: true},
		{l: o, r: o, err: true},
		{l: o, r: oplus, exp: o},
		{l: o, r: fn, err: true},
		{l: o, r: cusType{}, err: true},
		{l: fn, r: Nil, err: true},
		{l: fn, r: Number(2), err: true},
		{l: fn, r: String("test"), err: true},
		{l: fn, r: Bool(true), err: true},
		{l: fn, r: o, err: true},
		{l: fn, r: oplus, exp: fn},
		{l: fn, r: fn, err: true},
		{l: fn, r: cusType{}, err: true},
		{l: cus, r: Nil, err: true},
		{l: cus, r: Number(2), err: true},
		{l: cus, r: String("test"), err: true},
		{l: cus, r: Bool(true), err: true},
		{l: cus, r: o, err: true},
		{l: cus, r: oplus, exp: cus},
		{l: cus, r: fn, err: true},
		{l: cus, r: cusType{}, err: true},
	}

	// Add-specific cases
	adds = append(common, []arithCase{
		{l: Number(2), r: Number(5), exp: Number(7)},
		{l: Number(-2), r: Number(5.123), exp: Number(3.123)},
		{l: Number(2.24), r: Number(0.01), exp: Number(2.25)},
		{l: Number(0), r: Number(0.0), exp: Number(0)},
		{l: String("hi"), r: String("you"), exp: String("hiyou")},
		{l: String("0"), r: String("2"), exp: String("02")},
		{l: String(""), r: String(""), exp: String("")},
	}...)

	// Sub-specific cases
	subs = append(common, []arithCase{
		{l: Number(5), r: Number(2), exp: Number(3)},
		{l: Number(-2), r: Number(5.123), exp: Number(-7.123)},
		{l: Number(2.24), r: Number(0.01), exp: Number(2.23)},
		{l: Number(0), r: Number(0.0), exp: Number(0)},
		{l: String("hi"), r: String("you"), err: true},
	}...)

	// Mul-specific cases
	muls = append(common, []arithCase{
		{l: Number(5), r: Number(2), exp: Number(10)},
		{l: Number(-2), r: Number(5.123), exp: Number(-10.246)},
		{l: Number(2.24), r: Number(0.01), exp: Number(0.0224)},
		{l: Number(0), r: Number(0.0), exp: Number(0)},
		{l: String("hi"), r: String("you"), err: true},
	}...)

	// Div-specific cases
	divs = append(common, []arithCase{
		{l: Number(5), r: Number(2), exp: Number(2.5)},
		{l: Number(-2), r: Number(5.123), exp: Number(-0.390396252)},
		{l: Number(2.24), r: Number(0.01), exp: Number(224)},
		{l: Number(0), r: Number(0.0), exp: Number(math.NaN())},
		{l: String("hi"), r: String("you"), err: true},
	}...)

	// Mod-specific cases
	mods = append(common, []arithCase{
		{l: Number(5), r: Number(2), exp: Number(1)},
		{l: Number(-2), r: Number(5.123), exp: Number(-2)},
		{l: Number(2.24), r: Number(1.1), exp: Number(0)},
		{l: Number(0), r: Number(0.0), err: true},
		{l: String("hi"), r: String("you"), err: true},
	}...)

	// Unm-specific cases
	unms = []arithCase{
		{l: Nil, err: true},
		{l: Number(4), exp: Number(-4)},
		{l: Number(-3.1415), exp: Number(3.1415)},
		{l: Number(0), exp: Number(0)},
		{l: String("ok"), err: true},
		{l: Bool(false), err: true},
		{l: oplus, exp: Number(-1)},
		{l: o, err: true},
		{l: fn, err: true},
		{l: cus, err: true},
	}
)

func init() {
	fRetArg := NewNativeFunc(ctx, "", func(args ...Val) Val {
		ExpectAtLeastNArgs(2, args)
		return args[0]
	})
	fRetUnm := NewNativeFunc(ctx, "", func(args ...Val) Val {
		return Number(-1)
	})
	oplus.Set(String("__add"), fRetArg)
	oplus.Set(String("__sub"), fRetArg)
	oplus.Set(String("__mul"), fRetArg)
	oplus.Set(String("__div"), fRetArg)
	oplus.Set(String("__mod"), fRetArg)
	oplus.Set(String("__unm"), fRetUnm)
}

func TestType(t *testing.T) {
	cases := []struct {
		src Val
		exp string
	}{
		{src: Nil, exp: "nil"},
		{src: Bool(true), exp: "bool"},
		{src: Bool(false), exp: "bool"},
		{src: Number(1), exp: "number"},
		{src: Number(3.1415), exp: "number"},
		{src: Number(0.0), exp: "number"},
		{src: String("ok"), exp: "string"},
		{src: String(""), exp: "string"},
		{src: fn, exp: "func"},
		{src: o, exp: "object"},
		{src: oplus, exp: "object"},
		{src: cusType{}, exp: "custom"},
	}
	for i, c := range cases {
		got := Type(c.src)
		if got != c.exp {
			t.Errorf("[%d] - expected %s, got %s", i, c.exp, got)
		}
	}
}

func TestArithmetic(t *testing.T) {
	checkPanic := func(lbl string, i int, p bool) {
		if e := recover(); (e != nil) != p {
			if p {
				t.Errorf("[%s %d] - expected error, got none", lbl, i)
			} else {
				t.Errorf("[%s %d] - expected no error, got %s", lbl, i, e)
			}
		}
	}
	cases := map[string][]arithCase{
		"add": adds,
		"sub": subs,
		"mul": muls,
		"div": divs,
		"mod": mods,
		"unm": unms,
	}
	for k, v := range cases {
		for i, c := range v {
			func() {
				defer checkPanic(k, i, c.err)
				var ret Val
				switch k {
				case "add":
					ret = ari.Add(c.l, c.r)
				case "sub":
					ret = ari.Sub(c.l, c.r)
				case "mul":
					ret = ari.Mul(c.l, c.r)
				case "div":
					ret = ari.Div(c.l, c.r)
				case "mod":
					ret = ari.Mod(c.l, c.r)
				case "unm":
					ret = ari.Unm(c.l)
				}
				if _, ok := ret.(Number); ok {
					if math.Abs(ret.Float()-c.exp.Float()) > floatCompareBuffer {
						t.Errorf("[%s %d] - expected %s, got %s", k, i, c.exp, ret)
					}
				} else if ret != c.exp {
					t.Errorf("[%s %d] - expected %s, got %s", k, i, c.exp, ret)
				}
			}()
		}
	}
}
