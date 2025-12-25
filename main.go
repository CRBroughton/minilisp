package main

import "fmt"

func main() {
	// Print simple values
	fmt.Println(printExpr(makeNum(42)))  // "42"
	fmt.Println(printExpr(makeSym("x"))) // "x"

	// Print lists
	lst := list(makeSym("+"), makeNum(1), makeNum(2))
	fmt.Println(printExpr(lst)) // "(+ 1 2)"

	// Print nested structures
	nested := list(makeNum(1), list(makeNum(2), makeNum(3)))
	fmt.Println(printExpr(nested))

	// Parse text
	expr := readStr("(+ 1 2)")
	// Print text
	fmt.Println(printExpr(expr))

	// again slight more complex version
	input := "(define x (lambda (n) (* n 2)))"
	expr = readStr(input)
	output := printExpr(expr)
	fmt.Println(output)

	// eval
	env := NewEnv(nil)
	// Numbers evaluate to themselves
	result := eval(makeNum(42), env)
	fmt.Println(printExpr(result)) // "42"

	// Define and lookup variables
	env.Define("x", makeNum(10))
	result = eval(makeSym("x"), env)
	fmt.Println(printExpr(result)) // "10"

	// Nested scoping works
	child := NewEnv(env)
	child.Define("y", makeNum(20))
	result = eval(makeSym("x"), child) // Finds x in parent
	fmt.Println(printExpr(result))     // "10"

	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("*", makeBuiltin(builtinMul))

	result = eval(readStr("(+ (* 2 3) (* 4 5))"), env)
	fmt.Println(printExpr(result)) // "26"

	env.Define("<", makeBuiltin(builtinLt))
	result = eval(readStr("(< 3 10)"), env)
	fmt.Println(printExpr(result)) // "true"
}
