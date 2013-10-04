This document will evolve to eventually be a complete language reference, with a formal grammar definition. For this early release, it will remain an informal yet thorough coverage of the language's syntax and features.

## Source code representation

The source code files must be encoded in UTF-8.

## Lexical elements

### Comments

Agora supports two forms of comments:

* Line comments start with `//` and end at the end of the line
* Block comments start with `/*` and end at the next `*/`

### Semicolons

Statements are terminated with a semicolon, but the `;` may be omitted in the source code. The scanner stage of the compiler automatically inserts the semicolons if the last token on the line is:

* an identifier
* a literal value
* one of the keywords `return`, `yield`, `debug`, `break` or `continue`
* one of `++`, `--`, `)`, `]` or `}`

### Identifiers

An identifier is a sequence of letters, underscores and numbers. It must start with a letter or an underscore.

The following identifiers are keywords in the language and may not be used as identifiers:

* if
* else
* for
* func
* return
* debug
* break
* continue
* yield
* range

Additionally, the following identifiers are reserved and may not be used as variables:

* true
* false
* nil
* import
* panic
* recover
* len
* keys
* string
* number
* bool
* type
* status
* reset
* this
* args

### Operators and delimiters

The following symbols represent operators and delimiters in the language:

* ( ) [ ] { }
* . , ; :
* + - * / % ! && || ?
* == != < <= > >=
* = := += -= *= /= %=
* ++ --

### Number literals

Number literals can be represented as integers or floats. At the moment there is an inconsistency between what is accepted by the compiler and what can be used. Only base-10 notation should be used for now, i.e. `42`, and floating-points should use the integer - decimal point - fraction notation i.e. `3.1415`.

### String literals

At the moment there is an inconsistency between what is accepted by the compiler and what can be used. Only string literals within double quotes should be used, i.e. `"this is a string"`. It may not contain newlines, but escape characters can be used (i.e. `\n` for newline).

### Boolean literals

Booleans are represented with the `true` and `false` literal values. However, in addition to the true boolean values, agora treats some values as "truthy" and "falsy". It is easier to list the "falsy" values, everything else being "truthy":

* The `false` boolean value
* The empty string value
* The number zero (`0`)
* The `nil` value
* An object with a `__bool` meta-method that returns `false`

### Nil literal

The nil value is represented with `nil`.

### Object literal

Objects are represented using the `{key: value, otherkey: value}` notation, which may be used recursively. Using this literal notation, the keys are treated as strings.

## Defining variables

A variable must be defined before it can be used. A new variable is introduced using the `:=` operator, which also explicitly assigns its initial value. Variables are also implicitly defined when they appear as arguments of a function, or as name of a function in the *function statement* notation, explained later.

### Scopes

All variables are declared in the scope of the function where they are defined. All module-level variables are scoped in the top-level function (the module). Functions declared within another function can access variables in the parent functions, provided they are declared before the funtion that uses them. Closures are also supported.

The only way to expose information is to return a value. When a module imports another module, it only gets access to the value returned by the imported module. With the object type, using different keys, it is possible to expose multiple functions and values.

## Functions

An agora source file is called a "module" and it is implicitly a function, without the `func` keyword. It is often called the top-level function.

Other functions are introduced using the `func` keyword. It can be used as a statement and as an expression:

```
// Statement
func myFunc() { return "statement" }
// Expression
myVar := func() { return "expression" }
```

In statement form, the "name" of the function is in fact a variable in the scope of the parent of the function being declared. The above example is equivalent to this:

```
myFunc := func() { return "statement" }
```

The expression form is self-explanatory.

Functions are first-class values and can be stored in variables and passed around in function arguments and return values, or in object fields. It is possible to declare a function within a function, and returning a function from a function will close over the variables of the parent function(s).

Functions declare expected arguments by giving a list of identifiers within the parentheses of its definition. It can't declare a return value variable. Functions always return a single value, which is `nil` if there is no explicitly returned value or in case of a naked `return` statement.

```
func Add(x, y) {
    return x + y
}
```

Functions may receive more or less arguments than expected. In the former case, the extra arguments can be retrieved via the `args` reserved identifier, which is an array-like object that holds *all* arguments passed to the function, at keys `0` to `len(args)-1`. In the latter case, the extra argument variables have the `nil` value.

If the function was assigned to an object's field, and was called with the object notation, then its `this` reserved identifier is set to the object.

```
obj := {name: "Martin"}
obj.MyFunc = func() {
    return this.name
}
obj.MyFunc()
```

If the same function is stored in a variable and called *not* with the object notation, the `this` identifier is `nil`.

```
noThis := obj.MyFunc
noThis() // Error
```

Functions can be coroutines, meaning that they can `yield` a value and execution to a caller function, and re-enter execution at a later time, after the `yield` statement:

