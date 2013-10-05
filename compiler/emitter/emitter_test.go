package emitter

import (
	"fmt"
	"testing"

	"github.com/PuerkitoBio/agora/bytecode"
	"github.com/PuerkitoBio/agora/compiler/parser"
	"github.com/davecgh/go-spew/spew"
)

var (
	// The cases here match the files in /compiler/emitter/testdata/*
	// The symbols (src field of the case) are validated by running
	// `agora ast FILE`.
	emitcases = []struct {
		src []*parser.Symbol
		exp *bytecode.File
		err bool
	}{
		0: {
			// Assignment
			src: []*parser.Symbol{
				&parser.Symbol{Id: ":=", Ar: parser.ArBinary, First: &parser.Symbol{Id: "(name)", Val: "a"}, Second: &parser.Symbol{Id: "(literal)", Val: "5"}},
			},
			exp: &bytecode.File{
				Fns: []*bytecode.Fn{
					&bytecode.Fn{
						// := emits Second before First, so literal is K0
						Ks: []*bytecode.K{
							&bytecode.K{
								Type: bytecode.KtInteger,
								Val:  int64(5),
							},
							&bytecode.K{
								Type: bytecode.KtString,
								Val:  "a",
							},
						},
						Is: []bytecode.Instr{
							bytecode.NewInstr(bytecode.OP_PUSH, bytecode.FLG_K, 0),
							bytecode.NewInstr(bytecode.OP_POP, bytecode.FLG_V, 1),
						},
					},
				},
			},
		},
		1: {
			// return nil
			src: []*parser.Symbol{
				&parser.Symbol{Id: "return", Ar: parser.ArStatement, First: &parser.Symbol{Id: "nil", Ar: parser.ArName, Val: nil}},
			},
			exp: &bytecode.File{
				Fns: []*bytecode.Fn{
					&bytecode.Fn{
						Is: []bytecode.Instr{
							bytecode.NewInstr(bytecode.OP_PUSH, bytecode.FLG_N, 0),
							bytecode.NewInstr(bytecode.OP_RET, bytecode.FLG__, 0),
						},
					},
				},
			},
		},
		2: {
			// NOT (!) operator
			src: []*parser.Symbol{
				&parser.Symbol{Id: ":=", Ar: parser.ArBinary, First: &parser.Symbol{Id: "(name)", Val: "a"},
					Second: &parser.Symbol{Id: "!", Ar: parser.ArUnary, First: &parser.Symbol{Id: "true", Val: true, Ar: parser.ArLiteral}}},
			},
			exp: &bytecode.File{
				Fns: []*bytecode.Fn{
					&bytecode.Fn{
						Ks: []*bytecode.K{
							&bytecode.K{
								Type: bytecode.KtBoolean,
								Val:  int64(1),
							},
							&bytecode.K{
								Type: bytecode.KtString,
								Val:  "a",
							},
						},
						Is: []bytecode.Instr{
							bytecode.NewInstr(bytecode.OP_PUSH, bytecode.FLG_K, 0),
							bytecode.NewInstr(bytecode.OP_NOT, bytecode.FLG__, 0),
							bytecode.NewInstr(bytecode.OP_POP, bytecode.FLG_V, 1),
						},
					},
				},
			},
		},
		3: {
			// UNM (-) operator
			src: []*parser.Symbol{
				&parser.Symbol{Id: ":=", Ar: parser.ArBinary, First: &parser.Symbol{Id: "(name)", Val: "a"},
					Second: &parser.Symbol{Id: "-", Ar: parser.ArUnary, First: &parser.Symbol{Id: "(literal)", Val: "1", Ar: parser.ArLiteral}}},
			},
			exp: &bytecode.File{
				Fns: []*bytecode.Fn{
					&bytecode.Fn{
						Ks: []*bytecode.K{
							&bytecode.K{
								Type: bytecode.KtInteger,
								Val:  int64(1),
							},
							&bytecode.K{
								Type: bytecode.KtString,
								Val:  "a",
							},
						},
						Is: []bytecode.Instr{
							bytecode.NewInstr(bytecode.OP_PUSH, bytecode.FLG_K, 0),
							bytecode.NewInstr(bytecode.OP_UNM, bytecode.FLG__, 0),
							bytecode.NewInstr(bytecode.OP_POP, bytecode.FLG_V, 1),
						},
					},
				},
			},
		},
		4: {
			// ADD (+) operator
			src: []*parser.Symbol{
				&parser.Symbol{Id: ":=", Ar: parser.ArBinary, First: &parser.Symbol{Id: "(name)", Val: "a"},
					Second: &parser.Symbol{Id: "+", Ar: parser.ArBinary, First: &parser.Symbol{Id: "(literal)", Val: "5", Ar: parser.ArLiteral}, Second: &parser.Symbol{Id: "(literal)", Val: "2", Ar: parser.ArLiteral}}},
			},
			exp: &bytecode.File{
				Fns: []*bytecode.Fn{
					&bytecode.Fn{
						Ks: []*bytecode.K{
							&bytecode.K{
								Type: bytecode.KtInteger,
								Val:  int64(5),
							},
							&bytecode.K{
								Type: bytecode.KtInteger,
								Val:  int64(2),
							},
							&bytecode.K{
								Type: bytecode.KtString,
								Val:  "a",
							},
						},
						Is: []bytecode.Instr{
							bytecode.NewInstr(bytecode.OP_PUSH, bytecode.FLG_K, 0),
							bytecode.NewInstr(bytecode.OP_PUSH, bytecode.FLG_K, 1),
							bytecode.NewInstr(bytecode.OP_ADD, bytecode.FLG__, 0),
							bytecode.NewInstr(bytecode.OP_POP, bytecode.FLG_V, 2),
						},
					},
				},
			},
		},
	}

	isolateEmitCase = -1
)

