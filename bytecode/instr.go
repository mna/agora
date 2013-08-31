package bytecode

import (
	"fmt"
)

type Flag byte

const (
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
	FLG_Sn               // For DUMP, dump n frames
	FLG_Fn               // For NEW, set n fields
	FLG_INVL Flag = 0xFF // Invalid flag
)

var (
	FlagNames = [...]string{
		FLG__:  "_",
		FLG_K:  "K",
		FLG_V:  "V",
		FLG_F:  "F",
		FLG_A:  "A",
		FLG_N:  "N",
		FLG_T:  "T",
		FLG_An: "An",
		FLG_Jf: "Jf",
		FLG_Jb: "Jb",
		FLG_Sn: "Sn",
		FLG_Fn: "Fn",
	}

	FlagLookup = map[string]Flag{
		"_":  FLG__,
		"K":  FLG_K,
		"V":  FLG_V,
		"F":  FLG_F,
		"A":  FLG_A,
		"N":  FLG_N,
		"T":  FLG_T,
		"An": FLG_An,
		"Jf": FLG_Jf,
		"Jb": FLG_Jb,
		"Sn": FLG_Sn,
		"Fn": FLG_Fn,
	}
)

func NewFlag(nm string) Flag {
	t, ok := FlagLookup[nm]
	if !ok {
		return FLG_INVL
	}
	return t
}

func (ø Flag) String() string {
	return FlagNames[ø]
}

// A bytecode instruction is a sequence of 64 bits arranged like this (a single letter=a byte):
// `oabbbbbb`
// o: represents the opcode, on a single byte. See opcodes.go for the list of codes.
// a: indicates what the next value represents (Flag)
// b: represents the index of the data in the relevant table (K or V), on 6 bytes.
//    Gives a possibility of 2^48 items in each table.
type Instr uint64

func NewInstr(op Opcode, flg Flag, ix uint64) Instr {
	return Instr(uint64(op)<<56 | uint64(flg)<<48 | ix)
}

func (ø Instr) Opcode() Opcode {
	return Opcode(ø >> 56)
}

func (ø Instr) Flag() Flag {
	return Flag((ø << 8) >> 56)
}

func (ø Instr) Index() uint64 {
	return uint64((ø << 16) >> 16)
}

func (ø Instr) String() string {
	op, f, ix := ø.Opcode(), ø.Flag(), ø.Index()
	return fmt.Sprintf("%-4s %-2s %3d", op, f, ix)
}
