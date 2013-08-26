package compiler

import (
	"io"
	"io/ioutil"

	"github.com/PuerkitoBio/agora/bytecode"
	"github.com/PuerkitoBio/agora/compiler/emitter"
	"github.com/PuerkitoBio/agora/compiler/parser"
)

type Compiler struct{}

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