func TestEmit(t *testing.T) {
	// Arrange
	e := new(Emitter)
	for i, c := range emitcases {
		if isolateEmitCase >= 0 && isolateEmitCase != i {
			continue
		}
		if testing.Verbose() {
			fmt.Printf("testing emitter case %d...\n", i)
		}

		// Act
		f, err := e.Emit("", c.src, nil)

		// Assert
		if (err != nil) != c.err {
			if err == nil {
				t.Errorf("[%d] - expected an error, got none", i)
			} else {
				t.Errorf("[%d] - expected no error, got `%s`", i, err)
			}
		}
		if c.exp != nil {
			if !equal(i, f, c.exp) {
				t.Errorf("[%d] - expected\n", i)
				t.Error(spew.Sdump(c.exp))
				t.Error("got\n")
				t.Error(spew.Sdump(f))
			}
		}
		if !c.err && c.exp == nil {
			t.Errorf("[%d] - no assertion", i)
		}
	}
}

// Duplicate from decode_test.go in bytecode package, but different
// equal checks, here I don't care about Header and such
func equal(c int, f1, f2 *bytecode.File) bool {
	if f1 == nil && f2 == nil {
		return true
	}
	if f1 == nil || f2 == nil {
		if testing.Verbose() {
			fmt.Printf("[%d] - error: one *File is nil\n", c)
		}
		return false
	}
	// Ignore name, version
	if len(f1.Fns) != len(f2.Fns) {
		if testing.Verbose() {
			fmt.Printf("[%d] - error: f1 has %d funcs, f2 has %d\n", c, len(f1.Fns), len(f2.Fns))
		}
		return false
	}
	for i := 0; i < len(f1.Fns); i++ {
		fn1, fn2 := f1.Fns[i], f2.Fns[i]
		// Ignore function header, care only about instructions and Ks
		if len(fn1.Ks) != len(fn2.Ks) {
			if testing.Verbose() {
				fmt.Printf("[%d] - error: f1.func[%d] has %d Ks, f2.func[%d] has %d\n", c, i, len(fn1.Ks), i, len(fn2.Ks))
			}
			return false
		}
		for j := 0; j < len(fn1.Ks); j++ {
			k1, k2 := fn1.Ks[j], fn2.Ks[j]
			if k1.Type != k2.Type {
				if testing.Verbose() {
					fmt.Printf("[%d] - error: f1.func[%d].Ks[%d] has type %c, f2.func[%d].Ks[%d] has type %c\n", c, i, j, k1.Type, i, j, k2.Type)
				}
				return false
			}
			if k1.Val != k2.Val {
				if testing.Verbose() {
					fmt.Printf("[%d] - error: f1.func[%d].Ks[%d] has value %v (%T), f2.func[%d].Ks[%d] has value %v (%T)\n", c, i, j, k1.Val, k1.Val, i, j, k2.Val, k2.Val)
				}
				return false
			}
		}
		if len(fn1.Is) != len(fn2.Is) {
			if testing.Verbose() {
				fmt.Printf("[%d] - error: f1.func[%d] has %d Is, f2.func[%d] has %d\n", c, i, len(fn1.Is), i, len(fn2.Is))
			}
			return false
		}
		for j := 0; j < len(fn1.Is); j++ {
			if fn1.Is[j] != fn2.Is[j] {
				if testing.Verbose() {
					fmt.Printf("[%d] - error: f1.func[%d].Is[%d] has instr %s, f2.func[%d].Is[%d] has instr %s\n", c, i, j, fn1.Is[j], i, j, fn2.Is[j])
				}
				return false
			}
		}
	}
	return true
}
