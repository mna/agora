package parser

func makeFuncParser(p *Parser, prefix bool) func(*Symbol) *Symbol {
	return func(sym *Symbol) *Symbol {
		var a []*Symbol
		if !prefix && p.tkn.ar == arName { // Only for statement notation
			p.scp.define(p.tkn)
			sym.name = p.tkn.val.(string)
			p.advance(_SYM_ANY)
		}
		p.newScope()
		p.advance("(")
		if p.tkn.id != ")" {
			for {
				if p.tkn.ar != arName {
					p.error(p.tkn, "expected a parameter name")
				}
				p.scp.define(p.tkn)
				a = append(a, p.tkn)
				p.advance(_SYM_ANY)
				if p.tkn.id != "," {
					break
				}
				p.advance(",")
			}
		}
		sym.first = a
		p.advance(")")
		p.advance("{")
		sym.second = p.statements()
		p.advance("}")
		if !prefix { // Don't consume the ending semicolon when func is an expression
			p.advance(";")
		}
		sym.ar = arFunction
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
		sym.first = left
		sym.second = p.expression(0)
		p.advance(":")
		sym.third = p.expression(0)
		sym.ar = arTernary
		return sym
	})

	// The dot (selector) operator
	p.infix(".", 80, func(sym, left *Symbol) *Symbol {
		sym.first = left
		if p.tkn.ar != arName {
			p.error(p.tkn, "expected a field name")
		}
		p.tkn.ar = arLiteral
		sym.second = p.tkn
		sym.ar = arBinary
		p.advance(_SYM_ANY)
		return sym
	})

	// The array-notation field selector operator
	p.infix("[", 80, func(sym, left *Symbol) *Symbol {
		sym.first = left
		sym.second = p.expression(0)
		sym.ar = arBinary
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
		sym.first = p.expression(0)
		sym.second = p.block()
		p.advance(";")
		sym.ar = arStatement
		return sym
	})

	// If statement
	p.stmt("if", func(sym *Symbol) interface{} {
		sym.first = p.expression(0)
		sym.second = p.block()
		if p.tkn.id == "else" {
			p.scp.reserve(p.tkn)
			p.advance("else")
			if p.tkn.id == "if" {
				sym.third = p.statement()
			} else {
				sym.third = p.block()
				p.advance(";")
			}
		} else {
			p.advance(";")
		}
		sym.ar = arStatement
		return sym
	})

	// break statement
	p.stmt("break", func(sym *Symbol) interface{} {
		p.advance(";")
		if p.tkn.id != "}" && p.tkn.id != _SYM_END {
			p.error(p.tkn, "unreachable statement")
		}
		sym.ar = arStatement
		return sym
	})

	// return statement
	p.stmt("return", func(sym *Symbol) interface{} {
		if p.tkn.id != ";" {
			sym.first = p.expression(0)
		}
		p.advance(";")
		if p.tkn.id != "}" && p.tkn.id != _SYM_END {
			p.error(p.tkn, "unreachable statement")
		}
		sym.ar = arStatement
		return sym
	})

	// TODO : Must be the first statement(s) in a file
	// import statement
	p.stmt("import", func(sym *Symbol) interface{} {
		if p.tkn.id == "(" {
			p.advance("(")
			sym.first = p.importMany()
		} else {
			id, pth := p.importOne()
			sym.first = []*Symbol{id, pth}
		}
		sym.ar = arImport
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
		if p.tkn.id != ")" {
			for {
				a = append(a, p.expression(0))
				if p.tkn.id != "," {
					break
				}
				p.advance(",")
			}
		}
		p.advance(")")
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
				p.error(left, "expected a variable name")
			}
		}
		return sym
	})

	// The `this` keyword
	p.makeSymbol("this", 0).nudfn = func(sym *Symbol) *Symbol {
		p.scp.reserve(sym)
		sym.ar = arThis
		return sym
	}

	// The object literal notation
	p.prefix("{", func(sym *Symbol) *Symbol {
		var a []*Symbol
		if p.tkn.id != "}" {
			for {
				n := p.tkn
				if n.ar != arName && n.ar != arLiteral {
					p.error(n, "bad key")
				}
				p.advance(_SYM_ANY)
				p.advance(":")
				v := p.expression(0)
				v.key = n.val
				a = append(a, v)
				if p.tkn.id != "," {
					break
				}
				p.advance(",")
			}
		}
		p.advance("}")
		sym.first = a
		sym.ar = arUnary
		return sym
	})
	// TODO : No array literal ("[14, 83, "toto"]") for now

	// Increment/decrement statements
	p.suffix("--")
	p.suffix("++")
}
