// Top-down operator precedence parser
// Totally based on http://javascript.crockford.com/tdop/tdop.html
package tdop

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goblin/compiler/scanner"
	"github.com/PuerkitoBio/goblin/compiler/token"
)

var (
	// TODO : For now, package-level scope, but should be in a Parser struct
	symtbl  = make(map[string]*symbol) // Symbol table
	curTok  *symbol                    // Current token in symbol representation
	Scanner *scanner.Scanner
	curScp  *scope
)

func Parse(fn string, src []byte) {
	Scanner.Init(fn, src, nil)
	newScope()
	advance("")
	s := statements()
	advance("(end)")
	curScp.pop()
	for _, v := range s {
		fmt.Println(v)
	}
}

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
	arFunction
)

type scope struct {
	def    map[string]*symbol
	parent *scope
}

func newScope() *scope {
	curScp = &scope{
		make(map[string]*symbol),
		curScp,
	}
	return curScp
}

func itself(s *symbol) *symbol {
	return s
}

func (s *scope) define(n *symbol) *symbol {
	t, ok := s.def[n.val.(string)]
	if ok {
		if t.res {
			error("already reserved")
		} else {
			error("already defined")
		}
		panic("unreachable")
	}
	s.def[n.val.(string)] = n
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
	if t, ok := s.def[n.val.(string)]; ok {
		if t.res {
			return
		}
		if t.ar == arName {
			error("already defined")
		}
	}
	s.def[n.val.(string)] = n
	n.res = true
}

func clone(ori *symbol) *symbol {
	return &symbol{
		ori.id,
		ori.val,
		ori.name,
		ori.key,
		ori.lbp,
		ori.ar,
		ori.res,
		ori.asg,
		ori.scp,
		ori.first,
		ori.second,
		ori.third,
		ori.nudfn,
		ori.ledfn,
		ori.stdfn,
	}
}

type symbol struct {
	id     string
	val    interface{}
	name   string
	key    interface{}
	lbp    int
	ar     arity
	res    bool
	asg    bool
	scp    *scope
	first  interface{} // May all be []*symbol or *symbol
	second interface{}
	third  interface{}

	nudfn func(*symbol) *symbol
	ledfn func(*symbol, *symbol) *symbol
	stdfn func(*symbol) interface{} // May return []*symbol or *symbol
}

func (s *symbol) nud() *symbol {
	if s.nudfn == nil {
		error(fmt.Sprintf("undefined %s: %s", s.id, s.val))
	}
	return s.nudfn(s)
}

func (s *symbol) String() string {
	return s.indentString(0)
}

func (s *symbol) indentString(ind int) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%-20s; %s\n", s.id, s.val))

	fmtChild := func(idx int, child interface{}) {
		if child != nil {
			switch v := child.(type) {
			case []*symbol:
				for i, c := range v {
					buf.WriteString(fmt.Sprintf("%s[%d.%d] %s", strings.Repeat(" ", (ind+1)*3), idx, i+1, c.indentString(ind+1)))
				}
			case *symbol:
				buf.WriteString(fmt.Sprintf("%s[%d] %s", strings.Repeat(" ", (ind+1)*3), idx, v.indentString(ind+1)))
			}
		}
	}
	fmtChild(1, s.first)
	fmtChild(2, s.second)
	fmtChild(3, s.third)
	return buf.String()
}

func (s *symbol) led(left *symbol) *symbol {
	if s.ledfn == nil {
		error("missing operator")
	}
	return s.ledfn(s, left)
}

func (s *symbol) std() interface{} {
	if s.stdfn == nil {
		error("invalid operation")
	}
	return s.stdfn(s)
}

