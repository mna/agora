package parser

// TODO : Make "import" and "debug" built-in functions, import
// returns the imported module's value to a variable. `debug` adds
// a DUMP instruction at this location, takes a number as parameter
// (frame size).

func makeFuncParser(p *Parser, prefix bool) func(*Symbol) *Symbol {
	return func(sym *Symbol) *Symbol {
		var a []*Symbol
		if !prefix && p.tkn.Ar == ArName { // Only for statement notation
			p.scp.define(p.tkn)
			sym.Name = p.tkn.Val.(string)
			p.advance(_SYM_ANY)
		}
		p.newScope()
		p.advance("(")
		if p.tkn.Id != ")" {
			for {
				if p.tkn.Ar != ArName {
					p.error(p.tkn, "expected a parameter name")
				}
				p.scp.define(p.tkn)
				a = append(a, p.tkn)
				p.advance(_SYM_ANY)
				if p.tkn.Id != "," {
					break
				}
				p.advance(",")
			}
		}
		sym.First = a
		p.advance(")")
		p.advance("{")
		stmts := p.statements()
		stmts = p.appendReturnNil(stmts)
		sym.Second = stmts
		p.advance("}")
		if !prefix { // Don't consume the ending semicolon when func is an expression
			p.advance(";")
		}
		sym.Ar = ArFunction
		p.popScope()
		return sym
	}
}

func makeFuncParserIface(p *Parser, prefix bool) func(*Symbol) interface{} {
	f := makeFuncParser(p, prefix)
	return func(s *Symbol) interface{} {
		return f(s)
	}
}

func (p *Parser) appendReturnNil(s []*Symbol) []*Symbol {
	// Make sure the function ends with a return statement, adding a return nil otherwise
	if l := len(s); l == 0 || s[l-1].Id != "return" {
		ret := p.makeSymbol("return", 0).clone()
		ret.Ar = ArStatement
		ret.First = p.makeSymbol("nil", 0).clone()
		s = append(s, ret)
	}
	return s
}

