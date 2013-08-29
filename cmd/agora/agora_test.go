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
	testFilePath      = "./testdata/asm/"
	expResFilePath    = "./exp/" // relative to testFilePath, since a Chdir is made
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
		if filepath.Ext(fi.Name()) == ".asm" {
			if isolateFileTest != "" && !strings.HasPrefix(fi.Name(), isolateFileTest) {
				continue
			}
			runTestFile(t, fi.Name())
		}
	}
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
	// Create the execution command
	r := &run{
		Debug:   true,
		FromAsm: true,
	}
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
	// Setup the out buffer
	buf := bytes.NewBuffer(nil)
	stdout = buf

	// Execute the test file
	fmt.Println(ori)
	err = r.Execute([]string{ori})
	if err != nil {
		// TODO : Still leave the expErr check there too, for errors caught
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

	// Compare both outputs
	got := crc64.Checksum(buf.Bytes(), ecmaTbl)
	exp := crc64.Checksum(res, ecmaTbl)

	// Assert
	if exp != got {
		t.Errorf("unexpected result for %s", fnm)
		t.Log(buf.String())
		if testing.Verbose() {
			t.Log(string(res))
		}
	}
}
