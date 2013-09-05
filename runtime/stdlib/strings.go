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
// Args:
// 0..n - The strings to convert to upper case and concatenate
// Returns:
// The uppercase string
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
// Args:
// 0..n - The strings to convert to lower case and concatenate
// Returns:
// The lowercase string
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
// Args:
// 0 - The source string
// 1..n - The prefixes to test
// Returns:
// true if the source string starts with any of the specified prefixes
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
// Args:
// 0 - The source string
// 1..n - The suffixes to test
// Returns:
// true if the source string ends with any of the specified suffixes
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
func (s *StringsMod) strings_ByteAt(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	src := args[0].String()
	at := args[1].Int()
	if at >= len(src) {
		return runtime.String("")
	}
	return runtime.String(src[at : at+1])
}

// Args:
// 0..n - the strings to concatenate
// Returns:
// The concatenated string
func (s *StringsMod) strings_Concat(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	buf := bytes.NewBuffer(nil)
	for _, v := range args {
		_, err := buf.WriteString(v.String())
		if err != nil {
			panic(err)
		}
	}
	return runtime.String(buf.String())
}

// Args:
// 0 - the source string
// 1..n - the strings to test if they are found in the source string
// Returns:
// True if any of the strings are found in the source string, false otherwise.
func (s *StringsMod) strings_Contains(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	src := args[0].String()
	for _, v := range args[1:] {
		if strings.Contains(src, v.String()) {
			return runtime.Bool(true)
		}
	}
	return runtime.Bool(false)
}

// Args:
// 0 - The source string
// 1 - [Optional] the start index in the source string
// 2 (or 1) .. n - The substrings to search for in the source string.
// Returns:
// The index of the first found substring in source, if any is found, or -1
func (s *StringsMod) strings_Index(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	src := args[0].String()
	start := 0
	find := 1
	switch v := args[1].(type) {
	case runtime.Int, runtime.Float:
		runtime.ExpectAtLeastNArgs(3, args)
		start = v.Int()
		find = 2
	}
	src = src[start:]
	for _, v := range args[find:] {
		if ix := strings.Index(src, v.String()); ix >= 0 {
			return runtime.Int(ix)
		}
	}
	return runtime.Int(-1)
}

// Args:
// 0 - The source string
// 1 - [Optional] the start index in the source string
// 2 (or 1) .. n - The substrings to search for in the source string.
// Returns:
// The last index of the first found substring in source, if any is found, or -1
func (s *StringsMod) strings_LastIndex(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	src := args[0].String()
	start := 0
	find := 1
	switch v := args[1].(type) {
	case runtime.Int, runtime.Float:
		runtime.ExpectAtLeastNArgs(3, args)
		start = v.Int()
		find = 2
	}
	src = src[start:]
	for _, v := range args[find:] {
		if ix := strings.LastIndex(src, v.String()); ix >= 0 {
			return runtime.Int(ix)
		}
	}
	return runtime.Int(-1)
}

// Slice a string to get a substring. Basically the same as slicing in Go.
// Args:
// 0 - The source string
// 1 - The start index
// 2 [optional] - The high bound, such that the length of the resulting string is high-start
// Results:
// The sliced string.
func (s *StringsMod) strings_Slice(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	src := args[0].String()
	start := args[1].Int()
	end := len(src)
	if len(args) > 2 {
		end = args[2].Int()
	}
	return runtime.String(src[start:end])
}

// Args:
// 0 - the source string
// 1 - the separator
// 2 [optional] - the maximum number of splits, defaults to all
// Returns:
// An array-like object with splits as values and indices as keys.
func (s *StringsMod) strings_Split(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	src := args[0].String()
	sep := args[1].String()
	cnt := -1
	if len(args) > 2 {
		cnt = args[2].Int()
	}
	splits := strings.SplitN(src, sep, cnt)
	ob := runtime.NewObject()
	for i, v := range splits {
		ob.Set(runtime.Int(i), runtime.String(v))
	}
	return ob
}

// Args:
// 0 - The source object
// 1 - The separator, empty string by default
// Returns:
// The concatenated string of all the array-like indices of the source object.
func (s *StringsMod) strings_Join(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	ob := args[0].(*runtime.Object)
	sep := ""
	if len(args) > 1 {
		sep = args[1].String()
	}
	l := ob.Len().Int()
	buf := bytes.NewBuffer(nil)
	for i := 0; i < l; i++ {
		if _, err := buf.WriteString(ob.Get(runtime.Int(i)).String()); err != nil {
			panic(err)
		}
		if _, err := buf.WriteString(sep); err != nil {
			panic(err)
		}
	}
	return runtime.String(buf.String())
}

func (s *StringsMod) strings_Replace(args ...runtime.Val) runtime.Val {

}
