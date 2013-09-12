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

func TestStringsContains(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	ret := sm.strings_Contains(runtime.String("contains something"), runtime.String("what"), runtime.Nil, runtime.Number(3), runtime.String("some"))
	if !ret.Bool() {
		t.Errorf("expected true, got false")
	}
	ret = sm.strings_Contains(runtime.String("contains something"), runtime.String("no"), runtime.Nil, runtime.Number(3), runtime.String("hw"))
	if ret.Bool() {
		t.Errorf("expected false, got true")
	}
}

func TestStringsIndex(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	ret := sm.strings_Index(runtime.String("agora"), runtime.String("arg"), runtime.Nil, runtime.Number(3), runtime.String("go"))
	exp := 1
	if ret.Int() != int64(exp) {
		t.Errorf("expected %d, got %d", exp, ret.Int())
	}
	ret = sm.strings_Index(runtime.String("agora"), runtime.Number(2), runtime.String("arg"), runtime.Nil, runtime.Number(3), runtime.String("go"))
	exp = -1
	if ret.Int() != int64(exp) {
		t.Errorf("expected %d, got %d", exp, ret.Int())
	}
}

func TestStringsLastIndex(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	ret := sm.strings_LastIndex(runtime.String("agoragore"), runtime.String("arg"), runtime.Nil, runtime.Number(3), runtime.String("go"))
	exp := 5
	if ret.Int() != int64(exp) {
		t.Errorf("expected %d, got %d", exp, ret.Int())
	}
	ret = sm.strings_Index(runtime.String("agoragore"), runtime.Number(6), runtime.String("arg"), runtime.Nil, runtime.Number(3), runtime.String("go"))
	exp = -1
	if ret.Int() != int64(exp) {
		t.Errorf("expected %d, got %d", exp, ret.Int())
	}
}

func TestStringsSlice(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	ret := sm.strings_Slice(runtime.String("agora"), runtime.Number(2))
	exp := "ora"
	if ret.String() != exp {
		t.Errorf("expected %s, got %s", exp, ret)
	}
	ret = sm.strings_Slice(runtime.String("agora"), runtime.Number(2), runtime.Number(4))
	exp = "or"
	if ret.String() != exp {
		t.Errorf("expected %s, got %s", exp, ret)
	}
}

func TestStringsSplit(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	ret := sm.strings_Split(runtime.String("aa:bb::dd"), runtime.String(":"))
	ob := ret.(runtime.Object)
	exp := []string{"aa", "bb", "", "dd"}
	if l := ob.Len().Int(); l != int64(len(exp)) {
		t.Errorf("expected split length of %d, got %d", len(exp), l)
	}
	for i, v := range exp {
		got := ob.Get(runtime.Number(i))
		if got.String() != v {
			t.Errorf("expected split index %d to be %s, got %s", i, v, got)
		}
	}
	ret = sm.strings_Split(runtime.String("aa:bb::dd:ee:"), runtime.String(":"), runtime.Number(2))
	ob = ret.(runtime.Object)
	exp = []string{"aa", "bb::dd:ee:"}
	if l := ob.Len().Int(); l != int64(len(exp)) {
		t.Errorf("expected split length of %d, got %d", len(exp), l)
	}
	for i, v := range exp {
		got := ob.Get(runtime.Number(i))
		if got.String() != v {
			t.Errorf("expected split index %d to be %s, got %s", i, v, got)
		}
	}
}

func TestStringsJoin(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	parts := []string{"this", "is", "", "it!"}
	ob := runtime.NewObject()
	for i, v := range parts {
		ob.Set(runtime.Number(i), runtime.String(v))
	}
	ret := sm.strings_Join(ob)
	exp := "thisisit!"
	if ret.String() != exp {
		t.Errorf("expected %s, got %s", exp, ret)
	}
	ret = sm.strings_Join(ob, runtime.String("--"))
	exp = "this--is----it!"
	if ret.String() != exp {
		t.Errorf("expected %s, got %s", exp, ret)
	}
}
