The standard library is voluntarily small and minimal for this early release. As the language gains features and stabilizes, the right way to offer APIs will become more obvious, and the major use-cases of the language will be better known, allowing for better decisions regarding what makes sense to include in the stdlib.

There are currently five (5) stdlib modules:

* **conv** to provide value conversions and type inspection.
* **fmt** to provide formatted I/O, a subset of Go's `fmt` package.
* **math** to provide the usual mathematical functions, a subset of Go's `math` and `math/rand` packages.
* **os** to provide file access and process manipulation, a subset of Go's `os`, `os/exec` and `io/ioutil` packages.
* **strings** to provide string manipulation functions and regular expressions, a subset of Go's `strings` and `regexp` packages.

Next: [Native Go API](https://github.com/PuerkitoBio/agora/wiki/Native-Go-API)