func (p *Parser) defineGrammar() {
	// Ponctuation symbols
	p.makeSymbol(":", 0)
	p.makeSymbol(";", 0)
	p.makeSymbol(",", 0)
	p.makeSymbol(")", 0)
	p.makeSymbol("]", 0)
	p.makeSymbol("}", 0)
	p.makeSymbol("else", 0)

	// Infix operators
	p.infix("+", 50, nil)  // Add
	p.infix("-", 50, nil)  // Subtract
	p.infix("*", 60, nil)  // Multiply
	p.infix("/", 60, nil)  // Divide
	p.infix("%", 60, nil)  // Modulo
	p.infix("==", 40, nil) // Equals
	p.infix("<", 40, nil)  // Lower than
	p.infix(">", 40, nil)  // Greater than
	p.infix("!=", 40, nil) // Not equal
	p.infix("<=", 40, nil) // Lower than or equal
	p.infix(">=", 40, nil) // Greater than or equal

	// Ternary operator
	p.infix("?", 20, func(sym, left *Symbol) *Symbol {
		sym.First = left
		sym.Second = p.expression(0)
		p.advance(":")
		sym.Third = p.expression(0)
		sym.Ar = ArTernary
		return sym
	})

	// The dot (selector) operator
	p.infix(".", 80, func(sym, left *Symbol) *Symbol {
		sym.First = left
		if p.tkn.Ar != ArName {
			p.error(p.tkn, "expected a field name")
		}
		p.tkn.Ar = ArLiteral
		sym.Second = p.tkn
		sym.Ar = ArBinary
		p.advance(_SYM_ANY)
		return sym
	})

	// The array-notation field selector operator
	p.infix("[", 80, func(sym, left *Symbol) *Symbol {
		sym.First = left
		sym.Second = p.expression(0)
		sym.Ar = ArBinary
		p.advance("]")
		return sym
	})

	// The logical operators
	p.infixr("&&", 30, nil)
	p.infixr("||", 30, nil)

	// The unary operators
	p.prefix("-", nil) // Unary minus
	p.prefix("!", nil) // Not

	// The expression grouping operator
	p.prefix("(", func(sym *Symbol) *Symbol {
		e := p.expression(0)
		p.advance(")")
		return e
	})

	// The assignment operators
	p.assignment("=")
	p.assignment("+=")
	p.assignment("-=")
	p.assignment("*=")
	p.assignment("/=")
	p.assignment("%=")

	// Language constants
	p.constant("true", true)   // boolean true
	p.constant("false", false) // boolean false
	p.constant("nil", nil)     // nil value
	p.constant("args", "args") // The special variable args

	// Statement
	p.stmt("{", func(sym *Symbol) interface{} {
		a := p.statements()
		p.advance("}")
		return a
	})

	// The define operator, to declare-assign variables
	p.define(":=")

	// TODO : This supports the for [condition] notation, nothing else
	// For loop
	p.stmt("for", func(sym *Symbol) interface{} {
		sym.First = p.expression(0)
		sym.Second = p.block()
		p.advance(";")
		sym.Ar = ArStatement
		return sym
	})

	// If statement
	p.stmt("if", func(sym *Symbol) interface{} {
		sym.First = p.expression(0)
		sym.Second = p.block()
		if p.tkn.Id == "else" {
			p.scp.reserve(p.tkn)
			p.advance("else")
			if p.tkn.Id == "if" {
				sym.Third = p.statement()
			} else {
				sym.Third = p.block()
				p.advance(";")
			}
		} else {
			p.advance(";")
		}
		sym.Ar = ArStatement
		return sym
	})

	// break statement
	p.stmt("break", func(sym *Symbol) interface{} {
		p.advance(";")
		if p.tkn.Id != "}" && p.tkn.Id != _SYM_END {
			p.error(p.tkn, "unreachable statement")
		}
		sym.Ar = ArStatement
		return sym
	})

	// return statement
	p.stmt("return", func(sym *Symbol) interface{} {
		if p.tkn.Id != ";" {
			sym.First = p.expression(0)
		}
		p.advance(";")
		if p.tkn.Id != "}" && p.tkn.Id != _SYM_END {
			p.error(p.tkn, "unreachable statement")
		}
		sym.Ar = ArStatement
		return sym
	})

	// TODO : Must be the first statement(s) in a file
	// import statement
	p.stmt("import", func(sym *Symbol) interface{} {
		if p.tkn.Id == "(" {
			p.advance("(")
			sym.First = p.importMany()
		} else {
			id, pth := p.importOne()
			sym.First = []*Symbol{id, pth}
		}
		sym.Ar = ArImport
		return sym
	})

	// func can be both an expression prefix:
	//   fnAdd := func(x, y) {return x+y}
	// or a statement:
	//   func Add(x, y) {return x+y}
	p.prefix("func", makeFuncParser(p, true))
	p.stmt("func", makeFuncParserIface(p, false))

	// The function/method call parser
	p.infix("(", 80, func(sym, left *Symbol) *Symbol {
		var a []*Symbol
		if p.tkn.Id != ")" {
			for {
				a = append(a, p.expression(0))
				if p.tkn.Id != "," {
					break
				}
				p.advance(",")
			}
		}
		p.advance(")")
		if left.Id == "." || left.Id == "[" {
			sym.Ar = ArTernary
			sym.First = left.First
			sym.Second = left.Second
			sym.Third = a
		} else {
			sym.Ar = ArBinary
			sym.First = left
			sym.Second = a
			if (left.Ar != ArUnary || left.Id != "func") &&
				left.Ar != ArName && left.Id != "(" &&
				left.Id != "&&" && left.Id != "||" && left.Id != "?" {
				p.error(left, "expected a variable name")
			}
		}
		return sym
	})

	// The `this` keyword
	p.makeSymbol("this", 0).nudfn = func(sym *Symbol) *Symbol {
		p.scp.reserve(sym)
		sym.Ar = ArThis
		return sym
	}

	// The object literal notation
	p.prefix("{", func(sym *Symbol) *Symbol {
		var a []*Symbol
		if p.tkn.Id != "}" {
			for {
				n := p.tkn
				if n.Ar != ArName && n.Ar != ArLiteral {
					p.error(n, "bad key")
				}
				p.advance(_SYM_ANY)
				p.advance(":")
				v := p.expression(0)
				v.Key = n.Val
				a = append(a, v)
				if p.tkn.Id != "," {
					break
				}
				p.advance(",")
			}
		}
		p.advance("}")
		sym.First = a
		sym.Ar = ArUnary
		return sym
	})
	// TODO : No array literal ("[14, 83, "toto"]") for now

	// Increment/decrement statements
	p.suffix("--")
	p.suffix("++")
}
