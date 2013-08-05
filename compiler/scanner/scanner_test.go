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
		0: {
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
		1: {
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
		2: {
			src: []byte(`
func Add(x, y) { // Essentially means var Add = func ...
  return x + y
}
return Add(4, "198")
`),
			exp: []token.Token{
				token.FUNC,
				token.IDENT,
				token.LPAREN,
				token.IDENT,
				token.COMMA,
				token.IDENT,
				token.RPAREN,
				token.LBRACE,
				token.COMMENT,
				token.RETURN,
				token.IDENT,
				token.ADD,
				token.IDENT,
				token.SEMICOLON,
				token.RBRACE,
				token.SEMICOLON,
				token.RETURN,
				token.IDENT,
				token.LPAREN,
				token.INT,
				token.COMMA,
				token.STRING,
				token.RPAREN,
				token.SEMICOLON,
			},
		},
		3: {
			src: []byte(`
func Fib(n) {
  if n < 2 {
    return 1
  }
  return Fib(n-1) + Fib(n-2)
}
return Fib(30)
`),
			exp: []token.Token{
				token.FUNC,
				token.IDENT,
				token.LPAREN,
				token.IDENT,
				token.RPAREN,
				token.LBRACE,
				token.IF,
				token.IDENT,
				token.LSS,
				token.INT,
				token.LBRACE,
				token.RETURN,
				token.INT,
				token.SEMICOLON,
				token.RBRACE,
				token.SEMICOLON,
				token.RETURN,
				token.IDENT,
				token.LPAREN,
				token.IDENT,
				token.SUB,
				token.INT,
				token.RPAREN,
				token.ADD,
				token.IDENT,
				token.LPAREN,
				token.IDENT,
				token.SUB,
				token.INT,
				token.RPAREN,
				token.SEMICOLON,
				token.RBRACE,
				token.SEMICOLON,
				token.RETURN,
				token.IDENT,
				token.LPAREN,
				token.INT,
				token.RPAREN,
				token.SEMICOLON,
			},
		},
		4: {
			src: []byte(`
import "fmt" // implicit fmt variable
fmt.Println("Hello ", "world")
`),
			exp: []token.Token{
				token.IMPORT,
				token.STRING,
				token.SEMICOLON,
				token.COMMENT,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.LPAREN,
				token.STRING,
				token.COMMA,
				token.STRING,
				token.RPAREN,
				token.SEMICOLON,
			},
		},
		5: {
			src: []byte(`
import "fmt"
// range over 0..9, there's implicit 0 (start) and 1 (step) constants
for i := range 10 { 
  fmt.Println(i)
}
`),
			exp: []token.Token{
				token.IMPORT,
				token.STRING,
				token.SEMICOLON,
				token.COMMENT,
				token.FOR,
				token.IDENT,
				token.DEFINE,
				token.RANGE,
				token.INT,
				token.LBRACE,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.LPAREN,
				token.IDENT,
				token.RPAREN,
				token.SEMICOLON,
				token.RBRACE,
				token.SEMICOLON,
			},
		},
		6: {
			src: []byte(`
a := 5
sum := 0
for a > 0 { 
  sum += a
  a-- // implicit constant 1
}
return sum
`),
			exp: []token.Token{
				token.IDENT,
				token.DEFINE,
				token.INT,
				token.SEMICOLON,
				token.IDENT,
				token.DEFINE,
				token.INT,
				token.SEMICOLON,
				token.FOR,
				token.IDENT,
				token.GTR,
				token.INT,
				token.LBRACE,
				token.IDENT,
				token.ADD_ASSIGN,
				token.IDENT,
				token.SEMICOLON,
				token.IDENT,
				token.DEC,
				token.SEMICOLON,
				token.COMMENT,
				token.RBRACE,
				token.SEMICOLON,
				token.RETURN,
				token.IDENT,
				token.SEMICOLON,
			},
		},
		7: {
			src: []byte(`
import "fmt"
a := "ok"
if a { 
  fmt.Println("true")
} else {
  fmt.Println("false")
}
`),
			exp: []token.Token{
				token.IMPORT,
				token.STRING,
				token.SEMICOLON,
				token.IDENT,
				token.DEFINE,
				token.STRING,
				token.SEMICOLON,
				token.IF,
				token.IDENT,
				token.LBRACE,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.LPAREN,
				token.STRING,
				token.RPAREN,
				token.SEMICOLON,
				token.RBRACE,
				token.ELSE,
				token.LBRACE,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.LPAREN,
				token.STRING,
				token.RPAREN,
				token.SEMICOLON,
				token.RBRACE,
				token.SEMICOLON,
			},
		},
		8: {
			src: []byte(`
a := true
if (3 + 5) > 4 && ("foo" > "bar") && a { 
  return 1
} else {
  return -1
}
`),
			exp: []token.Token{
				token.IDENT,
				token.DEFINE,
				token.TRUE,
				token.SEMICOLON,
				token.IF,
				token.LPAREN,
				token.INT,
				token.ADD,
				token.INT,
				token.RPAREN,
				token.GTR,
				token.INT,
				token.AND,
				token.LPAREN,
				token.STRING,
				token.GTR,
				token.STRING,
				token.RPAREN,
				token.AND,
				token.IDENT,
				token.LBRACE,
				token.RETURN,
				token.INT,
				token.SEMICOLON,
				token.RBRACE,
				token.ELSE,
				token.LBRACE,
				token.RETURN,
				token.SUB,
				token.INT,
				token.SEMICOLON,
				token.RBRACE,
				token.SEMICOLON,
			},
		},
		9: {
			src: []byte(`
a := {}
a.b = "6"
a.c = 4
a.d = a.b + a.c
return a.d
`),
			exp: []token.Token{
				token.IDENT,
				token.DEFINE,
				token.LBRACE,
				token.RBRACE,
				token.SEMICOLON,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.ASSIGN,
				token.STRING,
				token.SEMICOLON,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.ASSIGN,
				token.INT,
				token.SEMICOLON,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.ASSIGN,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.ADD,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.SEMICOLON,
				token.RETURN,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.SEMICOLON,
			},
		},
		10: {
			src: []byte(`
import "fmt"
a := {}
a.b = func(greet) {
  fmt.Println(this.c)
  return this.c + ", " + greet
}
a.c = "hi"
return a.b("you")
`),
			exp: []token.Token{
				token.IMPORT,
				token.STRING,
				token.SEMICOLON,
				token.IDENT,
				token.DEFINE,
				token.LBRACE,
				token.RBRACE,
				token.SEMICOLON,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.ASSIGN,
				token.FUNC,
				token.LPAREN,
				token.IDENT,
				token.RPAREN,
				token.LBRACE,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.LPAREN,
				token.THIS,
				token.PERIOD,
				token.IDENT,
				token.RPAREN,
				token.SEMICOLON,
				token.RETURN,
				token.THIS,
				token.PERIOD,
				token.IDENT,
				token.ADD,
				token.STRING,
				token.ADD,
				token.IDENT,
				token.SEMICOLON,
				token.RBRACE,
				token.SEMICOLON,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.ASSIGN,
				token.STRING,
				token.SEMICOLON,
				token.RETURN,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.LPAREN,
				token.STRING,
				token.RPAREN,
				token.SEMICOLON,
			},
		},
		11: {
			src: []byte(`
import "fmt"
a := {}
a.__noSuchMethod = func(nm) {
  fmt.Println("not found:" + nm)
}
a.b(12)
`),
			exp: []token.Token{
				token.IMPORT,
				token.STRING,
				token.SEMICOLON,
				token.IDENT,
				token.DEFINE,
				token.LBRACE,
				token.RBRACE,
				token.SEMICOLON,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.ASSIGN,
				token.FUNC,
				token.LPAREN,
				token.IDENT,
				token.RPAREN,
				token.LBRACE,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.LPAREN,
				token.STRING,
				token.ADD,
				token.IDENT,
				token.RPAREN,
				token.SEMICOLON,
				token.RBRACE,
				token.SEMICOLON,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.LPAREN,
				token.INT,
				token.RPAREN,
				token.SEMICOLON,
			},
		},
		12: {
			src: []byte(`
a := 5
func b(delta) {
  a += delta
}
b(3)
return a
`),
			exp: []token.Token{
				token.IDENT,
				token.DEFINE,
				token.INT,
				token.SEMICOLON,
				token.FUNC,
				token.IDENT,
				token.LPAREN,
				token.IDENT,
				token.RPAREN,
				token.LBRACE,
				token.IDENT,
				token.ADD_ASSIGN,
				token.IDENT,
				token.SEMICOLON,
				token.RBRACE,
				token.SEMICOLON,
				token.IDENT,
				token.LPAREN,
				token.INT,
				token.RPAREN,
				token.SEMICOLON,
				token.RETURN,
				token.IDENT,
				token.SEMICOLON,
			},
		},
		13: {
			src: []byte(`
import "fmt"
func f() {
  fmt.Println(args[0], args[1], args[2]) // Compiler translates this to PUSH AA ix, no need for Ks
}
f(17, "foo", false)
`),
			exp: []token.Token{
				token.IMPORT,
				token.STRING,
				token.SEMICOLON,
				token.FUNC,
				token.IDENT,
				token.LPAREN,
				token.RPAREN,
				token.LBRACE,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.LPAREN,
				token.IDENT,
				token.LBRACK,
				token.INT,
				token.RBRACK,
				token.COMMA,
				token.IDENT,
				token.LBRACK,
				token.INT,
				token.RBRACK,
				token.COMMA,
				token.IDENT,
				token.LBRACK,
				token.INT,
				token.RBRACK,
				token.RPAREN,
				token.SEMICOLON,
				token.COMMENT,
				token.RBRACE,
				token.SEMICOLON,
				token.IDENT,
				token.LPAREN,
				token.INT,
				token.COMMA,
				token.STRING,
				token.COMMA,
				token.FALSE,
				token.RPAREN,
				token.SEMICOLON,
			},
		},
		14: {
			src: []byte(`
a := {b: {c: {d: "allo"}}}
return a.b.c.d
`),
			exp: []token.Token{
				token.IDENT,
				token.DEFINE,
				token.LBRACE,
				token.IDENT,
				token.COLON,
				token.LBRACE,
				token.IDENT,
				token.COLON,
				token.LBRACE,
				token.IDENT,
				token.COLON,
				token.STRING,
				token.RBRACE,
				token.RBRACE,
				token.RBRACE,
				token.SEMICOLON,
				token.RETURN,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.PERIOD,
				token.IDENT,
				token.SEMICOLON,
			},
		},
		15: {
			src: []byte{},
			exp: []token.Token{},
		},
		16: {
			src: []byte(`
						Ceci est du GROS n'importe quoi@!
`),
			exp: []token.Token{
				token.IDENT,
				token.IDENT,
				token.IDENT,
				token.IDENT,
				token.IDENT,
				token.ILLEGAL,
				token.IDENT,
				token.IDENT,
				token.ILLEGAL,
				token.NOT,
			},
		},
	}

	isolateCase = -1
)

func TestScan(t *testing.T) {
	s := new(Scanner)
	err := new(ErrorList)
	for i, c := range cases {
		if isolateCase >= 0 && i != isolateCase {
			continue
		}

		err.Reset()
		s.Init("test", c.src, err.Add)
		j := 0
		gotErr := false
		for tok, lit := s.Scan(); tok != token.EOF; tok, lit = s.Scan() {
			if isolateCase >= 0 && testing.Verbose() {
				// Print the results
				fmt.Println(tok, lit)
			}

			if j < len(c.exp) && tok != c.exp[j] {
				t.Errorf("case %d - expected token %s at index %d, got %s", i, c.exp[j], j, tok)
				gotErr = true
				break
			} else if j >= len(c.exp) {
				t.Errorf("case %d - unexpected superfluous token at index %d", i, j)
				gotErr = true
				break
			}
			j++
		}
		if !gotErr && len(c.exp) > j {
			t.Errorf("case %d - missing %d token(s)", i, len(c.exp)-j)
			gotErr = true
		}
		if e := err.Err(); e != nil && testing.Verbose() {
			fmt.Printf("case %d - got scanning errors: %s\n", i, e)
			//t.Errorf("case %d - got scanning errors: %s", i, e)
		}
	}
}
