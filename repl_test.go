package main

import "testing"

func TestMultipleExpressions(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("*", makeBuiltin(builtinMul))

	input := `(define x 10) (define y 20) (+ x y)`
	exprs := readMultipleExprs(input)

	if len(exprs) != 3 {
		t.Fatalf("readMultipleExprs = %d expressions, want 3", len(exprs))
	}

	// Evaluate all
	for _, expr := range exprs {
		eval(expr, env)
	}

	// Check result
	result := eval(readStr("(+ x y)"), env)
	if result.Num != 30 {
		t.Errorf("x + y = %d, want 30", result.Num)
	}
}

func TestFactorialProgram(t *testing.T) {
	env := setupFullEnv()

	program := `
		(define factorial
			(lambda (n)
				(if (= n 0)
					1
					(* n (factorial (- n 1))))))

		(factorial 5)
	`

	exprs := readMultipleExprs(program)
	var result *Expr

	for _, expr := range exprs {
		result = eval(expr, env)
	}

	if result.Num != 120 {
		t.Errorf("factorial(5) = %d, want 120", result.Num)
	}
}

func TestMacroProgram(t *testing.T) {
	env := setupFullEnv()

	program := `
		(define defmacro
			(macro (name params body)
				(cons 'define (cons name (cons (cons 'macro (cons params (cons body nil))) nil)))))

		(defmacro when (test body)
			(cons 'if (cons test (cons body (cons 'nil nil)))))

		(when true 42)
	`

	exprs := readMultipleExprs(program)
	var result *Expr

	for _, expr := range exprs {
		result = eval(expr, env)
	}

	if result.Num != 42 {
		t.Errorf("when macro = %d, want 42", result.Num)
	}
}

func TestCompleteProgram(t *testing.T) {
	env := setupFullEnv()

	program := `
		; Define defmacro
		(define defmacro
			(macro (name params body)
				(cons 'define (cons name (cons (cons 'macro (cons params (cons body nil))) nil)))))

		; Define let macro
		(defmacro let (bindings body)
			(cons (cons 'lambda (cons (cons (head (head bindings)) nil) (cons body nil)))
				(cons (head (tail (head bindings))) nil)))

		; Use let to compute something
		(let ((x 10)) (+ (* x 2) 5))
	`

	exprs := readMultipleExprs(program)
	var result *Expr

	for _, expr := range exprs {
		result = eval(expr, env)
	}

	// let x=10: (+ (* 10 2) 5) = (+ 20 5) = 25
	if result.Num != 25 {
		t.Errorf("complete program = %d, want 25", result.Num)
	}
}

// Helper to setup full environment
func setupFullEnv() *Env {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("-", makeBuiltin(builtinSub))
	env.Define("*", makeBuiltin(builtinMul))
	env.Define("/", makeBuiltin(builtinDiv))
	env.Define("=", makeBuiltin(builtinEq))
	env.Define("<", makeBuiltin(builtinLt))
	env.Define("cons", makeBuiltin(builtinCons))
	env.Define("head", makeBuiltin(builtinHead))
	env.Define("tail", makeBuiltin(builtinTail))
	env.Define("null?", makeBuiltin(builtinNullP))
	env.Define("print", makeBuiltin(builtinPrint))
	return env
}
