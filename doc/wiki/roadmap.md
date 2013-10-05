Like the v0.1, the v0.2 release is still a **very** rough, **alpha** release. There is still a *lot* to do before even thinking using agora in a production-like environment, and before calling the language *stable*.

Everything except the major goals of the language should be considered unstable, and everything could change. The native API is not frozen, the VM is not frozen, the compiler is definitely not frozen (thank god!) and even the syntax, runtime design and stdlib are not frozen. What *is* frozen is that this is and will remain a dynamically-typed, dynamically-loaded, garbage collected, embeddable language with a syntax close to Go.

The following is by no means a rigid contract, more like an overly optimistic overview of the projected roadmap.

## v0.3

The next version will mostly be about the language and the runtime too. Expected features:

* Support embedded/prototype/metatable inheritance, à-la Javascript or Lua.
* Support a literal "array" notation (i.e. `a := [10, true, "hi"]`).
* Optimize the object when used as an array (dense integer keys).
* Support comma-separated list of assignments (for `:=` and `=`).
* Introduce the `switch` statement.
* Multiple return values.
* Review the `this` behaviour.
* Better/more tests.
* Fix bugs.

## v0.4

The next.next version will focus on the compiler (most probably a full rewrite) and the command-line tool cleanup.

## Beyond

Once the language features have landed and the compiler is finally decent, focus should turn to benchmarks, profiling and optimizations. The stdlib and better cross-platform support will get some love at this point too.

## Development

Development takes place in the `next` branch. I will keep `master` stable and updated only with official releases (there *may* be v0.2.n releases in case of absolute disaster before v0.3 lands, and if this happens, it will be on `master`).

Next: [Back to wiki index](https://github.com/PuerkitoBio/agora/wiki)

