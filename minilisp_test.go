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

func TestCons(t *testing.T) {
	car := makeNum(1)
	cdr := makeNum(2)
	pair := cons(car, cdr)

	if pair.Type != Cons {
		t.Errorf("cons type = %v, want Cons", pair.Type)
	}
	if pair.Car != car {
		t.Errorf("cons car = %v, want %v", pair.Car, car)
	}
	if pair.Cdr != cdr {
		t.Errorf("cons cdr = %v, want %v", pair.Cdr, cdr)
	}
}

func TestList(t *testing.T) {
	lst := list(makeNum(1), makeNum(2), makeNum(3))

	if lst.Type != Cons {
		t.Fatalf("list type = %v, want Cons", lst.Type)
	}
	if lst.Car.Type != Number || lst.Car.Num != 1 {
		t.Errorf("first element = %v, want 1", lst.Car)
	}
	if lst.Cdr.Car.Num != 2 {
		t.Errorf("second element = %v, want 2", lst.Cdr.Car)
	}
	if lst.Cdr.Cdr.Car.Num != 3 {
		t.Errorf("third element = %v, want 3", lst.Cdr.Cdr.Car)
	}
	if lst.Cdr.Cdr.Cdr != nilExpr {
		t.Errorf("list should terminate with nil")
	}
}

func TestListToSlice(t *testing.T) {
	lst := list(makeNum(1), makeNum(2), makeNum(3))
	slice := listToSlice(lst)

	if len(slice) != 3 {
		t.Fatalf("listToSlice length = %d, want 3", len(slice))
	}
	if slice[0].Num != 1 || slice[1].Num != 2 || slice[2].Num != 3 {
		t.Errorf("listToSlice values incorrect")
	}
}

func TestListToSliceEmpty(t *testing.T) {
	slice := listToSlice(nilExpr)
	if len(slice) != 0 {
		t.Errorf("listToSlice(nil) length = %d, want 0", len(slice))
	}
}
