package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goblin/compiler"
	"github.com/PuerkitoBio/goblin/runtime"
	"github.com/PuerkitoBio/goblin/runtime/stdlib"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: grun FILE")
		os.Exit(1)
	}
	f := os.Args[1]
	ctx := runtime.NewCtx(new(runtime.FileResolver), new(compiler.Asm))
	// Register the standard lib's Fmt package
	ctx.RegisterModule(new(stdlib.FmtMod))
	res, err := ctx.Load(f)
	if err != nil {
		fmt.Printf("FAIL - %s", err)
	} else {
		fmt.Printf("PASS - %v", res)
	}
}
