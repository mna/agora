package scanner

import (
	"fmt"
	"testing"

	"github.com/PuerkitoBio/goblin/compiler/token"
)

var (
	cases = []struct {
		src []byte
		exp []token.Token
	}{
		{
			src: []byte(`aB := 5
					return aB`),
			exp: []token.Token{
				token.IDENT,
				token.DEFINE,
				token.INT,
				token.SEMICOLON,
				token.RETURN,
				token.IDENT,
				token.SEMICOLON,
			},
		},
		{
			src: []byte(`
a := 7
b := 10
add := a + b
sub := a - b
mul := a * b
div := a / b // TODO : Should div return a float even for Int?
mod := b % a
not := !a
unm := -a`),
			exp: []token.Token{
				token.IDENT,
				token.DEFINE,
				token.INT,
				token.SEMICOLON,
				token.IDENT,
				token.DEFINE,
				token.INT,
				token.SEMICOLON,
				token.IDENT,
				token.DEFINE,
				token.IDENT,
				token.ADD,
				token.IDENT,
				token.SEMICOLON,
				token.IDENT,
				token.DEFINE,
				token.IDENT,
				token.SUB,
				token.IDENT,
				token.SEMICOLON,
				token.IDENT,
				token.DEFINE,
				token.IDENT,
				token.MUL,
				token.IDENT,
				token.SEMICOLON,
				token.IDENT,
				token.DEFINE,
				token.IDENT,
				token.DIV,
				token.IDENT,
				token.SEMICOLON,
				token.COMMENT,
				token.IDENT,
				token.DEFINE,
				token.IDENT,
				token.MOD,
				token.IDENT,
				token.SEMICOLON,
				token.IDENT,
				token.DEFINE,
				token.NOT,
				token.IDENT,
				token.SEMICOLON,
				token.IDENT,
				token.DEFINE,
				token.SUB,
				token.IDENT,
				token.SEMICOLON,
			},
		},
	}
)

func TestScan(t *testing.T) {
	s := new(Scanner)
	err := new(ErrorList)
	for _, c := range cases {
		err.Reset()
		s.Init("test", c.src, err.Add)
		j := 0
		gotErr := false
		for tok, lit := s.Scan(); tok != token.EOF; tok, lit = s.Scan() {
			if testing.Verbose() {
				fmt.Println(tok, lit)
			}
			if j < len(c.exp) && tok != c.exp[j] {
				t.Errorf("expected token %s at index %d, got %s", c.exp[j], j, tok)
				gotErr = true
				break
			} else if j >= len(c.exp) {
				t.Errorf("unexpected superfluous token at index %d", j)
				gotErr = true
				break
			}
			j++
		}
		if !gotErr && len(c.exp) > j {
			t.Errorf("missing %d token(s)", len(c.exp)-j)
			gotErr = true
		}
		if e := err.Err(); e != nil {
			t.Errorf("got scanning errors: %s", e)
		}
	}
}
