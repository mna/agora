# for..range notation

The `for..range` should be able to range over just about anything. It should be versatile and feel natural. It should be able to work as expected over most agora types.

## Numbers

The range could take one, two and three-args.

`for i := range(5) {}` is equivalent to `for i := 0; i < 5; i++ {}`. So the one-arg form means to loop from zero up to but not including the arg0.

`for i := range(2, 7) {}` is equivalent to `for i:= 2; i < 7; i++ {}`. So the two-arg form means to loop from arg0 up to but not including arg1.

`for i := range(2, 10, 3) {}` is equivalent to `for i := 2; i < 10; i += 3 {}`. So the three-arg form means to loop from arg0 up to but not including arg1, incrementing i by arg2.

In all forms, all args may be negative or positive.

## Strings

The range could take one, two and three-args.

`for s := range("test") {}` loops over each byte of the string arg0.

`for s := range("this is a string", " ") {}` loops over each part of the string arg0 delimited by the separator arg1.

`for s := range("this is a very long string", " ", 3) {}` loops over each part of the string arg0 delimited by the separator arg1, up to a maximum of arg2 parts.

## Objects

`for entry := range(obj) {}` loops over each key-value pair of arg0. The `entry` iterator variable is an object with two fields, "k" and "v". If the object has a "__keys" metamethod, it is used to get the keys to loop over. This may return 2 values once the multiple return values are implemented in the language.

A two-arg form could be used to loop over all key-value pairs regardless of the "__keys" metamethod. Not sure about this one.

Support a `__range` meta-method to implement and iterator on the object, giving a result very similar to range with a func? No, not sure right now, there is a meta-method to control the keys, range is expected to loop over keys.

## Funcs

Functions can be used to create iterators, using coroutines. ~~`for v := range(fn) {}` is equivalent to `reset(fn); for v := fn(); true; v = fn() { ... body ... if status(fn) == "func" { break } }`. Meaning it always loops at least once (over the value returned by the function), and exits if there are no more values to be returned from the coroutine (status is now "func"). That's because all functions return at least one value (nil or other).~~

Scratch that, actually the last return value is ignored, thanks to the loop being constructed like this: `for v := fn(); status(fn) == "suspended"; v = fn() {}`. So if the function doesn't yield (if it is a "standard" function with only return statements), it doesn't range at all.

Maybe a two-arg form to allow a "continue" from an existing coroutine, which would not run reset? Not for now, we'll see it this can be useful.

## Bools, Nil, Custom

There is no range over those type of values, there is no natural loop for those values.

## Implementation

* Use a built-in coroutine implementation for numbers, strings and objects, so that the `for..range` is ultimately always like the `for..range` over iterator funcs?

* Or try a full agora implementation, provided as a builtin and stored in the Go code as embedded resource?

* Or treat the for loop body as a function, so that a `for v := range(fn) {body}` is syntactic sugar for `fn(func(){body})`, where `fn` calls the provided func once for each iteration?

Here is the tentative implementation, that seems both simple and generic enough to support all types of range-enabled values:

-1 		: RNGS An N 		; create coro, pop n args and pass them to the range.

 0 		: RNGP An N 		; push n values - only 1 for now - from the coro on top of the range stack (the values) then push the condition bool value (true if coro is still alive, false if it ended). Use RNGP instead of a normal PUSH with a new flag, so that the (frequent) PUSH instructions stay fast.

 1 		: TEST Jf N 		; test the condition value, jump forward n instructions if false, otherwise continue, pops the condition value from the stack.

 2 		: POP  V  X 		; pop the top value from the stack (the value that was pushed by line 0) into the variable at index x, which is the iterator variable (i.e. the `v` in `for v := range something`). Could be multiple pop instructions once multiple return values are implemented.

 3..n : <loop body> 	; the instructions in the `for` loop body.

 n+1 	: JMP Jb X 			; jump back to instruction 0, start a new iteration.

 n+2 	: RNGE _ 0 			; end the coro and its goroutine, freeing resources as required (may reach this line from a break, so the coro may still be alive). Panics cause the coros to be terminated automatically, by a defer statement in the `run` method.

The `break` and `continue` loop keywords are handled as follows:

`break` : JMP Jf n+2 	; jump to the RNGE instruction, clearing the coro.

`continue` : JMP Jb 0 ; jump back to instruction 0, which is different from a normal 3-part loop, because in the case of a range, the init and post statements are necessarily the same (call the coro to get the new value).
