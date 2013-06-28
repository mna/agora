package runtime

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrInvalidFile = errors.New("invalid bytecode file")

	SIG = [...]byte{'6', '0', 'B', '1', '1', '4'}
)

func Undump(r io.Reader) (m Module, err error) {
	defer func() {
		if e := recover(); e != nil && e != io.EOF {
			m = nil
			err = e.(error)
		}
	}()

	readSig(r)
	nm := readString(r)
	m = newGoblinModule(nm)
	for {
		// Read a function
		readFunc(r, m.(*goblinModule))
	}
}

func readFunc(r io.Reader, m *goblinModule) {
	f := newGoblinFunc()
	f.stackSz = int(readInt64(r))
	f.expArgs = int(readInt64(r))
	f.expVars = int(readInt64(r))
	// TODO : Skip line start and line end, unused at the moment
	readInt64(r)
	readInt64(r)
	f.name = readString(r)
	m.fns = append(m.fns, f)

	readK(r, f)
	readI(r, f)
}

func readI(r io.Reader, f *GoblinFunc) {
	cnt := readInt64(r)
	for i := int64(0); i < cnt; i++ {
		f.code = append(f.code, Instr(readUint64(r)))
	}
}

func readK(r io.Reader, f *GoblinFunc) {
	cnt := readInt64(r)
	for i := int64(0); i < cnt; i++ {
		t := readByte(r)
		switch t {
		case 'i':
			v := int(readInt64(r))
			f.kTable = append(f.kTable, Int(v))
		case 'b':
			v := readInt64(r)
			f.kTable = append(f.kTable, Bool(v == 1))
		case 'f':
			v := readFloat64(r)
			f.kTable = append(f.kTable, Float(v))
		case 's':
			v := readString(r)
			f.kTable = append(f.kTable, String(v))
		}
	}
}

func readFloat64(r io.Reader) float64 {
	var v float64
	if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
		panic(err)
	}
	return v
}

func readInt64(r io.Reader) int64 {
	var v int64
	if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
		panic(err)
	}
	return v
}

func readUint64(r io.Reader) uint64 {
	var v uint64
	if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
		panic(err)
	}
	return v
}

func readByte(r io.Reader) byte {
	var b byte
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		panic(err)
	}
	return b
}

func readString(r io.Reader) string {
	var l int64
	if err := binary.Read(r, binary.LittleEndian, &l); err != nil {
		panic(err)
	}
	s := make([]byte, l)
	if err := binary.Read(r, binary.LittleEndian, s); err != nil {
		panic(err)
	}
	return string(s)
}

func readSig(r io.Reader) {
	var b [6]byte
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		panic(err)
	}
	if b != SIG {
		panic(ErrInvalidFile)
	}
}
