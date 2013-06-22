package stdlib

import (
	"fmt"

	"github.com/PuerkitoBio/goblin/runtime"
)

var Fmt map[string]runtime.NativeFunc

func init() {
	Fmt = make(map[string]runtime.NativeFunc, 1)
	Fmt["fmt.Println"] = fmt_Println
	Fmt["fmt.Printf"] = fmt_Printf
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

func fmt_Println(s runtime.Streams, args ...runtime.Val) runtime.Val {
	ifs := toNative(args)
	n, err := fmt.Fprintln(s.Stdout(), ifs...)
	if err != nil {
		panic(err)
	}
	return runtime.Int(n)
}

func fmt_Printf(s runtime.Streams, args ...runtime.Val) runtime.Val {
	var ft string
	if len(args) > 0 {
		ft = args[0].String()
	}
	ifs := toNative(args[1:])
	n, err := fmt.Fprintf(s.Stdout(), ft, ifs)
	if err != nil {
		panic(err)
	}
	return runtime.Int(n)
}
