package emitter

import (
	"errors"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/agora/bytecode"
	"github.com/PuerkitoBio/agora/compiler/parser"
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
	// TODO : Line start and end
	f.Fns = append(f.Fns, fn)
	for _, sym := range syms {
		e.emitSymbol(f, fn, sym, false)
	}
	return f, e.err
}

func (e *Emitter) emitSymbol(f *bytecode.File, fn *bytecode.Fn, sym *parser.Symbol, asg bool) {
	if e.err != nil {
		return
	}
	switch sym.Id {
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
		kix := e.registerK(fn, sym.Val, false)
		e.addInstr(fn, bytecode.OP_PUSH, bytecode.FLG_K, kix)
	case ":=":
		e.assert(sym.Ar == parser.ArBinary, errors.New("expected `:=` to have binary arity"))
		e.emitSymbol(f, fn, sym.Second.(*parser.Symbol), false)
		e.emitSymbol(f, fn, sym.First.(*parser.Symbol), true)
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
	case bytecode.OP_POP, bytecode.OP_RET:
		e.stackSz[fn] -= 1
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