```
// Example of a coroutine
func fn(n) {
	i := yield n + 1
	i = (yield i * 2) + 1
	return i * 3
}
fmt := import("fmt")
fmt.Println(fn(1)) // outputs 2
fmt.Println(fn(2)) // outputs 4
fmt.Println(fn(3)) // outputs 12
fmt.Println(fn(4)) // outputs 5 (restarts the function)
```

## Operators

Most operators have the obvious meaning.

* `+` : adds two values
* `-` : subtracts two values, or unary minus of a single value, depending on context
* `*` : multiplies two values
* `/` : divides two values
* `%` : returns the modulo of two values
* `==` : compares two values for equality
* `!=` : compares two values for inequality
* `<` : compares two values for lower-than
* `>` : compares two values for greater-than
* `<=` : compares two values for lower-than or equal
* `>=` : compares two values for greater-than or equal
* `?:` : ternary operator, checks the initial condition before the `?`, if true, evaluates the expression after the `?`, if false, evaluates the expression after the `:`
* `&&` : boolean "and" of two values
* `||` : boolean "or" of two values
* `!` : boolean negation of a value

### Assignment operators

* `:=` : defines a new variable and assigns a value to it
* `=` : assigns a value to an existing variable (or a field of an existing variable, if it is an object)
* `+=` : adds a value to an existing variable, and assigns it to itself
* `-=` : subtracts a value from an existing variable, and assigns it to itself
* `*=` : multiplies a value by an existing variable, and assigns it to itself
* `/=` : divides a value from an existing variable, and assigns it to itself
* `%=` : computes the modulo of an existing variable with a value, and assigns it to itself
* `++` : adds 1 to an existing variable, and assigns it to itself
* `--` : subtracts 1 from an existing variable, and assigns it to itself

### Arithmetic and comparison operations

All binary arithmetic operations (`+`, `-`, `*`, `/`, `%`) are defined on numbers. The `+` is also defined on strings, resulting in a concatenation of both values. The unary minus operation is defined on numbers.

Also, all arithmetic operations can be defined on objects, using the relevant meta-method (i.e. `__div` for `/`). If any of the operands is an object with the correct meta-method, the operation will be executed via this meta-method, using the left operand's meta-method if applicable, otherwise the right operand's.

Using arithmetic operations with any other value type results in a runtime error.

All types of values can be compared. For values of the same type, numbers, strings and booleans have the expected ordering (for booleans, `true` is greater than `false`). Nil can only be equal to itself. Objects without the `__cmp` meta-method, functions and custom values can be equal, but always return the first operand as `lower than` if `<` or `>` is requested (there is no logical ordering possible).

As for arithmetic operations, if an object with the `__cmp` meta-method is an operand, this function is called to execute the comparison, regardless of the type of the other value. The left operand's meta-method is called if applicable, otherwise the right operand's.

The full matrix of arithmetic and comparison behaviour is available in this spreadsheet:
https://docs.google.com/spreadsheet/ccc?key=0Atx1KnJmATDcdEV1TGhYTmxGWjRTbjBvdy00aWczRHc&usp=sharing

## Statements

### Increment and decrement

