package bytecode

// The binary signature that must be present at the start of
// each compiled bytecode file.
const (
	_SIGNATURE int32 = 0x000A602A
)

var (
	// Vars only to allow for testing, but are really constants
	_MAJOR_VERSION = 0
	_MINOR_VERSION = 0
)

func encodeVersionByte(maj, min int) byte {
	return byte(maj)<<4 | byte(min)
}

func decodeVersionByte(v byte) (maj, min int) {
	return int(v >> 4), int((v << 4) >> 4)
}

// The type tag that defines the constant's type.
type KType byte

const (
	// The possible constant types
	KtInteger KType = 'i'
	KtBoolean KType = 'b'
	KtFloat   KType = 'f'
	KtString  KType = 's'
)

// The full representation of a bytecode file, as defined
// in /compiler/bytecode.md.
type File struct {
	Name         string
	MajorVersion int
	MinorVersion int
	Fns          []Fn
}

// The representation of a single function in a bytecode file.
type Fn struct {
	Header H
	Ks     []K
	Is     []Instr
}

// The function header representation.
type H struct {
	Name      string
	StackSz   int64
	ExpArgs   int64
	ExpVars   int64
	LineStart int64
	LineEnd   int64
}

// The representation of a single constant value.
type K struct {
	Type KType
	Val  interface{}
}
