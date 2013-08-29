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
	gm := &agoraModule{
		id: f.Name,
	}
	gm.fns = make([]*AgoraFunc, len(f.Fns))
	for i, fn := range f.Fns {
		gf := newAgoraFunc(gm, c)
		gf.name = fn.Header.Name
		gf.stackSz = fn.Header.StackSz
		gf.expArgs = fn.Header.ExpArgs
		gf.expVars = fn.Header.ExpVars
		// TODO : Ignore LineStart and LineEnd at the moment, unused.
		gm.fns[i] = gf
		gf.kTable = make([]Val, len(fn.Ks))
		for j, k := range fn.Ks {
			switch k.Type {
			case bytecode.KtBoolean:
				gf.kTable[j] = Bool(k.Val.(int64) != 0)
			case bytecode.KtInteger:
				gf.kTable[j] = Int(k.Val.(int64))
			case bytecode.KtFloat:
				gf.kTable[j] = Float(k.Val.(float64))
			case bytecode.KtString:
				gf.kTable[j] = String(k.Val.(string))
			default:
				panic("invalid constant value type")
			}
		}
		gf.code = make([]bytecode.Instr, len(fn.Is))
		for j, ins := range fn.Is {
			gf.code[j] = ins
		}
	}
	return gm
}

func (g *agoraModule) Run() (v Val, err error) {
	defer PanicToError(&err)
	if len(g.fns) == 0 {
		return Nil, ErrModuleHasNoFunc
	}
	if g.v == nil {
		g.v = g.fns[0].Call(nil)
	}
	return g.v, nil
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

func (g *agoraModule) ID() string {
	return g.id
}

type ModuleResolver interface {
	Resolve(string) (io.Reader, error)
}

type FileResolver struct{}

var (
	extensions = [...]string{".agorac", ".agoraa", ".agora"}
)

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
