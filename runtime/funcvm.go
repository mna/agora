package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/PuerkitoBio/agora/bytecode"
)

var (
	ErrValNotAnObject = errors.New("value is not an object")
)

type funcVM struct {
	proto *AgoraFunc
	pc    int
	vars  map[string]Val
	stack []Val
	sp    int
	this  Val
	args  Val
}

func newFuncVM(proto *AgoraFunc) *funcVM {
	return &funcVM{
		proto,
		0,
		make(map[string]Val, proto.expVars),
		make([]Val, 0, proto.stackSz),
		0,
		nil,
		nil,
	}
}

func (f *funcVM) push(v Val) {
	// Stack has to grow as needed, StackSz doesn't take into account the loops
	if f.sp == len(f.stack) {
		if f.proto.ctx.Debug && f.sp == cap(f.stack) {
			fmt.Fprintf(f.proto.ctx.Stdout, "DEBUG expanding stack of func %s, current size: %d\n", f.proto.name, len(f.stack))
		}
		f.stack = append(f.stack, v)
	} else {
		f.stack[f.sp] = v
	}
	f.sp++
}

func (ø *funcVM) pop() Val {
	ø.sp--
	v := ø.stack[ø.sp]
	ø.stack[ø.sp] = Nil // free this reference for gc
	return v
}

func (ø *funcVM) getVal(flg bytecode.Flag, ix uint64) Val {
	switch flg {
	case bytecode.FLG_K:
		return ø.proto.kTable[ix]
	case bytecode.FLG_V:
		// If not found, will return Nil, so the value is always fine
		v, _ := ø.proto.ctx.getVar(ø.proto.kTable[ix].String())
		return v
	case bytecode.FLG_N:
		return Nil
	case bytecode.FLG_T:
		return ø.this
	case bytecode.FLG_F:
		return ø.proto.mod.fns[ix]
	case bytecode.FLG_AA:
		return ø.args
	}
	panic(fmt.Sprintf("Func.getVal() - invalid flag value %d", flg))
}

func (ø *funcVM) dumpInstrInfo(w io.Writer, i bytecode.Instr) {
	switch i.Flag() {
	case bytecode.FLG_K:
		fmt.Fprintf(w, " ; %s", ø.proto.kTable[i.Index()].dump())
	case bytecode.FLG_V:
		fmt.Fprintf(w, " ; var %s", ø.proto.kTable[i.Index()])
	case bytecode.FLG_N:
		fmt.Fprintf(w, " ; %s", Nil.dump())
	case bytecode.FLG_T:
		fmt.Fprint(w, " ; [this]")
	case bytecode.FLG_F:
		fmt.Fprintf(w, " ; %s", ø.proto.mod.fns[i.Index()].dump())
	case bytecode.FLG_AA:
		fmt.Fprintf(w, " ; args[%d]", i.Index())
	}
}

func (ø *funcVM) dump() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "\n> %s\n", ø.proto.dump())
	// Constants
	fmt.Fprintf(buf, "  Constants:\n")
	for i, v := range ø.proto.kTable {
		fmt.Fprintf(buf, "    [%3d] %s\n", i, v.dump())
	}
	// Variables
	fmt.Fprintf(buf, "\n  Variables:\n")
	if ø.this != nil {
		fmt.Fprintf(buf, "    [this] = %s\n", ø.this.dump())
	}
	// Sort the vars for deterministic output
	sortedVars := make([]string, len(ø.vars))
	j := 0
	for k, _ := range ø.vars {
		sortedVars[j] = k
		j++
	}
	sort.Strings(sortedVars)
	for _, k := range sortedVars {
		fmt.Fprintf(buf, "    %s = %s\n", k, ø.vars[k].dump())
	}
	// Stack
	fmt.Fprintf(buf, "\n  Stack:\n")
	i := int(math.Max(0, float64(ø.sp-5)))
	for i <= ø.sp {
		if i == ø.sp {
			fmt.Fprint(buf, "sp->")
		} else {
			fmt.Fprint(buf, "    ")
		}
		v := Val(Nil)
		if i < len(ø.stack) {
			v = ø.stack[i]
		}
		fmt.Fprintf(buf, "[%3d] %s\n", i, v.dump())
		i++
	}
	// Instructions
	fmt.Fprintf(buf, "\n  Instructions:\n")
	i = int(math.Max(0, float64(ø.pc-10)))
	for i <= ø.pc+10 {
		if i == ø.pc {
			fmt.Fprintf(buf, "pc->")
		} else {
			fmt.Fprintf(buf, "    ")
		}
		if i < len(ø.proto.code) {
			fmt.Fprintf(buf, "[%3d] %s", i, ø.proto.code[i])
			ø.dumpInstrInfo(buf, ø.proto.code[i])
			fmt.Fprintln(buf)
		} else {
			break
		}
		i++
	}
	fmt.Fprintln(buf)
	return buf.String()
}

func (vm *funcVM) createArgsVal(args []Val) Val {
	if len(args) == 0 {
		return Nil
	}
	o := NewObject()
	for i, v := range args {
		o.Set(Int(i), v)
	}
	return o
}

