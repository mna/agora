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
				&Symbol{Id: "return"},
				&Symbol{Id: "(literal)", Val: "5"},
			},
		},
		1: {
			src: []byte(`aB := 5
return aB`),
			exp: []*Symbol{
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "aB"},
				&Symbol{Id: "(literal)", Val: "5"},
				&Symbol{Id: "return"},
				&Symbol{Id: "(name)", Val: "aB"},
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
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(literal)", Val: "7"},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: "(literal)", Val: "10"},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "add"},
				&Symbol{Id: "+"},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "sub"},
				&Symbol{Id: "-"},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "mul"},
				&Symbol{Id: "*"},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "div"},
				&Symbol{Id: "/"},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "mod"},
				&Symbol{Id: "%"},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "not"},
				&Symbol{Id: "!"},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "unm"},
				&Symbol{Id: "-"},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "return"},
				&Symbol{Id: "nil"},
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
				&Symbol{Id: "func", Name: "Add"},
				&Symbol{Id: "(name)", Val: "x"},
				&Symbol{Id: "(name)", Val: "y"},
				&Symbol{Id: "return"},
				&Symbol{Id: "+"},
				&Symbol{Id: "(name)", Val: "x"},
				&Symbol{Id: "(name)", Val: "y"},
				&Symbol{Id: "return"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "Add"},
				&Symbol{Id: "(literal)", Val: "4"},
				&Symbol{Id: "(literal)", Val: `"198"`},
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
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "Add"},
				&Symbol{Id: "func"},
				&Symbol{Id: "(name)", Val: "x"},
				&Symbol{Id: "(name)", Val: "y"},
				&Symbol{Id: "return"},
				&Symbol{Id: "+"},
				&Symbol{Id: "(name)", Val: "x"},
				&Symbol{Id: "(name)", Val: "y"},
				&Symbol{Id: "return"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "Add"},
				&Symbol{Id: "(literal)", Val: "4"},
				&Symbol{Id: "(literal)", Val: `"198"`},
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
				&Symbol{Id: "func", Name: "Fib"},
				&Symbol{Id: "(name)", Val: "n"},
				&Symbol{Id: "if"},
				&Symbol{Id: "<"},
				&Symbol{Id: "(name)", Val: "n"},
				&Symbol{Id: "(literal)", Val: "2"},
				&Symbol{Id: "return"},
				&Symbol{Id: "(literal)", Val: "1"},
				&Symbol{Id: "return"},
				&Symbol{Id: "+"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "Fib"},
				&Symbol{Id: "-"},
				&Symbol{Id: "(name)", Val: "n"},
				&Symbol{Id: "(literal)", Val: "1"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "Fib"},
				&Symbol{Id: "-"},
				&Symbol{Id: "(name)", Val: "n"},
				&Symbol{Id: "(literal)", Val: "2"},
				&Symbol{Id: "return"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "Fib"},
				&Symbol{Id: "(literal)", Val: "30"},
			},
		},
		6: {
			src: []byte(`
fmt := import("fmt") // implicit fmt variable
fmt.Println("Hello ", "world")
`),
			exp: []*Symbol{
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "fmt"},
				&Symbol{Id: "("},
				&Symbol{Id: "import", Ar: ArName, Val: "import"},
				&Symbol{Id: "(literal)", Val: `"fmt"`},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "fmt"},
				&Symbol{Id: "(name)", Val: "Println"},
				&Symbol{Id: "(literal)", Val: `"Hello "`},
				&Symbol{Id: "(literal)", Val: `"world"`},
				&Symbol{Id: "return"},
				&Symbol{Id: "nil"},
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
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(literal)", Val: "5"},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "sum"},
				&Symbol{Id: "(literal)", Val: "0"},
				&Symbol{Id: "for"},
				&Symbol{Id: ">"},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(literal)", Val: "0"},
				&Symbol{Id: "+="},
				&Symbol{Id: "(name)", Val: "sum"},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "--"},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "return"},
				&Symbol{Id: "(name)", Val: "sum"},
			},
		},
		8: {
			src: []byte(`
f := import("fmt")
a := "ok"
if a { 
  f.Println("true")
} else {
  f.Println("false")
}
`),
			exp: []*Symbol{
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "f"},
				&Symbol{Id: "("},
				&Symbol{Id: "import", Ar: ArName, Val: "import"},
				&Symbol{Id: "(literal)", Val: `"fmt"`},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(literal)", Val: `"ok"`},
				&Symbol{Id: "if"},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "f"},
				&Symbol{Id: "(name)", Val: "Println"},
				&Symbol{Id: "(literal)", Val: `"true"`},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "f"},
				&Symbol{Id: "(name)", Val: "Println"},
				&Symbol{Id: "(literal)", Val: `"false"`},
				&Symbol{Id: "return"},
				&Symbol{Id: "nil"},
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
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "true"},
				&Symbol{Id: "if"},
				&Symbol{Id: "&&"},
				&Symbol{Id: ">"},
				&Symbol{Id: "+"},
				&Symbol{Id: "(literal)", Val: "3"},
				&Symbol{Id: "(literal)", Val: "5"},
				&Symbol{Id: "(literal)", Val: "4"},
				&Symbol{Id: "&&"},
				&Symbol{Id: ">"},
				&Symbol{Id: "(literal)", Val: `"foo"`},
				&Symbol{Id: "(literal)", Val: `"bar"`},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "return"},
				&Symbol{Id: "(literal)", Val: "1"},
				&Symbol{Id: "return"},
				&Symbol{Id: "-"},
				&Symbol{Id: "(literal)", Val: "1"},
				&Symbol{Id: "return"},
				&Symbol{Id: "nil"},
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
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "{"},
				&Symbol{Id: "="},
				&Symbol{Id: "."},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: "(literal)", Val: `"6"`},
				&Symbol{Id: "="},
				&Symbol{Id: "."},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "c"},
				&Symbol{Id: "(literal)", Val: "4"},
				&Symbol{Id: "="},
				&Symbol{Id: "."},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "d"},
				&Symbol{Id: "+"},
				&Symbol{Id: "."},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: "."},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "c"},
				&Symbol{Id: "return"},
				&Symbol{Id: "."},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "d"},
			},
		},
		11: {
			src: []byte(`
fmt := import("fmt")
a := {}
a.b = func(greet) {
  fmt.Println(this.c)
  return this.c + ", " + greet
}
a.c = "hi"
return a.b("you")
`),
			exp: []*Symbol{
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "fmt"},
				&Symbol{Id: "("},
				&Symbol{Id: "import", Ar: ArName, Val: "import"},
				&Symbol{Id: "(literal)", Val: `"fmt"`},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "{"},
				&Symbol{Id: "="},
				&Symbol{Id: "."},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: "func"},
				&Symbol{Id: "(name)", Val: "greet"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "fmt"},
				&Symbol{Id: "(name)", Val: "Println"},
				&Symbol{Id: "."},
				&Symbol{Id: "this"},
				&Symbol{Id: "(name)", Val: "c"},
				&Symbol{Id: "return"},
				&Symbol{Id: "+"},
				&Symbol{Id: "+"},
				&Symbol{Id: "."},
				&Symbol{Id: "this"},
				&Symbol{Id: "(name)", Val: "c"},
				&Symbol{Id: "(literal)", Val: `", "`},
				&Symbol{Id: "(name)", Val: "greet"},
				&Symbol{Id: "="},
				&Symbol{Id: "."},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "c"},
				&Symbol{Id: "(literal)", Val: `"hi"`},
				&Symbol{Id: "return"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: "(literal)", Val: `"you"`},
			},
		},
		12: {
			src: []byte(`
			fmt := import("fmt")
a := {}
a.__noSuchMethod = func(nm) {
  fmt.Println("not found:" + nm)
}
a.b(12)
`),
			exp: []*Symbol{
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "fmt"},
				&Symbol{Id: "("},
				&Symbol{Id: "import", Ar: ArName, Val: "import"},
				&Symbol{Id: "(literal)", Val: `"fmt"`},
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "{"},
				&Symbol{Id: "="},
				&Symbol{Id: "."},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "__noSuchMethod"},
				&Symbol{Id: "func"},
				&Symbol{Id: "(name)", Val: "nm"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "fmt"},
				&Symbol{Id: "(name)", Val: "Println"},
				&Symbol{Id: "+"},
				&Symbol{Id: "(literal)", Val: `"not found:"`},
				&Symbol{Id: "(name)", Val: "nm"},
				&Symbol{Id: "return"},
				&Symbol{Id: "nil"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: "(literal)", Val: "12"},
				&Symbol{Id: "return"},
				&Symbol{Id: "nil"},
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
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(literal)", Val: "5"},
				&Symbol{Id: "func", Name: "b"},
				&Symbol{Id: "(name)", Val: "delta"},
				&Symbol{Id: "+="},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "delta"},
				&Symbol{Id: "return"},
				&Symbol{Id: "nil"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: "(literal)", Val: "3"},
				&Symbol{Id: "return"},
				&Symbol{Id: "(name)", Val: "a"},
			},
		},
		14: {
			src: []byte(`
			fmt := import("fmt")
func f() {
  fmt.Println(args[0], args[1], args[2]) // Compiler translates this to PUSH AA ix, no need for Ks
}
f(17, "foo", false)
`),
			exp: []*Symbol{
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "fmt"},
				&Symbol{Id: "("},
				&Symbol{Id: "import", Ar: ArName, Val: "import"},
				&Symbol{Id: "(literal)", Val: `"fmt"`},
				&Symbol{Id: "func", Name: "f"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "fmt"},
				&Symbol{Id: "(name)", Val: "Println"},
				&Symbol{Id: "["},
				&Symbol{Id: "args"},
				&Symbol{Id: "(literal)", Val: "0"},
				&Symbol{Id: "["},
				&Symbol{Id: "args"},
				&Symbol{Id: "(literal)", Val: "1"},
				&Symbol{Id: "["},
				&Symbol{Id: "args"},
				&Symbol{Id: "(literal)", Val: "2"},
				&Symbol{Id: "return"},
				&Symbol{Id: "nil"},
				&Symbol{Id: "("},
				&Symbol{Id: "(name)", Val: "f"},
				&Symbol{Id: "(literal)", Val: "17"},
				&Symbol{Id: "(literal)", Val: `"foo"`},
				&Symbol{Id: "false"},
				&Symbol{Id: "return"},
				&Symbol{Id: "nil"},
			},
		},
		15: {
			src: []byte(`
a := {b: {c: {d: "allo"}}}
return a.b.c.d
`),
			exp: []*Symbol{
				&Symbol{Id: ":="},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "{", Key: ""},
				&Symbol{Id: "{", Key: "b"},
				&Symbol{Id: "{", Key: "c"},
				&Symbol{Id: "(literal)", Val: `"allo"`, Key: "d"},
				&Symbol{Id: "return"},
				&Symbol{Id: "."},
				&Symbol{Id: "."},
				&Symbol{Id: "."},
				&Symbol{Id: "(name)", Val: "a"},
				&Symbol{Id: "(name)", Val: "b"},
				&Symbol{Id: "(name)", Val: "c"},
				&Symbol{Id: "(name)", Val: "d"},
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
				&Symbol{Id: "if"},
				&Symbol{Id: "true"},
				&Symbol{Id: "return"},
				&Symbol{Id: "(literal)", Val: "1"},
				&Symbol{Id: "if"},
				&Symbol{Id: "false"},
				&Symbol{Id: "return"},
				&Symbol{Id: "(literal)", Val: "2"},
				&Symbol{Id: "return"},
				&Symbol{Id: "(literal)", Val: "3"},
				&Symbol{Id: "return"},
				&Symbol{Id: "nil"},
			},
		},
		17: {
			src: []byte(`
			Ceci est un gros n'importe quoi! @ # 234r3.112@O#Ihwev92h f9238f
`),
			err: true,
		},
	}

	isolateCase = -1
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
					if v.Id != c.exp[ix].Id {
						t.Errorf("[%d] - expected symbol id %s, got %s", i, c.exp[ix].Id, v.Id)
					}
					if c.exp[ix].Val != nil && v.Val != c.exp[ix].Val {
						t.Errorf("[%d] - expected symbol value %v, got %v", i, c.exp[ix].Val, v.Val)
					}
					if c.exp[ix].Name != "" && v.Name != c.exp[ix].Name {
						t.Errorf("[%d] - expected symbol name %s, got %s", i, c.exp[ix].Name, v.Name)
					}
					if c.exp[ix].Key != "" && v.Key != c.exp[ix].Key {
						t.Errorf("[%d] - expected symbol key %s, got %s", i, c.exp[ix].Key, v.Key)
					}
				}
				check(v.First)
				check(v.Second)
				check(v.Third)
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
