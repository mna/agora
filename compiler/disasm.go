package compiler

import (
	"fmt"
	"io"
	"strconv"

	"github.com/PuerkitoBio/goblin/bytecode"
)

var (
	maj, min      = bytecode.Version()
	disasmComment = fmt.Sprintf("// Generated from the disassembler, v%d.%d", maj, min)
)

type Disasm struct {
	w   io.Writer
	err error
}

func (d *Disasm) Uncompile(r io.Reader, w io.Writer) error {
	f, err := bytecode.NewDecoder(r).Decode()
	if err != nil {
		return err
	}
	d.w = w
	d.err = nil
	// 1- Write the standard comment
	d.write(disasmComment, true)
	// 2- Write every function
	for _, fn := range f.Fns {
		d.write("[f]", true)
		d.write(fn.Header.Name, true)
		d.write(fn.Header.StackSz, true)
		d.write(fn.Header.ExpArgs, true)
		d.write(fn.Header.ExpVars, true)
		d.write(fn.Header.LineStart, true)
		d.write(fn.Header.LineEnd, true)

		// 3- Write the function's K section
		d.write("[k]", true)
		for _, k := range fn.Ks {
			d.write(k.Type, false)
			d.write(k.Val, true)
		}
		// 4- Write the function's I section
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
