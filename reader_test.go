package main

import "testing"

func TestReadNumber(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"0", 0},
		{"42", 42},
		{"-10", -10},
		{"999", 999},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		if expr.Type != Number {
			t.Errorf("readStr(%q) type = %v, want Number", tt.input, expr.Type)
		}
		if expr.Num != tt.want {
			t.Errorf("readStr(%q) = %d, want %d", tt.input, expr.Num, tt.want)
		}
	}
}

func TestReadSymbol(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"x", "x"},
		{"+", "+"},
		{"factorial", "factorial"},
		{"my-var", "my-var"},
		{"*global*", "*global*"},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		if expr.Type != Symbol {
			t.Errorf("readStr(%q) type = %v, want Symbol", tt.input, expr.Type)
		}
		if expr.Sym != tt.want {
			t.Errorf("readStr(%q) = %q, want %q", tt.input, expr.Sym, tt.want)
		}
	}
}

func TestReadNil(t *testing.T) {
	expr := readStr("nil")
	if expr != nilExpr {
		t.Error("readStr(\"nil\") should return nilExpr")
	}
}

func TestReadBool(t *testing.T) {
	tests := []string{"true"}
	for _, input := range tests {
		expr := readStr(input)
		if expr != trueExpr {
			t.Errorf("readStr(%q) should return trueExpr", input)
		}
	}
}

func TestReadEmptyList(t *testing.T) {
	expr := readStr("()")
	if expr != nilExpr {
		t.Error("readStr(\"()\") should return nilExpr")
	}
}

func TestReadSimpleList(t *testing.T) {
	// (1 2 3)
	expr := readStr("(1 2 3)")

	// Check structure
	if expr.Type != Pair {
		t.Fatalf("type = %v, want Pair", expr.Type)
	}

	// Check elements
	nums := listToSlice(expr)
	if len(nums) != 3 {
		t.Fatalf("length = %d, want 3", len(nums))
	}

	for i, want := range []int{1, 2, 3} {
		if nums[i].Num != want {
			t.Errorf("element %d = %d, want %d", i, nums[i].Num, want)
		}
	}
}

func TestReadNestedList(t *testing.T) {
	// (1 (2 3) 4)
	expr := readStr("(1 (2 3) 4)")

	elems := listToSlice(expr)
	if len(elems) != 3 {
		t.Fatalf("length = %d, want 3", len(elems))
	}

	// First element should be 1
	if elems[0].Num != 1 {
		t.Errorf("first = %d, want 1", elems[0].Num)
	}

	// Second element should be (2 3)
	if elems[1].Type != Pair {
		t.Fatalf("second type = %v, want Pair", elems[1].Type)
	}
	inner := listToSlice(elems[1])
	if len(inner) != 2 || inner[0].Num != 2 || inner[1].Num != 3 {
		t.Errorf("inner list incorrect")
	}

	// Third element should be 4
	if elems[2].Num != 4 {
		t.Errorf("third = %d, want 4", elems[2].Num)
	}
}

func TestReadQuote(t *testing.T) {
	// 'x should become (quote x)
	expr := readStr("'x")

	if expr.Type != Pair {
		t.Fatalf("type = %v, want Pair", expr.Type)
	}

	elems := listToSlice(expr)
	if len(elems) != 2 {
		t.Fatalf("length = %d, want 2", len(elems))
	}

	// First element should be 'quote
	if elems[0].Type != Symbol || elems[0].Sym != "quote" {
		t.Errorf("first = %v, want 'quote", elems[0])
	}

	// Second element should be 'x
	if elems[1].Type != Symbol || elems[1].Sym != "x" {
		t.Errorf("second = %v, want 'x", elems[1])
	}
}

func TestReadQuoteList(t *testing.T) {
	// '(1 2 3) should become (quote (1 2 3))
	expr := readStr("'(1 2 3)")

	elems := listToSlice(expr)
	if len(elems) != 2 {
		t.Fatalf("length = %d, want 2", len(elems))
	}

	if elems[0].Sym != "quote" {
		t.Errorf("first = %v, want 'quote", elems[0])
	}

	quoted := listToSlice(elems[1])
	if len(quoted) != 3 {
		t.Fatalf("quoted list length = %d, want 3", len(quoted))
	}
}