func makeSymbol(id string, bp int) *symbol {
	s, ok := symtbl[id]
	if ok {
		if bp >= s.lbp {
			s.lbp = bp
		}
	} else {
		s = &symbol{
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
	var (
		tok token.Token
		lit string
		pos token.Position
	)
	for tok, lit, pos = Scanner.Scan(); tok == token.ILLEGAL || tok == token.COMMENT; tok, lit, pos = Scanner.Scan() {
		// Skip Illegal and Comment tokens
	}
	fmt.Println("SCAN: ", tok, lit, pos)
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
		o, ok = symtbl[tok.String()]
		if !ok {
			error("unknown operator " + tok.String())
		}
	} else if tok.IsLiteral() { // Excluding IDENT, part of the first if
		ar = arLiteral
		o = symtbl["(literal)"]
	} else if tok == token.EOF {
		o = symtbl["(end)"]
		curTok = o
		return o
	} else {
		error("unexpected token " + tok.String())
	}
	curTok = clone(o)
	curTok.ar = ar
	curTok.val = lit
	return curTok
}

func expression(rbp int) *symbol {
	t := curTok
	advance("")
	// TODO : Special case if in the process of defining a new var:
	// a := x
	// then a.nudfn is nil, but will be defined once := is processed.
	var left *symbol
	if t.nudfn == nil && t.ar == arName && curTok.id == ":=" {
		left = t
	} else {
		left = t.nud()
	}
	for rbp < curTok.lbp {
		t = curTok
		advance("")
		left = t.led(left)
	}
	return left
}

func infix(id string, bp int, ledfn func(*symbol, *symbol) *symbol) *symbol {
	s := makeSymbol(id, bp)
	if ledfn != nil {
		s.ledfn = ledfn
	} else {
		s.ledfn = func(sym, left *symbol) *symbol {
			sym.first = left
			sym.second = expression(bp)
			sym.ar = arBinary
			return sym
		}
	}
	return s
}

func infixr(id string, bp int, ledfn func(*symbol, *symbol) *symbol) *symbol {
	s := makeSymbol(id, bp)
	if ledfn != nil {
		s.ledfn = ledfn
	} else {
		s.ledfn = func(sym, left *symbol) *symbol {
			sym.first = left
			sym.second = expression(bp - 1)
			sym.ar = arBinary
			return sym
		}
	}
	return s
}

func prefix(id string, nudfn func(*symbol) *symbol) *symbol {
	s := makeSymbol(id, 0)
	if nudfn != nil {
		s.nudfn = nudfn
	} else {
		s.nudfn = func(sym *symbol) *symbol {
			curScp.reserve(sym)
			sym.first = expression(70)
			sym.ar = arUnary
			return sym
		}
	}
	return s
}

func assignment(id string) *symbol {
	return infixr(id, 10, func(sym, left *symbol) *symbol {
		if left.id != "." && left.id != "[" && left.ar != arName {
			error("bad lvalue")
		}
		sym.first = left
		sym.second = expression(9)
		sym.asg = true
		sym.ar = arBinary
		return sym
	})
}

// TODO : For now, it doesn't support a list of vars followed by a matching list of expressions (a, b, c := 1, 2, 3)
func define(id string) *symbol {
	return infixr(id, 10, func(sym, left *symbol) *symbol {
		if left.ar != arName {
			error("expected variable name")
		}
		curScp.define(left)
		sym.first = left
		sym.second = expression(9)
		sym.ar = arBinary
		return sym
	})
}

func constant(id string, v interface{}) *symbol {
	s := makeSymbol(id, 0)
	s.nudfn = func(sym *symbol) *symbol {
		curScp.reserve(sym)
		sym.val = symtbl[sym.id].val
		sym.ar = arLiteral
		return sym
	}
	s.val = v
	return s
}

func statement() interface{} {
	n := curTok
	if n.stdfn != nil {
		advance("")
		curScp.reserve(n)
		return n.std()
	}
	v := expression(0)
	if !v.asg && v.id != "(" && v.id != ":=" {
		error("bad expression statement: " + v.id)
	}
	advance(";")
	return v
}

func statements() []*symbol {
	var a []*symbol
	for {
		if curTok.id == "}" || curTok.id == "(end)" {
			break
		}
		s := statement()
		switch v := s.(type) {
		case []*symbol:
			a = append(a, v...)
		case *symbol:
			a = append(a, v)
		default:
			panic("unexpected type")
		}
	}
	return a
}

func stmt(id string, stdfn func(*symbol) interface{}) *symbol {
	s := makeSymbol(id, 0)
	s.stdfn = stdfn
	return s
}

func block() interface{} {
	t := curTok
	advance("{")
	return t.std()
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

	infix("+", 50, nil)
	infix("-", 50, nil)
	infix("*", 60, nil)
	infix("/", 60, nil)
	infix("%", 60, nil)
	infix("==", 40, nil)
	infix("<", 40, nil)
	infix(">", 40, nil)
	infix("!=", 40, nil)
	infix("<=", 40, nil)
	infix(">=", 40, nil)
	// Ternary operator?
	infix("?", 20, func(sym, left *symbol) *symbol {
		sym.first = left
		sym.second = expression(0)
		advance(":")
		sym.third = expression(0)
		sym.ar = arTernary
		return sym
	})
	// The dot (selector) operator
	infix(".", 80, func(sym, left *symbol) *symbol {
		sym.first = left
		if curTok.ar != arName {
			error("expected a field name")
		}
		curTok.ar = arLiteral
		sym.second = curTok
		sym.ar = arBinary
		advance("")
		return sym
	})
	// The array-notation field selector operator
	infix("[", 80, func(sym, left *symbol) *symbol {
		sym.first = left
		sym.second = expression(0)
		sym.ar = arBinary
		advance("]")
		return sym
	})
	// The logical operators
	infixr("&&", 30, nil)
	infixr("||", 30, nil)

	prefix("-", nil)
	prefix("!", nil)
	prefix("(", func(sym *symbol) *symbol {
		e := expression(0)
		advance(")")
		return e
	})

	assignment("=")
	assignment("+=")
	assignment("-=")
	assignment("*=")
	assignment("/=")
	assignment("%=")

	constant("true", true)
	constant("false", false)
	constant("nil", nil)

	makeSymbol("(literal)", 0).nudfn = itself

	stmt("{", func(sym *symbol) interface{} {
		// TODO : New scope? Runtime would have to handle this, so for now, no new scope in blocks.
		a := statements()
		advance("}")
		return a
	})
	define(":=")
	// TODO : This supports the for [condition] notation, nothing else
	stmt("for", func(sym *symbol) interface{} {
		sym.first = expression(0)
		sym.second = block()
		sym.ar = arStatement
		return sym
	})
	stmt("if", func(sym *symbol) interface{} {
		sym.first = expression(0)
		sym.second = block()
		if curTok.id == "else" {
			curScp.reserve(curTok)
			advance("else")
			if curTok.id == "if" {
				sym.third = statement()
			} else {
				sym.third = block()
			}
		}
		sym.ar = arStatement
		return sym
	})
	stmt("break", func(sym *symbol) interface{} {
		advance(";")
		if curTok.id != "}" && curTok.id != "(end)" {
			error("unreachable statement")
		}
		sym.ar = arStatement
		return sym
	})
	stmt("return", func(sym *symbol) interface{} {
		fmt.Println("return1 ", curTok.id)
		if curTok.id != ";" {
			sym.first = expression(0)
			fmt.Println("return2 ", curTok.id)
		}
		advance(";")
		fmt.Println("return3 ", curTok.id)
		if curTok.id != "}" && curTok.id != "(end)" {
			error("unreachable statement: " + curTok.id)
		}
		sym.ar = arStatement
		return sym
	})
	// func can be both an expression prefix:
	//   fnAdd := func(x, y) {return x+y}
	// or a statement:
	//   func Add(x, y) {return x+y}
	// TODO : Make this DRY and much cleaner
	prefix("func", func(sym *symbol) *symbol {
		var a []*symbol
		fmt.Println("FUNC PREFIX")
		if curTok.ar == arName {
			fmt.Println("FUNC define in scope name " + curTok.val.(string))
			curScp.define(curTok)
			sym.name = curTok.val.(string)
			advance("")
		}
		newScope()
		advance("(")
		if curTok.id != ")" {
			for {
				if curTok.ar != arName {
					error("expected a parameter name")
				}
				curScp.define(curTok)
				a = append(a, curTok)
				advance("")
				if curTok.id != "," {
					break
				}
				advance(",")
			}
		}
		sym.first = a
		advance(")")
		advance("{")
		sym.second = statements()
		advance("}")
		// Don't consume the ending prefix when func is an expression
		sym.ar = arFunction
		curScp.pop()
		return sym
	})
	stmt("func", func(sym *symbol) interface{} {
		var a []*symbol
		fmt.Println("FUNC STMT")
		// The func name (e.g. func Add(x, y)...) should be defined in both
		// the parent scope and the inner scope of the function. But then, just
		// define in the parent scope, which will make it available in the inner scope.
		if curTok.ar == arName {
			fmt.Println("FUNC define in scope name " + curTok.val.(string))
			curScp.define(curTok)
			sym.name = curTok.val.(string)
			advance("")
		}
		newScope()
		advance("(")
		if curTok.id != ")" {
			for {
				if curTok.ar != arName {
					error("expected a parameter name")
				}
				curScp.define(curTok)
				a = append(a, curTok)
				advance("")
				if curTok.id != "," {
					break
				}
				advance(",")
			}
		}
		sym.first = a
		advance(")")
		advance("{")
		sym.second = statements()
		advance("}")
		advance(";")
		sym.ar = arFunction
		curScp.pop()
		return sym
	})
	infix("(", 80, func(sym, left *symbol) *symbol {
		var a []*symbol
		if curTok.id != ")" {
			for {
				a = append(a, expression(0))
				if curTok.id != "," {
					break
				}
				advance(",")
			}
		}
		advance(")")
		if left.id == "." || left.id == "[" {
			sym.ar = arTernary
			sym.first = left.first
			sym.second = left.second
			sym.third = a
		} else {
			sym.ar = arBinary
			sym.first = left
			sym.second = a
			if (left.ar != arUnary || left.id != "func") &&
				left.ar != arName && left.id != "(" &&
				left.id != "&&" && left.id != "||" && left.id != "?" {
				error("expected a variable name")
			}
		}
		return sym
	})
	makeSymbol("this", 0).nudfn = func(sym *symbol) *symbol {
		curScp.reserve(sym)
		sym.ar = arThis
		return sym
	}
	prefix("{", func(sym *symbol) *symbol {
		var a []*symbol
		if curTok.id != "}" {
			for {
				n := curTok
				if n.ar != arName && n.ar != arLiteral {
					error("bad key")
				}
				advance("")
				advance(":")
				v := expression(0)
				v.key = n.val
				a = append(a, v)
				if curTok.id != "," {
					break
				}
				advance(",")
			}
		}
		advance("}")
		sym.first = a
		sym.ar = arUnary
		return sym
	})
	// TODO : No array literal ("[14, 83, "toto"]") for now
}
