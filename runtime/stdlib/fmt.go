package stdlib

import (
	"fmt"
	"github.com/PuerkitoBio/goblin/runtime"
)

func fmt_Println(ctx *runtime.Ctx, args ...Val) Val {
	if len(args) > 0 {
		ifs := make([]interface{}, len(args))
		for i, v := range args {
			ifs[i] = v.Native()
		}
	}
	n, err := fmt.Fprintln(ctx.Stdout, ifs...)
	if err != nil {
		panic(err)
	}
	return runtime.Int(n)
}
