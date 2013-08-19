package bytecode

// The binary signature that must be present at the start of
// each compiled bytecode file.
const (
	_MAJOR_VERSION       = 0
	_MINOR_VERSION       = 0
	_SIGNATURE     int32 = 0x000A602A
)

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
	Is     []I
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

// The representation of an instruction. For the bytecode format
// definition (and to avoid depending on the runtime package),
// there is no need to understand the internal composition of
// the instruction.
type I uint64
