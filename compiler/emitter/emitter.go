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
	e.emitFn(f, syms, scps)
	return f, nil
}

func (e *Emitter) emitFn(f *bytecode.File, syms []*parser.Symbol, scps []*parser.Scope) {
	i := 0
	fn := new(bytecode.Fn)
	if len(f.Fns) == 1 {
		fn.Header.Name = f.Name
	} else {
		e.assert(syms[i].Ar == parser.ArFunction, ErrExpectedFunc)
		fn.Header.Name = syms[i].Name
	}
}

func (e *Emitter) assert(cond bool, err error) {
	if !cond {
		e.err = err
	}
}
