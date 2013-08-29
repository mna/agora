package stdlib

import (
	"fmt"

	"github.com/PuerkitoBio/agora/runtime"
)

type fmtMod struct {
	ctx *runtime.Ctx
	ob  *runtime.Object
}

func NewFmt(c *runtime.Ctx) runtime.Module {
	f := &fmtMod{ctx: c}
	// Prepare the object
	f.ob = runtime.NewObject()
	f.ob.Set(runtime.String("Println"), runtime.NewNativeFunc(f.ctx, "fmt.Println", f.fmt_Println))
	f.ob.Set(runtime.String("Printf"), runtime.NewNativeFunc(f.ctx, "fmt.Printf", f.fmt_Printf))
	return f
}

func (f fmtMod) ID() string {
	return "fmt"
}

func (f fmtMod) Run() (v runtime.Val, err error) {
	return f.ob, nil
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

func (ø fmtMod) fmt_Println(args ...runtime.Val) runtime.Val {
	ifs := toNative(args)
	n, err := fmt.Fprintln(ø.ctx.Stdout, ifs...)
	if err != nil {
		panic(err)
	}
	return runtime.Int(n)
}

func (ø fmtMod) fmt_Printf(args ...runtime.Val) runtime.Val {
	var ft string
	if len(args) > 0 {
		ft = args[0].String()
	}
	ifs := toNative(args[1:])
	n, err := fmt.Fprintf(ø.ctx.Stdout, ft, ifs)
	if err != nil {
		panic(err)
	}
	return runtime.Int(n)
}
