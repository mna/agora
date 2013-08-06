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
