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
* one of the keywords `return`, `debug`, `break` or `continue`
* one of `++`, `--`, `)`, `]` or `}`

### Identifiers

An identifier is a sequence of letters, underscores and numbers. It must start with a letter.

The following identifiers are keywords in the language and may not be used as identifiers:

* if
* else
* for
* func
* return
* debug
* break
* continue

Additionally, the following identifiers are reserved and may not be used as variables:

* true
* false
* nil
* import
* panic
* recover
* len
* keys
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
* An object with a `__toBool` meta-method that returns `false`

### Nil literal

The nil value is represented with `nil`.

### Object literal

Objects are represented using the `{key: value, otherkey: value}` notation, which may be used recursively.

## Defining variables

A variable must be defined before it can be used. A new variable is introduced using the `:=` operator, which also explicitly assigns its initial value. Variables are also implicitly defined when they appear as arguments of a function, or as name of a function in the *function statement* notation, explained later.

### Scopes

All variables are declared in the scope of the function where they are defined. All module-level variables are scoped in the top-level function (the module).

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

Functions are first-class values and can be stored in variables and passed around in function arguments and return values, or in object fields. It is possible to declare a function within a function, although **closures are not supported at the moment** (returning a function from a function will not close over the variables of the parent function). This is a feature that will be added eventually.

Functions declare expected arguments by giving a list of identifiers within the parentheses of its definition. It can't declare a return value variable. Functions always return a single value, which is `nil` if there is no explicitly returned value.

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

## Statements

### Increment and decrement

Unlike in some languages such as C, and like Go, the `++` and `--` operators are statements and not expressions, they do not produce a value on the stack. So the statement `a := b++` is invalid. Those are postfix operators, they cannot be used as prefix.

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

The `for` statement can currently take three different forms: an infinite loop, a `while` equivalent, and a traditional 3-part `for`. In future versions, there will be support for a custom iterator-based `for range` notation.

The infinite loop is the most simple, it is equivalent to `for true {}` and takes the form `for { }`.

The `while` equivalent takes the form of `for <condition> { }` where `condition` evaluates to truthy or falsy. The loop continues while the condition is "truthy".

Finally, the 3-part `for` is the most traditional form, that looks like `for <init>; <condition>; <post> { }`. The `init` part is evaluated before entering the loop, then the `condition` part is evaluated, and if it is "truthy", the body of the `for` is executed. At the end of the body, the `post` statement is evaluated before returning to the `condition`, until the `condition` evaluates to "falsy".

```
for i := 0; i < 10; i++ {
    // Body
}
```

### The return statement

A return statement exits the current function. The return statement of the top-level function of the module terminates the module's execution, returning its return value to the caller. The return statement of the top-level function of the initial module returns the value to the Go host.

A function is not required to have a return statement, a default `return nil` statement is automatically added by the compiler if the last statement of the function is not a `return`.

A `return` must be followed by an expression, i.e. `return true`. This is the value that is going to be returned by the function. Only a single value can be returned.

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

A `continue` statement skips the rest of the `for` body and jumps to the execution of the `post` statement of the 3-part `for`, or to the execution of the `condition` in a `while`-equivalent `for` loop, or to the first statement of the `for` body in an infinite loop.

It is an invalid statement outside a `for` loop.

## Built-in functions

Agora has five (5) predeclared built-in functions. They are first-class function values like any other agora function, but their reserved identifier cannot be overridden.

* **import** : takes a single string value as argument, identifying a module to load and run, and returns the return value of the imported module.
* **panic** : takes a single value as argument, and if it is "truthy", raises a runtime error (a "panic") with this value. If the value is "falsy", it is a no-op and returns `nil`.
* **recover** : takes at least a single value as argument, which must be a function. If more values are provided, they are passed as arguments to the function. It executes the function and catches any error (panic) that the function may raise (it runs the function in *protected mode*). If an error is caught, it returns it, otherwise it returns `nil`.
* **len** : takes a single value as argument. If it is `nil`, returns `0`. If it is an object, returns the number of fields defined on the object (this behaviour may be overridden if the object has a `__len` meta-method). Otherwise it returns the length of the string value.
* **keys** : takes a single value as argument, which must be an object (it panics otherwise). Returns an array-like object holding all the keys of the object passed as argument. If the object has a `__keys` meta-method, it is called and its return value is returned. The order of the keys are undefined, even for an array-like object.

Because `recover` returns the eventual error, it cannot return the return value of the function that is executed. So if required, the function passed to `recover` should be a function value that stores its return value in an outer-scoped variable (eventually a closure, when the feature is added), like so:

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

Next: [Standard library](https://github.com/PuerkitoBio/agora/wiki/Standard-library)

