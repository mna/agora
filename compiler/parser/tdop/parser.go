// Top-down operator precedence parser
// Totally based on http://javascript.crockford.com/tdop/tdop.html
package tdop

import (
	"fmt"

	"github.com/PuerkitoBio/goblin/compiler/scanner"
	"github.com/PuerkitoBio/goblin/compiler/token"
)

var (
	// TODO : For now, package-level scope, but should be in a Parser struct
	symtbl  map[string]*symbol // Symbol table
	curTok  *symbol            // Current token in symbol representation
	Scanner *scanner.Scanner
	curScp  *scope
)

type arity int

const (
	// Initial possible arities, until we know more about the context
	arName arity = iota
	arLiteral
	arOperator

	// Then it can be set to something more precise
	arUnary
	arBinary
	arTernary
	arStatement
	arThis
)

type scope struct {
	def    map[string]*symbol
	parent *scope
}

func itself(s *symbol) *symbol {
	return s
}

func (s *scope) define(n *symbol) *symbol {
	t, ok := s.def[n.val]
	if ok {
		if t.res {
			error("already reserved")
		} else {
			error("already defined")
		}
		panic("unreachable")
	}
	s.def[n.val] = n
	n.res = false
	n.lbp = 0
	n.scp = s
	n.nudfn = itself
	n.ledfn = nil
	n.stdfn = nil
	return n
}

// The find method is used to find the definition of a name. It starts with the
// current scope and seeks, if necessary, back through the chain of parent scopes
// and ultimately to the symbol table. It returns symbol_table["(name)"] if it
// cannot find a definition.
func (s *scope) find(id string) *symbol {
	for scp := s; scp != nil; scp = scp.parent {
		if o, ok := scp.def[id]; ok {
			return o
		}
	}
	if o, ok := symtbl[id]; ok {
		return o
	}
	return symtbl["(name)"]
}

func (s *scope) pop() {
	curScp = s.parent
}

func (s *scope) reserve(n *symbol) {
	if n.ar != arName || n.res {
		return
	}
	if t, ok := s.def[n.val]; ok {
		if t.res {
			return
		}
		if t.ar == arName {
			error("already defined")
		}
	}
	s.def[n.val] = n
	n.res = true
}

func clone(ori *symbol) *symbol {
	return &symbol{
		ori.id,
		ori.val,
		ori.lbp,
		ori.ar,
		ori.res,
		ori.scp,
		ori.nudfn,
		ori.ledfn,
		ori.stdfn,
	}
}

type symbol struct {
	id  string
	val string
	lbp int
	ar  arity
	res bool
	scp *scope

	nudfn func(*symbol) *symbol
	ledfn func(*symbol, *symbol) *symbol
	stdfn func(*symbol) *symbol
}

func makeSymbol(id string, bp int) *symbol {
	s, ok := symtbl[id]
	if ok {
		fmt.Println("SYMBOL REDEFINED: ", id)
		if bp >= s.lbp {
			s.lbp = bp
		}
	} else {
		s := &symbol{
			id:  id,
			val: id,
			lbp: bp,
		}
		symtbl[id] = s
	}
	return s
}

func advance(id string) *symbol {
	if id != "" && curTok.id != id {
		error("expected " + id)
	}
	tok, lit, _ := Scanner.Scan()
	// If the token is IDENT or any keyword, treat as "name" in Crockford's impl
	var (
		o  *symbol
		ok bool
		ar arity
	)
	if tok == token.IDENT || tok.IsKeyword() {
		o = curScp.find(lit)
		ar = arName
	} else if tok.IsOperator() {
		ar = arOperator
		o, ok = symtbl[id]
		if !ok {
			error("unknown operator " + id)
		}
	} else if tok.IsLiteral() { // Excluding IDENT, part of the first if
		ar = arLiteral
		o = symtbl["(literal)"]
	} else if tok == token.EOF {
		o = symtbl["(end)"]
		return o
	} else {
		error("unexpected token " + id)
	}
	curTok = clone(o)
	curTok.ar = ar
	curTok.val = lit
	return curTok
}

func error(msg string) {
	panic(msg)
}

func init() {
	makeSymbol(":", 0)
	makeSymbol(";", 0)
	makeSymbol(",", 0)
	makeSymbol(")", 0)
	makeSymbol("]", 0)
	makeSymbol("}", 0)
	makeSymbol("else", 0)
	makeSymbol("(end)", 0)
	makeSymbol("(name)", 0)
}
