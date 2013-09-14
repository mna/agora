package stdlib

import (
	"testing"

	"github.com/PuerkitoBio/agora/runtime"
)

func TestStringsMatches(t *testing.T) {
	cases := []struct {
		args []runtime.Val
		exp  [][]string
	}{
		0: {
			args: []runtime.Val{
				runtime.String("this is a string"),
				runtime.String(`^.+$`),
			},
			exp: [][]string{
				0: []string{
					0: "this is a string",
				},
			},
		},
		1: {
			args: []runtime.Val{
				runtime.String("this is a string"),
				runtime.String(".*?(is)"),
			},
			exp: [][]string{
				0: []string{
					0: "this",
					1: "is",
				},
				1: []string{
					0: " is",
					1: "is",
				},
			},
		},
		2: {
			args: []runtime.Val{
				runtime.String("what whatever who where"),
				runtime.String(`(w.)\w+`),
				runtime.Number(2),
			},
			exp: [][]string{
				0: []string{
					0: "what",
					1: "wh",
				},
				1: []string{
					0: "whatever",
					1: "wh",
				},
			},
		},
	}
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	for i, c := range cases {
		ret := sm.strings_Matches(c.args...)
		ob := ret.(runtime.Object)
		if int64(len(c.exp)) != ob.Len().Int() {
			t.Errorf("[%f] - expected %d matches, got %d", i, len(c.exp), ob.Len().Int())
		} else {
			for j := int64(0); j < ob.Len().Int(); j++ {
				// For each match, there's 0..n number of matches (0 is the full match)
				mtch := ob.Get(runtime.Number(j))
				mo := mtch.(runtime.Object)
				if int64(len(c.exp[j])) != mo.Len().Int() {
					t.Errorf("[%d] - expected %d groups in match %d, got %d", i, len(c.exp[j]), j, mo.Len().Int())
				} else {
					for k := int64(0); k < mo.Len().Int(); k++ {
						grp := mo.Get(runtime.Number(k))
						gro := grp.(runtime.Object)
						st := gro.Get(runtime.String("Start"))
						e := gro.Get(runtime.String("End"))
						if e.Int() != st.Int()+int64(len(c.exp[j][k])) {
							t.Errorf("[%d] - expected end %d for group %d of match %d, got %d", i, st.Int()+int64(len(c.exp[j][k])), k, j, e.Int())
						}
						s := gro.Get(runtime.String("Text"))
						if s.String() != c.exp[j][k] {
							t.Errorf("[%d] - expected text '%s' for group %d of match %d, got '%s'", i, c.exp[j][k], k, j, s)
						}
					}
				}
			}
		}
	}
}

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

func TestStringsReplace(t *testing.T) {
	cases := []struct {
		args []runtime.Val
		exp  string
	}{
		0: {
			args: []runtime.Val{
				runtime.String("this is the source"),
				runtime.String("th"),
			},
			exp: "is is e source",
		},
		1: {
			args: []runtime.Val{
				runtime.String("this is the source"),
				runtime.String("th"),
				runtime.Number(1),
			},
			exp: "is is the source",
		},
		2: {
			args: []runtime.Val{
				runtime.String("this is the source"),
				runtime.String("t"),
				runtime.String("T"),
			},
			exp: "This is The source",
		},
		3: {
			args: []runtime.Val{
				runtime.String("this is the source"),
				runtime.String("t"),
				runtime.String("T"),
				runtime.Number(1),
			},
			exp: "This is the source",
		},
	}
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	for i, c := range cases {
		ret := sm.strings_Replace(c.args...)
		if ret.String() != c.exp {
			t.Errorf("[%d] - expected %s, got %s", i, c.exp, ret)
		}
	}
}

func TestStringsTrim(t *testing.T) {
	cases := []struct {
		args []runtime.Val
		exp  string
	}{
		0: {
			args: []runtime.Val{
				runtime.String(" "),
			},
			exp: "",
		},
		1: {
			args: []runtime.Val{
				runtime.String("\n  \t   hi \r"),
			},
			exp: "hi",
		},
		2: {
			args: []runtime.Val{
				runtime.String("xoxolovexox"),
				runtime.String("xo"),
			},
			exp: "love",
		},
	}
	ctx := runtime.NewCtx(nil, nil)
	sm := new(StringsMod)
	sm.SetCtx(ctx)
	for i, c := range cases {
		ret := sm.strings_Trim(c.args...)
		if ret.String() != c.exp {
			t.Errorf("[%d] - expected %s, got %s", i, c.exp, ret)
		}
	}
}
