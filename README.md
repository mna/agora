# The agora programming language

Agora is a dynamically typed, garbage collected, embeddable programming language. It is built in the Go programming language, and is meant to provide a syntactically similar, loose, dynamic and fun companion to the statically typed, machine compiled Go language - somewhat like Lua is to C.

## Installation

`go get github.com/PuerkitoBio/agora/...`

This will install the agora packages as well as the `agora` command-line tool. See `agora -h` for help, provided the `$GOPATH/bin` path is in your exported path.

## Example

## Resources

## Changelog

### v0.1.0 / 2013-09-05 (?)

* Initial release
* Explicit goal is to take the project off the ground, nothing fancy
* The compiler is there mostly to check/test the runtime
* No optimization, start with a decent runtime design

## License

Agora is licensed under the [BSD 3-Clause License][bsd], the same as the Go programming language. The full text of the license is available in the LICENSE file at the root of the repository.

[bsd]: http://opensource.org/licenses/BSD-3-Clause
