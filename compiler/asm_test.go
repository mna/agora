package compiler

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goblin/bytecode"
	. "github.com/PuerkitoBio/goblin/bytecode/testing"
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
			exp: AppendAny(SigVer(bytecode.Version()), Int64ToByteSlice(4), 't', 'e', 's', 't', ExpZeroInt64, ExpZeroInt64, ExpZeroInt64, ExpZeroInt64, ExpZeroInt64, ExpZeroInt64, ExpZeroInt64),
		},
		2: {
			id: "test",
			src: `
[f]
test
1
0
1
0
2
[k]
sa
i5
[i]
PUSH K 1 // Push constant value 5 on the stack
POP V 0  // Pop the value from the stack into variable identified by constant 0 (a)
PUSH V 0 // Push value of variable identified by constant 0 on the stack (a)
DUMP S 1
RET _ 0
`,
			exp: AppendAny(SigVer(bytecode.Version()), Int64ToByteSlice(4), 't', 'e', 's', 't', Int64ToByteSlice(1), ExpZeroInt64, Int64ToByteSlice(1), ExpZeroInt64, Int64ToByteSlice(2),
				Int64ToByteSlice(2),
				's', Int64ToByteSlice(1), 'a', 'i', Int64ToByteSlice(5),
				Int64ToByteSlice(5),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("PUSH"), bytecode.NewFlag("K"), 1))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("POP"), bytecode.NewFlag("V"), 0))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("PUSH"), bytecode.NewFlag("V"), 0))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("DUMP"), bytecode.NewFlag("S"), 1))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("RET"), bytecode.NewFlag("_"), 0))),
			),
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
