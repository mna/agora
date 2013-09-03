package stdlib

import (
	"bytes"
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
	}
	return s.ob, nil
}

func (s *StringsMod) SetCtx(c *runtime.Ctx) {
	s.ctx = c
}

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
