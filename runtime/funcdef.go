package runtime

import (
	"fmt"

	"github.com/PuerkitoBio/agora/bytecode"
)

// FuncFn represents the Func signature for native functions.
type FuncFn func(...Val) []Val

// A Func value in Agora is a Val that also implements the Func interface.
type Func interface {
	Val
	Call(this Val, args ...Val) []Val
}

// An agoraFuncDef represents an agora function's prototype.
type agoraFuncDef struct {
	ctx *Ctx
	mod *agoraModule
	// Internal fields filled by the compiler
	name    string
	stackSz int64
	expArgs int64
	kTable  []Val
	lTable  []string
	code    []bytecode.Instr
}

func newAgoraFuncDef(mod *agoraModule, c *Ctx) *agoraFuncDef {
	return &agoraFuncDef{
		ctx: c,
		mod: mod,
	}
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

// Get1 is a utility function to extract a single Val from a slice of Vals.
// Since many methods return a single value even though multiple return values
// are supported, this is a convenient way to extract the single value.
//
// As a special case, if the slice is empty, it returns nil.
func Get1(args []Val) Val {
	if len(args) > 0 {
		return args[0]
	}
	return nil
}

// Set1 is a utility function to return a single Val as a slice of Vals.
// Since many methods return a single value even though multiple return values
// are supported, this is a convenient way to wrap the single value in the format
// expected by the function signature.
//
// As a special case, if there is no value to wrap (arg is nil), it returns nil.
func Set1(arg Val) []Val {
	if arg == nil {
		return nil
	}
	return []Val{arg}
}

// Native returns the Go native representation of the native function type.
func (n *NativeFunc) Native() interface{} {
	return n
}

// Call executes the native function and returns its return value.
func (n *NativeFunc) Call(_ Val, args ...Val) []Val {
	n.ctx.pushFn(n, nil)
	defer n.ctx.popFn()
	return n.fn(args...)
}
