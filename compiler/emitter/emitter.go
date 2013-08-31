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
		"&&": bytecode.OP_AND,
		"||": bytecode.OP_OR,
	}
	binAsgSym2op = map[string]bytecode.Opcode{
		"+=": bytecode.OP_ADD,
		"-=": bytecode.OP_SUB,
		"*=": bytecode.OP_MUL,
		"/=": bytecode.OP_DIV,
		"%=": bytecode.OP_MOD,
	}
	unrSym2op = map[string]bytecode.Opcode{
		"++": bytecode.OP_ADD,
		"--": bytecode.OP_SUB,
		"!":  bytecode.OP_NOT,
		"-":  bytecode.OP_UNM,
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
	e.emitBlock(f, fn, syms)
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
	e.emitBlock(f, fn, stmts)
}

func (e *Emitter) emitAny(f *bytecode.File, fn *bytecode.Fn, sym *parser.Symbol, any interface{}) {
	switch v := any.(type) {
	case *parser.Symbol:
		e.emitSymbol(f, fn, v, false)
	case []*parser.Symbol:
		e.emitBlock(f, fn, v)
	default:
		e.assert(false, errors.New("expected branch of `"+sym.Id+"` to be a symbol or a slice of symbols"))
	}
}

func (e *Emitter) emitBlock(f *bytecode.File, fn *bytecode.Fn, syms []*parser.Symbol) {
	for _, sym := range syms {
		e.emitSymbol(f, fn, sym, false)
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
	case "(name)", "import", "panic", "recover", "len": // TODO : Cleaner way to handle all builtins?
		// TODO : For expected vars, the correct scope is required
		// Register the symbol
		e.assert(sym.Ar == parser.ArName || sym.Ar == parser.ArLiteral, errors.New("expected `"+sym.Id+"` to have name or literal arity"))
		kix := e.registerK(fn, sym.Val, true)
		if asg {
			e.addInstr(fn, bytecode.OP_POP, bytecode.FLG_V, kix)
		} else if sym.Ar == parser.ArLiteral {
			e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_K, kix)
		} else {
			e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_V, kix)
		}
	case "(literal)", "true", "false":
		// Register the symbol
		e.assert(!asg, errors.New("invalid assignment to a literal"))
		e.assert(sym.Ar == parser.ArLiteral, errors.New("expected `"+sym.Id+"` to have literal arity"))
		kix := e.registerK(fn, sym.Val, false)
		e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_K, kix)
	case "this":
		e.assert(!asg, errors.New("invalid assignment to the `this` keyword"))
		e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_T, 0)
	case "args":
		e.assert(!asg, errors.New("invalid assignment to the `args` keyword"))
		e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_A, 0)
	case ".", "[":
		e.assert(sym.Ar == parser.ArBinary, errors.New("expected `"+sym.Id+"` to have binary arity"))
		e.emitSymbol(f, fn, sym.Second.(*parser.Symbol), false)
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		if asg {
			e.addInstr(fn, bytecode.OP_SFLD, bytecode.FLG__, 0)
		} else {
			e.addInstr(fn, bytecode.OP_GFLD, bytecode.FLG__, 0)
		}
	case ":=":
		e.assert(sym.Ar == parser.ArBinary, errors.New("expected `:=` to have binary arity"))
		e.emitSymbol(f, fn, sym.Second.(*parser.Symbol), false)
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), true)
	case "!":
		e.assert(sym.Ar == parser.ArUnary, errors.New("expected `!` to have unary arity"))
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		e.addInstr(fn, unrSym2op[sym.Id], bytecode.FLG__, 0)
	case "-":
		if sym.Ar == parser.ArUnary {
			e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
			e.addInstr(fn, unrSym2op[sym.Id], bytecode.FLG__, 0)
			break
		}
		fallthrough
	case "+", "*", "/", "%", "<", ">", "<=", ">=", "==":
		e.assert(sym.Ar == parser.ArBinary, errors.New("expected `"+sym.Id+"` to have binary arity"))
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		e.emitSymbol(f, fn, sym.Second.(*parser.Symbol), false)
		e.addInstr(fn, binSym2op[sym.Id], bytecode.FLG__, 0)
	case "&&", "||":
		e.assert(sym.Ar == parser.ArBinary, errors.New("expected `"+sym.Id+"` to have binary arity"))
		e.emitAny(f, fn, sym, sym.First)
		e.emitAny(f, fn, sym, sym.Second)
		e.addInstr(fn, binSym2op[sym.Id], bytecode.FLG__, 0)
	case "=":
		e.assert(sym.Ar == parser.ArBinary, errors.New("expected `+` to have binary arity"))
		e.emitSymbol(f, fn, sym.Second.(*parser.Symbol), false)
		left := sym.First.(*parser.Symbol)
		if left.Id == "." {
			// Emit left, which will generate a SFLD
			e.emitSymbol(f, fn, left, true)
		} else {
			// Emit a standard POP instruction
			e.emitSymbol(f, fn, left, true)
		}
	case "+=", "-=", "*=", "/=", "%=":
		e.assert(sym.Ar == parser.ArBinary, errors.New("expected `"+sym.Id+"` to have binary arity"))
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		e.emitSymbol(f, fn, sym.Second.(*parser.Symbol), false)
		e.addInstr(fn, binAsgSym2op[sym.Id], bytecode.FLG__, 0)
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), true)
	case "++", "--":
		e.assert(sym.Ar == parser.ArUnary, errors.New("expected `"+sym.Id+"` to have unary arity"))
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		// Implicit `1` constant
		ix := e.registerK(fn, "1", false)
		e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_K, ix)
		e.addInstr(fn, unrSym2op[sym.Id], bytecode.FLG__, 0)
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), true)
	case "func":
		if sym.Name != "" {
			// Function defined as a statement, register the name as a K,
			// and push the function's value into this variable.
			kix := e.registerK(fn, sym.Name, true)
			e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_F, uint64(len(f.Fns))) // New Fn will be added at this index
			e.addInstr(fn, bytecode.OP_POP, bytecode.FLG_V, kix)
		}
		e.emitFn(f, sym)
		if sym.Name == "" {
			fix := len(f.Fns) - 1 // TODO : Will not work if there's a func within the func
			// Func defined as an expression, must be pushed on the stack
			e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_F, uint64(fix))
		}
	case "(":
		e.assert(sym.Ar == parser.ArBinary || sym.Ar == parser.ArTernary, errors.New("expected `(` to have binary or ternary arity"))
		// Push parameters
		var parms []*parser.Symbol
		var op bytecode.Opcode
		if sym.Ar == parser.ArBinary {
			parms = sym.Second.([]*parser.Symbol)
			op = bytecode.OP_CALL
		} else {
			parms = sym.Third.([]*parser.Symbol)
			op = bytecode.OP_CFLD
		}
		for _, parm := range parms {
			e.emitSymbol(f, fn, parm, false)
		}
		// If ternary, push field (Second)
		if sym.Ar == parser.ArTernary {
			e.emitSymbol(f, fn, sym.Second.(*parser.Symbol), false)
		}
		// Push function name (or parent object of the field if ternary)
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		// Call
		e.addInstr(fn, op, bytecode.FLG_An, uint64(len(parms)))
	case "{":
		e.assert(sym.Ar == parser.ArUnary, errors.New("expected `{` to have unary arity"))
		e.addInstr(fn, bytecode.OP_NEW, bytecode.FLG__, 0)
		// TODO : Pass number of fields to set to operator NEW
	case "?":
		// Similar to if, but yields a value
		e.assert(sym.Ar == parser.ArTernary, errors.New("expected `?` to have ternary arity"))
		// First is the condition, always a *Symbol
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		// Next comes the TEST
		tstIx := e.addTempInstr(fn)
		// Then the true expression, always a *Symbol
		e.emitSymbol(f, fn, sym.Second.(*parser.Symbol), false)
		// The a jump over the false expression
		jmpIx := e.addTempInstr(fn)
		// Update the test instruction, here starts the false part
		e.updateTestInstr(fn, tstIx)
		// Emit the false expression, always a *Symbol
		e.emitSymbol(f, fn, sym.Third.(*parser.Symbol), false)
		// Update the jump instruction, to after the false part
		e.updateJumpfInstr(fn, jmpIx)
	case "if":
		e.assert(sym.Ar == parser.ArStatement, errors.New("expected `if` to have statement arity"))
		// First is the condition, always a *Symbol
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		// Next comes the TEST, but we don't know yet how many instructions to jump
		// insert a placeholder (invalid op) so that it fails explicitly should it ever make it to
		// the VM.
		tstIx := e.addTempInstr(fn)
		// Then comes the body
		e.emitBlock(f, fn, sym.Second.([]*parser.Symbol))
		// Update the test instruction, now that we know where to jump to
		e.updateTestInstr(fn, tstIx)
		// Then comes the ELSE/ELSE IF, maybe
		if sym.Third != nil {
			// If so, insert a jump over the else part
			jmpIx := e.addTempInstr(fn)
			// And re-update the test instruction, since an instr was added
			e.updateTestInstr(fn, tstIx)
			// Emit the else or else-if part
			e.emitAny(f, fn, sym, sym.Third)
			// Update the jump instruction now that we know how many instrs to jump over
			e.updateJumpfInstr(fn, jmpIx)
		}
	case "for":
		// TODO : Supports only the `for [condition]` at the moment (`while` equivalent)
		// Emit the condition
		start := len(fn.Is)
		e.emitAny(f, fn, sym, sym.First)
		// Add a test instruction placeholder
		tstIx := e.addTempInstr(fn)
		// Emit the body
		e.emitAny(f, fn, sym, sym.Second)
		// Add the jump-back to for condition instruction
		e.addInstr(fn, bytecode.OP_JMP, bytecode.FLG_Jb, uint64(len(fn.Is)-start))
		// Update the test instruction
		e.updateTestInstr(fn, tstIx)
	case "debug":
		var err error
		var ix int64 = 1 // Default to 1 stack to dump
		if sym.First != nil {
			// If present, it must be a literal number
			ix, err = strconv.ParseInt(sym.First.(*parser.Symbol).Val.(string), 10, 64)
			e.assert(err == nil, errors.New("invalid number literal"))
		}
		e.addInstr(fn, bytecode.OP_DUMP, bytecode.FLG_Sn, uint64(ix))
	case "return":
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), false)
		e.addInstr(fn, bytecode.OP_RET, bytecode.FLG__, 0)
	default:
		e.err = errors.New("unexpected symbol id: " + sym.Id)
	}
	// After treating the symbol, if it had a Key value, push the Key name
	if sym.Key != nil {
		// Can be on name, literal, func call, any operator, hard to assert...
		kix := e.registerK(fn, sym.Key, true)
		e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_K, kix)
	}
}

