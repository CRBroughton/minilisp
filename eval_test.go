package main

import "testing"

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
		t.Error("eval(#t) should return trueExpr")
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

	// Nested if: (if (< 3 5) (if #t 1 2) 3)
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
