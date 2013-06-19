package runtime

type debug struct {
	Name      string
	File      string
	LineStart int
	LineEnd   int
}

type Var struct {
	debug
}

type FuncProto struct {
	KTable []Val
	VTable []Var
	Code   []Instr
	debug
}

type Func struct {
	*FuncProto
	pc    int
	vars  []Val
	stack []Val
	sp    int
}
