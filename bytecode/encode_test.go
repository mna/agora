package bytecode

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/PuerkitoBio/agora/bytecode/testing"
)

var (
	enccases = []struct {
		maj int
		min int
		f   *File
		exp []byte
		err error
	}{
		0: {
			// Simplest case, encodes the file header only
			maj: defMaj,
			min: defMin,
			f: &File{
				MajorVersion: defMaj,
				MinorVersion: defMin,
			},
			exp: SigVer(_MAJOR_VERSION, _MINOR_VERSION),
		},
		1: {
			// Check version encoding, matching with compiler version
			maj: 1,
			min: 2,
			f:   &File{MajorVersion: 1, MinorVersion: 2},
			exp: append(ExpSig, 0x12),
		},
		2: {
			// Version mismatch error
			maj: defMaj,
			min: defMin,
			f:   &File{MinorVersion: defMin + 1},
			err: ErrVersionMismatch,
		},
		3: {
			// Top-level function gets the file name
			maj: defMaj,
			min: defMin,
			f: &File{
				MajorVersion: defMaj,
				MinorVersion: defMin,
				Name:         "test", Fns: []*Fn{&Fn{}}},
			exp: AppendAny(SigVer(_MAJOR_VERSION, _MINOR_VERSION), Int64ToByteSlice(4), 't', 'e', 's', 't',
				// StackSz - ExpArgs - ParentFnIx - LineStart - LineEnd
				ExpZeroInt64, ExpZeroInt64, ExpZeroInt64, ExpZeroInt64, ExpZeroInt64,
				// Ks - Ls - Is
				ExpZeroInt64, ExpZeroInt64, ExpZeroInt64),
		},
		4: {
			maj: defMaj,
			min: defMin,
			f: &File{
				MajorVersion: defMaj,
				MinorVersion: defMin,
				Name:         "test", Fns: []*Fn{
					&Fn{
						Header: H{
							StackSz:    2,
							ExpArgs:    3,
							ParentFnIx: 0,
							LineStart:  5,
							LineEnd:    6,
						},
						Ks: []*K{
							&K{
								Type: KtInteger,
								Val:  int64(7),
							},
						},
					},
				}},
			exp: AppendAny(SigVer(_MAJOR_VERSION, _MINOR_VERSION), Int64ToByteSlice(4), 't', 'e', 's', 't',
				// StackSz - ExpArgs - ParentFnIx - LineStart - LineEnd
				Int64ToByteSlice(2), Int64ToByteSlice(3), ExpZeroInt64, Int64ToByteSlice(5), Int64ToByteSlice(6),
				// Ks - Ls - Is
				Int64ToByteSlice(1), byte(KtInteger), Int64ToByteSlice(7), ExpZeroInt64, ExpZeroInt64),
		},
		5: {
			// Invalid KType
			maj: defMaj,
			min: defMin,
			f: &File{
				MajorVersion: defMaj,
				MinorVersion: defMin,
				Name:         "test", Fns: []*Fn{
					&Fn{
						Header: H{
							StackSz:    2,
							ExpArgs:    3,
							ParentFnIx: 0,
							LineStart:  5,
							LineEnd:    6,
						},
						Ks: []*K{
							&K{
								Type: KType('z'),
								Val:  int64(7),
							},
						},
					},
				}},
			err: ErrInvalidKType,
		},
		6: {
			// Invalid K value type
			maj: defMaj,
			min: defMin,
			f: &File{
				MajorVersion: defMaj,
				MinorVersion: defMin,
				Name:         "test", Fns: []*Fn{
					&Fn{
						Header: H{
							StackSz:   2,
							ExpArgs:   3,
							LineStart: 5,
							LineEnd:   6,
						},
						Ks: []*K{
							&K{
								Type: KtInteger,
								Val:  float64(3.5),
							},
						},
					},
				}},
			err: ErrUnexpectedKValType,
		},
		7: {
			maj: defMaj,
			min: defMin,
			f: &File{
				MajorVersion: defMaj,
				MinorVersion: defMin,
				Name:         "test", Fns: []*Fn{
					&Fn{
						Header: H{
							StackSz:    2,
							ExpArgs:    3,
							ParentFnIx: 4,
							LineStart:  5,
							LineEnd:    6,
						},
						Ks: []*K{
							&K{
								Type: KtInteger,
								Val:  int64(7),
							},
						},
						Is: []Instr{
							NewInstr(OP_ADD, FLG_K, 12),
							NewInstr(OP_DUMP, FLG_Sn, 0),
						},
					},
				}},
			exp: AppendAny(SigVer(_MAJOR_VERSION, _MINOR_VERSION), Int64ToByteSlice(4), 't', 'e', 's', 't',
				// StackSz - ExpArgs - ParentFnIx - LineStart - LineEnd
				Int64ToByteSlice(2), Int64ToByteSlice(3), Int64ToByteSlice(4), Int64ToByteSlice(5), Int64ToByteSlice(6),
				// Ks - Ls - Is
				Int64ToByteSlice(1), byte(KtInteger), Int64ToByteSlice(7), ExpZeroInt64, Int64ToByteSlice(2),
				// 2 ops
				0x0C, 0x00, 0x00, 0x00, 0x00, 0x00, byte(FLG_K), byte(OP_ADD), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, byte(FLG_Sn), byte(OP_DUMP)),
		},
		// Invalid opcode
		8: {
			maj: defMaj,
			min: defMin,
			f: &File{
				MajorVersion: defMaj,
				MinorVersion: defMin,
				Name:         "test", Fns: []*Fn{
					&Fn{
						Header: H{
							StackSz:   2,
							ExpArgs:   3,
							LineStart: 5,
							LineEnd:   6,
						},
						Ks: []*K{
							&K{
								Type: KtInteger,
								Val:  int64(7),
							},
						},
						Is: []Instr{
							NewInstr(Opcode(op_max+1), FLG_K, 12),
						},
					},
				}},
			err: ErrUnknownOpcode,
		},
		9: {
			// Multiple functions
			maj: defMaj,
			min: defMin,
			f: &File{
				MajorVersion: defMaj,
				MinorVersion: defMin,
				Name:         "test", Fns: []*Fn{
					&Fn{
						Header: H{
							StackSz:    2,
							ExpArgs:    3,
							ParentFnIx: 4,
							LineStart:  5,
							LineEnd:    6,
						},
						Ks: []*K{
							&K{
								Type: KtInteger,
								Val:  int64(7),
							},
						},
						Is: []Instr{
							NewInstr(OP_ADD, FLG_K, 12),
							NewInstr(OP_DUMP, FLG_Sn, 0),
						},
					},
					&Fn{
						Header: H{
							Name:       "f2",
							StackSz:    2,
							ExpArgs:    3,
							ParentFnIx: 0,
							LineStart:  5,
							LineEnd:    6,
						},
						Ks: []*K{
							&K{
								Type: KtString,
								Val:  "const",
							},
						},
						Is: []Instr{
							NewInstr(OP_RET, FLG__, 0),
						},
					},
				}},
			exp: AppendAny(SigVer(_MAJOR_VERSION, _MINOR_VERSION), Int64ToByteSlice(4), 't', 'e', 's', 't',
				// StackSz - ExpArgs - ParentFnIx - LineStart - LineEnd
				Int64ToByteSlice(2), Int64ToByteSlice(3), Int64ToByteSlice(4), Int64ToByteSlice(5), Int64ToByteSlice(6),
				// Ks - Ls - Is
				Int64ToByteSlice(1), byte(KtInteger), Int64ToByteSlice(7), ExpZeroInt64, Int64ToByteSlice(2),
				// 2 ops
				0x0C, 0x00, 0x00, 0x00, 0x00, 0x00, byte(FLG_K), byte(OP_ADD), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, byte(FLG_Sn), byte(OP_DUMP),
				// Fn 2
				Int64ToByteSlice(2), 'f', '2',
				// StackSz - ExpArgs - ParentFnIx - LineStart - LineEnd
				Int64ToByteSlice(2), Int64ToByteSlice(3), ExpZeroInt64, Int64ToByteSlice(5), Int64ToByteSlice(6),
				// Ks - Ls - Is
				Int64ToByteSlice(1), byte(KtString), Int64ToByteSlice(5), 'c', 'o', 'n', 's', 't', ExpZeroInt64, Int64ToByteSlice(1),
				// 1 op
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		},
	}

	isolateEncCase = -1
)

func TestEncode(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	for i, c := range enccases {
		if isolateEncCase >= 0 && isolateEncCase != i {
			continue
		}
		if testing.Verbose() {
			fmt.Printf("testing encode case %d...\n", i)
		}

		// Arrange
		_MAJOR_VERSION = c.maj
		_MINOR_VERSION = c.min
		buf.Reset()

		// Act
		err := enc.Encode(c.f)

		// Assert
		if err != c.err {
			if c.err == nil {
				t.Errorf("[%d] - expected no error, got `%s`", i, err)
			} else {
				t.Errorf("[%d] - expected error `%s`, got `%s`", i, c.err, err)
			}
		}
		if c.exp != nil {
			got := buf.Bytes()
			if bytes.Compare(got, c.exp) != 0 {
				t.Errorf("[%d] - expected \n%x\n, got \n%x\n", i, c.exp, got)
			}
		}
		if c.err == nil && c.exp == nil {
			t.Errorf("[%d] - no assertion", i)
		}
	}
}
