package stdlib

import (
	"path/filepath"
	"testing"

	"github.com/PuerkitoBio/agora/runtime"
)

func TestFilepathAbs(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	fm := new(FilepathMod)
	fm.SetCtx(ctx)
	p := "./testdata"
	// Abs
	exp, e := filepath.Abs(p)
	if e != nil {
		panic(e)
	}
	ret := runtime.Get1(fm.filepath_Abs(runtime.String(p)))
	if ret.String() != exp {
		t.Errorf("expected '%s', got '%s'", exp, ret.String())
	}
	// IsAbs
	{
		exp := filepath.IsAbs(p)
		ret := runtime.Get1(fm.filepath_IsAbs(runtime.String(p)))
		if ret.Bool() != exp {
			t.Errorf("expected '%v', got '%v'", exp, ret.Bool())
		}
	}
}

func TestFilepathBaseDirExt(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	fm := new(FilepathMod)
	fm.SetCtx(ctx)
	p, e := filepath.Abs("./testdata/readfile.txt")
	if e != nil {
		panic(e)
	}
	// Base
	exp := filepath.Base(p)
	ret := runtime.Get1(fm.filepath_Base(runtime.String(p)))
	if ret.String() != exp {
		t.Errorf("expected base '%s', got '%s'", exp, ret.String())
	}
	// Dir
	exp = filepath.Dir(p)
	ret = runtime.Get1(fm.filepath_Dir(runtime.String(p)))
	if ret.String() != exp {
		t.Errorf("expected dir '%s', got '%s'", exp, ret.String())
	}
	// Ext
	exp = filepath.Ext(p)
	ret = runtime.Get1(fm.filepath_Ext(runtime.String(p)))
	if ret.String() != exp {
		t.Errorf("expected extension '%s', got '%s'", exp, ret.String())
	}
}

func TestFilepathJoin(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	fm := new(FilepathMod)
	fm.SetCtx(ctx)
	parts := []string{"./testdata", "..", "../../compiler", "test"}
	exp := filepath.Join(parts...)
	vals := make([]runtime.Val, len(parts))
	for i, s := range parts {
		vals[i] = runtime.String(s)
	}
	ret := runtime.Get1(fm.filepath_Join(vals...))
	if ret.String() != exp {
		t.Errorf("expected '%s', got '%s'", exp, ret.String())
	}
}
