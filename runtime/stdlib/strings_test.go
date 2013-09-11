package stdlib

import (
	"testing"

	"github.com/PuerkitoBio/agora/runtime"
)

func TestStringsToUpper(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	ret := sm.strings_ToUpper(runtime.String("this"), runtime.String("Is"), runtime.String("A"), runtime.String("... strInG"))
	exp := "THISISA... STRING"
	if ret.String() != exp {
		t.Errorf("expected %s, got %s", exp, ret)
	}
}

func TestStringsToLower(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	ret := sm.strings_ToLower(runtime.String("this"), runtime.String("Is"), runtime.String("A"), runtime.String("... strInG"))
	exp := "thisisa... string"
	if ret.String() != exp {
		t.Errorf("expected %s, got %s", exp, ret)
	}
}

func TestStringsHasPrefix(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	ret := sm.strings_HasPrefix(runtime.String("what prefix?"), runtime.String("no"), runtime.Nil, runtime.Number(3), runtime.String("wh"))
	if !ret.Bool() {
		t.Errorf("expected true, got false")
	}
	ret = sm.strings_HasPrefix(runtime.String("what prefix?"), runtime.String("no"), runtime.Nil, runtime.Number(3), runtime.String("hw"))
	if ret.Bool() {
		t.Errorf("expected false, got true")
	}
}

func TestStringsHasSuffix(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	ret := sm.strings_HasSuffix(runtime.String("suffix, you say"), runtime.String("ay"), runtime.Nil, runtime.Number(3), runtime.String("wh"))
	if !ret.Bool() {
		t.Errorf("expected true, got false")
	}
	ret = sm.strings_HasSuffix(runtime.String("suffix, you say"), runtime.String("no"), runtime.Nil, runtime.Number(3), runtime.String("hw"))
	if ret.Bool() {
		t.Errorf("expected false, got true")
	}
}

func TestStringsByteAt(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	src := "some string"
	ix := 0
	ret := sm.strings_ByteAt(runtime.String(src), runtime.Number(ix))
	if ret.String() != string(src[ix]) {
		t.Errorf("expected byte %s at index %d, got %s", string(src[ix]), ix, ret)
	}
	ix = 3
	ret = sm.strings_ByteAt(runtime.String(src), runtime.Number(ix))
	if ret.String() != string(src[ix]) {
		t.Errorf("expected byte %s at index %d, got %s", string(src[ix]), ix, ret)
	}
	ix = 22
	ret = sm.strings_ByteAt(runtime.String(src), runtime.Number(ix))
	if ret.String() != "" {
		t.Errorf("expected byte %s at index %d, got %s", "", ix, ret)
	}
}

func TestStringsConcat(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	ret := sm.strings_Concat(runtime.String("hello"), runtime.Number(12), runtime.Bool(true), runtime.String("end"))
	exp := "hello12trueend"
	if ret.String() != exp {
		t.Errorf("expected %s, got %s", exp, ret)
	}
}
