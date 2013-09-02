// Package compiler provides the agora source code compiler.
package compiler

import (
	"io"
	"io/ioutil"

	"github.com/PuerkitoBio/agora/bytecode"
	"github.com/PuerkitoBio/agora/compiler/emitter"
	"github.com/PuerkitoBio/agora/compiler/parser"
)

// A Compiler represents the source code compiler. It implements the runtime.Compiler
// interface so that it is suitable for runtime.Ctx.
type Compiler struct{}

// Compile takes a module identifier and a reader, and compiles its source date
// to an in-memory representation of agora bytecode, ready to be executed.
// If an error is encountered, it is returned as second value, otherwise it is
// nil.
func (c *Compiler) Compile(id string, r io.Reader) (*bytecode.File, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	p := parser.New()
	syms, scps, err := p.Parse(id, b)
	if err != nil {
		return nil, err
	}
	e := new(emitter.Emitter)
	return e.Emit(id, syms, scps)
}
