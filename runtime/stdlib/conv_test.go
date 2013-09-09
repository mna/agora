package stdlib

import (
	"testing"

	"github.com/PuerkitoBio/agora/runtime"
)

func TestConvNumber(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	// For case 10 below
	ob := runtime.NewObject()
	ob.Set(runtime.String("__toFloat"), runtime.NewNativeFunc(ctx, "", func(args ...runtime.Val) runtime.Val {
		return runtime.Number(22)
	}))

	cases := []struct {
		src runtime.Val
		exp runtime.Val
		err bool
	}{
		0: {
			src: runtime.Nil,
			err: true,
		},
		1: {
			src: runtime.Number(1),
			exp: runtime.Number(1),
		},
		2: {
			src: runtime.Bool(true),
			exp: runtime.Number(1),
		},
		3: {
			src: runtime.Bool(false),
			exp: runtime.Number(0),
		},
		4: {
			src: runtime.String(""),
			err: true,
		},
		5: {
			src: runtime.String("not a number"),
			err: true,
		},
		6: {
			src: runtime.String("17"),
			exp: runtime.Number(17),
		},
		7: {
			src: runtime.String("3.1415"),
			exp: runtime.Number(3.1415),
		},
		8: {
			src: runtime.NewObject(),
			err: true,
		},
		9: {
			src: runtime.NewNativeFunc(ctx, "", func(args ...runtime.Val) runtime.Val { return runtime.Nil }),
			err: true,
		},
		10: {
			src: ob,
			exp: runtime.Number(22),
		},
	}

	cm := new(ConvMod)
	cm.SetCtx(ctx)
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
			ret := cm.conv_Number(c.src)
			if ret != c.exp {
				t.Errorf("[%d] - expected %v, got %v", i, c.exp, ret)
			}
		}()
	}
}

func TestConvType(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)

	cases := []struct {
		src runtime.Val
		exp string
	}{
		0: {
			src: runtime.Nil,
			exp: "nil",
		},
		1: {
			src: runtime.Number(0),
			exp: "number",
		},
		2: {
			src: runtime.Bool(false),
			exp: "bool",
		},
		3: {
			src: runtime.String(""),
			exp: "string",
		},
		4: {
			src: runtime.NewNativeFunc(ctx, "", func(args ...runtime.Val) runtime.Val { return runtime.Nil }),
			exp: "func",
		},
		5: {
			src: runtime.NewObject(),
			exp: "object",
		},
	}
	cm := new(ConvMod)
	cm.SetCtx(ctx)
	for i, c := range cases {
		ret := cm.conv_Type(c.src)
		if ret.String() != c.exp {
			t.Errorf("[%d] - expected %s, got %s", i, c.exp, ret)
		}
	}
}
