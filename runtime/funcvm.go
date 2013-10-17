package runtime

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/PuerkitoBio/agora/bytecode"
	"github.com/PuerkitoBio/gocoro"
)

// An agoraFuncVM is a runnable instance of a function value. It holds the virtual machine
// required to execute the instructions.
type agoraFuncVM struct {
	// Func info
	val   *agoraFuncVal
	proto *agoraFuncDef
	debug bool

	// Stacks and counters
	pc    int   // program counter
	stack []Val // function stack
	sp    int
	rng   rangeStack

	// Variables
	vars map[string]Val
	this Val
	args Val
}

// Instantiate a runnable representation of the function prototype.
func newFuncVM(fv *agoraFuncVal) *agoraFuncVM {
	p := fv.proto
	return &agoraFuncVM{
		val:   fv,
		proto: p,
		debug: p.ctx.Debug,
		stack: make([]Val, 0, p.stackSz),
		vars:  make(map[string]Val, len(p.lTable)),
	}
}

// Push a value onto the stack.
func (f *agoraFuncVM) push(v Val) {
	// Stack has to grow as needed, StackSz doesn't take into account the loops
	if f.sp == len(f.stack) {
		if f.debug && f.sp == cap(f.stack) {
			fmt.Fprintf(f.proto.ctx.Stdout, "DEBUG expanding stack of func %s, current size: %d\n", f.val.name, len(f.stack))
		}
		f.stack = append(f.stack, v)
	} else {
		f.stack[f.sp] = v
	}
	f.sp++
}

// Pop a value from the stack.
func (f *agoraFuncVM) pop() Val {
	f.sp--
	v := f.stack[f.sp]
	f.stack[f.sp] = Nil // free this reference for gc
	return v
}

// Get a value from *somewhere*, depending on the flag.
func (f *agoraFuncVM) getVal(flg bytecode.Flag, ix uint64) Val {
	switch flg {
	case bytecode.FLG_K:
		return f.proto.kTable[ix]
	case bytecode.FLG_V:
		// Fail if variable cannot be found
		varNm := f.proto.kTable[ix].String()
		v, ok := f.proto.ctx.getVar(varNm, f)
		if !ok {
			panic("variable not found: " + varNm)
		}
		return v
	case bytecode.FLG_N:
		return Nil
	case bytecode.FLG_T:
		return f.this
	case bytecode.FLG_F:
		return newAgoraFuncVal(f.proto.mod.fns[ix], f)
	case bytecode.FLG_A:
		return f.args
	}
	panic(fmt.Sprintf("invalid flag value %d", flg))
}

// Pretty-print an instruction.
func (f *agoraFuncVM) dumpInstrInfo(w io.Writer, i bytecode.Instr) {
	switch i.Flag() {
	case bytecode.FLG_K:
		v := f.proto.kTable[i.Index()]
		fmt.Fprintf(w, " ; %s", dumpVal(v))
	case bytecode.FLG_V:
		fmt.Fprintf(w, " ; var %s", f.proto.kTable[i.Index()])
	case bytecode.FLG_N:
		fmt.Fprintf(w, " ; %s", Nil.Dump())
	case bytecode.FLG_T:
		fmt.Fprint(w, " ; [this]")
	case bytecode.FLG_F:
		fmt.Fprintf(w, " ; [func %s]", f.proto.mod.fns[i.Index()].name)
	case bytecode.FLG_A:
		fmt.Fprint(w, " ; [args]")
	}
}

