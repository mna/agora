package runtime

import (
	"io"
	"os"

	"github.com/PuerkitoBio/goblin/runtime/nfi"
)

type Ctx struct {
	Protos []*FuncProto
	stdout io.ReadWriter
	stdin  io.ReadWriter
	stderr io.ReadWriter

	// Native funcs table
	nTable map[string]nfi.NativeFunc
}

func NewCtx() *Ctx {
	return &Ctx{
		nil,
		os.Stdout,
		os.Stdin,
		os.Stderr,
		make(map[string]func(...Val) Val),
	}
}

func (ø *Ctx) RegisterNativeFuncs(fs map[string]nfi.NativeFunc) {
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
