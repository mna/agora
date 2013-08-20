package bytecode

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrInvalidData = errors.New("the data to decode is not valid goblin bytecode")
)

type Decoder struct {
	r   io.Reader
	err error
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (dec *Decoder) Decode() (*File, error) {
	dec.err = nil
	// 1- Read and assert the signature
	sig := dec.readSignature()
	dec.assertSignature(sig)
	// 2- Read and assert the version
	ver := dec.readVersion()
	dec.assertVersion(ver)
	// Do not create useless structures if the header is invalid
	if dec.err != nil {
		return nil, dec.err
	}

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
	return f, dec.err
}

func (dec *Decoder) assertSignature(sig int32) {
	if dec.err != nil {
		return
	}
	if sig != _SIGNATURE {
		dec.err = ErrInvalidData
	}
}

func (dec *Decoder) assertVersion(ver byte) {
	if dec.err != nil {
		return
	}
	if ver != encodeVersionByte(_MAJOR_VERSION, _MINOR_VERSION) {
		dec.err = ErrVersionMismatch
	}
}

func (dec *Decoder) readFunc() (Fn, bool) {
	var fn Fn
	if dec.err != nil {
		return fn, false
	}
	// This first read *may* return io.EOF, this means that
	// there is no more function to read, return and cancel the error.
	nm := dec.readString()
	if dec.err == io.EOF {
		dec.err = nil
		return fn, false
	}
	fn.Header.Name = nm
	return fn, true
}

func (dec *Decoder) readVersion() byte {
	var ver byte
	dec.read(&ver)
	return ver
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

func (dec *Decoder) readSignature() int32 {
	var sig int32
	dec.read(&sig)
	return sig
}

func (dec *Decoder) read(v interface{}) {
	if dec.err != nil {
		return
	}
	dec.err = binary.Read(dec.r, binary.LittleEndian, v)
}
