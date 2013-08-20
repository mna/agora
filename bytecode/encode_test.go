package bytecode

import (
	"bytes"
	"fmt"
	"testing"
)

var (
	cases = []struct {
		maj int
		min int
		f   *File
		exp []byte
		err error
	}{
		0: {
			f:   &File{},
			exp: []byte{0x2A, 0x60, 0x0A, 0x00, 0x00},
		},
		1: {
			maj: 1,
			min: 2,
			f:   &File{MajorVersion: 1, MinorVersion: 2},
			exp: []byte{0x2A, 0x60, 0x0A, 0x00, 0x12},
		},
	}

	isolateCase = -1
)

func TestEncode(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	for i, c := range cases {
		if isolateCase >= 0 && isolateCase != i {
			continue
		}
		if testing.Verbose() {
			fmt.Printf("testing case %d...\n", i)
		}

		// Arrange
		_MAJOR_VERSION = c.maj
		_MINOR_VERSION = c.min
		buf.Reset()

		// Act
		err := enc.Encode(c.f)

		// Assert
		if err != c.err {
			if c.err == nil {
				t.Errorf("expected no error, got `%s`", err)
			} else {
				t.Errorf("expected error `%s`, got `%s`", c.err, err)
			}
		}
		if c.exp != nil {
			got := buf.Bytes()
			if bytes.Compare(got, c.exp) != 0 {
				t.Errorf("expected %x, got %x", c.exp, got)
			}
		}
	}
}
