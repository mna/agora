// Top-down operator precedence parser
// Totally based on http://javascript.crockford.com/tdop/tdop.html
package parser

import (
	"fmt"

	"github.com/PuerkitoBio/agora/compiler/scanner"
	"github.com/PuerkitoBio/agora/compiler/token"
)

const (
	_SYM_END  = "(end)"
	_SYM_NAME = "(name)"
	_SYM_LIT  = "(literal)"
	_SYM_ANY  = ""
	_SYM_BAD  = "(bad)"
)

var (
	reqSymbols = []string{
		_SYM_END,
		_SYM_NAME,
		_SYM_LIT,
	}
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
	s = p.appendReturnNil(s)
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
	for _, s := range reqSymbols {
		sym := p.makeSymbol(s, 0)
		if s == _SYM_LIT {
			sym.nudfn = itselfNud
		}
	}
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
			Id:  id,
			Val: id,
			lbp: bp,
		}
		p.tbl[id] = s
	}
	return s
}

func (p *Parser) advance(id string) *Symbol {
	if id != _SYM_ANY && p.tkn.Id != id {
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
		ar Arity
	)
	if tok == token.IDENT || tok.IsKeyword() {
		o = p.scp.find(lit)
		ar = ArName
	} else if tok.IsOperator() {
		ar = ArOperator
		o, ok = p.tbl[tok.String()]
		if !ok {
			p.err.Add(pos, "unknown operator "+tok.String())
			goto scan
		}
	} else if tok.IsLiteral() { // Excluding IDENT, part of the first if
		ar = ArLiteral
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
	p.tkn.Ar = ar
	p.tkn.Val = lit
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
	if t.nudfn == nil && t.Ar == ArName && p.tkn.Id == ":=" {
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
			sym.First = left
			sym.Second = p.expression(bp)
			sym.Ar = ArBinary
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
			sym.First = left
			sym.Second = p.expression(bp - 1)
			sym.Ar = ArBinary
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
			sym.First = p.expression(70)
			sym.Ar = ArUnary
			return sym
		}
	}
	return s
}

func (p *Parser) suffix(id string) *Symbol {
	return p.infixr(id, 10, func(sym, left *Symbol) *Symbol {
		if left.Id != "." && left.Id != "[" && left.Ar != ArName {
			p.error(left, "bad lvalue")
		}
		sym.First = left
		sym.asg = true
		sym.Ar = ArStatement
		return sym
	})
}

func (p *Parser) assignment(id string) *Symbol {
	return p.infixr(id, 10, func(sym, left *Symbol) *Symbol {
		if left.Id != "." && left.Id != "[" && left.Ar != ArName {
			p.error(left, "bad lvalue")
		}
		if left.res {
			p.error(left, "cannot assign to a reserved identifier")
		}
		sym.First = left
		sym.Second = p.expression(9)
		sym.asg = true
		sym.Ar = ArBinary
		return sym
	})
}

// TODO : For now, it doesn't support a list of vars followed by a
// matching list of expressions (a, b, c := 1, 2, 3)
func (p *Parser) define(id string) *Symbol {
	return p.infixr(id, 10, func(sym, left *Symbol) *Symbol {
		if left.Ar != ArName {
			p.error(left, "expected variable name")
		}
		p.scp.define(left)
		sym.First = left
		sym.Second = p.expression(9)
		sym.Ar = ArBinary
		return sym
	})
}

func (p *Parser) constant(id string, v interface{}) *Symbol {
	s := p.makeSymbol(id, 0)
	s.nudfn = func(sym *Symbol) *Symbol {
		p.scp.reserve(sym)
		sym.Val = p.tbl[sym.Id].Val
		sym.Ar = ArLiteral
		return sym
	}
	s.Val = v
	return s
}

func (p *Parser) builtin(id string) *Symbol {
	s := p.makeSymbol(id, 0)
	s.nudfn = func(sym *Symbol) *Symbol {
		p.scp.reserve(sym)
		sym.Val = p.tbl[sym.Id].Val
		sym.Ar = ArName
		return sym
	}
	s.Val = id
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
	if !v.asg && v.Id != "(" && v.Id != ":=" {
		p.error(v, "bad expression statement")
	}
	p.advance(";")
	return v
}

func (p *Parser) statements() []*Symbol {
	var a []*Symbol
	for {
		if p.tkn.Id == "}" || p.tkn.Id == _SYM_END {
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

func (p *Parser) error(s *Symbol, msg string) {
	if s.Id != _SYM_END {
		p.err.Add(s.pos, fmt.Sprintf("[tok: %s ; sym: %s ; val: %v] %s", s.tok, s.Id, s.Val, msg))
		// Change the symbol to a (bad) symbol, returning itself in all conditions
		s.Id = _SYM_BAD
		s.ledfn = itselfLed
		s.nudfn = itselfNud
		s.stdfn = itselfStd
		if p.Debug {
			fmt.Println((*p.err)[p.err.Len()-1])
		}
	}
}
