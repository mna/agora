package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goblin/compiler"
	"github.com/PuerkitoBio/goblin/runtime"
	"github.com/PuerkitoBio/goblin/runtime/stdlib"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: grun FILE")
		os.Exit(1)
	}
	f := os.Args[1]
	fr, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	// Assemble the file
	ctx := compiler.Asm(fr)
	// Run the function at index 0 (the main)
	if len(ctx.Protos) == 0 {
		panic("no function in specified file")
	}
	ctx.RegisterNativeFuncs(stdlib.Fmt)
	// Print the function prototypes
	spew.Dump(ctx)

	// Execute the program
	fn := runtime.NewFunc(ctx, ctx.Protos[0])
	ret := fn.Call()

	// Print the resulting stack and variables
	spew.Dump(fn)

	fmt.Printf("PASS - %v", ret)
}
