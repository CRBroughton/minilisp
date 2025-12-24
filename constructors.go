package main

type ExprType int

const (
	Nil ExprType = iota
	Bool
	Number
	Symbol
	Cons
	Builtin
	Lambda
	Macro
)

type Expr struct {
	Type   ExprType
	Num    int
	Sym    string
	Car    *Expr // head | first
	Cdr    *Expr // tail | last
	Fn     func([]*Expr) *Expr
	Params *Expr
	Body   *Expr
	Env    *Env
}

var nilExpr = &Expr{Type: Nil}
var trueExpr = &Expr{Type: Bool}
var falseExpr = &Expr{Type: Bool}

// Basic construtors for the various types
func makeNum(n int) *Expr {
	return &Expr{Type: Number, Num: n}
}

func makeSym(s string) *Expr {
	if s == "nil" {
		return nilExpr
	}
	if s == "true" {
		return trueExpr
	}
	if s == "false" {
		return falseExpr
	}
	return &Expr{Type: Symbol, Sym: s}
}

func cons(car, cdr *Expr) *Expr {
	return &Expr{Type: Cons, Car: car, Cdr: cdr}
}

func makeBuiltin(fn func([]*Expr) *Expr) *Expr {
	return &Expr{Type: Builtin, Fn: fn}
}

// some helpers Ill need for lists

func list(exprs ...*Expr) *Expr {
	result := nilExpr

	for i := len(exprs) - 1; i >= 0; i-- {
		result = cons(exprs[i], result)
	}
	return result
}

func listToSlice(e *Expr) []*Expr {
	var result []*Expr
	for e != nilExpr && e.Type == Cons {
		result = append(result, e.Car)
		e = e.Cdr
	}
	return result
}
