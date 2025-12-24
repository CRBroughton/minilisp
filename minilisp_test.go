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
