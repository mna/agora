The "assembly" source code format is an easy-to-parse assembly-style syntax that provides a human-readable representation of the low-level bytecode that the virtual machine executes. This articles documents the format of this representation, which closely matches the bytecode format. For a more thorough explanation of the various fields, see the [bytecode format](https://github.com/PuerkitoBio/agora/wiki/Bytecode-format) article. For examples of assembly source code files, look in the /testdata/asm directory.

Note that anywhere in the code, whitespace-only lines and comment-only lines are skipped. The only comment notation allowed is the up-to-the-end-of-line `//` style.

## The function

An assembly source must have at least one function section. The first function section represents the top-level (module) function. The function section is identified by the string `[f]`.

Then comes the function header, with the following fields, one per line:

1. The function's name. The top-level function's name should be the name of the file or the identifier of the module.
2. The expected stack size.
3. The expected arguments count.
4. The expected variables count.
5. The starting line of the function in the source code.
6. The ending line of the function in the source code.

This is followed by the constant section, or the K section.

## The K section

Each function must have a K section, which may be empty, identified by the string `[k]`. This section lists the various constants or symbols required by the function, one per line. The K information follows this format:

* The first character is the constant's type. It must be one of `i` for integer, `f` for float, `b` for boolean, and `s` for string.
* The remaining characters represent the constant's value. Booleans are represented as `0` for `false` and `1` for true. Floats must be in a format understood by `strconv.ParseFloat()`. Integers must be in base-10.

Next comes the instruction section, or the I section.

## The I section

Each function must have an I section, which may be empty, identified by the string `[i]`. This section lists the instructions required to execute the function, one per line. Each instruction follows this format, separated by one space, and each part is required:

1. The operation code. See /bytecode/opcodes.go for the list of valid identifiers (the string literal representation of the opcode is used, i.e. the keys of the `OpLookup` variable).
2. The operation flag. See /bytecode/instr.go for the list of valid identifiers (the string literal representation of the flag is used, i.e. the keys of the `FlagLookup` variable).
3. The index value. This is an integer in base-10.

## Repeat

Multiple `[f]` sections can then follow, each with its own K and I section. When an instruction refers to a function (for example `PUSH F 3`), the index value is the index of the function in the assembly code, starting at 0.

The same goes for instructions that refer to a constant or symbol (for example, `PUSH K 2` or `POP V 3` - push value of constant at index 2; pop into variable identified by the constant at index 3). The index is the position of the constant or symbol in the K section of the assembly code.

Next: [Virtual machine](https://github.com/PuerkitoBio/agora/wiki/Virtual-machine)

