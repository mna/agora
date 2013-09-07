package stdlib

import (
	"github.com/PuerkitoBio/agora/runtime"
)

type ConvMod struct {
	ctx *runtime.Ctx
	ob  *runtime.Object
}

func (c *ConvMod) ID() string {
	return "conv"
}

func (c *ConvMod) Run(_ ...runtime.Val) (v runtime.Val, err error) {
	defer runtime.PanicToError(&err)
	if c.ob == nil {
		// Prepare the object
		c.ob = runtime.NewObject()
		c.ob.Set(runtime.String("Number"), runtime.NewNativeFunc(c.ctx, "conv.Number", c.conv_Number))
		c.ob.Set(runtime.String("String"), runtime.NewNativeFunc(c.ctx, "conv.String", c.conv_String))
		c.ob.Set(runtime.String("Bool"), runtime.NewNativeFunc(c.ctx, "conv.Bool", c.conv_Bool))
		c.ob.Set(runtime.String("Type"), runtime.NewNativeFunc(c.ctx, "conv.Type", c.conv_Type))
	}
	return c.ob, nil
}

func (c *ConvMod) SetCtx(ctx *runtime.Ctx) {
	c.ctx = ctx
}

func (c *ConvMod) conv_Number(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Number(args[0].Float())
}

func (c *ConvMod) conv_String(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.String(args[0].String())
}

func (c *ConvMod) conv_Bool(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Bool(args[0].Bool())
}

func (c *ConvMod) conv_Type(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	switch args[0].(type) {
	case runtime.String:
		return runtime.String("string")
	case runtime.Number:
		return runtime.String("number")
	case runtime.Bool:
		return runtime.String("bool")
	case runtime.Func:
		return runtime.String("func")
	case *runtime.Object:
		return runtime.String("object")
	default:
		return runtime.String("nil")
	}
}
