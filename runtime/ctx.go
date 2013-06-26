package runtime

import (
	"fmt"
	"io"
	"os"
)

var (
	DefaultCtx = NewCtx()
)

type Module interface {
	Load(*Ctx) Val
}

type Ctx struct {
	// Public fields
	FuncProtos []*FuncProto   // TODO : Move into module or unnecessary?...
	Stdout     io.ReadWriter  // The standard streams
	Stdin      io.ReadWriter  // ...
	Stderr     io.ReadWriter  // ...
	Logic      LogicProcessor // The boolean logic processor (And, Or, Not)

	// Call stack
	callstack []*Func
	callsp    int

	// Modules management
	loadedMods map[string]Val    // Modules export a Val
	nativeMods map[string]Module // List of available native modules
}

func NewCtx() *Ctx {
	return &Ctx{
		nil,
		os.Stdout,
		os.Stdin,
		os.Stderr,
		defaultLogic{},
		nil,
		0,
		make(map[string]Val),
		make(map[string]Module),
	}
}

func Run() interface{} {
	return DefaultCtx.Run()
}

func (ø *Ctx) Run() interface{} {
	if len(ø.Protos) == 0 {
		panic("no function available in this context")
	}
	f := newFunc(ø, ø.Protos[0])
	return f.Call().Native()
}

func RegisterModule(nm string, m Module) {
	DefaultCtx.RegisterModule(nm, m)
}

func (ø *Ctx) RegisterModule(nm string, m Module) {
	ø.nativeMods[nm] = m
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
