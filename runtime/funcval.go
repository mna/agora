package runtime

import (
	"fmt"
)

// funcVal implements most of the Val interface's methods, except
// Native() which must be done on the actual type.
type funcVal struct {
	ctx  *Ctx
	name string
}

func (f *funcVal) dump() string {
	return fmt.Sprintf("%s (Func)", f.name)
}

// Int is an invalid conversion.
func (f *funcVal) Int() int64 {
	panic(NewTypeError("func", "", "int"))
}

// Float is an invalid conversion.
func (f *funcVal) Float() float64 {
	panic(NewTypeError("func", "", "float"))
}

// String prints the function representation
func (f *funcVal) String() string {
	return fmt.Sprintf("<func %s (%p)>", f.name, f)
}

// Bool returns true.
func (f *funcVal) Bool() bool {
	return true
}
