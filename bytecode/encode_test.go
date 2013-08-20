package bytecode

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

var (
	// Useful prefilled expected byte slices
	expSig          = []byte{0x2A, 0x60, 0x0A, 0x00}
	expSigAndDefVer = append(expSig, encodeVersionByte(_MAJOR_VERSION, _MINOR_VERSION))
	expZeroInt64    = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	enccases = []struct {
		maj int
		min int
		f   *File
		exp []byte
		err error
	}{
		0: {
			// Simplest case, encodes the file header only
			f:   &File{},
			exp: expSigAndDefVer,
		},
		1: {
			// Check version encoding, matching with compiler version
			maj: 1,
			min: 2,
			f:   &File{MajorVersion: 1, MinorVersion: 2},
			exp: append(expSig, 0x12),
		},
		2: {
			// Version mismatch error
			f:   &File{MinorVersion: 1},
			err: ErrVersionMismatch,
		},
		3: {
			// Top-level function gets the file name
			f:   &File{Name: "test", Fns: []Fn{Fn{}}},
			exp: appendAny(expSigAndDefVer, int64ToByteSlice(4), 't', 'e', 's', 't', expZeroInt64, expZeroInt64, expZeroInt64, expZeroInt64, expZeroInt64, expZeroInt64, expZeroInt64),
		},
		4: {
			f: &File{Name: "test", Fns: []Fn{
				Fn{
					Header: H{
						StackSz:   2,
						ExpArgs:   3,
						ExpVars:   4,
						LineStart: 5,
						LineEnd:   6,
					},
					Ks: []K{
						K{
							Type: KtInteger,
							Val:  int64(7),
						},
					},
				},
			}},
			exp: appendAny(expSigAndDefVer, int64ToByteSlice(4), 't', 'e', 's', 't', int64ToByteSlice(2), int64ToByteSlice(3), int64ToByteSlice(4), int64ToByteSlice(5), int64ToByteSlice(6), int64ToByteSlice(1), byte(KtInteger), int64ToByteSlice(7), expZeroInt64),
		},
		5: {
			// Invalid KType
			f: &File{Name: "test", Fns: []Fn{
				Fn{
					Header: H{
						StackSz:   2,
						ExpArgs:   3,
						ExpVars:   4,
						LineStart: 5,
						LineEnd:   6,
					},
					Ks: []K{
						K{
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
			f: &File{Name: "test", Fns: []Fn{
				Fn{
					Header: H{
						StackSz:   2,
						ExpArgs:   3,
						ExpVars:   4,
						LineStart: 5,
						LineEnd:   6,
					},
					Ks: []K{
						K{
							Type: KtInteger,
							Val:  float64(3.5),
						},
					},
				},
			}},
			err: ErrUnexpectedKValType,
		},
		7: {
			f: &File{Name: "test", Fns: []Fn{
				Fn{
					Header: H{
						StackSz:   2,
						ExpArgs:   3,
						ExpVars:   4,
						LineStart: 5,
						LineEnd:   6,
					},
					Ks: []K{
						K{
							Type: KtInteger,
							Val:  int64(7),
						},
					},
					Is: []Instr{
						NewInstr(OP_ADD, FLG_K, 12),
						NewInstr(OP_DUMP, FLG_S, 0),
					},
				},
			}},
			exp: appendAny(expSigAndDefVer, int64ToByteSlice(4), 't', 'e', 's', 't', int64ToByteSlice(2), int64ToByteSlice(3), int64ToByteSlice(4), int64ToByteSlice(5), int64ToByteSlice(6), int64ToByteSlice(1), byte(KtInteger), int64ToByteSlice(7), int64ToByteSlice(2), 0x0C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x1B),
		},
		// Invalid opcode
		8: {
			f: &File{Name: "test", Fns: []Fn{
				Fn{
					Header: H{
						StackSz:   2,
						ExpArgs:   3,
						ExpVars:   4,
						LineStart: 5,
						LineEnd:   6,
					},
					Ks: []K{
						K{
							Type: KtInteger,
							Val:  int64(7),
						},
					},
					Is: []Instr{
						NewInstr(Opcode(250), FLG_K, 12),
					},
				},
			}},
			err: ErrUnknownOpcode,
		},
		9: {
			// Multiple functions
			f: &File{Name: "test", Fns: []Fn{
				Fn{
					Header: H{
						StackSz:   2,
						ExpArgs:   3,
						ExpVars:   4,
						LineStart: 5,
						LineEnd:   6,
					},
					Ks: []K{
						K{
							Type: KtInteger,
							Val:  int64(7),
						},
					},
					Is: []Instr{
						NewInstr(OP_ADD, FLG_K, 12),
						NewInstr(OP_DUMP, FLG_S, 0),
					},
				},
				Fn{
					Header: H{
						Name:      "f2",
						StackSz:   2,
						ExpArgs:   3,
						ExpVars:   4,
						LineStart: 5,
						LineEnd:   6,
					},
					Ks: []K{
						K{
							Type: KtString,
							Val:  "const",
						},
					},
					Is: []Instr{
						NewInstr(OP_RET, FLG__, 0),
					},
				},
			}},
			exp: appendAny(expSigAndDefVer, int64ToByteSlice(4), 't', 'e', 's', 't', int64ToByteSlice(2), int64ToByteSlice(3), int64ToByteSlice(4), int64ToByteSlice(5), int64ToByteSlice(6), int64ToByteSlice(1), byte(KtInteger), int64ToByteSlice(7), int64ToByteSlice(2), 0x0C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x1B, int64ToByteSlice(2), 'f', '2', int64ToByteSlice(2), int64ToByteSlice(3), int64ToByteSlice(4), int64ToByteSlice(5), int64ToByteSlice(6), int64ToByteSlice(1), byte(KtString), int64ToByteSlice(5), 'c', 'o', 'n', 's', 't', int64ToByteSlice(1), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		},
	}

	isolateEncCase = -1
)

func appendAny(b []byte, vals ...interface{}) []byte {
	for i, v := range vals {
		switch v := v.(type) {
		case []byte:
			b = append(b, v...)
		case byte:
			b = append(b, v)
		case int32:
			b = append(b, byte(v))
		case int:
			b = append(b, byte(v))
		default:
			panic(fmt.Sprintf("invalid type to append at pos %d", i))
		}
	}
	return b
}

func int64ToByteSlice(i int64) []byte {
	buf := bytes.NewBuffer(nil)
	if err := binary.Write(buf, binary.LittleEndian, i); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

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
