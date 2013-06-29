package main

import (
	"fmt"
	"io"
	"os"

	"github.com/PuerkitoBio/goblin/compiler"
	"github.com/PuerkitoBio/goblin/runtime"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: asm FILE")
		os.Exit(1)
	}
	f := os.Args[1]
	res := new(runtime.FileResolver)
	r, err := res.Resolve(f)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(2)
	}
	if rc, ok := r.(io.ReadCloser); ok {
		defer rc.Close()
	}
	c := new(compiler.Asm)
	b, err := c.Compile(f, r)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(3)
	}
	fmt.Print(b)
}
