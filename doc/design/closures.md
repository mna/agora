# Closures design

A variable *v* must be closed-over once the function *f1* in which it is declared goes out of scope (returns), iif:

* There is at least one function *f2* declared inside *f1*;
* The function *f2* refers to the variable *v*;
* The function *f2* is still accessible after *f1* returns.

*Closing-over* may refer to (at least) two behaviours:

* The value of *v* is **copied** to *f2*'s environment;
* The value of *v* is stored *on the heap* and is **referenced** from *f2*'s environment.

The first approach means that two closures referring to the same *v* inherit independent copies of *v* after the closure is created.

The second approach allows two closures referring to the same *v* to share its value, which is probably the expected behaviour, or at least it is the behaviour in Go (http://play.golang.org/p/qhiXgw9Bfn).

## Solution

This solution is simple, easy to implement, and plays well with the GC. It is O(n) where *n* is the depth of the function, but it will be <= 3 in common cases.

It involves splitting the current single representation of `Func` for both the prototype and the actual value into two separate representations.

* `agoraFunc` will represent the prototype, which means all the static information about the function.

* `agoraFuncVal` will represent a function value, and will be created each time a new function value is produced at runtime and will hold both its prototype and its runtime context - a linked list of reachable local variables (i.e. allows for closures, currying, and normal funcs). It will be created something like `newFuncVal(proto, parentVM)`.

* `funcVM` will still represent one running instance of a function value, with its stack, state, local variables, etc.

Native functions are left untouched, since they can't close over agora variables.

The get/setVar methods of the `Ctx` will be impacted. Resolving a variable will now mean looking at the function VM's list of locals, and its function value's runtime context, then its function value's parent's runtime context, up the hierarchy. There will be dead code to remove in `Ctx` regarding the lexical scope/current variable resolution.

The upside of this solution is that closures come for free (no need to do "escape analysis" and move upvalues at a specific time), as long as the function value is reachable, it keeps its full context.

The downside is that *all* variables are kept, not just the ones closed-over. This will have to be done ultimately, so that unreferenced values get GCed. Hence this is a *temporary* solution to make the feature available, but will need to be corrected in a future version.

