package parser

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	src := []byte(`
import abc "this/is/a/path"
`)

	p := new(Parser)
	a, err := p.Parse("test", src)
	if err != nil {
		panic(err)
	}
	fmt.Println(a)
}
