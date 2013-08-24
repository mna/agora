package compiler

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/agora/bytecode"
	. "github.com/PuerkitoBio/agora/bytecode/testing"
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
			// Full valid func
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
		3: {
			// Unknown opcode
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
PRUT V 0 
DUMP S 1
RET _ 0
`,
			err: bytecode.ErrUnknownOpcode,
		},
		4: {
			// Invalid instruction
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
ADD X
DUMP S 1
RET _ 0
`,
			err: ErrInvalidInstruction,
		},
		5: {
			// Many functions, valid
			id: "test",
			src: `
//
// func Add(x, y) { // Essentially means var Add = func ...
//   return x + y
// }
// return Add(4, "198")
//
[f]
test
3
0
1
0
3
[k]
sAdd
i4
s198
[i]
PUSH F 1
POP V 0
PUSH K 1
PUSH K 2
PUSH V 0
CALL A 2
DUMP S 1
RET _ 0
[f]
Add
2
2
2
0
2
[k]
sx
sy
[i]
PUSH V 0
PUSH V 1
ADD _ 0
RET _ 0
`,
			exp: AppendAny(SigVer(bytecode.Version()), Int64ToByteSlice(4), 't', 'e', 's', 't', Int64ToByteSlice(3), ExpZeroInt64, Int64ToByteSlice(1), ExpZeroInt64, Int64ToByteSlice(3),
				Int64ToByteSlice(3),
				's', Int64ToByteSlice(3), 'A', 'd', 'd', 'i', Int64ToByteSlice(4), 's', Int64ToByteSlice(3), '1', '9', '8',
				Int64ToByteSlice(8),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("PUSH"), bytecode.NewFlag("F"), 1))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("POP"), bytecode.NewFlag("V"), 0))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("PUSH"), bytecode.NewFlag("K"), 1))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("PUSH"), bytecode.NewFlag("K"), 2))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("PUSH"), bytecode.NewFlag("V"), 0))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("CALL"), bytecode.NewFlag("A"), 2))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("DUMP"), bytecode.NewFlag("S"), 1))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("RET"), bytecode.NewFlag("_"), 0))),
				Int64ToByteSlice(3), 'A', 'd', 'd', Int64ToByteSlice(2), Int64ToByteSlice(2), Int64ToByteSlice(2), ExpZeroInt64, Int64ToByteSlice(2),
				Int64ToByteSlice(2),
				's', Int64ToByteSlice(1), 'x', 's', Int64ToByteSlice(1), 'y',
				Int64ToByteSlice(4),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("PUSH"), bytecode.NewFlag("V"), 0))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("PUSH"), bytecode.NewFlag("V"), 1))),
				UInt64ToByteSlice(uint64(bytecode.NewInstr(bytecode.NewOpcode("ADD"), bytecode.NewFlag("_"), 0))),
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
