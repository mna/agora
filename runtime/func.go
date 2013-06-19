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

func (ø *Func) getVal(tbl Table, ix uint64) Val {
	switch tbl {
	case TBL_K:
		return ø.KTable[ix]
	case TBL_V:
		return ø.vars[ix]
	}
	panic(fmt.Sprintf("Func.getVal() - unknown tbl value %d", tbl))
}

func (ø *Func) setVal(tbl Table, ix uint64, v Val) {
	switch tbl {
	case TBL_K:
		panic("Func.setVal() - invalid set value on KTable")
	case TBL_V:
		ø.vars[ix] = v
	default:
		panic(fmt.Sprintf("Func.setVal() - unknown tbl value %d", tbl))
	}
}

func (ø *Func) Run() Val {
	for {
		// Get the instruction to process
		i := ø.Code[ø.pc]
		// Decode the instruction
		op, tbl, ix := i.Opcode(), i.Table(), i.Index()
		switch op {
		case OP_RET:
			// End this function call, return the value on top of the stack
			return ø.pop()

		case OP_PUSH:
			ø.push(ø.getVal(tbl, ix))
			ø.pc++

		case OP_POP:
			ø.setVal(tbl, ix, ø.pop())
			ø.pc++
		}
	}
}
