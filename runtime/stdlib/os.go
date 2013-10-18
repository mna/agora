package stdlib

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/PuerkitoBio/agora/runtime"
)

// The os module, as documented in
// https://github.com/PuerkitoBio/agora/wiki/Standard-library
type OsMod struct {
	ctx *runtime.Ctx
	ob  runtime.Object
}

type file struct {
	runtime.Object
	f *os.File
	s *bufio.Scanner
}

func (o *OsMod) newFile(f *os.File) *file {
	ob := runtime.NewObject()
	of := &file{
		ob,
		f,
		nil,
	}
	ob.Set(runtime.String("Name"), runtime.String(f.Name()))
	ob.Set(runtime.String("Close"), runtime.NewNativeFunc(o.ctx, "os.File.Close", of.closeFile))
	ob.Set(runtime.String("ReadLine"), runtime.NewNativeFunc(o.ctx, "os.File.ReadLine", of.readLine))
	ob.Set(runtime.String("Seek"), runtime.NewNativeFunc(o.ctx, "os.File.Seek", of.seek))
	ob.Set(runtime.String("Write"), runtime.NewNativeFunc(o.ctx, "os.File.Write", of.write))
	ob.Set(runtime.String("WriteLine"), runtime.NewNativeFunc(o.ctx, "os.File.WriteLine", of.writeLine))
	return of
}

func (of *file) closeFile(args ...runtime.Val) []runtime.Val {
	e := of.f.Close()
	if e != nil {
		panic(e)
	}
	return nil
}

func (of *file) readLine(args ...runtime.Val) []runtime.Val {
	if of.s == nil {
		of.s = bufio.NewScanner(of.f)
	}
	if of.s.Scan() {
		return runtime.Set1(runtime.String(of.s.Text()))
	}
	if e := of.s.Err(); e != nil {
		panic(e)
	}
	return runtime.Set1(runtime.Nil)
}

func (of *file) seek(args ...runtime.Val) []runtime.Val {
	off := int64(0)
	if len(args) > 0 {
		off = args[0].Int()
	}
	rel := 0
	if len(args) > 1 {
		rel = int(args[1].Int())
	}
	n, e := of.f.Seek(off, rel)
	if e != nil {
		panic(e)
	}
	return runtime.Set1(runtime.Number(n))
}

func (of *file) write(args ...runtime.Val) []runtime.Val {
	n := 0
	for _, v := range args {
		m, e := of.f.WriteString(v.String())
		if e != nil {
			panic(e)
		}
		n += m
	}
	return runtime.Set1(runtime.Number(n))
}

func (of *file) writeLine(args ...runtime.Val) []runtime.Val {
	n := runtime.Get1(of.write(args...))
	m, e := of.f.WriteString("\n")
	if e != nil {
		panic(e)
	}
	return runtime.Set1(runtime.Number(int(n.Int()) + m))
}

func (o *OsMod) ID() string {
	return "os"
}

func (o *OsMod) Run(_ ...runtime.Val) (v []runtime.Val, err error) {
	defer runtime.PanicToError(&err)
	if o.ob == nil {
		// Prepare the object
		o.ob = runtime.NewObject()
		o.ob.Set(runtime.String("TempDir"), runtime.String(os.TempDir()))
		o.ob.Set(runtime.String("PathSeparator"), runtime.String(os.PathSeparator))
		o.ob.Set(runtime.String("PathListSeparator"), runtime.String(os.PathListSeparator))
		o.ob.Set(runtime.String("DevNull"), runtime.String(os.DevNull))
		o.ob.Set(runtime.String("Exec"), runtime.NewNativeFunc(o.ctx, "os.Exec", o.os_Exec))
		o.ob.Set(runtime.String("Exit"), runtime.NewNativeFunc(o.ctx, "os.Exit", o.os_Exit))
		o.ob.Set(runtime.String("Getenv"), runtime.NewNativeFunc(o.ctx, "os.Getenv", o.os_Getenv))
		o.ob.Set(runtime.String("Getwd"), runtime.NewNativeFunc(o.ctx, "os.Getwd", o.os_Getwd))
		o.ob.Set(runtime.String("ReadFile"), runtime.NewNativeFunc(o.ctx, "os.ReadFile", o.os_ReadFile))
		o.ob.Set(runtime.String("WriteFile"), runtime.NewNativeFunc(o.ctx, "os.WriteFile", o.os_WriteFile))
		o.ob.Set(runtime.String("Open"), runtime.NewNativeFunc(o.ctx, "os.Open", o.os_Open))
		o.ob.Set(runtime.String("TryOpen"), runtime.NewNativeFunc(o.ctx, "os.TryOpen", o.os_TryOpen))
		o.ob.Set(runtime.String("Mkdir"), runtime.NewNativeFunc(o.ctx, "os.Mkdir", o.os_Mkdir))
		o.ob.Set(runtime.String("Remove"), runtime.NewNativeFunc(o.ctx, "os.Remove", o.os_Remove))
		o.ob.Set(runtime.String("RemoveAll"), runtime.NewNativeFunc(o.ctx, "os.RemoveAll", o.os_RemoveAll))
		o.ob.Set(runtime.String("Rename"), runtime.NewNativeFunc(o.ctx, "os.Rename", o.os_Rename))
		o.ob.Set(runtime.String("ReadDir"), runtime.NewNativeFunc(o.ctx, "os.ReadDir", o.os_ReadDir))
	}
	return runtime.Set1(o.ob), nil
}

