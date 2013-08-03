package runtime

// The opcode takes 8 bytes, leaving 256 possible codes.
type Opcode byte

const (
	OP_RET Opcode = iota
	OP_LOAD
	OP_PUSH
	OP_POP
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD
	OP_POW
	OP_NOT
	OP_UNM
	OP_CALL
	OP_EQ
	OP_LT
	OP_LTE
	OP_GT
	OP_GTE
	OP_AND
	OP_OR
	OP_TEST
	OP_JMPB
	OP_JMPF
	OP_NEW
	OP_SFLD
	OP_GFLD
	OP_CFLD
	// Debugging
	OP_DUMP
	OP_INVL Opcode = 0xFF
)

var (
	OpNames = [...]string{
		OP_RET:  "RET",
		OP_LOAD: "LOAD",
		OP_PUSH: "PUSH",
		OP_POP:  "POP",
		OP_ADD:  "ADD",
		OP_SUB:  "SUB",
		OP_MUL:  "MUL",
		OP_DIV:  "DIV",
		OP_MOD:  "MOD",
		OP_POW:  "POW",
		OP_NOT:  "NOT",
		OP_UNM:  "UNM",
		OP_CALL: "CALL",
		OP_EQ:   "EQ",
		OP_LT:   "LT",
		OP_LTE:  "LTE",
		OP_GT:   "GT",
		OP_GTE:  "GTE",
		OP_AND:  "AND",
		OP_OR:   "OR",
		OP_TEST: "TEST",
		OP_JMPB: "JMPB",
		OP_JMPF: "JMPF",
		OP_NEW:  "NEW",
		OP_SFLD: "SFLD",
		OP_GFLD: "GFLD",
		OP_CFLD: "CFLD",
		OP_DUMP: "DUMP",
	}

	OpLookup = map[string]Opcode{
		"RET":  OP_RET,
		"LOAD": OP_LOAD,
		"PUSH": OP_PUSH,
		"POP":  OP_POP,
		"ADD":  OP_ADD,
		"SUB":  OP_SUB,
		"MUL":  OP_MUL,
		"DIV":  OP_DIV,
		"MOD":  OP_MOD,
		"POW":  OP_POW,
		"NOT":  OP_NOT,
		"UNM":  OP_UNM,
		"CALL": OP_CALL,
		"EQ":   OP_EQ,
		"LT":   OP_LT,
		"LTE":  OP_LTE,
		"GT":   OP_GT,
		"GTE":  OP_GTE,
		"AND":  OP_AND,
		"OR":   OP_OR,
		"TEST": OP_TEST,
		"JMPB": OP_JMPB,
		"JMPF": OP_JMPF,
		"NEW":  OP_NEW,
		"SFLD": OP_SFLD,
		"GFLD": OP_GFLD,
		"CFLD": OP_CFLD,
		"DUMP": OP_DUMP,
	}
)

func NewOpcode(nm string) Opcode {
	o, ok := OpLookup[nm]
	if !ok {
		return OP_INVL
	}
	return o
}

func (ø Opcode) String() string {
	return OpNames[ø]
}
