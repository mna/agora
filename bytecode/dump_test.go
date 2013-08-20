package bytecode

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	_MAJOR_VERSION = 0
	_MINOR_VERSION = 0

	f := &File{
		Name:         "test",
		MajorVersion: 0,
		MinorVersion: 0,
	}
	// Little endian...
	exp := []byte{
		0x2A,
		0x60,
		0x0A,
		0x00,
		0x00,
	}

	buf := bytes.NewBuffer(nil)
	err := NewEncoder(buf).Encode(f)
	if err != nil {
		t.Error(err)
	} else if bytes.Compare(buf.Bytes(), exp) != 0 {
		t.Errorf("expected %x, got %x", exp, buf.Bytes())
	}
}

func TestEncodeVersion(t *testing.T) {
	_MAJOR_VERSION = 1
	_MINOR_VERSION = 2

	f := &File{
		Name:         "test",
		MajorVersion: 1,
		MinorVersion: 2,
	}
	// Little endian...
	exp := []byte{
		0x2A,
		0x60,
		0x0A,
		0x00,
		0x12,
	}

	buf := bytes.NewBuffer(nil)
	err := NewEncoder(buf).Encode(f)
	if err != nil {
		t.Error(err)
	} else if bytes.Compare(buf.Bytes(), exp) != 0 {
		t.Errorf("expected %x, got %x", exp, buf.Bytes())
	}
}
