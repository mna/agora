package compiler

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var (
	asmcases = []struct {
		id  string
		src string
		exp []byte
		err error
	}{
		0: {
			// No input
			id:  "test",
			err: ErrNoInput,
		},
		1: {
			// Empty func
			id: "test",
			src: `
// Some comment

[f]
`,
			exp: []byte{},
		},
	}

	isolateAsmCase = -1
)

func TestAsm(t *testing.T) {
	a := new(Asm)
	for i, c := range asmcases {
		if isolateAsmCase >= 0 && isolateAsmCase != i {
			continue
		}
		if testing.Verbose() {
			fmt.Printf("testing assembler case %d...\n", i)
		}

		// Act
		got, err := a.Compile(c.id, strings.NewReader(c.src))

		// Assert
		if err != c.err {
			if c.err == nil {
				t.Errorf("[%d] - expected no error, got `%s`", i, err)
			} else {
				t.Errorf("[%d] - expected error `%s`, got `%s`", i, c.err, err)
			}
		}
		if c.exp != nil {
			if bytes.Compare(got, c.exp) != 0 {
				t.Errorf("[%d] - expected \n%x\n, got \n%x\n", i, c.exp, got)
			}
		}
		if c.err == nil && c.exp == nil {
			t.Errorf("[%d] - no assertion", i)
		}
	}
}
