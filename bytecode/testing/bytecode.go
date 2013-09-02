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

func EqualFile(f1, f2 *File) bool {
	if f1 == nil && f2 == nil {
		return true
	}
	if f1 == nil || f2 == nil {
		return false
	}
	if f1.Name != f2.Name {
		return false
	}
	if f1.MajorVersion != f2.MajorVersion {
		return false
	}
	if f1.MinorVersion != f2.MinorVersion {
		return false
	}
	if len(f1.Fns) != len(f2.Fns) {
		return false
	}
	for i := 0; i < len(f1.Fns); i++ {
		fn1, fn2 := f1.Fns[i], f2.Fns[i]
		if fn1.Header.Name != fn2.Header.Name {
			return false
		}
		if fn1.Header.StackSz != fn2.Header.StackSz {
			return false
		}
		if fn1.Header.ExpArgs != fn2.Header.ExpArgs {
			return false
		}
		if fn1.Header.ExpVars != fn2.Header.ExpVars {
			return false
		}
		if fn1.Header.LineStart != fn2.Header.LineStart {
			return false
		}
		if fn1.Header.LineEnd != fn2.Header.LineEnd {
			return false
		}
		if len(fn1.Ks) != len(fn2.Ks) {
			return false
		}
		for j := 0; j < len(fn1.Ks); j++ {
			k1, k2 := fn1.Ks[j], fn2.Ks[j]
			if k1.Type != k2.Type {
				return false
			}
			if k1.Val != k2.Val {
				return false
			}
		}
		if len(fn1.Is) != len(fn2.Is) {
			return false
		}
		for j := 0; j < len(fn1.Is); j++ {
			if fn1.Is[j] != fn2.Is[j] {
				return false
			}
		}
	}
	return true
}