func TestReadWithWhitespace(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"  42  ", "42"},
		{"\n\n100\n\n", "100"},
		{"\t(\t1\t2\t)\t", "(1 2)"},
		{"(  +   1   2  )", "(+ 1 2)"},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		got := printExpr(expr)
		if got != tt.want {
			t.Errorf("readStr(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestReadComments(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"; comment\n42", "42"},
		{"42 ; end of line comment", "42"},
		{"(1 ; middle comment\n 2)", "(1 2)"},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		got := printExpr(expr)
		if got != tt.want {
			t.Errorf("readStr(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestReadComplexExpression(t *testing.T) {
	input := "(define factorial (lambda (n) (if (= n 0) 1 (* n (factorial (- n 1))))))"
	expr := readStr(input)

	// Just verify it parses without error
	if expr.Type != Pair {
		t.Fatalf("type = %v, want Pair", expr.Type)
	}

	// Verify it round-trips correctly
	printed := printExpr(expr)
	reparsed := readStr(printed)
	printed2 := printExpr(reparsed)

	if printed != printed2 {
		t.Errorf("round-trip failed:\noriginal:  %s\nreparsed: %s", printed, printed2)
	}
}

func TestReadString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`"hello"`, "hello"},
		{`"world"`, "world"},
		{`""`, ""},
		{`"hello world"`, "hello world"},
		{`"std.lisp"`, "std.lisp"},
		{`"/path/to/file.lisp"`, "/path/to/file.lisp"},
		{`"with-dashes"`, "with-dashes"},
		{`"with_underscores"`, "with_underscores"},
		{`"123"`, "123"},
		{`"symbols: + - * /"`, "symbols: + - * /"},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)
		if expr.Type != String {
			t.Errorf("readStr(%q) type = %v, want String", tt.input, expr.Type)
		}
		if expr.Str != tt.want {
			t.Errorf("readStr(%q) = %q, want %q", tt.input, expr.Str, tt.want)
		}
	}
}

func TestReadStringInList(t *testing.T) {
	// (load "file.lisp")
	expr := readStr(`(load "file.lisp")`)

	if expr.Type != Pair {
		t.Fatalf("type = %v, want Pair", expr.Type)
	}

	elems := listToSlice(expr)
	if len(elems) != 2 {
		t.Fatalf("length = %d, want 2", len(elems))
	}

	// First element: load (symbol)
	if elems[0].Type != Symbol || elems[0].Sym != "load" {
		t.Errorf("first element = %v (type %v), want symbol 'load'", elems[0].Sym, elems[0].Type)
	}

	// Second element: "file.lisp" (string)
	if elems[1].Type != String || elems[1].Str != "file.lisp" {
		t.Errorf("second element = %q (type %v), want string \"file.lisp\"", elems[1].Str, elems[1].Type)
	}
}

func TestReadMultipleStrings(t *testing.T) {
	// (print "hello" "world")
	expr := readStr(`(print "hello" "world")`)

	elems := listToSlice(expr)
	if len(elems) != 3 {
		t.Fatalf("length = %d, want 3", len(elems))
	}

	if elems[0].Sym != "print" {
		t.Errorf("first = %v, want 'print'", elems[0])
	}

	if elems[1].Type != String || elems[1].Str != "hello" {
		t.Errorf("second = %q (type %v), want \"hello\"", elems[1].Str, elems[1].Type)
	}

	if elems[2].Type != String || elems[2].Str != "world" {
		t.Errorf("third = %q (type %v), want \"world\"", elems[2].Str, elems[2].Type)
	}
}

func TestReadStringWithWhitespace(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`  "hello"  `, "hello"},
		{"\n\"test\"\n", "test"},
		{`(  "spaced"  )`, "spaced"},
	}

	for _, tt := range tests {
		expr := readStr(tt.input)

		// Extract string from list if needed
		var str *Expr
		if expr.Type == Pair {
			str = listToSlice(expr)[0]
		} else {
			str = expr
		}

		if str.Type != String {
			t.Errorf("readStr(%q) type = %v, want String", tt.input, str.Type)
		}
		if str.Str != tt.want {
			t.Errorf("readStr(%q) = %q, want %q", tt.input, str.Str, tt.want)
		}
	}
}
