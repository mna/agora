package emitter

import (
	"errors"

	"github.com/PuerkitoBio/agora/bytecode"
	"github.com/PuerkitoBio/agora/compiler/parser"
)

var (
	ErrExpectedFunc = errors.New("expected a function")
)

type Emitter struct {
	err error
}

func (e *Emitter) Emit(id string, syms []*parser.Symbol, scps []*parser.Scope) (*bytecode.File, error) {
	f := bytecode.NewFile(id)
	e.emitRoot(f, syms, scps)
	return f, nil
}

func (e *Emitter) emitRoot(f *bytecode.File, syms []*parser.Symbol, scps []*parser.Scope) {
	fn := new(bytecode.Fn)
	fn.Header.Name = f.Name
	f.Fns = append(f.Fns, fn)
	for _, sym := range syms {
		e.emitSymbol(f, fn, sym)
	}
}

func (e *Emitter) emitSymbol(f *bytecode.File, fn *bytecode.Fn, sym *parser.Symbol) {
	switch sym.Ar {
	case parser.ArBinary:
		switch sym.Id {
		case ":=":
			// Start by emitting the rvalue

		}
	case parser.ArImport:

	case parser.ArFunction:
		chfn := new(bytecode.Fn)
		f.Fns = append(f.Fns, chfn)
		e.emitFn(chfn, sym)
	}
}

func (e *Emitter) emitFn(fn *bytecode.Fn, syms *parser.Symbol) {
}

func (e *Emitter) assert(cond bool, err error) {
	if !cond {
		e.err = err
	}
}