func (e *Emitter) isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	switch s := v.(type) {
	case *parser.Symbol:
		return s == nil
	case []*parser.Symbol:
		return len(s) == 0
	}
	return true
}

func (e *Emitter) addTempInstr(fn *bytecode.Fn) int {
	e.addInstr(fn, bytecode.OP_INVL, bytecode.FLG_INVL, 0)
	return len(fn.Is) - 1
}

func (e *Emitter) updateTestInstr(fn *bytecode.Fn, ix int) {
	fn.Is[ix] = bytecode.NewInstr(bytecode.OP_TEST, bytecode.FLG_Jf, uint64(len(fn.Is)-ix-1))
}

func (e *Emitter) updateJumpfInstr(fn *bytecode.Fn, ix int) {
	fn.Is[ix] = bytecode.NewInstr(bytecode.OP_JMP, bytecode.FLG_Jf, uint64(len(fn.Is)-ix-1))
}

func (e *Emitter) addInstr(fn *bytecode.Fn, op bytecode.Opcode, flg bytecode.Flag, ix uint64) {
	if e.err != nil {
		return
	}
	switch op {
	case bytecode.OP_PUSH, bytecode.OP_LOAD, bytecode.OP_NEW:
		e.stackSz[fn] += 1
	case bytecode.OP_POP, bytecode.OP_RET, bytecode.OP_UNM, bytecode.OP_NOT, bytecode.OP_TEST,
		bytecode.OP_LT, bytecode.OP_LTE, bytecode.OP_GT, bytecode.OP_GTE, bytecode.OP_EQ,
		bytecode.OP_AND, bytecode.OP_OR, bytecode.OP_ADD, bytecode.OP_SUB, bytecode.OP_MUL,
		bytecode.OP_DIV, bytecode.OP_MOD, bytecode.OP_GFLD:
		e.stackSz[fn] -= 1
	case bytecode.OP_SFLD:
		e.stackSz[fn] -= 3
	case bytecode.OP_CALL:
		e.stackSz[fn] -= (int64(ix) + 1)
	case bytecode.OP_CFLD:
		e.stackSz[fn] -= (int64(ix) + 2)
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
			val = int64(1)
		} else {
			s = "false"
			val = int64(0)
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
