package runtime

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrModuleNotFound = errors.New("module not found")
)

type Module interface {
	Load(*Ctx) Val
}

type ModuleResolver interface {
	Resolve(string) (io.Reader, error)
}

type FileResolver struct{}

func (ø FileResolver) Resolve(id string) (io.Reader, error) {
	var nm string
	if filepath.IsAbs(id) {
		nm = id
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		nm = filepath.Join(pwd, id)
	}
	if !strings.HasSuffix(nm, ".goblin") {
		nm += ".goblin"
	}
	return os.Open(nm)
}

type Compiler interface {
	Compile(string, io.Reader) (Module, error)
}

/*
Sequence for loading, compiling, and bootstrapping execution:

* Get or create a Ctx (DefaultCtx or NewCtx())
* ctx.LoadFile(id string) (Val, error)
* If module is cached (ctx.loadedMods), return the Val, done.
* If module is native (ctx.nativeMods), call Module.Load(ctx), cache and return the value, done.
* If module is not cached, call ModuleResolver.Resolve(id string) (io.Reader, error)
* If Resolve returns an error, return nil, error, done.
* Call Compiler.Compile(id string, r io.Reader) (Module, error)
* If Compile returns an error, return nil, error, done.
* Call Module.Load(ctx), cache and return the value, done.
*/

type Ctx struct {
	// Public fields
	Protos   []Func
	Stdout   io.ReadWriter  // The standard streams
	Stdin    io.ReadWriter  // ...
	Stderr   io.ReadWriter  // ...
	Logic    LogicProcessor // The boolean logic processor (And, Or, Not)
	Resolver ModuleResolver
	Compiler Compiler

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

func NewCtx(resolver ModuleResolver, comp Compiler) *Ctx {
	return &Ctx{
		Stdout:     os.Stdout,
		Stdin:      os.Stdin,
		Stderr:     os.Stderr,
		Logic:      defaultLogic{},
		Resolver:   resolver,
		Compiler:   comp,
		loadedMods: make(map[string]Val),
		nativeMods: make(map[string]Module),
	}
}

func (ø *Ctx) Load(id string) (Val, error) {
	if id == "" {
		return nil, ErrModuleNotFound
	}
	if v, ok := ø.loadedMods[id]; ok {
		return v, nil
	}
	if m, ok := ø.nativeMods[id]; ok {
		ø.loadedMods[id] = m.Load(ø)
		return ø.loadedMods[id], nil
	}
	r, err := ø.Resolver.Resolve(id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rc, ok := r.(io.ReadCloser); ok {
			rc.Close()
		}
	}()
	m, err := ø.Compiler.Compile(id, r)
	if err != nil {
		return nil, err
	}
	ø.loadedMods[id] = m.Load(ø)
	return ø.loadedMods[id], nil
}

func (ø *Ctx) RegisterModule(id string, m Module) {
	ø.nativeMods[id] = m
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
