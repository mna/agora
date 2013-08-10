// +build ignore

package parser

import (
	"fmt"
	"testing"
)

var (
	cases = []struct {
		nm  string
		src []byte
		exp string
		err bool
	}{
		0: {
			nm: "single-with-ident",
			src: []byte(`
import abc "this/is/a/path"
`),
			exp: "module single-with-ident\n  \"this/is/a/path\" (abc)\n",
		},
		1: {
			nm: "single-without-ident",
			src: []byte(`
import "this/is/a/path"
`),
			exp: "module single-without-ident\n  \"this/is/a/path\" (path)\n",
		},
		2: {
			nm: "single-empty",
			src: []byte(`
import
`),
			exp: "module single-empty\n",
			err: true,
		},
		3: {
			nm: "multi-empty",
			src: []byte(`
import ()
`),
			exp: "module multi-empty\n",
		},
		4: {
			nm: "multi-one-noident",
			src: []byte(`
import (
	"this/is/a/path"
)
`),
			exp: "module multi-one-noident\n  \"this/is/a/path\" (path)\n",
		},
		5: {
			nm: "multi-one-ident",
			src: []byte(`
import (
	xyz   	"this/is/a/path"
)
`),
			exp: "module multi-one-ident\n  \"this/is/a/path\" (xyz)\n",
		},
		6: {
			nm: "multi-many",
			src: []byte(`
import (
	"path/one"
	xyz   	"this/is/a/path"
	"path/two"
)
`),
			exp: "module multi-many\n  \"path/one\" (one)\n  \"this/is/a/path\" (xyz)\n  \"path/two\" (two)\n",
		},
	}
)

func TestParser(t *testing.T) {
	p := new(Parser)
	for _, c := range cases {
		if testing.Verbose() {
			fmt.Printf("testing %s...\n", c.nm)
		}
		a, err := p.Parse(c.nm, c.src)
		if a.String() != c.exp {
			t.Errorf("expected %s, got %s", c.exp, a)
		}
		if err != nil {
			if !c.err {
				t.Errorf("expected no error, got %s", err)
			}
			if testing.Verbose() {
				fmt.Println(err)
			}
		}
	}
}
