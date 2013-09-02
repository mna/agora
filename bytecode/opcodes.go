package bytecode

import (
	"strconv"
)

// The opcode takes one byte, leaving 256 possible codes.
type Opcode byte

const (
	// The possible opcodes
	OP_RET  Opcode = iota // return
	OP_PUSH               // push a value onto the stack
	OP_POP                // pop a value from the stack
	OP_ADD                // add two values from the stack, push the result
	OP_SUB                // subtract two values from the stack, push the result
	OP_MUL                // multiply two values from the stack, push the result
	OP_DIV                // divide two values from the stack, push the result
	OP_MOD                // compute the modulo of two values from the stack, push the result
	OP_NOT                // boolean negation of one value from the stack, push the result
	OP_UNM                // unary minus of one value from the stack, push the result
	OP_EQ                 // check equality of two values from the stack, push the result
	OP_LT                 // lower than on two values from the stack, push the result
	OP_LTE                // lower than or equal on two values from the stack, push the result
	OP_GT                 // greater than on two values from the stack, push the result
	OP_GTE                // greater than or equal on two values from the stack, push the result
	OP_AND                // boolean `and` on two values from the stack, push the result
	OP_OR                 // boolean `or` on two values from the stack, push the result
	OP_TEST               // check the boolean value on top of the stack, if false jump n instructions
	OP_JMP                // perform an unconditional jump (forward or backward, depending on the flag)
	OP_NEW                // create and initialize a new object, push the result
	OP_SFLD               // set the value of an object's field, using 3 values from the stack (object variable, key and value)
	OP_GFLD               // get the value of an object's field, push the result, using 2 values from the stack (object variable and key)
	OP_CFLD               // call a method on an object, push the result, using 2 values + n arguments from the stack (object variable and key)
	OP_CALL               // call a function, push the result, using 1 value + n arguments from the stack
	op_dbgstart
	OP_DUMP               // print the execution context, if the Ctx is in debug mode
	op_max                // Indicates the maximum legal opcode
	OP_INVL Opcode = 0xFF // Invalid opcode
)

var (
	// Lookup table of opcodes to literal name
	OpNames = [...]string{
		OP_RET:  "RET",
		OP_PUSH: "PUSH",
		OP_POP:  "POP",
		OP_ADD:  "ADD",
		OP_SUB:  "SUB",
		OP_MUL:  "MUL",
		OP_DIV:  "DIV",
		OP_MOD:  "MOD",
		OP_NOT:  "NOT",
		OP_UNM:  "UNM",
		OP_EQ:   "EQ",
		OP_LT:   "LT",
		OP_LTE:  "LTE",
		OP_GT:   "GT",
		OP_GTE:  "GTE",
		OP_AND:  "AND",
		OP_OR:   "OR",
		OP_TEST: "TEST",
		OP_JMP:  "JMP",
		OP_NEW:  "NEW",
		OP_SFLD: "SFLD",
		OP_GFLD: "GFLD",
		OP_CFLD: "CFLD",
		OP_CALL: "CALL",
		OP_DUMP: "DUMP",
	}

	// Loopup table of literal opcode names to Opcode value
	OpLookup = map[string]Opcode{
		"RET":  OP_RET,
		"PUSH": OP_PUSH,
		"POP":  OP_POP,
		"ADD":  OP_ADD,
		"SUB":  OP_SUB,
		"MUL":  OP_MUL,
		"DIV":  OP_DIV,
		"MOD":  OP_MOD,
		"NOT":  OP_NOT,
		"UNM":  OP_UNM,
		"EQ":   OP_EQ,
		"LT":   OP_LT,
		"LTE":  OP_LTE,
		"GT":   OP_GT,
		"GTE":  OP_GTE,
		"AND":  OP_AND,
		"OR":   OP_OR,
		"TEST": OP_TEST,
		"JMP":  OP_JMP,
		"NEW":  OP_NEW,
		"SFLD": OP_SFLD,
		"GFLD": OP_GFLD,
		"CFLD": OP_CFLD,
		"CALL": OP_CALL,
		"DUMP": OP_DUMP,
	}
)

// NewOpcode returns the opcode value corresponding to the provided name, or
// the invalid opcode if the name is unknown.
func NewOpcode(nm string) Opcode {
	o, ok := OpLookup[nm]
	if !ok {
		return OP_INVL
	}
	return o
}

// String returns the literal string representation of the opcode value, or the
// string representation of the opcode number if it is unknown.
func (o Opcode) String() string {
	if o >= 0 && int(o) < len(OpNames) {
		return OpNames[o]
	}
	return strconv.Itoa(int(o))
}
