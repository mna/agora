package runtime

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/agora/bytecode"
)

type EmptyModuleError string

func (e EmptyModuleError) Error() string {
	return string(e)
}
func NewEmptyModuleError(id string) EmptyModuleError {
	return EmptyModuleError(fmt.Sprintf("empty module: %s", id))
}

// The Module interface defines the required behaviours for a Module.
type Module interface {
	ID() string
	Run(...Val) (Val, error)
}

// A NativeModule is a Module with added behaviour required for supporting
// native Go modules.
type NativeModule interface {
	Module
	SetCtx(*Ctx)
}

type agoraModule struct {
	id  string
	fns []*agoraFunc
	v   Val
}

func newAgoraModule(f *bytecode.File, c *Ctx) *agoraModule {
	m := &agoraModule{
		id: f.Name,
	}
	// Define all functions
	m.fns = make([]*agoraFunc, len(f.Fns))
	for i, fn := range f.Fns {
		af := newAgoraFunc(m, c)
		af.name = fn.Header.Name
		af.stackSz = fn.Header.StackSz
		af.expArgs = fn.Header.ExpArgs
		if i != 0 { // No parent for root func
			af.parent = m.fns[fn.Header.ParentFnIx]
		}
		// TODO : Ignore LineStart and LineEnd at the moment, unused.
		m.fns[i] = af
		af.kTable = make([]Val, len(fn.Ks))
		for j, k := range fn.Ks {
			switch k.Type {
			case bytecode.KtBoolean:
				af.kTable[j] = Bool(k.Val.(int64) != 0)
			case bytecode.KtInteger:
				af.kTable[j] = Number(k.Val.(int64))
			case bytecode.KtFloat:
				af.kTable[j] = Number(k.Val.(float64))
			case bytecode.KtString:
				af.kTable[j] = String(k.Val.(string))
			default:
				panic("invalid constant value type")
			}
		}
		af.lTable = make([]string, len(fn.Ls))
		for j, l := range fn.Ls {
			af.lTable[j] = af.kTable[l].String()
		}
		af.code = make([]bytecode.Instr, len(fn.Is))
		for j, ins := range fn.Is {
			af.code[j] = ins
		}
	}
	return m
}

// Run executes the module and returns its return value, or an error.
func (m *agoraModule) Run(args ...Val) (v Val, err error) {
	defer PanicToError(&err)
	if len(m.fns) == 0 {
		return Nil, NewEmptyModuleError(m.ID())
	}
	// Do not re-run a module if it has already been imported. Use the cached value.
	if m.v == nil {
		fn := m.fns[0]
		fn.ctx.pushModule(m.ID())
		defer fn.ctx.popModule(m.ID())
		m.v = m.fns[0].Call(nil, args...)
	}
	return m.v, nil
}

// PanicToError is a utility function for modules implementations to catch panics
// and translate them to an error interface. It should be called in a defer statement,
// with the address of an error variable (usually a named return value) as argument.
func PanicToError(err *error) {
	if p := recover(); p != nil {
		if e, ok := p.(error); ok {
			*err = e
		} else {
			*err = fmt.Errorf("%s", p)
		}
	}
}

// ID returns the identifier of the module.
func (m *agoraModule) ID() string {
	return m.id
}

// A ModuleResolver interface represents the required behaviour for the component
// responsible for matching a module identifier to actual source code.
// Various implementations can be provided, for example by loading modules
// in a database, over http, compressed, secured and signed, etc.
type ModuleResolver interface {
	Resolve(string) (io.Reader, error)
}

// A FileResolver is a ModuleResolver that turns the module identifier into
// a file path to find the matching source code.
type FileResolver struct{}

var (
	extensions = [...]string{".agorac", ".agoraa", ".agora"}
)

// Resolve matches the provided identifier with a source file.
//
// If the identifier has no extension (which is recommended), Resolve looks
// for files in the following order:
//
// 1- .agorac (compiled bytecode)
// 2- .agoraa (agora assembly code)
// 3- .agora  (agora source code)
//
// TODO : This doesn't work, the Ctx has a single compiler, that may
// compile assembly or source, but not both. The Resolver should look
// for compiled bytecode or the same source code as the initial Ctx.Load.
func (f FileResolver) Resolve(id string) (io.Reader, error) {
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
	if filepath.Ext(nm) == "" {
		for _, ext := range extensions {
			if _, err := os.Stat(nm + ext); err != nil {
				if !os.IsNotExist(err) {
					return nil, err
				}
			} else {
				nm += ext
				break
			}
		}
	}
	return os.Open(nm)
}
