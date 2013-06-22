package stdlib

import (
	"fmt"

	"github.com/PuerkitoBio/goblin/runtime"
)

var Fmt map[string]runtime.NativeFunc

func init() {
	Fmt = make(map[string]runtime.NativeFunc, 1)
	Fmt["fmt.Println"] = fmt_Println
}

func fmt_Println(s runtime.Streams, args ...runtime.Val) runtime.Val {
	var ifs []interface{}

	if len(args) > 0 {
		ifs = make([]interface{}, len(args))
		for i, v := range args {
			ifs[i] = v.Native()
		}
	}
	n, err := fmt.Fprintln(s.Stdout(), ifs...)
	if err != nil {
		panic(err)
	}
	return runtime.Int(n)
}
