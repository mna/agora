package main

import (
	"bytes"
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/PuerkitoBio/agora/compiler"
	"github.com/PuerkitoBio/agora/runtime"
	"github.com/PuerkitoBio/agora/runtime/stdlib"
)

const (
	testFilePath      = "./runtime/testdata/"
	expResFilePath    = "./exp/"
	expectErrPrefix   = "x"
	skipOnShortPrefix = "*"

	isolateFileTest = ""
)

var (
	comp    = new(compiler.Asm)
	resolv  = new(runtime.FileResolver)
	fmtMod  = new(stdlib.FmtMod)
	ecmaTbl = crc64.MakeTable(crc64.ECMA)
)

func TestFiles(t *testing.T) {
	// Get list of files
	files, err := ioutil.ReadDir(testFilePath)
	if err != nil {
		panic(err)
	}

	// Change current directory to testFilePath, otherwise relative imports
	// will not work.
	err = os.Chdir(testFilePath)
	if err != nil {
		panic(err)
	}

	// Test all files
	for _, fi := range files {
		if filepath.Ext(fi.Name()) == ".agora" {
			if isolateFileTest != "" && isolateFileTest != fi.Name() {
				continue
			}
			runTestFile(t, fi.Name())
		}
	}
}

func createCtx() *runtime.Ctx {
	// Create context and Stdout buffer
	ctx := runtime.NewCtx(resolv, comp)
	ctx.RegisterNativeModule(fmtMod)
	buf := bytes.NewBuffer(nil)
	ctx.Stdout = buf
	return ctx
}

func runTestFile(t *testing.T, fnm string) {
	var expErrMsg string

	ori := fnm
	// Is an error expected?
	expErr := strings.HasPrefix(fnm, expectErrPrefix)
	if expErr {
		fnm = fnm[1:]
	}
	// Should it be skipped in -short mode?
	skipOnShort := strings.HasPrefix(fnm, skipOnShortPrefix)
	if skipOnShort {
		fnm = fnm[1:]
	}
	if skipOnShort && testing.Short() {
		if testing.Verbose() {
			fmt.Println("skipping ", fnm, ".")
		}
		return
	} else {
		if testing.Verbose() {
			fmt.Println("testing ", fnm, "...")
		}
	}
	// Create the execution context
	ctx := createCtx()
	// Load the expected result
	res, err := ioutil.ReadFile(filepath.Join(expResFilePath, strings.Replace(fnm, filepath.Ext(fnm), ".exp", 1)))
	if err != nil {
		panic(err)
	}
	if expErr {
		expErrMsg = string(res)
		expErrMsg = strings.Trim(expErrMsg, "\n\t\r ")
	}
	// TODO: Because LOAD currently panics when ctx.Load returns an error, cyclic
	// dependencies cause panics, so catch in a defer.
	defer func() {
		if e := recover(); e != nil && expErr {
			// An error was expected, check if this is the correct error message
			if fmt.Sprintf("%s", e) != expErrMsg {
				t.Errorf("expected error %s, got %s", expErrMsg, e)
			}
		}
	}()

	// Execute the test file
	v, err := ctx.Load(ori)
	if err != nil {
		// TODO : Still leave the expErr check there too, for errors on Load caught
		// before the runtime.
		if expErr {
			// An error was expected, check if this is the correct error message
			if err.Error() != expErrMsg {
				t.Errorf("expected error %s, got %s", expErrMsg, err)
			}
		} else {
			t.Errorf("failed with error %s for %s", err, fnm)
		}
		return
	} else if expErr {
		t.Errorf("expected error %s, got no error", expErrMsg)
	}

	// Add the PASS string to the output, since this will be printed by the execution
	fmt.Fprintf(ctx.Stdout, "PASS - %v\n", v)
	// Then compare both outputs
	got := crc64.Checksum(ctx.Stdout.(*bytes.Buffer).Bytes(), ecmaTbl)
	exp := crc64.Checksum(res, ecmaTbl)

	// Assert
	if exp != got {
		t.Errorf("unexpected result for %s", fnm)
		t.Log(ctx.Stdout)
		if testing.Verbose() {
			t.Log(string(res))
		}
	}
}
