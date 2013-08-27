package emitter

import (
	"errors"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/agora/bytecode"
	"github.com/PuerkitoBio/agora/compiler/parser"
)

var (
	binSym2op = map[string]bytecode.Opcode{
		"+":  bytecode.OP_ADD,
		"-":  bytecode.OP_SUB,
		"*":  bytecode.OP_MUL,
		"/":  bytecode.OP_DIV,
		"%":  bytecode.OP_MOD,
		"<":  bytecode.OP_LT,
		"<=": bytecode.OP_LTE,
		">":  bytecode.OP_GT,
		">=": bytecode.OP_GTE,
		"==": bytecode.OP_EQ,
	}
)

type Emitter struct {
	err     error
	kMap    map[*bytecode.Fn]map[string]int
	stackSz map[*bytecode.Fn]int64
}

func (e *Emitter) Emit(id string, syms []*parser.Symbol, scps *parser.Scope) (*bytecode.File, error) {
	// Reset the internal fields
	e.err = nil
	e.kMap = make(map[*bytecode.Fn]map[string]int)
	e.stackSz = make(map[*bytecode.Fn]int64)

	// Create the bytecode representation structure
	f := bytecode.NewFile(id)
	fn := new(bytecode.Fn)
	fn.Header.Name = f.Name // Expected args is always 0 for top-level func
	// TODO : Line start and end, ExpVars
	f.Fns = append(f.Fns, fn)
	for _, sym := range syms {
		e.emitSymbol(f, fn, sym, false)
	}
	return f, e.err
}

func (e *Emitter) emitFn(f *bytecode.File, sym *parser.Symbol) {
	if e.err != nil {
		return
	}
	e.assert(sym.Ar == parser.ArFunction, errors.New("expected `"+sym.Id+"` to have function arity"))
	fn := new(bytecode.Fn)
	fn.Header.Name = sym.Name
	args := sym.First.([]*parser.Symbol)
	fn.Header.ExpArgs = int64(len(args))
	// TODO : ExpVars, Line Start, Line End
	f.Fns = append(f.Fns, fn)
	// Define the expected args in the K table - *MUST* be defined in spots 0..ExpArgs - 1
	for _, arg := range args {
		e.assert(arg.Ar == parser.ArName, errors.New("expected argument to have name arity"))
		e.registerK(fn, arg.Val, true)
	}
	stmts := sym.Second.([]*parser.Symbol)
	for _, stmt := range stmts {
		e.emitSymbol(f, fn, stmt, false)
	}
}

