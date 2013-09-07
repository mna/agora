package stdlib

import (
	"os"
	"os/exec"

	"github.com/PuerkitoBio/agora/runtime"
)

type OsMod struct {
	ctx *runtime.Ctx
	ob  *runtime.Object
}

func (o *OsMod) ID() string {
	return "os"
}

func (o *OsMod) Run(_ ...runtime.Val) (v runtime.Val, err error) {
	defer runtime.PanicToError(&err)
	if o.ob == nil {
		// Prepare the object
		o.ob = runtime.NewObject()
		o.ob.Set(runtime.String("PathSeparator"), runtime.String(os.PathSeparator))
		o.ob.Set(runtime.String("PathListSeparator"), runtime.String(os.PathListSeparator))
		o.ob.Set(runtime.String("DevNull"), runtime.String(os.DevNull))
		o.ob.Set(runtime.String("Exit"), runtime.NewNativeFunc(o.ctx, "os.Exit", o.os_Exit))
		o.ob.Set(runtime.String("Getenv"), runtime.NewNativeFunc(o.ctx, "os.Getenv", o.os_Getenv))
		o.ob.Set(runtime.String("Getwd"), runtime.NewNativeFunc(o.ctx, "os.Getwd", o.os_Getwd))
	}
	return o.ob, nil
}

func (o *OsMod) SetCtx(ctx *runtime.Ctx) {
	o.ctx = ctx
}

func (o *OsMod) os_Exit(args ...runtime.Val) runtime.Val {
	if len(args) == 0 {
		os.Exit(0)
	}
	os.Exit(int(args[0].Int()))
	return runtime.Nil
}

func (o *OsMod) os_Getenv(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.String(os.Getenv(args[0].String()))
}

func (o *OsMod) os_Getwd(args ...runtime.Val) runtime.Val {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return runtime.String(pwd)
}

func (o *OsMod) os_Exec(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	c := exec.Command(args[0].String(), toString(args[1:])...)
	b, e := c.CombinedOutput()
	if e != nil {
		panic(e)
	}
	return runtime.String(b)
}

func (o *OsMod) os_Open(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	nm := args[0].String()
	flg := "r" // defaults to read-only
	if len(args) > 1 {
		// Second arg is the flag (r - w - a)
		flg = args[1].String()
	}
	perm := "" // defaults to none
	if len(args) > 2 {
		// Third arg is the permission (see os.FileMode abbreviations)
		perm = args[2].String()
	}
}

func toString(args []runtime.Val) []string {
	s := make([]string, len(args))
	for i, a := range args {
		s[i] = a.String()
	}
	return s
}
