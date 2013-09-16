package compiler

import (
	"fmt"
	"io"
	"strconv"

	"github.com/PuerkitoBio/agora/bytecode"
)

var (
	maj, min      = bytecode.Version()
	disasmComment = fmt.Sprintf("// Generated from the disassembler, v%d.%d", maj, min)
)

// A Disasm translates a bytecode representation into assembly source code.
type Disasm struct {
	w   io.Writer
	err error
}

// ToAsm takes the in-memory bytecode File structure and translates it to
// assembly source code, writing the results to the provided writer. If an
// error is encountered, it is returned, otherwise it returns nil.
func (d *Disasm) ToAsm(f *bytecode.File, w io.Writer) error {
	d.w = w
	d.err = nil
	// 1- Write the standard comment
	d.write(disasmComment, true)
	// 2- Write every function
	for _, fn := range f.Fns {
		d.write("[f]", true)
		// If the func name is empty, set it to <anon>
		if fn.Header.Name == "" {
			d.write("<anon>", true)
		} else {
			d.write(fn.Header.Name, true)
		}
		d.write(fn.Header.StackSz, true)
		d.write(fn.Header.ExpArgs, true)
		d.write(fn.Header.ParentFnIx, true)
		d.write(fn.Header.LineStart, true)
		d.write(fn.Header.LineEnd, true)

		// 3- Write the function's K section
		d.write("[k]", true)
		for _, k := range fn.Ks {
			d.write(k.Type, false)
			d.write(k.Val, true)
		}
		// 4- Write the function's L section
		d.write("[l]", true)
		for _, l := range fn.Ls {
			d.write(l, true)
		}
		// 5- Write the function's I section
		d.write("[i]", true)
		for _, i := range fn.Is {
			op, flg, ix := i.Opcode(), i.Flag(), i.Index()
			d.write(op.String(), false)
			d.write(" ", false)
			d.write(flg.String(), false)
			d.write(" ", false)
			d.write(ix, true)
		}
	}
	return d.err
}

// Uncompile reads the bytecode source data from the provided reader, and translates
// it to assembly source code written into the writer. If an error is encountered, it
// is returned, otherwise it returns nil.
func (d *Disasm) Uncompile(r io.Reader, w io.Writer) error {
	f, err := bytecode.NewDecoder(r).Decode()
	if err != nil {
		return err
	}
	return d.ToAsm(f, w)
}

func (d *Disasm) write(i interface{}, newLine bool) {
	if d.err != nil {
		return
	}
	switch v := i.(type) {
	case int64:
		d.write(strconv.FormatInt(v, 10), newLine)
	case uint64:
		d.write(strconv.FormatUint(v, 10), newLine)
	case float64:
		d.write(strconv.FormatFloat(v, 'f', -1, 64), newLine)
	case bytecode.KType:
		d.write(string(v), newLine)
	case string:
		_, d.err = io.WriteString(d.w, v)
		if newLine {
			d.write("\n", false)
		}
	default:
		panic(fmt.Sprintf("unexpected type to write: %T", i))
	}
}
