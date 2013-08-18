// Top-down operator precedence parser
// Totally based on http://javascript.crockford.com/tdop/tdop.html
package parser

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goblin/compiler/scanner"
	"github.com/PuerkitoBio/goblin/compiler/token"
)

const (
	_SYM_END  = "(end)"
	_SYM_NAME = "(name)"
	_SYM_LIT  = "(literal)"
	_SYM_ANY  = ""
)

type Parser struct {
	// Created with the Parser
	scn *scanner.Scanner // the Scanner

	// Parse state reinitialized at each .Parse() call
	tkn *Symbol            // current token in Symbol representation
	tbl map[string]*Symbol // Symbol table
	scp *Scope             // the top-level (universe) scope
	err *scanner.ErrorList // the error handler

	// Exported fields
	Debug bool
}

// Create a new Parser
func New() *Parser {
	return &Parser{
		scn: new(scanner.Scanner),
	}
}

// Parse the provided source code and returns the AST along with the
// various scopes and an error (corresponding to the scanner.ErrorList)
func (p *Parser) Parse(filename string, src []byte) ([]*Symbol, *Scope, error) {
	// Initialize parsing state
	p.tbl = make(map[string]*Symbol)
	p.err = new(scanner.ErrorList)
	u := p.newScope()
	p.defineRequiredSymbols()
	p.defineGrammar()

	// Initialize the scanner
	p.scn.Init(filename, src, p.err.Add)

	// advance() automatically sets the current token (p.tkn)
	p.advance(_SYM_ANY)

	// Parse all statements
	s := p.statements()
	// Consume the final token
	p.advance(_SYM_END)
	// Pop the universe scope
	p.popScope()

	if p.Debug {
		for _, v := range s {
			fmt.Println(v)
		}
	}
	return s, u, p.err.Err()
}

// Those tokens *must* exist in the symbol table, they are required
// regardless of the grammar of the language:
// - (name)
// - (literal)
// - (end)
func (p *Parser) defineRequiredSymbols() {
	p.makeSymbol(_SYM_END, 0)
	p.makeSymbol(_SYM_NAME, 0)
	p.makeSymbol(_SYM_LIT, 0).nudfn = itself
}

// Create a new scope, as a child of the current scope of the parser.
func (p *Parser) newScope() *Scope {
	p.scp = &Scope{
		make(map[string]*Symbol),
		p.scp,
		p,
	}
	return p.scp
}

// Exit the current scope, making its parent the new current scope.
func (p *Parser) popScope() *Scope {
	p.scp = p.scp.parent
	return p.scp
}

// Create a symbol in the symbol table.
func (p *Parser) makeSymbol(id string, bp int) *Symbol {
	s, ok := p.tbl[id]
	if ok {
		if bp >= s.lbp {
			s.lbp = bp
		}
	} else {
		s = &Symbol{
			p:   p,
			id:  id,
			val: id,
			lbp: bp,
		}
		p.tbl[id] = s
	}
	return s
}

func (p *Parser) advance(id string) *Symbol {
	if id != _SYM_ANY && p.tkn.id != id {
		p.error(p.tkn, "expected "+id)
	}
	var (
		tok token.Token
		lit string
		pos token.Position
	)
scan:
	for tok, lit, pos = p.scn.Scan(); tok == token.ILLEGAL || tok == token.COMMENT; tok, lit, pos = p.scn.Scan() {
		// Skip Illegal and Comment tokens
	}
	if p.Debug {
		fmt.Println("SCAN: ", tok, lit, pos)
	}
	// If the token is IDENT or any keyword, treat as "name" in Crockford's impl
	var (
		o  *Symbol
		ok bool
		ar arity
	)
	if tok == token.IDENT || tok.IsKeyword() {
		o = p.scp.find(lit)
		ar = arName
	} else if tok.IsOperator() {
		ar = arOperator
		o, ok = p.tbl[tok.String()]
		if !ok {
			p.err.Add(pos, "unknown operator "+tok.String())
			goto scan
		}
	} else if tok.IsLiteral() { // Excluding IDENT, part of the first if
		ar = arLiteral
		o = p.tbl[_SYM_LIT]
	} else if tok == token.EOF {
		o = p.tbl[_SYM_END]
		o.tok = token.EOF
		o.pos = pos
		p.tkn = o
		return o
	} else {
		p.err.Add(pos, "unexpected token "+tok.String())
		goto scan
	}
	p.tkn = o.clone()
	p.tkn.ar = ar
	p.tkn.val = lit
	p.tkn.tok = tok
	p.tkn.pos = pos
	return p.tkn
}

func (p *Parser) expression(rbp int) *Symbol {
	t := p.tkn
	p.advance(_SYM_ANY)
	// Special case if in the process of defining a new var:
	//   `a := x`
	// then a.nudfn is nil, but will be defined once := is processed.
	var left *Symbol
	if t.nudfn == nil && t.ar == arName && p.tkn.id == ":=" {
		left = t
	} else {
		left = t.nud()
	}
	for rbp < p.tkn.lbp {
		t = p.tkn
		p.advance(_SYM_ANY)
		left = t.led(left)
	}
	return left
}

func (p *Parser) infix(id string, bp int, ledfn func(*Symbol, *Symbol) *Symbol) *Symbol {
	s := p.makeSymbol(id, bp)
	if ledfn != nil {
		s.ledfn = ledfn
	} else {
		s.ledfn = func(sym, left *Symbol) *Symbol {
			sym.first = left
			sym.second = p.expression(bp)
			sym.ar = arBinary
			return sym
		}
	}
	return s
}

