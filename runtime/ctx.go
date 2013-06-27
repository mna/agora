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
	Protos []Func
	Stdout io.ReadWriter  // The standard streams
	Stdin  io.ReadWriter  // ...
	Stderr io.ReadWriter  // ...
	Logic  LogicProcessor // The boolean logic processor (And, Or, Not)

	// Call stack
	callstack []Func
	callsp    int

	// Nested scopes - like call stack, but only for funcVM.run() calls
	scopes  []*funcVM
	scopesp int

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
	f := ø.Protos[0]
	return f.Call().Native()
}

func RegisterModule(nm string, m Module) {
	DefaultCtx.RegisterModule(nm, m)
}

func (ø *Ctx) RegisterModule(nm string, m Module) {
	ø.nativeMods[nm] = m
}

func (ø *Ctx) push(f Func, fvm *funcVM) {
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

	if fvm != nil { // fvm may be nil, if f is native
		if ø.scopesp == len(ø.scopes) {
			if ø.scopesp == cap(ø.scopes) {
				fmt.Printf("DEBUG expanding scopes of ctx, current size: %d\n", len(ø.scopes))
			}
			ø.scopes = append(ø.scopes, fvm)
		} else {
			ø.scopes[ø.scopesp] = fvm
		}
		ø.scopesp++
	}
}

func (ø *Ctx) pop(scopeToo bool) {
	ø.callsp--
	ø.callstack[ø.callsp] = nil // free this reference for gc

	if scopeToo {
		ø.scopesp--
		ø.scopes[ø.scopesp] = nil // free this reference for gc
	}
}

func (ø *Ctx) getVar(nm string) (Val, bool) {
	// Current scope is ø.scopesp - 1
	for i := ø.scopesp - 1; i >= 0; i-- {
		f := ø.scopes[i]
		if v, ok := f.vars[nm]; ok {
			return v, true
		}
	}
	return Nil, false
}

func (ø *Ctx) setVar(nm string, v Val) bool {
	// Current scope is ø.scopesp - 1
	for i := ø.scopesp - 1; i >= 0; i-- {
		f := ø.scopes[i]
		if _, ok := f.vars[nm]; ok {
			f.vars[nm] = v
			return true
		}
	}
	return false
}

func (ø *Ctx) dump(n int) {
	if n < 0 {
		return
	}
	for i, cnt := ø.callsp, ø.callsp-n; i > 0 && i > cnt; i-- {
		fmt.Fprintf(ø.Stdout, "[Call Stack %3d]\n===============\n", i-1)
		fmt.Fprint(ø.Stdout, ø.callstack[i-1].(dumper).dump())
	}
}
