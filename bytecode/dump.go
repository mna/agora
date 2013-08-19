package bytecode

import (
	"io"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func catch(err *error) {
	if p := recover(); p != nil {
		if e, ok := p.(error); ok {
			*err = e
		} else {
			panic(p)
		}
	}
}
func (enc *Encoder) Encode(f *File) (err error) {
	defer catch(&err)

	return nil
}
