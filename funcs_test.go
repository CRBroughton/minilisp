package main

import "testing"

func TestBuiltinAdd(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))

	tests := []struct {
		input string
		want  int
	}{
		{"(+ 1 2)", 3},
		{"(+ 1 2 3)", 6},
		{"(+ 1 2 3 4 5)", 15},
		{"(+ 0)", 0},
		{"(+ -5 5)", 0},
		{"(+ 100 -50)", 50},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)

		if result.Type != Number {
			t.Errorf("%s: type = %v, want Number", tt.input, result.Type)
		}
		if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestBuiltinSub(t *testing.T) {
	env := NewEnv(nil)
	env.Define("-", makeBuiltin(builtinSub))

	tests := []struct {
		input string
		want  int
	}{
		{"(- 10 5)", 5},
		{"(- 100 50 25)", 25},
		{"(- 0 10)", -10},
		{"(- 5 5)", 0},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)

		if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestBuiltinMul(t *testing.T) {
	env := NewEnv(nil)
	env.Define("*", makeBuiltin(builtinMul))

	tests := []struct {
		input string
		want  int
	}{
		{"(* 2 3)", 6},
		{"(* 2 3 4)", 24},
		{"(* 5)", 5},
		{"(* 10 0)", 0},
		{"(* -2 3)", -6},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)
		if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestBuiltinDiv(t *testing.T) {
	env := NewEnv(nil)
	env.Define("/", makeBuiltin(builtinDiv))

	tests := []struct {
		input string
		want  int
	}{
		{"(/ 10 2)", 5},
		{"(/ 100 10)", 10},
		{"(/ 7 2)", 3}, // Integer division
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)
		if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}

func TestBuiltinEq(t *testing.T) {
	env := NewEnv(nil)
	env.Define("=", makeBuiltin(builtinEq))

	tests := []struct {
		input    string
		wantTrue bool
	}{
		{"(= 5 5)", true},
		{"(= 5 6)", false},
		{"(= 0 0)", true},
		{"(= -5 -5)", true},
		{"(= 10 5)", false},
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

func TestBuiltinLt(t *testing.T) {
	env := NewEnv(nil)
	env.Define("<", makeBuiltin(builtinLt))

	tests := []struct {
		input    string
		wantTrue bool
	}{
		{"(< 3 5)", true},
		{"(< 5 3)", false},
		{"(< 5 5)", false},
		{"(< -10 0)", true},
		{"(< 0 -10)", false},
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

func TestBuiltinPairs(t *testing.T) {
	env := NewEnv(nil)
	env.Define("pair", makeBuiltin(builtinPair))

	expr := readStr("(pair 1 2)")
	result := eval(expr, env)

	if result.Type != Pair {
		t.Fatalf("pair result type = %v, want Pair", result.Type)
	}
	if result.Head.Num != 1 {
		t.Errorf("head = %d, want 1", result.Head.Num)
	}
	if result.Tail.Num != 2 {
		t.Errorf("tail = %d, want 2", result.Tail.Num)
	}
}

func TestBuiltinHead(t *testing.T) {
	env := NewEnv(nil)
	env.Define("head", makeBuiltin(builtinHead))
	env.Define("pair", makeBuiltin(builtinPair))

	expr := readStr("(head (pair 1 2))")
	result := eval(expr, env)

	if result.Num != 1 {
		t.Errorf("head = %d, want 1", result.Num)
	}
}

func TestBuiltinTail(t *testing.T) {
	env := NewEnv(nil)
	env.Define("tail", makeBuiltin(builtinTail))
	env.Define("pair", makeBuiltin(builtinPair))

	expr := readStr("(tail (pair 1 2))")
	result := eval(expr, env)

	if result.Num != 2 {
		t.Errorf("tail = %d, want 2", result.Num)
	}
}

func TestBuiltinNullP(t *testing.T) {
	env := NewEnv(nil)
	env.Define("null?", makeBuiltin(builtinNullP))

	tests := []struct {
		input    string
		wantTrue bool
	}{
		{"(null? nil)", true},
		{"(null? 42)", false},
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

func TestEvalList(t *testing.T) {
	env := NewEnv(nil)
	env.Define("x", makeNum(10))
	env.Define("y", makeNum(20))

	// Evaluate (x y) should give [10, 20]
	lst := list(makeSym("x"), makeSym("y"))
	results := evalList(lst, env)

	if len(results) != 2 {
		t.Fatalf("evalList length = %d, want 2", len(results))
	}
	if results[0].Num != 10 || results[1].Num != 20 {
		t.Errorf("evalList = [%d, %d], want [10, 20]", results[0].Num, results[1].Num)
	}
}

func TestComplexArithmetic(t *testing.T) {
	env := NewEnv(nil)
	env.Define("+", makeBuiltin(builtinAdd))
	env.Define("*", makeBuiltin(builtinMul))
	env.Define("-", makeBuiltin(builtinSub))

	tests := []struct {
		input string
		want  int
	}{
		{"(+ (* 2 3) (* 4 5))", 26}, // (+ 6 20) = 26
		{"(+ 1 (+ 2 (+ 3 4)))", 10}, // 1 + 2 + 3 + 4
		{"(- (+ 10 5) (* 2 3))", 9}, // 15 - 6
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		result := eval(expr, env)
		if result.Num != tt.want {
			t.Errorf("%s = %d, want %d", tt.input, result.Num, tt.want)
		}
	}
}
