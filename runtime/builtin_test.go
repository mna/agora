package runtime

import (
	"io"
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
				o.Set(String("__keys"), NewNativeFunc(ctx, "", func(args ...Val) Val {
					k := NewObject()
					k.Set(Number(0), String("b"))
					k.Set(Number(1), Number(1))
					return k
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
		ret := bi._keys(c.src)
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
		ret := bi._len(c.src)
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
					String("__bool"): NewNativeFunc(ctx, "", func(args ...Val) Val {
						return Bool(false)
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
			src: NewNativeFunc(ctx, "", func(args ...Val) Val { return Nil }),
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
					String("__toBool"): NewNativeFunc(ctx, "", func(args ...Val) Val {
						return Bool(true)
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
		f := NewNativeFunc(ctx, "", func(args ...Val) Val {
			if c.panicWith != nil {
				panic(c.panicWith)
			}
			return Nil
		})
		ret := bi._recover(f)
		if c.exp != ret {
			t.Errorf("[%d] - expected %v, got %v", i, c.exp, ret)
		}
	}
}
