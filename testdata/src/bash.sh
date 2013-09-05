#!/usr/bin/env agora run -R
// Output: Hello, Agora !
func greet(name) {
	fmt := import("fmt")
	fmt.Println("Hello,", name, "!")
}
greet("Agora")
