Agora comes with a command-line tool to quickly build and run programs without the need for a custom Go host (remember, agora is an embeddable language). Provided your `$GOPATH/bin` environment variable is included in your `$PATH`, you can type `agora -h` to print the help screen of the tool.

The tool supports the following sub-commands:

* asm : compile assembly source to bytecode
* ast : pretty-print the abstract syntax tree of agora source
* build : compile agora source to bytecode
* dasm : disassemble bytecode to assembly source
* run : compile and execute agora source
* version : print the current agora version

## Shebang #!

It is also possible to run scripts with the [shebang notation][shebang] on Unix-y systems. For example, provided that the `agora` command-line is in your `$PATH`:

```
#!/usr/bin/env agora run -R
fmt := import("fmt")
fmt.Println("Hello from agora!")
```

You just need to set your script file to be executable (`chmod +x FILE`) and then it can be run like any shell script. The agora scanner will recognize the `#!` if it is on the first line, at the very beginning of the file, and will treat it like a comment.

## asm

`agora asm [OPTIONS] FILE`

The `asm` sub-command compiles an [assembly source file][assembly] into the bytecode binary format.

Options:

```
-o (--output) : save to this output file
-x (--hexadecimal) : produce hexadecimal output instead of raw binary
```

## ast

`agora ast [OPTIONS] FILE`

The `ast` sub-command prints the abstract syntax tree of an agora source code file.

Options:

```
-o (--output) : save to this output file
-e (--all-errors) : print all errors, not just a summary
```

## build

`agora build [OPTIONS] FILE`

The `build` sub-command compiles an agora source file to bytecode.

Options:

```
-o (--output) : save to this output file
-a (--assembly) : build to assembly source instead of bytecode
```

## dasm

`agora dasm [OPTIONS] FILE`

The `dasm` sub-command disassembles a bytecode file to assembly source.

Options:

```
-o (--output) : save to this output file
```

## run

`agora run [OPTIONS] FILE [args...]`

The `run` sub-command compiles and executes an agora source file, and prints the result. Additional values after the file are passed as arguments to the agora module.

Options:

```
-a (--from-asm) : compile and execute from an assembly source file
-d (--debug) : run in debug mode
-o (--output) : save to this output file
-R (--no-result) : do not print the result value
-S (--no-stdlib) : do not register the stdlib in the execution context
```

## version

`agora version`

The `version` sub-command prints the current agora version.

Next: [Roadmap][next]

[next]: https://github.com/PuerkitoBio/agora/wiki/Roadmap
[shebang]: http://en.wikipedia.org/wiki/Shebang_(Unix)

