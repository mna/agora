package runtime

import (
	"io"
)

type Streams interface {
	Stdout() io.ReadWriter
	Stdin() io.ReadWriter
	Stderr() io.ReadWriter
}

type NativeFunc func(Streams, ...Val) Val
