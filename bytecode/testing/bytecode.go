package testing

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var (
	// Useful prefilled expected byte slices
	ExpSig       = []byte{0x2A, 0x60, 0x0A, 0x00}
	ExpZeroInt64 = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

func SigVer(maj, min int) []byte {
	return append(ExpSig, byte(maj<<4)|byte(min))
}

func AppendAny(b []byte, vals ...interface{}) []byte {
	for i, v := range vals {
		switch v := v.(type) {
		case []byte:
			b = append(b, v...)
		case byte:
			b = append(b, v)
		case int32:
			b = append(b, byte(v))
		case int:
			b = append(b, byte(v))
		default:
			panic(fmt.Sprintf("invalid type to append at pos %d", i))
		}
	}
	return b
}

func Int64ToByteSlice(i int64) []byte {
	buf := bytes.NewBuffer(nil)
	if err := binary.Write(buf, binary.LittleEndian, i); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func UInt64ToByteSlice(u uint64) []byte {
	buf := bytes.NewBuffer(nil)
	if err := binary.Write(buf, binary.LittleEndian, u); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
