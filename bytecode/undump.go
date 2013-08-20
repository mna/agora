package bytecode

import (
	"io"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (dec *Decoder) Decode() (*File, error) {
	return nil, nil
}
