package runtime

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/PuerkitoBio/agora/bytecode"
)

var (
	ErrModuleNotFound  = errors.New("module not found")
	ErrModuleHasNoFunc = errors.New("module has no function")
	ErrCyclicDepFound  = errors.New("cyclic module dependency found")
)

type Compiler interface {
	Compile(string, io.Reader) (*bytecode.File, error)
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
	Resolver ModuleResolver // The module loading resolver (match a module to a string literal)
	Compiler Compiler       // The source code compiler
	Debug    bool           // Debug mode outputs helpful messages

	// Call stack
	frames []*frame
	frmsp  int

	// Modules management
	loadingMods map[string]bool // Modules currently being loaded
	loadedMods  map[string]Module
	builtin     *Object
}

func NewCtx(resolver ModuleResolver, comp Compiler) *Ctx {
	c := &Ctx{
		Stdout:      os.Stdout,
		Stdin:       os.Stdin,
		Stderr:      os.Stderr,
		Logic:       defaultLogic{},
		Resolver:    resolver,
		Compiler:    comp,
		loadingMods: make(map[string]bool),
		loadedMods:  make(map[string]Module),
	}
	b := new(builtinMod)
	b.SetCtx(c)
	if v, err := b.Run(); err != nil {
		panic("error loading angora builtin module: " + err.Error())
	} else {
		c.builtin = v.(*Object)
	}
	return c
}

/*
Sequence for loading, compiling, and bootstrapping execution:

* Get or create a Ctx (DefaultCtx or NewCtx())
* ctx.LoadFile(id string) (Val, error)
* If module is cached (ctx.loadedMods), return the Val, done.
* If module is native (ctx.nativeMods), call Module.Load(ctx), cache and return the value, done.
* If module is not cached, call ModuleResolver.Resolve(id string) (io.Reader, error)
* If Resolve returns an error, return nil, error, done.
* If file is already bytecode, just decode
* Otherwise call Compiler.Compile(id string, r io.Reader) (*bytecode.File, error)
* If Compile returns an error, return nil, error, done.
* Create module from *bytecode.File
* Cache module and return, do NOT execute the module.
*/
func (c *Ctx) Load(id string) (Module, error) {
	if id == "" {
		return nil, ErrModuleNotFound
	}
	if c.loadingMods[id] {
		return nil, ErrCyclicDepFound
	}
	c.loadingMods[id] = true
	defer delete(c.loadingMods, id)
	// If already loaded, return from cache
	if m, ok := c.loadedMods[id]; ok {
		return m, nil
	}
	// Else, resolve the matching file from the module id
	r, err := c.Resolver.Resolve(id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rc, ok := r.(io.ReadCloser); ok {
			rc.Close()
		}
	}()
	// If already bytecode, just decode
	var f *bytecode.File
	if rs, ok := r.(io.ReadSeeker); ok && bytecode.IsBytecode(rs) {
		// TODO : Eventually come up with a better solution, or at least a
		// failover if r is not a ReadSeeker.
		dec := bytecode.NewDecoder(r)
		f, err = dec.Decode()
	} else {
		// Compile to bytecode
		f, err = c.Compiler.Compile(id, r)
	}
	if err != nil {
		return nil, err
	}
	mod := newAgoraModule(f, c)
	// cache and return
	c.loadedMods[id] = mod
	return mod, nil
}

func (c *Ctx) RegisterNativeModule(m NativeModule) {
	m.SetCtx(c)
	c.loadedMods[m.ID()] = m
}

func (c *Ctx) push(f Func, fvm *funcVM) {
	// Stack has to grow as needed
	if c.frmsp == len(c.frames) {
		if c.Debug && c.frmsp == cap(c.frames) {
			fmt.Fprintf(c.Stdout, "DEBUG expanding frames of ctx, current size: %d\n", len(c.frames))
		}
		c.frames = append(c.frames, &frame{f, fvm})
	} else {
		c.frames[c.frmsp] = &frame{f, fvm}
	}
	c.frmsp++
}

func (ø *Ctx) pop() {
	ø.frmsp--
	ø.frames[ø.frmsp] = nil // free this reference for gc
}

func (c *Ctx) getVar(nm string) (Val, bool) {
	// Current frame is c.frmsp - 1
	for i := c.frmsp - 1; i >= 0; i-- {
		frm := c.frames[i]
		if frm.fvm != nil {
			if v, ok := frm.fvm.vars[nm]; ok {
				return v, true
			}
		}
	}
	// Finally, look if the identifier refers to a built-in function.
	// This will return Nil if it doesn't match any built-in.
	b := c.builtin.Get(String(nm))
	return b, b != Nil
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
