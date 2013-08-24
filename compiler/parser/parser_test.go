package parser

import (
	"fmt"
	"testing"
)

var (
	cases = []struct {
		src []byte
		exp []*Symbol
		err bool
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
			exp: []*Symbol{
				&Symbol{id: "func", name: "Fib"},
				&Symbol{id: "(name)", val: "n"},
				&Symbol{id: "if"},
				&Symbol{id: "<"},
				&Symbol{id: "(name)", val: "n"},
				&Symbol{id: "(literal)", val: "2"},
				&Symbol{id: "return"},
				&Symbol{id: "(literal)", val: "1"},
				&Symbol{id: "return"},
				&Symbol{id: "+"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "Fib"},
				&Symbol{id: "-"},
				&Symbol{id: "(name)", val: "n"},
				&Symbol{id: "(literal)", val: "1"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "Fib"},
				&Symbol{id: "-"},
				&Symbol{id: "(name)", val: "n"},
				&Symbol{id: "(literal)", val: "2"},
				&Symbol{id: "return"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "Fib"},
				&Symbol{id: "(literal)", val: "30"},
			},
		},
		6: {
			src: []byte(`
import "fmt" // implicit fmt variable
fmt.Println("Hello ", "world")
`),
			exp: []*Symbol{
				&Symbol{id: "import"},
				&Symbol{id: "(name)", val: "fmt"},
				&Symbol{id: "(literal)", val: `"fmt"`},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "fmt"},
				&Symbol{id: "(name)", val: "Println"},
				&Symbol{id: "(literal)", val: `"Hello "`},
				&Symbol{id: "(literal)", val: `"world"`},
			},
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
			exp: []*Symbol{
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(literal)", val: "5"},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "sum"},
				&Symbol{id: "(literal)", val: "0"},
				&Symbol{id: "for"},
				&Symbol{id: ">"},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(literal)", val: "0"},
				&Symbol{id: "+="},
				&Symbol{id: "(name)", val: "sum"},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "--"},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "return"},
				&Symbol{id: "(name)", val: "sum"},
			},
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
			exp: []*Symbol{
				&Symbol{id: "import"},
				&Symbol{id: "(name)", val: "fmt"},
				&Symbol{id: "(literal)", val: `"fmt"`},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(literal)", val: `"ok"`},
				&Symbol{id: "if"},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "fmt"},
				&Symbol{id: "(name)", val: "Println"},
				&Symbol{id: "(literal)", val: `"true"`},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "fmt"},
				&Symbol{id: "(name)", val: "Println"},
				&Symbol{id: "(literal)", val: `"false"`},
			},
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
			exp: []*Symbol{
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "true"},
				&Symbol{id: "if"},
				&Symbol{id: "&&"},
				&Symbol{id: ">"},
				&Symbol{id: "+"},
				&Symbol{id: "(literal)", val: "3"},
				&Symbol{id: "(literal)", val: "5"},
				&Symbol{id: "(literal)", val: "4"},
				&Symbol{id: "&&"},
				&Symbol{id: ">"},
				&Symbol{id: "(literal)", val: `"foo"`},
				&Symbol{id: "(literal)", val: `"bar"`},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "return"},
				&Symbol{id: "(literal)", val: "1"},
				&Symbol{id: "return"},
				&Symbol{id: "-"},
				&Symbol{id: "(literal)", val: "1"},
			},
		},
		10: {
			src: []byte(`
a := {}
a.b = "6"
a.c = 4
a.d = a.b + a.c
return a.d
`),
			exp: []*Symbol{
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "{"},
				&Symbol{id: "="},
				&Symbol{id: "."},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: "(literal)", val: `"6"`},
				&Symbol{id: "="},
				&Symbol{id: "."},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "c"},
				&Symbol{id: "(literal)", val: "4"},
				&Symbol{id: "="},
				&Symbol{id: "."},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "d"},
				&Symbol{id: "+"},
				&Symbol{id: "."},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: "."},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "c"},
				&Symbol{id: "return"},
				&Symbol{id: "."},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "d"},
			},
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
			exp: []*Symbol{
				&Symbol{id: "import"},
				&Symbol{id: "(name)", val: "fmt"},
				&Symbol{id: "(literal)", val: `"fmt"`},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "{"},
				&Symbol{id: "="},
				&Symbol{id: "."},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: "func"},
				&Symbol{id: "(name)", val: "greet"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "fmt"},
				&Symbol{id: "(name)", val: "Println"},
				&Symbol{id: "."},
				&Symbol{id: "this"},
				&Symbol{id: "(name)", val: "c"},
				&Symbol{id: "return"},
				&Symbol{id: "+"},
				&Symbol{id: "+"},
				&Symbol{id: "."},
				&Symbol{id: "this"},
				&Symbol{id: "(name)", val: "c"},
				&Symbol{id: "(literal)", val: `", "`},
				&Symbol{id: "(name)", val: "greet"},
				&Symbol{id: "="},
				&Symbol{id: "."},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "c"},
				&Symbol{id: "(literal)", val: `"hi"`},
				&Symbol{id: "return"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: "(literal)", val: `"you"`},
			},
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
			exp: []*Symbol{
				&Symbol{id: "import"},
				&Symbol{id: "(name)", val: "fmt"},
				&Symbol{id: "(literal)", val: `"fmt"`},
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "{"},
				&Symbol{id: "="},
				&Symbol{id: "."},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "__noSuchMethod"},
				&Symbol{id: "func"},
				&Symbol{id: "(name)", val: "nm"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "fmt"},
				&Symbol{id: "(name)", val: "Println"},
				&Symbol{id: "+"},
				&Symbol{id: "(literal)", val: `"not found:"`},
				&Symbol{id: "(name)", val: "nm"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: "(literal)", val: "12"},
			},
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
			exp: []*Symbol{
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(literal)", val: "5"},
				&Symbol{id: "func", name: "b"},
				&Symbol{id: "(name)", val: "delta"},
				&Symbol{id: "+="},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "delta"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: "(literal)", val: "3"},
				&Symbol{id: "return"},
				&Symbol{id: "(name)", val: "a"},
			},
		},
		14: {
			src: []byte(`
import "fmt"
func f() {
  fmt.Println(args[0], args[1], args[2]) // Compiler translates this to PUSH AA ix, no need for Ks
}
f(17, "foo", false)
`),
			exp: []*Symbol{
				&Symbol{id: "import"},
				&Symbol{id: "(name)", val: "fmt"},
				&Symbol{id: "(literal)", val: `"fmt"`},
				&Symbol{id: "func", name: "f"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "fmt"},
				&Symbol{id: "(name)", val: "Println"},
				&Symbol{id: "["},
				&Symbol{id: "args"},
				&Symbol{id: "(literal)", val: "0"},
				&Symbol{id: "["},
				&Symbol{id: "args"},
				&Symbol{id: "(literal)", val: "1"},
				&Symbol{id: "["},
				&Symbol{id: "args"},
				&Symbol{id: "(literal)", val: "2"},
				&Symbol{id: "("},
				&Symbol{id: "(name)", val: "f"},
				&Symbol{id: "(literal)", val: "17"},
				&Symbol{id: "(literal)", val: `"foo"`},
				&Symbol{id: "false"},
			},
		},
		15: {
			src: []byte(`
a := {b: {c: {d: "allo"}}}
return a.b.c.d
`),
			exp: []*Symbol{
				&Symbol{id: ":="},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "{", key: ""},
				&Symbol{id: "{", key: "b"},
				&Symbol{id: "{", key: "c"},
				&Symbol{id: "(literal)", val: `"allo"`, key: "d"},
				&Symbol{id: "return"},
				&Symbol{id: "."},
				&Symbol{id: "."},
				&Symbol{id: "."},
				&Symbol{id: "(name)", val: "a"},
				&Symbol{id: "(name)", val: "b"},
				&Symbol{id: "(name)", val: "c"},
				&Symbol{id: "(name)", val: "d"},
			},
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
			exp: []*Symbol{
				&Symbol{id: "if"},
				&Symbol{id: "true"},
				&Symbol{id: "return"},
				&Symbol{id: "(literal)", val: "1"},
				&Symbol{id: "if"},
				&Symbol{id: "false"},
				&Symbol{id: "return"},
				&Symbol{id: "(literal)", val: "2"},
				&Symbol{id: "return"},
				&Symbol{id: "(literal)", val: "3"},
			},
		},
		17: {
			src: []byte(`
			Ceci est un gros n'importe quoi! @ # 234r3.112@O#Ihwev92h f9238f
`),
			err: true,
		},
	}

	isolateCase = 17
)

