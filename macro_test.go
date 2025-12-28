package main

import "testing"

func setupMacroTestEnv() *Env {
	env := NewEnv(nil)

	// Arithmetic
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("-", makeBuiltin(builtinSub))
	env.Define("*", makeBuiltin(builtinMul))
	env.Define("/", makeBuiltin(builtinDiv))
	env.Define("=", makeBuiltin(builtinEq))
	env.Define("<", makeBuiltin(builtinLt))
	env.Define(">", makeBuiltin(builtinGt))

	// String operations
	env.Define("string-append", makeBuiltin(builtinStringAppend))

	// List operations
	env.Define("pair", makeBuiltin(builtinPair))
	env.Define("list", makeBuiltin(builtinList))
	env.Define("head", makeBuiltin(builtinHead))
	env.Define("tail", makeBuiltin(builtinTail))

	// Hash operations
	env.Define("hash", makeBuiltin(builtinHash))
	env.Define("hash-get", makeBuiltin(builtinHashGet))
	env.Define("hash-set", makeBuiltin(builtinHashSet))

	// Bootstrap defmacro
	defmacroCode := "(define defmacro (macro (name params body) (pair 'define (pair name (pair (pair 'macro (pair params (pair body nil))) nil)))))"
	eval(readStr(defmacroCode), env)

	eval(readStr(`(load "std/macro.lisp")`), env)

	return env
}

func TestThreadFirstMacro(t *testing.T) {
	env := setupMacroTestEnv()

	// Test single form: (-> 5 (* 2)) should be 10
	code := `(-> 5 (* 2))`
	result := eval(readStr(code), env)
	if result.Num != 10 {
		t.Errorf("(-> 5 (* 2)) = %d, want 10", result.Num)
	}

	// Test multiple forms: (-> 5 (* 2) (+ 3)) should be 13
	// Expands to: (+ (* 5 2) 3) = (+ 10 3) = 13
	code = `(-> 5 (* 2) (+ 3))`
	result = eval(readStr(code), env)
	if result.Num != 13 {
		t.Errorf("(-> 5 (* 2) (+ 3)) = %d, want 13", result.Num)
	}

	// Test three forms: (-> 10 (- 3) (* 2)) should be 14
	// Expands to: (* (- 10 3) 2) = (* 7 2) = 14
	code = `(-> 10 (- 3) (* 2))`
	result = eval(readStr(code), env)
	if result.Num != 14 {
		t.Errorf("(-> 10 (- 3) (* 2)) = %d, want 14", result.Num)
	}
}

func TestThreadLastMacro(t *testing.T) {
	env := setupMacroTestEnv()

	// Test single form: (->> 5 (* 2)) should be 10
	code := `(->> 5 (* 2))`
	result := eval(readStr(code), env)
	if result.Num != 10 {
		t.Errorf("(->> 5 (* 2)) = %d, want 10", result.Num)
	}

	// Test difference from ->: (->> 3 (- 10)) should be 7
	// Expands to: (- 10 3) = 7
	code = `(->> 3 (- 10))`
	result = eval(readStr(code), env)
	if result.Num != 7 {
		t.Errorf("(->> 3 (- 10)) = %d, want 7", result.Num)
	}

	// Test string-append: (->> "World" (string-append "Hello, "))
	code = `(->> "World" (string-append "Hello, "))`
	result = eval(readStr(code), env)
	if result.Str != "Hello, World" {
		t.Errorf("(->> \"World\" (string-append \"Hello, \")) = %q, want \"Hello, World\"", result.Str)
	}
}

func TestWhenMacro(t *testing.T) {
	env := setupMacroTestEnv()

	// Test when with true - should execute body
	code := `(when true 42)`
	result := eval(readStr(code), env)
	if result.Num != 42 {
		t.Errorf("(when true 42) = %d, want 42", result.Num)
	}

	// Test when with nil - should return nil
	code = `(when nil 42)`
	result = eval(readStr(code), env)
	if result != nilExpr {
		t.Errorf("(when nil 42) should return nil")
	}

	// Test when with truthy value (non-zero number)
	code = `(when 1 100)`
	result = eval(readStr(code), env)
	if result.Num != 100 {
		t.Errorf("(when 1 100) = %d, want 100", result.Num)
	}
}

func TestCondMacro(t *testing.T) {
	env := setupMacroTestEnv()

	// Test cond with first clause true
	env.Define("y", makeNum(0))
	code := `(cond
		((= y 0) 100)
		((< y 0) 200)
		(true 300))`
	result := eval(readStr(code), env)
	if result.Num != 100 {
		t.Errorf("cond with y=0 should return 100, got %d", result.Num)
	}

	// Test cond with second clause true
	env.Define("z", makeNum(-5))
	code = `(cond
		((= z 0) 100)
		((< z 0) 200)
		(true 300))`
	result = eval(readStr(code), env)
	if result.Num != 200 {
		t.Errorf("cond with z=-5 should return 200, got %d", result.Num)
	}

	// Test cond with default clause
	env.Define("x", makeNum(5))
	code = `(cond
		((= x 0) 100)
		((< x 0) 200)
		(true 300))`
	result = eval(readStr(code), env)
	if result.Num != 300 {
		t.Errorf("cond with x=5 should return 300, got %d", result.Num)
	}
}

func TestThreadMacrosWithHash(t *testing.T) {
	env := setupMacroTestEnv()

	// Test -> with hash-get (from fetch.lisp pattern)
	code := `(-> (hash "name" "Alice" "age" 30) (hash-get "name"))`
	result := eval(readStr(code), env)
	if result.Str != "Alice" {
		t.Errorf("(-> hash (hash-get \"name\")) = %q, want \"Alice\"", result.Str)
	}

	// Test chaining multiple hash-gets
	code = `(define user (hash "profile" (hash "name" "Bob")))`
	eval(readStr(code), env)
	code = `(-> user (hash-get "profile") (hash-get "name"))`
	result = eval(readStr(code), env)
	if result.Str != "Bob" {
		t.Errorf("chained hash-get = %q, want \"Bob\"", result.Str)
	}
}

func TestThreadMacrosWithList(t *testing.T) {
	env := setupMacroTestEnv()

	// Test -> with list operations
	code := `(-> (list 1 2 3) (head))`
	result := eval(readStr(code), env)
	if result.Num != 1 {
		t.Errorf("(-> (list 1 2 3) (head)) = %d, want 1", result.Num)
	}

	// Test ->> with pair construction
	// (->> 3 (pair 2) (pair 1)) expands to (pair 1 (pair 2 3))
	// This creates (1 . (2 . 3)) which is (1 2 . 3) - not a proper list
	code = `(->> 3 (pair 2) (pair 1))`
	result = eval(readStr(code), env)
	if result.Head.Num != 1 {
		t.Errorf("head should be 1, got %d", result.Head.Num)
	}
	if result.Tail.Head.Num != 2 {
		t.Errorf("second element should be 2, got %d", result.Tail.Head.Num)
	}
	if result.Tail.Tail.Num != 3 {
		t.Errorf("tail should be 3, got %d", result.Tail.Tail.Num)
	}
}
