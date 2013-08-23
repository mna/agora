package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/PuerkitoBio/goblin/bytecode"
)

var (
	ErrModuleNotFound  = errors.New("module not found")
	ErrModuleHasNoFunc = errors.New("module has no function")
	ErrCyclicDepFound  = errors.New("cyclic module dependency found")
)

type Compiler interface {
	Compile(string, io.Reader) ([]byte, error)
}

type frame struct {
	f   Func
	fvm *funcVM
}

type Ctx struct {
	// Public fields
	Stdout   io.ReadWriter  // The standard streams
	Stdin    io.ReadWriter  // ...
	Stderr   io.ReadWriter  // ...
	Logic    LogicProcessor // The boolean logic processor (And, Or, Not)
	Resolver ModuleResolver
	Compiler Compiler

	// Call stack
	frames []*frame
	frmsp  int

	// Modules management
	loadingMods map[string]bool   // Modules currently being loaded
	loadedMods  map[string]Val    // Modules export a Val
	nativeMods  map[string]Module // List of available native modules
}

func NewCtx(resolver ModuleResolver, comp Compiler) *Ctx {
	return &Ctx{
		Stdout:      os.Stdout,
		Stdin:       os.Stdin,
		Stderr:      os.Stderr,
		Logic:       defaultLogic{},
		Resolver:    resolver,
		Compiler:    comp,
		loadingMods: make(map[string]bool),
		loadedMods:  make(map[string]Val),
		nativeMods:  make(map[string]Module),
	}
}

/*
Sequence for loading, compiling, and bootstrapping execution:

* Get or create a Ctx (DefaultCtx or NewCtx())
* ctx.LoadFile(id string) (Val, error)
* If module is cached (ctx.loadedMods), return the Val, done.
* If module is native (ctx.nativeMods), call Module.Load(ctx), cache and return the value, done.
* If module is not cached, call ModuleResolver.Resolve(id string) (io.Reader, error)
* If Resolve returns an error, return nil, error, done.
* TODO : if file is already bytecode, just decode
* Call Compiler.Compile(id string, r io.Reader) ([]byte, error)
* If Compile returns an error, return nil, error, done.
* Call Undump(b) (Module, error)
* If Undump returns an error, return nil, error, done.
* Call Module.Load(ctx), cache and return the value, done.
*/
func (ø *Ctx) Load(id string) (Val, error) {
	if id == "" {
		return nil, ErrModuleNotFound
	}
	if ø.loadingMods[id] {
		return nil, ErrCyclicDepFound
	}
	ø.loadingMods[id] = true
	defer delete(ø.loadingMods, id)

	// If already loaded, return from cache
	if v, ok := ø.loadedMods[id]; ok {
		return v, nil
	}
	// If native module, get from native table
	if m, ok := ø.nativeMods[id]; ok {
		loaded := m.Load(ø)
		ø.loadedMods[id] = loaded
		return loaded, nil
	}
	// Else, resolve the matching file from the module id
	r, err := ø.Resolver.Resolve(id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rc, ok := r.(io.ReadCloser); ok {
			rc.Close()
		}
	}()
	// TODO : If already bytecode, skip
	// Compile to bytecode
	b, err := ø.Compiler.Compile(id, r)
	if err != nil {
		return nil, err
	}
	// TODO : Compile should return *bytecode.File, makes no sense to return bytecode,
	// then decode it back to *bytecode.File...
	// Load the bytecode in memory
	f, err := bytecode.NewDecoder(bytes.NewReader(b)).Decode()
	if err != nil {
		return nil, err
	}
	gm := newGoblinModule(f)
	// Load the module, cache and return
	loaded := gm.Load(ø)
	ø.loadedMods[id] = loaded
	return loaded, nil
}

func (ø *Ctx) RegisterNativeModule(m Module) {
	ø.nativeMods[m.ID()] = m
}

func (ø *Ctx) push(f Func, fvm *funcVM) {
	// Stack has to grow as needed
	if ø.frmsp == len(ø.frames) {
		if ø.frmsp == cap(ø.frames) {
			fmt.Fprintf(ø.Stdout, "DEBUG expanding frames of ctx, current size: %d\n", len(ø.frames))
		}
		ø.frames = append(ø.frames, &frame{f, fvm})
	} else {
		ø.frames[ø.frmsp] = &frame{f, fvm}
	}
	ø.frmsp++
}

func (ø *Ctx) pop() {
	ø.frmsp--
	ø.frames[ø.frmsp] = nil // free this reference for gc
}

func (ø *Ctx) getVar(nm string) (Val, bool) {
	// Current frame is ø.frmsp - 1
	for i := ø.frmsp - 1; i >= 0; i-- {
		frm := ø.frames[i]
		if frm.fvm != nil {
			if v, ok := frm.fvm.vars[nm]; ok {
				return v, true
			}
		}
	}
	return Nil, false
}

func (ø *Ctx) setVar(nm string, v Val) bool {
	// Current frame is ø.frmsp - 1
	for i := ø.frmsp - 1; i >= 0; i-- {
		frm := ø.frames[i]
		if frm.fvm != nil {
			if _, ok := frm.fvm.vars[nm]; ok {
				frm.fvm.vars[nm] = v
				return true
			}
		}
	}
	return false
}

func (ø *Ctx) dump(n int) {
	if n < 0 {
		return
	}
	for i, cnt := ø.frmsp, ø.frmsp-n; i > 0 && i > cnt; i-- {
		fmt.Fprintf(ø.Stdout, "\n[Frame %3d]\n===========", i-1)
		if frm := ø.frames[i-1]; frm.fvm != nil {
			fmt.Fprintln(ø.Stdout, frm.fvm.dump())
		} else {
			fmt.Fprintln(ø.Stdout, frm.f.(dumper).dump())
		}
	}
}