func TestParse(t *testing.T) {
	p := New()
	if isolateCase >= 0 && testing.Verbose() {
		p.Debug = true
	}
	for i, c := range cases {
		if isolateCase >= 0 && i != isolateCase {
			continue
		}
		if testing.Verbose() {
			fmt.Printf("testing parser case %d...\n", i)
		}

		syms, _, err := p.Parse("test", c.src)
		ix := -1
		var check func(interface{})
		check = func(root interface{}) {
			switch v := root.(type) {
			case *Symbol:
				ix++
				if ix < len(c.exp) {
					if v.id != c.exp[ix].id {
						t.Errorf("[%d] - expected symbol id %s, got %s", i, c.exp[ix].id, v.id)
					}
					if c.exp[ix].val != nil && v.val != c.exp[ix].val {
						t.Errorf("[%d] - expected symbol value %v, got %v", i, c.exp[ix].val, v.val)
					}
					if c.exp[ix].name != "" && v.name != c.exp[ix].name {
						t.Errorf("[%d] - expected symbol name %s, got %s", i, c.exp[ix].name, v.name)
					}
					if c.exp[ix].key != "" && v.key != c.exp[ix].key {
						t.Errorf("[%d] - expected symbol key %s, got %s", i, c.exp[ix].key, v.key)
					}
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
			check(syms)
			if len(c.exp) != (ix + 1) {
				t.Errorf("[%d] - expected %d symbols, got %d", i, len(c.exp), ix+1)
			}
		}
		if (err != nil) != c.err {
			if c.err {
				t.Errorf("[%d] - expected error(s), got none", i)
			} else {
				t.Errorf("[%d] - expected no error, got %s", i, err)
			}
		}
		if c.exp == nil && !c.err {
			t.Errorf("[%d] - no assertion", i)
		}
	}
}
