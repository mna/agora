package runtime

import (
	"testing"
)

const (
	floatCompareBuffer = 1e-6
)

func TestArithmetic(t *testing.T) {
	ctx := NewCtx(nil, nil)
	ari := defaultArithmetic{}
	o := NewObject()
	oplus := NewObject()
	oplus.Set(String("__add"), NewNativeFunc(ctx, "", func(args ...Val) Val {
		ExpectAtLeastNArgs(2, args)
		return args[0]
	}))
	fn := NewNativeFunc(ctx, "", func(_ ...Val) Val { return Nil })
	checkPanic := func(lbl string, i int, p bool) {
		if e := recover(); (e != nil) != p {
			if p {
				t.Errorf("[%s %d] - expected error, got none", lbl, i)
			} else {
				t.Errorf("[%s %d] - expected no error, got %s", lbl, i, e)
			}
		}
	}

	// Add
	adds := []struct {
		l, r, exp Val
		err       bool
	}{
		{l: Number(2), r: Number(5), exp: Number(7)},
		{l: Number(-2), r: Number(5.123), exp: Number(3.123)},
		{l: Number(2.24), r: Number(0.01), exp: Number(2.25)},
		{l: Number(0), r: Number(0.0), exp: Number(0)},
		{l: String("hi"), r: String("you"), exp: String("hiyou")},
		{l: String("0"), r: String("2"), exp: String("02")},
		{l: String(""), r: String(""), exp: String("")},
		{l: Nil, r: Nil, err: true},
		{l: Nil, r: Number(2), err: true},
		{l: Nil, r: String("test"), err: true},
		{l: Nil, r: Bool(true), err: true},
		{l: Nil, r: o, err: true},
		{l: Nil, r: oplus, exp: Nil},
		{l: Nil, r: fn, err: true},
		// TODO: Custom
		{l: Number(2), r: Nil, err: true},
		{l: Number(2), r: String("test"), err: true},
		{l: Number(2), r: Bool(true), err: true},
		{l: Number(2), r: o, err: true},
		{l: Number(2), r: oplus, exp: Number(2)},
		{l: Number(2), r: fn, err: true},
		// TODO: Custom
		{l: String("ok"), r: Nil, err: true},
		{l: String("ok"), r: Number(2), err: true},
		{l: String("ok"), r: Bool(true), err: true},
		{l: String("ok"), r: o, err: true},
		{l: String("ok"), r: oplus, exp: String("ok")},
		{l: String("ok"), r: fn, err: true},
		// TODO: Custom
		{l: Bool(true), r: Nil, err: true},
		{l: Bool(true), r: Number(2), err: true},
		{l: Bool(true), r: String("test"), err: true},
		{l: Bool(true), r: Bool(true), err: true},
		{l: Bool(true), r: o, err: true},
		{l: Bool(true), r: oplus, exp: Bool(true)},
		{l: Bool(true), r: fn, err: true},
		// TODO: Custom
		{l: oplus, r: Nil, exp: Nil},
		{l: oplus, r: Number(2), exp: Number(2)},
		{l: oplus, r: String("test"), exp: String("test")},
		{l: oplus, r: Bool(true), exp: Bool(true)},
		{l: oplus, r: o, exp: o},
		{l: oplus, r: oplus, exp: oplus},
		{l: oplus, r: fn, exp: fn},
		// TODO: Custom
		{l: o, r: Nil, err: true},
		{l: o, r: Number(2), err: true},
		{l: o, r: String("test"), err: true},
		{l: o, r: Bool(true), err: true},
		{l: o, r: o, err: true},
		{l: o, r: oplus, exp: o},
		{l: o, r: fn, err: true},
		// TODO: Custom
		{l: fn, r: Nil, err: true},
		{l: fn, r: Number(2), err: true},
		{l: fn, r: String("test"), err: true},
		{l: fn, r: Bool(true), err: true},
		{l: fn, r: o, err: true},
		{l: fn, r: oplus, exp: fn},
		{l: fn, r: fn, err: true},
		// TODO: Custom
	}
	for i, c := range adds {
		func() {
			defer checkPanic("add", i, c.err)
			ret := ari.Add(c.l, c.r)
			if ret != c.exp {
				t.Errorf("[add %d] - expected %s, got %s", i, c.exp, ret)
			}
		}()
	}
}
