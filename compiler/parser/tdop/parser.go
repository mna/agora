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
	symtbl  map[string]symbol // Symbol table
	curTok  tokenSymbol       // Current token in symbol representation
	Scanner *scanner.Scanner
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

type tokenSymbol struct {
	symbol
	ar arity
}

type symbol struct {
	id  string
	val string
	lbp int
}

func (s symbol) nud() symbol {
	s.error("undefined")
	panic("unreachable")
}

func (s symbol) led(left symbol) symbol {
	s.error("missing operator")
	panic("unreachable")
}

func (s symbol) error(msg string) {
	panic(msg)
}

func makeSymbol(id string, bp int) symbol {
	s, ok := symtbl[id]
	if ok {
		fmt.Println("SYMBOL REDEFINED: ", id)
		if bp >= s.lbp {
			s.lbp = bp
		}
	} else {
		s := symbol{
			id,
			id,
			bp,
		}
		symtbl[id] = s
	}
	return s
}

func advance(id string) tokenSymbol {
	if id != "" && curTok.id != id {
		error("expected " + id)
	}
	tok, lit, pos := Scanner.Scan()
	// If the token is IDENT or any keyword, treat as "name" in Crockford's impl
	var (
		o  symbol
		ok bool
		ar arity
	)
	if tok == token.IDENT || tok.IsKeyword() {
		o = scope.find(lit)
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
	} else {
		error("unexpected token " + id)
	}
	curTok = tokenSymbol{
		o,
		ar,
	}
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
