package stdlib

import (
	"math"
	"testing"
	"time"

	"github.com/PuerkitoBio/agora/runtime"
)

func TestPi(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	ob, err := mm.Run()
	if err != nil {
		panic(err)
	}
	ret := ob.(runtime.Object).Get(runtime.String("Pi"))
	exp := math.Pi
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMax(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)

	cases := []struct {
		src []runtime.Val
		exp runtime.Val
	}{
		0: {
			src: []runtime.Val{runtime.Number(3), runtime.Number(0), runtime.Number(-12.74), runtime.Number(1)},
			exp: runtime.Number(3),
		},
		1: {
			src: []runtime.Val{runtime.String("24"), runtime.Bool(true), runtime.Number(12.74)},
			exp: runtime.Number(24),
		},
		2: {
			src: []runtime.Val{runtime.Number(0), runtime.String("0")},
			exp: runtime.Number(0),
		},
	}

	for i, c := range cases {
		ret := mm.math_Max(c.src...)
		if ret != c.exp {
			t.Errorf("[%d] - expected %f, got %f", i, c.exp.Float(), ret.Float())
		}
	}
}

func TestMin(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)

	cases := []struct {
		src []runtime.Val
		exp runtime.Val
	}{
		0: {
			src: []runtime.Val{runtime.Number(3), runtime.Number(0), runtime.Number(-12.74), runtime.Number(1)},
			exp: runtime.Number(-12.74),
		},
		1: {
			src: []runtime.Val{runtime.String("24"), runtime.Bool(true), runtime.Number(12.74)},
			exp: runtime.Number(1),
		},
		2: {
			src: []runtime.Val{runtime.Number(0), runtime.String("0")},
			exp: runtime.Number(0),
		},
	}

	for i, c := range cases {
		ret := mm.math_Min(c.src...)
		if ret != c.exp {
			t.Errorf("[%d] - expected %f, got %f", i, c.exp.Float(), ret.Float())
		}
	}
}

func TestRand(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)

	mm.math_RandSeed(runtime.Number(time.Now().UnixNano()))
	// no-arg form
	ret := mm.math_Rand()
	if ret.Int() < 0 {
		t.Errorf("expected no-arg to procude non-negative value, got %d", ret.Int())
	}
	// one-arg form
	ret = mm.math_Rand(runtime.Number(10))
	if ret.Int() < 0 || ret.Int() >= 10 {
		t.Errorf("expected one-arg to produce non-negative value lower than 10, got %d", ret.Int())
	}
	// two-args form
	ret = mm.math_Rand(runtime.Number(3), runtime.Number(9))
	if ret.Int() < 3 || ret.Int() >= 9 {
		t.Errorf("expected two-args to produce value >= 3 and < 9, got %d", ret.Int())
	}
}
func TestMathAbs(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := -3.5
	ret := mm.math_Abs(runtime.Number(val))
	exp := math.Abs(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathAcos(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 0.12
	ret := mm.math_Acos(runtime.Number(val))
	exp := math.Acos(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathAcosh(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 1.12
	ret := mm.math_Acosh(runtime.Number(val))
	exp := math.Acosh(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathAsin(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 0.12
	ret := mm.math_Asin(runtime.Number(val))
	exp := math.Asin(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathAsinh(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 0.12
	ret := mm.math_Asinh(runtime.Number(val))
	exp := math.Asinh(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathAtan(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 0.12
	ret := mm.math_Atan(runtime.Number(val))
	exp := math.Atan(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathAtan2(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 0.12
	val2 := 1.12
	ret := mm.math_Atan2(runtime.Number(val), runtime.Number(val2))
	exp := math.Atan2(val, val2)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathAtanh(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 0.12
	ret := mm.math_Atanh(runtime.Number(val))
	exp := math.Atanh(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathCeil(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 6.12
	ret := mm.math_Ceil(runtime.Number(val))
	exp := math.Ceil(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathCos(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 0.12
	ret := mm.math_Cos(runtime.Number(val))
	exp := math.Cos(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathCosh(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 0.12
	ret := mm.math_Cosh(runtime.Number(val))
	exp := math.Cosh(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathExp(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 0.12
	ret := mm.math_Exp(runtime.Number(val))
	exp := math.Exp(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathFloor(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 4.12
	ret := mm.math_Floor(runtime.Number(val))
	exp := math.Floor(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathInf(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 1
	ret := mm.math_Inf(runtime.Number(val))
	exp := math.Inf(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathIsInf(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 3.12
	val2 := 1
	ret := mm.math_IsInf(runtime.Number(val), runtime.Number(val2))
	exp := math.IsInf(val, val2)
	if ret.Bool() != exp {
		t.Errorf("expected %v, got %v", exp, ret.Bool())
	}
}

func TestMathIsNaN(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 0.12
	ret := mm.math_IsNaN(runtime.Number(val))
	exp := math.IsNaN(val)
	if ret.Bool() != exp {
		t.Errorf("expected %v, got %v", exp, ret.Bool())
	}
}

func TestMathNaN(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	ret := mm.math_NaN()
	exp := math.NaN()
	if math.IsNaN(ret.Float()) != math.IsNaN(exp) {
		t.Errorf("expected NaN, got %f", ret.Float())
	}
}

func TestMathPow(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 1.12
	val2 := 3.45
	ret := mm.math_Pow(runtime.Number(val), runtime.Number(val2))
	exp := math.Pow(val, val2)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathSin(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 1.12
	ret := mm.math_Sin(runtime.Number(val))
	exp := math.Sin(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathSinh(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 1.12
	ret := mm.math_Sinh(runtime.Number(val))
	exp := math.Sinh(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathSqrt(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 1.12
	ret := mm.math_Sqrt(runtime.Number(val))
	exp := math.Sqrt(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathTan(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 1.12
	ret := mm.math_Tan(runtime.Number(val))
	exp := math.Tan(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}

func TestMathTanh(t *testing.T) {
	// This is just an interface to Go's function, so just a quick simple test
	ctx := runtime.NewCtx(nil, nil)
	mm := new(MathMod)
	mm.SetCtx(ctx)
	val := 1.12
	ret := mm.math_Tanh(runtime.Number(val))
	exp := math.Tanh(val)
	if ret.Float() != exp {
		t.Errorf("expected %f, got %f", exp, ret.Float())
	}
}
