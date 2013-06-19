package runtime

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func Asm(r io.Reader) []*FuncProto {
	var p *FuncProto
	var fps []*FuncProto

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		switch line {
		case "[f]":
			if p != nil {
				fps = append(fps, p)
			}
			p = loadFunc(s)
		case "[v]":
			loadVars(s, p)
		}
	}
	fps = append(fps, p)
	return fps
}

func loadInstrs(s *bufio.Scanner, p *FuncProto) {
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		parts := strings.Fields(line)
		l := len(parts)
		var (
			op  Opcode
			flg Flag
			ix  int
		)
		op = NewOpcode(parts[0])
		if l > 1 {
			flg = NewFlag(parts[1])
			ix, _ = strconv.Atoi(parts[2])
		}
		p.Code = append(p.Code, NewInstr(op, flg, uint64(ix)))
	}
}

func loadVars(s *bufio.Scanner, p *FuncProto) {
	var v Var
	i := 0
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "[k]") {
			loadKs(s, p)
			return
		}
		switch i {
		case 0:
			// Var name
			v = Var{}
			v.Name = s.Text()
		case 1:
			// Var file
			v.File = s.Text()
		case 2:
			// Var line start
			v.LineStart, _ = strconv.Atoi(s.Text())
		case 3:
			// Var line end
			v.LineEnd, _ = strconv.Atoi(s.Text())
			p.VTable = append(p.VTable, v)
			i = -1
		}
		i++
	}
	panic("missing constant section [k]")
}

func loadKs(s *bufio.Scanner, p *FuncProto) {
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "[i]") {
			loadInstrs(s, p)
			return
		}
		line := strings.TrimSpace(s.Text())
		switch line[0] {
		case 'i':
			// Integer
			i := String(line[1:]).Int()
			p.KTable = append(p.KTable, Int(i))
		case 'f':
			// Float
			f := String(line[1:]).Float()
			p.KTable = append(p.KTable, Float(f))
		case 's':
			// String
			p.KTable = append(p.KTable, String(line[1:]))
		case 'b':
			// Boolean
			p.KTable = append(p.KTable, Bool(line[1] == 1))
		case 'n':
			// Nil
			p.KTable = append(p.KTable, Nil)
		default:
			panic("invalid constant value type")
		}
	}
	panic("missing instructions section [i]")
}

func loadFunc(s *bufio.Scanner) *FuncProto {
	p := &FuncProto{}
	i := 0
	for s.Scan() {
		switch i {
		case 0:
			// Stack size
			p.StackSz, _ = strconv.Atoi(s.Text())
		case 1:
			// Func name
			p.Name = s.Text()
		case 2:
			// File name
			p.File = s.Text()
		case 3:
			// Line start
			p.LineStart, _ = strconv.Atoi(s.Text())
		case 4:
			// Line end
			p.LineEnd, _ = strconv.Atoi(s.Text())
			return p
		}
		i++
	}
	panic("missing func fields or scan error")
}
