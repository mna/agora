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
	}

	isolateCase = 2
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
