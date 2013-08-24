package runtime

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/agora/bytecode"
)

type Module interface {
	ID() string
	Load(*Ctx) Val
}

type agoraModule struct {
	id  string
	fns []*AgoraFunc
}

func newAgoraModule(f *bytecode.File) *agoraModule {
	gm := &agoraModule{
		id: f.Name,
	}
	gm.fns = make([]*AgoraFunc, len(f.Fns))
	for i, fn := range f.Fns {
		gf := newAgoraFunc(gm)
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

func (g *agoraModule) Load(ctx *Ctx) Val {
	if len(g.fns) == 0 {
		panic(ErrModuleHasNoFunc)
	}
	for i, _ := range g.fns {
		g.fns[i].ctx = ctx
	}
	return g.fns[0].Call(nil)
}

func (g *agoraModule) ID() string {
	return g.id
}

type ModuleResolver interface {
	Resolve(string) (io.Reader, error)
}

type FileResolver struct{}

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
	if !strings.HasSuffix(nm, ".agora") {
		nm += ".agora"
	}
	return os.Open(nm)
}
