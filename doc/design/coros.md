# Coroutines design

```
func coro1(v) {
	n := yield v+1
	n = yield n * 2
	return n + 7
}

ret := coro1(3)
ret = coro1(ret+2)
ret = coro1(ret+3)
```

The `funcVM` already keeps its state and could be easily adjusted to return and reenter at its current position and state. At first glance, it looks like the `agoraFuncVal` should be modified to add a `vm` field that would be set only when a yield statement is executed, so that the vm is kept alive. Successive calls to the same function value would check to see if this vm field is set, and it would act as a resume instead of an initial call, with the given argument (only 1 for now) that would be passed as the return value of the call to yield from within the coroutine.

`yield` should probably be a keyword, this is very semantically close to `return`, and it produces its own opcode. However, this means that there can be ambiguities in the syntax: `a := yield 1 + 2 * 3`. Well, this is not strictly speaking ambiguous, it means yield the value 7, but to yield 3 and assign the result of yield * 3 to a, parentheses are required: `a := (yield 1 + 2) * 3`.

There is no need for a `resume` keyword or builtin, simply calling a function that has a yield statement will resume this function. The return statement in the coroutine causes the vm field to be unset, so that future calls to the function result in a new call, creating a new vm.

There probably should be a `status(fn)` built-in or keyword, to retrieve the state of the coroutine, which could be `running` (from within the coroutine, which probably requires a state field on the agora function value), `suspended` (from outside the coroutine), `dead` (once the function has returned), and `func` if the function is not an actual coroutine or has not been called yet (or the equivalent of Type(fn) if the value is not a function?).

A `cancel(fn)` built-in or keyword is also needed.

## Iterator

The `for range` notation can use these coroutines as iterators. Internally, it would be equivalent to the following calls:

- check status to make sure it returns `func`
- call the function, use as loop variable's value
- execute the loop's body
- check status to make sure it returns `suspended`
- repeat steps 2 through 4 until 4 is false (returns `dead`)

A `for range` could also receive a number as argument, in which case it would loop like a `for i := 0; i < n.Int(); i++` construct. It would also work with a string, looping over each byte.