// Pretty-print a function's execution context.
func (f *agoraFuncVM) dump() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "\n> %s\n", f.val.Dump())
	// Constants
	fmt.Fprintf(buf, "  Constants:\n")
	for i, v := range f.proto.kTable {
		fmt.Fprintf(buf, "    [%3d] %s\n", i, dumpVal(v))
	}
	// Variables
	fmt.Fprintf(buf, "\n  Variables:\n")
	if f.this != nil {
		fmt.Fprintf(buf, "    [this] = %s\n", dumpVal(f.this))
	}
	if f.args != nil {
		fmt.Fprintf(buf, "    [args] = %s\n", dumpVal(f.args))
	}
	// Sort the vars for deterministic output
	sortedVars := make([]string, len(f.vars))
	j := 0
	for k, _ := range f.vars {
		sortedVars[j] = k
		j++
	}
	sort.Strings(sortedVars)
	for _, k := range sortedVars {
		fmt.Fprintf(buf, "    %s = %s\n", k, dumpVal(f.vars[k]))
	}
	// Stack
	fmt.Fprintf(buf, "\n  Stack:\n")
	i := int(math.Max(0, float64(f.sp-5)))
	for i <= f.sp {
		if i == f.sp {
			fmt.Fprint(buf, "sp->")
		} else {
			fmt.Fprint(buf, "    ")
		}
		v := Val(Nil)
		if i < len(f.stack) {
			v = f.stack[i]
		}
		fmt.Fprintf(buf, "[%3d] %s\n", i, dumpVal(v))
		i++
	}
	// Instructions
	fmt.Fprintf(buf, "\n  Instructions:\n")
	i = int(math.Max(0, float64(f.pc-10)))
	for i <= f.pc+10 {
		if i == f.pc {
			fmt.Fprintf(buf, "pc->")
		} else {
			fmt.Fprintf(buf, "    ")
		}
		if i < len(f.proto.code) {
			fmt.Fprintf(buf, "[%3d] %s", i, f.proto.code[i])
			f.dumpInstrInfo(buf, f.proto.code[i])
			fmt.Fprintln(buf)
		} else {
			break
		}
		i++
	}
	fmt.Fprintln(buf)
	return buf.String()
}

// Create the reserved identifier `args` value, as an Object.
func (vm *agoraFuncVM) createArgsVal(args []Val) Val {
	if len(args) == 0 {
		return Nil
	}
	o := NewObject()
	for i, v := range args {
		o.Set(Number(i), v)
	}
	return o
}

// Create the local variables all initialized to nil
func (vm *agoraFuncVM) createLocals() {
	for _, s := range vm.proto.lTable {
		vm.vars[s] = Nil
	}
}

