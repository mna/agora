package stdlib

import (
	"math"
	"math/rand"

	"github.com/PuerkitoBio/agora/runtime"
)

type MathMod struct {
	ctx *runtime.Ctx
	ob  *runtime.Object
}

func (m *MathMod) ID() string {
	return "math"
}

func (m *MathMod) Run(_ ...runtime.Val) (v runtime.Val, err error) {
	defer runtime.PanicToError(&err)
	if m.ob == nil {
		// Prepare the object
		m.ob = runtime.NewObject()
		m.ob.Set(runtime.String("Pi"), runtime.Number(math.Pi))
		m.ob.Set(runtime.String("Abs"), runtime.NewNativeFunc(m.ctx, "math.Abs", m.math_Abs))
		m.ob.Set(runtime.String("Acos"), runtime.NewNativeFunc(m.ctx, "math.Acos", m.math_Acos))
		m.ob.Set(runtime.String("Acosh"), runtime.NewNativeFunc(m.ctx, "math.Acosh", m.math_Acosh))
		m.ob.Set(runtime.String("Asin"), runtime.NewNativeFunc(m.ctx, "math.Asin", m.math_Asin))
		m.ob.Set(runtime.String("Asinh"), runtime.NewNativeFunc(m.ctx, "math.Asinh", m.math_Asinh))
		m.ob.Set(runtime.String("Atan"), runtime.NewNativeFunc(m.ctx, "math.Atan", m.math_Atan))
		m.ob.Set(runtime.String("Atan2"), runtime.NewNativeFunc(m.ctx, "math.Atan2", m.math_Atan2))
		m.ob.Set(runtime.String("Atanh"), runtime.NewNativeFunc(m.ctx, "math.Atanh", m.math_Atanh))
		m.ob.Set(runtime.String("Ceil"), runtime.NewNativeFunc(m.ctx, "math.Ceil", m.math_Ceil))
		m.ob.Set(runtime.String("Cos"), runtime.NewNativeFunc(m.ctx, "math.Cos", m.math_Cos))
		m.ob.Set(runtime.String("Cosh"), runtime.NewNativeFunc(m.ctx, "math.Cosh", m.math_Cosh))
		m.ob.Set(runtime.String("Exp"), runtime.NewNativeFunc(m.ctx, "math.Exp", m.math_Exp))
		m.ob.Set(runtime.String("Floor"), runtime.NewNativeFunc(m.ctx, "math.Floor", m.math_Floor))
		m.ob.Set(runtime.String("Inf"), runtime.NewNativeFunc(m.ctx, "math.Inf", m.math_Inf))
		m.ob.Set(runtime.String("IsInf"), runtime.NewNativeFunc(m.ctx, "math.IsInf", m.math_IsInf))
		m.ob.Set(runtime.String("IsNaN"), runtime.NewNativeFunc(m.ctx, "math.IsNaN", m.math_IsNaN))
		m.ob.Set(runtime.String("Max"), runtime.NewNativeFunc(m.ctx, "math.Max", m.math_Max))
		m.ob.Set(runtime.String("Min"), runtime.NewNativeFunc(m.ctx, "math.Min", m.math_Min))
		m.ob.Set(runtime.String("NaN"), runtime.NewNativeFunc(m.ctx, "math.NaN", m.math_NaN))
		m.ob.Set(runtime.String("Pow"), runtime.NewNativeFunc(m.ctx, "math.Pow", m.math_Pow))
		m.ob.Set(runtime.String("Sin"), runtime.NewNativeFunc(m.ctx, "math.Sin", m.math_Sin))
		m.ob.Set(runtime.String("Sinh"), runtime.NewNativeFunc(m.ctx, "math.Sinh", m.math_Sinh))
		m.ob.Set(runtime.String("Sqrt"), runtime.NewNativeFunc(m.ctx, "math.Sqrt", m.math_Sqrt))
		m.ob.Set(runtime.String("Tan"), runtime.NewNativeFunc(m.ctx, "math.Tan", m.math_Tan))
		m.ob.Set(runtime.String("Tanh"), runtime.NewNativeFunc(m.ctx, "math.Tanh", m.math_Tanh))
		m.ob.Set(runtime.String("RandSeed"), runtime.NewNativeFunc(m.ctx, "math.RandSeed", m.math_RandSeed))
		m.ob.Set(runtime.String("Rand"), runtime.NewNativeFunc(m.ctx, "math.Rand", m.math_Rand))
	}
	return m.ob, nil
}

func (m *MathMod) SetCtx(ctx *runtime.Ctx) {
	m.ctx = ctx
}

func (m *MathMod) math_Abs(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Abs(args[0].Float()))
}

func (m *MathMod) math_Acos(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Acos(args[0].Float()))
}

func (m *MathMod) math_Acosh(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Acosh(args[0].Float()))
}

func (m *MathMod) math_Asin(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Asin(args[0].Float()))
}

func (m *MathMod) math_Asinh(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Asinh(args[0].Float()))
}

func (m *MathMod) math_Atan(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Atan(args[0].Float()))
}

func (m *MathMod) math_Atan2(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	return runtime.Number(math.Atan2(args[0].Float(), args[1].Float()))
}

func (m *MathMod) math_Atanh(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Atanh(args[0].Float()))
}

func (m *MathMod) math_Ceil(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Ceil(args[0].Float()))
}

func (m *MathMod) math_Cos(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Cos(args[0].Float()))
}

func (m *MathMod) math_Cosh(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Cosh(args[0].Float()))
}

func (m *MathMod) math_Exp(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Exp(args[0].Float()))
}

func (m *MathMod) math_Floor(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Floor(args[0].Float()))
}

func (m *MathMod) math_Inf(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Inf(int(args[0].Int())))
}

func (m *MathMod) math_IsInf(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	return runtime.Bool(math.IsInf(args[0].Float(), int(args[1].Int())))
}

func (m *MathMod) math_IsNaN(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Bool(math.IsNaN(args[0].Float()))
}

func (m *MathMod) math_Max(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	max := args[len(args)-1].Float()
	for i := len(args) - 2; i >= 0; i-- {
		max = math.Max(max, args[i].Float())
	}
	return runtime.Number(max)
}

func (m *MathMod) math_Min(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	min := args[len(args)-1].Float()
	for i := len(args) - 2; i >= 0; i-- {
		min = math.Min(min, args[i].Float())
	}
	return runtime.Number(min)
}

func (m *MathMod) math_NaN(_ ...runtime.Val) runtime.Val {
	return runtime.Number(math.NaN())
}

func (m *MathMod) math_Pow(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	return runtime.Number(math.Pow(args[0].Float(), args[1].Float()))
}

func (m *MathMod) math_Sin(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Sin(args[0].Float()))
}

func (m *MathMod) math_Sinh(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Sinh(args[0].Float()))
}

func (m *MathMod) math_Sqrt(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Sqrt(args[0].Float()))
}

func (m *MathMod) math_Tan(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Tan(args[0].Float()))
}

func (m *MathMod) math_Tanh(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(math.Tanh(args[0].Float()))
}

func (m *MathMod) math_RandSeed(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	rand.Seed(args[0].Int())
	return runtime.Nil
}

func (m *MathMod) math_Rand(args ...runtime.Val) runtime.Val {
	switch len(args) {
	case 0:
		return runtime.Number(rand.Int())
	case 1:
		return runtime.Number(rand.Intn(int(args[0].Int())))
	default:
		low := args[0].Int()
		high := args[1].Int()
		n := rand.Intn(int(high - low))
		return runtime.Number(int64(n) + low)
	}
}
