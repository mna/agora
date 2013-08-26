package bytecode

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrInvalidData = errors.New("input data is not valid bytecode")
)

type Decoder struct {
	r       io.Reader
	err     error
	sigRead bool
	sigOk   bool
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (dec *Decoder) IsBytecode() bool {
	if !dec.sigRead {
		dec.sigRead = true
		// Do not consume the bytes on the reader, use a bufio.Reader
		br := bufio.NewReader(dec.r)
		b, err := br.Peek(4)
		if err != nil {
			dec.sigOk = false
			return false
		}
		sig, n := binary.Varint(b)
		if n <= 0 {
			dec.sigOk = false
			return false
		}
		dec.assertSignature(int32(sig))
	}
	return dec.sigOk
}

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
		dec.sigOk = true
		if sig != _SIGNATURE {
			dec.sigOk = false
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
	fn.Header.ExpVars = dec.readInt64()
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
	dec.sigRead = true
	return sig
}

func (dec *Decoder) read(v interface{}) {
	dec.guard(func() {
		dec.err = binary.Read(dec.r, binary.LittleEndian, v)
	})
}
