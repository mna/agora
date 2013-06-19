package runtime

import (
	"fmt"
)

type debug struct {
	Name      string
	File      string
	LineStart int
	LineEnd   int
}

type Var struct {
	debug
}

type FuncProto struct {
	StackSz int
	KTable  []Val
	VTable  []Var
	Code    []Instr
	debug
}

func NewFunc(proto *FuncProto) *Func {
	// Initialize all variables to the goblin Nil (not Go's nil interface)
	vars := make([]Val, len(proto.VTable))
	for i, _ := range vars {
		vars[i] = Nil
	}
	return &Func{
		proto,
		0,
		vars,
		make([]Val, proto.StackSz),
		0,
	}
}

type Func struct {
	*FuncProto
	pc    int
	vars  []Val
	stack []Val
	sp    int
}

func (ø *Func) push(v Val) {
	ø.stack[ø.sp] = v
	ø.sp++
}

func (ø *Func) pop() Val {
	ø.sp--
	v := ø.stack[ø.sp]
	ø.stack[ø.sp] = nil // free this reference for gc
	return v
}

func (ø *Func) getVal(flg Flag, ix uint64) Val {
	switch flg {
	case FLG_K:
		return ø.KTable[ix]
	case FLG_V:
		return ø.vars[ix]
	case FLG_N:
		return Nil
	}
	panic(fmt.Sprintf("Func.getVal() - invalid flag value %d", flg))
}

func (ø *Func) setVal(flg Flag, ix uint64, v Val) {
	switch flg {
	case FLG_V:
		ø.vars[ix] = v
	default:
		panic(fmt.Sprintf("Func.setVal() - invalid flag value %d", flg))
	}
}

func (ø *Func) Run() Val {
	for {
		// Get the instruction to process
		i := ø.Code[ø.pc]
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
			ø.setVal(flg, ix, ø.pop())

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
			ø.push(x.Not())

		case OP_UNM:
			x := ø.pop()
			ø.push(x.Unm())
		}
	}
}
