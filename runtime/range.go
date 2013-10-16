package runtime

import (
	"strings"

	"github.com/PuerkitoBio/gocoro"
)

type rangeStack struct {
	st []gocoro.Caller
	sp int
}

var (
	factory = map[string]func(...Val) gocoro.Caller{
		"number": newNumberCoro,
		"string": newStringCoro,
		"object": newObjectCoro,
	}
)

func newNumberCoro(args ...Val) gocoro.Caller {
	l := len(args)
	start := int64(0)
	max := args[0].Int()
	inc := int64(1)
	if l > 1 {
		start = max
		max = args[1].Int()
	}
	if l > 2 {
		inc = args[2].Int()
	}
	return gocoro.New(func(y gocoro.Yielder, args ...interface{}) interface{} {
		if inc >= 0 {
			for i := start; i < max; i += inc {
				y.Yield(Number(i))
			}
		} else {
			for i := start; i > max; i += inc {
				y.Yield(Number(i))
			}
		}
		panic(gocoro.ErrEndOfCoro)
	})
}

func newStringCoro(args ...Val) gocoro.Caller {
	l := len(args)
	src := args[0].String()
	sep := ""
	if l > 1 && args[1].Bool() {
		sep = args[1].String()
	}
	max := int64(-1)
	if l > 2 {
		max = args[2].Int()
	}
	return gocoro.New(func(y gocoro.Yielder, args ...interface{}) interface{} {
		if max == 0 {
			panic(gocoro.ErrEndOfCoro)
		}
		if sep == "" {
			cnt := int64(len(src))
			if max >= 0 && max < cnt {
				cnt = max
			}
			for i := int64(0); i < cnt; i++ {
				y.Yield(String(src[i]))
			}
		} else {
			cnt := int64(0)
			for max < 0 || cnt < max {
				splits := strings.SplitN(src, sep, 2)
				if len(splits) == 0 {
					break
				}
				y.Yield(String(splits[0]))
				cnt++
				if len(splits) == 1 {
					break
				}
				src = splits[1]
			}
		}
		panic(gocoro.ErrEndOfCoro)
	})
}

func newObjectCoro(args ...Val) gocoro.Caller {
	ob := args[0].(Object)
	return gocoro.New(func(y gocoro.Yielder, args ...interface{}) interface{} {
		ks := ob.Keys().(Object)
		for i := int64(0); i < ks.Len().Int(); i++ {
			val := NewObject()
			key := ks.Get(Number(i))
			val.Set(String("k"), key)
			val.Set(String("v"), ob.Get(key))
			y.Yield(val)
		}
		panic(gocoro.ErrEndOfCoro)
	})
}

func newFuncCoro(args ...Val) gocoro.Caller {
	fn := args[0].(Func)
	if afn, ok := fn.(*agoraFuncVal); ok {
		afn.reset()
		return gocoro.New(func(y gocoro.Yielder, _ ...interface{}) interface{} {
			for v := afn.Call(Nil, args[1:]...); afn.status() == "suspended"; v = afn.Call(Nil) {
				y.Yield(v)
			}
			panic(gocoro.ErrEndOfCoro)
		})
	} else {
		panic(NewTypeError("native func", "", "range"))
	}
}

func (rs rangeStack) push(args ...Val) {
	ExpectAtLeastNArgs(1, args)
	t := Type(args[0])
	if fn, ok := factory[t]; !ok {
		panic(NewTypeError(t, "", "range"))
	} else {
		c := fn(args...)
		if rs.sp == len(rs.st) {
			rs.st = append(rs.st, c)
		} else {
			rs.st[rs.sp] = c
		}
		rs.sp++
	}
}

func (rs rangeStack) pop() {
	rs.sp--
	c := rs.st[rs.sp]
	rs.st[rs.sp] = nil
	if c.Status() == gocoro.StSuspended {
		c.Cancel()
	}
}

func (rs rangeStack) clear() {
	for rs.sp > 0 {
		rs.pop()
	}
}