Unlike in some languages such as C, and like Go, the `++` and `--` operators are statements and not expressions, they do not produce a value on the stack. So the statement `a := b++` is invalid. Those are postfix operators, they cannot be used as prefix. This is subject to change and there is an open issue about this (#5).

### The if statement

The `if` statement evaluates the condition next to the `if` keyword, and if it is "truthy", it executes the statements in the body of the `if` (note that "truthy" is different than the stricter "boolean true").

An optional `else` statement may be present. The statements within the `else` block are executed if the `if` condition is "falsy". The `else` part may introduce another `if`.

Parentheses are not required around the `if` condition.

```
if "mystring" && 0 {
    // This won't execute because 0 is falsy
} else if myVar > 38 {
    // This depends on the value of myVar
} else {
    // Otherwise this is executed
}
```

### The for statement

The `for` statement can take four different forms: an infinite loop, a `while` equivalent, a traditional 3-part `for` and a `for range`.

The infinite loop is the most simple, it is equivalent to `for true {}` and takes the form `for { }`.

The `while` equivalent takes the form of `for <condition> { }` where `condition` evaluates to truthy or falsy. The loop continues while the condition is "truthy".

The 3-part `for` is the most traditional form, that looks like `for <init>; <condition>; <post> { }`. The `init` part is evaluated before entering the loop, then the `condition` part is evaluated, and if it is "truthy", the body of the `for` is executed. At the end of the body, the `post` statement is evaluated before returning to the `condition`, until the `condition` evaluates to "falsy".

```
for i := 0; i < 10; i++ {
    // Body
}
```

The `for range` notation allows iteration over the following types of value:

* Number
* String
* Func
* Object

It panics if the value is of another type. The range over numbers supports 3 different args:

`for v := range [start,] max[, increment]`

The range over strings also supports 3 different args:

`for v := range str[, sep[, max]]`

It loops over each byte of the string if `sep` is empty or nil, otherwise it loops over parts of the string separated by the specified separator. In any case, it loops over a maximum of `max` values if it is >= 0.

The range over functions calls the iteration function until the `return` statement is reached, excluding the value returned by `return`. In other words, it loops over all values returned by `yield` statements. This is necessary because all functions have an implicit `return nil` statement, so otherwise it wouldn't be possible to have such a range loop 0 time. Any subsequent values after the function value get passed as argument to the function.

The range over objects loops over the keys of the object, returning an object with two keys, `k` and `v` (holding the key and value, respectively).

### The return statement

A return statement exits the current function. The return statement of the top-level function of the module terminates the module's execution, returning its return value to the caller. The return statement of the top-level function of the initial module returns the value to the Go host.

A function is not required to have a return statement, a default `return nil` statement is automatically added by the compiler if the last statement of the function is not a `return`.

A `return` can be followed by an expression, i.e. `return true`. This is the value that is going to be returned by the function. Only a single value can be returned. An empty `return` is equivalent to `return nil`.

### The break statement

A `break` statement terminates the execution of the innermost `for` loop. Agora does not support labels, so it cannot break multiple embedded loops. It is an invalid statement outside a `for` loop.

```
for {
    if age > 40 {
        break
    }
}
```

### The continue statement

A `continue` statement skips the rest of the `for` body and jumps to the execution of the `post` statement of the 3-part `for`, or to the execution of the `condition` in a `while`-equivalent `for` loop (or a `for range` loop), or to the first statement of the `for` body in an infinite loop.

It is an invalid statement outside a `for` loop.

### The range statement

The `range` statement is used in `for` loops and is explained in the `for` statement section.

### The yield statement

The `yield` statement is used to return values to the caller and suspend a function's execution, while waiting to resume after this statement. This effectively turns the function into a coroutine. `yield` returns a value to the caller, but also returns a value to the coroutine once it is resumed.

A coroutine is resumed simply by calling the function again.

## Built-in functions

Agora has eleven (11) predeclared built-in functions. They are first-class function values like any other agora function, but their reserved identifier cannot be overridden.

* **import** : takes a single string value as argument, identifying a module to load and run, and returns the return value of the imported module.
* **panic** : takes a single value as argument, and if it is "truthy", raises a runtime error (a "panic") with this value. If the value is "falsy", it is a no-op and returns `nil`.
* **recover** : takes at least a single value as argument, which must be a function. If more values are provided, they are passed as arguments to the function. It executes the function and catches any error (panic) that the function may raise (it runs the function in *protected mode*). If an error is caught, it returns it, otherwise it returns `nil`.
* **len** : takes a single value as argument. If it is `nil`, returns `0`. If it is an object, returns the number of fields defined on the object (this behaviour may be overridden if the object has a `__len` meta-method). Otherwise it returns the length of the string value.
* **keys** : takes a single value as argument, which must be an object (it panics otherwise). Returns an array-like object holding all the keys of the object passed as argument. If the object has a `__keys` meta-method, it is called and its return value is returned. The order of the keys are undefined, even for an array-like object.
* **number** : converts a value to a number.
* **string** : converts a value to a string.
* **bool** : converts a value to a boolean.
* **type** : returns the type of a value, namely `number`, `string`, `bool`, `func`, `object`, `nil` or `custom`.
* **status** : returns the coroutine status of a function, which can be empty string ("") if it isn't a coroutine, `running` if the coroutine is currently in execution, and `suspended` if it is in `yield` state, waiting to resume.
* **reset** : resets a coroutine function so that the next call to the function restarts its execution from the beginning.

Because `recover` returns the eventual error, it cannot return the return value of the function that is executed. So if required, the function passed to `recover` should be a function value that stores its return value in an outer-scoped variable, or a closure, like so:

```
a := nil
err := recover(func() {
	a = "storing return value"
})
if err {
	// Handle error
}
return a
```

## Objects

An object can have keys of any value except `nil`. The dot notation implicitly creates a string key, so `obj.key = 3` is equivalent to `obj["key"] = 3`. The `[]` notation is required to create keys of other types. Assigning `nil` to an object's key removes the key from the object.

The following meta-methods are currently supported, so that an object's behaviour can be overridden:

* **__int** : converts the object to an integer value.
* **__float** : converts the object to a float value.
* **__bool** : converts the object to a boolean value.
* **__string** : converts the object to a string value.
* **__native** : converts the object to a native Go value.
* **__cmp** : compares the object to another value, returning 1, 0 or -1 if greater, equal or lower than the value.
* **__add** : adds a value to the object.
* **__sub** : subtract a value from the object.
* **__mul** : multiply a value with the object.
* **__div** : divide a value from the object.
* **__mod** : gets the module of the object divided by a value.
* **__unm** : gets the unary minus operation of the object.
* **__len** : gets the length of the object.
* **__keys** : gets the keys of the object.
* **__noSuchMethod** : defines a method to call on the object if an unknown method is called.


Next: [Standard library](https://github.com/PuerkitoBio/agora/wiki/Standard-library)

