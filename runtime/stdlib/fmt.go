package stdlib

import (
	"fmt"

	"github.com/PuerkitoBio/goblin/runtime"
)

type FmtMod struct {
	ctx *runtime.Ctx
}

func (ø FmtMod) Load(ctx *runtime.Ctx) runtime.Val {
	ø.ctx = ctx
	ob := runtime.NewObject()
	ob.Set(runtime.String("Println"), runtime.NewNativeFunc(ctx, "fmt.Println", ø.fmt_Println))
	ob.Set(runtime.String("Printf"), runtime.NewNativeFunc(ctx, "fmt.Printf", ø.fmt_Printf))
	return ob
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

func (ø FmtMod) fmt_Println(args ...runtime.Val) runtime.Val {
	ifs := toNative(args)
	n, err := fmt.Fprintln(ø.ctx.Stdout, ifs...)
	if err != nil {
		panic(err)
	}
	return runtime.Int(n)
}

func (ø FmtMod) fmt_Printf(args ...runtime.Val) runtime.Val {
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
