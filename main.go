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
	env.Define("-", makeBuiltin(builtinSub))
	env.Define("*", makeBuiltin(builtinMul))
	env.Define("=", makeBuiltin(builtinEq))

	result = eval(readStr("(+ (* 2 3) (* 4 5))"), env)
	fmt.Println(printExpr(result)) // "26"

	env.Define("<", makeBuiltin(builtinLt))
	result = eval(readStr("(< 3 10)"), env)
	fmt.Println(printExpr(result)) // "true"

	// Quote
	result = eval(readStr("'(+ 1 2)"), env)
	fmt.Println(printExpr(result)) // "(+ 1 2)" - not evaluated!

	// If
	result = eval(readStr("(if (< 3 5) 10 20)"), env)
	fmt.Println(printExpr(result)) // "10"

	// Define
	eval(readStr("(define x 42)"), env)
	result = eval(readStr("x"), env)
	fmt.Println(printExpr(result)) // "42"

	// Begin
	result = eval(readStr("(begin (define x 1) (define y 2) (+ x y))"), env)
	fmt.Println(printExpr(result)) // "3"

	// Anonymous functions
	result = eval(readStr("((lambda (x) (* x 2)) 21)"), env)
	fmt.Println(printExpr(result)) // 42

	// Named functions
	eval(readStr("(define double (lambda (x) (* x 2)))"), env)
	result = eval(readStr("(double 21)"), env)
	fmt.Println(printExpr(result)) // 42

	// Closures
	program := `
  (begin
    (define make-adder (lambda (n) (lambda (x) (+ x n))))
    (define add5 (make-adder 5))
    (add5 10))
`
	result = eval(readStr(program), env)
	fmt.Println(printExpr(result)) // 15

	// Recursion
	eval(readStr("(define factorial (lambda (n) (if (= n 0) 1 (* n (factorial (- n 1))))))"), env)
	result = eval(readStr("(factorial 5)"), env)
	fmt.Println(printExpr(result)) // 120
}
