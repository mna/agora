The standard library is voluntarily small and minimal for this early release. As the language gains features and stabilizes, the right way to offer APIs will become more obvious, and the major use-cases of the language will be better known, allowing for better decisions regarding what makes sense to include in the stdlib.

There are currently seven (7) stdlib modules:

* **conv** to provide value conversions and type inspection.
* **filepath** to provide file path manipulation functions, a subset of Go's `path/filepath` package.
* **fmt** to provide formatted I/O, a subset of Go's `fmt` package.
* **math** to provide the usual mathematical functions, a subset of Go's `math` and `math/rand` packages.
* **os** to provide file access and process manipulation, a subset of Go's `os`, `os/exec` and `io/ioutil` packages.
* **strings** to provide string manipulation functions and regular expressions, a subset of Go's `strings` and `regexp` packages.
* **time** to provide date and time functions and types, a subset of Go's `time` package.

## conv

The `conv` module exposes the following methods:

* **Number(val)** : converts val to a number, returns the number.
* **String(val)** : converts val to a string, returns the string.
* **Bool(val)** : converts val to a boolean, returns the boolean.
* **Type(val)** : checks the type of val, returns a string representing the type. The possible return values are `string`, `number`, `bool`, `func`, `object` or `nil`.

## filepath

* **Abs(val)** : returns the absolute path of val. It may panic.
* **Base(val)** : returns the last element of val.
* **Dir(val)** : returns all but the last element of val.
* **Ext(val)** : returns the extension of the last element of val. The extension is the suffix of the last element starting at the last dot.
* **IsAbs(val)** : returns true if val is an absolute path.
* **Join(vals...)** : joins any number of path elements into a single path, and returns the resulting path.

## fmt

* **Print(vals...)** : prints the vals to stdout.
* **Println(vals...)** : prints the vals to stdout, then prints a newline.
* **Scanln()** : reads text up to a newline character from stdin.
* **Scanint()** : reads and returns an integer value from stdin.

## math

* **Pi** : number field that holds the Pi value.
* **Abs(val)** : returns the absolute value of val.
* **Acos(val)** : returns the arccosine of val.
* **Acosh(val)** : returns the inverse hyperbolic cosine of val.
* **Asin(val)** : returns the arcsine of val.
* **Asinh(val)** : returns the inverse hyperbolic sine of val.
* **Atan(val)** : returns the arctangent of val.
* **Atan2(val1, val2)** : returns the arctangent of val1/val2.
* **Atanh(val)** : returns the inverse hyperbolic tangent of val.
* **Ceil(val)** : returns the ceiling of val.
* **Cos(val)** : returns the cosine of val.
* **Cosh(val)** : returns the hyperbolic cosine of val.
* **Exp(val)** : returns the base-e exponential of val.
* **Floor(val)** : returns the floor of val.
* **Inf(val)** : returns positive infinity if val >= 0, negative infinity otherwise.
* **IsInf(val1, val2)** : returns true if val1 is infinity according to the sign of val2.
* **IsNaN(val)** : returns true if val is not a number (NaN).
* **Max(vals...)** : returns the maximum value of all vals.
* **Min(vals...)** : returns the minimum value of all vals.
* **NaN()** : returns the not-a-number (NaN) value.
* **Pow(val1, val2)** : returns the base-val1 exponential of val2.
* **Sin(val)** : returns the sine of val.
* **Sinh(val)** : returns the hyperbolic sine of val.
* **Sqrt(val)** : returns the square root of val.
* **Tan(val)** : returns the tangent of val.
* **Tanh(val)** : returns the hyperbolic tangent of val.
* **RandSeed(val)** : initializes the random generator with the val seed.
* **Rand([val1[, val2]])** : returns a random value >= 0. If val1 is provided, it is used as the higher bound. If both val1 and val2 are provided, val1 is the inclusive lower bound, val2 is the higher bound.

Next: [Native Go API](https://github.com/PuerkitoBio/agora/wiki/Native-Go-API)

