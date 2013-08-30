package runtime

import (
	"fmt"
	"github.com/PuerkitoBio/agora/bytecode"
)

// FuncFn represents the Func signature for native functions.
type FuncFn func(...Val) Val

// A Func value in Agora is a Val that also implements the Func interface.
type Func interface {
	Val
	Call(this Val, args ...Val) Val
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

type AgoraFunc struct {
	// Expose the default Func value's behaviour
	*funcVal

	// Internal fields filled by the compiler
	mod     *agoraModule
	stackSz int64
	expArgs int64
	expVars int64
	kTable  []Val
	code    []bytecode.Instr
}

func newAgoraFunc(mod *agoraModule, c *Ctx) *AgoraFunc {
	return &AgoraFunc{
		&funcVal{ctx: c},
		mod,
		0,
		0,
		0,
		nil,
		nil,
	}
}

func (ø *AgoraFunc) Native() interface{} {
	return ø
}

func (ø *AgoraFunc) Cmp(v Val) int {
	if ø == v {
		return 0
	}
	return -1
}

func (ø *AgoraFunc) Call(this Val, args ...Val) Val {
	vm := newFuncVM(ø)
	vm.this = this
	ø.ctx.push(ø, vm)
	defer ø.ctx.pop()
	return vm.run(args...)
}

type NativeFunc struct {
	// Expose the default Func value's behaviour
	*funcVal

	// Internal fields
	fn FuncFn
}

func ExpectAtLeastNArgs(n int, args []Val) {
	if len(args) < n {
		panic(fmt.Sprintf("expected at least %d argument(s), got %d", n, len(args)))
	}
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

func (ø *NativeFunc) Call(_ Val, args ...Val) Val {
	ø.ctx.push(ø, nil)
	defer ø.ctx.pop()
	return ø.fn(args...)
}
