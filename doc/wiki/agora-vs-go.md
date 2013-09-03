# Similarities and differences between Go and Agora

Agora is syntactically very similar to Go. But its goal is obviously not to be a clone, so there are also important differences. Most of those differences support the goal of making a *looser*, dynamic and simple interpreted language that complements Go. This article documents these similarities and differences.

## Similarities

* The first and most immediately obvious similarity is the syntax. Like Go, Agora uses curly braces to delimit code blocks, and is semicolon-free (semicolons are *optional* and are automatically added at the scanner stage of the compiler - like Go).

* The agora keywords are all the same as the Go keywords - although agora has less. They are `if`, `else`, `for`, `func`, `return`, `debug`, `break` and `continue`. 

* Many operators are also the same. They are `+, -, *, /, %, ==, !=, >, <, >=, <=, !`, the logical operators `&&` and `||`, the increment and decrement `++` and `--`, and the assignment operators `=, :=, +=, -=, *=, /=, %=`. Operator precedence is also the same as Go.

* Field access can be done via the dot `.` operator, like Go.

* Function calls use parentheses around the parameter list (even if there is no parameter), while `if` and `for` statements don't require parentheses.

* The `nil` literal represents a null value, and boolean literals are `true` and `false` as expected.

* `len` is a built-in that returns the length of a string, or the number of items in an object (roughly similar to a map).

* `import`, `panic` and `recover` are present and do mostly what is expected, but with subtle differences, see the next section.

* Comments use the same notation as Go, namely `//` for one-liners, and `/* */` for multi-line.

## Differences

* The most obvious and important difference is probably the lack of types. Variables, arguments and functions have no types. Values have a type (a value is either a `string`, a `number`, a `bool`, an `object`, a `func` or `nil`).

* Code is organized in single-file *modules*. Statements can (and must, essentially!) appear at the top-level of the module. This is implicitly the top-level function of the module, the code that gets executed when the module is imported.

* There is no uppercase rule for what is or isn't exported by a module. The only thing a module exports is the return value of the top-level function. Anything else is private inside the module. To return multiple functions (say, like the `fmt` module), an object is returned. However, the uppercase-is-public is still used as a convention (for instance, all stdlib's functions start with an uppercase).

* There is a ternary `condition ? iftrue : iffalse` operator.

* Unlike Go, `import` is a built-in function, not a keyword that must appear at the top of the package. So it can be called wherever makes most sense, since this can be a costly operation (loading from a file, compiling, executing). As mentioned previously, it returns the value returned by the module and must be stored in a variable (there is no implicit "variable" derived from the import path).

* There are no multiple return values, functions can only return a single value.

* There are *truthy* and *falsy* values other than the boolean `true` and `false`. Namely, the `0` number, the empty string, `nil` and `false` are all false. An object can also be false if it provides a meta-method `__toBool`.

* The `panic` built-in function takes a value and raises an error with it. However, it *doesn't* raise if the value is falsy. This is symmetric with the behaviour of `recover`.

* The `recover` built-in function takes a function as parameter, and executes it in protected mode. By default, a runtime error is a panic, and stops all agora code execution to return the error to the Go host program (the second value returned from `module.Run()`). To *catch* errors in agora code, the `recover` function must be used, this is what is called the protected mode (similar to Lua error handling). `recover` runs the provided function, and if an error occurs, it catches it and returns it. Otherwise it returns `nil`.

* Field access can *also* be done using an array-like syntax, `object["key"] = value`. Any type can be used as the key, and when assigning to a field using the `.` operator, the key is implicitly a string (denoted by the identifier of the field), so `object.key` is equivalent to `object["key"]`.

* There are no bitwise operators.

* There is no multi-thread (goroutine) or channel support. Agora code must be single-threaded, although different execution contexts *can* be run in parallel.

* There are no slices or maps, the only compound data structure is the object. It can represent a slice and a map, and eventually will be optimized for when used as a slice/array.
'''
