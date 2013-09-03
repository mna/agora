# The agora programming language

Agora is a dynamically typed, garbage collected, embeddable programming language. It is built in the Go programming language, and is meant to provide a syntactically similar, loose, dynamic and fun companion to the statically typed, machine compiled Go language - somewhat like Lua is to C.

## Installation

`go get github.com/PuerkitoBio/agora/...`

This will install the agora packages as well as the `agora` command-line tool. See `agora -h` for help, provided the `$GOPATH/bin` path is in your exported path.

## Example

More examples are available in the wiki and the source code, but to give a taste of the syntax, here is the usual `hello world`:

```
// Output: Hello, Agora !
func greet(name) {
	fmt := import("fmt")
	fmt.Println("Hello,", name, "!")
}
greet("Agora")
```

A few things to note:

* It looks *very* similar to Go, minus the types.
* `import` is a built-in function, not a keyword. This is important with dynamically-loaded modules, it gives you control of where this overhead of loading the code is done. It returns the value exported by the module - in this case, an object that exposes methods like `Println`.
* Obviously, since this is a dynamically-typed language, arguments have no types.
* `:=` introduces a new variable. Using an undefined variable is an error, so this statement could not have been `=`.
* Statements are valid in the top-level (module) scope. That's because a module is an implicit (top-level) function.
* Semicolons are managed just like in Go, so although they are inserted in the scanning stage, they are optional (and usually omitted) in the source code.

## Resources

* Source code documentation on GoDoc: http://godoc.org/github.com/PuerkitoBio/agora

## Changelog

### v0.1.0 / 2013-09-05 (?)

* Initial release
* Explicit goal is to take the project off the ground, nothing fancy, not much more
* The compiler is there mostly to check/test the runtime, it is ugly
* No optimization, start with a decent runtime design

## License

Agora is licensed under the [BSD 3-Clause License][bsd], the same as the Go programming language. The full text of the license is available in the LICENSE file at the root of the repository.

[bsd]: http://opensource.org/licenses/BSD-3-Clause
