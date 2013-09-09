package stdlib

import (
	"os"
)

import (
	"testing"

	"github.com/PuerkitoBio/agora/runtime"
)

func TestOsFields(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	om := new(OsMod)
	om.SetCtx(ctx)
	ob, err := om.Run()
	if err != nil {
		panic(err)
	}
	{
		ob := ob.(runtime.Object)
		ret := ob.Get(runtime.String("PathSeparator"))
		exp := string(os.PathSeparator)
		if ret.String() != exp {
			t.Errorf("expected path separator %s, got %s", exp, ret.String())
		}
		ret = ob.Get(runtime.String("PathListSeparator"))
		exp = string(os.PathListSeparator)
		if ret.String() != exp {
			t.Errorf("expected path list separator %s, got %s", exp, ret.String())
		}
		ret = ob.Get(runtime.String("DevNull"))
		exp = os.DevNull
		if ret.String() != exp {
			t.Errorf("expected dev/null %s, got %s", exp, ret)
		}
	}
}

func TestOsExect(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	om := new(OsMod)
	om.SetCtx(ctx)
	exp := "hello"
	ret := om.os_Exec(runtime.String("echo"), runtime.String(exp))
	// Shell adds a \n after output
	if ret.String() != exp+"\n" {
		t.Errorf("expected '%s', got '%s'", exp, ret)
	}
}
