package runtime

import (
	"io"
	"os"
)

type Ctx struct {
	Protos []*FuncProto
	stdout io.ReadWriter
	stdin  io.ReadWriter
	stderr io.ReadWriter
	logic  Logic

	// Native funcs table
	nTable map[string]NativeFunc
}

func NewCtx() *Ctx {
	return &Ctx{
		nil,
		os.Stdout,
		os.Stdin,
		os.Stderr,
		defaultLogic{},
		make(map[string]NativeFunc),
	}
}

func (ø *Ctx) SetStdStreams(stdin, stdout, stderr io.ReadWriter) {
	ø.stdin, ø.stdout, ø.stderr = stdin, stdout, stderr
}

func (ø *Ctx) Run() interface{} {
	if len(ø.Protos) == 0 {
		panic("no function available in this context")
	}
	f := NewFunc(ø, ø.Protos[0])
	return f.Call().Native()
}

func (ø *Ctx) RegisterNativeFuncs(fs map[string]NativeFunc) {
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
