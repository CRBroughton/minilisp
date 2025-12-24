package minilisp

import "testing"

func TestExprTypes(t *testing.T) {
	tests := []struct {
		name string
		expr *Expr
		want ExprType
	}{
		{"nil is Nil type", nilExpr, Nil},
		{"true is Bool", trueExpr, Bool},
		{"number is Number type", makeNum(42), Number},
		{"symbol is Symbol type", makeSym("x"), Symbol},
		{"cons is Cons type", cons(makeNum(1), nilExpr), Cons},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expr.Type != tt.want {
				t.Errorf("got type %v, want %v", tt.expr.Type, tt.want)
			}
		})
	}
}

func TestMakeNum(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{0, 0},
		{42, 42},
		{-10, -10},
		{999, 999},
	}

	for _, tt := range tests {
		expr := makeNum(tt.input)
		if expr.Type != Number {
			t.Errorf("makeNum(%d) type = %v, want Number", tt.input, expr.Type)
		}

		if expr.Num != tt.want {
			t.Errorf("makeNum(%d) = %d, want %d", tt.input, expr.Num, tt.want)
		}
	}
}

func TestMakeSym(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"x", "x"},
		{"+", "+"},
		{"factorial", "factorial"},
		{"my-var", "my-var"},
	}

	for _, tt := range tests {
		expr := makeSym(tt.input)
		if expr.Type != Symbol {
			t.Errorf("makeSym(%q) type = %v, want Symbol", tt.input, expr.Type)
		}
		if expr.Sym != tt.want {
			t.Errorf("makeSym(%q) = %q, want %q", tt.input, expr.Sym, tt.want)
		}
	}
}

func TestMakeSymSpecialCases(t *testing.T) {
	tests := []struct {
		input    string
		wantType ExprType
	}{
		{"nil", Nil},
		{"true", Bool},
		{"false", Bool},
	}

	for _, tt := range tests {
		expr := makeSym(tt.input)
		if expr.Type != tt.wantType {
			t.Errorf("makeSym(%q) type = %v, want %v", tt.input, expr.Type, tt.wantType)
		}
	}
}
