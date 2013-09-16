package bytecode

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	// Predefined errors
	ErrVersionMismatch    = errors.New("the specified file version does not match the compiler's")
	ErrUnexpectedKValType = errors.New("unexpected constant value type")
	ErrInvalidKType       = errors.New("invalid constant type tag")
	ErrUnknownOpcode      = errors.New("unknown instruction opcode")
)

// An encoder takes an in-memory representation of agora code and encodes it into
// the agora bytecode representation.
type Encoder struct {
	w   io.Writer
	err error
}

// NewEncoder returns an Encoder that will write to the provided writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode encodes the provided in-memory File structure into the bytecode format,
// written to the encoder's writer. If an error is encountered, it is returned, otherwise
// it returns nil.
func (enc *Encoder) Encode(f *File) (err error) {
	// Reset error
	enc.err = nil
	// 1- Signature
	enc.write(_SIGNATURE)
	// 2- Version (must match exactly that of the compiler)
	enc.assertVersion(f)
	enc.write(encodeVersionByte(f.MajorVersion, f.MinorVersion))
	// 3- Each function
	for i, fn := range f.Fns {
		// 4- Function header
		if i == 0 {
			enc.write(f.Name) // The top-level function gets its name from the source file
		} else {
			enc.write(fn.Header.Name)
		}
		enc.write(fn.Header.StackSz)
		enc.write(fn.Header.ExpArgs)
		enc.write(fn.Header.ParentFnIx)
		enc.write(fn.Header.LineStart)
		enc.write(fn.Header.LineEnd)

		// 5- The K section
		enc.write(int64(len(fn.Ks)))
		for _, k := range fn.Ks {
			enc.assertKType(k.Type)
			enc.write(k)
		}

		// 6- The L section
		enc.write(int64(len(fn.Ls)))
		for _, l := range fn.Ls {
			enc.write(l)
		}

		// 7- The I section
		enc.write(int64(len(fn.Is)))
		for _, ins := range fn.Is {
			enc.assertOpcode(ins)
			enc.write(uint64(ins))
		}
	}
	return enc.err
}

func (enc *Encoder) guard(fn func()) {
	if enc.err != nil {
		return
	}
	fn()
}

func (enc *Encoder) assertOpcode(ins Instr) {
	enc.guard(func() {
		if ins.Opcode() >= op_max {
			enc.err = ErrUnknownOpcode
		}
	})
}

func (enc *Encoder) assertKType(kt KType) {
	enc.guard(func() {
		if _, ok := validKtypes[kt]; !ok {
			enc.err = ErrInvalidKType
		}
	})
}

func (enc *Encoder) assertVersion(f *File) {
	enc.guard(func() {
		if f.MajorVersion != _MAJOR_VERSION || f.MinorVersion != _MINOR_VERSION {
			enc.err = ErrVersionMismatch
		}
	})
}

func (enc *Encoder) write(v interface{}) {
	enc.guard(func() {
		switch val := v.(type) {
		case *K:
			enc.write(byte(val.Type))
			switch kval := val.Val.(type) {
			case string:
				if val.Type != KtString {
					enc.err = ErrUnexpectedKValType
					return
				}
				enc.write(kval)
			case int64:
				if val.Type != KtInteger && val.Type != KtBoolean {
					enc.err = ErrUnexpectedKValType
					return
				}
				enc.write(kval)
			case float64:
				if val.Type != KtFloat {
					enc.err = ErrUnexpectedKValType
					return
				}
				enc.write(kval)
			default:
				enc.err = ErrUnexpectedKValType
				return
			}
		case string:
			enc.write(int64(len(val)))
			enc.write([]byte(val))
		default:
			enc.err = binary.Write(enc.w, binary.LittleEndian, val)
		}
	})
}
