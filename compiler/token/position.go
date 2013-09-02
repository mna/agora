package token

import (
	"fmt"
)

// A Position is a specific point in a source file where a token is found.
type Position struct {
	Filename string
	Offset   int // Byte offset of the position in the file
	Line     int
	Column   int
}

func (p Position) IsValid() bool {
	return p.Line > 0
}

func (p Position) String() string {
	s := p.Filename
	if p.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d:%d", p.Line, p.Column)
	}
	if s == "" {
		s = "-"
	}
	return s
}
