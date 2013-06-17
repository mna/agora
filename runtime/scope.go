package runtime

/*
Scope is an interface.

Universe is the global scope. It holds FTable (table of functions, callable via index with the CALL opcode), VTable (table of global variables, which should be a map for direct access by name, since it will never be accessed by index, will go through the KTable first), KTable (table of constants in the global scope - constants are variable names and literal values).

Function is the function scope. It holds the VTable and KTable, and points to its parent scope (maybe another function, or the Universe scope).

*/

type Scope interface {
	Parent() Scope
	Execute()
}

type commonScope struct {
	vTable map[string]Val
	kTable []Val
	stack  []Val
	sp     int
	code   []instruction
	pc     int
	parent Scope
}

func (ø *commonScope) Parent() Scope {
	return ø.parent
}

type UScope struct {
	*commonScope
	fTable []Function
}

type FScope struct {
	*commonScope
}