func (p *Parser) infixr(id string, bp int, ledfn func(*Symbol, *Symbol) *Symbol) *Symbol {
	s := p.makeSymbol(id, bp)
	if ledfn != nil {
		s.ledfn = ledfn
	} else {
		s.ledfn = func(sym, left *Symbol) *Symbol {
			sym.first = left
			sym.second = p.expression(bp - 1)
			sym.ar = arBinary
			return sym
		}
	}
	return s
}

func (p *Parser) prefix(id string, nudfn func(*Symbol) *Symbol) *Symbol {
	s := p.makeSymbol(id, 0)
	if nudfn != nil {
		s.nudfn = nudfn
	} else {
		s.nudfn = func(sym *Symbol) *Symbol {
			p.scp.reserve(sym)
			sym.first = p.expression(70)
			sym.ar = arUnary
			return sym
		}
	}
	return s
}

func (p *Parser) suffix(id string) *Symbol {
	return p.infixr(id, 10, func(sym, left *Symbol) *Symbol {
		if left.id != "." && left.id != "[" && left.ar != arName {
			p.error(left, "bad lvalue")
		}
		sym.first = left
		sym.asg = true
		sym.ar = arUnary
		return sym
	})
}

func (p *Parser) assignment(id string) *Symbol {
	return p.infixr(id, 10, func(sym, left *Symbol) *Symbol {
		if left.id != "." && left.id != "[" && left.ar != arName {
			p.error(left, "bad lvalue")
		}
		sym.first = left
		sym.second = p.expression(9)
		sym.asg = true
		sym.ar = arBinary
		return sym
	})
}

// TODO : For now, it doesn't support a list of vars followed by a
// matching list of expressions (a, b, c := 1, 2, 3)
func (p *Parser) define(id string) *Symbol {
	return p.infixr(id, 10, func(sym, left *Symbol) *Symbol {
		if left.ar != arName {
			p.error(left, "expected variable name")
		}
		p.scp.define(left)
		sym.first = left
		sym.second = p.expression(9)
		sym.ar = arBinary
		return sym
	})
}

func (p *Parser) constant(id string, v interface{}) *Symbol {
	s := p.makeSymbol(id, 0)
	s.nudfn = func(sym *Symbol) *Symbol {
		p.scp.reserve(sym)
		sym.val = p.tbl[sym.id].val
		sym.ar = arLiteral
		return sym
	}
	s.val = v
	return s
}

func (p *Parser) statement() interface{} {
	n := p.tkn
	if n.stdfn != nil {
		p.advance(_SYM_ANY)
		p.scp.reserve(n)
		return n.std()
	}
	v := p.expression(0)
	if !v.asg && v.id != "(" && v.id != ":=" {
		p.error(v, "bad expression statement")
	}
	p.advance(";")
	return v
}

func (p *Parser) statements() []*Symbol {
	var a []*Symbol
	for {
		if p.tkn.id == "}" || p.tkn.id == _SYM_END {
			break
		}
		tok := p.tkn
		s := p.statement()
		switch v := s.(type) {
		case []*Symbol:
			a = append(a, v...)
		case *Symbol:
			a = append(a, v)
		default:
			p.error(tok, "unexpected statement type")
		}
	}
	return a
}

func (p *Parser) stmt(id string, stdfn func(*Symbol) interface{}) *Symbol {
	s := p.makeSymbol(id, 0)
	s.stdfn = stdfn
	return s
}

func (p *Parser) block() interface{} {
	t := p.tkn
	p.advance("{")
	return t.std()
}

// Returns a slice of imports, in pairs (one import = 2 items, first the identifier,
// then the path).
func (p *Parser) importMany() []*Symbol {
	var a []*Symbol
	for p.tkn.id != ")" {
		id, pth := p.importOne()
		a = append(a, id, pth)
	}
	p.advance(")")
	p.advance(";")
	return a
}

// Return a pair of Symbols, the identifier and the path
func (p *Parser) importOne() (id *Symbol, pth *Symbol) {
	if p.tkn.ar == arName {
		// Define in scope
		p.scp.define(p.tkn)
		id = p.tkn
		p.advance(_SYM_ANY)
	}
	var path string
	var ok bool
	if path, ok = p.tkn.val.(string); p.tkn.ar != arLiteral || !ok {
		p.error(p.tkn, "import path must be a string literal")
	}
	if id == nil {
		// No explicit identifier for the import, use the last portion of the import path
		path = path[1 : len(path)-1] // Remove \"
		if strings.HasSuffix(path, "/") {
			path = path[:len(path)-1]
		}
		idx := strings.LastIndex(path, "/")
		nm := path[idx+1:]
		if len(nm) == 0 {
			p.error(p.tkn, "invalid import path")
		}
		// Create new name Symbol for this identifier
		o := p.tbl[_SYM_NAME]
		sym := o.clone()
		sym.ar = arName
		sym.val = nm
		p.scp.define(sym)
		id = sym
	}
	pth = p.tkn
	p.advance(_SYM_ANY)
	p.advance(";")
	return
}

func (p *Parser) error(s *Symbol, msg string) {
	p.err.Add(s.pos, fmt.Sprintf("[tok: %s ; sym: %s ; val: %v] %s", s.tok, s.id, s.val, msg))
	if p.Debug {
		fmt.Println(p.err)
	}
}
