package runtime

import (
	"io"
	"strings"
	"testing"
)

func TestKeys(t *testing.T) {
	bi := new(builtinMod)
	ctx := NewCtx(nil, nil)
	bi.SetCtx(ctx)

	cases := []struct {
		src Val
		exp []Val
	}{
		0: {
			src: func() Object {
				o := NewObject()
				return o
			}(),
			exp: []Val{},
		},
		1: {
			src: func() Object {
				o := NewObject()
				o.Set(String("a"), Number(0))
				return o
			}(),
			exp: []Val{String("a")},
		},
		2: {
			src: func() Object {
				o := NewObject()
				o.Set(String("a"), Number(0))
				o.Set(String("b"), Number(0))
				o.Set(Number(1), Number(0))
				return o
			}(),
			exp: []Val{
				String("a"),
				String("b"),
				Number(1),
			},
		},
		3: {
			src: func() Object {
				o := NewObject()
				o.Set(String("a"), Number(0))
				o.Set(String("b"), Number(0))
				o.Set(Number(1), Number(0))
				o.Set(String("__keys"), NewNativeFunc(ctx, "", func(args ...Val) []Val {
					k := NewObject()
					k.Set(Number(0), String("b"))
					k.Set(Number(1), Number(1))
					return Set1(k)
				}))
				return o
			}(),
			exp: []Val{
				String("b"),
				Number(1),
			},
		},
	}

	for i, c := range cases {
		ret := Get1(bi._keys(c.src))
		ob := ret.(Object)
		l := ob.Len().Int()
		if l != int64(len(c.exp)) {
			t.Errorf("[%d] - expected %d keys, got %d", i, len(c.exp), l)
		} else {
			for _, key := range c.exp {
				// Cannot assume an ordering
				found := false
				for k := int64(0); k < ob.Len().Int(); k++ {
					if ob.Get(Number(k)) == key {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("[%d] - expected key %v to exist, got nil", i, key)
				}
			}
		}
	}
}

func TestLen(t *testing.T) {
	cases := []struct {
		src Val
		exp int64
	}{
		0: {
			src: Nil,
			exp: 0,
		},
		1: {
			src: Number(3.14),
			exp: 4,
		},
		2: {
			src: String("hi, there"),
			exp: 9,
		},
		3: {
			src: Bool(true),
			exp: 4,
		},
		4: {
			src: String(`this
has
new
lines`),
			exp: 18,
		},
		5: {
			src: &object{
				map[Val]Val{
					Number(1):      String("val1"),
					String("name"): Bool(false),
					String("subobj"): &object{
						map[Val]Val{
							String("key"): Number(10),
						},
					},
				},
			},
			exp: 3,
		},
		6: {
			src: &object{},
			exp: 0,
		},
		7: {
			src: String(""),
			exp: 0,
		},
	}

	bi := new(builtinMod)
	ctx := NewCtx(nil, nil)
	bi.SetCtx(ctx)
	for i, c := range cases {
		ret := Get1(bi._len(c.src))
		if c.exp != ret.Int() {
			t.Errorf("[%d] - expected %d, got %d", i, c.exp, ret.Int())
		}
	}
}

func TestPanic(t *testing.T) {
	ctx := NewCtx(nil, nil)

	cases := []struct {
		src Val
		err bool
	}{
		0: {
			src: Nil,
			err: false,
		},
		1: {
			src: Bool(false),
			err: false,
		},
		2: {
			src: String(""),
			err: false,
		},
		3: {
			src: Number(0),
			err: false,
		},
		4: {
			src: Number(0.0),
			err: false,
		},
		5: {
			src: &object{
				map[Val]Val{
					String("__bool"): NewNativeFunc(ctx, "", func(args ...Val) []Val {
						return Set1(Bool(false))
					}),
				},
			},
			err: false,
		},
		6: {
			src: Number(0.1),
			err: true,
		},
		7: {
			src: Bool(true),
			err: true,
		},
		8: {
			src: String("error"),
			err: true,
		},
		9: {
			src: NewNativeFunc(ctx, "", func(args ...Val) []Val { return nil }),
			err: true,
		},
		10: {
			src: Number(-1),
			err: true,
		},
		11: {
			src: &object{},
			err: true,
		},
		12: {
			src: &object{
				map[Val]Val{
					String("__bool"): NewNativeFunc(ctx, "", func(args ...Val) []Val {
						return Set1(Bool(true))
					}),
				},
			},
			err: true,
		},
	}

	bi := new(builtinMod)
	bi.SetCtx(ctx)
	for i, c := range cases {
		func() {
			defer func() {
				if e := recover(); (e != nil) != c.err {
					if c.err {
						t.Errorf("[%d] - expected a panic, got none", i)
					} else {
						t.Errorf("[%d] - expected no panic, got %v", i, e)
					}
				}
			}()
			bi._panic(c.src)
		}()
	}
}

func TestRecover(t *testing.T) {
	ctx := NewCtx(nil, nil)

	cases := []struct {
		panicWith interface{}
		exp       Val
	}{
		0: {
			panicWith: nil,
			exp:       Nil,
		},
		1: {
			panicWith: Number(1),
			exp:       Number(1),
		},
		2: {
			panicWith: io.EOF,
			exp:       String("EOF"),
		},
		3: {
			panicWith: Bool(true),
			exp:       Bool(true),
		},
		4: {
			panicWith: String("test"),
			exp:       String("test"),
		},
		5: {
			panicWith: "not an error interface",
			exp:       String("not an error interface"),
		},
		6: {
			panicWith: 666,
			exp:       String("666"),
		},
		7: {
			panicWith: false,
			exp:       String("false"),
		},
	}

	bi := new(builtinMod)
	bi.SetCtx(ctx)
	for i, c := range cases {
		f := NewNativeFunc(ctx, "", func(args ...Val) []Val {
			if c.panicWith != nil {
				panic(c.panicWith)
			}
			return Set1(Nil)
		})
		ret := Get1(bi._recover(f))
		if c.exp != ret {
			t.Errorf("[%d] - expected %v, got %v", i, c.exp, ret)
		}
	}
}

func TestConvBool(t *testing.T) {
	ctx := NewCtx(nil, nil)
	// For case 9 below
	ob := NewObject()
	ob.Set(String("__bool"), NewNativeFunc(ctx, "", func(args ...Val) []Val {
		return Set1(Bool(false))
	}))

	cases := []struct {
		src Val
		exp Val
		err bool
	}{
		0: {
			src: Nil,
			exp: Bool(false),
		},
		1: {
			src: Number(1),
			exp: Bool(true),
		},
		2: {
			src: Number(3.1415),
			exp: Bool(true),
		},
		3: {
			src: Number(0),
			exp: Bool(false),
		},
		4: {
			src: Bool(true),
			exp: Bool(true),
		},
		5: {
			src: Bool(false),
			exp: Bool(false),
		},
		6: {
			src: String("some string"),
			exp: Bool(true),
		},
		7: {
			src: NewObject(),
			exp: Bool(true),
		},
		8: {
			src: NewNativeFunc(ctx, "", func(args ...Val) []Val { return nil }),
			exp: Bool(true),
		},
		9: {
			src: ob,
			exp: Bool(false),
		},
	}

	bm := new(builtinMod)
	bm.SetCtx(ctx)
	for i, c := range cases {
		func() {
			defer func() {
				if e := recover(); (e != nil) != c.err {
					if c.err {
						t.Errorf("[%d] - expected a panic, got none", i)
					} else {
						t.Errorf("[%d] - expected no panic, got %v", i, e)
					}
				}
			}()
			ret := Get1(bm._bool(c.src))
			if ret != c.exp {
				t.Errorf("[%d] - expected %v, got %v", i, c.exp, ret)
			}
		}()
	}
}

func TestConvString(t *testing.T) {
	ctx := NewCtx(nil, nil)
	// For case 8 below
	ob := NewObject()
	ob.Set(String("__string"), NewNativeFunc(ctx, "", func(args ...Val) []Val {
		return Set1(String("ok"))
	}))
	// For case 7 below
	fn := NewNativeFunc(ctx, "nm", func(args ...Val) []Val { return nil })

	cases := []struct {
		src   Val
		exp   Val
		start bool
	}{
		0: {
			src: Nil,
			exp: String("nil"),
		},
		1: {
			src: Number(1),
			exp: String("1"),
		},
		2: {
			src: Number(3.1415),
			exp: String("3.1415"),
		},
		3: {
			src: Bool(true),
			exp: String("true"),
		},
		4: {
			src: Bool(false),
			exp: String("false"),
		},
		5: {
			src: String("some string"),
			exp: String("some string"),
		},
		6: {
			src: NewObject(),
			exp: String("{}"),
		},
		7: {
			src:   fn,
			exp:   String("<func nm ("),
			start: true,
		},
		8: {
			src: ob,
			exp: String("ok"),
		},
	}

	bm := new(builtinMod)
	bm.SetCtx(ctx)
	for i, c := range cases {
		func() {
			defer func() {
				if e := recover(); e != nil {
					t.Errorf("[%d] - expected no panic, got %v", i, e)
				}
			}()
			ret := Get1(bm._string(c.src))
			if (c.start && !strings.HasPrefix(ret.String(), c.exp.String())) || (!c.start && ret != c.exp) {
				t.Errorf("[%d] - expected %v, got %v", i, c.exp, ret)
			}
		}()
	}
}

func TestConvNumber(t *testing.T) {
	ctx := NewCtx(nil, nil)
	// For case 10 below
	ob := NewObject()
	ob.Set(String("__float"), NewNativeFunc(ctx, "", func(args ...Val) []Val {
		return Set1(Number(22))
	}))

	cases := []struct {
		src Val
		exp Val
		err bool
	}{
		0: {
			src: Nil,
			err: true,
		},
		1: {
			src: Number(1),
			exp: Number(1),
		},
		2: {
			src: Bool(true),
			exp: Number(1),
		},
		3: {
			src: Bool(false),
			exp: Number(0),
		},
		4: {
			src: String(""),
			err: true,
		},
		5: {
			src: String("not a number"),
			err: true,
		},
		6: {
			src: String("17"),
			exp: Number(17),
		},
		7: {
			src: String("3.1415"),
			exp: Number(3.1415),
		},
		8: {
			src: NewObject(),
			err: true,
		},
		9: {
			src: NewNativeFunc(ctx, "", func(args ...Val) []Val { return nil }),
			err: true,
		},
		10: {
			src: ob,
			exp: Number(22),
		},
	}

	bm := new(builtinMod)
	bm.SetCtx(ctx)
	for i, c := range cases {
		func() {
			defer func() {
				if e := recover(); (e != nil) != c.err {
					if c.err {
						t.Errorf("[%d] - expected a panic, got none", i)
					} else {
						t.Errorf("[%d] - expected no panic, got %v", i, e)
					}
				}
			}()
			ret := Get1(bm._number(c.src))
			if ret != c.exp {
				t.Errorf("[%d] - expected %v, got %v", i, c.exp, ret)
			}
		}()
	}
}

func TestConvType(t *testing.T) {
	ctx := NewCtx(nil, nil)

	cases := []struct {
		src Val
		exp string
	}{
		0: {
			src: Nil,
			exp: "nil",
		},
		1: {
			src: Number(0),
			exp: "number",
		},
		2: {
			src: Bool(false),
			exp: "bool",
		},
		3: {
			src: String(""),
			exp: "string",
		},
		4: {
			src: NewNativeFunc(ctx, "", func(args ...Val) []Val { return nil }),
			exp: "func",
		},
		5: {
			src: NewObject(),
			exp: "object",
		},
	}
	bm := new(builtinMod)
	bm.SetCtx(ctx)
	for i, c := range cases {
		ret := Get1(bm._type(c.src))
		if ret.String() != c.exp {
			t.Errorf("[%d] - expected %s, got %s", i, c.exp, ret)
		}
	}
}
