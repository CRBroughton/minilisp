package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestEvalNumber(t *testing.T) {
	env := NewEnv(nil)
	tests := []int{0, 42, -10, 999}

	for _, num := range tests {
		expr := makeNum(num)
		result := eval(expr, env)

		if result.Type != Number {
			t.Errorf("eval(%d) type = %v, want Number", num, result.Type)
		}
		if result.Num != num {
			t.Errorf("eval(%d) = %d, want %d", num, result.Num, num)
		}
	}
}

func TestEvalBool(t *testing.T) {
	env := NewEnv(nil)

	result := eval(trueExpr, env)
	if result != trueExpr {
		t.Error("eval(true) should return trueExpr")
	}
}

func TestEvalNil(t *testing.T) {
	env := NewEnv(nil)

	result := eval(nilExpr, env)
	if result != nilExpr {
		t.Error("eval(nil) should return nilExpr")
	}
}

func TestEvalSymbol(t *testing.T) {
	env := NewEnv(nil)
	env.Define("x", makeNum(42))
	env.Define("y", makeNum(99))

	tests := []struct {
		symbol string
		want   int
	}{
		{"x", 42},
		{"y", 99},
	}

	for _, tt := range tests {
		expr := makeSym(tt.symbol)
		result := eval(expr, env)

		if result.Type != Number {
			t.Errorf("eval(%s) type = %v, want Number", tt.symbol, result.Type)
		}
		if result.Num != tt.want {
			t.Errorf("eval(%s) = %d, want %d", tt.symbol, result.Num, tt.want)
		}
	}
}

func TestEvalUndefinedSymbol(t *testing.T) {
	env := NewEnv(nil)

	defer func() {
		if r := recover(); r == nil {
			t.Error("eval(undefined) should panic")
		}
	}()

	eval(makeSym("undefined"), env)
}

func TestEvalSymbolInNestedScope(t *testing.T) {
	parent := NewEnv(nil)
	parent.Define("x", makeNum(10))

	child := NewEnv(parent)
	child.Define("y", makeNum(20))

	// Child should see both x and y
	xResult := eval(makeSym("x"), child)
	if xResult.Num != 10 {
		t.Errorf("eval(x) in child = %d, want 10", xResult.Num)
	}

	yResult := eval(makeSym("y"), child)
	if yResult.Num != 20 {
		t.Errorf("eval(y) in child = %d, want 20", yResult.Num)
	}
}

func TestEvalWithShadowing(t *testing.T) {
	parent := NewEnv(nil)
	parent.Define("x", makeNum(10))

	child := NewEnv(parent)
	child.Define("x", makeNum(99)) // Shadow parent's x

	result := eval(makeSym("x"), child)
	if result.Num != 99 {
		t.Errorf("eval(x) in child = %d, want 99 (shadowed value)", result.Num)
	}
}