// run executes the instructions of the function. This is the actual implementation
// of the Virtual Machine.
func (f *agoraFuncVM) run(args ...Val) []Val {
	// Register the defer to release all `for range` coroutines created
	// by the VM and possibly still alive from a resume of this VM.
	clearRange := true
	defer func() {
		if clearRange {
			f.rng.clear()
		}
	}()

	// Keep reference to arithmetic and comparer
	arith := f.proto.ctx.Arithmetic
	cmp := f.proto.ctx.Comparer

	// If the program counter is 0, this is an initial run, not a resume as
	// a coroutine.
	if f.pc == 0 {
		// Create local variables
		f.createLocals()

		// Expected args are defined in constant table spots 0 to ExpArgs - 1.
		for j, l := int64(0), int64(len(args)); j < f.proto.expArgs; j++ {
			if j < l {
				f.vars[f.proto.kTable[j].String()] = args[j]
			} else {
				f.vars[f.proto.kTable[j].String()] = Nil
			}
		}
		// Keep the args array
		f.args = f.createArgsVal(args)
	} else {
		// This is a resume for a coroutine, push the received arg (only one) on the stack
		var a0 Val = Nil
		if len(args) > 0 {
			a0 = args[0]
		}
		f.push(a0)
	}

	// Execute the instructions
	for {
		// Get the instruction to process
		i := f.proto.code[f.pc]
		// Decode the instruction
		op, flg, ix := i.Opcode(), i.Flag(), i.Index()
		// Increment the PC, if a jump requires a different PC delta, it will set it explicitly
		f.pc++
		switch op {
		case bytecode.OP_RET:
			// End this function call, return the value on top of the stack and remove
			// the vm if it was set on the value
			f.val.coroState = nil
			return Set1(f.pop())

		case bytecode.OP_YLD:
			// Yield n value(s), save the vm so it can be called back, and return
			f.val.coroState = f
			clearRange = false // Keep active range coros, so that they can continue on a resume
			return Set1(f.pop())

		case bytecode.OP_PUSH:
			f.push(f.getVal(flg, ix))

		case bytecode.OP_POP:
			if nm, v := f.proto.kTable[ix].String(), f.pop(); !f.proto.ctx.setVar(nm, v, f) {
				// Not found anywhere, panic
				panic("unknown variable: " + nm)
			}

		case bytecode.OP_ADD:
			y, x := f.pop(), f.pop()
			f.push(arith.Add(x, y))

		case bytecode.OP_SUB:
			y, x := f.pop(), f.pop()
			f.push(arith.Sub(x, y))

		case bytecode.OP_MUL:
			y, x := f.pop(), f.pop()
			f.push(arith.Mul(x, y))

		case bytecode.OP_DIV:
			y, x := f.pop(), f.pop()
			f.push(arith.Div(x, y))

		case bytecode.OP_MOD:
			y, x := f.pop(), f.pop()
			f.push(arith.Mod(x, y))

		case bytecode.OP_NOT:
			x := f.pop()
			f.push(Bool(!x.Bool()))

		case bytecode.OP_UNM:
			x := f.pop()
			f.push(arith.Unm(x))

		case bytecode.OP_EQ:
			y, x := f.pop(), f.pop()
			f.push(Bool(cmp.Cmp(x, y) == 0))

		case bytecode.OP_NEQ:
			y, x := f.pop(), f.pop()
			f.push(Bool(cmp.Cmp(x, y) != 0))

		case bytecode.OP_LT:
			y, x := f.pop(), f.pop()
			f.push(Bool(cmp.Cmp(x, y) < 0))

		case bytecode.OP_LTE:
			y, x := f.pop(), f.pop()
			f.push(Bool(cmp.Cmp(x, y) <= 0))

		case bytecode.OP_GT:
			y, x := f.pop(), f.pop()
			f.push(Bool(cmp.Cmp(x, y) > 0))

		case bytecode.OP_GTE:
			y, x := f.pop(), f.pop()
			f.push(Bool(cmp.Cmp(x, y) >= 0))

		case bytecode.OP_TEST:
			if !f.pop().Bool() {
				// Do the jump over ix instructions
				f.pc += int(ix)
			}

		case bytecode.OP_JMP:
			if flg == bytecode.FLG_Jf {
				f.pc += int(ix)
			} else {
				f.pc -= (int(ix) + 1) // +1 because pc is already on next instr
			}

		case bytecode.OP_NEW:
			ob := NewObject()
			for j := ix; j > 0; j-- {
				key, val := f.pop(), f.pop()
				ob.Set(key, val)
			}
			f.push(ob)

		case bytecode.OP_SFLD:
			vr, k, vl := f.pop(), f.pop(), f.pop()
			if ob, ok := vr.(Object); ok {
				ob.Set(k, vl)
			} else {
				panic(NewTypeError(Type(vr), "", "object"))
			}

		case bytecode.OP_GFLD:
			vr, k := f.pop(), f.pop()
			if ob, ok := vr.(Object); ok {
				f.push(ob.Get(k))
			} else {
				panic(NewTypeError(Type(vr), "", "object"))
			}

		case bytecode.OP_CFLD:
			vr, k := f.pop(), f.pop()
			// Pop the arguments in reverse order
			args := make([]Val, ix)
			for j := ix; j > 0; j-- {
				args[j-1] = f.pop()
			}
			if ob, ok := vr.(Object); ok {
				vals := ob.callMethod(k, args...)
				for _, v := range vals {
					f.push(v)
				}
			} else {
				panic(NewTypeError(Type(vr), "", "object"))
			}

		case bytecode.OP_CALL:
			// ix is the number of args
			// Pop the function itself, ensure it is a function
			x := f.pop()
			fn, ok := x.(Func)
			if !ok {
				panic(NewTypeError(Type(x), "", "func"))
			}
			// Pop the arguments in reverse order
			args := make([]Val, ix)
			for j := ix; j > 0; j-- {
				args[j-1] = f.pop()
			}
			// Call the function, and store the return value(s) on the stack
			vals := fn.Call(nil, args...)
			for _, v := range vals {
				f.push(v)
			}

		case bytecode.OP_RNGS:
			// Pop the arguments in reverse order
			args := make([]Val, ix)
			for j := ix; j > 0; j-- {
				args[j-1] = f.pop()
			}
			// Create the range coroutine
			f.rng.push(args...)

		case bytecode.OP_RNGP:
			coro := f.rng.st[f.rng.sp-1]
			v, e := coro.Resume()
			var vals []interface{}
			if sl, ok := v.([]interface{}); ok {
				vals = sl
			} else {
				vals = []interface{}{v}
			}
			// Push the values
			if e == nil {
				for j := uint64(0); j < ix; j++ {
					if j < uint64(len(vals)) {
						f.push(vals[j].(Val))
					} else {
						f.push(Nil)
					}
				}
			} else if e != gocoro.ErrEndOfCoro {
				panic(e)
			}
			// Push the condition
			f.push(Bool(e == nil))

		case bytecode.OP_RNGE:
			// Release the range coroutine
			f.rng.pop()

		case bytecode.OP_DUMP:
			if f.debug {
				// Dumps `ix` number of stack traces
				f.proto.ctx.dump(int(ix))
			}

		default:
			panic(fmt.Sprintf("unknown opcode %s", op))
		}
	}
}
