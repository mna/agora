package runtime

import (
	"fmt"
	"io"
	"os"
)

type Ctx struct {
	Protos []*FuncProto
	stdout io.ReadWriter
	stdin  io.ReadWriter
	stderr io.ReadWriter
	logic  Logic

	// Native funcs table
	nTable map[string]NativeFunc

	// Call stack
	callstack []*Func
	callsp    int
}

func NewCtx() *Ctx {
	return &Ctx{
		nil,
		os.Stdout,
		os.Stdin,
		os.Stderr,
		defaultLogic{},
		make(map[string]NativeFunc),
		make([]*Func, 2),
		0,
	}
}

func (ø *Ctx) SetStdStreams(stdin, stdout, stderr io.ReadWriter) {
	ø.stdin, ø.stdout, ø.stderr = stdin, stdout, stderr
}

func (ø *Ctx) Run() interface{} {
	if len(ø.Protos) == 0 {
		panic("no function available in this context")
	}
	f := newFunc(ø, ø.Protos[0])
	return f.Call().Native()
}

func (ø *Ctx) RegisterNativeFuncs(fs map[string]NativeFunc) {
	for k, v := range fs {
		ø.nTable[k] = v
	}
}

func (ø *Ctx) Stdout() io.ReadWriter {
	return ø.stdout
}

func (ø *Ctx) Stdin() io.ReadWriter {
	return ø.stdin
}

func (ø *Ctx) Stderr() io.ReadWriter {
	return ø.stderr
}

func (ø *Ctx) push(f *Func) {
	// Stack has to grow as needed
	if ø.callsp == len(ø.callstack) {
		if ø.callsp == cap(ø.callstack) {
			fmt.Printf("DEBUG expanding call stack of ctx, current size: %d\n", len(ø.callstack))
		}
		ø.callstack = append(ø.callstack, f)
	} else {
		ø.callstack[ø.callsp] = f
	}
	ø.callsp++

}

func (ø *Ctx) pop() *Func {
	ø.callsp--
	f := ø.callstack[ø.callsp]
	ø.callstack[ø.callsp] = nil // free this reference for gc
	return f
}

func (ø *Ctx) getVar(nm string) (Val, bool) {
	// Current call is ø.callsp - 1
	for i := ø.callsp - 1; i >= 0; i-- {
		f := ø.callstack[i]
		if v, ok := f.vars[nm]; ok {
			return v, true
		}
	}
	return Nil, false
}

func (ø *Ctx) setVar(nm string, v Val) bool {
	// Current call is ø.callsp - 1
	for i := ø.callsp - 1; i >= 0; i-- {
		f := ø.callstack[i]
		if _, ok := f.vars[nm]; ok {
			f.vars[nm] = v
			return true
		}
	}
	return false
}
