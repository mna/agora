package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/PuerkitoBio/agora/bytecode"
	"github.com/PuerkitoBio/agora/compiler"
	"github.com/PuerkitoBio/agora/runtime"
	"github.com/PuerkitoBio/agora/runtime/stdlib"
	"github.com/jessevdk/go-flags"
)

// The assembler command struct.
type asm struct {
	Output string `short:"o" long:"output" description:"output file"`
	Hexa   bool   `short:"x" long:"hexadecimal" description:"print output as hexadecimal"`
}

// Execute the assembler command
func (a *asm) Execute(args []string) error {
	// Validate input
	if len(args) != 1 {
		return fmt.Errorf("expected an input file name")
	}
	// Open input file
	inf, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer inf.Close()
	// Compile to bytecode File
	f, err := new(compiler.Asm).Compile(args[0], inf)
	if err != nil {
		return err
	}
	// Write output
	var outF *os.File
	outF = os.Stdout
	if a.Output != "" {
		outF, err = os.Create(a.Output)
		if err != nil {
			return err
		}
		defer outF.Close()
	}
	// Encode to bytecode
	buf := bytes.NewBuffer(nil)
	err = bytecode.NewEncoder(buf).Encode(f)
	if err != nil {
		return err
	}
	if a.Hexa {
		_, err = io.WriteString(outF, fmt.Sprintf("%x", buf.Bytes()))
	} else {
		_, err = outF.Write(buf.Bytes())
	}
	if err != nil {
		return err
	}
	return nil
}

// The disassembler command struct
type dasm struct {
	Output string `short:"o" long:"output" description:"output file"`
}

// Execute the disassembler command
func (d *dasm) Execute(args []string) error {
	// Validate input
	if len(args) != 1 {
		return fmt.Errorf("expected an input file name")
	}
	// Open input file
	inf, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer inf.Close()
	// Open output file
	var f *os.File
	f = os.Stdout
	if d.Output != "" {
		f, err = os.Create(d.Output)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	// Compile to assembly
	return new(compiler.Disasm).Uncompile(inf, f)
}

// The run command struct
type run struct {
	Output   string `short:"o" long:"output" description:"output file"`
	FromAsm  bool   `short:"a" long:"from-asm" description:"run an assembly input"`
	NoStdlib bool   `short:"S" long:"nostdlib" description:"do not import the stdlib"`
}

func (r *run) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected an input file")
	}
	var c runtime.Compiler
	if r.FromAsm {
		c = new(compiler.Asm)
	} else {
		c = new(compiler.Compiler)
	}
	ctx := runtime.NewCtx(new(runtime.FileResolver), c)
	if !r.NoStdlib {
		// Register the standard lib's Fmt package
		ctx.RegisterNativeModule(new(stdlib.FmtMod))
	}
	res, err := ctx.Load(args[0])
	if err == nil {
		fmt.Printf("\n\n= %v\n", res)
	}
	return err
}

// The build command struct
type build struct{}

func main() {
	a, d, r := new(asm), new(dasm), new(run)
	p := flags.NewParser(nil, flags.Default)
	p.AddCommand("asm", "assembler", "compile source assembler to bytecode", a)
	p.AddCommand("dasm", "disassembler", "disassemble bytecode to source assembly", d)
	p.AddCommand("run", "run", "execute a source program", r)
	// In case of errors, usage text is automatically displayed. In case of
	// success, the Execute() method of the matching command is called.
	p.Parse()
}
