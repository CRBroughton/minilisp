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
