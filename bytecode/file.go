package bytecode

// The binary signature that must be present at the start of
// each compiled bytecode file.
const (
	_SIGNATURE int32 = 0x000A602A
)

var (
	// Vars only to allow for testing, but are really constants
	_MAJOR_VERSION = 0
	_MINOR_VERSION = 1
)

// Version returns the major and minor version of the bytecode format.
func Version() (int, int) {
	return _MAJOR_VERSION, _MINOR_VERSION
}

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

var (
	validKtypes = map[KType]struct{}{
		KtInteger: struct{}{},
		KtBoolean: struct{}{},
		KtFloat:   struct{}{},
		KtString:  struct{}{},
	}
)

// A File is an in-memory representation of a bytecode file, as defined
// in /doc/bytecode.md.
type File struct {
	Name         string
	MajorVersion int
	MinorVersion int
	Fns          []*Fn
}

// NewFile returns a File structure initialized with the specified name and
// the current version.
func NewFile(nm string) *File {
	return &File{
		Name:         nm,
		MajorVersion: _MAJOR_VERSION,
		MinorVersion: _MINOR_VERSION,
	}
}

// A Fn is the representation of a single function in a bytecode file.
type Fn struct {
	Header H
	Ks     []*K
	Ls     []int64 // locals, as indexes into the K table
	Is     []Instr
}

// An H is the function header representation.
type H struct {
	Name       string
	StackSz    int64
	ExpArgs    int64
	ParentFnIx int64 // Lexical scope parent function, as index into the Fn table
	LineStart  int64
	LineEnd    int64
}

// A K is the representation of a single constant value.
type K struct {
	Type KType
	Val  interface{}
}
