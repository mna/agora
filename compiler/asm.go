package compiler

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goblin/runtime"
	"io"
	"strconv"
	"strings"
)

var (
	ErrMissingStackSz   = errors.New("missing stack size")
	ErrMissingArgsCnt   = errors.New("missing arguments count")
	ErrMissingVarsCnt   = errors.New("missing variables count")
	ErrMissingLineStart = errors.New("missing line start")
	ErrMissingLineEnd   = errors.New("missing line end")
	ErrMissingFnNm      = errors.New("missing function name")
)

type Asm struct {
	ended bool
}

func (ø *Asm) Compile(id string, r io.Reader) ([]byte, error) {
	var line string

	s := bufio.NewScanner(r)
	buf := bytes.NewBuffer(nil)

	if err := binary.Write(buf, binary.LittleEndian, _SIG); err != nil {
		return nil, err
	}
	for ø.ended = !s.Scan(); !ø.ended; ø.ended = !s.Scan() {
		line = strings.TrimSpace(s.Text())
		if line == "[f]" || line == "[k]" || line == "[i]" {
			break
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	for !ø.ended {
		switch line {
		case "[f]":
			// Read a func
			f := new(fn)
			if !ø.loadInt64(s, &f.stackSz) {
				return nil, ErrMissingStackSz
			}
			if !ø.loadInt64(s, &f.args) {
				return nil, ErrMissingArgsCnt
			}
			if !ø.loadInt64(s, &f.vars) {
				return nil, ErrMissingVarsCnt
			}
			if !ø.loadInt64(s, &f.lineStart) {
				return nil, ErrMissingLineStart
			}
			if !ø.loadInt64(s, &f.lineEnd) {
				return nil, ErrMissingLineEnd
			}
			if !s.Scan() {
				return nil, ErrMissingFnNm
			}
			nm := s.Text()
			f.nmSz = int64(len(nm))
			if err := binary.Write(buf, binary.LittleEndian, f); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, []byte(nm)); err != nil {
				return nil, err
			}
			if !s.Scan() {
				break
			}
			line = strings.TrimSpace(s.Text())

		case "[k]":
			for {
				if t, v, ok := ø.loadK(s); ok {
					// Write the K type
					if err := binary.Write(buf, binary.LittleEndian, t); err != nil {
						return nil, err
					}
					if t == 's' {
						// Write the string length
						if err := binary.Write(buf, binary.LittleEndian, int64(len(v.([]byte)))); err != nil {
							return nil, err
						}
					}
					if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
						return nil, err
					}
				} else {
					line = strings.TrimSpace(s.Text())
					break
				}
			}

		case "[i]":
			for {
				if i, ok := ø.loadI(s); ok {
					if err := binary.Write(buf, binary.LittleEndian, i); err != nil {
						return nil, err
					}
				} else {
					line = strings.TrimSpace(s.Text())
					break
				}
			}
		}
	}
	return buf.Bytes(), nil
}

func (ø *Asm) loadI(s *bufio.Scanner) (uint64, bool) {
	if !s.Scan() {
		ø.ended = true
		return 0, false
	}
	line := strings.TrimSpace(s.Text())
	if line[0] == '[' {
		return 0, false
	}
	flds := strings.Fields(s.Text())
	op := runtime.NewOpcode(flds[0])
	var f runtime.Flag
	var ix uint64
	if len(flds) > 1 {
		f = runtime.NewFlag(flds[1])
		i, _ := strconv.Atoi(flds[2])
		ix = uint64(i)
	}
	return uint64(runtime.NewInstr(op, f, ix)), true
}

func (ø *Asm) loadK(s *bufio.Scanner) (byte, interface{}, bool) {
	if !s.Scan() {
		ø.ended = true
		return byte(0), nil, false
	}
	line := strings.TrimSpace(s.Text())
	switch line[0] {
	case 'i', 'b':
		v, _ := strconv.Atoi(line[1:])
		return line[0], v, true
	case 'f':
		v, _ := strconv.ParseFloat(line[1:], 64)
		return line[0], v, true
	case 's':
		return line[0], []byte(line[1:]), true
	}
	return byte(0), nil, false
}

func (ø *Asm) loadInt64(s *bufio.Scanner, dest *int64) bool {
	if !s.Scan() {
		ø.ended = true
		return false
	}
	i, _ := strconv.Atoi(s.Text())
	*dest = int64(i)
	return true
}
