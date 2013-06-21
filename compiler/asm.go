package compiler

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goblin/runtime"
)

var (
	s *bufio.Scanner
	m map[string]func(*runtime.FuncProto)
)

func Asm(r io.Reader) *runtime.Ctx {
	ctx := runtime.NewCtx()
	s = bufio.NewScanner(r)

	m = map[string]func(*runtime.FuncProto){
		"[f]": func(_ *runtime.FuncProto) {
			p := &runtime.FuncProto{}
			i := 0
			for s.Scan() {
				switch i {
				case 0:
					if s.Text() == "true" {
						p.Native = true
						i++
					} else {
						// Stack size
						p.StackSz, _ = strconv.Atoi(s.Text())
					}
				case 1:
					// Expected args count
					p.ExpArgs, _ = strconv.Atoi(s.Text())
				case 2:
					if p.Native {
						p.NativeName = s.Text()
						i = 5
					} else {
						// Func name
						p.Name = s.Text()
					}
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

		"[k]": func(p *runtime.FuncProto) {
			for s.Scan() {
				line := strings.TrimSpace(s.Text())
				if f, ok := m[line]; ok {
					f(p)
					return
				}

				switch line[0] {
				case 'i':
					// Integer
					i := runtime.String(line[1:]).Int()
					p.KTable = append(p.KTable, runtime.Int(i))
				case 'f':
					// Float
					f := runtime.String(line[1:]).Float()
					p.KTable = append(p.KTable, runtime.Float(f))
				case 's':
					// String
					p.KTable = append(p.KTable, runtime.String(line[1:]))
				case 'b':
					// Boolean
					p.KTable = append(p.KTable, runtime.Bool(line[1] == 1))
				case 'n':
					// Nil
					p.KTable = append(p.KTable, runtime.Nil)
				default:
					panic("invalid constant value type")
				}
			}
			panic("missing instructions section [i]")
		},

		"[v]": func(p *runtime.FuncProto) {
			var v runtime.Var
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
					v = runtime.Var{}
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

		"[i]": func(p *runtime.FuncProto) {
			for s.Scan() {
				line := strings.TrimSpace(s.Text())
				if f, ok := m[line]; ok {
					f(p)
					return
				}
				parts := strings.Fields(line)
				l := len(parts)
				var (
					op  runtime.Opcode
					flg runtime.Flag
					ix  int64
				)
				op = runtime.NewOpcode(parts[0])
				if l > 1 {
					flg = runtime.NewFlag(parts[1])
					ix, _ = strconv.ParseInt(parts[2], 10, 64)
				}
				p.Code = append(p.Code, runtime.NewInstr(op, flg, uint64(ix)))
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
