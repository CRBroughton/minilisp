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

func TestBuiltinCons(t *testing.T) {
	env := NewEnv(nil)
	env.Define("cons", makeBuiltin(builtinCons))

	expr := readStr("(cons 1 2)")
	result := eval(expr, env)

	if result.Type != Cons {
		t.Fatalf("cons result type = %v, want Cons", result.Type)
	}
	if result.Car.Num != 1 {
		t.Errorf("car = %d, want 1", result.Car.Num)
	}
	if result.Cdr.Num != 2 {
		t.Errorf("cdr = %d, want 2", result.Cdr.Num)
	}
}

func TestBuiltinCar(t *testing.T) {
	env := NewEnv(nil)
	env.Define("car", makeBuiltin(builtinCar))
	env.Define("cons", makeBuiltin(builtinCons))

	expr := readStr("(car (cons 1 2))")
	result := eval(expr, env)

	if result.Num != 1 {
		t.Errorf("car = %d, want 1", result.Num)
	}
}

func TestBuiltinCdr(t *testing.T) {
	env := NewEnv(nil)
	env.Define("cdr", makeBuiltin(builtinCdr))
	env.Define("cons", makeBuiltin(builtinCons))

	expr := readStr("(cdr (cons 1 2))")
	result := eval(expr, env)

	if result.Num != 2 {
		t.Errorf("cdr = %d, want 2", result.Num)
	}
}
