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
}
