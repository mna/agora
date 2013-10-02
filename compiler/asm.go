package compiler

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/agora/bytecode"
)

var (
	// Predefined errors
	ErrInvalidInstruction = errors.New("invalid instruction")
	ErrNoInput            = errors.New("no input provided")
)

// An Asm is an assembly source code compiler. It implements the runtime.Compiler
// interface, so that it is suitable for runtime.Ctx.
type Asm struct {
	s     *bufio.Scanner
	f     *bytecode.File
	ended bool
	err   error
}

// Compile takes a module identifier and a reader, and compiles its assembly source
// code to an in-memory representation of agora bytecode, ready for execution.
// If an error is encounted, it is returned as second value, otherwise it is nil.
func (a *Asm) Compile(id string, r io.Reader) (*bytecode.File, error) {
	a.ended = false
	a.err = nil
	a.s = bufio.NewScanner(r)
	// Ignore everything before the [f] section
	a.findSection("[f]")
	// Edge case: if no func section (empty input), don't create the File, return
	if a.ended {
		return nil, ErrNoInput
	}
	a.f = bytecode.NewFile(id)
	// Read the first section, the other ones get called recursively as needed
	a.readFn()
	return a.f, a.err
}

func (a *Asm) findSection(s string) {
	for line, ok := a.getLine(false); ok; line, ok = a.getLine(false) {
		if line == s {
			break
		}
	}
}

func (a *Asm) readFn() {
	fn := new(bytecode.Fn)
	fn.Header.Name, _ = a.getLine(false)
	fn.Header.StackSz = a.getInt64()
	fn.Header.ExpArgs = a.getInt64()
	fn.Header.ParentFnIx = a.getInt64()
	fn.Header.LineStart = a.getInt64()
	fn.Header.LineEnd = a.getInt64()
	// Step to the K section (must be present, even if empty)
	a.findSection("[k]")
	a.f.Fns = append(a.f.Fns, fn)
	a.readKs(fn)
}

func (a *Asm) readKs(fn *bytecode.Fn) {
	// While the L section is not reached
	for l, ok := a.getLine(true); ok && l != "[l]"; l, ok = a.getLine(true) {
		var err error
		k := new(bytecode.K)
		// The K Type is the first character of the line
		k.Type = bytecode.KType(l[0])
		switch k.Type {
		case bytecode.KtInteger, bytecode.KtBoolean:
			// Finish the trim
			val := strings.TrimRight(l[1:], " \t")
			k.Val, err = strconv.ParseInt(val, 10, 64)
		case bytecode.KtFloat:
			val := strings.TrimRight(l[1:], " \t")
			k.Val, err = strconv.ParseFloat(val, 64)
		default:
			// Untrimmed string value
			k.Val = l[1:]
		}
		fn.Ks = append(fn.Ks, k)
		if err != nil && a.err == nil {
			a.err = err
		}
	}
	a.readLs(fn)
}

func (a *Asm) readLs(fn *bytecode.Fn) {
	// While the L section is not reached
	for l, ok := a.getLine(false); ok && l != "[i]"; l, ok = a.getLine(false) {
		var i int64
		i, a.err = strconv.ParseInt(l, 10, 64)
		fn.Ls = append(fn.Ls, i)
	}
	a.readIs(fn)
}

func (a *Asm) readIs(fn *bytecode.Fn) {
	var l string
	var ok bool
	// While a new F section is not reached
	for l, ok = a.getLine(false); ok && l != "[f]"; l, ok = a.getLine(false) {
		// Split in three parts
		parts := strings.SplitN(l, " ", 3)
		if a.assertIParts(parts) {
			var ix uint64
			o := bytecode.NewOpcode(parts[0])
			f := bytecode.NewFlag(parts[1])
			ix, a.err = strconv.ParseUint(parts[2], 10, 64)
			fn.Is = append(fn.Is, bytecode.NewInstr(o, f, ix))
		}
	}
	if ok {
		a.readFn()
	}
}

func (a *Asm) assertIParts(p []string) bool {
	if a.err != nil || a.ended {
		return false
	}
	if len(p) != 3 {
		a.err = ErrInvalidInstruction
		return false
	}
	return true
}

func (a *Asm) getInt64() int64 {
	if v, ok := a.getLine(false); ok {
		var i int64
		i, a.err = strconv.ParseInt(v, 10, 64)
		return i
	}
	return 0
}

func (a *Asm) getLine(kSect bool) (string, bool) {
	if a.err != nil || a.ended {
		return "", false
	}
	var l string
	for l == "" { // Skip empty lines or comment-only lines
		ok := a.s.Scan()
		if !ok {
			// In case of EOF, s.Scan() returns false, but s.Err() returns nil
			a.err = a.s.Err()
			a.ended = true
			return "", false
		}
		// Ignore comments
		l = a.s.Text()
		i := strings.Index(l, "//")
		if i >= 0 {
			l = l[:i]
		}
		//  For the K section, unless the line is empty, do not trim (trim left only)
		trimmed := strings.TrimSpace(l)
		if kSect && len(trimmed) > 0 {
			l = strings.TrimLeft(l, " \t")
		} else {
			l = trimmed
		}
	}
	return l, true
}
