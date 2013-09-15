package parser

// This function defines the whole grammar of the language.
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

	// Statement
	p.stmt("{", func(sym *Symbol) interface{} {
		a := p.statements()
		p.advance("}")
		return a
	})

	// The define operator, to declare-assign variables
	p.define(":=")

	// For loop
	p.stmt("for", func(sym *Symbol) interface{} {
		// Check for the infinite loop form (i.e. `for {}`). If this is the case,
		// sym.First is nil, while sym.Second holds the body.
		sym.First = nil
		if p.tkn.Id != "{" {
			f := p.expression(0)
			if p.tkn.Id == "{" {
				// Single expression form (i.e. `while`)
				sym.First = f
			} else {
				// 3-part for (for stmt ; expr ; stmt {})
				pt1 := f
				p.advance(";")
				pt2 := p.expression(0)
				p.advance(";")
				pt3 := p.expression(0)
				// Special case for the 3-part for, each part is in a slice of interface{}
				sym.First = []interface{}{pt1, pt2, pt3}
			}
		}
		sym.Second = p.block()
		p.advance(";")
		sym.Ar = ArStatement
		return sym
	})

	// If statement
	p.stmt("if", func(sym *Symbol) interface{} {
		sym.First = p.expression(0)
		sym.Second = p.block()
		sym.Third = nil
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

	// continue statement
	p.stmt("continue", func(sym *Symbol) interface{} {
		p.advance(";")
		sym.Ar = ArStatement
		return sym
	})

	// debug statement
	p.stmt("debug", func(sym *Symbol) interface{} {
		sym.First = nil
		if p.tkn.Id != ";" {
			// Evaluate the number of stack traces to print
			if p.tkn.Id != "(literal)" {
				p.error(p.tkn, "expected number literal")
			}
			sym.First = p.tkn
			p.advance(_SYM_ANY)
		}
		p.advance(";")
		sym.Ar = ArStatement
		return sym
	})

	// return statement
	p.stmt("return", func(sym *Symbol) interface{} {
		sym.First = p.expression(0)
		p.advance(";")
		if p.tkn.Id != "}" && p.tkn.Id != _SYM_END {
			p.error(p.tkn, "unreachable statement")
		}
		sym.Ar = ArStatement
		return sym
	})

	// import builtin
	p.builtin("import")
	p.builtin("panic")
	p.builtin("recover")
	p.builtin("len")
	p.builtin("keys")

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
			sym.Third = nil
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
		sym.Ar = ArName
		return sym
	}
	// The `args` keyword
	p.makeSymbol("args", 0).nudfn = func(sym *Symbol) *Symbol {
		p.scp.reserve(sym)
		sym.Ar = ArName
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

	// Increment/decrement statements
	p.suffix("--")
	p.suffix("++")
}

func makeFuncParser(p *Parser, prefix bool) func(*Symbol) *Symbol {
	return func(sym *Symbol) *Symbol {
		var a []*Symbol
		sym.Name = ""
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
