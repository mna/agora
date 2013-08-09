package parser

import (
	"fmt"

	"github.com/PuerkitoBio/goblin/compiler/ast"
	"github.com/PuerkitoBio/goblin/compiler/scanner"
	"github.com/PuerkitoBio/goblin/compiler/token"
)

type Parser struct {
	scanner.Scanner
	filename   string
	tok        token.Token
	lit        string
	matchedLit string
	ast        *ast.Module
}

func (p *Parser) Parse(fn string, src []byte) (*ast.Module, error) {
	errl := new(scanner.ErrorList)
	p.Init(fn, src, errl.Add)
	p.filename = fn
	p.tok, p.lit = p.Scan()
	p.module()
	return p.ast, errl.Err()
}

func (p *Parser) match(t token.Token) bool {
	if p.tok == t {
		p.matchedLit = p.lit
		p.tok, p.lit = p.Scan()
		return true
	}
	return false
}

func (p *Parser) expect(t token.Token) bool {
	if p.match(t) {
		return true
	}
	p.Error(fmt.Sprintf("expected %s, found %s", t, p.tok))
	return false
}

func (p *Parser) module() {
	p.ast = ast.NewModule(p.filename)
	p.importStmt()
	p.stmtList()
}

func (p *Parser) importStmt() bool {
	if p.match(token.IMPORT) {
		if p.match(token.LPAREN) {
			// Possible list of imports
			for p.importSpec() {
			}
			p.expect(token.RPAREN)
			p.expect(token.SEMICOLON)
		} else {
			// Single import
			if !p.importSpec() {
				p.Error("ImportStmt : expected import path")
			}
		}
		return true
	}
	return false
}

func (p *Parser) importSpec() bool {
	var id, path string
	if p.match(token.IDENT) {
		id = p.matchedLit
	}
	if p.match(token.STRING) {
		path = p.matchedLit
	} else if id != "" {
		p.Error("ImportSpec : missing import path")
	}
	if path != "" || id != "" {
		p.expect(token.SEMICOLON)
		p.ast.AddImport(path, id)
		return true
	}
	return false
}

func (p *Parser) stmtList() bool {
	one := false
	for p.statement() {
		one = true
	}
	return one
}

func (p *Parser) statement() bool {
	return p.simpleStmt || p.returnStmt || p.breakStmt || p.continueStmt || p.ifStmt || p.forStmt
}

func (p *Parser) simpleStmt() bool {
	// TODO : Later...
}

func (p *Parser) returnStmt() bool {
	if p.match(token.RETURN) {
		p.expression()
		return true
	}
	return false
}

func (p *Parser) breakStmt() bool {
	return p.match(token.BREAK)
}

func (p *Parser) continueStmt() bool {
	return p.match(token.CONTINUE)
}

func (p *Parser) ifStmt() bool {
	if p.match(token.IF) {
		if !p.expression() {
			p.Error("IfStmt : expected an expression")
		}
		if !p.block() {
			p.Error("IfStmt : expected a body")
		}
		if p.match(token.ELSE) {
			if !p.ifStmt() && !p.block() {
				p.Error("IfStmt : expected an if or a body after the else keyword")
			}
		}
		return true
	}
	return false
}

func (p *Parser) forStmt() bool {
	if p.match(token.FOR) {
		// For Clause must be checked first, range next, and condition last
		if !(p.forClause() || p.rangeClause() || p.condition) {
			p.Error("ForStmt : expected a condition, a for clause or a range clause")
		}
		if !p.block() {
			p.Error("ForStmt : expected a for body")
		}
		return true
	}
	return false
}

func (p *Parser) condition() bool {
	return p.expression()
}

func (p *Parser) forClause() bool {
	p.initStmt() // Optional
	if !p.match(token.SEMICOLON) {
		// TODO : Rewind the initStmt!
		return false
	}
	// Else, assume this is a for clause, the rest is expected
	p.condition() // Optional
	p.expect(token.SEMICOLON)
	p.postStmt() // Optional
	return true
}

func (p *Parser) rangeClause() bool {
	var ok bool
	if p.expressionList() {
		ok = p.expect(token.ASSIGN)
	} else if p.identifierList() {
		ok = p.expect(token.DEFINE)
	}
	if !ok {
		// TODO : Rewind?!?!?!
		return false
	}
	p.expect(token.RANGE)
	if !p.expression() {
		p.Error("RangeClause : expected an expression")
	}
	return true
}
