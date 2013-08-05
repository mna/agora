package compiler

import (
	"bytes"
	"github.com/PuerkitoBio/goblin/runtime"
	"os"
	"testing"
)

// TODO : This is not a test at all...
func TestAsm(t *testing.T) {
	f, err := os.Open("../runtime/testdata/05-nativefunc.goblin")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	a := &Asm{}
	b, err := a.Compile("05-nativefunc", f)
	if err != nil {
		panic(err)
	}
	t.Log(b)
	m, err := runtime.Undump(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	t.Log(m)
}
