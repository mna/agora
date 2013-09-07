This article explains the internals of the agora virtual machine. All agora code is interpreted at runtime by a stack-based virtual machine that loops through the instructions until an unrecovered error (a panic in the runtime) is raised or until a return statement is executed.

All executable code is contained within a function, even the top-level module statements, which are implicitly part of the top-level function. This explains why the virtual machine is located in the /runtime/funcvm.go file.

When the execution context starts running an agora program (via `runtime.Module.Run(...)`), what happens is that the top-level function of this module (the one at index 0 in the list of function prototypes) is called. Calling an agora function results in the following steps:

* A function VM is instantiated for the function prototype. This is similar to the relationship between a class and an instance of the class - the many instances of the same function prototype may be created and run, and each one maintains its own separate state, like the *program counter*, *stack* and *stack pointer*.
* The `this` keyword is set for the instance of the function. It is `nil` unless the function is called on an object (i.e. with the syntax `obj.FunctionField(args)`).
* The function is pushed onto the frames stack of the execution context.
* The function is executed.
* On return, the function is popped from the frames stack of the execution context.

The rest of this article will focus on "the function is executed" part.

## The funcVM type

The `runtime.funcVM` struct is defined like this:

```Go
type funcVM struct {
	proto *AgoraFunc
	pc    int
	vars  map[string]Val
	stack []Val
	sp    int
	this  Val
	args  Val
}
```

It holds a reference to its function prototype `proto` (which itself holds a reference to the execution context). It also has a program counter field (`pc`) that points to the next instruction to process. The `vars` field holds the local variables. The `stack` is the central place where values are manipulated, and `sp` is the stack pointer that points to the next free index in the stack. Finally, `this` and `args` hold the values of these special reserved identifiers.

The `run(...Val) Val` method is where execution takes place. The first thing it does is assign the values of the parameters' variables. This is why the *expected arguments* function header field is so important, the VM assigns the first *n* values received as arguments to those variables stored in the K table at indices 0..n-1 (the function's arguments variables must *always* be stored as the first K symbols, starting at index 0). If the function received less arguments than expected, the remaining variables are set to `nil`.

Then it creates the `args` reserved identifier's value, which is an array-like object holding all received arguments. This is stored in the `funcVM.args` field.

And now it is ready to enter the execution loop, which is an infinite loop that processes instructions. It starts at the instruction at index 0 in the I section and decodes it into and opcode (`op`), a flag (`flg`) and an index (`ix`), and immediately increment the `pc` field to point to the next expected instruction (if there is a jump, it will override this value). An instruction is a 64-bit value where the most significant byte is the opcode, the second-most significant byte is the flag, and the remaining 6 bytes is the index.

Then comes the `switch` on the opcode. The only one that can exit the execution loop is `OP_RET` which is the return statement, which is why the compiler automatically adds a `return nil` at the end of each function if the last instruction is not a `return`.

The full list of opcodes is available in /bytecode/opcodes.go, while the list of flags is in /bytecode/instr.go. The next section explains the behaviour of each opcode.

## The opcodes

* **RET** : pops one value from the stack and returns it, ending the function's execution.
* **PUSH** : gets the value identified by `flg` and `ix`, depending on the flag, and pushes it on the stack:
    - **K** : the constant value at index `ix` in the K table.
    - **V** : the variable identified by the string at index `ix` in the K table. It can be a local variable, or a variable reachable in the current scope (defined in a function in the outer-scope). There are no closures at the moment.
    - **N** : the value `nil`.
    - **T** : the `this` reserved identifier.
    - **F** : the function at in dex `ix` in the module's function table.
    - **A** : the `args` reserved identifier.
* **POP** : pops a value from the stack, stores it in the variable identified by the string at index `ix` in the K table. If the variable does not already exist, it is created as a local variable.
* **ADD | SUB | MUL | DIV | MOD | AND | OR** : pops two values from the stack, performs the operation, and pushes the result on the stack.
* **NOT | UNM** : pops one value from the stack, performs the operation, and pushes the result on the stack.
* **EQ | LT | LTE | GT | GTE** : pops two values from the stack, compares them, and pushes the boolean result for the operation (the comparison returns 1 if greater, 0 if equal and -1 if lower).
* **TEST** : pops one value from the stack, tests its boolean representation, if it is `false`, jumps forward `ix` instructions.
* **JMP** : if the flag is `Jf`, jumps forward `ix` instructions, if it is `Jb`, jumps backward `ix + 1` instructions (because the `pc` is already pointing on the next instruction).
* **NEW** : creates a new object and pushes it on the stack. If `ix` is greater than 0, pops `2*ix` values from the stack, initializing fields on the object in `ix` pair of values representing the key and the value.
* **SFLD** : pops three values from the stack (`object`, `key` and `value` in order of pops) and sets the `object`'s `key` to `value`. It panics if `object` is not an object.
* **GFLD** : pops two values from the stack (`object` and `key` in order of pops) and pushes the value of the `object`'s `key` onto the stack. It panics if `object` is not an object.
* **CFLD** : pops two values from the stjumps forward `ix` instructions, if it is `Jb`, jumps backward `ix + 1` instructions (because the `pc` is already pointing on the next instruction).
* **NEW** : creates a new object and pushes it on the stack. If `ix` is greater than 0, pops `2*ix` values frolue at `object.key` is not a function, or if it does not exist and there is no `__noSuchMethod` function on the object.
* **CALL** : pops one value from the stack, and `ix` additional values representing the arguments, and calls the function, pushing the return value of the function on the stack. It panics if the expected function is not a function.
* **DUMP** : pretty-prints `ix` number of frames, starting at the current executing frame, to the execution context's `Stdout` stream. It is a no-op if the execution context is not in debug mode. This is the instruction generated by `debug` statements in the agora source code.

Next: [Back to the index](https://github.com/PuerkitoBio/agora/wiki)