func TestQuote(t *testing.T) {
	env := NewEnv(nil)

	tests := []struct {
		input string
		want  string
	}{
		{"(quote x)", "x"},
		{"(quote 42)", "42"},
		{"(quote (+ 1 2))", "(+ 1 2)"},
		{"(quote (quote x))", "(quote x)"},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)
		got := printExpr(result)

		if got != tt.want {
			t.Errorf("%s = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestQuoteSugar(t *testing.T) {
	env := NewEnv(nil)

	tests := []struct {
		input string
		want  string
	}{
		{"'x", "x"},
		{"'42", "42"},
		{"'(+ 1 2)", "(+ 1 2)"},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)
		got := printExpr(result)

		if got != tt.want {
			t.Errorf("%s = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIf(t *testing.T) {
	env := NewEnv(nil)

	tests := []struct {
		input string
		want  int
	}{
		{"(if true 1 2)", 1},
		{"(if nil 1 2)", 2},
		{"(if 42 10 20)", 10}, // Any non-nil is truthy
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)

		if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestIfOnlyEvaluatesOneBranch(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))

	// Only the true branch should evaluate
	expr := readStr("(if true (+ 1 2) (+ 3 undefined))")
	result := eval(expr, env)

	if result.Num != 3 {
		t.Errorf("if true branch = %d, want 3", result.Num)
	}

	// Only the false branch should evaluate
	expr = readStr("(if nil (+ 1 undefined) (+ 2 3))")
	result = eval(expr, env)

	if result.Num != 5 {
		t.Errorf("if false branch = %d, want 5", result.Num)
	}
}

func TestDefine(t *testing.T) {
	env := NewEnv(nil)

	// Define a variable
	expr := readStr("(define x 42)")
	eval(expr, env)

	// Look it up
	val, ok := env.Lookup("x")
	if !ok {
		t.Fatal("x should be defined")
	}
	if val.Num != 42 {
		t.Errorf("x = %d, want 42", val.Num)
	}
}

func TestDefineWithExpression(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("*", makeBuiltin(builtinMul))

	// Define using an expression
	expr := readStr("(define result (+ (* 2 3) (* 4 5)))")
	eval(expr, env)

	val, _ := env.Lookup("result")
	if val.Num != 26 {
		t.Errorf("result = %d, want 26", val.Num)
	}
}

func TestBegin(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))

	// begin evaluates multiple expressions, returns last
	expr := readStr("(begin (define x 10) (define y 20) (+ x y))")
	result := eval(expr, env)

	if result.Num != 30 {
		t.Errorf("begin result = %d, want 30", result.Num)
	}

	// x and y should be defined
	if val, _ := env.Lookup("x"); val.Num != 10 {
		t.Error("x should be 10")
	}
	if val, _ := env.Lookup("y"); val.Num != 20 {
		t.Error("y should be 20")
	}
}

func TestNestedIf(t *testing.T) {
	env := NewEnv(nil)
	env.Define("<", makeBuiltin(builtinLt))

	// Nested if: (if (< 3 5) (if true 1 2) 3)
	expr := readStr("(if (< 3 5) (if true 1 2) 3)")
	result := eval(expr, env)

	if result.Num != 1 {
		t.Errorf("nested if = %d, want 1", result.Num)
	}
}

func TestComplexProgram(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("*", makeBuiltin(builtinMul))
	env.Define("<", makeBuiltin(builtinLt))

	program := `
		(begin
			(define x 10)
			(define y (+ x 5))
			(if (< x y)
				(* x y)
				0))
	`

	expr := readStr(program)
	result := eval(expr, env)

	// x=10, y=15, x<y is true, so (* 10 15) = 150
	if result.Num != 150 {
		t.Errorf("program result = %d, want 150", result.Num)
	}
}

func TestLambda(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("*", makeBuiltin(builtinMul))

	// Create a lambda
	expr := readStr("(lambda (x) (* x 2))")
	result := eval(expr, env)

	if result.Type != Lambda {
		t.Fatalf("lambda type = %v, want Lambda", result.Type)
	}
}

func TestLambdaApplication(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("*", makeBuiltin(builtinMul))

	// ((lambda (x) (* x 2)) 21) → 42
	expr := readStr("((lambda (x) (* x 2)) 21)")
	result := eval(expr, env)

	if result.Num != 42 {
		t.Errorf("lambda application = %d, want 42", result.Num)
	}
}

func TestLambdaMultipleParams(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))

	// ((lambda (x y) (+ x y)) 10 20) → 30
	expr := readStr("((lambda (x y) (+ x y)) 10 20)")
	result := eval(expr, env)

	if result.Num != 30 {
		t.Errorf("lambda with 2 params = %d, want 30", result.Num)
	}
}

func TestDefineLambda(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("*", makeBuiltin(builtinMul))

	// Define a function
	eval(readStr("(define double (lambda (x) (* x 2)))"), env)

	// Use it
	result := eval(readStr("(double 21)"), env)

	if result.Num != 42 {
		t.Errorf("double(21) = %d, want 42", result.Num)
	}
}

func TestSimpleClosure(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))

	program := `
		(begin
			(define x 10)
			(define add-x (lambda (y) (+ x y)))
			(add-x 5))
	`

	result := eval(readStr(program), env)

	if result.Num != 15 {
		t.Errorf("closure result = %d, want 15", result.Num)
	}
}

func TestNestedClosure(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))

	program := `
		(begin
			(define make-adder
				(lambda (n)
					(lambda (x) (+ x n))))
			(define add5 (make-adder 5))
			(add5 10))
	`

	result := eval(readStr(program), env)

	if result.Num != 15 {
		t.Errorf("nested closure = %d, want 15", result.Num)
	}
}

