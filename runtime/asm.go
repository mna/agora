package runtime

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

var (
	s *bufio.Scanner
	m map[string]func(*FuncProto)
)

func Asm(r io.Reader) *Ctx {
	ctx := &Ctx{}
	s = bufio.NewScanner(r)

	m = map[string]func(*FuncProto){
		"[f]": func(_ *FuncProto) {
			p := &FuncProto{}
			i := 0
			for s.Scan() {
				switch i {
				case 0:
					// Stack size
					p.StackSz, _ = strconv.Atoi(s.Text())
				case 1:
					// Expected args count
					p.ExpArgs, _ = strconv.Atoi(s.Text())
				case 2:
					// Func name
					p.Name = s.Text()
				case 3:
					// File name
					p.File = s.Text()
				case 4:
					// Line start
					p.LineStart, _ = strconv.Atoi(s.Text())
				case 5:
					// Line end
					p.LineEnd, _ = strconv.Atoi(s.Text())
				default:
					ctx.Protos = append(ctx.Protos, p)
					// Find where to go from here
					f := m[s.Text()]
					f(p)
				}
				i++
			}
		},

		"[k]": func(p *FuncProto) {
			for s.Scan() {
				line := strings.TrimSpace(s.Text())
				if f, ok := m[line]; ok {
					f(p)
					return
				}

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
		},

		"[v]": func(p *FuncProto) {
			var v Var
			i := 0
			for s.Scan() {
				line := strings.TrimSpace(s.Text())
				if f, ok := m[line]; ok {
					f(p)
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
		},

		"[i]": func(p *FuncProto) {
			for s.Scan() {
				line := strings.TrimSpace(s.Text())
				if f, ok := m[line]; ok {
					f(p)
					return
				}
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
		},
	}

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if f, ok := m[line]; ok {
			f(nil)
		}
	}
	return ctx
}
