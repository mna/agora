package runtime

type builtinMod struct {
	ctx *Ctx
	ob  *Object
}

func (b *builtinMod) ID() string {
	return "<builtin>"
}

func (b *builtinMod) Run() (v Val, err error) {
	defer PanicToError(&err)
	if b.ob == nil {
		b.ob = NewObject()
		b.ob.Set(String("import"), NewNativeFunc(b.ctx, "import", b._import))
	}
	return b.ob, nil
}

func (b *builtinMod) SetCtx(c *Ctx) {
	b.ctx = c
}

func (b *builtinMod) _import(args ...Val) Val {
	m, err := b.ctx.Load(args[0].String()) // Will panic if no parameter received
	if err != nil {
		panic(err)
	}
	v, err := m.Run()
	if err != nil {
		panic(err)
	}
	return v
}
