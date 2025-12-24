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
		{"cons is Cons type", cons(makeNum(1)), Cons},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expr.Type != tt.want {
				t.Errorf("got type %v, want %v", tt.expr.Type, tt.want)
			}
		})
	}
}