func (ø *funcVM) run(args ...Val) Val {
	// Expected args are defined in constant table spots 0 to ExpArgs - 1.
	for j, l := int64(0), int64(len(args)); j < ø.proto.expArgs; j++ {
		if j < l {
			ø.vars[ø.proto.kTable[j].String()] = args[j]
		} else {
			ø.vars[ø.proto.kTable[j].String()] = Nil
		}
	}
	// Keep the args array
	ø.args = ø.createArgsVal(args)

	// Execute the instructions
	for {
		// Get the instruction to process
		i := ø.proto.code[ø.pc]
		// Decode the instruction
		op, flg, ix := i.Opcode(), i.Flag(), i.Index()
		// Increment the PC, if a jump requires a different PC delta, it will set it explicitly
		ø.pc++
		switch op {
		case bytecode.OP_RET:
			// End this function call, return the value on top of the stack
			return ø.pop()

		case bytecode.OP_LOAD:
			v, err := ø.proto.ctx.Load(ø.getVal(flg, ix).String())
			if err != nil {
				panic(err) // TODO : Better error management
			}
			ø.push(v)

		case bytecode.OP_PUSH:
			ø.push(ø.getVal(flg, ix))

		case bytecode.OP_POP:
			if nm, v := ø.proto.kTable[ix].String(), ø.pop(); !ø.proto.ctx.setVar(nm, v) {
				// Not found anywhere, create variable locally
				ø.vars[nm] = v
			}

		case bytecode.OP_ADD:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Add(y))

		case bytecode.OP_SUB:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Sub(y))

		case bytecode.OP_MUL:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Mul(y))

		case bytecode.OP_DIV:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Div(y))

		case bytecode.OP_MOD:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Mod(y))

		case bytecode.OP_POW:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Pow(y))

		case bytecode.OP_NOT:
			x := ø.pop()
			ø.push(ø.proto.ctx.Logic.Not(x))

		case bytecode.OP_UNM:
			x := ø.pop()
			ø.push(x.Unm())

		case bytecode.OP_CALL:
			// ix is the number of args
			// Pop the function itself, ensure it is a function
			x := ø.pop()
			f, ok := x.(Func)
			if !ok {
				// TODO : Make an ErrXxx
				panic("call on a non-function value")
			}
			// Pop the arguments in reverse order
			args := make([]Val, ix)
			for j := ix; j > 0; j-- {
				args[j-1] = ø.pop()
			}
			// Call the function, and store the return value on the stack
			ø.push(f.Call(nil, args...))

		case bytecode.OP_EQ:
			y, x := ø.pop(), ø.pop()
			cmp := x.Cmp(y)
			ø.push(Bool(cmp == 0))

		case bytecode.OP_LT:
			y, x := ø.pop(), ø.pop()
			cmp := x.Cmp(y)
			ø.push(Bool(cmp < 0))

		case bytecode.OP_LTE:
			y, x := ø.pop(), ø.pop()
			cmp := x.Cmp(y)
			ø.push(Bool(cmp <= 0))

		case bytecode.OP_GT:
			y, x := ø.pop(), ø.pop()
			cmp := x.Cmp(y)
			ø.push(Bool(cmp > 0))

		case bytecode.OP_GTE:
			y, x := ø.pop(), ø.pop()
			cmp := x.Cmp(y)
			ø.push(Bool(cmp >= 0))

		case bytecode.OP_AND:
			y, x := ø.pop(), ø.pop()
			ø.push(ø.proto.ctx.Logic.And(x, y))

		case bytecode.OP_OR:
			y, x := ø.pop(), ø.pop()
			ø.push(ø.proto.ctx.Logic.Or(x, y))

		case bytecode.OP_TEST:
			if !ø.pop().Bool() {
				// Do the jump over ix instructions
				ø.pc += int(ix)
			}

		case bytecode.OP_JMPB:
			// TODO : Eventually change to a single JMP with signed value
			ø.pc -= (int(ix) + 1) // +1 because pc is already on next instr

		case bytecode.OP_JMPF:
			ø.pc += int(ix)

		case bytecode.OP_NEW:
			ø.push(NewObject())

		case bytecode.OP_DUMP:
			if ø.proto.ctx.Debug {
				// Dumps `ix` number of stack traces
				ø.proto.ctx.dump(int(ix)) // TODO : check int value
			}

		case bytecode.OP_SFLD:
			vr, k, vl := ø.pop(), ø.pop(), ø.pop()
			if ob, ok := vr.(*Object); ok {
				ob.Set(k, vl)
				// Then push back on stack
				if flg == bytecode.FLG_Push {
					ø.push(ob)
				}
			} else {
				panic(ErrValNotAnObject)
			}

		case bytecode.OP_GFLD:
			vr, k := ø.pop(), ø.pop()
			if ob, ok := vr.(*Object); ok {
				ø.push(ob.Get(k))
			} else {
				panic(ErrValNotAnObject)
			}

		case bytecode.OP_CFLD:
			vr, k := ø.pop(), ø.pop()
			// Pop the arguments in reverse order
			args := make([]Val, ix)
			for j := ix; j > 0; j-- {
				args[j-1] = ø.pop()
			}
			if ob, ok := vr.(*Object); ok {
				ø.push(ob.callMethod(k, args...))
			} else {
				panic(ErrValNotAnObject)
			}

		default:
			panic(fmt.Sprintf("unknown opcode %s", op))
		}
	}
}
