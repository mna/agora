package stdlib

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/agora/runtime"
)

type StringsMod struct {
	ctx *runtime.Ctx
	ob  *runtime.Object
}

func (s *StringsMod) ID() string {
	return "strings"
}

func (s *StringsMod) Run(_ ...runtime.Val) (v runtime.Val, err error) {
	defer runtime.PanicToError(&err)
	if s.ob == nil {
		// Prepare the object
		s.ob = runtime.NewObject()
		s.ob.Set(runtime.String("ToLower"), runtime.NewNativeFunc(s.ctx, "strings.ToLower", s.strings_ToLower))
		s.ob.Set(runtime.String("ToUpper"), runtime.NewNativeFunc(s.ctx, "strings.ToUpper", s.strings_ToUpper))
		s.ob.Set(runtime.String("HasPrefix"), runtime.NewNativeFunc(s.ctx, "strings.HasPrefix", s.strings_HasPrefix))
		s.ob.Set(runtime.String("HasSuffix"), runtime.NewNativeFunc(s.ctx, "strings.HasSuffix", s.strings_HasSuffix))
		s.ob.Set(runtime.String("Matches"), runtime.NewNativeFunc(s.ctx, "strings.Matches", s.strings_Matches))
		s.ob.Set(runtime.String("CharAt"), runtime.NewNativeFunc(s.ctx, "strings.CharAt", s.strings_CharAt))
	}
	return s.ob, nil
}

func (s *StringsMod) SetCtx(c *runtime.Ctx) {
	s.ctx = c
}

// Converts strings to uppercase, concatenating all strings.
func (s *StringsMod) strings_ToUpper(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	buf := bytes.NewBuffer(nil)
	for _, v := range args {
		_, err := buf.WriteString(strings.ToUpper(v.String()))
		if err != nil {
			panic(err)
		}
	}
	return runtime.String(buf.String())
}

// Converts strings to lowercase, concatenating all strings.
func (s *StringsMod) strings_ToLower(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	buf := bytes.NewBuffer(nil)
	for _, v := range args {
		_, err := buf.WriteString(strings.ToLower(v.String()))
		if err != nil {
			panic(err)
		}
	}
	return runtime.String(buf.String())
}

// Returns true if the string at arg0 starts with any of the following strings.
func (s *StringsMod) strings_HasPrefix(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	src := args[0].String()
	for _, v := range args[1:] {
		if strings.HasPrefix(src, v.String()) {
			return runtime.Bool(true)
		}
	}
	return runtime.Bool(false)
}

// Returns true if the string at arg0 ends with any of the following strings.
func (s *StringsMod) strings_HasSuffix(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	src := args[0].String()
	for _, v := range args[1:] {
		if strings.HasSuffix(src, v.String()) {
			return runtime.Bool(true)
		}
	}
	return runtime.Bool(false)
}

// Args:
// 0 - The string
// 1 - The regexp pattern
// 2 - (optional) a maximum number of matches to return
//
// Returns:
// An object holding all the matches, or nil if no match.
// Each match contains:
// n - The nth match group (when n=0, the full text of the match)
// Each match group contains:
// start - the index of the start of the match
// length - the length of the match
// text - the string of the match
func (s *StringsMod) strings_Matches(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	src := args[0].String()
	rx := regexp.MustCompile(args[1].String())
	n := -1 // By default, return all matches
	if len(args) > 2 {
		n = args[2].Int()
	}
	strmtch := rx.FindAllStringSubmatch(src, n)
	if strmtch == nil {
		return runtime.Nil
	}
	ixmtch := rx.FindAllStringSubmatchIndex(src, n)
	ob := runtime.NewObject()
	for i, mtches := range strmtch {
		obch := runtime.NewObject()
		for j, mtch := range mtches {
			leaf := runtime.NewObject()
			leaf.Set(runtime.String("text"), runtime.String(mtch))
			leaf.Set(runtime.String("start"), runtime.Int(ixmtch[i][2*j]))
			leaf.Set(runtime.String("length"), runtime.Int(ixmtch[i][2*j+1]))
			obch.Set(runtime.Int(j), leaf)
		}
		ob.Set(runtime.Int(i), obch)
	}
	return ob
}

// Args:
// 0 - The source string
// 1 - The 0-based index number
//
// Returns:
// The character at that position, as a string, or an empty string if
// the index is out of bounds.
func (s *StringsMod) strings_CharAt(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	src := args[0].String()
	at := args[1].Int()
	if at >= len(src) {
		return runtime.String("")
	}
	return runtime.String(src[at : at+1])
}
