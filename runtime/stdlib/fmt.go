package stdlib

import (
	"fmt"

	"github.com/PuerkitoBio/agora/runtime"
)

type FmtMod struct {
	ctx *runtime.Ctx
	ob  *runtime.Object
}

func (f *FmtMod) ID() string {
	return "fmt"
}

func (f *FmtMod) Run(_ ...runtime.Val) (v runtime.Val, err error) {
	defer runtime.PanicToError(&err)
	if f.ob == nil {
		// Prepare the object
		f.ob = runtime.NewObject()
		f.ob.Set(runtime.String("Print"), runtime.NewNativeFunc(f.ctx, "fmt.Print", f.fmt_Print))
		f.ob.Set(runtime.String("Println"), runtime.NewNativeFunc(f.ctx, "fmt.Println", f.fmt_Println))
		f.ob.Set(runtime.String("Printf"), runtime.NewNativeFunc(f.ctx, "fmt.Printf", f.fmt_Printf))
	}
	return f.ob, nil
}

func (f *FmtMod) SetCtx(c *runtime.Ctx) {
	f.ctx = c
}

func toNative(args []runtime.Val) []interface{} {
	var ifs []interface{}

	if len(args) > 0 {
		ifs = make([]interface{}, len(args))
		for i, v := range args {
			ifs[i] = v.Native()
		}
	}
	return ifs
}

func (f *FmtMod) fmt_Print(args ...runtime.Val) runtime.Val {
	ifs := toNative(args)
	n, err := fmt.Fprint(f.ctx.Stdout, ifs...)
	if err != nil {
		panic(err)
	}
	return runtime.Float(n)
}

func (f *FmtMod) fmt_Println(args ...runtime.Val) runtime.Val {
	ifs := toNative(args)
	n, err := fmt.Fprintln(f.ctx.Stdout, ifs...)
	if err != nil {
		panic(err)
	}
	return runtime.Float(n)
}

func (f *FmtMod) fmt_Printf(args ...runtime.Val) runtime.Val {
	var ft string
	if len(args) > 0 {
		ft = args[0].String()
	}
	ifs := toNative(args[1:])
	n, err := fmt.Fprintf(f.ctx.Stdout, ft, ifs)
	if err != nil {
		panic(err)
	}
	return runtime.Float(n)
}
