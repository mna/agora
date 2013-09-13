package runtime

import (
	"fmt"
)

type builtinMod struct {
	ctx *Ctx
	ob  Object
}

func (b *builtinMod) ID() string {
	return "<builtin>"
}

func (b *builtinMod) Run(_ ...Val) (v Val, err error) {
	defer PanicToError(&err)
	if b.ob == nil {
		b.ob = NewObject()
		b.ob.Set(String("import"), NewNativeFunc(b.ctx, "import", b._import))
		b.ob.Set(String("panic"), NewNativeFunc(b.ctx, "panic", b._panic))
		b.ob.Set(String("recover"), NewNativeFunc(b.ctx, "recover", b._recover))
		b.ob.Set(String("len"), NewNativeFunc(b.ctx, "len", b._len))
		b.ob.Set(String("keys"), NewNativeFunc(b.ctx, "keys", b._keys))
	}
	return b.ob, nil
}

func (b *builtinMod) SetCtx(c *Ctx) {
	b.ctx = c
}

func (b *builtinMod) _import(args ...Val) Val {
	ExpectAtLeastNArgs(1, args)
	m, err := b.ctx.Load(args[0].String())
	if err != nil {
		panic(err)
	}
	v, err := m.Run()
	if err != nil {
		panic(err)
	}
	return v
}

func (b *builtinMod) _panic(args ...Val) Val {
	ExpectAtLeastNArgs(1, args)
	if args[0].Bool() {
		panic(args[0])
	}
	return Nil
}

func (b *builtinMod) _recover(args ...Val) (ret Val) {
	// Do not catch panics if args are invalid
	ExpectAtLeastNArgs(1, args)
	f, ok := args[0].(Func)
	if !ok {
		panic("first parameter must be a function")
	}
	// Catch panics in running the function. Cannot use PanicToError, because
	// it needs the true type of the panic'd value.
	ret = Nil
	defer func() {
		if err := recover(); err != nil {
			switch v := err.(type) {
			case Val:
				ret = v
			case error:
				ret = String(v.Error())
			default:
				ret = String(fmt.Sprintf("%v", v))
			}
		}
	}()
	// Return value is discarded, because recover returns the error, if any, or Nil.
	// The function to run in recovery mode must be a closure or assign its return
	// value to an outer-scope variable.

	// TODO : This would lose the `this` keyword in case of recover being called
	// on an object's method.
	f.Call(Nil, args[1:]...)
	return ret
}

func (b *builtinMod) _len(args ...Val) Val {
	ExpectAtLeastNArgs(1, args)
	switch v := args[0].(type) {
	case Object:
		return v.Len()
	case null:
		return Number(0)
	default:
		return Number(len(v.String()))
	}
}

func (b *builtinMod) _keys(args ...Val) Val {
	ExpectAtLeastNArgs(1, args)
	ob := args[0].(Object)
	return ob.Keys()
}
