package bytecode

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	// Predefined errors
	ErrInvalidData = errors.New("input data is not valid bytecode")
)

// A Decoder reads a bytecode-encoded source into a structured representation in memory.
type Decoder struct {
	r   io.Reader
	err error
}

// NewDecoder returns a Decoder that reads from the provided reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// IsBytecode checks if the provided reader reads from a bytecode-encoded source.
// It checks if the agora bytecode signature is present at the start of the data.
func IsBytecode(rs io.ReadSeeker) bool {
	var i int32
	if err := binary.Read(rs, binary.LittleEndian, i); err != nil {
		return false
	}
	defer rs.Seek(0, 0)
	return i == _SIGNATURE
}

// Decode reads the bytecode-encoded source into an in-memory data structure, and
// returns the File structure containing the translated bytecode, or an error.
func (dec *Decoder) Decode() (*File, error) {
	// 1- Read and assert the signature
	sig := dec.readSignature()
	dec.assertSignature(sig)
	if dec.err != nil {
		return nil, ErrInvalidData
	}
	// 2- Read and assert the version
	ver := dec.readByte()
	dec.assertVersion(ver)
	// 3- Create the File structure
	f := new(File)
	f.MajorVersion, f.MinorVersion = decodeVersionByte(ver)
	for {
		fn, ok := dec.readFunc()
		if !ok {
			break
		}
		f.Fns = append(f.Fns, fn)
		if len(f.Fns) == 1 {
			f.Name = fn.Header.Name
		}
	}
	// If the error is EOF at this point, cancel it
	if dec.err == io.EOF {
		dec.err = nil
	}
	return f, dec.err
}

func (dec *Decoder) guard(fn func()) {
	if dec.err != nil {
		return
	}
	fn()
}

func (dec *Decoder) assertSignature(sig int32) {
	dec.guard(func() {
		if sig != _SIGNATURE {
			dec.err = ErrInvalidData
		}
	})
}

func (dec *Decoder) assertVersion(ver byte) {
	dec.guard(func() {
		if ver != encodeVersionByte(_MAJOR_VERSION, _MINOR_VERSION) {
			dec.err = ErrVersionMismatch
		}
	})
}

func (dec *Decoder) assertKType(kt KType) {
	dec.guard(func() {
		if _, ok := validKtypes[kt]; !ok {
			dec.err = ErrInvalidKType
		}
	})
}

func (dec *Decoder) assertOpcode(ins Instr) {
	dec.guard(func() {
		if ins.Opcode() >= op_max {
			dec.err = ErrUnknownOpcode
		}
	})
}

func (dec *Decoder) readFunc() (*Fn, bool) {
	nm := dec.readString()
	if dec.err != nil {
		return nil, false
	}
	fn := new(Fn)

	// Function header
	fn.Header.Name = nm
	fn.Header.StackSz = dec.readInt64()
	fn.Header.ExpArgs = dec.readInt64()
	fn.Header.ParentFnIx = dec.readInt64()
	fn.Header.LineStart = dec.readInt64()
	fn.Header.LineEnd = dec.readInt64()

	// K section
	ks := dec.readInt64()
	if ks > 0 {
		fn.Ks = make([]*K, ks)
		for i := int64(0); i < ks; i++ {
			fn.Ks[i] = dec.readK()
		}
	}

	// L section
	ls := dec.readInt64()
	if ls > 0 {
		fn.Ls = make([]int64, ls)
		for i := int64(0); i < ls; i++ {
			fn.Ls[i] = dec.readInt64()
		}
	}

	// I section
	is := dec.readInt64()
	if is > 0 {
		fn.Is = make([]Instr, is)
		for i := int64(0); i < is; i++ {
			fn.Is[i] = Instr(dec.readUInt64())
			dec.assertOpcode(fn.Is[i])
		}
	}
	return fn, true
}

func (dec *Decoder) readK() *K {
	k := new(K)
	k.Type = KType(dec.readByte())
	dec.assertKType(k.Type)
	switch k.Type {
	case KtInteger, KtBoolean:
		k.Val = dec.readInt64()
	case KtFloat:
		k.Val = dec.readFloat64()
	case KtString:
		k.Val = dec.readString()
	}

	return k
}

func (dec *Decoder) readByte() byte {
	var b byte
	dec.read(&b)
	return b
}

func (dec *Decoder) readString() string {
	l := dec.readInt64()
	if l <= 0 {
		return ""
	}
	buf := make([]byte, l)
	dec.read(buf)
	return string(buf)
}

func (dec *Decoder) readInt64() int64 {
	var i int64
	dec.read(&i)
	return i
}

func (dec *Decoder) readUInt64() uint64 {
	var u uint64
	dec.read(&u)
	return u
}

func (dec *Decoder) readFloat64() float64 {
	var f float64
	dec.read(&f)
	return f
}

func (dec *Decoder) readSignature() int32 {
	var sig int32
	dec.read(&sig)
	return sig
}

func (dec *Decoder) read(v interface{}) {
	dec.guard(func() {
		dec.err = binary.Read(dec.r, binary.LittleEndian, v)
	})
}