func (e *Emitter) emitSymbol(f *bytecode.File, fn *bytecode.Fn, sym *parser.Symbol, asg bool) {
	if e.err != nil {
		return
	}
	switch sym.Id {
	case "nil":
		e.assert(!asg, errors.New("invalid assignment to nil"))
		e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_N, 0)
	case "(name)":
		// TODO : For expected vars, the correct scope is required
		// Register the symbol
		kix := e.registerK(fn, sym.Val, true)
		if asg {
			e.addInstr(fn, bytecode.OP_POP, bytecode.FLG_V, kix)
		} else {
			e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_V, kix)
		}
	case "(literal)":
		// Register the symbol
		e.assert(!asg, errors.New("invalid assignment to a literal"))
		kix := e.registerK(fn, sym.Val, false)
		e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_K, kix)
	case ":=":
		e.assert(sym.Ar == parser.ArBinary, errors.New("expected `:=` to have binary arity"))
		e.emitSymbol(f, fn, sym.Second.(*parser.Symbol), false)
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), true)
	case "!":
		e.assert(sym.Ar == parser.ArUnary, errors.New("expected `!` to have unary arity"))
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		e.addInstr(fn, bytecode.OP_NOT, bytecode.FLG__, 0)
	case "-":
		if sym.Ar == parser.ArUnary {
			e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
			e.addInstr(fn, bytecode.OP_UNM, bytecode.FLG__, 0)
			break
		}
		fallthrough
	case "+", "*", "/", "%", "<", ">", "<=", ">=", "==":
		e.assert(sym.Ar == parser.ArBinary, errors.New("expected `"+sym.Id+"` to have binary arity"))
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		e.emitSymbol(f, fn, sym.Second.(*parser.Symbol), false)
		e.addInstr(fn, binSym2op[sym.Id], bytecode.FLG__, 0)
	case "func":
		if sym.Name != "" {
			// Function defined as a statement, register the name as a K,
			// and push the function's value into this variable.
			kix := e.registerK(fn, sym.Name, true)
			e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_F, uint64(len(f.Fns))) // New Fn will be added at this index
			e.addInstr(fn, bytecode.OP_POP, bytecode.FLG_V, kix)
		}
		e.emitFn(f, sym)
	case "(":
		if sym.Ar == parser.ArBinary {
			// Push parameters
			parms := sym.Second.([]*parser.Symbol)
			for _, parm := range parms {
				e.emitSymbol(f, fn, parm, false)
			}
			// Push function name
			e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
			// Call
			e.addInstr(fn, bytecode.OP_CALL, bytecode.FLG_nA, uint64(len(parms)))
		} else if sym.Ar == parser.ArTernary {

		} else {
			e.assert(false, errors.New("expected `(` to have binary or ternary arity"))
		}
	case "if":
		e.assert(sym.Ar == parser.ArStatement, errors.New("expected `if` to have statement arity"))
		// First is the condition, always a *Symbol
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		// Next comes the TEST, but we don't know yet how many instructions to jump
		// insert a placeholder (invalid op) so that it fails explicitly should it ever make it to
		// the VM.
		e.addInstr(fn, bytecode.OP_INVL, bytecode.FLG_INVL, 0)
		// Then comes the body
		// TODO : Emit block ([]*Symbol)? Reused by emitRoot and emitFn?
		// TODO : Keep count of instrs in map, or simply emit, then count - index of INVL?
		// Then comes the ELSE/ELSE IF, maybe
	case "return":
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		e.addInstr(fn, bytecode.OP_RET, bytecode.FLG__, 0)
	default:
		e.err = errors.New("unexpected symbol id: " + sym.Id)
	}
}

func (e *Emitter) addInstr(fn *bytecode.Fn, op bytecode.Opcode, flg bytecode.Flag, ix uint64) {
	if e.err != nil {
		return
	}
	switch op {
	case bytecode.OP_PUSH:
		e.stackSz[fn] += 1
	case bytecode.OP_POP, bytecode.OP_RET, bytecode.OP_UNM, bytecode.OP_NOT:
		e.stackSz[fn] -= 1
	case bytecode.OP_ADD, bytecode.OP_SUB, bytecode.OP_MUL, bytecode.OP_DIV, bytecode.OP_MOD:
		e.stackSz[fn] -= 2
	}
	if e.stackSz[fn] > fn.Header.StackSz {
		fn.Header.StackSz = e.stackSz[fn]
	}
	fn.Is = append(fn.Is, bytecode.NewInstr(op, flg, ix))
}

func (e *Emitter) registerK(fn *bytecode.Fn, val interface{}, isName bool) uint64 {
	var kt bytecode.KType
	s, ok := val.(string)
	if ok {
		if isName {
			val = s
			kt = bytecode.KtString
		} else if s[0] == '"' || s[0] == '`' {
			// Strip the quotes
			s = s[1 : len(s)-1]
			val = s
			kt = bytecode.KtString
		} else if strings.Index(s, ".") >= 0 {
			val, e.err = strconv.ParseFloat(s, 64)
			kt = bytecode.KtFloat
		} else {
			val, e.err = strconv.ParseInt(s, 10, 64)
			kt = bytecode.KtInteger
		}
	} else {
		kt = bytecode.KtBoolean
		if v := val.(bool); v {
			s = "true"
			val = 1
		} else {
			s = "false"
			val = 0
		}
	}
	m, ok := e.kMap[fn]
	if !ok {
		m = make(map[string]int)
		e.kMap[fn] = m
	}
	i, ok := m[s]
	if !ok {
		i = len(m)
		m[s] = i
		fn.Ks = append(fn.Ks, &bytecode.K{Type: kt, Val: val})
	}
	return uint64(i)
}

func (e *Emitter) assert(cond bool, err error) {
	if !cond {
		e.err = err
	}
}
