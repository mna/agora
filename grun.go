package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goblin/compiler"
	"github.com/PuerkitoBio/goblin/runtime/stdlib"
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
	// Register the standard lib's Fmt package
	ctx.RegisterNativeFuncs(stdlib.Fmt)
	// Run the main function
	res := ctx.Run()
	fmt.Printf("PASS - %v", res)
}
