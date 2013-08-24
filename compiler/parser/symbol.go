package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/agora/compiler/token"
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
	arFunction
	arImport
)

func itselfLed(s, left *Symbol) *Symbol {
	return left
}

func itselfNud(s *Symbol) *Symbol {
	return s
}

func itselfStd(s *Symbol) interface{} {
	return s
}

type Symbol struct {
	p      *Parser
	id     string
	val    interface{}
	name   string
	key    interface{}
	lbp    int
	ar     arity
	res    bool
	asg    bool
	tok    token.Token
	pos    token.Position
	first  interface{} // May all be []*Symbol or *Symbol
	second interface{}
	third  interface{}

	nudfn func(*Symbol) *Symbol
	ledfn func(*Symbol, *Symbol) *Symbol
	stdfn func(*Symbol) interface{} // May return []*Symbol or *Symbol
}

func (s Symbol) clone() *Symbol {
	return &Symbol{
		s.p,
		s.id,
		s.val,
		s.name,
		s.key,
		s.lbp,
		s.ar,
		s.res,
		s.asg,
		s.tok,
		s.pos,
		s.first,
		s.second,
		s.third,
		s.nudfn,
		s.ledfn,
		s.stdfn,
	}
}

func (s *Symbol) led(left *Symbol) *Symbol {
	if s.ledfn == nil {
		s.p.error(s, "missing operator")
	}
	return s.ledfn(s, left)
}

func (s *Symbol) std() interface{} {
	if s.stdfn == nil {
		s.p.error(s, "invalid operation")
	}
	return s.stdfn(s)
}

func (s *Symbol) nud() *Symbol {
	if s.nudfn == nil {
		s.p.error(s, "undefined")
	}
	return s.nudfn(s)
}

func (s *Symbol) String() string {
	return s.indentString(0)
}

func (s *Symbol) indentString(ind int) string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%-20s; %v", s.id, s.val))
	if s.name != "" {
		buf.WriteString(fmt.Sprintf(" (nm: %s)", s.name))
	} else if s.key != nil {
		buf.WriteString(fmt.Sprintf(" (key: %s)", s.key))
	}
	buf.WriteString("\n")

	fmtChild := func(idx int, child interface{}) {
		if child != nil {
			switch v := child.(type) {
			case []*Symbol:
				for i, c := range v {
					buf.WriteString(fmt.Sprintf("%s[%d.%d] %s", strings.Repeat(" ", (ind+1)*3), idx, i+1, c.indentString(ind+1)))
				}
			case *Symbol:
				buf.WriteString(fmt.Sprintf("%s[%d] %s", strings.Repeat(" ", (ind+1)*3), idx, v.indentString(ind+1)))
			}
		}
	}
	fmtChild(1, s.first)
	fmtChild(2, s.second)
	fmtChild(3, s.third)
	return buf.String()
}
