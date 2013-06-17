package runtime

/*

** STACK MANIPULATION AND MEMORY STORAGE **

PUSH index : pushes onto the stack the value held by the variable identified by KTable[index] (found in the VTable, recursively through the scopes, all the way to the global variables) or the value itself from KTable[index] if the index is negative.

POP index : pops the value on top of the stack and stores it into the variable identified by the name held in KTable[index] (found in the VTable, recursively through the scopes).

TODO : How to manage objects with fields, i.e. `obj.field1.subfield3 = 10`. Store the whole string in the KTable, and point it to the top-level value in the VTable, and strip at points to find sub-fields? Would not allow the use of any characters in the field names, like JSON does.

** ARITHMETIC AND LOGICAL OPERATIONS **

ADD : adds the two top-most values on the stack, stores the result in their place.
SUB : subtracts the two top-most values on the stack, stores the result in their place.
MUL : multiplies the two top-most values on the stack, stores the result in their place.
DIV : divides the two top-most values on the stack, stores the result in their place.
POW : computes the exponent using the two top-most values on the stack, stores the result in their place.
MOD : computes the remainder using the two top-most values on the stack, stores the result in their place.
NOT : negates the value of the top-most value on the stack, stores the result in its place.
UNM : switches the sign of the top-most value on the stack, stores the result in its place.

** CONTROL-FLOW OPERATIONS **

JMP [delta] : performs an unconditional jump. By default, jumps over the next instruction, but can specify a delta (positive or negative). A delta of 1 jumps over 1 instruction, -1 returns to the previous instruction, etc.

*/
