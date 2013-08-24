package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/agora/compiler"
	"github.com/jessevdk/go-flags"
)

type asm struct {
	Output string `short:"o" long:"output" description:"output file"`
}

func (a *asm) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected an input file name")
	}
	var f *os.File
	var err error
	f = os.Stdout
	if a.Output != "" {
		f, err = os.Create(a.Output)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	inf, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer inf.Close()
	b, err := new(compiler.Asm).Compile(args[0], inf)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

type dasm struct {
	Output string `short:"o" long:"output" description:"output file"`
}

func (d *dasm) Execute(args []string) error {
	return nil
}

func main() {
	a, d := new(asm), new(dasm)
	p := flags.NewParser(nil, flags.Default)
	p.AddCommand("asm", "assembler", "compile source assembler to bytecode", a)
	p.AddCommand("dasm", "disassembler", "disassemble bytecode to source assembly", d)
	// In case of errors, usage text is automatically displayed. In case of
	// success, the Execute() method of the matching command is called.
	p.Parse()
}
