// Dump defines the types to serialize a compiled source code module
// to its bytecode representation, executable by the VM.
package compiler

var (
	_SIG = [...]byte{'6', '0', 'B', '1', '1', '4'}
)

type fn struct {
	stackSz   int64
	args      int64
	vars      int64
	lineStart int64
	lineEnd   int64
	nmSz      int64
}
