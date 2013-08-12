package tdop

import (
	"testing"

	"github.com/PuerkitoBio/goblin/compiler/scanner"
)

var (
	cases = []struct {
		src []byte
	}{
		0: {
			src: []byte(`return 5`),
		},
		1: {
			src: []byte(`aB := 5
return aB`),
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
		},
		3: {
			src: []byte(`
func Add(x, y) { // Essentially means var Add = func ...
  return x + y
}
return Add(4, "198")
`),
		},
		4: {
			src: []byte(`
Add := func(x, y) { // Essentially means var Add = func ...
  return x + y
}
return Add(4, "198")
`),
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
	}

	isolateCase = 14
)

func TestParse(t *testing.T) {
	Scanner = new(scanner.Scanner)
	for i, c := range cases {
		if isolateCase >= 0 && i != isolateCase {
			continue
		}

		Parse("test", c.src)
	}
}