func TestClosureCapture(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))

	// Two closures capturing different values
	program := `
		(begin
			(define make-adder (lambda (n) (lambda (x) (+ x n))))
			(define add3 (make-adder 3))
			(define add7 (make-adder 7))
			(+ (add3 10) (add7 10)))
	`

	result := eval(readStr(program), env)

	// (add3 10) = 13, (add7 10) = 17, 13 + 17 = 30
	if result.Num != 30 {
		t.Errorf("two closures = %d, want 30", result.Num)
	}
}

func TestRecursion(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("-", makeBuiltin(builtinSub))
	env.Define("*", makeBuiltin(builtinMul))
	env.Define("=", makeBuiltin(builtinEq))

	// Define factorial
	factorial := `
		(define factorial
			(lambda (n)
				(if (= n 0)
					1
					(* n (factorial (- n 1))))))
	`
	eval(readStr(factorial), env)

	tests := []struct {
		input int
		want  int
	}{
		{0, 1},
		{1, 1},
		{5, 120},
		{6, 720},
	}

	for _, tt := range tests {
		expr := readStr(fmt.Sprintf("(factorial %d)", tt.input))
		result := eval(expr, env)

		if result.Num != tt.want {
			t.Errorf("factorial(%d) = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestFibonacci(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("-", makeBuiltin(builtinSub))
	env.Define("<", makeBuiltin(builtinLt))

	fib := `
		(define fib
			(lambda (n)
				(if (< n 2)
					n
					(+ (fib (- n 1)) (fib (- n 2))))))
	`
	eval(readStr(fib), env)

	tests := []struct {
		input int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
		{7, 13},
	}

	for _, tt := range tests {
		expr := readStr(fmt.Sprintf("(fib %d)", tt.input))
		result := eval(expr, env)

		if result.Num != tt.want {
			t.Errorf("fib(%d) = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestHigherOrderFunction(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("*", makeBuiltin(builtinMul))

	// Function that takes a function as argument
	program := `
		(begin
			(define apply-twice
				(lambda (f x)
					(f (f x))))
			(define double (lambda (n) (* n 2)))
			(apply-twice double 5))
	`

	result := eval(readStr(program), env)

	// double(double(5)) = double(10) = 20
	if result.Num != 20 {
		t.Errorf("higher-order function = %d, want 20", result.Num)
	}
}

func TestMacroCreation(t *testing.T) {
	env := NewEnv(nil)

	// Create a macro
	expr := readStr("(macro (x) x)")
	result := eval(expr, env)

	if result.Type != Macro {
		t.Fatalf("macro type = %v, want Macro", result.Type)
	}
}

func TestSimpleMacro(t *testing.T) {
	env := NewEnv(nil)
	env.Define("pair", makeBuiltin(builtinPair))

	// Define a macro that quotes its argument
	eval(readStr("(define my-quote (macro (x) (pair 'quote (pair x nil))))"), env)

	// Use it
	result := eval(readStr("(my-quote (+ 1 2))"), env)

	// Should return (+ 1 2) unevaluated
	got := printExpr(result)
	want := "(+ 1 2)"

	if got != want {
		t.Errorf("my-quote = %q, want %q", got, want)
	}
}

func TestMacroVsFunction(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("pair", makeBuiltin(builtinPair))

	// Function version - evaluates argument
	eval(readStr("(define func-quote (lambda (x) x))"), env)
	funcResult := eval(readStr("(func-quote (+ 1 2))"), env)

	// Macro version - doesn't evaluate argument
	eval(readStr("(define macro-quote (macro (x) (pair 'quote (pair x nil))))"), env)
	macroResult := eval(readStr("(macro-quote (+ 1 2))"), env)

	// Function should return 3
	if funcResult.Num != 3 {
		t.Errorf("function = %v, want 3", funcResult)
	}

	// Macro should return (+ 1 2)
	if printExpr(macroResult) != "(+ 1 2)" {
		t.Errorf("macro = %q, want \"(+ 1 2)\"", printExpr(macroResult))
	}
}

func TestUnlessMacro(t *testing.T) {
	env := NewEnv(nil)
	env.Define("pair", makeBuiltin(builtinPair))

	// Define unless macro
	unless := `
		(define unless
			(macro (test body)
				(pair 'if (pair test (pair 'nil (pair body nil))))))
	`
	eval(readStr(unless), env)

	// Test it
	tests := []struct {
		input string
		want  int
	}{
		{"(unless nil 42)", 42},
		{"(unless true 42)", 0}, // Returns nil → 0
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)

		// Handle nil case
		if result == nilExpr {
			if tt.want != 0 {
				t.Errorf("%s = nil, want %d", tt.input, tt.want)
			}
		} else if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestMacroExpansion(t *testing.T) {
	env := NewEnv(nil)
	env.Define("pair", makeBuiltin(builtinPair))

	// Define unless macro
	unless := `
		(define unless
			(macro (test body)
				(pair 'if (pair test (pair 'nil (pair body nil))))))
	`
	eval(readStr(unless), env)

	// Use it
	result := eval(readStr("(unless nil 42)"), env)

	if result.Num != 42 {
		t.Errorf("unless nil 42 = %v, want 42", result)
	}
}

func TestDefmacro(t *testing.T) {
	env := NewEnv(nil)
	env.Define("pair", makeBuiltin(builtinPair))

	// Bootstrap defmacro
	defmacroCode := "(define defmacro (macro (name params body) (pair 'define (pair name (pair (pair 'macro (pair params (pair body nil))) nil)))))"
	eval(readStr(defmacroCode), env)

	// Use defmacro to define unless
	unless := "(defmacro unless (test body) (pair 'if (pair test (pair 'nil (pair body nil)))))"
	eval(readStr(unless), env)

	// Test it
	result := eval(readStr("(unless nil 99)"), env)

	if result.Num != 99 {
		t.Errorf("defmacro unless = %v, want 99", result)
	}
}

func TestAndMacro(t *testing.T) {
	env := NewEnv(nil)
	env.Define("pair", makeBuiltin(builtinPair))

	defmacro := "(define defmacro (macro (name params body) (pair 'define (pair name (pair (pair 'macro (pair params (pair body nil))) nil)))))"
	eval(readStr(defmacro), env)

	andMacro := "(defmacro and (a b) (pair 'if (pair a (pair b (pair 'nil nil)))))"
	eval(readStr(andMacro), env)

	tests := []struct {
		input    string
		wantTrue bool
	}{
		{"(and true true)", true},
		{"(and true nil)", false},
		{"(and nil true)", false},
		{"(and nil nil)", false},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)

		isTrue := result == trueExpr
		if isTrue != tt.wantTrue {
			t.Errorf("%s = %v, want %v", tt.input, isTrue, tt.wantTrue)
		}
	}
}

func TestOrMacro(t *testing.T) {
	env := NewEnv(nil)
	env.Define("pair", makeBuiltin(builtinPair))

	defmacro := "(define defmacro (macro (name params body) (pair 'define (pair name (pair (pair 'macro (pair params (pair body nil))) nil)))))"
	eval(readStr(defmacro), env)

	orMacro := "(defmacro or (a b) (pair 'if (pair a (pair a (pair b nil)))))"
	eval(readStr(orMacro), env)

	tests := []struct {
		input    string
		wantTrue bool
	}{
		{"(or true true)", true},
		{"(or true nil)", true},
		{"(or nil true)", true},
		{"(or nil nil)", false},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)

		isTrue := result == trueExpr
		if isTrue != tt.wantTrue {
			t.Errorf("%s = %v, want %v", tt.input, isTrue, tt.wantTrue)
		}
	}
}

func TestLetMacro(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("head", makeBuiltin(builtinHead))
	env.Define("tail", makeBuiltin(builtinTail))
	env.Define("pair", makeBuiltin(builtinPair))

	defmacro := "(define defmacro (macro (name params body) (pair 'define (pair name (pair (pair 'macro (pair params (pair body nil))) nil)))))"
	eval(readStr(defmacro), env)

	// let macro (simplified - one binding only)
	letMacro := `
		(defmacro let (bindings body)
			(pair (pair 'lambda
					(pair (pair (head (head bindings)) nil)
						(pair body nil)))
				(pair (head (tail (head bindings))) nil)))
	`
	eval(readStr(letMacro), env)

	// Test it
	result := eval(readStr("(let ((x 10)) (+ x 5))"), env)

	if result.Num != 15 {
		t.Errorf("let = %v, want 15", result)
	}
}

func TestRecursiveMacroExpansion(t *testing.T) {
	env := NewEnv(nil)
	env.Define("pair", makeBuiltin(builtinPair))

	defmacro := "(define defmacro (macro (name params body) (pair 'define (pair name (pair (pair 'macro (pair params (pair body nil))) nil)))))"
	eval(readStr(defmacro), env)

	// when expands to if
	when := "(defmacro when (test body) (pair 'if (pair test (pair body (pair 'nil nil)))))"
	eval(readStr(when), env)

	// unless expands to when (which expands to if)
	unless := "(defmacro unless (test body) (pair 'when (pair (pair 'if (pair test (pair 'nil (pair true nil)))) (pair body nil))))"
	eval(readStr(unless), env)

	result := eval(readStr("(unless nil 77)"), env)

	if result.Num != 77 {
		t.Errorf("recursive expansion = %v, want 77", result)
	}
}

func TestLoadSimpleFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test-*.lisp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write test code
	content := `(define x 42)
(define y 99)`
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Load the file
	env := setupFullEnv()

	loadExpr := readStr(fmt.Sprintf("(load \"%s\")", tmpfile.Name()))
	eval(loadExpr, env)

	// Check that variables are defined
	if val, ok := env.Lookup("x"); !ok || val.Num != 42 {
		t.Errorf("x should be 42, got %v", val)
	}
	if val, ok := env.Lookup("y"); !ok || val.Num != 99 {
		t.Errorf("y should be 99, got %v", val)
	}
}

func TestLoadWithMacros(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "macros-*.lisp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Define a macro in the file
	content := `(define defmacro
		(macro (name params body)
			(pair 'define (pair name (pair (pair 'macro (pair params (pair body nil))) nil)))))

(defmacro when (test body)
	(pair 'if (pair test (pair body (pair 'nil nil)))))`

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Load and use the macro
	env := setupFullEnv()

	loadExpr := readStr(fmt.Sprintf("(load \"%s\")", tmpfile.Name()))
	eval(loadExpr, env)

	// Use the when macro
	result := eval(readStr("(when true 42)"), env)
	if result.Num != 42 {
		t.Errorf("when macro should return 42, got %v", result)
	}
}

func TestLoadReturnsLastValue(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "return-*.lisp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := `(define x 10)
(define y 20)
(+ x y)`

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	env := setupFullEnv()

	loadExpr := readStr(fmt.Sprintf("(load \"%s\")", tmpfile.Name()))
	result := eval(loadExpr, env)

	// Should return 30 (the result of (+ x y))
	if result.Num != 30 {
		t.Errorf("load should return 30, got %v", result)
	}
}

func TestLoadNonexistentFile(t *testing.T) {
	env := setupFullEnv()

	defer func() {
		if r := recover(); r == nil {
			t.Error("loading nonexistent file should panic")
		}
	}()

	loadExpr := readStr("(load \"nonexistent.lisp\")")
	eval(loadExpr, env)
}

func TestLoadRelativePath(t *testing.T) {
	// Create subdirectory
	tmpdir, err := os.MkdirTemp("", "testdir-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// Create file in subdirectory
	filepath := filepath.Join(tmpdir, "helper.lisp")
	content := `(define helper-value 123)`

	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	env := setupFullEnv()

	// Change to temp directory
	oldDir, _ := os.Getwd()
	os.Chdir(tmpdir)
	defer os.Chdir(oldDir)

	// Load with relative path
	loadExpr := readStr("(load \"helper.lisp\")")
	eval(loadExpr, env)

	if val, ok := env.Lookup("helper-value"); !ok || val.Num != 123 {
		t.Errorf("helper-value should be 123, got %v", val)
	}
}

func TestLoadPreventsDuplicates(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "counter-*.lisp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// File that increments a counter
	content := `(define counter (if (null? counter) 1 (+ counter 1)))`

	if err := os.WriteFile(tmpfile.Name(), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	env := setupFullEnv()
	env.Define("counter", nilExpr) // Start with nil

	// Load twice
	loadExpr := readStr(fmt.Sprintf("(load \"%s\")", tmpfile.Name()))
	eval(loadExpr, env)
	eval(loadExpr, env)

	// counter should be 2 (loaded twice)
	if val, ok := env.Lookup("counter"); !ok || val.Num != 2 {
		t.Errorf("counter should be 2 (file loaded twice), got %v", val)
	}
}
