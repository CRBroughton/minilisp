package main

import "testing"

func TestPrintNil(t *testing.T) {
	if got := printExpr(nilExpr); got != "nil" {
		t.Errorf("printExpr(nil) = %q, want \"nil\"", got)
	}
}

func TestPrintBool(t *testing.T) {
	if got := printExpr(trueExpr); got != "true" {
		t.Errorf("printExpr(true) = %q, want \"true\"", got)
	}
}

func TestPrintNumber(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{42, "42"},
		{-10, "-10"},
		{999, "999"},
	}

	for _, tt := range tests {
		got := printExpr(makeNum(tt.input))
		if got != tt.want {
			t.Errorf("printExpr(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestPrintSymbol(t *testing.T) {
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
		got := printExpr(makeSym(tt.input))
		if got != tt.want {
			t.Errorf("printExpr(sym %q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestPrintBuiltin(t *testing.T) {
	fn := makeBuiltin(func(args []*Expr) *Expr { return nilExpr })
	got := printExpr(fn)
	if got != "<builtin>" {
		t.Errorf("printExpr(builtin) = %q, want \"<builtin>\"", got)
	}
}

func TestPrintLambda(t *testing.T) {
	lambda := &Expr{Type: Lambda}
	got := printExpr(lambda)
	if got != "<lambda>" {
		t.Errorf("printExpr(lambda) = %q, want \"<lambda>\"", got)
	}
}

func TestPrintMacro(t *testing.T) {
	macro := &Expr{Type: Macro}
	got := printExpr(macro)
	if got != "<macro>" {
		t.Errorf("printExpr(macro) = %q, want \"<macro>\"", got)
	}
}

func TestPrintEmptyList(t *testing.T) {
	// Empty list is just nil
	got := printExpr(nilExpr)
	if got != "nil" {
		t.Errorf("printExpr(empty list) = %q, want \"nil\"", got)
	}
}

func TestPrintSimpleList(t *testing.T) {
	// (1 2 3)
	lst := list(makeNum(1), makeNum(2), makeNum(3))
	got := printExpr(lst)
	want := "(1 2 3)"
	if got != want {
		t.Errorf("printExpr((1 2 3)) = %q, want %q", got, want)
	}
}

func TestPrintNestedList(t *testing.T) {
	// (1 (2 3) 4)
	inner := list(makeNum(2), makeNum(3))
	lst := list(makeNum(1), inner, makeNum(4))
	got := printExpr(lst)
	want := "(1 (2 3) 4)"
	if got != want {
		t.Errorf("printExpr((1 (2 3) 4)) = %q, want %q", got, want)
	}
}

func TestPrintSymbolList(t *testing.T) {
	// (+ 1 2)
	lst := list(makeSym("+"), makeNum(1), makeNum(2))
	got := printExpr(lst)
	want := "(+ 1 2)"
	if got != want {
		t.Errorf("printExpr((+ 1 2)) = %q, want %q", got, want)
	}
}

func TestPrintImproperList(t *testing.T) {
	// (1 . 2) - improper list (dotted pair)
	pair := cons(makeNum(1), makeNum(2))
	got := printExpr(pair)
	want := "(1 . 2)"
	if got != want {
		t.Errorf("printExpr((1 . 2)) = %q, want %q", got, want)
	}
}
