package runtime

import (
	"fmt"
)

// funcVal implements most of the Val interface's methods, except
// Native() which must be done on the actual type.
type funcVal struct {
	ctx  *Ctx
	name string
}

func (f *funcVal) dump() string {
	return fmt.Sprintf("%s (Func)", f.name)
}

// Int is an invalid conversion.
func (f *funcVal) Int() int64 {
	panic(NewTypeError("func", "", "int"))
}

// Float is an invalid conversion.
func (f *funcVal) Float() float64 {
	panic(NewTypeError("func", "", "float"))
}

// String prints the function representation
func (f *funcVal) String() string {
	return fmt.Sprintf("<func %s (%p)>", f.name, f)
}

// Bool returns true.
func (f *funcVal) Bool() bool {
	return true
}

// The environment for a given func value. This is a linked list.
type env struct {
	upvals map[string]Val
	parent *env
}

// An agoraFuncVal is a func's value, capturing its environment.
type agoraFuncVal struct {
	*funcVal
	proto *agoraFuncDef
	env   *env
}

// Create a new function value from the specified function prototype,
// with the given function instance (VM) as environment.
func newAgoraFuncVal(def *agoraFuncDef, vm *funcVM) *agoraFuncVal {
	var e *env
	if vm != nil {
		e = &env{
			vm.vars,
			vm.val.env,
		}
	}
	return &agoraFuncVal{
		&funcVal{
			def.ctx,
			def.name,
		},
		def,
		e,
	}
}

// Call instantiates an executable function intance from this agora function
// value, sets the `this` value and executes the function's instructions.
// It returns the agora function's return value.
func (a *agoraFuncVal) Call(this Val, args ...Val) Val {
	vm := newFuncVM(a)
	vm.this = this
	a.proto.ctx.pushFn(a, vm)
	defer a.proto.ctx.popFn()
	return vm.run(args...)
}

// Native returns the Go native representation of an agora function.
func (a *agoraFuncVal) Native() interface{} {
	return a
}
