package bytecode

import (
	"fmt"
)

// A Flag indicates the meaning of the index (or value) part of an instruction.
type Flag byte

const (
	// The possible values of Flag
	FLG__    Flag = iota // Ignored
	FLG_K                // Constant table index
	FLG_V                // Variable table index
	FLG_F                // Function prototype index
	FLG_A                // Arguments array
	FLG_N                // Nil value
	FLG_T                // `this` keyword
	FLG_An               // Args count in a CALL or CFLD instruction
	FLG_Jf               // Jump forward over n instructions
	FLG_Jb               // Jump back over n instructions
	FLG_Sn               // Dump n frames
	FLG_Fn               // Set n fields
	FLG_Fn2              // Set n field pairs
	FLG_INVL Flag = 0xFF // Invalid flag
)

var (
	// The lookup table of Flag values to literal flag names
	FlagNames = [...]string{
		FLG__:   "_",
		FLG_K:   "K",
		FLG_V:   "V",
		FLG_F:   "F",
		FLG_A:   "A",
		FLG_N:   "N",
		FLG_T:   "T",
		FLG_An:  "An",
		FLG_Jf:  "Jf",
		FLG_Jb:  "Jb",
		FLG_Sn:  "Sn",
		FLG_Fn:  "Fn",
		FLG_Fn2: "Fn2",
	}

	// The lookup table of literal flag names to Flag values
	FlagLookup = map[string]Flag{
		"_":   FLG__,
		"K":   FLG_K,
		"V":   FLG_V,
		"F":   FLG_F,
		"A":   FLG_A,
		"N":   FLG_N,
		"T":   FLG_T,
		"An":  FLG_An,
		"Jf":  FLG_Jf,
		"Jb":  FLG_Jb,
		"Sn":  FLG_Sn,
		"Fn":  FLG_Fn,
		"Fn2": FLG_Fn2,
	}
)

// NewFlag returns the Flag value identified by the provided literal name, or the
// invalid Flag value if the name is unknown.
func NewFlag(nm string) Flag {
	t, ok := FlagLookup[nm]
	if !ok {
		return FLG_INVL
	}
	return t
}

// String returns the literal name corresponding to the current Flag value.
func (f Flag) String() string {
	if f < 0 || int(f) > len(FlagNames)-1 {
		return ""
	}
	return FlagNames[f]
}

// An Instr is an agora instruction to be executed by the virtual machine at runtime.
//
// A bytecode instruction is a sequence of 64 bits arranged like this (a single letter=a byte):
// `ofvvvvvv`
// o: represents the opcode, on a single byte. See /bytecode/opcodes.go for the list of codes.
// f: indicates what the next value represents (Flag)
// v: represents the instruction's value (or index), on 6 bytes. Its meaning depends on the flag.
//    Gives a maximum possible value of 2^48.
type Instr uint64

// NewInstr returns an instruction value constructed from the provided opcode, flag and ix.
func NewInstr(op Opcode, flg Flag, ix uint64) Instr {
	return Instr(uint64(op)<<56 | uint64(flg)<<48 | ix)
}

// Opcode returns the opcode part of the instruction (the most significant byte).
func (i Instr) Opcode() Opcode {
	return Opcode(i >> 56)
}

// Flag returns the flag part of the instruction (the second most significant byte).
func (i Instr) Flag() Flag {
	return Flag((i << 8) >> 56)
}

// Index returns the index (or value) part of the instruction, that is, the 6 least
// significant bytes.
func (i Instr) Index() uint64 {
	return uint64((i << 16) >> 16)
}

// String returns a literal representation of the instruction.
func (i Instr) String() string {
	op, f, ix := i.Opcode(), i.Flag(), i.Index()
	return fmt.Sprintf("%-4s %-2s %3d", op, f, ix)
}
