# Getting started

This article will get you up to speed in installing and using the agora programming language.

## Installation

The first step to use agora is to install it! Provided you already have the [Go language][go] installed and [your `$GOPATH` environment variable setup correctly][gopath], this is a simple matter of running this command:

`go get github.com/PuerkitoBio/agora/...`

And yes, the three dots at the end are part of the command, literally. The agora repository is a collection of multiple packages, and this command will instruct `go get` to install all of them.

To test the installation, run the following command (`$` represents the command prompt):

```
$ cd $GOPATH/src/github.com/PuerkitoBio/agora
$ go test -short ./...
```

Again, the dot-slash-three-dots is part of the command, literally. This will run the tests of the various agora packages - skipping the long-running ones - and should all succeed, displaying `ok` at the start of each line (some subdirectories may not have tests, in which case it will display a `?`, this is normal).

## Introduction to agora

### Agora is dynamically-typed

This means that variables, parameters and functions (the return value) have no types, only *values* have a type. The following types are supported:

* String, i.e. `"hello, I'm a string!"`
* Number, i.e. `17` or `3.1415`
* Boolean, i.e. `true` or `false`
* Function, i.e. `func add(x, y) { return x + y }`
* Object, i.e. `{name: "Martin", age: 38}`
* Nil, i.e. `nil`

### Agora is embeddable

This means that it can be part of a bigger program. Since agora is built in the Go programming language as a collection of packages (think *libraries*), it is designed to be easily embedded in a Go program to provide dynamic scripting capabilities to the otherwise static Go executable (see the [native Go API][native] for more information on how to call agora from Go). Go can execute agora code, and agora code can call Go code via *native modules* exposed to agora.

But thanks to the `agora` command-line tool (which we'll get to in a minute), it is also possible to run agora programs directly, without a host. Well, this is more or less true: the tool *is* a general-purpose Go host.

### Agora is garbage-collected

You don't have to manually manage the memory in an agora program (in fact, you can't). All values are automatically garbage-collected even though there is no garbage collector in agora *per se*! It is the native Go GC that is responsible for reclaiming unused memory, and the agora runtime makes sure to release any reference to values it doesn't need anymore, as soon as possible.

### Agora is dynamically loaded

Agora modules are loaded dynamically at runtime, as opposed to compiled and linked into the binary executable as is the case with Go. The `import(string)` built-in function is responsible to load those modules in agora code, while the `runtime.Ctx.Load(string)` method is the native Go API to load the initial module to bootstrap the agora execution.

The execution context (`runtime.Ctx`) uses a *module resolver* to find the matching agora source code for a given module identifier (a string literal). Out-of-the-box, agora provides a `runtime.FileResolver` that looks up the module in the file system, but any type that implements the `runtime.ModuleResolver` interface can be used to resolve a module ID.

### Agora is interpreted

Once compiled, the agora bytecode is interpreted by a stack-based virtual machine. The bytecode is essentially a list of *opcodes* (operations) and some metadata, such as `OP_PUSH K 1`. This instruction is composed of the operation `OP_PUSH`, the flag `K` and the index `1` which instructs the VM to push onto the stack the value of the constant (`K`) at index 1 in the constant table.

See the [design and architecture][design] article for more information on the VM and the internals of agora.

### Agora has 4 representation formats

* The agora source code is most likely the code that humans write, and is the one with a syntax similar to Go.
* The agora assembly code is a *pseudo-assembly* language that is basically a human-readable representation of the binary bytecode. It is possible to compile this format using the `agora asm` command, and to disassemble the bytecode format to assembly code using `agora dasm`. It is not meant to be written directly by a human, although it is definitely possible.
* The bytecode format is the binary, compiled code that can be persisted, for example, in a file.
* Finally, the bytecode format also has a matching in-memory representation, used at runtime to execute the instructions in the virtual machine. This in-memory representation is provided by the data structure `bytecode.File`.

For more information on the various formats, see the [bytecode format][bytecode] and [assembly code format][assembly] articles in the wiki.

## The command-line tool

Agora comes with a command-line tool to quickly build and run programs without the need for a custom Go host (remember, agora is an *embeddable* language). Provided your `$GOPATH/bin` environment variable is included in your `$PATH`, you can type `agora -h` to print the help screen of the tool.

The tool supports many sub-commands:

* ast : pretty-print the abstract syntax tree of an agora source file, as generated by the parser.
* asm : compile an assembly source code to bytecode format.
* dasm : disassemble a bytecode file to an assembly source code.
* build : compile an agora source code to bytecode format.
* run : compile and execute an agora source code.

Each of these sub-commands support various flags, use `agora <subcommand> -h` for a contextual help screen describing the applicable options.

For this *getting started* article, we will use only `agora run` to run a simple program.

## Your first agora program

Let's write a simple program that converts the case of its command-line arguments based on the case of the first letter. If the word starts with an uppercase, the whole word is converted to lowercase, and vice-versa. Admittedly, this is not terribly useful, but the goal is to get acquainted with the language.

```
// Import the required native modules
s := import("strings")
f := import("fmt")

// Declare a function to change the case of a word
func changeCase(word) {
	if word >= "a" {
		return s.ToUpper(word)
	}
  return s.ToLower(word)
}

// Loop over all received arguments
for i := 0; i < len(args); i++ {
	f.Println(changeCase(args[i]))
}
```

Now save this program to a file, and run it with `agora run <file> Welcome to Agora!`. If all goes well, it should print:

```
welcome
TO
agora!
```

Ok, so what happened here? First, the `agora run` command passes all strings following the file name to the program, as input arguments.

Then, the source code imports two stdlib modules, `strings` to convert to upper- and lowercase, and `fmt` to print to stdout.

The `changeCase` function takes a single word as argument, and checks if it is greater to or equal to "a". This is a simple way to check if it starts with a lowercase letter. If this is the case, it returns the word converted to uppercase. Otherwise it returns the word converted to lowercase.

Then there is a loop over all received arguments (obtained via the `args` reserved identifier). Each blank-separated word is sent to `changeCase` for conversion, and printed on the screen, one per line. The result of `agora run` also displays `= <nil> (runtime.null)`, because each function has an explicit `return nil` statement added if it doesn't end with a `return`.

That's it! In future versions, a `for range` construct will be available, and 	possibly a functional-style `map`, but for v0.1, that's the way to do it. Note that the `changeCase` function could also be written using the ternary `?:` operator. This is left as an exercise for the reader.

## More resources

* The [godoc source code documentation][godoc]
* The [language reference][ref]
* The [index of the documentation wiki][wiki]

[go]: http://golang.org/doc/install
[gopath]: http://golang.org/doc/code.html#GOPATH
[native]: 
[design]: 
[godoc]: http://godoc.org/github.com/PuerkitoBio/agora
[ref]: 
[wiki]: https://github.com/PuerkitoBio/agora/wiki
[bytecode]:
[assembly]:
