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