func (o *OsMod) SetCtx(ctx *runtime.Ctx) {
	o.ctx = ctx
}

func (o *OsMod) os_Exit(args ...runtime.Val) []runtime.Val {
	if len(args) == 0 {
		os.Exit(0)
	}
	os.Exit(int(args[0].Int()))
	return nil
}

func (o *OsMod) os_Getenv(args ...runtime.Val) []runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	return runtime.Set1(runtime.String(os.Getenv(args[0].String())))
}

func (o *OsMod) os_Getwd(args ...runtime.Val) []runtime.Val {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return runtime.Set1(runtime.String(pwd))
}

func (o *OsMod) os_Exec(args ...runtime.Val) []runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	c := exec.Command(args[0].String(), toString(args[1:])...)
	b, e := c.CombinedOutput()
	if e != nil {
		panic(e)
	}
	return runtime.Set1(runtime.String(b))
}

func (o *OsMod) os_Mkdir(args ...runtime.Val) []runtime.Val {
	// No-op if no arg
	if len(args) == 0 {
		return nil
	}
	perm := os.FileMode(0777)
	// Last args *may* be the permissions to use if it is a number
	if l, ok := args[len(args)-1].(runtime.Number); ok {
		perm = os.FileMode(l.Int())
		args = args[:len(args)-1]
	}
	// Use the mkdir-all version, to create all missing dirs as required
	for _, v := range args {
		if e := os.MkdirAll(v.String(), perm); e != nil {
			panic(e)
		}
	}
	return nil
}

func createFileInfo(fi os.FileInfo) []runtime.Val {
	o := runtime.NewObject()
	o.Set(runtime.String("Name"), runtime.String(fi.Name()))
	o.Set(runtime.String("Size"), runtime.Number(fi.Size()))
	o.Set(runtime.String("IsDir"), runtime.Bool(fi.IsDir()))
	return runtime.Set1(o)
}

func (o *OsMod) os_ReadDir(args ...runtime.Val) []runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	fis, e := ioutil.ReadDir(args[0].String())
	if e != nil {
		panic(e)
	}
	ob := runtime.NewObject()
	for i, fi := range fis {
		ob.Set(runtime.Number(i), runtime.Get1(createFileInfo(fi)))
	}
	return runtime.Set1(ob)
}

func (o *OsMod) os_Remove(args ...runtime.Val) []runtime.Val {
	for _, v := range args {
		if e := os.Remove(v.String()); e != nil {
			panic(e)
		}
	}
	return nil
}

func (o *OsMod) os_RemoveAll(args ...runtime.Val) []runtime.Val {
	for _, v := range args {
		if e := os.RemoveAll(v.String()); e != nil {
			panic(e)
		}
	}
	return nil
}

func (o *OsMod) os_Rename(args ...runtime.Val) []runtime.Val {
	runtime.ExpectAtLeastNArgs(2, args)
	if e := os.Rename(args[0].String(), args[1].String()); e != nil {
		panic(e)
	}
	return nil
}

func (o *OsMod) os_ReadFile(args ...runtime.Val) []runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	b, e := ioutil.ReadFile(args[0].String())
	if e != nil {
		panic(e)
	}
	return runtime.Set1(runtime.String(b))
}

func (o *OsMod) os_WriteFile(args ...runtime.Val) []runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	f, e := os.Create(args[0].String())
	if e != nil {
		panic(e)
	}
	defer f.Close()
	n := 0
	for _, v := range args[1:] {
		m, e := f.WriteString(v.String())
		if e != nil {
			panic(e)
		}
		n += m
	}
	return runtime.Set1(runtime.Number(n))
}

func (o *OsMod) os_TryOpen(args ...runtime.Val) (ret []runtime.Val) {
	defer func() {
		if e := recover(); e != nil {
			ret = runtime.Set1(runtime.Nil)
		}
	}()
	return o.os_Open(args...)
}

func (o *OsMod) os_Open(args ...runtime.Val) []runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	nm := args[0].String()
	flg := "r" // defaults to read-only
	if len(args) > 1 {
		// Second arg is the flag (same as C's fopen)
		// r  - open for reading
		// w  - open for writing (file need not exist)
		// a  - open for appending (file need not exist)
		// r+ - open for reading and writing, start at beginning
		// w+ - open for reading and writing (overwrite file)
		// a+ - open for reading and writing (append if file exists)
		flg = args[1].String()
	}
	var flgi int
	switch flg {
	case "r":
		flgi = os.O_RDONLY
	case "w":
		flgi = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	case "a":
		flgi = os.O_APPEND | os.O_CREATE
	case "r+":
		flgi = os.O_RDWR
	case "w+":
		flgi = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	case "a+":
		flgi = os.O_RDWR | os.O_APPEND | os.O_CREATE
	default:
		panic("invalid file flag mode: " + flg)
	}
	f, e := os.OpenFile(nm, flgi, 0666)
	if e != nil {
		panic(e)
	}
	return runtime.Set1(o.newFile(f))
}

func toString(args []runtime.Val) []string {
	s := make([]string, len(args))
	for i, a := range args {
		s[i] = a.String()
	}
	return s
}
