# Multiple return values

The goal of this feature is to allow multiple return values to be returned from functions and methods in agora. This document explores the design possibilities coupled with the expression list change (i.e. the list of lhs variables and rhs expressions like `a, b, c := some(), exp + ress + ions`).

## Various uses

The multiple return values can be used in the following cases:

* `return a, b, c` : multiple return values from a function;

* `a, b, c = yield d, e` : multiple yielded values, multiple values returned from `yield` (multiple arguments received in a resume call of a coroutine);

* `for v := range returnTwoValues(), 10, false` : multiple values passed to `range` - all values are passed to the internal range handler, but extra values are ignored;

* `a, b := returnThreeValues()` : more return values than lhs assignments - in which case the extra return values are ignored;

* `a, b, c := returnTwoValues()` : less return values than lhs assignments - in which case the extra lhs variables are set to `nil`;

* `a, b, c := returnTwoValues(), returnThreeValues(), 10` : same scenario, just to illustrate that the same behaviour applies here as well, `a` and `b` get the return values of the first function call, `c` gets the first return value of the second function call, all other rhs values are ignored;

* `a = returnTwoValues(), returnThreeValues()`

* `returnTwoValues(), returnThreeValues()` : allow? Probably not. If the statement is not `=`, `:=`, `return`, `range` or `yield`, no comma-separated list of expressions is allowed (and function calls, of course).

* `call( returnTwoValues(), returnThreeValues(), 10 )` : in function and method calls, all values are used and passed to the function - extra arguments (all arguments) are always available via the `args` special identifier;

* `sum := 2 + returnThreeValues() + 5` : only the first value is used, and it is nil if there is no return value.

* `sum := 2 + yield + 5` : same as above.

* `obj := {a: returnThreeValues(), b: yield}` : same as above.

* `obj.a, obj.b = returnThreeValues()` : for this to work, required to push values, winnow to keep only the right number, then push b, then obj, then pop, then push a, then push obj, then pop again.

Support a `var...` notation to collect all remaining values? As in:

```
a, b, allOthers... = returnFiveValues()
```

In this case, `allOthers` would be an array-like object, indexed from 0 to 2. The same notation could be done for function arguments (and assignments from yielded values, but this is obviously the same case as the above example), and it provides Go-like variadic arguments support:

```
func anyArgs(expected, unknown...) {
	// In the call below, unknown would have a length of 3, 
	// 0: "test", 1: {b: 3}, 2: false
}
anyArgs(5, "test", {b: 3}, false)
```

## Implementation

All *producers* of value (expressions, rhs) push their values on the stack, regardless of how many will actually be used and discarded. That is because it is impossible to know at compile-time how many values from a specific function call will be required, and it would be possible but at a high complexity cost to do at runtime. This proposed solution looks simple and appears to cover all cases the same way.

Then, before assigning the values to their variables (before the POP instructions), the stack is adjusted so that only relevant values are kept, and the other ones are discarded. If there are missing values, Nils are pushed on the stack.

How this is done is that when an assignment is emitted (strictly, `=` or `:=`), prior to emitting the code for the rhs (which is itself emitted *before* the code for the lhs, so that the rhs is always fully evaluated before it is assigned to the lhs), emit a `BMKS _ 0` instruction. This instruction adds the current stack index into a stack of such bookmarks (LIFO).

Then it emits the code for the rhs, and prior to emitting the code for the lhs, it emits a `BMKE An X` instruction. This instruction pops the top bookmark from the bookmark stack, and makes sure there are exactly X values on the stack starting at the bookmark. If there are more values, they are popped until the right number remains, and if there are less, `nil`s are pushed on the stack up to that number.

The number is known by the compiler, this is the number of lhs variables in the assignment statement.

For `return`, `yield`, `range` and function call statements, the emitted code is a little different. Like assignments, prior to emitting the code for the expression list, a `BMKS _ 0` instruction is emitted. Then the expression list is emitted, and it is the `RET`, `YLD`, `RNGS`, `CALL` or `CFLD` opcode that pops the bookmark and consumes all values up to this stack index. There is no need for a `BMKE` instruction.

The signature of functions will change to return `[]Val` instead of `Val`. This means `runtime.agoraFuncVM.run()`, `runtime.Func.Call()` and `runtime.Module.Run()`, at a minimum.

It also means that functions could now return no value, so that the implicit `return nil` could be dropped (it would be changed for an implicit naked `return` instead). This should be investigated.

## Implementation v2

The previously documented implementation doesn't work well for cases such as `sum := a + returnTwoVals()`. The function would be called first and would push two values on the stack. The following ADD opcode would add those two values instead of `a` and the first returned value. This means that every function call ever would need a bookmark, more or less.

Instead, it would seem better to use the flag and index (that are now unused thanks to the bookmark approach for the number of args) to specify the number of return values to push on the stack. If the function returns less values, it pushes nils up to this number. If it returns more, exceeding values are ignored.

To recap:

* `=` and `:=` : use BKMS before executing the rhs, push all produced values on the stack, then execute a BKME with the number of values to keep before executing the lhs. Compiler knows how many values must be kept.

* `return` statement : use BKMS to mark the origin of the stack position, then return (pop) all values on the stack from this point upward. No BKME needed, the RET opcode acts as the BKME.

* `yield`, `range`, function and method arguments : use BKMS to mark the origin of the stack position for the arguments, then pop all values on the stack from this point upward in the corresponding YLD, RNGS, CALL or CFLD opcodes. No BKME needed.

* `range` return values : a range can only be used in a `for`, as the rhs of a `=` or `:=` statement, so its return values will always be managed via the BKMS and BKME of the assignment operators.

* `yield`, `?`, function and method return values : unlike the `range`, those return values can be used at any place, such as `1 + returnThree() + 2`, in which case only a single value must be used. Same in object literal: `obj := {a: returnThree()}`, or in key resolution: `obj[returnThree()] = true`. Return values may also be completly ignored, such as in `returnThree()`. It would probably be inefficient to use BMKS approach every time a function is called, and since the YLD, CALL and CFLD opcodes now have unused flag and index, it should be used on the same opcode to specify how many return values should be pushed to the stack. The emitter must look at the immediate parent of the call or yield in the AST to see how many values are required:
    - nil parent : the call is a standalone statement, return values are ignored, nothing is pushed (flag is An, ix is 0).
    - any binary op, unary op, cmp, &&/||, += type assignment, [] key resolution : only one value is pushed (nil if none). Flag is An, ix is 1.
    - := or = : stack is controlled by BKMS and BKME, so push all values. Flag is _, ix is ignored, means push all.
    - () function call : all values are pushed, they will be used as arguments, and all values are used in arguments to a function call. Flag is _, ix is ignored.
    - {} literal object : only one value must be used, to set the field of the object. Flag is An, ix is 1.
    - ? ternary op : if the call is in the first part, only one value must be pushed (the condition). In the parts 2 and 3, multiple values should be allowed to be returned, and follow the same conditions as function calls.
    - if statement : only one value.
    - for statement : condition part must use only one value. Part 1 or 3 is either zero (ignored return values) or controlled by assignments BKMS and BKME.
    - yield statement : return all values, much like when parent is a function call.
    - return statement : return all values.
