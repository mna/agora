package runtime

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/agora/bytecode"
)

type Module interface {
	ID() string
	Run() (Val, error)
}

type NativeModule interface {
	Module
	SetCtx(*Ctx)
}

type agoraModule struct {
	id  string
	fns []*AgoraFunc
	v   Val
}

func newAgoraModule(f *bytecode.File, c *Ctx) *agoraModule {
	m := &agoraModule{
		id: f.Name,
	}
	m.fns = make([]*AgoraFunc, len(f.Fns))
	for i, fn := range f.Fns {
		af := newAgoraFunc(m, c)
		af.name = fn.Header.Name
		af.stackSz = fn.Header.StackSz
		af.expArgs = fn.Header.ExpArgs
		af.expVars = fn.Header.ExpVars
		// TODO : Ignore LineStart and LineEnd at the moment, unused.
		m.fns[i] = af
		af.kTable = make([]Val, len(fn.Ks))
		for j, k := range fn.Ks {
			switch k.Type {
			case bytecode.KtBoolean:
				af.kTable[j] = Bool(k.Val.(int64) != 0)
			case bytecode.KtInteger:
				af.kTable[j] = Int(k.Val.(int64))
			case bytecode.KtFloat:
				af.kTable[j] = Float(k.Val.(float64))
			case bytecode.KtString:
				af.kTable[j] = String(k.Val.(string))
			default:
				panic("invalid constant value type")
			}
		}
		af.code = make([]bytecode.Instr, len(fn.Is))
		for j, ins := range fn.Is {
			af.code[j] = ins
		}
	}
	return m
}

func (m *agoraModule) Run() (v Val, err error) {
	defer PanicToError(&err)
	if len(m.fns) == 0 {
		return Nil, ErrModuleHasNoFunc
	}
	// Do not re-run a module if it has already been imported. Use the cached value.
	if m.v == nil {
		m.v = m.fns[0].Call(nil)
	}
	return m.v, nil
}

func PanicToError(err *error) {
	if p := recover(); p != nil {
		if e, ok := p.(error); ok {
			*err = e
		} else {
			*err = fmt.Errorf("%s", p)
		}
	}
}

func (m *agoraModule) ID() string {
	return m.id
}

type ModuleResolver interface {
	Resolve(string) (io.Reader, error)
}

type FileResolver struct{}

var (
	extensions = [...]string{".agorac", ".agoraa", ".agora"}
)

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
	// If there is no extension, try files in the following order:
	// 1- .agorac (compiled bytecode)
	// 2- .agoraa (agora assembly code)
	// 3- .agora  (agora source code)
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
