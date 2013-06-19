package main

import (
	"fmt"
	"github.com/PuerkitoBio/goblin/runtime"
	"github.com/davecgh/go-spew/spew"
	"os"
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
	fps := runtime.Asm(fr)
	// Run the function at index 0 (the main)
	if len(fps) == 0 {
		panic("no function in specified file")
	}
	spew.Dump(fps)
	fn := runtime.NewFunc(fps[0])
	ret := fn.Run()
	fmt.Printf("PASS - %v", ret)
}
