package runtime

import (
	"io"
	"os"
)

type Ctx struct {
	Protos []*FuncProto
	Stdout io.Writer
	Stdin  io.Reader
	Stderr io.Writer

	// Native funcs table
	nTable map[string]func(...Val) Val
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
