package main

import (
	"bytes"
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goblin/compiler"
	"github.com/PuerkitoBio/goblin/runtime"
	"github.com/PuerkitoBio/goblin/runtime/stdlib"
)

const (
	testFilePath   = "./runtime/testdata/"
	expResFilePath = "./runtime/testdata/exp/"
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

	// Test all files
	for _, fi := range files {
		if filepath.Ext(fi.Name()) == ".goblin" {
			runTestFile(t, fi.Name())
		}
	}
}

func createCtx() *runtime.Ctx {
	// Create context and Stdout buffer
	ctx := runtime.NewCtx(resolv, comp)
	ctx.RegisterModule(fmtMod)
	buf := bytes.NewBuffer(nil)
	ctx.Stdout = buf
	return ctx
}

func runTestFile(t *testing.T, fnm string) {
	ctx := createCtx()

	// Execute the test file
	fmt.Println("testing ", fnm, "...")
	v, err := ctx.Load(filepath.Join(testFilePath, fnm))
	if err != nil {
		t.Errorf("failed with error %s for %s", err, fnm)
		return
	}
	fmt.Fprintf(ctx.Stdout, "PASS - %v", v)
	got := crc64.Checksum(ctx.Stdout.(*bytes.Buffer).Bytes(), ecmaTbl)

	// Load the expected result
	res, err := ioutil.ReadFile(filepath.Join(expResFilePath, strings.Replace(fnm, filepath.Ext(fnm), ".exp", 1)))
	if err != nil {
		panic(err)
	}
	exp := crc64.Checksum(res, ecmaTbl)

	// Assert
	if exp != got {
		t.Errorf("unexpected result for %s", fnm)
		t.Log(ctx.Stdout)
	}
}
