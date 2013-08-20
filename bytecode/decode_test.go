package bytecode

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

var (
	deccases = []struct {
		maj int
		min int
		src []byte
		exp *File
		err error
	}{
		0: {
			// Simplest case, encodes the file header only
			src: expSigAndDefVer,
			exp: &File{},
		},
	}

	isolateDecCase = -1
)

func TestDecode(t *testing.T) {
	for i, c := range deccases {
		if isolateDecCase >= 0 && isolateDecCase != i {
			continue
		}
		if testing.Verbose() {
			fmt.Printf("testing decode case %d...\n", i)
		}

		// Arrange
		_MAJOR_VERSION = c.maj
		_MINOR_VERSION = c.min

		// Act
		f, err := NewDecoder(bytes.NewBuffer(c.src)).Decode()

		// Assert
		if err != c.err {
			if c.err == nil {
				t.Errorf("[%d] - expected no error, got `%s`", i, err)
			} else {
				t.Errorf("[%d] - expected error `%s`, got `%s`", i, c.err, err)
			}
		}
		if c.exp != nil {
			if !reflect.DeepEqual(f, c.exp) {
				t.Errorf("[%d] - expected \n%#v\n, got \n%#v\n", i, c.exp, f)
			}
		}
		if c.err == nil && c.exp == nil {
			t.Errorf("[%d] - no assertion", i)
		}
	}
}
