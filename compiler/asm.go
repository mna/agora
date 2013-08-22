package compiler

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goblin/bytecode"
)

var (
	ErrInvalidInstruction = errors.New("invalid instruction")
	ErrNoInput            = errors.New("no input provided")
)

type Asm struct {
	s     *bufio.Scanner
	f     *bytecode.File
	ended bool
	err   error
}

func (a *Asm) Compile(id string, r io.Reader) ([]byte, error) {
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
	// Do not compile unnecessarily if there is an error
	if a.err != nil {
		return nil, a.err
	}
	// Compile to bytecode
	buf := bytes.NewBuffer(nil)
	err := bytecode.NewEncoder(buf).Encode(a.f)
	return buf.Bytes(), err
}

func (a *Asm) findSection(s string) {
	for line, ok := a.getLine(); ok; line, ok = a.getLine() {
		if line == s {
			break
		}
	}
}

func (a *Asm) readFn() {
	fn := new(bytecode.Fn)
	fn.Header.Name, _ = a.getLine()
	fn.Header.StackSz = a.getInt64()
	fn.Header.ExpArgs = a.getInt64()
	fn.Header.ExpVars = a.getInt64()
	fn.Header.LineStart = a.getInt64()
	fn.Header.LineEnd = a.getInt64()
	// Step to the K section (must be present, even if empty)
	a.findSection("[k]")
	a.f.Fns = append(a.f.Fns, fn)
	a.readKs(fn)
}

func (a *Asm) readKs(fn *bytecode.Fn) {
	// While the I section is not reached
	for l, ok := a.getLine(); ok && l != "[i]"; l, ok = a.getLine() {
		var err error
		k := new(bytecode.K)
		// The K Type is the first character of the line
		k.Type = bytecode.KType(l[0])
		switch k.Type {
		case bytecode.KtInteger, bytecode.KtBoolean:
			k.Val, err = strconv.ParseInt(l[1:], 10, 64)
		case bytecode.KtFloat:
			k.Val, err = strconv.ParseFloat(l[1:], 64)
		default:
			k.Val = l[1:]
		}
		fn.Ks = append(fn.Ks, k)
		if err != nil && a.err == nil {
			a.err = err
		}
	}
	a.readIs(fn)
}

func (a *Asm) readIs(fn *bytecode.Fn) {
	var l string
	var ok bool
	// While a new F section is not reached
	for l, ok = a.getLine(); ok && l != "[f]"; l, ok = a.getLine() {
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
	if v, ok := a.getLine(); ok {
		var i int64
		i, a.err = strconv.ParseInt(v, 10, 64)
		return i
	}
	return 0
}

func (a *Asm) getLine() (string, bool) {
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
		l = strings.TrimSpace(l)
	}
	return l, true
}
