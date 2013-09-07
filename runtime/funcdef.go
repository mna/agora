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

// NewNativeFunc returns a native function initialized with the specified context,
// name and function implementation.
func NewNativeFunc(ctx *Ctx, nm string, fn FuncFn) *NativeFunc {
	return &NativeFunc{
		&funcVal{
			ctx,
			nm,
		},
		fn,
	}
}

// An AgoraFunc represents an agora function.
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

// Native returns the Go native representation of an agora function.
func (a *AgoraFunc) Native() interface{} {
	return a
}

// Cmp compares an Agora function to another value.
func (a *AgoraFunc) Cmp(v Val) int {
	if a == v {
		return 0
	}
	return -1
}

// Call instantiates an executable function intance from this agora function
// prototype, sets the `this` value and executes the function's instructions.
// It returns the agora function's return value.
func (a *AgoraFunc) Call(this Val, args ...Val) Val {
	vm := newFuncVM(a)
	vm.this = this
	a.ctx.pushFn(a, vm)
	defer a.ctx.popFn()
	return vm.run(args...)
}

// A NativeFunc represents a Go function exposed to agora.
type NativeFunc struct {
	// Expose the default Func value's behaviour
	*funcVal

	// Internal fields
	fn FuncFn
}

// ExpectAtLeastNArgs is a utility function for native modules implementation
// to ensure that the minimum number of arguments required are provided. It panics
// otherwise, which is the correct way to raise errors in the agora runtime.
func ExpectAtLeastNArgs(n int, args []Val) {
	if len(args) < n {
		panic(fmt.Sprintf("expected at least %d argument(s), got %d", n, len(args)))
	}
}

// Native returns the Go native representation of the native function type.
func (n *NativeFunc) Native() interface{} {
	return n
}

// Cmp compares the native function with another value.
func (n *NativeFunc) Cmp(v Val) int {
	if n == v {
		return 0
	}
	return -1
}

// Call executes the native function and returns its return value.
func (n *NativeFunc) Call(_ Val, args ...Val) Val {
	n.ctx.pushFn(n, nil)
	defer n.ctx.popFn()
	return n.fn(args...)
}
