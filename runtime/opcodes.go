package runtime

// The opcode takes 8 bytes, leaving 256 possible codes.
type Opcode byte

const (
	OP_RET Opcode = iota
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
	OP_LT
	OP_GT
	OP_TEST
	OP_JMPB
	OP_JMPF
	OP_INVL Opcode = 0xFF
)

var (
	OpNames = [...]string{
		OP_RET:  "RET ",
		OP_PUSH: "PUSH",
		OP_POP:  "POP ",
		OP_ADD:  "ADD ",
		OP_SUB:  "SUB ",
		OP_MUL:  "MUL ",
		OP_DIV:  "DIV ",
		OP_MOD:  "MOD ",
		OP_POW:  "POW ",
		OP_NOT:  "NOT ",
		OP_UNM:  "UNM ",
		OP_CALL: "CALL",
		OP_LT:   "LT  ",
		OP_GT:   "GT  ",
		OP_TEST: "TEST",
		OP_JMPB: "JMPB",
		OP_JMPF: "JMPF",
	}

	OpLookup = map[string]Opcode{
		"RET":  OP_RET,
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
		"LT":   OP_LT,
		"GT":   OP_GT,
		"TEST": OP_TEST,
		"JMPB": OP_JMPB,
		"JMPF": OP_JMPF,
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
