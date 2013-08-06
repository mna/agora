package ast

import (
	"bytes"
	"fmt"
	"strings"
)

type Module struct {
	filename string
	imports  []Import
}

func NewModule(fn string) *Module {
	return &Module{
		filename: fn,
	}
}

func (m *Module) AddImport(path, ident string) {
	m.imports = append(m.imports, NewImport(path, ident))
}

func (m *Module) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("module %s\n", m.filename))
	for _, i := range m.imports {
		buf.WriteString(fmt.Sprintf("  %s\n", i))
	}
	return buf.String()
}

type Import struct {
	path  string
	ident string
}

func NewImport(path, ident string) Import {
	if len(ident) == 0 {
		i := strings.LastIndex(path, "/")
		ident = path[i+1:]
	}
	return Import{
		path,
		ident,
	}
}

func (i Import) String() string {
	return fmt.Sprintf("%s (%s)", i.path, i.ident)
}
