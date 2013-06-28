package compiler

import (
	"strings"
	"testing"
)

func TestAsm(t *testing.T) {
	src := `
[f]
3
0
1
0
1
<main>
[k]
sfmt
sPrintln
sHello 
sworld
[i]
LOAD K 0 // <-fmt
POP V 0 // ->fmt
PUSH K 2 // <-"Hello "
PUSH K 3 // <-"world"
PUSH K 1 // <-"Println"
PUSH K 0 // <-"fmt"
CFLD A 2
PUSH N 0 // <-Nil
DUMP
RET
`

	a := &Asm{}
	b, err := a.Compile("test", strings.NewReader(src))
	if err != nil {
		panic(err)
	}
	t.Log(b)
}
