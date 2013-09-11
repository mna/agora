The v0.1 release was a huge step to take this project off the ground, but it is also a **very** rough, **alpha** release. There is still a *lot* to do before even thinking using agora in a production-like environment, and before calling the language *stable*.

Everything except the major goals of the language should be considered unstable, and everything could change. The native API is not frozen, the VM is not frozen, the compiler is definitely not frozen (thank god!) and even the syntax, runtime design and stdlib are not frozen. What *is* frozen is that this is and will remain a dynamically-typed, dynamically-loaded, garbage collected, embeddable language with a syntax close to Go.

The following is by no means a rigid contract, more like an overly optimistic overview of the projected roadmap.

## v0.2

The next version will yet again focus on the language features and the runtime. The compiler could definitely use some work (and the command-line tool could be cleaned-up), but I don't feel this is a priority as long as the language is not yet stabilized. The compiler is *barely* good enough, but it *is* good enough for running and testing the language. It will probably only change when it is required for new features to be supported.

The features envisioned for v0.2 are:

* Change the generated instructions so that `&&` and `||` are short-circuiting.
* Support the immediately-invoked function expression syntax (IIFE, i.e. `func(){}()`).
* Support closures.
* Support coroutines.
* Using coroutines, support a flexible, extensible `for ... range` notation.
* Better/more tests.
* Fix bugs.

This will probably be enough for this version.

## v0.3

The backlog will no doubt have time to fill up until then, but as it stands today, the next.next version will mostly be about the language and the runtime too. Expected features:

* Support embedded/prototype/metatable inheritance, à-la Javascript or Lua.
* Support a literal "array" notation (i.e. `a := [10, true, "hi"]`).
* Optimize the object when used as an array (dense integer keys).
* Support comma-separated list of assignments (for `:=` and `=`).
* Introduce the `switch` statement.
* Support interacting with Go channels.
* Better/more tests.
* Fix bugs.

## Beyond

Once the language features have landed, focus should turn to refactoring the compiler, returning better error messages (both in the compiler and the runtime), and stabilizing/enhancing the stdlib. Testing on Windows, BSD and other Go-supported platforms should also start somewhere around here, I don't have a Windows nor BSD machine and I don't want to spend time debugging for an OS-specific problem while the project is evolving fast. As far as I can tell, at the moment at least, there is no OS-specific code in agora so it shouldn't be too hard to get it to work. The only thing that comes to mind is the directory separator for Windows.

## Development

Development will take place in the `next` branch. I will keep `master` stable and updated only with official releases (there *may* be v0.1.n releases in case of absolute disaster before v0.2 lands, and if this happens, it will be on `master`).

Next: [Virtual machine](https://github.com/PuerkitoBio/agora/wiki/Virtual-machine)

