package runtime

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrInvalidFile = errors.New("invalid bytecode file")

	// TODO : Add Major.Minor version in the compiled chunk
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
	f := newGoblinFunc(m)
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
			f.kTable = append(f.kTable, Bool(v != 0))
		case 'f':
			v := readFloat64(r)
			f.kTable = append(f.kTable, Float(v))
		case 's':
			v := readString(r)
			f.kTable = append(f.kTable, String(v))
		}
	}
}
