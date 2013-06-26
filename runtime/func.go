package runtime

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrNativeFuncNotFound = errors.New("native function not found")

	ErrInvalidConvFuncToInt    = errors.New("cannot convert Func to Int")
	ErrInvalidConvFuncToFloat  = errors.New("cannot convert Func to Float")
	ErrInvalidConvFuncToString = errors.New("cannot convert Func to String")

	ErrInvalidOpAddOnFunc = errors.New("cannot apply Add on a Func value")
	ErrInvalidOpSubOnFunc = errors.New("cannot apply Sub on a Func value")
	ErrInvalidOpMulOnFunc = errors.New("cannot apply Mul on a Func value")
	ErrInvalidOpDivOnFunc = errors.New("cannot apply Div on a Func value")
	ErrInvalidOpPowOnFunc = errors.New("cannot apply Pow on a Func value")
	ErrInvalidOpModOnFunc = errors.New("cannot apply Mod on a Func value")
	ErrInvalidOpUnmOnFunc = errors.New("cannot apply Unm on a Func value")
)

type debug struct {
	File      string
	LineStart int
	LineEnd   int
}

type FuncProto struct {
	IsNative bool
	Name     string
	StackSz  int
	ExpArgs  int
	ExpVars  int
	KTable   []Val
	Code     []Instr
	Dbg      debug
}

type Func struct {
	proto *FuncProto
	ctx   *Ctx // TODO : Call stack in context?
	pc    int
	vars  map[string]Val
	stack []Val
	sp    int
	this  Val
	args  []Val
}

func newFunc(ctx *Ctx, proto *FuncProto) *Func {
	return &Func{
		proto,
		ctx,
		0,
		make(map[string]Val, proto.ExpVars),
		make([]Val, 0, proto.StackSz), // Initial cap of StackSz
		0,
		nil,
		nil,
	}
}

// Int is an invalid conversion.
func (ø *FuncProto) Int() int {
	panic(ErrInvalidConvFuncToInt)
}

// Float is an invalid conversion.
func (ø *FuncProto) Float() float64 {
	panic(ErrInvalidConvFuncToFloat)
}

// String is an invalid conversion.
func (ø *FuncProto) String() string {
	panic(ErrInvalidConvFuncToString)
}

// Bool returns true.
func (ø *FuncProto) Bool() bool {
	return true
}

func (ø *FuncProto) Native() interface{} {
	return ø
}

func (ø *FuncProto) Cmp(v Val) int {
	if ø == v {
		// Point to same function
		return 0
	}
	// Otherwise, always return -1 (no rational way to compare 2 functions)
	return -1
}

// Add is an invalid operation.
func (ø *FuncProto) Add(v Val) Val {
	panic(ErrInvalidOpAddOnFunc)
}

// Sub is an invalid operation.
func (ø *FuncProto) Sub(v Val) Val {
	panic(ErrInvalidOpSubOnFunc)
}

// Mul is an invalid operation.
func (ø *FuncProto) Mul(v Val) Val {
	panic(ErrInvalidOpMulOnFunc)
}

// Div is an invalid operation.
func (ø *FuncProto) Div(v Val) Val {
	panic(ErrInvalidOpDivOnFunc)
}

// Mod is an invalid operation.
func (ø *FuncProto) Mod(v Val) Val {
	panic(ErrInvalidOpModOnFunc)
}

// Pow is an invalid operation.
func (ø *FuncProto) Pow(v Val) Val {
	panic(ErrInvalidOpPowOnFunc)
}

// Unm is an invalid operation.
func (ø *FuncProto) Unm() Val {
	panic(ErrInvalidOpUnmOnFunc)
}

func (ø *Func) push(v Val) {
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

func (ø *Func) pop() Val {
	ø.sp--
	v := ø.stack[ø.sp]
	ø.stack[ø.sp] = Nil // free this reference for gc
	return v
}

func (ø *Func) getVal(flg Flag, ix uint64) Val {
	switch flg {
	case FLG_K:
		return ø.proto.KTable[ix]
	case FLG_V:
		// If not found, will return Nil, so the value is always fine
		v, _ := ø.ctx.getVar(ø.proto.KTable[ix].String())
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

func (ø *FuncProto) dump() string {
	return fmt.Sprintf("%s (Func)", ø.Name)
}

func (ø *Func) dumpAll() {
	if ø.proto.IsNative {
		fmt.Printf("\nfunc %s (native)\n", ø.proto.Name)
	} else {
		fmt.Printf("\nfunc %s (file: %s)\n", ø.proto.Name, ø.proto.Dbg.File)
	}
	// Constants
	fmt.Printf("  Constants:\n")
	for i, v := range ø.proto.KTable {
		fmt.Printf("    [%3d] %s\n", i, v.dump())
	}
	// Variables
	fmt.Printf("\n  Variables:\n")
	for k, v := range ø.vars {
		fmt.Printf("    %s = %s\n", k, v.dump())
	}
	// Stack
	fmt.Printf("\n  Stack:\n")
	i := int(math.Max(0, float64(ø.sp-5)))
	for i <= ø.sp {
		if i == ø.sp {
			fmt.Print("sp->")
		} else {
			fmt.Print("    ")
		}
		v := Val(Nil)
		if i < len(ø.stack) {
			v = ø.stack[i]
		}
		fmt.Printf("[%3d] %s\n", i, v.dump())
		i++
	}
	// Instructions
	fmt.Printf("\n  Instructions:\n")
	i = int(math.Max(0, float64(ø.pc-3)))
	for i <= ø.pc+3 {
		if i == ø.pc {
			fmt.Printf("pc->")
		} else {
			fmt.Printf("    ")
		}
		if i < len(ø.proto.Code) {
			fmt.Printf("[%3d] %s\n", i, ø.proto.Code[i])
		} else {
			break
		}
		i++
	}
	fmt.Println()
}

func (ø *Func) Call(args ...Val) Val {
	ø.ctx.push(ø)
	defer ø.ctx.pop()

	if ø.proto.IsNative {
		f, ok := ø.ctx.nTable[ø.proto.Name]
		if !ok {
			panic(ErrNativeFuncNotFound)
		}
		return f(ø.ctx, args...)
	} else {
		return ø.callVM(args...)
	}
}

func (ø *Func) callVM(args ...Val) Val {
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
