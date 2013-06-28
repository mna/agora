package compiler

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goblin/runtime"
)

var (
	ErrMissingStackSz   = errors.New("missing stack size")
	ErrMissingArgsCnt   = errors.New("missing arguments count")
	ErrMissingVarsCnt   = errors.New("missing variables count")
	ErrMissingLineStart = errors.New("missing line start")
	ErrMissingLineEnd   = errors.New("missing line end")
	ErrMissingFnNm      = errors.New("missing function name")
)

type fn struct {
	stackSz   int64
	args      int64
	vars      int64
	lineStart int64
	lineEnd   int64
}

type Asm struct {
	ended bool
}

func (ø *Asm) Compile(id string, r io.Reader) ([]byte, error) {
	var line string

	s := bufio.NewScanner(r)
	buf := bytes.NewBuffer(nil)

	// Write the header signature (bytes 0x60B114)
	if err := binary.Write(buf, binary.LittleEndian, runtime.SIG); err != nil {
		return nil, err
	}
	// Write the module's id
	if err := ø.writeString(buf, id); err != nil {
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
			if err := binary.Write(buf, binary.LittleEndian, f); err != nil {
				return nil, err
			}
			if !s.Scan() {
				return nil, ErrMissingFnNm
			}
			nm := s.Text()
			if err := ø.writeString(buf, nm); err != nil {
				return nil, err
			}
			if !s.Scan() {
				break
			}
			line = strings.TrimSpace(s.Text())

		case "[k]":
			cntK := int64(0)
			bufK := bytes.NewBuffer(nil)
			for {
				if t, v, ok := ø.loadK(s); ok {
					cntK++
					// Write the K type
					if err := binary.Write(bufK, binary.LittleEndian, t); err != nil {
						return nil, err
					}
					if t == 's' {
						// Write the string length
						if err := binary.Write(bufK, binary.LittleEndian, int64(len(v.([]byte)))); err != nil {
							return nil, err
						}
					}
					if err := binary.Write(bufK, binary.LittleEndian, v); err != nil {
						return nil, err
					}
				} else {
					// Write the number of Ks
					if err := binary.Write(buf, binary.LittleEndian, cntK); err != nil {
						return nil, err
					}
					// Append the bufK
					if _, err := buf.Write(bufK.Bytes()); err != nil {
						return nil, err
					}
					line = strings.TrimSpace(s.Text())
					break
				}
			}

		case "[i]":
			cntI := int64(0)
			bufI := bytes.NewBuffer(nil)
			for {
				if i, ok := ø.loadI(s); ok {
					cntI++
					if err := binary.Write(bufI, binary.LittleEndian, i); err != nil {
						return nil, err
					}
				} else {
					//Write the number of instructions
					if err := binary.Write(buf, binary.LittleEndian, cntI); err != nil {
						return nil, err
					}
					// Append the bufI
					if _, err := buf.Write(bufI.Bytes()); err != nil {
						return nil, err
					}
					line = strings.TrimSpace(s.Text())
					break
				}
			}
		}
	}
	return buf.Bytes(), nil
}

func (ø *Asm) writeString(w io.Writer, s string) error {
	if err := binary.Write(w, binary.LittleEndian, int64(len(s))); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, []byte(s))
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
