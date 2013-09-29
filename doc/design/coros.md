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

The `funcVM` already keeps its state and can be easily adjusted to return and reenter at its current position and state. The `agoraFuncVal` should be modified to add a `coroState` field that would be set to its `funcVM` only when a yield statement is executed, so that the vm is kept alive. Successive calls to the same function value check to see if this vm field is set, and it acts as a resume instead of an initial call, with the given argument (only 1 for now) that is passed as the return value of the call to yield from within the coroutine.

`yield` should probably be a keyword, this is very semantically close to `return`, and it produces its own opcode. However, this means that there can be ambiguities in the syntax: `a := yield 1 + 2 * 3`. Well, this is not strictly speaking ambiguous, it means yield the value 7, but to yield 3 and assign the result of yield * 3 to a, parentheses are required: `a := (yield 1 + 2) * 3`.

There is no need for a `resume` keyword or builtin, simply calling a function that has a yield statement will resume this function. To put it another way, *all* agora functions are coroutines - it's just that, like Donald Knuth said, a subroutine is just a special of a coroutine, one without a yield statement. The return statement in the coroutine causes the `coroState` field to be unset, so that future calls to the function result in a new call, creating a new vm.

There probably should be a `status(fn)` built-in or keyword, to retrieve the state of the coroutine, which could be `running` (from within the coroutine, which probably requires a state field on the agora function value - no, in hindsight, a `IsRunning` method on the `Ctx` does the trick), `suspended` (from outside the coroutine), and `func` if the function has not been called yet or is ready to restart. It panics if the value is not a function.

A `cancel(fn)` built-in or keyword is also needed. Or not. After further reflection, the built-in should be `reset(fn)`, since `cancel` sounds like the thing is now dead and cannot be used, which is not the case with agora, it just nils the `coroState` field so that the next call will be an initial function call instead of a continuation from the latest yield.
