package nfi

import (
	"io"

	"github.com/PuerkitoBio/goblin/runtime"
)

type Streams interface {
	Stdout() io.ReadWriter
	Stdin() io.ReadWriter
	Stderr() io.ReadWriter
}

type NativeFunc func(Streams, ...runtime.Val) runtime.Val
