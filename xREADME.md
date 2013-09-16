# The agora programming language

Agora is a dynamically typed, garbage collected, embeddable programming language. It is built with the Go programming language, and is meant to provide a syntactically similar, loose and dynamic companion to the statically typed, machine compiled Go language - somewhat like Lua is to C.

## Installation

`go get -t github.com/PuerkitoBio/agora/...`

This will install the agora packages as well as the `agora` command-line tool. See `agora -h` for help, provided the `$GOPATH/bin` path is in your exported path.

## Example

More examples are available in the wiki and the source code under /testdata/src, but to give a taste of the syntax, here is the usual `hello world`:

```
// Output: Hello, Agora !
fmt := import("fmt")
func greet(name) {
	fmt.Println("Hello,", name, "!")
}
greet("Agora")
```

A few things to note:

* It looks *very* similar to Go, minus the types.
* `import` is a built-in function, not a keyword. This is important with dynamically-loaded modules, it gives you control of where this overhead of loading and compiling the code is done. It returns the value exported by the module - in this case, an object that exposes methods like `Println`.
* Obviously, since this is a dynamically-typed language, arguments have no types.
* `:=` introduces a new variable. Using an undefined variable is an error, so this statement could not have been `=`.
* Statements are valid in the top-level (module) scope. That's because a module (the name for an agora file) is an implicit (top-level) function.
* Semicolons are managed just like in Go, so although they are inserted in the scanning stage, they are optional (and usually omitted) in the source code.

## Resources

* Getting started: https://github.com/PuerkitoBio/agora/wiki/Getting-started
* Documentation's index in the wiki: https://github.com/PuerkitoBio/agora/wiki
* Source code documentation on GoDoc: http://godoc.org/github.com/PuerkitoBio/agora

## Changelog

### v0.1.0 / 2013-09-17

* Initial release

## License

Agora is licensed under the [BSD 3-Clause License][bsd], the same as the Go programming language. The full text of the license is available in the LICENSE file at the root of the repository.

[bsd]: http://opensource.org/licenses/BSD-3-Clause
