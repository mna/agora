package stdlib

import (
	"path/filepath"

	"github.com/PuerkitoBio/agora/runtime"
)

type FilepathMod struct {
	ctx *runtime.Ctx
	ob  runtime.Object
}

func (fp *FilepathMod) ID() string {
	return "filepath"
}

func (fp *FilepathMod) Run(_ ...runtime.Val) (v runtime.Val, err error) {
	defer runtime.PanicToError(&err)
	if fp.ob == nil {
		// Prepare the object
		fp.ob = runtime.NewObject()
		fp.ob.Set(runtime.String("Abs"), runtime.NewNativeFunc(fp.ctx, "filepath.Abs", fp.filepath_Abs))
		fp.ob.Set(runtime.String("Base"), runtime.NewNativeFunc(fp.ctx, "filepath.Base", fp.filepath_Base))
		fp.ob.Set(runtime.String("Dir"), runtime.NewNativeFunc(fp.ctx, "filepath.Dir", fp.filepath_Dir))
		fp.ob.Set(runtime.String("Ext"), runtime.NewNativeFunc(fp.ctx, "filepath.Ext", fp.filepath_Ext))
		fp.ob.Set(runtime.String("IsAbs"), runtime.NewNativeFunc(fp.ctx, "filepath.IsAbs", fp.filepath_IsAbs))
		fp.ob.Set(runtime.String("Join"), runtime.NewNativeFunc(fp.ctx, "filepath.Join", fp.filepath_Join))
	}
	return fp.ob, nil
}

func (fp *FilepathMod) SetCtx(c *runtime.Ctx) {
	fp.ctx = c
}

func (fp *FilepathMod) filepath_Abs(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	s, e := filepath.Abs(args[0].String())
	if e != nil {
		panic(e)
	}
	return runtime.String(s)
}

func (fp *FilepathMod) filepath_Base(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.String(filepath.Base(args[0].String()))
}

func (fp *FilepathMod) filepath_Dir(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.String(filepath.Dir(args[0].String()))
}

func (fp *FilepathMod) filepath_Ext(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.String(filepath.Ext(args[0].String()))
}

func (fp *FilepathMod) filepath_IsAbs(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Bool(filepath.IsAbs(args[0].String()))
}

func (fp *FilepathMod) filepath_Join(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	s := toString(args)
	return runtime.String(filepath.Join(s...))
}
