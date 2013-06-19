package runtime

type Table byte

const (
	TBL_K Table = iota
	TBL_V
	TBL_INVL Table = 0xFF
)

var (
	TableNames = [...]string{
		TBL_K: "K",
		TBL_V: "V",
	}

	TableLookup = map[string]Table{
		"K": TBL_K,
		"V": TBL_V,
	}
)

func NewTable(nm string) Table {
	t, ok := TableLookup[nm]
	if !ok {
		return TBL_INVL
	}
	return t
}

// A bytecode instruction is a sequence of 64 bits arranged like this (a single letter=a byte):
// `oabbbbbb`
// o: represents the opcode, on a single byte. See opcodes.go for the list of codes.
// a: represents on what data to operate, on a single byte.
//    - 0 means from the KTable (constant)
//    - 1 means from the VTable (variable)
// b: represents the index of the data in the relevant table (K or V), on 6 bytes.
//    Gives a possibility of 2^48 items in each table.
type Instr uint64

func NewInstr(op Opcode, tbl Table, ix uint64) Instr {
	return Instr(uint64(op)<<56 | uint64(tbl)<<48 | ix)
}

func (ø Instr) Opcode() Opcode {
	return Opcode(ø >> 56)
}

func (ø Instr) Table() Table {
	return Table((ø << 8) >> 56)
}

func (ø Instr) Index() uint64 {
	return uint64((ø << 16) >> 16)
}
