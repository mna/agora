package agora

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/PuerkitoBio/agora/compiler"
	"github.com/PuerkitoBio/agora/runtime"
	"github.com/PuerkitoBio/agora/runtime/stdlib"
)

// This test runs all source files in ../testdata/src/*.agora and checks if
// the results are as expected.
//
// The header of each source code file can define a YAML front-matter block
// with the following fields:
// * output: the expected output (may contain \n for newlines)
// * result: the expected result value
// * long: if true, this test is skipped if the -short flag is set
// * args: the command-line arguments to pass to the test file
// * error: the expected error message (omit if no error is expected)

const (
	srcDir = "./testdata/src"
)

func TestSourceFiles(t *testing.T) {
	// Change working directory to where the source files are
	os.Chdir(srcDir)
	fis, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, fi := range fis {
		if filepath.Ext(fi.Name()) == ".agora" {
			testFile(t, fi)
		}
	}
}

func testFile(t *testing.T, fi os.FileInfo) {
	f, e := os.Open(fi.Name())
	if e != nil {
		panic(e)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	m := readFrontMatter(s)
	if len(m) == 0 {
		if testing.Verbose() {
			fmt.Printf("no front matter, skipping file %s...\n", fi.Name())
		}
		return
	}
	// Keep the rest of the file in a reader
	buf := bytes.NewBuffer(nil)
	for s.Scan() {
		buf.WriteString(s.Text())
		buf.WriteString("\n")
	}
	if s.Err() != nil {
		panic(s.Err())
	}
	// And actually run and test the file
	if _, ok := m["long"]; ok {
		if testing.Short() {
			if testing.Verbose() {
				fmt.Printf("skipping long test file %s...\n", fi.Name())
			}
			return
		}
	}
	if testing.Verbose() {
		fmt.Printf("testing file %s...\n", fi.Name())
	}
	runAndAssertFile(t, strings.TrimSuffix(fi.Name(), filepath.Ext(fi.Name())), bytes.NewReader(buf.Bytes()), m)
}

type testResolver struct {
	r  io.Reader
	mr runtime.ModuleResolver
}

func (t *testResolver) Resolve(id string) (io.Reader, error) {
	if t.r != nil {
		r := t.r
		t.r = nil
		return r, nil
	}
	return t.mr.Resolve(id)
}

func runAndAssertFile(t *testing.T, id string, r io.Reader, m map[string]string) {
	// Use the custom test resolver to return the reader
	buf := bytes.NewBuffer(nil)
	ctx := runtime.NewCtx(&testResolver{
		r,
		new(runtime.FileResolver),
	}, new(compiler.Compiler))
	ctx.Stdout = buf
	ctx.RegisterNativeModule(new(stdlib.FilepathMod))
	ctx.RegisterNativeModule(new(stdlib.FmtMod))
	ctx.RegisterNativeModule(new(stdlib.MathMod))
	ctx.RegisterNativeModule(new(stdlib.OsMod))
	ctx.RegisterNativeModule(new(stdlib.StringsMod))
	ctx.RegisterNativeModule(new(stdlib.TimeMod))

	mod, err := ctx.Load(id)
	var ret runtime.Val
	if err == nil {
		var args []runtime.Val
		if v, ok := m["args"]; ok {
			s := strings.Split(v, " ")
			args = make([]runtime.Val, len(s))
			for i, arg := range s {
				args[i] = runtime.String(arg)
			}
		}
		ret, err = mod.Run(args...)
	}

	assert := false
	if v, ok := m["error"]; ok {
		assert = true
		if err == nil {
			t.Errorf("[%s] - expected error '%s', got none", id, v)
		} else if err.Error() != v {
			t.Errorf("[%s] - expected error '%s', got '%s'", id, v, err)
		}
	} else if err != nil {
		t.Errorf("[%s] - expected no error, got '%s'", id, err)
	}
	if v, ok := m["result"]; ok {
		assert = true
		v = strings.Replace(v, "\\n", "\n", -1)
		v = strings.Replace(v, "\\t", "\t", -1)
		switch retv := ret.(type) {
		// comapre runtime.Object with special function
		case runtime.Object:
			if compareStringifiedObjects(retv.String(), v) {
				str := fmt.Sprintf("%s", retv)
				t.Errorf("[%s] - expected result '%s', got '%s'", id, v, str)
			}
		case runtime.Func:
			str := fmt.Sprintf("%s", retv)
			if str != v {
				t.Errorf("[%s] - expected result '%s', got '%s'", id, v, str)
			}
		default:
			if retv.String() != v {
				t.Errorf("[%s] - expected result '%s', got '%s'", id, v, retv)
			}
		}
	}
	if v, ok := m["output"]; ok {
		assert = true
		v = strings.Replace(v, "\\n", "\n", -1)
		v = strings.Replace(v, "\\t", "\t", -1)
		// compare output with special function
		if got := buf.String(); compareOutputsWithObjects(got, v) {
			t.Errorf("[%s] - expected output '%s', got '%s'", id, v, got)
		}
	}
	if !assert {
		t.Errorf("[%s] - no assert", id)
	}
}

/*
# Reason
--------------------------------------------------------------------------------
$ go test
--- FAIL: TestSourceFiles (38.34s)
        agora_test.go:161: [66-yield-complex] - expected output '1
                nil
                nil
                {f:test,e:7,d:4,c:3,b:2,a:1}    <--------------+
                ', got '1                                      |
                nil                                            |  just the same
                nil                                            |
                {c:3,b:2,a:1,f:test,e:7,d:4}    <--------------+
                '
--------------------------------------------------------------------------------
$ go test
--- FAIL: TestSourceFiles (42.53s)
        agora_test.go:148: [37-literal-obj] - expected result '{_:_-value!!,4:4-value!,s:s-value,i:i-value,4:int-4-value,5:int-5-value-via-i-var,string:int-s-value-via-s-var}', got '{4:4-value!,s:s-value,i:i-value,4:int-4-value,5:int-5-value-via-i-var,string:int-s-value-via-s-var,_:_-value!!}'
--------------------------------------------------------------------------------
"abc" and "abc"                                          //=> should be true
"abc, {a:1,b:2,c:3}" and "abc, {b:2,c:3,a:1}"            //=> should be true too
--------------------------------------------------------------------------------
Because runtime.Object is map[Val]Val, but map doesn't support ordering.
Playground implementation: http://play.golang.org/p/ENgKLm798i
*/
// res = result of script (runtime.Object.String())
// fmr  = front matter result
// return: false if equal, otherwise true
func compareStringifiedObjects(res string /*(runtime.Object).String()*/, fmr string) bool {
	type item struct {
		key, value string
	}
	parse_string_to_object := func(out string) ([]item, bool /*success*/) {
		// out: "{i:i-string,_:_-string,4:4-string}"
		out = strings.TrimPrefix(out, "{")
		// out: "i:i-string,_:_-string,4:4-string}"
		out = strings.TrimSuffix(out, "}")
		// out: "i:i-string,_:_-string,4:4-string"
		outs := strings.FieldsFunc(out, func(r rune) bool {
			return r == ','
		})
		// outs: [ "i:i-string", "_:_-string", "4:4-string" ]
		out_obj := make([]item, 0, len(outs))
		// fill out_obj
		for _, o := range outs {
			// o: "i:i-string"
			okv := strings.FieldsFunc(o, func(r rune) bool {
				return r == ':'
			})
			// okv: [ "i", "i-string" ]
			// check okv length
			if len(okv) != 2 {
				return nil, false
			}
			out_obj = append(out_obj, item{okv[0], okv[1]})
			// out_obj: [item{"i", "i-string"}, item{"4","4-string"}]
		}
		return out_obj, true
	}
	// front matter
	fm, ok := parse_string_to_object(fmr)
	if !ok {
		if testing.Verbose() {
			fmt.Println("unrecognized front matter object")
		}
		return true // not equal
	}
	// script result
	sr, ok := parse_string_to_object(res) // obj.String())
	if !ok {
		if testing.Verbose() {
			fmt.Println("unrecognized script object")
		}
		return true // not equal
	}
	// len(fm) == len(sr) ?
	if len(fm) != len(sr) {
		if testing.Verbose() {
			fmt.Printf(
				"improper length of the object, expected: %d, got: %d\n",
				len(fm),
				len(sr),
			)
		}
		return true // not equal
	}
	// compare objects
	for _, fmi := range fm {
		for i := 0; i < len(sr); i++ {
			if sr[i].key == fmi.key && sr[i].value == fmi.value {
				// delete item from sr slice
				sr = append(sr[:i], sr[i+1:]...)
				break
			}
		}
	}
	if len(sr) != 0 {
		return true // not equal
	}
	// all right
	return false // equal
}

// out = ouput of script
// fm  = front matter output
// return: false if equal, otherwise true
func compareOutputsWithObjects(out, fm string) bool {
	type slices struct {
		even bool
		body []string
	}
	// object regexp: {a:value,b:value-of-b!!!}
	rxp := regexp.MustCompile(`\{([^:}]+\:[^,}]+,)*([^:}]+\:[^,}]+)\}`)
	// func example: http://play.golang.org/p/fqzoSzaZqb
	cut := func(s string) (*slices, bool /*cuted*/) {
		ix := rxp.FindAllStringIndex(s, -1)
		if ix == nil {
			return nil, false
		}
		var ss []string
		var cur int
		for _, ixx := range ix {
			if cur != ixx[0] {
				ss = append(ss, s[cur:ixx[0]])
			}
			ss = append(ss, s[ixx[0]:ixx[1]])
			cur = ixx[1]
		}
		if cur != len(s) {
			ss = append(ss, s[cur:len(s)])
		}
		even := ix[0][0] != 0
		return &slices{
			even: even,
			body: ss,
		}, true
	}
	// cut out and fm
	outs, cuted := cut(out)
	if !cuted {
		// outs not cuted
		return out != fm // simple string cmp
	}
	fms, cuted := cut(fm)
	// check all
	if !cuted || outs.even != fms.even || len(outs.body) != len(fms.body) {
		return true // not equal
	}
	// loop
	for i := 0; i < len(outs.body); i++ {
		// if (ous.even and even) or (not outs.even and not even)
		if (outs.even && i%2 != 0) || (!outs.even && i%2 == 0) {
			if compareStringifiedObjects(outs.body[i], fms.body[i]) {
				return true // not equal
			}
		} else {
			if outs.body[i] != fms.body[i] {
				return true // not equal
			}
		}
	}
	// all right
	return false // equal
}

func readFrontMatter(s *bufio.Scanner) map[string]string {
	m := make(map[string]string)
	infm := false
	for s.Scan() {
		l := strings.Trim(s.Text(), " ")
		if l == "/*---" || l == "---*/" { // The front matter is delimited by 3 dashes and in a block comment
			if infm {
				// This signals the end of the front matter
				return m
			} else {
				// This is the start of the front matter
				infm = true
			}
		} else if infm {
			sections := strings.SplitN(l, ":", 2)
			if len(sections) != 2 {
				// Invalid front matter line
				return nil
			}
			m[sections[0]] = strings.Trim(sections[1], " ")
		} else if l != "" {
			// No front matter, quit
			return nil
		}
	}
	if err := s.Err(); err != nil {
		// The scanner stopped because of an error
		return nil
	}
	return nil
}
