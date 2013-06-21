package stdlib

import (
	"fmt"

	"github.com/PuerkitoBio/goblin/runtime"
	"github.com/PuerkitoBio/goblin/runtime/nfi"
)

var Fmt map[string]nfi.NativeFunc

func init() {
	Fmt = make(map[string]nfi.NativeFunc, 1)
	Fmt["fmt.Println"] = fmt_Println
}

func fmt_Println(s nfi.Streams, args ...runtime.Val) runtime.Val {
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
