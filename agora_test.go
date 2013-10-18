package agora

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	gort "runtime"
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
	if testing.Verbose() {
		fmt.Printf("[%d goroutine(s) at startup]\n", gort.NumGoroutine())
	}
	for _, fi := range fis {
		if filepath.Ext(fi.Name()) == ".agora" {
			testFile(t, fi)
		}
	}
	if testing.Verbose() {
		fmt.Printf("[%d goroutine(s) after tests]\n", gort.NumGoroutine())
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
	var ret []runtime.Val
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
		switch retv := runtime.Get1(ret).(type) {
		case runtime.Object, runtime.Func:
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
		if got := buf.String(); got != v {
			t.Errorf("[%s] - expected output '%s', got '%s'", id, v, got)
		}
	}
	if !assert {
		t.Errorf("[%s] - no assert", id)
	}
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
