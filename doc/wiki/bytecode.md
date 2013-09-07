This document presents the binary format of bytecode-compiled files. Each source file is compiled to a corresponding, standalone bytecode format file. Data is always stored in *little-endian* ordering.

## Header

The file starts with a header with this format:

* 4 bytes : the signature used to identify the bytecode file format, which is 0x000A602A (AGORA, more or less).
* 1 byte  : the version number of the compiler used to generate the bytecode file, i.e. 0x12 for v1.2 (high hexadecimal digit is the major, low hexadecimal digit is the minor version number).

The header is always exactly 5 bytes long.

## Functions

The rest of the file is made up of 1 or many function representations. Each function has this format:

* The function's header
* The function's constants or symbols (referred to as the K section)
* The function's instructions (referred to as the I section)

A **string** is encoded as follows:

* **int64**  : the length in bytes of the string.
* **bytes**  : *n* bytes representing the string.

### The function header

* **string** : the name of the function. For the top-level function, this is the name of the source file.
* **int64**  : the initial **stack size** required by the function. This is merely a hint to the VM so that a reasonable initial stack is allocated, but it may grow as needed (for example, the compiler may not take into account loops in the stack size).
* **int64**  : the number of **expected arguments** that the function may receive. Being a dynamic language, more or less actual arguments may be passed, but this represents the number of arguments that have corresponding parameters acting as local variables for these arguments inside the function. Unlike the stack size, this must be exactly the number of defined arguments on the function's signature. This value is always 0 for the top-level function.
* **int64**  : the number of **expected variables** in the function. Much like the stack size, this is merely a hint to the VM so that a reasonable initial variable map size is allocated, but it may grow as needed. The expected arguments of the function count towards the number of variables.
* **int64**  : the starting line number in the source code file where this function is defined, starting at 1. This is for debugging purpose only.
* **int64**  : the ending line number in the source code file where this function is defined, starting at 1. This is for debugging purpose only.

### The K section

There is a *header* of the K section, namely:

* **int64**  : the first field in this section represents the number of constants (or symbols) that make up the K section. For this *n* number of times, the following section is present.

Then comes *n* times the definition of a single constant:

* **byte**   : indicates the type of the constant, where `i` indicates an integer ( **int64** ), `b` a boolean (stored as an **int64**, and where `0` means `false`, any other value is `true`), `f` indicates a float ( **float64** ), and finally `s` indicates a **string**.
* **variable** : the following field depends on the type of the constant.

### The I section

There is a *header* of the I section, namely:

* **int64**  : the first field in this section represents the number of instructions that make up the I section. For this *n* number of times, the following section is present.

Then comes *n* times the definition of a single instruction:

* **uint64** : each instruction is encoded as an **uint64**.

An instruction is thus 8 bytes, composed of the following fields, from the most significant byte to the least significant:

* **1 byte**  : the first byte represents the opcode. See /runtime/opcodes.go for the definition of opcodes.
* **1 byte**  : the second byte is the *flag*, that gives meaning to the following bytes or give precisions to the opcode action. See /runtime/instr.go for the definition of flags.
* **6 bytes** : the remaining bytes contain an index into either the constant table, the `args` array or the function prototype table, or an explicit value (i.e. the number of instructions to jump over).

Next: [Assembly code format][asm]

[asm]: https://github.com/PuerkitoBio/agora/wiki/Assembly-code-format

