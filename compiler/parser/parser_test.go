package parser

import (
	"testing"

	"github.com/PuerkitoBio/goblin/compiler/scanner"
)

var (
	cases = []struct {
		src []byte
		exp []*Symbol
	}{
		0: {
			src: []byte(`return 5`),
			exp: []*Symbol{
				&Symbol{id: "return"},
				&Symbol{id: "(literal)", val: "5"},
			},
		},
		1: {
			src: []byte(`aB := 5
return aB`),
			exp: []*Symbol{
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "aB"},
				&Symbol{id: "(literal)", val: "5"},
				&Symbol{id: "return"},
				&Symbol{id: "(name)", val: "aB"},
			},
		},
		2: {
			src: []byte(`
a := 7
b := 10
add := a + b
sub := a - b
mul := a * b
div := a / b // TODO : Should div return a float even for Int?
mod := b % a
not := !a
unm := -a
`),
			exp: []*Symbol{
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(literal)", val: "7"},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: "(literal)", val: "10"},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "add"},
				&Symbol{id: "+"},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "sub"},
				&Symbol{id: "-"},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "mul"},
				&Symbol{id: "*"},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "div"},
				&Symbol{id: "/"},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "mod"},
				&Symbol{id: "%"},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "not"},
				&Symbol{id: "!"},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "unm"},
				&Symbol{id: "-"},
				&Symbol{id: "(name)", val: "a"},
			},
		},
		3: {
			src: []byte(`
func Add(x, y) { // Essentially means var Add = func ...
  return x + y
}
return Add(4, "198")
`),
			exp: []*Symbol{
				&Symbol{id: "func", name: "Add"},
				&Symbol{id: "(name)", val: "x"},
				&Symbol{id: "(name)", val: "y"},
				&Symbol{id: "return"},
				&Symbol{id: "+"},
				&Symbol{id: "(name)", val: "x"},
				&Symbol{id: "(name)", val: "y"},
				&Symbol{id: "return"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "Add"},
				&Symbol{id: "(literal)", val: "4"},
				&Symbol{id: "(literal)", val: `"198"`},
			},
		},
		4: {
			src: []byte(`
Add := func(x, y) { // Essentially means var Add = func ...
  return x + y
}
return Add(4, "198")
`),
			exp: []*Symbol{
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "Add"},
				&Symbol{id: "func"},
				&Symbol{id: "(name)", val: "x"},
				&Symbol{id: "(name)", val: "y"},
				&Symbol{id: "return"},
				&Symbol{id: "+"},
				&Symbol{id: "(name)", val: "x"},
				&Symbol{id: "(name)", val: "y"},
				&Symbol{id: "return"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "Add"},
				&Symbol{id: "(literal)", val: "4"},
				&Symbol{id: "(literal)", val: `"198"`},
			},
		},
		5: {
			src: []byte(`
func Fib(n) {
  if n < 2 {
    return 1
  }
  return Fib(n-1) + Fib(n-2)
}
return Fib(30)
`),
		},
		6: {
			src: []byte(`
import "fmt" // implicit fmt variable
fmt.Println("Hello ", "world")
`),
		},
		7: {
			src: []byte(`
a := 5
sum := 0
for a > 0 { 
  sum += a
  a-- // implicit constant 1
}
return sum
`),
		},
		8: {
			src: []byte(`
import "fmt"
a := "ok"
if a { 
  fmt.Println("true")
} else {
  fmt.Println("false")
}
`),
		},
		9: {
			src: []byte(`
a := true
if (3 + 5) > 4 && ("foo" > "bar") && a { 
  return 1
} else {
  return -1
}
`),
		},
		10: {
			src: []byte(`
a := {}
a.b = "6"
a.c = 4
a.d = a.b + a.c
return a.d
`),
		},
		11: {
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
		},
		12: {
			src: []byte(`
import "fmt"
a := {}
a.__noSuchMethod = func(nm) {
  fmt.Println("not found:" + nm)
}
a.b(12)
`),
		},
		13: {
			src: []byte(`
a := 5
func b(delta) {
  a += delta
}
b(3)
return a
`),
		},
		14: {
			src: []byte(`
import "fmt"
func f() {
  fmt.Println(args[0], args[1], args[2]) // Compiler translates this to PUSH AA ix, no need for Ks
}
f(17, "foo", false)
`),
		},
		15: {
			src: []byte(`
a := {b: {c: {d: "allo"}}}
return a.b.c.d
`),
		},
		16: {
			src: []byte(`
if true {
	return 1
} else if false {
	return 2
} else {
	return 3
}
`),
		},
	}

	isolateCase = 5
)

func TestParse(t *testing.T) {
	Scanner = new(scanner.Scanner)
	for i, c := range cases {
		if isolateCase >= 0 && i != isolateCase {
			continue
		}

		s := Parse("test", c.src)
		ix := -1
		var check func(interface{})
		check = func(root interface{}) {
			switch v := root.(type) {
			case *Symbol:
				ix++
				if v.id != c.exp[ix].id {
					t.Errorf("[%d] - expected symbol id %s, got %s", i, c.exp[ix].id, v.id)
				}
				if c.exp[ix].val != nil && v.val != c.exp[ix].val {
					t.Errorf("[%d] - expected symbol value %s, got %s", i, c.exp[ix].val, v.val)
				}
				if c.exp[ix].name != "" && v.name != c.exp[ix].name {
					t.Errorf("[%d] - expected symbol name %s, got %s", i, c.exp[ix].name, v.name)
				}
				check(v.first)
				check(v.second)
				check(v.third)
			case []*Symbol:
				for _, s := range v {
					check(s)
				}
			case nil:
			default:
				panic("unkown type")
			}
		}
		if c.exp != nil {
			check(s)
			if len(c.exp) != (ix + 1) {
				t.Errorf("[%d] - expected %d symbols, got %d", i, len(c.exp), ix+1)
			}
		}
	}
}
