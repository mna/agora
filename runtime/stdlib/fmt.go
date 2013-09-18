package stdlib

import (
	"bufio"
	"fmt"

	"github.com/PuerkitoBio/agora/runtime"
)

// The fmt module, as documented in
// https://github.com/PuerkitoBio/agora/wiki/Standard-library
type FmtMod struct {
	ctx *runtime.Ctx
	ob  runtime.Object
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
		f.ob.Set(runtime.String("Scanln"), runtime.NewNativeFunc(f.ctx, "fmt.Scanln", f.fmt_Scanln))
		f.ob.Set(runtime.String("Scanint"), runtime.NewNativeFunc(f.ctx, "fmt.Scanint", f.fmt_Scanint))
	}
	return f.ob, nil
}

func (f *FmtMod) SetCtx(c *runtime.Ctx) {
	f.ctx = c
}

func toStringIface(args []runtime.Val) []interface{} {
	var ifs []interface{}

	if len(args) > 0 {
		ifs = make([]interface{}, len(args))
		for i, v := range args {
			ifs[i] = v.String()
		}
	}
	return ifs
}

func (f *FmtMod) fmt_Print(args ...runtime.Val) runtime.Val {
	ifs := toStringIface(args)
	n, err := fmt.Fprint(f.ctx.Stdout, ifs...)
	if err != nil {
		panic(err)
	}
	return runtime.Number(n)
}

func (f *FmtMod) fmt_Println(args ...runtime.Val) runtime.Val {
	ifs := toStringIface(args)
	n, err := fmt.Fprintln(f.ctx.Stdout, ifs...)
	if err != nil {
		panic(err)
	}
	return runtime.Number(n)
}

func (f *FmtMod) fmt_Scanln(args ...runtime.Val) runtime.Val {
	var (
		b, l []byte
		e    error
		pre  bool
	)
	r := bufio.NewReader(f.ctx.Stdin)
	for l, pre, e = r.ReadLine(); pre && e == nil; l, pre, e = r.ReadLine() {
		b = append(b, l...)
	}
	if e != nil {
		panic(e)
	}
	b = append(b, l...)
	return runtime.String(b)
}

func (f *FmtMod) fmt_Scanint(args ...runtime.Val) runtime.Val {
	var i int
	if _, e := fmt.Fscanf(f.ctx.Stdin, "%d", &i); e != nil {
		panic(e)
	}
	return runtime.Number(i)
}
