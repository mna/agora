package runtime

import (
	"bytes"
	"fmt"
	"math"
)

type funcVM struct {
	proto *GoblinFunc
	pc    int
	vars  map[string]Val
	stack []Val
	sp    int
	this  Val
	args  []Val
}

func newFuncVM(proto *GoblinFunc) *funcVM {
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

func (ø *funcVM) push(v Val) {
	// Stack has to grow as needed, StackSz doesn't take into account the loops
	if ø.sp == len(ø.stack) {
		if ø.sp == cap(ø.stack) {
			fmt.Printf("DEBUG expanding stack of func %s, current size: %d\n", ø.proto.Name, len(ø.stack))
		}
		ø.stack = append(ø.stack, v)
	} else {
		ø.stack[ø.sp] = v
	}
	ø.sp++
}

func (ø *funcVM) pop() Val {
	ø.sp--
	v := ø.stack[ø.sp]
	ø.stack[ø.sp] = nil // free this reference for gc
	return v
}

func (ø *funcVM) getVal(flg Flag, ix uint64) Val {
	switch flg {
	case FLG_K:
		return ø.proto.kTable[ix]
	case FLG_V:
		// If not found, will return Nil, so the value is always fine
		v, _ := ø.proto.ctx.getVar(ø.proto.kTable[ix].String())
		return v
	case FLG_N:
		return Nil
	case FLG_T:
		return ø.this
	case FLG_F:
		return ø.ctx.Protos[ix]
	case FLG_AA:
		return ø.args[ix]
	}
	panic(fmt.Sprintf("Func.getVal() - invalid flag value %d", flg))
}

func (ø *funcVM) dump() string {
	buf := bytes.New(nil)
	fmt.Fprintf(buf, "\n> %s\n", ø.proto.dump())
	// Constants
	fmt.Fprintf(buf, "  Constants:\n")
	for i, v := range ø.proto.KTable {
		fmt.Fprintf(buf, "    [%3d] %s\n", i, v.dump())
	}
	// Variables
	fmt.Fprintf(buf, "\n  Variables:\n")
	for k, v := range ø.vars {
		fmt.Fprintf(buf, "    %s = %s\n", k, v.dump())
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
	i = int(math.Max(0, float64(ø.pc-3)))
	for i <= ø.pc+3 {
		if i == ø.pc {
			fmt.Fprintf(buf, "pc->")
		} else {
			fmt.Fprintf(buf, "    ")
		}
		if i < len(ø.proto.Code) {
			fmt.Fprintf(buf, "[%3d] %s\n", i, ø.proto.Code[i])
		} else {
			break
		}
		i++
	}
	fmt.Fprintln(buf)
	return buf.String()
}

func (ø *funcVM) run(args ...Val) Val {
	// Expected args are defined in constant table spots 0 to ExpArgs - 1.
	for j, l := 0, len(args); j < ø.proto.ExpArgs; j++ {
		if j < l {
			ø.vars[ø.proto.KTable[j].String()] = args[j]
		} else {
			ø.vars[ø.proto.KTable[j].String()] = Nil
		}
	}
	// Keep the args array
	ø.args = args

	// Execute the instructions
	for {
		// Get the instruction to process
		i := ø.proto.Code[ø.pc]
		// Decode the instruction
		op, flg, ix := i.Opcode(), i.Flag(), i.Index()
		// Increment the PC, if a jump requires a different PC delta, it will set it explicitly
		ø.pc++
		switch op {
		case OP_RET:
			// End this function call, return the value on top of the stack
			return ø.pop()

		case OP_PUSH:
			ø.push(ø.getVal(flg, ix))

		case OP_POP:
			if nm, v := ø.proto.KTable[ix].String(), ø.pop(); !ø.ctx.setVar(nm, v) {
				// Not found anywhere, create variable locally
				ø.vars[nm] = v
			}

		case OP_ADD:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Add(y))

		case OP_SUB:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Sub(y))

		case OP_MUL:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Mul(y))

		case OP_DIV:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Div(y))

		case OP_MOD:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Mod(y))

		case OP_POW:
			y, x := ø.pop(), ø.pop()
			ø.push(x.Pow(y))

		case OP_NOT:
			x := ø.pop()
			ø.push(ø.ctx.logic.Not(x))

		case OP_UNM:
			x := ø.pop()
			ø.push(x.Unm())

		case OP_CALL:
			// ix is the number of args
			// Pop the function itself, ensure it is a function
			x := ø.pop()
			fn := newFunc(ø.ctx, x.(*FuncProto))
			// Pop the arguments in reverse order
			args := make([]Val, ix)
			for j := ix; j > 0; j-- {
				args[j-1] = ø.pop()
			}
			// Call the function, and store the return value on the stack
			ø.push(fn.Call(args...))

		case OP_EQ:
			y, x := ø.pop(), ø.pop()
			cmp := x.Cmp(y)
			ø.push(Bool(cmp == 0))

		case OP_LT:
			y, x := ø.pop(), ø.pop()
			cmp := x.Cmp(y)
			ø.push(Bool(cmp < 0))

		case OP_LTE:
			y, x := ø.pop(), ø.pop()
			cmp := x.Cmp(y)
			ø.push(Bool(cmp <= 0))

		case OP_GT:
			y, x := ø.pop(), ø.pop()
			cmp := x.Cmp(y)
			ø.push(Bool(cmp > 0))

		case OP_GTE:
			y, x := ø.pop(), ø.pop()
			cmp := x.Cmp(y)
			ø.push(Bool(cmp >= 0))

		case OP_AND:
			y, x := ø.pop(), ø.pop()
			ø.push(ø.ctx.logic.And(x, y))

		case OP_OR:
			y, x := ø.pop(), ø.pop()
			ø.push(ø.ctx.logic.Or(x, y))

		case OP_TEST:
			if !ø.pop().Bool() {
				// Do the jump over ix instructions
				ø.pc += int(ix)
			}

		case OP_JMPB:
			// TODO : Eventually change to a single JMP with signed value
			ø.pc -= (int(ix) + 1) // +1 because pc is already on next instr

		case OP_JMPF:
			ø.pc += int(ix)

		case OP_NEW:
			ø.push(newObject(ø.ctx))

		case OP_DUMP:
			ø.dumpAll()

		case OP_SFLD:
			vr, k, vl := ø.pop(), ø.pop(), ø.pop()
			vr.(*Object).set(k.String(), vl) // TODO : Detect valid type to give good error message

		case OP_GFLD:
			vr, k := ø.pop(), ø.pop()
			ø.push(vr.(*Object).get(k.String())) // TODO : Detect valid type to give good error message

		case OP_CFLD:
			vr, k := ø.pop(), ø.pop()
			// Pop the arguments in reverse order
			args := make([]Val, ix)
			for j := ix; j > 0; j-- {
				args[j-1] = ø.pop()
			}
			ø.push(vr.(*Object).callMethod(k.String(), args...))

		default:
			panic(fmt.Sprintf("unknown opcode %s", op))
		}
	}
}
