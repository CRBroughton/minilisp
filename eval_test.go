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
