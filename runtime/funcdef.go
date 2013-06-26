package runtime

// FuncFn represents the Func signature for native functions.
type FuncFn func(...Val) Val

// A Func value in Goblin is a Val that also implements the Func interface.
type Func interface {
	Call(...Val) Val
}

func NewNativeFunc(ctx *Ctx, nm string, fn FuncFn) *NativeFunc {
	return &NativeFunc{
		&funcVal{
			ctx,
			nm,
		},
		fn,
	}
}

func NewGoblinFunc(ctx *Ctx, nm string) *GoblinFunc {
	return &GoblinFunc{
		&funcVal{
			ctx,
			nm,
		},
		0,
		0,
		0,
		nil,
		nil,
	}
}

type GoblinFunc struct {
	// Expose the default Func value's behaviour
	*funcVal

	// Internal fields filled by the compiler
	stackSz int
	expArgs int
	expVars int
	kTable  []Val
	code    []Instr
}

func (ø *GoblinFunc) Native() interface{} {
	return ø
}

func (ø *GoblinFunc) Cmp(v Val) int {
	if ø == v {
		return 0
	}
	return -1
}

func (ø *GoblinFunc) Call(args ...Val) Val {
	ø.ctx.push(ø)
	defer ø.ctx.pop()
	vm := newFuncVM(ø)
	return vm.run(args...)
}

type NativeFunc struct {
	// Expose the default Func value's behaviour
	*funcVal

	// Internal fields
	fn FuncFn
}

func (ø *NativeFunc) Native() interface{} {
	return ø
}

func (ø *NativeFunc) Cmp(v Val) int {
	if ø == v {
		return 0
	}
	return -1
}

func (ø *NativeFunc) Call(args ...Val) Val {
	ø.ctx.push(ø)
	defer ø.ctx.pop()
	return ø.fn(args...)
}
